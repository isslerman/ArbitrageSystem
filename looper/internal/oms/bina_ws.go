package oms

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type BinaWSConsumer struct {
	conn      *websocket.Conn
	pub       string
	priv      ed25519.PrivateKey
	connMutex sync.Mutex
	wg        sync.WaitGroup
}

func NewBinaWSConsumer(pub, priv string) (*BinaWSConsumer, error) {
	privER := loadEdKeys(priv)

	// Conectar ao endpoint WebSocket da Binance (Spot)
	wsURL := "wss://ws-api.binance.com:443/ws-api/v3"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Falha ao conectar ao WebSocket da Binance: %v", err)
	}
	defer conn.Close()
	log.Println("Conexão WebSocket estabelecida com sucesso.")

	return &BinaWSConsumer{
		pub:  pub,
		priv: privER,
	}, nil
}

func (w *BinaWSConsumer) Start(ctx context.Context) (<-chan *ExecutionReport, error) {
	out := make(chan *ExecutionReport, 64)

	// Start connection manager goroutine
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer close(out)

		backoff := time.Second
		maxBackoff := time.Minute

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := w.establishConnection(); err != nil {
					log.Printf("Connection error: %v", err)
					time.Sleep(backoff)
					backoff = min(backoff*2, maxBackoff)
					continue
				}

				// Reset backoff on successful connection
				backoff = time.Second

				// Start message receiver
				if err := w.receiveMessages(ctx, out); err != nil {
					log.Printf("Message receiver error: %v", err)
					w.closeConnection()
				}
			}
		}
	}()

	return out, nil
}

func (w *BinaWSConsumer) establishConnection() error {
	w.connMutex.Lock()
	defer w.connMutex.Unlock()

	// Close existing connection if any
	if w.conn != nil {
		w.conn.Close()
	}
	wsURL := "wss://ws-api.binance.com:443/ws-api/v3"
	// Create new connection
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	w.conn = c

	// Perform logon
	if err := w.logon(); err != nil {
		c.Close()
		return fmt.Errorf("logon failed: %w", err)
	}

	// Subscribe to user data stream
	if err := w.subscribe(); err != nil {
		c.Close()
		return fmt.Errorf("subscribe failed: %w", err)
	}

	log.Println("Successfully connected and authenticated")
	return nil
}

// ORDER CREATED:
// executionReport: &{EventType:executionReport
// EventTime:1754157728242
// Symbol:USDTBRL
// ClientOrderID:web_4776588d7c004ffc9b3521e023bfe856
// Side:BUY OrderType:LIMIT
// TimeInForce:GTC
// Quantity:5.00000000
// Price:5.57100000
// StopPrice:0.00000000
// IcebergQuantity:0.00000000
// OrderListID:-1 OrigClientOrderID:
// ExecutionType:NEW
// OrderStatus:NEW
// OrderRejectReason:NONE OrderID:215764785 LastExecutedQuantity:0.00000000 CumulativeFilledQty:0.00000000 LastExecutedPrice:0.00000000 CommissionAmount:0 CommissionAsset:<nil> TransactionTime:1754157728242 TradeID:-1 PreventedMatchID:<nil> ExecutionID:504306765 IsOrderOnBook:true IsMakerSide:false IgnoreField:false OrderCreationTime:1754157728242 CumulativeQuoteQty:0.00000000 LastQuoteQty:0.00000000 QuoteOrderQty:0.00000000 WorkingTime:1754157728242 SelfTradePreventionMode:EXPIRE_MAKER}

// ORDER CANCELLED
// executionReport: &{EventType:executionReport
// EventTime:1754157757036
// Symbol:USDTBRL
// ClientOrderID:web_b0493b70524b4de2ba32942eb421554a
// Side:BUY
// OrderType:LIMIT
// TimeInForce:GTC
// Quantity:5.00000000
// Price:5.57100000
// StopPrice:0.00000000
// IcebergQuantity:0.00000000
// OrderListID:-1
// OrigClientOrderID:web_4776588d7c004ffc9b3521e023bfe856
// ExecutionType:CANCELED
// OrderStatus:CANCELED
// OrderRejectReason:NONE
// OrderID:215764785
// LastExecutedQuantity:0.00000000
// CumulativeFilledQty:0.00000000
// LastExecutedPrice:0.00000000
// CommissionAmount:0
// CommissionAsset:<nil>
// TransactionTime:1754157757036
// TradeID:-1 PreventedMatchID:<nil>
// ExecutionID:504306861
// IsOrderOnBook:false
// IsMakerSide:false
// IgnoreField:false
// OrderCreationTime:1754157728242
// CumulativeQuoteQty:0.00000000
// LastQuoteQty:0.00000000
// QuoteOrderQty:0.00000000
// WorkingTime:1754157728242
// SelfTradePreventionMode:EXPIRE_MAKER
// }

func (w *BinaWSConsumer) receiveMessages(ctx context.Context, out chan<- *ExecutionReport) error {

	// opcional: renova deadline a cada pong
	const wsTimeout = 5 * time.Minute
	w.conn.SetReadDeadline(time.Now().Add(wsTimeout))
	w.conn.SetPongHandler(func(string) error {
		w.conn.SetReadDeadline(time.Now().Add(wsTimeout))
		return nil
	})

	for {
		// leitura bruta
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			return err
		}

		er, ok, err := parseExecReport(message)
		if err != nil {
			log.Printf("parseExecReport: %v [msg]: %s", err, string(message))
			continue
		}
		if ok { // is executionReport
			select {
			case out <- er: // envia struct já decodificado
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// cancelar externo
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (w *BinaWSConsumer) logon() error {
	// Preparar mensagem de login (session.logon) com apiKey, timestamp atual e assinatura
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // timestamp em ms
	// Construir payload de assinatura: apiKey=<APIKEY>&timestamp=<TIMESTAMP>
	payload := fmt.Sprintf("apiKey=%s&timestamp=%d", w.pub, timestamp)
	// Assinar o payload com chave privada Ed25519 e codificar resultado em Base64
	signatureBytes := ed25519.Sign(w.priv, []byte(payload))
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	// Montar estrutura JSON para request de logon
	loginReq := map[string]interface{}{
		"id":     1,
		"method": "session.logon",
		"params": map[string]interface{}{
			"apiKey":    w.pub,
			"timestamp": timestamp,
			"signature": signatureBase64,
		},
	}
	// Enviar requisição de logon
	if err := w.conn.WriteJSON(loginReq); err != nil {
		log.Fatalf("Erro ao enviar session.logon: %v", err)
	}
	log.Println("Autenticando no WebSocket...")

	// Ler resposta do login
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		log.Fatalf("Erro ao ler resposta de login: %v", err)
	}
	// Verificar se login teve status 200
	var resp map[string]interface{}
	if err := json.Unmarshal(message, &resp); err != nil {
		log.Fatalf("Erro ao decodificar resposta de login: %v", err)
	}
	if status, ok := resp["status"].(float64); !ok || int(status) != 200 {
		log.Fatalf("Falha na autenticação WebSocket: %s", message)
	}
	log.Println("Autenticação bem-sucedida no WebSocket.")

	return nil
}

func (w *BinaWSConsumer) subscribe() error {
	// Assinar o User Data Stream da conta
	subscribeReq := map[string]interface{}{
		"id":     2,
		"method": "userDataStream.subscribe",
	}
	if err := w.conn.WriteJSON(subscribeReq); err != nil {
		log.Fatalf("Erro ao enviar userDataStream.subscribe: %v", err)
	}
	log.Println("Enviada requisição de inscrição no User Data Stream...")

	// Ler resposta da inscrição
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		log.Fatalf("Erro ao ler resposta de subscrição: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(message, &resp); err != nil {
		log.Fatalf("Erro ao decodificar resposta de subscrição: %v", err)
	}
	if status, ok := resp["status"].(float64); !ok || int(status) != 200 {
		log.Fatalf("Falha ao subscrever User Data Stream: %s", message)
	}
	log.Println("Subscrição ao User Data Stream realizada com sucesso. Aguardando eventos...")

	return nil
}

func (w *BinaWSConsumer) closeConnection() {
	w.connMutex.Lock()
	defer w.connMutex.Unlock()

	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
}

func (w *BinaWSConsumer) Wait() {
	w.wg.Wait()
}

// Helper functions remain the same as previous version
func loadEdKeys(derB64 string) (priv ed25519.PrivateKey) {
	switch len(derB64) % 4 {
	case 2:
		derB64 += "=="
	case 3:
		derB64 += "="
	}
	// Convert a BASE64 string to bytes - [00 00 00 00 ... 00 00]
	der, err := base64.StdEncoding.DecodeString(derB64)
	if err != nil {
		log.Fatalf("Base64 inválido: %v", err)
	}
	if len(der) != 48 { // 16 cabeçalho + 32 seed = 48 bytes
		log.Fatalf("DER inesperado: %d bytes (esperado 48)", len(der))
	}

	// parse PKCS#8
	privKeyAny, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		log.Fatalf("ParsePKCS8: %v", err)
	}

	privKey, ok := privKeyAny.(ed25519.PrivateKey)
	if !ok {
		log.Fatal("Chave privada fornecida não é do tipo Ed25519")
	}
	return privKey
}

// func min(a, b time.Duration) time.Duration {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func mustEnv(k string) string {
// 	v := os.Getenv(k)
// 	if v == "" {
// 		log.Fatalf("env %s not set", k)
// 	}
// 	return v
// }

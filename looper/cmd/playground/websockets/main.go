package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	apiKey := "wpCbnuCGOYuKsS24LA5en4k0GVDahHxxW3cJAld3ldkbjwLwhWEfdsni7ssfmDMK"
	derB64 := "MC4CAQAwBQYDK2VwBCIEINzxjOfc9ged4OA6HgOIM0gNVOivlf900toUeSToVpmg"

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

	// Conectar ao endpoint WebSocket da Binance (Spot)
	wsURL := "wss://ws-api.binance.com:443/ws-api/v3"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Falha ao conectar ao WebSocket da Binance: %v", err)
	}
	defer conn.Close()
	log.Println("Conexão WebSocket estabelecida com sucesso.")

	// Preparar mensagem de login (session.logon) com apiKey, timestamp atual e assinatura
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // timestamp em ms
	// Construir payload de assinatura: apiKey=<APIKEY>&timestamp=<TIMESTAMP>
	payload := fmt.Sprintf("apiKey=%s&timestamp=%d", apiKey, timestamp)
	// Assinar o payload com chave privada Ed25519 e codificar resultado em Base64
	signatureBytes := ed25519.Sign(privKey, []byte(payload))
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	// Montar estrutura JSON para request de logon
	loginReq := map[string]interface{}{
		"id":     1,
		"method": "session.logon",
		"params": map[string]interface{}{
			"apiKey":    apiKey,
			"timestamp": timestamp,
			"signature": signatureBase64,
		},
	}
	// Enviar requisição de logon
	if err := conn.WriteJSON(loginReq); err != nil {
		log.Fatalf("Erro ao enviar session.logon: %v", err)
	}
	log.Println("Autenticando no WebSocket...")

	// Ler resposta do login
	_, message, err := conn.ReadMessage()
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

	// Assinar o User Data Stream da conta
	subscribeReq := map[string]interface{}{
		"id":     2,
		"method": "userDataStream.subscribe",
	}
	if err := conn.WriteJSON(subscribeReq); err != nil {
		log.Fatalf("Erro ao enviar userDataStream.subscribe: %v", err)
	}
	log.Println("Enviada requisição de inscrição no User Data Stream...")

	// Ler resposta da inscrição
	_, message, err = conn.ReadMessage()
	if err != nil {
		log.Fatalf("Erro ao ler resposta de subscrição: %v", err)
	}
	if err := json.Unmarshal(message, &resp); err != nil {
		log.Fatalf("Erro ao decodificar resposta de subscrição: %v", err)
	}
	if status, ok := resp["status"].(float64); !ok || int(status) != 200 {
		log.Fatalf("Falha ao subscrever User Data Stream: %s", message)
	}
	log.Println("Subscrição ao User Data Stream realizada com sucesso. Aguardando eventos...")

	// Loop para ler e processar eventos de ordem executada
	for {
		_, message, err = conn.ReadMessage()
		if err != nil {
			log.Printf("Conexão encerrada ou erro de leitura: %v", err)
			break
		}
		// Decodificar a mensagem JSON recebida
		var evt map[string]interface{}
		if err := json.Unmarshal(message, &evt); err != nil {
			log.Printf("Mensagem não-JSON recebida: %s", string(message))
			continue
		}
		// Verificar se é um evento de User Data Stream (possui campo "event")
		eventData, ok := evt["event"]
		if !ok {
			// Não é um evento (pode ser resposta a outras requisições, ignorar)
			continue
		}
		// Checar tipo de evento "e"
		if eventMap, ok := eventData.(map[string]interface{}); ok {
			if evtType, ok := eventMap["e"].(string); ok && evtType == "executionReport" {
				// Imprimir evento de ordem executada (executionReport)
				fmt.Printf("Evento de ordem: %s\n", string(message))
			}
		}
	}
}

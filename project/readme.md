## How to install
1. gh repo clone isslerman/ArbitrageSystem
2. Get Go installed -> https://go.dev/doc/install
3. Gcc req in Ubuntu -> apt-get install build-essential
4. Linux build all pods: 
```make -f makefile.nix build-all```
1. Linux build all pods: 
```bash make -f makefile.nix build-server```

## ToDo:

- Add Exchanges: 
1. https://bitcointoyou.com/#servicos
2. https://bitnuvem.com - não tem SOL
3. https://bitypreco.com
4. ver lista aqui -> https://cointradermonitor.com/arbitragem
5. + algumas -> [BingX], [Alpha Vantage],[Tiingo]

 [ ]- Bus rule filter for exchanges: if price, volume == 0, don't send msg.   
- Ver Erro:

--------
    panic: runtime error: invalid memory address or nil pointer dereference
    [signal SIGSEGV: segmentation violation code=0x2 addr=0x18 pc=0x102ae6454]

    goroutine 6 [running]:
    pods/internal/pod.(*Pod).sendBestOrderViaGRPC(0x140002dff88)
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/internal/pod/pod.go:73 +0x44
    pods/internal/pod.(*Pod).Run(0x140002dff88)
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/internal/pod/pod.go:52 +0x9c
    pods/infra/api.(*Server).runPod(0x14000188b38, 0x0?)
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/infra/api/handlers.go:53 +0xe8
    created by pods/infra/api.(*Server).StartPod in goroutine 4
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/infra/api/handlers.go:18 +0xc8

--------

    panic: runtime error: invalid memory address or nil pointer dereference
    [signal SIGSEGV: segmentation violation code=0x2 addr=0x10 pc=0x10285a3d8]

    goroutine 22 [running]:
    pods/internal/pod.(*Pod).Run(0x14000213f88)
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/internal/pod/pod.go:53 +0x98
    pods/infra/api.(*Server).runPod(0x14000188b38, 0x0?)
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/infra/api/handlers.go:53 +0xc8
    created by pods/infra/api.(*Server).StartPod in goroutine 20
        /Users/marcosissler/projects/202404-ArbitrageSystem/pods/infra/api/handlers.go:18 +0xc8

--------

- Others CEX to explore: 
  - NovaDax (https://www.novadax.com.br/product/orderbook?pair=SOL_BRL)
  - Coinnext (https://trade.coinext.com.br/)
- Poder ver uma amostra de dados do Pod para verificação de dados. 
- How to know if any pod stopped
- How to know if any pod is giving any error?
- Create a dockerfile with metabase ? stay at our notebook or let's go to our server? 
- How much is the balance in each wallet?
- How much we have earned each day?
- Add zap logging to pods. Each pod has a sample data to a specific logfile /var/log/pod.log - possible with rotate file. 
- Possible new tokens: Fetch.AI, Filecoin (Foxbit)
- IMPORTANT! Tecnica para pegar nosso melhor preço de compra na Binance e deixar uma ordem de venda na exchange com 0,4% de profit e sempre ajustar conforme o preço alterar. Reajuste por evento e/ou tempo. Qts segundos ou movimentação do preço. Caso seja efetuada a venda, o evento dispara uma compra na Binance ( market price? qual melhor estratégia? )


--------
PODS - PORTS USED
--------
15000 - BINA
15001 - BITP
15002 - FOXB
15003 - MBTC
15004 - RIPI


--------
ERROR: [bina order error]: <APIError> code=-1013, msg=Filter failure: LOT_SIZE - &{0.1943 SOLBRL 872.6 buy limit}
RULE: API URL -> {{url}}/api/v3/exchangeInfo?symbol=SOLBRL        
{
    "filterType": "LOT_SIZE",
    "minQty": "0.00100000",
    "maxQty": "9222449.00000000",
    "stepSize": "0.00100000"
},
FIX: amount need to have 3 decimals.

----- 
ERROR sem fundos na binance
2024/06/11 12:47:31 ERROR Error: !BADKEY="Funds insufficient"
-----
ERRO sem internet:
ERROR error sending ntfy msg status=429
2024/06/11 19:13:15 Error reading message: read tcp [2804:7f0:6401:a167:d0d5:84f2:5a2b:77e5]:62365->[2606:4700::6812:1fbd]:443: read: operation timed out
-----
tabela de erros: 
1771 - Error reading message: websocket: close 1006 (abnormal closure): unexpected EOF - 2024-06-12 14:30:42.903327-03
-----
Log file: 
2024/06/12 15:52:20 ERROR Error: LoggerInfoRepo: can't insert Info Log -  !BADKEY="pq: value too long for type character varying(255)"
-----
Error for Binance buying less than minimum: 
2024/06/12 16:41:05 Error handling message: <APIError> code=-1013, msg=Invalid quantity., [123 34 98 111 100 121 34 58 123 34 97 109 111 117 110 116 34 58 48 46 54 53 44 34 97 118 101 114 97 103 101 95 101 120 101 99 117 116 105 111 110 95 112 114 105 99 101 34 58 56 52 53 44 34 105 100 34 58 34 69 69 57 51 53 65 52 66 45 54 69 67 66 45 52 56 66 50 45 66 66 68 66 45 57 49 55 51 55 67 66 68 65 55 49 56 34 44 34 99 114 101 97 116 101 100 95 97 116 34 58 34 50 48 50 52 45 48 54 45 49 50 84 49 57 58 52 49 58 48 50 46 50 50 48 90 34 44 34 101 120 101 99 117 116 101 100 95 97 109 111 117 110 116 34 58 48 46 48 48 48 52 44 34 101 120 116 101 114 110 97 108 95 105 100 34 58 110 117 108 108 44 34 112 97 105 114 34 58 34 83 79 76 95 66 82 76 34 44 34 112 114 105 99 101 34 58 56 52 52 46 52 56 44 34 114 101 109 97 105 110 105 110 103 95 97 109 111 117 110 116 34 58 48 46 54 52 57 53 54 56 50 50 44 34 115 105 100 101 34 58 34 115 101 108 108 34 44 34 115 116 97 116 117 115 34 58 34 101 120 101 99 117 116 101 100 95 112 97 114 116 105 97 108 108 121 34 44 34 116 121 112 101 34 58 34 108 105 109 105 116 34 44 34 117 112 100 97 116 101 100 95 97 116 34 58 34 50 48 50 52 45 48 54 45 49 50 84 49 57 58 52 49 58 48 52 46 48 49 48 90 34 44 34 117 115 101 114 95 105 100 34 58 34 57 69 66 56 51 55 65 53 45 54 65 56 49 45 52 65 68 68 45 56 66 57 55 45 49 54 54 55 50 66 49 51 50 54 55 65 34 125 44 34 116 105 109 101 115 116 97 109 112 34 58 49 55 49 56 50 50 49 50 54 52 55 53 57 44 34 116 111 112 105 99 34 58 34 111 114 100 101 114 95 115 116 97 116 117 115 34 44 34 105 100 34 58 49 48 54 54 53 53 51 51 51 57 125], {1066553339 order_status 1718221264759 {0.65 845 EE935A4B-6ECB-48B2-BBDB-91737CBDA718 2024-06-12 19:41:02.22 +0000 UTC 0.0004 <nil> SOL_BRL 844.48 0.64956822 sell executed_partially limit 2024-06-12 19:41:04.01 +0000 UTC 9EB837A5-6A81-4ADD-8B97-16672B13267A}}
make: *** [start-server] Error 1
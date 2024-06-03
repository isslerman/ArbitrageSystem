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

- Bus rule filter for exchanges: if price, volume == 0, don't send msg.   
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
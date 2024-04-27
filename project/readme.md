ToDo:

- Add Exchanges: 
1. https://bitcointoyou.com/#servicos
2. https://bitnuvem.com - não tem SOL
3. https://bitypreco.com
4. ver lista aqui -> https://cointradermonitor.com/arbitragem
5. + algumas -> [BingX], [Alpha Vantage],[Tiingo]

- Bus rule filter for exchanges: if price, volume == 0, don't send msg.   
- Ver Erro:

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

- Poder ver uma amostra de dados do Pod para verificação de dados. 
- 
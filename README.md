# ArbitrageSystem

Licença de Uso Restrito

Copyright (c) 2024 Marcos Carneiro Issler

Todos os direitos reservados. O uso, cópia, modificação ou distribuição deste software, seja na íntegra ou em parte, é estritamente proibido sem a autorização expressa do autor.


your-project/
├── cmd/
│   └── engine/                # Entrypoint binary for the engine
│       └── main.go
│
├── config/
│   └── config.go              # Loads envs, files, CLI flags etc.
│
├── internal/
│   ├── engine/                # Main app coordinator (was App)
│   │   ├── engine.go
│   │   └── start.go
│   │
│   ├── marketdata/            # gRPC receiver, order book, trades
│   │   ├── receiver.go
│   │   ├── models.go
│   │   └── orderbook.go
│   │
│   ├── strategy/              # Strategy logic
│   │   ├── interface.go       # Strategy interface
│   │   └── makerwick/         # Specific strategy module
│   │       └── strategy.go
│   │
│   ├── oms/                   # Order Management System (exchange API)
│   │   └── client.go
│   │
│   ├── riskmodule/            # Position sizing, exposure limits
│   │   └── checker.go
│   │
│   ├── monitor/               # Logging, alerting, metrics
│   │   └── logger.go
│   │
│   └── repository/            # DB layer (if needed)
│       └── store.go
│
├── proto/                     # gRPC Protobuf files
│   └── order.proto
│
├── scripts/                   # Startup scripts, Dockerfiles, etc.
│
├── go.mod
└── README.md


---
The OMS is responsible for:

- Placing orders
- Tracking orders
- Handling order status (open, filled, canceled)
- Reacting to order lifecycle events (e.g. partially filled → rebalance)

🔄 Flow Overview
1. Engine starts → initializes Strategy, MarketData, OMS, etc.
2. OMS starts → spawns WebSocket listener.
3. WebSocket handler listens to Binance executionReport events:
   1. Detects fill/cancel
   2. Updates internal order state
   3. Notifies the strategy (via channel, callback, or internal flag)
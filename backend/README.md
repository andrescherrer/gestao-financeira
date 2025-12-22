# Backend - Sistema de Gestão Financeira

Backend desenvolvido em **Go** com **Fiber**, seguindo **Domain-Driven Design (DDD)**.

## 📁 Estrutura de Pastas

```
backend/
├── cmd/
│   └── api/                    # Ponto de entrada da aplicação
│       └── main.go
│
├── internal/                   # Código interno (DDD)
│   ├── shared/                 # Shared Kernel
│   │   ├── domain/
│   │   │   ├── valueobjects/   # Money, Currency, AccountContext, etc.
│   │   │   └── events/         # Domain Events base
│   │   └── infrastructure/
│   │       └── eventbus/        # Event Bus
│   │
│   ├── identity/               # Identity Context
│   │   ├── domain/             # Entidades, Value Objects, Services
│   │   ├── application/        # Use Cases, DTOs
│   │   ├── infrastructure/      # Repositórios, Serviços externos
│   │   └── presentation/       # Handlers HTTP, DTOs
│   │
│   ├── account/                # Account Management Context
│   ├── transaction/            # Transaction Context (Core Domain)
│   ├── category/                # Category Context
│   ├── budget/                 # Budget Context
│   ├── reporting/              # Reporting Context
│   ├── investment/             # Investment Context
│   ├── goal/                   # Goal Context
│   └── notification/           # Notification Context
│
└── pkg/                        # Pacotes compartilhados
    ├── database/               # Configuração do banco
    ├── logger/                 # Logger
    ├── validator/              # Validação
    └── middleware/             # Middlewares HTTP
```

## 🚀 Como Executar

```bash
# Instalar dependências
go mod download

# Executar
go run cmd/api/main.go
```

## 📚 Documentação

Veja o [PLANEJAMENTO_GO.md](../../docs/planejamento/PLANEJAMENTO_GO.md) para mais detalhes sobre a arquitetura.


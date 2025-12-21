# Sistema de Gestão Financeira

Sistema de gestão financeira pessoal e profissional desenvolvido em **Go** (backend) e **Next.js** (frontend) seguindo **Domain-Driven Design (DDD)**.

## 🚀 Stack Tecnológico

### Backend
- **Go 1.21+** com **Fiber** (framework web)
- **PostgreSQL** (banco de dados)
- **Redis** (cache e rate limiting)
- **GORM** (ORM)
- **OpenTelemetry** (observabilidade)
- **Prometheus + Grafana** (monitoramento)

### Frontend
- **Next.js 14+** com **TypeScript**
- **shadcn/ui** (componentes UI)
- **Tailwind CSS** (styling)
- **TanStack Query** (server state)
- **React Hook Form + Zod** (formulários)

## 📁 Estrutura do Projeto

```
gestao-financeira/
├── backend/              # Backend Go
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── migrations/
├── frontend/             # Frontend Next.js
│   ├── app/
│   ├── components/
│   └── lib/
├── docs/                 # Documentação
└── scripts/              # Scripts utilitários
```

## 🛠️ Desenvolvimento

### Pré-requisitos
- Go 1.21+
- Node.js 20+
- Docker e Docker Compose
- PostgreSQL 15+
- Redis 7+

### Iniciando o Projeto

1. Clone o repositório
2. Configure as variáveis de ambiente (veja `.env.example`)
3. Execute `docker-compose up` para subir os serviços
4. Execute as migrations
5. Inicie o backend e frontend

## 📚 Documentação

- [Planejamento Completo](./PLANEJAMENTO_GO.md)
- [Tarefas do Projeto](./TAREFAS.md)
- [Explicação do Planejamento](./EXPLICACAO_GO.md)

## 📝 Licença

Este projeto é privado.


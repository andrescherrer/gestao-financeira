# Sistema de Gestão Financeira

Sistema de gestão financeira pessoal e profissional desenvolvido em **Go** (backend) e **Vue 3** (frontend) seguindo **Domain-Driven Design (DDD)**.

## 🚀 Stack Tecnológico

### Backend
- **Go 1.21+** com **Fiber** (framework web)
- **PostgreSQL** (banco de dados)
- **Redis** (cache e rate limiting)
- **GORM** (ORM)
- **OpenTelemetry** (observabilidade)
- **Prometheus + Grafana** (monitoramento)

### Frontend
- **Vue 3** com **TypeScript**
- **Vite** (build tool)
- **shadcn-vue** (componentes UI baseados em reka-ui)
- **Tailwind CSS** (styling)
- **TanStack Query** (server state)
- **Vee-Validate + Zod** (formulários)
- **Pinia** (state management)

## 📁 Estrutura do Projeto

```
gestao-financeira/
├── backend/              # Backend Go
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── migrations/
├── frontend/             # Frontend Vue 3
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── api/
│   │   └── utils/
│   ├── public/
│   └── cypress/
├── docs/                 # Documentação
├── deploy/               # Scripts de deploy
└── monitoring/           # Configurações de monitoramento
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

- [Planejamento Completo](./docs/planejamento/PLANEJAMENTO_GO.md)
- [Tarefas do Projeto](./TAREFAS.md)
- [Guia de Deploy](./docs/DEPLOY.md)
- [Configuração da API](./backend/CONFIG.md)

## 📝 Licença

Este projeto é privado.


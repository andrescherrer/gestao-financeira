# Planejamento de Tarefas - Sistema de Gestão Financeira

Este documento detalha as tarefas práticas para implementação do sistema, organizadas por fases e prioridades.

## 📊 Status Geral do Projeto

**Última verificação:** 2025-01-27

### ✅ Concluído
- **Setup Inicial** (SETUP-001 a SETUP-006): ✅ Completo
- **Sprint 1.1: Setup Backend** (BE-001 a BE-008): ✅ Completo
- **Sprint 1.2: Shared Kernel** (SK-001 a SK-006): ✅ Completo
- **Sprint 1.3: Identity Context** (ID-001 a ID-013): ✅ Completo

### ⏳ Em Progresso / Pendente
- **Sprint 1.4: Account Management** (AC-001 a AC-011): ⏳ Estrutura criada, implementação pendente
- **Sprint 1.5: Transaction Context** (TX-001 a TX-015): ⏳ Estrutura criada, implementação pendente
- **Sprint 1.6: Swagger** (DOC-001 a DOC-006): ⏳ Não iniciado
- **Sprint 1.7: Setup Frontend** (FE-001 a FE-009): ⏳ Estrutura criada, implementação pendente
- **Demais sprints**: ⏳ Não iniciadas

### 📈 Progresso
- **Fase 1 (Fundação e MVP)**: ~30% concluído
  - Backend base: ✅ 100%
  - Identity Context: ✅ 100%
  - Account Management: ⏳ 0%
  - Transaction Context: ⏳ 0%
  - Frontend: ⏳ 0%

---

## 📋 Legenda

- **Prioridade**: 🔴 Alta | 🟡 Média | 🟢 Baixa
- **Tipo**: 🔵 Backend | 🟣 Frontend | 🟠 DevOps | 🟤 Testes | ⚪ Documentação
- **Esforço**: Estimativa em horas (1h, 2h, 4h, 8h, 16h)
- **Status**: ⏳ Pendente | 🚧 Em Progresso | ✅ Concluído | ❌ Bloqueado

---

## 🎯 Pré-requisitos e Setup Inicial

### Tarefas de Configuração do Ambiente

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| SETUP-001 | Configurar repositório Git e estrutura inicial | ⚪ | 🔴 | 2h | - | ✅ |
| SETUP-002 | Configurar Docker e docker-compose para desenvolvimento | 🟠 | 🔴 | 4h | SETUP-001 | ✅ |
| SETUP-003 | Configurar PostgreSQL no Docker | 🟠 | 🔴 | 2h | SETUP-002 | ✅ |
| SETUP-004 | Configurar Redis no Docker | 🟠 | 🔴 | 2h | SETUP-002 | ✅ |
| SETUP-005 | Configurar variáveis de ambiente (.env.example) | ⚪ | 🔴 | 1h | SETUP-002 | ✅ |
| SETUP-006 | Configurar CI/CD básico (GitHub Actions) | 🟠 | 🟡 | 4h | SETUP-001 | ✅ |

---

## 📦 Fase 1: Fundação e MVP (3-4 semanas)

### Sprint 1.1: Setup Backend (Semana 1)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| BE-001 | Criar estrutura de pastas Go (cmd, internal, pkg) | 🔵 | 🔴 | 1h | SETUP-001 | ✅ |
| BE-002 | Inicializar go.mod e dependências básicas (Fiber, GORM) | 🔵 | 🔴 | 2h | BE-001 | ✅ |
| BE-003 | Configurar Fiber com middlewares básicos (logger, recover, CORS) | 🔵 | 🔴 | 4h | BE-002 | ✅ |
| BE-004 | Configurar conexão com PostgreSQL (GORM) | 🔵 | 🔴 | 4h | SETUP-003, BE-002 | ✅ |
| BE-005 | Implementar health check endpoint (/health) | 🔵 | 🔴 | 2h | BE-003 | ✅ |
| BE-006 | Configurar logger estruturado (zerolog) | 🔵 | 🟡 | 2h | BE-002 | ✅ |
| BE-007 | Criar Dockerfile multi-stage para backend | 🟠 | 🔴 | 2h | BE-002 | ✅ |
| BE-008 | Testar build e execução em Docker | 🟤 | 🔴 | 2h | BE-007, SETUP-002 | ✅ |

**Entregável Sprint 1.1:** API rodando em Docker com health check funcionando

---

### Sprint 1.2: Shared Kernel (Semana 1-2)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| SK-001 | Criar value object Money (amount, currency) | 🔵 | 🔴 | 4h | BE-001 | ✅ |
| SK-002 | Criar value object Currency (BRL, USD, EUR) | 🔵 | 🔴 | 2h | SK-001 | ✅ |
| SK-003 | Criar value object AccountContext (Personal, Business) | 🔵 | 🔴 | 2h | BE-001 | ✅ |
| SK-004 | Criar domain events base (DomainEvent interface) | 🔵 | 🔴 | 2h | BE-001 | ✅ |
| SK-005 | Implementar Event Bus simples | 🔵 | 🟡 | 4h | SK-004 | ✅ |
| SK-006 | Testes unitários para value objects | 🟤 | 🟡 | 4h | SK-001, SK-002, SK-003 | ✅ |

**Entregável Sprint 1.2:** Shared Kernel completo com testes

---

### Sprint 1.3: Identity Context - Backend (Semana 2)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| ID-001 | Criar value object Email com validação | 🔵 | 🔴 | 2h | BE-001 | ✅ |
| ID-002 | Criar value object PasswordHash (bcrypt) | 🔵 | 🔴 | 4h | BE-001 | ✅ |
| ID-003 | Criar value object UserName | 🔵 | 🔴 | 2h | BE-001 | ✅ |
| ID-004 | Criar entidade User (agregado raiz) | 🔵 | 🔴 | 4h | ID-001, ID-002, ID-003 | ✅ |
| ID-005 | Criar interface UserRepository | 🔵 | 🔴 | 1h | ID-004 | ✅ |
| ID-006 | Implementar GormUserRepository | 🔵 | 🔴 | 6h | ID-005, BE-004 | ✅ |
| ID-007 | Criar migration para tabela users | 🔵 | 🔴 | 2h | ID-006 | ✅ |
| ID-008 | Implementar RegisterUserUseCase | 🔵 | 🔴 | 4h | ID-004, ID-005 | ✅ |
| ID-009 | Implementar LoginUseCase com JWT | 🔵 | 🔴 | 6h | ID-004, ID-005 | ✅ |
| ID-010 | Criar AuthHandler (Register, Login) | 🔵 | 🔴 | 4h | ID-008, ID-009 | ✅ |
| ID-011 | Criar middleware de autenticação JWT | 🔵 | 🔴 | 4h | ID-009 | ✅ |
| ID-012 | Configurar rotas de autenticação (/api/v1/auth/*) | 🔵 | 🔴 | 2h | ID-010 | ✅ |
| ID-013 | Testes unitários para Identity Context | 🟤 | 🟡 | 8h | ID-004, ID-008, ID-009 | ✅ |

**Entregável Sprint 1.3:** API de autenticação funcionando (registro e login)

---

### Sprint 1.4: Account Management Context - Backend (Semana 2-3)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| AC-001 | Criar value object AccountID | 🔵 | 🔴 | 1h | BE-001 | ✅ |
| AC-002 | Criar entidade Account (agregado raiz) | 🔵 | 🔴 | 4h | SK-001, SK-003, AC-001 | ⏳ |
| AC-003 | Criar interface AccountRepository | 🔵 | 🔴 | 1h | AC-002 | ⏳ |
| AC-004 | Implementar GormAccountRepository | 🔵 | 🔴 | 6h | AC-003, BE-004 | ⏳ |
| AC-005 | Criar migration para tabela accounts | 🔵 | 🔴 | 2h | AC-004 | ⏳ |
| AC-006 | Implementar CreateAccountUseCase | 🔵 | 🔴 | 4h | AC-002, AC-003 | ⏳ |
| AC-007 | Implementar ListAccountsUseCase | 🔵 | 🔴 | 2h | AC-003 | ⏳ |
| AC-008 | Implementar GetAccountUseCase | 🔵 | 🔴 | 2h | AC-003 | ⏳ |
| AC-009 | Criar AccountHandler (CRUD) | 🔵 | 🔴 | 6h | AC-006, AC-007, AC-008 | ⏳ |
| AC-010 | Configurar rotas de accounts (/api/v1/accounts/*) | 🔵 | 🔴 | 2h | AC-009, ID-011 | ⏳ |
| AC-011 | Testes unitários para Account Context | 🟤 | 🟡 | 6h | AC-002, AC-006 | ⏳ |

**Entregável Sprint 1.4:** API de contas funcionando (CRUD completo)

---

### Sprint 1.5: Transaction Context - Backend (Semana 3)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| TX-001 | Criar value object TransactionID | 🔵 | 🔴 | 1h | BE-001 | ⏳ |
| TX-002 | Criar value object TransactionType (Income, Expense) | 🔵 | 🔴 | 2h | BE-001 | ⏳ |
| TX-003 | Criar value object TransactionDescription | 🔵 | 🔴 | 1h | BE-001 | ⏳ |
| TX-004 | Criar entidade Transaction (agregado raiz) | 🔵 | 🔴 | 6h | SK-001, TX-001, TX-002, TX-003 | ⏳ |
| TX-005 | Criar interface TransactionRepository | 🔵 | 🔴 | 1h | TX-004 | ⏳ |
| TX-006 | Implementar GormTransactionRepository | 🔵 | 🔴 | 6h | TX-005, BE-004 | ⏳ |
| TX-007 | Criar migration para tabela transactions | 🔵 | 🔴 | 2h | TX-006 | ⏳ |
| TX-008 | Implementar CreateTransactionUseCase | 🔵 | 🔴 | 4h | TX-004, TX-005, AC-003 | ⏳ |
| TX-009 | Implementar ListTransactionsUseCase | 🔵 | 🔴 | 4h | TX-005 | ⏳ |
| TX-010 | Implementar GetTransactionUseCase | 🔵 | 🔴 | 2h | TX-005 | ⏳ |
| TX-011 | Implementar UpdateTransactionUseCase | 🔵 | 🔴 | 4h | TX-004, TX-005 | ⏳ |
| TX-012 | Implementar DeleteTransactionUseCase | 🔵 | 🔴 | 2h | TX-004, TX-005 | ⏳ |
| TX-013 | Criar TransactionHandler (CRUD completo) | 🔵 | 🔴 | 6h | TX-008, TX-009, TX-010, TX-011, TX-012 | ⏳ |
| TX-014 | Configurar rotas de transactions (/api/v1/transactions/*) | 🔵 | 🔴 | 2h | TX-013, ID-011 | ⏳ |
| TX-015 | Testes unitários para Transaction Context | 🟤 | 🟡 | 8h | TX-004, TX-008 | ⏳ |

**Entregável Sprint 1.5:** API de transações funcionando (CRUD completo)

---

### Sprint 1.6: Swagger e Documentação Básica (Semana 3)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| DOC-001 | Instalar e configurar swaggo/swag | 🔵 | 🔴 | 2h | BE-002 | ⏳ |
| DOC-002 | Adicionar anotações Swagger nos handlers de Auth | 🔵 | 🔴 | 2h | ID-010, DOC-001 | ⏳ |
| DOC-003 | Adicionar anotações Swagger nos handlers de Account | 🔵 | 🔴 | 2h | AC-009, DOC-001 | ⏳ |
| DOC-004 | Adicionar anotações Swagger nos handlers de Transaction | 🔵 | 🔴 | 2h | TX-013, DOC-001 | ⏳ |
| DOC-005 | Configurar rota /swagger/* no Fiber | 🔵 | 🔴 | 1h | DOC-001 | ⏳ |
| DOC-006 | Gerar e testar documentação Swagger | 🔵 | 🔴 | 1h | DOC-002, DOC-003, DOC-004, DOC-005 | ⏳ |

**Entregável Sprint 1.6:** Swagger documentado e acessível

---

### Sprint 1.7: Setup Frontend (Semana 3-4)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-001 | Criar projeto Next.js 14 com TypeScript | 🟣 | 🔴 | 2h | - | ⏳ |
| FE-002 | Configurar Tailwind CSS | 🟣 | 🔴 | 2h | FE-001 | ⏳ |
| FE-003 | Instalar e configurar shadcn/ui | 🟣 | 🔴 | 4h | FE-002 | ⏳ |
| FE-004 | Instalar dependências (React Hook Form, Zod, Axios, TanStack Query) | 🟣 | 🔴 | 1h | FE-001 | ⏳ |
| FE-005 | Configurar estrutura de pastas (app, components, lib) | 🟣 | 🔴 | 2h | FE-001 | ⏳ |
| FE-006 | Criar layout base (Header, Sidebar, Footer) | 🟣 | 🔴 | 4h | FE-003 | ⏳ |
| FE-007 | Configurar cliente API (Axios) | 🟣 | 🔴 | 4h | FE-004 | ⏳ |
| FE-008 | Configurar variáveis de ambiente (.env.local) | 🟣 | 🔴 | 1h | FE-001 | ⏳ |
| FE-009 | Criar Dockerfile para frontend | 🟠 | 🟡 | 2h | FE-001 | ⏳ |

**Entregável Sprint 1.7:** Frontend configurado e rodando

---

### Sprint 1.8: Módulo de Autenticação - Frontend (Semana 4)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-AUTH-001 | Criar página de Login (/login) | 🟣 | 🔴 | 4h | FE-006, FE-007 | ⏳ |
| FE-AUTH-002 | Criar página de Registro (/register) | 🟣 | 🔴 | 4h | FE-006, FE-007 | ⏳ |
| FE-AUTH-003 | Criar hook useAuth para gerenciar autenticação | 🟣 | 🔴 | 4h | FE-007 | ⏳ |
| FE-AUTH-004 | Implementar proteção de rotas (middleware) | 🟣 | 🔴 | 4h | FE-AUTH-003 | ⏳ |
| FE-AUTH-005 | Criar componente de formulário de login (React Hook Form + Zod) | 🟣 | 🔴 | 4h | FE-003, FE-004 | ⏳ |
| FE-AUTH-006 | Criar componente de formulário de registro | 🟣 | 🔴 | 4h | FE-003, FE-004 | ⏳ |
| FE-AUTH-007 | Integrar com API de autenticação (login) | 🟣 | 🔴 | 2h | FE-AUTH-001, DOC-006 | ⏳ |
| FE-AUTH-008 | Integrar com API de autenticação (registro) | 🟣 | 🔴 | 2h | FE-AUTH-002, DOC-006 | ⏳ |
| FE-AUTH-009 | Implementar tratamento de erros e loading states | 🟣 | 🟡 | 2h | FE-AUTH-007, FE-AUTH-008 | ⏳ |
| FE-AUTH-010 | Testar fluxo completo de autenticação | 🟤 | 🔴 | 2h | FE-AUTH-007, FE-AUTH-008 | ⏳ |

**Entregável Sprint 1.8:** Autenticação funcionando no frontend

---

### Sprint 1.9: Módulo de Contas - Frontend (Semana 4)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-ACC-001 | Criar hook useAccounts (TanStack Query) | 🟣 | 🔴 | 2h | FE-007 | ⏳ |
| FE-ACC-002 | Criar página de lista de contas (/accounts) | 🟣 | 🔴 | 4h | FE-006, FE-ACC-001 | ⏳ |
| FE-ACC-003 | Criar componente AccountCard | 🟣 | 🔴 | 2h | FE-003 | ⏳ |
| FE-ACC-004 | Criar página de detalhes da conta (/accounts/[id]) | 🟣 | 🔴 | 4h | FE-006, FE-ACC-001 | ⏳ |
| FE-ACC-005 | Criar página de criação de conta (/accounts/new) | 🟣 | 🔴 | 4h | FE-006, FE-ACC-001 | ⏳ |
| FE-ACC-006 | Criar formulário de conta (React Hook Form + Zod) | 🟣 | 🔴 | 4h | FE-003, FE-004 | ⏳ |
| FE-ACC-007 | Integrar com API de contas (listar) | 🟣 | 🔴 | 2h | FE-ACC-002, DOC-006 | ⏳ |
| FE-ACC-008 | Integrar com API de contas (criar) | 🟣 | 🔴 | 2h | FE-ACC-005, DOC-006 | ⏳ |
| FE-ACC-009 | Integrar com API de contas (detalhes) | 🟣 | 🔴 | 2h | FE-ACC-004, DOC-006 | ⏳ |
| FE-ACC-010 | Implementar loading e error states | 🟣 | 🟡 | 2h | FE-ACC-007, FE-ACC-008 | ⏳ |

**Entregável Sprint 1.9:** Módulo de contas funcionando no frontend

---

### Sprint 1.10: Módulo de Transações - Frontend (Semana 4)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-TX-001 | Criar hook useTransactions (TanStack Query) | 🟣 | 🔴 | 2h | FE-007 | ⏳ |
| FE-TX-002 | Criar página de lista de transações (/transactions) | 🟣 | 🔴 | 4h | FE-006, FE-TX-001 | ⏳ |
| FE-TX-003 | Criar componente TransactionTable | 🟣 | 🔴 | 4h | FE-003 | ⏳ |
| FE-TX-004 | Criar página de detalhes da transação (/transactions/[id]) | 🟣 | 🔴 | 4h | FE-006, FE-TX-001 | ⏳ |
| FE-TX-005 | Criar página de criação de transação (/transactions/new) | 🟣 | 🔴 | 4h | FE-006, FE-TX-001 | ⏳ |
| FE-TX-006 | Criar formulário de transação (React Hook Form + Zod) | 🟣 | 🔴 | 6h | FE-003, FE-004 | ⏳ |
| FE-TX-007 | Integrar com API de transações (listar) | 🟣 | 🔴 | 2h | FE-TX-002, DOC-006 | ⏳ |
| FE-TX-008 | Integrar com API de transações (criar) | 🟣 | 🔴 | 2h | FE-TX-005, DOC-006 | ⏳ |
| FE-TX-009 | Integrar com API de transações (detalhes) | 🟣 | 🔴 | 2h | FE-TX-004, DOC-006 | ⏳ |
| FE-TX-010 | Implementar loading e error states | 🟣 | 🟡 | 2h | FE-TX-007, FE-TX-008 | ⏳ |

**Entregável Sprint 1.10:** Módulo de transações funcionando no frontend

---

## 🎯 Fase 2: Core Domain e Integrações (3-4 semanas)

### Sprint 2.1: Integração Transaction ↔ Account (Semana 5)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| INT-001 | Implementar atualização de saldo ao criar transação | 🔵 | 🔴 | 6h | TX-008, AC-003 | ⏳ |
| INT-002 | Implementar atualização de saldo ao atualizar transação | 🔵 | 🔴 | 6h | TX-011, AC-003 | ⏳ |
| INT-003 | Implementar atualização de saldo ao deletar transação | 🔵 | 🔴 | 4h | TX-012, AC-003 | ⏳ |
| INT-004 | Criar domain event TransactionCreated | 🔵 | 🔴 | 2h | SK-004, TX-004 | ⏳ |
| INT-005 | Criar handler para atualizar saldo via event bus | 🔵 | 🔴 | 4h | INT-004, SK-005 | ⏳ |
| INT-006 | Testes de integração Transaction ↔ Account | 🟤 | 🔴 | 4h | INT-001, INT-002, INT-003 | ⏳ |
| FE-INT-001 | Atualizar saldo em tempo real no frontend | 🟣 | 🟡 | 4h | FE-ACC-002, FE-TX-008 | ⏳ |

**Entregável Sprint 2.1:** Transações atualizam saldo das contas automaticamente

---

### Sprint 2.2: Event Bus e Domain Events (Semana 5)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| EVT-001 | Expandir Event Bus com retry e error handling | 🔵 | 🟡 | 4h | SK-005 | ⏳ |
| EVT-002 | Criar domain events para User (UserRegistered, etc.) | 🔵 | 🟡 | 2h | SK-004, ID-004 | ⏳ |
| EVT-003 | Criar domain events para Account (AccountCreated, etc.) | 🔵 | 🟡 | 2h | SK-004, AC-002 | ⏳ |
| EVT-004 | Implementar publicação automática de eventos nos use cases | 🔵 | 🟡 | 4h | EVT-002, EVT-003, INT-004 | ⏳ |
| EVT-005 | Criar event handlers para logging | 🔵 | 🟢 | 2h | EVT-001 | ⏳ |

**Entregável Sprint 2.2:** Sistema de eventos de domínio funcionando

---

### Sprint 2.3: Category Context - Backend (Semana 6)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| CAT-001 | Criar value object CategoryID | 🔵 | 🔴 | 1h | BE-001 | ⏳ |
| CAT-002 | Criar entidade Category (agregado raiz) | 🔵 | 🔴 | 4h | CAT-001 | ⏳ |
| CAT-003 | Criar interface CategoryRepository | 🔵 | 🔴 | 1h | CAT-002 | ⏳ |
| CAT-004 | Implementar GormCategoryRepository | 🔵 | 🔴 | 6h | CAT-003, BE-004 | ⏳ |
| CAT-005 | Criar migration para tabela categories | 🔵 | 🔴 | 2h | CAT-004 | ⏳ |
| CAT-006 | Implementar use cases de Category (CRUD) | 🔵 | 🔴 | 6h | CAT-002, CAT-003 | ⏳ |
| CAT-007 | Criar CategoryHandler (CRUD completo) | 🔵 | 🔴 | 4h | CAT-006 | ⏳ |
| CAT-008 | Configurar rotas de categories (/api/v1/categories/*) | 🔵 | 🔴 | 2h | CAT-007, ID-011 | ⏳ |
| CAT-009 | Adicionar anotações Swagger para Category | 🔵 | 🟡 | 2h | CAT-007, DOC-001 | ⏳ |
| CAT-010 | Testes unitários para Category Context | 🟤 | 🟡 | 6h | CAT-002, CAT-006 | ⏳ |

**Entregável Sprint 2.3:** API de categorias funcionando

---

### Sprint 2.4: Módulo de Categorias - Frontend (Semana 6-7)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-CAT-001 | Criar hook useCategories (TanStack Query) | 🟣 | 🔴 | 2h | FE-007 | ⏳ |
| FE-CAT-002 | Criar página de lista de categorias (/categories) | 🟣 | 🔴 | 4h | FE-006, FE-CAT-001 | ⏳ |
| FE-CAT-003 | Criar formulário de categoria | 🟣 | 🔴 | 4h | FE-003, FE-004 | ⏳ |
| FE-CAT-004 | Integrar com API de categorias | 🟣 | 🔴 | 4h | FE-CAT-002, CAT-009 | ⏳ |
| FE-CAT-005 | Adicionar seleção de categoria no formulário de transação | 🟣 | 🔴 | 4h | FE-TX-006, FE-CAT-001 | ⏳ |
| FE-CAT-006 | Criar componente de seleção de categoria (combobox) | 🟣 | 🟡 | 4h | FE-003, FE-CAT-001 | ⏳ |

**Entregável Sprint 2.4:** Módulo de categorias funcionando no frontend

---

### Sprint 2.5: Melhorias Frontend (Semana 7)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-IMP-001 | Implementar atualização de saldo em tempo real | 🟣 | 🟡 | 4h | FE-ACC-002, INT-001 | ⏳ |
| FE-IMP-002 | Adicionar filtros avançados em transações (data, tipo, categoria) | 🟣 | 🟡 | 6h | FE-TX-002 | ⏳ |
| FE-IMP-003 | Implementar paginação em listas | 🟣 | 🟡 | 4h | FE-TX-002, FE-ACC-002 | ⏳ |
| FE-IMP-004 | Implementar ordenação em tabelas | 🟣 | 🟡 | 2h | FE-IMP-003 | ⏳ |
| FE-IMP-005 | Criar componente Toast para notificações | 🟣 | 🟡 | 2h | FE-003 | ⏳ |
| FE-IMP-006 | Criar componente Dialog de confirmação | 🟣 | 🟡 | 2h | FE-003 | ⏳ |
| FE-IMP-007 | Criar componente EmptyState | 🟣 | 🟡 | 2h | FE-003 | ⏳ |
| FE-IMP-008 | Melhorar loading states em todos os módulos | 🟣 | 🟡 | 4h | FE-IMP-005 | ⏳ |
| FE-IMP-009 | Melhorar error handling em todos os módulos | 🟣 | 🟡 | 4h | FE-IMP-005 | ⏳ |

**Entregável Sprint 2.5:** Interface melhorada com filtros, paginação e feedback visual

---

### Sprint 2.6: Validações e Error Handling (Semana 7-8)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| VAL-001 | Implementar validações customizadas no backend | 🔵 | 🟡 | 4h | BE-002 | ⏳ |
| VAL-002 | Melhorar error handling no backend (error types) | 🔵 | 🟡 | 4h | BE-003 | ⏳ |
| VAL-003 | Criar middleware de tratamento de erros global | 🔵 | 🟡 | 4h | VAL-002 | ⏳ |
| VAL-004 | Implementar validações no frontend (Zod schemas) | 🟣 | 🟡 | 4h | FE-004 | ⏳ |
| VAL-005 | Melhorar mensagens de erro no frontend | 🟣 | 🟡 | 2h | VAL-004 | ⏳ |
| LOG-001 | Configurar logging estruturado completo | 🔵 | 🟡 | 4h | BE-006 | ⏳ |
| LOG-002 | Adicionar request ID em todas as requisições | 🔵 | 🟡 | 2h | LOG-001 | ⏳ |

**Entregável Sprint 2.6:** Sistema robusto de validação e tratamento de erros

---

### Sprint 2.7: Testes de Integração (Semana 8)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| TEST-INT-001 | Criar testes de integração para Identity Context | 🟤 | 🟡 | 4h | ID-013 | ⏳ |
| TEST-INT-002 | Criar testes de integração para Account Context | 🟤 | 🟡 | 4h | AC-011 | ⏳ |
| TEST-INT-003 | Criar testes de integração para Transaction Context | 🟤 | 🟡 | 4h | TX-015 | ⏳ |
| TEST-INT-004 | Criar testes de integração para Category Context | 🟤 | 🟡 | 4h | CAT-010 | ⏳ |
| TEST-INT-005 | Criar testes E2E básicos (autenticação → criar conta → criar transação) | 🟤 | 🟡 | 8h | FE-AUTH-010, FE-ACC-010, FE-TX-010 | ⏳ |

**Entregável Sprint 2.7:** Suite de testes de integração completa

---

## 📊 Fase 3: Funcionalidades Essenciais (4-5 semanas)

### Sprint 3.1: Budget Context - Backend (Semana 9-10)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| BUD-001 | Criar value object BudgetID | 🔵 | 🔴 | 1h | BE-001 | ⏳ |
| BUD-002 | Criar entidade Budget (agregado raiz) | 🔵 | 🔴 | 6h | BUD-001, SK-001, CAT-001 | ⏳ |
| BUD-003 | Criar interface BudgetRepository | 🔵 | 🔴 | 1h | BUD-002 | ⏳ |
| BUD-004 | Implementar GormBudgetRepository | 🔵 | 🔴 | 6h | BUD-003, BE-004 | ⏳ |
| BUD-005 | Criar migration para tabela budgets | 🔵 | 🔴 | 2h | BUD-004 | ⏳ |
| BUD-006 | Implementar use cases de Budget (CRUD) | 🔵 | 🔴 | 8h | BUD-002, BUD-003 | ⏳ |
| BUD-007 | Implementar cálculo de progresso do orçamento | 🔵 | 🔴 | 4h | BUD-002, TX-005 | ⏳ |
| BUD-008 | Criar BudgetHandler | 🔵 | 🔴 | 4h | BUD-006 | ⏳ |
| BUD-009 | Configurar rotas de budgets (/api/v1/budgets/*) | 🔵 | 🔴 | 2h | BUD-008, ID-011 | ⏳ |
| BUD-010 | Adicionar anotações Swagger para Budget | 🔵 | 🟡 | 2h | BUD-008, DOC-001 | ⏳ |
| BUD-011 | Testes unitários para Budget Context | 🟤 | 🟡 | 6h | BUD-002, BUD-006 | ⏳ |

**Entregável Sprint 3.1:** API de orçamentos funcionando

---

### Sprint 3.2: Recurring Transactions - Backend (Semana 10)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| REC-001 | Adicionar campos de recorrência na entidade Transaction | 🔵 | 🟡 | 4h | TX-004 | ⏳ |
| REC-002 | Criar serviço de processamento de transações recorrentes | 🔵 | 🟡 | 8h | TX-004, TX-008 | ⏳ |
| REC-003 | Criar job/cron para processar transações recorrentes | 🔵 | 🟡 | 4h | REC-002 | ⏳ |
| REC-004 | Testes para transações recorrentes | 🟤 | 🟡 | 4h | REC-002 | ⏳ |

**Entregável Sprint 3.2:** Sistema de transações recorrentes funcionando

---

### Sprint 3.3: Reporting Context - Backend (Semana 10-11)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| REP-001 | Criar use case para relatório mensal | 🔵 | 🔴 | 6h | TX-005, CAT-003 | ⏳ |
| REP-002 | Criar use case para relatório anual | 🔵 | 🔴 | 4h | REP-001 | ⏳ |
| REP-003 | Criar use case para relatório por categoria | 🔵 | 🔴 | 4h | REP-001 | ⏳ |
| REP-004 | Criar use case para receitas vs despesas | 🔵 | 🔴 | 4h | REP-001 | ⏳ |
| REP-005 | Criar ReportHandler | 🔵 | 🔴 | 4h | REP-001, REP-002, REP-003, REP-004 | ⏳ |
| REP-006 | Configurar rotas de reports (/api/v1/reports/*) | 🔵 | 🔴 | 2h | REP-005, ID-011 | ⏳ |
| REP-007 | Adicionar anotações Swagger para Reports | 🔵 | 🟡 | 2h | REP-005, DOC-001 | ⏳ |
| REP-008 | Implementar cache de relatórios (Redis) | 🔵 | 🟡 | 4h | REP-001, SETUP-004 | ⏳ |
| REP-009 | Testes para Reporting Context | 🟤 | 🟡 | 6h | REP-001 | ⏳ |

**Entregável Sprint 3.3:** API de relatórios funcionando com cache

---

### Sprint 3.4: Cache e Performance - Backend (Semana 11)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| PERF-001 | Configurar Redis no backend | 🔵 | 🟡 | 2h | SETUP-004 | ⏳ |
| PERF-002 | Implementar cache em AccountRepository | 🔵 | 🟡 | 4h | AC-004, PERF-001 | ⏳ |
| PERF-003 | Implementar cache em CategoryRepository | 🔵 | 🟡 | 4h | CAT-004, PERF-001 | ⏳ |
| PERF-004 | Implementar paginação no backend | 🔵 | 🟡 | 4h | TX-009, AC-007 | ⏳ |
| PERF-005 | Implementar rate limiting | 🔵 | 🟡 | 4h | BE-003, PERF-001 | ⏳ |
| PERF-006 | Criar índices no banco de dados | 🔵 | 🟡 | 4h | BE-004 | ⏳ |

**Entregável Sprint 3.4:** Sistema otimizado com cache e paginação

---

### Sprint 3.5: Módulo de Orçamento - Frontend (Semana 11-12)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-BUD-001 | Criar hook useBudgets (TanStack Query) | 🟣 | 🔴 | 2h | FE-007 | ⏳ |
| FE-BUD-002 | Criar página de dashboard de orçamentos (/budget) | 🟣 | 🔴 | 6h | FE-006, FE-BUD-001 | ⏳ |
| FE-BUD-003 | Criar componente de progresso de orçamento | 🟣 | 🔴 | 4h | FE-003 | ⏳ |
| FE-BUD-004 | Criar formulário de orçamento | 🟣 | 🔴 | 4h | FE-003, FE-004 | ⏳ |
| FE-BUD-005 | Integrar com API de budgets | 🟣 | 🔴 | 4h | FE-BUD-002, BUD-010 | ⏳ |
| FE-BUD-006 | Implementar alertas de limite de orçamento | 🟣 | 🟡 | 4h | FE-BUD-003 | ⏳ |

**Entregável Sprint 3.5:** Módulo de orçamento funcionando no frontend

---

### Sprint 3.6: Módulo de Relatórios - Frontend (Semana 12-13)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-REP-001 | Instalar e configurar Recharts | 🟣 | 🔴 | 2h | FE-001 | ⏳ |
| FE-REP-002 | Criar hook useReports (TanStack Query) | 🟣 | 🔴 | 2h | FE-007 | ⏳ |
| FE-REP-003 | Criar página de relatórios (/reports) | 🟣 | 🔴 | 6h | FE-006, FE-REP-002 | ⏳ |
| FE-REP-004 | Criar componente de gráfico receitas vs despesas | 🟣 | 🔴 | 4h | FE-REP-001 | ⏳ |
| FE-REP-005 | Criar componente de gráfico por categoria | 🟣 | 🔴 | 4h | FE-REP-001 | ⏳ |
| FE-REP-006 | Criar componente de gráfico de tendências temporais | 🟣 | 🔴 | 4h | FE-REP-001 | ⏳ |
| FE-REP-007 | Criar filtros de período (mensal, anual) | 🟣 | 🔴 | 4h | FE-REP-003 | ⏳ |
| FE-REP-008 | Integrar com API de relatórios | 🟣 | 🔴 | 4h | FE-REP-003, REP-007 | ⏳ |
| FE-REP-009 | Implementar exportação CSV | 🟣 | 🟡 | 4h | FE-REP-003 | ⏳ |
| FE-REP-010 | Implementar exportação PDF | 🟣 | 🟢 | 6h | FE-REP-003 | ⏳ |

**Entregável Sprint 3.6:** Módulo de relatórios com gráficos funcionando

---

### Sprint 3.7: Melhorias Gerais Frontend (Semana 13)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-GEN-001 | Implementar dark mode (shadcn/ui) | 🟣 | 🟡 | 4h | FE-003 | ⏳ |
| FE-GEN-002 | Melhorar responsividade mobile | 🟣 | 🟡 | 8h | FE-006 | ⏳ |
| FE-GEN-003 | Implementar lazy loading de rotas | 🟣 | 🟡 | 2h | FE-001 | ⏳ |
| FE-GEN-004 | Implementar code splitting | 🟣 | 🟡 | 2h | FE-001 | ⏳ |
| FE-GEN-005 | Adicionar ARIA labels para acessibilidade | 🟣 | 🟡 | 4h | FE-003 | ⏳ |
| FE-GEN-006 | Otimizar imagens (Next.js Image) | 🟣 | 🟢 | 2h | FE-001 | ⏳ |

**Entregável Sprint 3.7:** Interface otimizada e acessível

---

## 🚀 Fase 4: Produção e Performance (3-4 semanas)

### Sprint 4.1: Observabilidade Backend (Semana 14)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| OBS-001 | Configurar Prometheus para métricas | 🔵 | 🔴 | 4h | BE-003 | ⏳ |
| OBS-002 | Criar middleware de métricas HTTP | 🔵 | 🔴 | 4h | OBS-001 | ⏳ |
| OBS-003 | Adicionar métricas de negócio (transações criadas, etc.) | 🔵 | 🟡 | 4h | OBS-001 | ⏳ |
| OBS-004 | Configurar OpenTelemetry para tracing | 🔵 | 🟡 | 6h | BE-002 | ⏳ |
| OBS-005 | Configurar Grafana para visualização | 🟠 | 🟡 | 4h | OBS-001 | ⏳ |
| OBS-006 | Criar dashboards no Grafana | 🟠 | 🟡 | 4h | OBS-005 | ⏳ |

**Entregável Sprint 4.1:** Sistema de observabilidade completo

---

### Sprint 4.2: Segurança e Produção (Semana 14-15)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| SEC-001 | Configurar headers de segurança (Helmet) | 🔵 | 🔴 | 2h | BE-003 | ⏳ |
| SEC-002 | Implementar rate limiting robusto | 🔵 | 🔴 | 4h | PERF-005 | ⏳ |
| SEC-003 | Configurar CORS para produção | 🔵 | 🔴 | 2h | BE-003 | ⏳ |
| SEC-004 | Implementar graceful shutdown | 🔵 | 🔴 | 4h | BE-003 | ⏳ |
| SEC-005 | Configurar health checks robustos (liveness/readiness) | 🔵 | 🔴 | 4h | BE-005 | ⏳ |
| SEC-006 | Revisar e melhorar validações de segurança | 🔵 | 🟡 | 4h | VAL-001 | ⏳ |

**Entregável Sprint 4.2:** Sistema seguro e pronto para produção

---

### Sprint 4.3: CI/CD Completo (Semana 15)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| CI-001 | Configurar GitHub Actions para testes | 🟠 | 🔴 | 4h | SETUP-006 | ⏳ |
| CI-002 | Configurar build e push de Docker image | 🟠 | 🔴 | 4h | CI-001, BE-007 | ⏳ |
| CI-003 | Configurar deploy automático (staging) | 🟠 | 🟡 | 4h | CI-002 | ⏳ |
| CI-004 | Configurar CI/CD para frontend | 🟠 | 🟡 | 4h | FE-009 | ⏳ |
| CI-005 | Configurar deploy frontend (Vercel/Netlify) | 🟠 | 🟡 | 2h | CI-004 | ⏳ |

**Entregável Sprint 4.3:** Pipeline CI/CD completo

---

### Sprint 4.4: Testes Frontend (Semana 15-16)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| FE-TEST-001 | Configurar Vitest e React Testing Library | 🟤 | 🟡 | 4h | FE-001 | ⏳ |
| FE-TEST-002 | Criar testes unitários para componentes | 🟤 | 🟡 | 8h | FE-TEST-001 | ⏳ |
| FE-TEST-003 | Criar testes de integração frontend-backend | 🟤 | 🟡 | 8h | FE-TEST-001 | ⏳ |
| FE-TEST-004 | Configurar Playwright para E2E | 🟤 | 🟡 | 4h | FE-001 | ⏳ |
| FE-TEST-005 | Criar testes E2E principais | 🟤 | 🟡 | 8h | FE-TEST-004 | ⏳ |
| FE-TEST-006 | Testes de acessibilidade | 🟤 | 🟢 | 4h | FE-TEST-001 | ⏳ |

**Entregável Sprint 4.4:** Suite de testes frontend completa

---

### Sprint 4.5: Otimizações e Deploy (Semana 16)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| OPT-001 | Otimizar queries do banco de dados | 🔵 | 🟡 | 4h | PERF-006 | ⏳ |
| OPT-002 | Implementar backup automático | 🟠 | 🔴 | 4h | SETUP-003 | ⏳ |
| OPT-003 | Otimizar bundle size do frontend | 🟣 | 🟡 | 4h | FE-001 | ⏳ |
| OPT-004 | Configurar error tracking (Sentry) | 🟣 | 🟡 | 4h | FE-001 | ⏳ |
| OPT-005 | Configurar PWA (opcional) | 🟣 | 🟢 | 6h | FE-001 | ⏳ |
| DEPLOY-001 | Configurar ambiente de produção | 🟠 | 🔴 | 8h | CI-002, CI-005 | ⏳ |
| DEPLOY-002 | Documentação de deploy | ⚪ | 🔴 | 4h | DEPLOY-001 | ⏳ |

**Entregável Sprint 4.5:** Sistema deployado em produção

---

## 🎨 Fase 5: Funcionalidades Avançadas (4-5 semanas)

### Sprint 5.1: Investment Context (Semana 17-18)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| INV-001 | Criar value objects para Investment | 🔵 | 🟡 | 4h | BE-001 | ⏳ |
| INV-002 | Criar entidade Investment | 🔵 | 🟡 | 6h | INV-001 | ⏳ |
| INV-003 | Implementar repositório e use cases | 🔵 | 🟡 | 8h | INV-002 | ⏳ |
| INV-004 | Criar handlers e rotas | 🔵 | 🟡 | 4h | INV-003 | ⏳ |
| FE-INV-001 | Criar módulo de investimentos no frontend | 🟣 | 🟡 | 12h | INV-004 | ⏳ |

**Entregável Sprint 5.1:** Módulo de investimentos completo

---

### Sprint 5.2: Goal Context (Semana 18-19)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| GOAL-001 | Criar value objects para Goal | 🔵 | 🟡 | 4h | BE-001 | ⏳ |
| GOAL-002 | Criar entidade Goal | 🔵 | 🟡 | 6h | GOAL-001 | ⏳ |
| GOAL-003 | Implementar repositório e use cases | 🔵 | 🟡 | 8h | GOAL-002 | ⏳ |
| GOAL-004 | Criar handlers e rotas | 🔵 | 🟡 | 4h | GOAL-003 | ⏳ |
| FE-GOAL-001 | Criar módulo de metas no frontend | 🟣 | 🟡 | 12h | GOAL-004 | ⏳ |

**Entregável Sprint 5.2:** Módulo de metas completo

---

### Sprint 5.3: Notification Context (Semana 19-20)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| NOT-001 | Criar entidade Notification | 🔵 | 🟡 | 4h | BE-001 | ⏳ |
| NOT-002 | Implementar repositório e use cases | 🔵 | 🟡 | 6h | NOT-001 | ⏳ |
| NOT-003 | Criar handlers e rotas | 🔵 | 🟡 | 4h | NOT-002 | ⏳ |
| NOT-004 | Implementar WebSocket para notificações em tempo real | 🔵 | 🟡 | 8h | NOT-002 | ⏳ |
| FE-NOT-001 | Criar módulo de notificações no frontend | 🟣 | 🟡 | 8h | NOT-003 | ⏳ |
| FE-NOT-002 | Integrar WebSocket no frontend | 🟣 | 🟡 | 4h | FE-NOT-001, NOT-004 | ⏳ |

**Entregável Sprint 5.3:** Sistema de notificações em tempo real

---

### Sprint 5.4: Dashboard Completo (Semana 20)

| ID | Tarefa | Tipo | Prioridade | Esforço | Dependências | Status |
|----|--------|------|------------|---------|--------------|--------|
| DASH-001 | Criar API de dashboard (métricas agregadas) | 🔵 | 🟡 | 6h | REP-001, AC-003 | ⏳ |
| FE-DASH-001 | Criar dashboard principal no frontend | 🟣 | 🟡 | 12h | DASH-001, FE-REP-001 | ⏳ |
| FE-DASH-002 | Adicionar cards de métricas principais | 🟣 | 🟡 | 4h | FE-DASH-001 | ⏳ |
| FE-DASH-003 | Adicionar gráficos resumidos | 🟣 | 🟡 | 4h | FE-DASH-001 | ⏳ |

**Entregável Sprint 5.4:** Dashboard completo e funcional

---

## 📝 Notas Importantes

### Priorização
- 🔴 **Alta**: Essencial para MVP e funcionalidade básica
- 🟡 **Média**: Importante mas pode ser feito depois
- 🟢 **Baixa**: Nice to have, pode ser opcional

### Estimativas
- Estimativas são em horas de trabalho focado
- Ajuste conforme sua velocidade e experiência
- Considere tempo para debugging e imprevistos

### Dependências
- Respeite as dependências entre tarefas
- Algumas tarefas podem ser paralelizadas
- Backend geralmente deve estar pronto antes do frontend correspondente

### Testes
- Testes devem ser escritos junto com o código
- Não deixe testes para o final
- Cobertura mínima recomendada: 70%

### Documentação
- Documente decisões importantes
- Mantenha README atualizado
- Documente APIs no Swagger

---

## 🎯 Próximos Passos

1. **Revisar este documento** e ajustar conforme necessário
2. **Priorizar tarefas** baseado em suas necessidades
3. **Criar issues no GitHub** para cada tarefa
4. **Começar pela Fase 1, Sprint 1.1**
5. **Revisar e atualizar** este documento conforme o progresso

---

**Última atualização:** Baseado no PLANEJAMENTO_GO.md


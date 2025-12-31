# 🧪 Guia de Testes E2E com Cypress

## 📋 Pré-requisitos

1. **Servidor de desenvolvimento rodando**
   - Para testes E2E, você precisa do frontend rodando
   - Opções:
     - `npm run preview` (produção) - porta 4173
     - `npm run dev` (desenvolvimento) - porta 3000

2. **Variável de ambiente (opcional)**
   - Por padrão, Cypress usa `http://localhost:4173`
   - Para mudar: `CYPRESS_BASE_URL=http://localhost:3000 npm run test:e2e:open`

## 🚀 Executando os Testes

### Interface Gráfica (Recomendado para Desenvolvimento)

A interface gráfica do Cypress permite:
- Ver os testes executando em tempo real
- Debug interativo
- Time-travel debugging
- Screenshots e vídeos automáticos

```bash
# Opção 1: Via npm script
npm run test:e2e:open

# Opção 2: Via script helper
./scripts/cypress-open.sh

# Opção 3: Diretamente
npx cypress open
```

**Nota**: A interface gráfica requer um ambiente com display gráfico (X11, Wayland, etc.)

### Modo Headless (CI/CD)

Executa os testes sem interface gráfica:

```bash
# Execução completa
npm run test:e2e

# Com browser visível (headed)
npm run test:e2e:headed

# Browser específico
npx cypress run --browser chrome
npx cypress run --browser firefox
npx cypress run --browser edge
```

## 📁 Estrutura de Testes

```
cypress/
├── e2e/                    # Testes E2E
│   ├── auth.cy.ts         # Testes de autenticação
│   ├── transactions.cy.ts # Testes de transações
│   ├── accounts.cy.ts    # Testes de contas
│   └── accessibility.cy.ts # Testes de acessibilidade
└── support/
    ├── e2e.ts            # Configuração global
    └── commands.ts       # Comandos customizados
```

## 🛠️ Comandos Customizados

O Cypress foi configurado com comandos customizados para facilitar os testes:

### `cy.login()`
Faz login automaticamente com token mock:
```typescript
cy.login()
```

### `cy.setAuthToken(token)`
Define um token de autenticação:
```typescript
cy.setAuthToken('my-token-123')
```

### `cy.mockApi(method, url, response, status)`
Mocka uma resposta de API:
```typescript
cy.mockApi('GET', '/api/v1/accounts', { accounts: [] }, 200)
```

## 🎯 Testes Disponíveis

### Autenticação (`auth.cy.ts`)
- ✅ Redirecionamento para login quando não autenticado
- ✅ Exibição do formulário de login
- ✅ Validação de erros
- ✅ Login bem-sucedido
- ✅ Logout

### Transações (`transactions.cy.ts`)
- ✅ Exibição da lista de transações
- ✅ Filtro por tipo
- ✅ Abertura do formulário de criação
- ✅ Criação de nova transação

### Contas (`accounts.cy.ts`)
- ✅ Exibição da lista de contas
- ✅ Abertura do formulário de criação
- ✅ Criação de nova conta

### Acessibilidade (`accessibility.cy.ts`)
- ✅ Sem violações de acessibilidade na página de login
- ✅ Sem violações no dashboard
- ✅ Sem violações na página de transações
- ✅ Hierarquia de headings correta
- ✅ Labels de formulário corretos
- ✅ Labels de botões corretos

## 🐳 Executando no Docker

### Opção 1: Docker Compose (Recomendado)

A forma mais fácil de executar os testes E2E em containers é usando Docker Compose:

```bash
# Na raiz do projeto
docker-compose -f docker-compose.test.yml up frontend-e2e --build
```

Este comando irá:
1. Construir e iniciar o container `frontend-preview` (servidor da aplicação)
2. Aguardar o servidor estar saudável
3. Executar os testes E2E no container `frontend-e2e`

**Vantagens:**
- ✅ Não precisa instalar Node.js, npm ou Cypress localmente
- ✅ Ambiente isolado e reproduzível
- ✅ Configuração automática de rede entre containers
- ✅ Healthcheck garante que o servidor está pronto antes dos testes

Para mais detalhes, consulte: [`docs/TESTES_E2E_DOCKER.md`](../../docs/TESTES_E2E_DOCKER.md)

### Opção 2: Docker Manual

Para executar os testes E2E manualmente no Docker:

```bash
# Build da imagem de teste
docker build -t gestao-financeira-frontend-test -f frontend/Dockerfile.test frontend/

# Executar testes (requer servidor rodando)
docker run --rm --network host \
  -v "$(pwd)/frontend:/app" \
  -w /app \
  -e CYPRESS_BASE_URL=http://localhost:4173 \
  gestao-financeira-frontend-test \
  npm run test:e2e
```

## 📸 Screenshots e Vídeos

- Screenshots são salvos automaticamente em falhas: `cypress/screenshots/`
- Vídeos estão desabilitados por padrão (configurável em `cypress.config.ts`)

## 🔧 Configuração

A configuração principal está em `cypress.config.ts`:

- **baseUrl**: URL base da aplicação (padrão: `http://localhost:4173`)
- **viewportWidth/Height**: Tamanho da viewport (1280x720)
- **timeouts**: Timeouts padrão para comandos (10s)

## 🐛 Debug

### Modo Debug
Adicione `.debug()` em qualquer comando:
```typescript
cy.get('button').debug().click()
```

### Pausar Execução
Use `cy.pause()` para pausar a execução:
```typescript
cy.pause()
```

### Logs no Console
Cypress mostra logs detalhados no console durante a execução.

## 📚 Recursos

- [Documentação Oficial do Cypress](https://docs.cypress.io/)
- [Best Practices](https://docs.cypress.io/guides/references/best-practices)
- [Comandos Customizados](https://docs.cypress.io/api/cypress-api/custom-commands)


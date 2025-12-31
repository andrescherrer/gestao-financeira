# 🐳 Executando Testes E2E em Containers Docker

Este documento explica como executar os testes E2E (End-to-End) da aplicação usando Docker, **sem precisar instalar nada na sua máquina local**, exceto Docker e Docker Compose.

## 📋 Pré-requisitos

- **Docker** 20.10+ instalado
- **Docker Compose** 2.0+ instalado
- **Git** (para clonar o repositório)

**Não é necessário:**
- ❌ Node.js instalado
- ❌ npm/yarn instalado
- ❌ Cypress instalado
- ❌ Navegadores instalados

## 🏗️ Arquitetura

A configuração de testes E2E em containers utiliza:

1. **`frontend-preview`**: Container que serve a aplicação frontend compilada
2. **`frontend-e2e`**: Container que executa os testes Cypress
3. **Rede Docker**: Comunicação entre containers

```
┌─────────────────────┐
│  frontend-preview   │  ← Serve a aplicação em http://frontend-preview:4173
│  (Porta 4173:4173)  │
└──────────┬──────────┘
           │
           │ HTTP
           │
┌──────────▼──────────┐
│   frontend-e2e      │  ← Executa testes Cypress
│   (Cypress)         │
└─────────────────────┘
```

## 🚀 Executando os Testes

### Opção 1: Executar Todos os Testes E2E

```bash
# Na raiz do projeto
docker-compose -f docker-compose.test.yml up frontend-e2e --build
```

Este comando irá:
1. Construir a imagem do frontend-preview
2. Construir a imagem do frontend-e2e
3. Iniciar o frontend-preview e aguardar estar saudável
4. Executar todos os testes E2E
5. Mostrar os resultados no terminal

### Opção 2: Executar com Logs Detalhados

```bash
docker-compose -f docker-compose.test.yml up frontend-e2e --build --abort-on-container-exit
```

O flag `--abort-on-container-exit` faz com que o Docker Compose pare automaticamente quando os testes terminarem.

### Opção 3: Executar e Remover Containers Após Conclusão

```bash
docker-compose -f docker-compose.test.yml up frontend-e2e --build --abort-on-container-exit --remove-orphans
```

### Opção 4: Executar Apenas o Frontend Preview (para debug)

Se você quiser apenas iniciar o servidor frontend para testar manualmente:

```bash
docker-compose -f docker-compose.test.yml up frontend-preview --build
```

O frontend estará disponível em: `http://localhost:4173`

## 📊 Visualizando Resultados

### Screenshots de Falhas

Quando um teste falha, o Cypress salva automaticamente screenshots em:
```
frontend/cypress/screenshots/
```

### Logs dos Testes

Os logs dos testes aparecem diretamente no terminal durante a execução.

### Exemplo de Saída

```
frontend-e2e-test  |   Running:  auth.cy.ts                                    (1 of 4)
frontend-e2e-test  | 
frontend-e2e-test  |   ✓ should redirect to login when not authenticated (1234ms)
frontend-e2e-test  |   ✓ should show login form (567ms)
frontend-e2e-test  |   ✓ should display validation errors on invalid login (890ms)
frontend-e2e-test  |   ✓ should login successfully with valid credentials (2345ms)
frontend-e2e-test  |   ✓ should logout successfully (1234ms)
frontend-e2e-test  | 
frontend-e2e-test  |   Running:  transactions.cy.ts                            (2 of 4)
frontend-e2e-test  |   ...
```

## 🔧 Configuração

### Base da Imagem

A imagem de teste usa **Debian Slim** (`node:20-slim`) ao invés de Alpine Linux para melhor compatibilidade com Cypress e suas dependências. Isso garante que todas as bibliotecas necessárias (GTK, NSS, X11, etc.) estejam disponíveis sem problemas de compatibilidade.

### Servidor Frontend

O servidor frontend usa **`serve`** ao invés de `vite preview` para evitar problemas de validação de hostname em containers Docker. O `serve` é mais simples e não tem restrições de acesso baseadas em hostname.

### Mock do React

A aplicação usa a biblioteca `sonner` que tenta importar React dinamicamente. Para resolver isso, foi criado um plugin do Vite (`vite.config.plugins.ts`) que intercepta as importações de React e as substitui por um mock, permitindo que a aplicação Vue funcione corretamente.

### Variáveis de Ambiente

As variáveis de ambiente podem ser configuradas no `docker-compose.test.yml` ou via arquivo `.env`:

```bash
# .env (opcional)
CYPRESS_BASE_URL=http://frontend-preview:4173
NODE_ENV=test
```

### Portas

- **Frontend Preview**: `4173:4173` (acessível em `http://localhost:4173`)
- **Cypress**: Executa internamente, não expõe portas

### Volumes

Os seguintes volumes são montados:

- `./frontend:/app` - Código fonte (para desenvolvimento)
- `/app/node_modules` - Dependências (volume anônimo para performance)
- `./frontend/cypress/screenshots:/app/cypress/screenshots` - Screenshots de falhas
- `./frontend/cypress/videos:/app/cypress/videos` - Vídeos (se habilitado)

## 🐛 Troubleshooting

### Problema: "Cannot connect to frontend-preview"

**Solução**: Verifique se o `frontend-preview` está saudável:

```bash
docker-compose -f docker-compose.test.yml ps
```

Aguarde até que o healthcheck passe antes de executar os testes.

### Problema: "Cypress binary not found"

**Solução**: Reconstrua a imagem:

```bash
docker-compose -f docker-compose.test.yml build --no-cache frontend-e2e
```

### Problema: Testes muito lentos

**Solução**: Os testes podem ser lentos na primeira execução devido ao build. Execuções subsequentes serão mais rápidas devido ao cache do Docker.

### Problema: Screenshots não aparecem

**Solução**: Verifique as permissões do diretório:

```bash
chmod -R 755 frontend/cypress/screenshots
```

### Problema: "Xvfb: command not found"

**Solução**: A imagem já inclui Xvfb. Se o erro persistir, reconstrua:

```bash
docker-compose -f docker-compose.test.yml build --no-cache frontend-e2e
```

### Problema: "Cypress failed to start" ou "Missing library or dependency"

**Solução**: A imagem usa Debian Slim que inclui todas as dependências necessárias. Se o erro persistir:

1. Reconstrua a imagem sem cache:
```bash
docker-compose -f docker-compose.test.yml build --no-cache frontend-e2e
```

2. Verifique se o Xvfb está rodando:
```bash
docker-compose -f docker-compose.test.yml run --rm frontend-e2e ps aux | grep Xvfb
```

3. Execute os testes manualmente para ver o erro completo:
```bash
docker-compose -f docker-compose.test.yml run --rm frontend-e2e sh -c "Xvfb :99 -screen 0 1280x720x24 > /dev/null 2>&1 & export DISPLAY=:99 && sleep 5 && npm run test:e2e"
```

### Problema: "403 Forbidden" ao acessar o frontend

**Solução**: O servidor frontend usa `serve` ao invés de `vite preview` para evitar problemas de validação de hostname. Se o erro persistir:

1. Verifique se o frontend-preview está rodando:
```bash
docker-compose -f docker-compose.test.yml ps frontend-preview
```

2. Teste o acesso manualmente:
```bash
docker-compose -f docker-compose.test.yml run --rm frontend-e2e wget -O- http://frontend-preview:4173
```

### Problema: "Failed to resolve module specifier 'react'"

**Solução**: Este problema foi resolvido com um plugin do Vite que mocka React para a biblioteca `sonner`. Se o erro persistir:

1. Verifique se o arquivo `vite.config.plugins.ts` existe
2. Reconstrua a imagem do frontend-preview:
```bash
docker-compose -f docker-compose.test.yml build --no-cache frontend-preview
```

## 🔄 Workflow de Desenvolvimento

### 1. Desenvolvimento Local (sem Docker)

Se você tem Node.js instalado localmente:

```bash
cd frontend
npm run preview  # Terminal 1
npm run test:e2e:open  # Terminal 2
```

### 2. Testes em CI/CD (com Docker)

Use o Docker Compose para testes automatizados:

```bash
docker-compose -f docker-compose.test.yml up frontend-e2e --build --abort-on-container-exit
```

### 3. Debug de Testes Específicos

Para executar apenas um arquivo de teste específico, você pode modificar temporariamente o comando:

```bash
docker-compose -f docker-compose.test.yml run --rm frontend-e2e npx cypress run --spec "cypress/e2e/auth.cy.ts"
```

## 📝 Estrutura de Arquivos

```
.
├── docker-compose.test.yml          # Configuração Docker Compose para testes
├── frontend/
│   ├── Dockerfile.test              # Dockerfile para ambiente de testes
│   ├── vite.config.plugins.ts       # Plugin do Vite para mockar React
│   ├── scripts/
│   │   └── run-e2e-docker.sh        # Script helper para executar E2E no Docker
│   ├── cypress/
│   │   ├── e2e/
│   │   │   ├── auth.cy.ts
│   │   │   ├── transactions.cy.ts
│   │   │   ├── accounts.cy.ts
│   │   │   └── accessibility.cy.ts
│   │   ├── screenshots/             # Screenshots de falhas
│   │   └── videos/                  # Vídeos (se habilitado)
│   └── cypress.config.cjs           # Configuração do Cypress
└── docs/
    └── TESTES_E2E_DOCKER.md         # Este documento
```

## 🎯 Comandos Úteis

### Limpar Containers e Volumes

```bash
# Parar e remover containers
docker-compose -f docker-compose.test.yml down

# Remover também volumes
docker-compose -f docker-compose.test.yml down -v

# Limpar imagens não utilizadas
docker image prune -f
```

### Reconstruir do Zero

```bash
# Remover tudo e reconstruir
docker-compose -f docker-compose.test.yml down -v
docker-compose -f docker-compose.test.yml build --no-cache
docker-compose -f docker-compose.test.yml up frontend-e2e --abort-on-container-exit
```

### Ver Logs em Tempo Real

```bash
docker-compose -f docker-compose.test.yml up frontend-e2e --build --follow
```

### Executar em Modo Interativo

Para debug avançado, você pode entrar no container:

```bash
docker-compose -f docker-compose.test.yml run --rm frontend-e2e sh
```

Dentro do container:

```bash
# Verificar se o frontend está acessível
curl http://frontend-preview:4173

# Executar testes manualmente
npm run test:e2e

# Executar teste específico
npx cypress run --spec "cypress/e2e/auth.cy.ts"
```

## 🔐 Segurança

- Os containers de teste não expõem portas desnecessárias
- As imagens usam usuários não-root quando possível
- Os volumes são montados apenas para os diretórios necessários

## 📚 Recursos Adicionais

- [Documentação do Cypress](https://docs.cypress.io/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Cypress Best Practices](https://docs.cypress.io/guides/references/best-practices)

## ✅ Checklist de Execução

Antes de executar os testes, certifique-se de:

- [ ] Docker e Docker Compose estão instalados e funcionando
- [ ] Você está na raiz do projeto
- [ ] O arquivo `docker-compose.test.yml` existe
- [ ] Você tem permissões de escrita em `frontend/cypress/screenshots`

## 🎉 Exemplo Completo

```bash
# 1. Navegar para a raiz do projeto
cd /home/andre/Projetos/gestao-financeira

# 2. Executar testes E2E
docker-compose -f docker-compose.test.yml up frontend-e2e --build --abort-on-container-exit

# 3. Ver resultados
# Os testes serão executados e os resultados aparecerão no terminal

# 4. Limpar após conclusão
docker-compose -f docker-compose.test.yml down
```

---

**Última atualização**: 2025-12-31  
**Versão**: 1.2.0

### Changelog

- **v1.2.0** (2025-12-31): 
  - Substituição de `vite preview` por `serve` para resolver problemas de 403 Forbidden
  - Implementação de plugin do Vite para mockar React (necessário para biblioteca `sonner`)
  - Correção do comando Xvfb no docker-compose
  - Melhorias na documentação de troubleshooting
- **v1.1.0** (2025-12-31): Migração de Alpine Linux para Debian Slim para melhor compatibilidade com Cypress
- **v1.0.0** (2025-12-31): Versão inicial da documentação

### Status Atual dos Testes

**Última execução**: 2025-12-31

- ✅ **Infraestrutura**: Funcionando corretamente
  - Xvfb iniciando corretamente
  - Servidor frontend acessível
  - Plugin React mock funcionando
  
- 📊 **Resultados dos Testes**:
  - **Accessibility**: 3/6 passando (violações de acessibilidade reais detectadas)
  - **Accounts**: 2/3 passando
  - **Auth**: 3/5 passando
  - **Transactions**: 0/4 passando (requer ajustes nos mocks ou na aplicação)

**Total**: 8/18 testes passando (44%)

> **Nota**: Os testes que falham são principalmente devido a problemas de lógica da aplicação ou violações de acessibilidade reais, não problemas de infraestrutura. A infraestrutura de testes E2E está funcionando corretamente.


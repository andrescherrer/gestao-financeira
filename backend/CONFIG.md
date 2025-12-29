# Configuração da API

Este documento descreve todas as variáveis de ambiente disponíveis para configurar a API.

## Variáveis de Ambiente

### Ambiente

- **ENV**: Ambiente da aplicação (`dev`, `development`, `staging`, `production`, `prod`)
  - Padrão: `dev`
  - Em produção, use `production` ou `prod`

### Servidor

- **API_PORT**: Porta do servidor HTTP
  - Padrão: `8080`

- **API_READ_TIMEOUT**: Timeout de leitura (ex: `10s`, `1m`)
  - Padrão: `10s`

- **API_WRITE_TIMEOUT**: Timeout de escrita (ex: `10s`, `1m`)
  - Padrão: `10s`

- **API_IDLE_TIMEOUT**: Timeout de inatividade (ex: `120s`, `2m`)
  - Padrão: `120s`

- **API_BODY_LIMIT**: Limite de tamanho do body em bytes (ex: `10485760` para 10MB)
  - Padrão: `10485760` (10MB)

### Banco de Dados (PostgreSQL)

- **POSTGRES_HOST**: Host do PostgreSQL
  - Padrão: `localhost`

- **POSTGRES_PORT**: Porta do PostgreSQL
  - Padrão: `5432`

- **POSTGRES_USER**: Usuário do PostgreSQL
  - Padrão: `postgres`

- **POSTGRES_PASSWORD**: Senha do PostgreSQL
  - Padrão: `postgres`
  - ⚠️ **IMPORTANTE**: Altere em produção!

- **POSTGRES_DB**: Nome do banco de dados
  - Padrão: `gestao_financeira`

- **POSTGRES_SSLMODE**: Modo SSL (`disable`, `require`, `verify-ca`, `verify-full`)
  - Padrão: `disable`

- **POSTGRES_MAX_OPEN_CONNS**: Máximo de conexões abertas no pool
  - Padrão: `25`

- **POSTGRES_MAX_IDLE_CONNS**: Máximo de conexões idle no pool
  - Padrão: `5`

- **POSTGRES_CONN_MAX_LIFETIME**: Tempo máximo de vida da conexão (ex: `5m`, `1h`)
  - Padrão: `5m`

### JWT

- **JWT_SECRET**: Chave secreta para assinatura de tokens JWT
  - Padrão: `your-secret-key-change-in-production`
  - ⚠️ **IMPORTANTE**: Use uma chave forte e aleatória em produção!

- **JWT_EXPIRATION**: Tempo de expiração do token (ex: `24h`, `7d`)
  - Padrão: `24h`

- **JWT_ISSUER**: Emissor do token JWT
  - Padrão: `gestao-financeira-api`

- **JWT_SIGNING_METHOD**: Método de assinatura (`HS256`, `HS384`, `HS512`, etc.)
  - Padrão: `HS256`

### Redis (Opcional)

- **REDIS_URL**: URL de conexão do Redis (ex: `redis://localhost:6379`)
  - Padrão: (vazio - desabilitado)
  - Se não configurado, recursos de cache serão desabilitados

- **REDIS_TTL**: Tempo de vida do cache (ex: `5m`, `1h`)
  - Padrão: `5m`

### Logging

- **LOG_LEVEL**: Nível de log (`debug`, `info`, `warn`, `error`, `fatal`, `panic`, `trace`)
  - Padrão: `info`

- **LOG_FORMAT**: Formato de log (`json` para produção, `console` para desenvolvimento)
  - Padrão: `console`

### Migrations

- **MIGRATIONS_PATH**: Caminho para o diretório de migrations
  - Padrão: `file://migrations`

### CORS

- **ALLOWED_ORIGINS**: Origens permitidas (separadas por vírgula para múltiplas origens)
  - Padrão: `http://localhost:3000`

- **CORS_MAX_AGE**: Tempo de cache do preflight em segundos
  - Padrão: `86400` (24 horas)

## Exemplo de Arquivo .env

```bash
# Ambiente
ENV=dev

# Servidor
API_PORT=8080
API_READ_TIMEOUT=10s
API_WRITE_TIMEOUT=10s
API_IDLE_TIMEOUT=120s
API_BODY_LIMIT=10485760

# Banco de Dados
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=gestao_financeira
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h
JWT_ISSUER=gestao-financeira-api

# Redis (opcional)
REDIS_URL=
REDIS_TTL=5m

# Logging
LOG_LEVEL=info
LOG_FORMAT=console

# Migrations
MIGRATIONS_PATH=file://migrations

# CORS
ALLOWED_ORIGINS=http://localhost:3000
CORS_MAX_AGE=86400
```

## Configuração para Produção

Em produção, certifique-se de:

1. **ENV=production** ou **ENV=prod**
2. **JWT_SECRET** forte e aleatório
3. **POSTGRES_PASSWORD** seguro
4. **LOG_FORMAT=json** para logs estruturados
5. **POSTGRES_SSLMODE=require** ou superior
6. **ALLOWED_ORIGINS** configurado corretamente para seu domínio
7. Todas as senhas e secrets devem ser diferentes dos padrões

## Validação

A configuração é validada no startup da aplicação. Se alguma configuração obrigatória estiver faltando ou inválida, a aplicação não iniciará e mostrará uma mensagem de erro clara.

## Notas

- Todos os valores de duração suportam: `s` (segundos), `m` (minutos), `h` (horas), `d` (dias)
- Valores numéricos são interpretados como inteiros
- Valores booleanos são interpretados como strings (`true`/`false`)
- Nunca commite arquivos `.env` no controle de versão

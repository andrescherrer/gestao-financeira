# Segurança - Gestão Financeira API

Este documento descreve as medidas de segurança implementadas na API.

## Proteção contra SQL Injection

### GORM e Prepared Statements

A API utiliza **GORM** como ORM, que automaticamente protege contra SQL Injection através de **Prepared Statements**.

#### Como Funciona

1. **Prepared Statements**: GORM usa prepared statements por padrão para todas as queries
2. **Parâmetros Escapados**: Todos os parâmetros são automaticamente escapados antes de serem inseridos nas queries
3. **Type Safety**: O uso de structs e value objects garante type safety

#### Exemplo de Proteção

```go
// ✅ SEGURO: GORM usa prepared statements automaticamente
db.Where("user_id = ?", userID).Find(&accounts)

// ✅ SEGURO: Parâmetros são escapados
db.Where("name = ?", accountName).First(&account)

// ✅ SEGURO: Value objects garantem type safety
accountID := valueobjects.NewAccountID(id)
db.Where("id = ?", accountID.Value()).First(&account)
```

#### Validações Adicionais

Além da proteção do GORM, a API implementa:

1. **Validação de UUIDs**: Todos os IDs são validados como UUIDs antes de serem usados em queries
2. **Validação de Input**: DTOs são validados usando `go-playground/validator`
3. **Value Objects**: Uso de value objects garante que apenas valores válidos chegam ao banco

#### Recomendações

- ✅ **Sempre use** prepared statements (GORM faz isso automaticamente)
- ✅ **Nunca** construa queries SQL manualmente com concatenação de strings
- ✅ **Sempre** valide inputs antes de usar em queries
- ✅ **Use** value objects para garantir type safety

#### Verificação

Para verificar se uma query está usando prepared statements:

```go
// GORM loga queries com placeholders (?)
// Exemplo de log:
// SELECT * FROM accounts WHERE user_id = $1 AND name = $2
// Parâmetros: [user-id-uuid, "Account Name"]
```

## Proteção contra XSS (Cross-Site Scripting)

### Backend

O backend não renderiza HTML, então não há risco de XSS no servidor. Todas as respostas são JSON.

### Frontend

O frontend Vue 3 escapa automaticamente conteúdo em templates:

```vue
<!-- ✅ SEGURO: Vue escapa automaticamente -->
<template>
  <div>{{ userInput }}</div>
</template>
```

**⚠️ ATENÇÃO**: Se precisar renderizar HTML, use sanitização:

```vue
<!-- ❌ PERIGOSO: Não use v-html sem sanitização -->
<div v-html="userInput"></div>

<!-- ✅ SEGURO: Use sanitização -->
<div v-html="sanitizeHtml(userInput)"></div>
```

## Rate Limiting

### Configuração Global

A API implementa rate limiting global configurável:

- **Padrão**: 100 requisições por minuto por IP
- **Configurável**: Via variáveis de ambiente
- **Headers**: Inclui headers `X-RateLimit-*` nas respostas

### Rate Limiting por Endpoint

Endpoints críticos têm limites mais restritivos:

- **Auth endpoints** (`/auth/login`, `/auth/register`): 10 req/min
- **Write operations** (POST, PUT, DELETE): 30 req/min
- **Read operations** (GET): 100 req/min

Ver `backend/pkg/middleware/ratelimit.go` para detalhes.

## Autenticação e Autorização

### JWT Tokens

- Tokens são assinados com HS256
- Expiração configurável (padrão: 24h)
- Validação de assinatura e expiração
- Verificação de existência de usuário no banco

### Middleware de Autenticação

- Valida token JWT
- Verifica se usuário existe e está ativo
- Cache de verificação (30s TTL) para performance

## Validação de Input

### Múltiplas Camadas

1. **Frontend**: Validação com Zod schemas
2. **Backend DTOs**: Validação com `go-playground/validator`
3. **Domain**: Validação em value objects

### Tipos de Validação

- **UUIDs**: Todos os IDs são validados como UUIDs
- **Emails**: Validação de formato de email
- **Valores Monetários**: Validação de valores positivos
- **Enums**: Validação de valores permitidos

## CORS

### Configuração

- Origens permitidas configuráveis
- Credenciais habilitadas apenas para origens permitidas
- Headers permitidos: `Origin`, `Content-Type`, `Accept`, `Authorization`

## Headers de Segurança

A API inclui os seguintes headers de segurança:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`

## Logging e Monitoramento

### Logs Estruturados

- Logs estruturados com zerolog
- Request ID em todas as requisições
- Níveis de log configuráveis

### Informações Sensíveis

- Senhas nunca são logadas
- Tokens JWT são logados apenas parcialmente (primeiros caracteres)
- Dados sensíveis são mascarados nos logs

## Recomendações para Produção

1. ✅ Use HTTPS em produção
2. ✅ Configure CORS corretamente
3. ✅ Use secrets fortes para JWT
4. ✅ Configure rate limiting apropriado
5. ✅ Monitore logs para atividades suspeitas
6. ✅ Mantenha dependências atualizadas
7. ✅ Use prepared statements (já implementado via GORM)

---

**Última atualização**: 2025-12-29

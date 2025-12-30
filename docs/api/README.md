# Documentação da API - Gestão Financeira

Esta pasta contém a documentação completa da API e recursos para facilitar a integração.

## 📚 Conteúdo

- **Swagger UI**: Documentação interativa disponível em `http://localhost:8080/swagger/index.html`
- **Postman Collection**: Collection completa com todos os endpoints
- **Postman Environment**: Variáveis de ambiente para facilitar testes

## 🚀 Como Usar

### Swagger UI

1. Inicie a API: `docker-compose up` ou `go run cmd/api/main.go`
2. Acesse: `http://localhost:8080/swagger/index.html`
3. Explore os endpoints interativamente
4. Teste requisições diretamente no navegador

### Postman Collection

1. **Importar Collection e Environment:**
   - Abra o Postman
   - Clique em "Import"
   - Selecione `Gestao_Financeira_API.postman_collection.json`
   - Selecione `Gestao_Financeira_API.postman_environment.json`
   - Selecione o environment "Gestão Financeira API - Local"

2. **Configurar Variáveis:**
   - A variável `base_url` já está configurada para `http://localhost:8080`
   - Para outros ambientes, altere o valor da variável

3. **Autenticação Automática:**
   - Execute a requisição "Login" na pasta "Authentication"
   - O token será automaticamente salvo na variável `api_token`
   - Todas as requisições protegidas usarão automaticamente este token

4. **Fluxo Recomendado:**
   ```
   1. Register User (ou Login se já tiver conta)
   2. Create Account (account_id será salvo automaticamente)
   3. Create Category (category_id será salvo automaticamente)
   4. Create Transaction (transaction_id será salvo automaticamente)
   5. Create Budget (budget_id será salvo automaticamente)
   6. Explorar outros endpoints
   ```

## 🔐 Autenticação

A API utiliza autenticação JWT (JSON Web Tokens):

1. **Obter Token:**
   ```http
   POST /api/v1/auth/login
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "SecurePass123"
   }
   ```

2. **Usar Token:**
   ```http
   Authorization: Bearer <token>
   ```

3. **Token Expira:**
   - Padrão: 24 horas
   - Configurável via variável de ambiente `JWT_EXPIRATION`

## 📋 Endpoints Disponíveis

### Públicos
- `POST /api/v1/auth/register` - Registrar novo usuário
- `POST /api/v1/auth/login` - Fazer login

### Protegidos (requerem autenticação)

#### Accounts
- `POST /api/v1/accounts` - Criar conta
- `GET /api/v1/accounts` - Listar contas (com paginação)
- `GET /api/v1/accounts/:id` - Obter conta por ID

#### Transactions
- `POST /api/v1/transactions` - Criar transação
- `GET /api/v1/transactions` - Listar transações (com filtros e paginação)
- `GET /api/v1/transactions/:id` - Obter transação por ID
- `PUT /api/v1/transactions/:id` - Atualizar transação
- `DELETE /api/v1/transactions/:id` - Deletar transação (soft delete)
- `POST /api/v1/transactions/:id/restore` - Restaurar transação deletada

#### Categories
- `POST /api/v1/categories` - Criar categoria
- `GET /api/v1/categories` - Listar categorias (com paginação)
- `GET /api/v1/categories/:id` - Obter categoria por ID
- `PUT /api/v1/categories/:id` - Atualizar categoria
- `DELETE /api/v1/categories/:id` - Deletar categoria (soft delete)
- `POST /api/v1/categories/:id/restore` - Restaurar categoria deletada

#### Budgets
- `POST /api/v1/budgets` - Criar orçamento
- `GET /api/v1/budgets` - Listar orçamentos (com filtros e paginação)
- `GET /api/v1/budgets/:id` - Obter orçamento por ID
- `GET /api/v1/budgets/:id/progress` - Obter progresso do orçamento
- `PUT /api/v1/budgets/:id` - Atualizar orçamento
- `DELETE /api/v1/budgets/:id` - Deletar orçamento

#### Reports
- `GET /api/v1/reports/monthly` - Relatório mensal
- `GET /api/v1/reports/annual` - Relatório anual
- `GET /api/v1/reports/category` - Relatório por categoria
- `GET /api/v1/reports/income-vs-expense` - Comparação receitas vs despesas

## 📝 Exemplos de Requisições

### Criar Conta

```http
POST /api/v1/accounts
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Conta Corrente",
  "type": "BANK",
  "initial_balance": "1000.00",
  "currency": "BRL",
  "context": "PERSONAL"
}
```

### Criar Transação

```http
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "INCOME",
  "amount": "150.50",
  "currency": "BRL",
  "description": "Salário",
  "date": "2025-12-30"
}
```

## 🔢 Códigos de Resposta HTTP

- `200 OK` - Operação bem-sucedida
- `201 Created` - Recurso criado com sucesso
- `400 Bad Request` - Dados inválidos ou validação falhou
- `401 Unauthorized` - Token ausente ou inválido
- `403 Forbidden` - Acesso negado (sem permissão)
- `404 Not Found` - Recurso não encontrado
- `409 Conflict` - Conflito (ex: recurso já existe)
- `422 Unprocessable Entity` - Erro de validação de domínio
- `429 Too Many Requests` - Rate limit excedido
- `500 Internal Server Error` - Erro interno do servidor

## ⚡ Rate Limiting

A API implementa rate limiting para proteger contra abuso:

- **Limite padrão**: 100 requisições por minuto por IP
- **Headers de resposta**:
  - `X-RateLimit-Limit`: Limite total
  - `X-RateLimit-Remaining`: Requisições restantes
  - `X-RateLimit-Reset`: Timestamp de reset

## 📄 Paginação

Endpoints de listagem suportam paginação:

- `page`: Número da página (1-based, padrão: 1)
- `limit`: Itens por página (padrão: 10, máximo: 100)

**Exemplo:**
```
GET /api/v1/transactions?page=2&limit=20
```

**Resposta:**
```json
{
  "message": "Transactions retrieved successfully",
  "data": {
    "transactions": [...],
    "count": 20,
    "pagination": {
      "page": 2,
      "limit": 20,
      "total": 45,
      "total_pages": 3,
      "has_next": true,
      "has_prev": true
    }
  }
}
```

## 🛠️ Características Técnicas

- **Arquitetura**: Domain-Driven Design (DDD) com Clean Architecture
- **Atomicidade**: Operações críticas garantidas por Unit of Work pattern
- **Soft Delete**: Exclusão lógica com possibilidade de restauração
- **Validação**: Validação em múltiplas camadas (frontend, backend, domain)
- **Tratamento de Erros**: Erros tipados e consistentes em toda a API

## 📞 Suporte

- **Email**: support@gestaofinanceira.com
- **GitHub**: https://github.com/gestao-financeira
- **Documentação Swagger**: http://localhost:8080/swagger/index.html

## 📜 Licença

Apache 2.0


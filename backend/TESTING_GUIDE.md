# Guia de Testes - Sistema de Gestão Financeira

## 📋 Pré-requisitos

### Obrigatórios

1. **Go 1.21+**
   ```bash
   go version
   # Deve mostrar: go version go1.21.x ou superior
   ```

2. **Docker e Docker Compose**
   ```bash
   docker --version
   docker compose version
   ```

3. **Git**
   ```bash
   git --version
   ```

### Opcionais (mas recomendados)

4. **Make** (para facilitar comandos)
   ```bash
   make --version
   ```

---

## 🚀 Como Testar

### 1. Testes Unitários

#### Executar todos os testes do Identity Context:
```bash
cd backend
go test ./internal/identity/... -v
```

#### Executar testes de um componente específico:
```bash
# Value Objects
go test ./internal/identity/domain/valueobjects/... -v

# Entities
go test ./internal/identity/domain/entities/... -v

# Use Cases
go test ./internal/identity/application/usecases/... -v

# Services
go test ./internal/identity/infrastructure/services/... -v

# Handlers
go test ./internal/identity/presentation/handlers/... -v

# Middleware
go test ./pkg/middleware/... -v
```

#### Ver cobertura de testes:
```bash
# Cobertura geral
go test ./internal/identity/... -cover

# Cobertura detalhada
go test ./internal/identity/... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Relatório HTML (abre no navegador)
go tool cover -html=coverage.out
```

#### Executar testes com race detector:
```bash
go test ./internal/identity/... -race
```

---

### 2. Testar a API (Manual)

#### Passo 1: Iniciar os serviços
```bash
# Na raiz do projeto
docker compose up -d

# Verificar se os serviços estão rodando
docker compose ps
```

#### Passo 2: Verificar health check
```bash
curl http://localhost:8080/health
```

#### Passo 3: Testar registro de usuário
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "senha123456",
    "first_name": "João",
    "last_name": "Silva"
  }'
```

#### Passo 4: Testar login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "senha123456"
  }'
```

#### Passo 5: Testar rota protegida (com token JWT)
```bash
# Primeiro, faça login e copie o token da resposta
TOKEN="seu-token-jwt-aqui"

curl -X GET http://localhost:8080/api/v1/ \
  -H "Authorization: Bearer $TOKEN"
```

---

### 3. Testes de Integração (quando implementados)

#### Com banco de dados em memória (SQLite):
```bash
# Os testes de integração usarão SQLite em memória
# Não requerem Docker rodando
go test ./internal/identity/infrastructure/persistence/... -v
```

---

## 🛠️ Comandos Úteis

### Verificar se tudo está configurado:
```bash
# Verificar Go
go version

# Verificar dependências
cd backend
go mod verify

# Baixar dependências
go mod download

# Verificar formatação
gofmt -l ./internal/identity/...

# Verificar build
go build ./cmd/api/...
```

### Executar a aplicação localmente:
```bash
cd backend

# Copiar variáveis de ambiente
cp ../env.example .env

# Editar .env com suas configurações (se necessário)

# Executar
go run cmd/api/main.go
```

### Executar com Docker:
```bash
# Na raiz do projeto
docker compose up --build
```

---

## 📊 Verificar Cobertura de Testes

### Cobertura por componente:
```bash
cd backend

# Cobertura geral
go test ./internal/identity/... -cover

# Cobertura detalhada
go test ./internal/identity/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep identity

# Gerar relatório HTML
go tool cover -html=coverage.out -o coverage.html
# Abrir coverage.html no navegador
```

### Cobertura esperada:
- Value Objects: ~89%
- Services: ~89%
- Handlers: ~88%
- Use Cases: ~87%
- Entities: ~81%
- **Total: ~75%** (sem testes de integração)

---

## 🧪 Exemplos de Testes

### Teste de registro bem-sucedido:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "novo@example.com",
    "password": "senha123456",
    "first_name": "Maria",
    "last_name": "Santos"
  }'
```

**Resposta esperada (201):**
```json
{
  "message": "User registered successfully",
  "data": {
    "user_id": "uuid-gerado",
    "email": "novo@example.com",
    "first_name": "Maria",
    "last_name": "Santos",
    "full_name": "Maria Santos"
  }
}
```

### Teste de login bem-sucedido:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "novo@example.com",
    "password": "senha123456"
  }'
```

**Resposta esperada (200):**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": "uuid",
    "email": "novo@example.com",
    "first_name": "Maria",
    "last_name": "Santos",
    "full_name": "Maria Santos",
    "expires_in": 86400
  }
}
```

### Teste de erro (email já existe):
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "novo@example.com",
    "password": "senha123456",
    "first_name": "Outro",
    "last_name": "Usuario"
  }'
```

**Resposta esperada (409):**
```json
{
  "error": "User with this email already exists",
  "code": 409
}
```

---

## 🔍 Troubleshooting

### Problema: "go: command not found"
**Solução:** Instalar Go
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Ou baixar de: https://golang.org/dl/
```

### Problema: "docker: command not found"
**Solução:** Instalar Docker
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker
```

### Problema: "connection refused" ao testar API
**Solução:** Verificar se os serviços estão rodando
```bash
docker compose ps
docker compose logs api
docker compose logs postgres
```

### Problema: "database connection failed"
**Solução:** Verificar variáveis de ambiente
```bash
# Verificar .env
cat .env

# Verificar se PostgreSQL está saudável
docker compose exec postgres pg_isready -U postgres
```

### Problema: Testes falhando
**Solução:** Verificar dependências
```bash
cd backend
go mod tidy
go mod download
go test ./... -v
```

---

## 📝 Checklist de Testes

### Antes de fazer commit:
- [ ] Todos os testes unitários passando
- [ ] Cobertura acima de 70%
- [ ] Código formatado (`gofmt`)
- [ ] Build sem erros
- [ ] API responde corretamente

### Antes de fazer deploy:
- [ ] Todos os testes passando
- [ ] Testes de integração passando
- [ ] Cobertura acima de 80%
- [ ] Health checks funcionando
- [ ] Testes E2E básicos passando

---

## 🎯 Próximos Passos

1. **Executar testes unitários** (já disponíveis)
2. **Testar API manualmente** (com Docker)
3. **Aguardar Sprint 2.7** para testes de integração
4. **Implementar testes E2E** (Sprint 4.4)

---

## 📚 Recursos Adicionais

- [Documentação Go Testing](https://pkg.go.dev/testing)
- [Documentação Fiber](https://docs.gofiber.io/)
- [Documentação GORM](https://gorm.io/docs/)


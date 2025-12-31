# 🚀 Guia Rápido de Testes

## ✅ Verificação Rápida de Pré-requisitos

Execute estes comandos para verificar se tudo está instalado:

```bash
# Verificar Go
go version
# Deve mostrar: go version go1.21.x ou superior

# Verificar Docker
docker --version
docker compose version

# Verificar Git
git --version
```

---

## 🧪 Testes Unitários (Mais Rápido)

### Opção 1: Com Make (se disponível)
```bash
cd backend
make test-identity
```

### Opção 2: Com Go diretamente
```bash
cd backend

# Todos os testes do Identity Context
go test ./internal/identity/... -v

# Ver cobertura
go test ./internal/identity/... -cover
```

**Tempo estimado:** 1-2 segundos  
**Não requer:** Docker ou banco de dados

---

## 🌐 Testar API (Requer Docker)

### Passo 1: Iniciar serviços
```bash
# Na raiz do projeto
docker compose up -d

# Verificar status
docker compose ps
```

### Passo 2: Testar endpoints

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Registrar usuário
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

#### Fazer login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "senha123456"
  }'
```

**Tempo estimado:** 30-60 segundos (primeira vez)  
**Requer:** Docker rodando

---

## 📊 Ver Cobertura de Testes

```bash
cd backend

# Cobertura simples
go test ./internal/identity/... -cover

# Cobertura detalhada
go test ./internal/identity/... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Relatório visual (HTML)
go tool cover -html=coverage.out
# Abre no navegador automaticamente
```

---

## 🛠️ Comandos Úteis

### Verificar se código compila
```bash
cd backend
go build ./cmd/api/...
```

### Verificar formatação
```bash
cd backend
gofmt -l ./internal/identity/...
```

### Limpar arquivos de teste
```bash
cd backend
rm -f coverage.out coverage.html
```

---

## ⚡ Testes Mais Comuns

### Testar apenas value objects
```bash
cd backend
go test ./internal/identity/domain/valueobjects/... -v
```

### Testar apenas use cases
```bash
cd backend
go test ./internal/identity/application/usecases/... -v
```

### Testar apenas handlers
```bash
cd backend
go test ./internal/identity/presentation/handlers/... -v
```

---

## 🐛 Problemas Comuns

### "go: command not found"
**Solução:** Instalar Go
- Ubuntu/Debian: `sudo apt install golang-go`
- Ou baixar: https://golang.org/dl/

### "docker: command not found"
**Solução:** Instalar Docker
- Ubuntu/Debian: `sudo apt install docker.io docker-compose`

### "connection refused" na API
**Solução:** Verificar se Docker está rodando
```bash
docker compose ps
docker compose logs api
```

### Testes falhando
**Solução:** Atualizar dependências
```bash
cd backend
go mod tidy
go mod download
```

---

## 📝 Checklist Rápido

Antes de fazer commit:
- [ ] `go test ./internal/identity/... -v` ✅
- [ ] `go build ./cmd/api/...` ✅
- [ ] `gofmt -l ./internal/identity/...` (sem saída = OK)

---

## 🎯 Próximos Passos

1. ✅ **Agora:** Testes unitários (já funcionam)
2. ⏳ **Sprint 2.7:** Testes de integração (planejado)
3. ⏳ **Sprint 4.4:** Testes E2E (planejado)


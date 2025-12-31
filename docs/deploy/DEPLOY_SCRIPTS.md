# Guia de Deploy em Produção

Este diretório contém scripts e configurações para deploy em produção.

## 📋 Pré-requisitos

1. Docker e Docker Compose instalados
2. Acesso ao servidor de produção
3. Domínio configurado (opcional, mas recomendado)
4. Certificados SSL (para HTTPS)

## 🚀 Deploy Rápido

### 1. Preparar Ambiente

```bash
# Copiar arquivo de exemplo
cp .env.production.example .env.production

# Editar variáveis de ambiente
nano .env.production
```

### 2. Configurar Variáveis Críticas

**IMPORTANTE:** Configure pelo menos estas variáveis:

- `POSTGRES_PASSWORD` - Senha forte do PostgreSQL
- `JWT_SECRET` - Chave secreta JWT (gere com `openssl rand -base64 32`)
- `REDIS_PASSWORD` - Senha do Redis
- `GRAFANA_ADMIN_PASSWORD` - Senha do Grafana
- `GRAFANA_SECRET_KEY` - Chave secreta do Grafana
- `ALLOWED_ORIGINS` - Domínios permitidos para CORS

### 3. Build das Imagens

```bash
# Build da API
docker build -t gestao-financeira/api:latest ./backend

# Build do Frontend
docker build -t gestao-financeira/frontend:latest ./frontend
```

### 4. Deploy

```bash
# Iniciar serviços
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# Verificar logs
docker-compose -f docker-compose.prod.yml logs -f

# Verificar status
docker-compose -f docker-compose.prod.yml ps
```

## 📝 Scripts de Deploy

### deploy.sh

Script automatizado de deploy:

```bash
chmod +x deploy/deploy.sh
./deploy/deploy.sh
```

### backup.sh

Script para criar backup antes do deploy:

```bash
chmod +x deploy/backup.sh
./deploy/backup.sh
```

## 🔒 Segurança

### Checklist de Segurança

- [ ] Todas as senhas foram alteradas dos valores padrão
- [ ] JWT_SECRET é forte e aleatório
- [ ] POSTGRES_SSLMODE está configurado (require ou verify-full)
- [ ] Redis está protegido com senha
- [ ] CORS está configurado apenas para domínios permitidos
- [ ] Firewall configurado para permitir apenas portas necessárias
- [ ] Certificados SSL configurados (HTTPS)
- [ ] Logs não contêm informações sensíveis
- [ ] Backups automáticos configurados

### Geração de Secrets

```bash
# JWT Secret
openssl rand -base64 32

# Redis Password
openssl rand -base64 24

# Grafana Secret Key
openssl rand -base64 32
```

## 📊 Monitoramento

Após o deploy, acesse:

- **API Health Check:** `http://your-server:8080/health`
- **API Swagger:** `http://your-server:8080/swagger/index.html`
- **Prometheus:** `http://your-server:9090`
- **Grafana:** `http://your-server:3001` (admin/admin - altere!)

## 🔄 Atualização

Para atualizar a aplicação:

```bash
# 1. Fazer backup
./deploy/backup.sh

# 2. Pull das novas imagens
docker-compose -f docker-compose.prod.yml pull

# 3. Recriar containers
docker-compose -f docker-compose.prod.yml up -d --force-recreate

# 4. Verificar logs
docker-compose -f docker-compose.prod.yml logs -f api
```

## 🛠️ Troubleshooting

### Verificar Logs

```bash
# Logs da API
docker-compose -f docker-compose.prod.yml logs api

# Logs do PostgreSQL
docker-compose -f docker-compose.prod.yml logs postgres

# Logs de todos os serviços
docker-compose -f docker-compose.prod.yml logs
```

### Verificar Saúde dos Serviços

```bash
# Status dos containers
docker-compose -f docker-compose.prod.yml ps

# Health check da API
curl http://localhost:8080/health

# Health check detalhado
curl http://localhost:8080/health/detailed
```

### Reiniciar Serviços

```bash
# Reiniciar API
docker-compose -f docker-compose.prod.yml restart api

# Reiniciar todos os serviços
docker-compose -f docker-compose.prod.yml restart
```

## 📚 Documentação Adicional

- [Configuração da API](../configuracao/CONFIG.md)
- [Guia de Backup](../docs/tarefas_concluidas/20251231_123000_OPT-002.md)
- [Health Check Avançado](../docs/tarefas_concluidas/20251231_072100_HEALTH-001.md)


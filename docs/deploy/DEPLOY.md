# Guia Completo de Deploy

Este documento fornece um guia completo para deploy do sistema de Gestão Financeira em diferentes ambientes.

## 📋 Índice

1. [Pré-requisitos](#pré-requisitos)
2. [Deploy em Produção](#deploy-em-produção)
3. [Deploy em Staging](#deploy-em-staging)
4. [Deploy em Desenvolvimento](#deploy-em-desenvolvimento)
5. [Configuração de Variáveis](#configuração-de-variáveis)
6. [Segurança](#segurança)
7. [Monitoramento](#monitoramento)
8. [Troubleshooting](#troubleshooting)
9. [Manutenção](#manutenção)

---

## Pré-requisitos

### Software Necessário

- **Docker** 20.10+ e **Docker Compose** 2.0+
- **Git** para clonar o repositório
- **OpenSSL** para gerar secrets (opcional, mas recomendado)
- Acesso SSH ao servidor (para deploy remoto)

### Requisitos do Sistema

- **CPU:** Mínimo 2 cores, recomendado 4+
- **RAM:** Mínimo 4GB, recomendado 8GB+
- **Disco:** Mínimo 20GB livre, recomendado 50GB+
- **Rede:** Portas 80, 443, 8080, 5432, 6379, 9090, 3001 disponíveis

### Verificação

```bash
# Verificar Docker
docker --version
docker-compose --version

# Verificar recursos
free -h
df -h
nproc
```

---

## Deploy em Produção

### 1. Preparação

```bash
# Clonar repositório
git clone <repository-url>
cd gestao-financeira

# Copiar arquivo de exemplo
cp deploy/env.production.example .env.production

# Editar variáveis
nano .env.production
```

### 2. Configuração de Variáveis Críticas

**IMPORTANTE:** Configure estas variáveis antes do deploy:

```bash
# Gerar secrets
JWT_SECRET=$(openssl rand -base64 32)
REDIS_PASSWORD=$(openssl rand -base64 24)
GRAFANA_SECRET_KEY=$(openssl rand -base64 32)

# Adicionar ao .env.production
echo "JWT_SECRET=$JWT_SECRET" >> .env.production
echo "REDIS_PASSWORD=$REDIS_PASSWORD" >> .env.production
echo "GRAFANA_SECRET_KEY=$GRAFANA_SECRET_KEY" >> .env.production
```

### 3. Deploy Automatizado

```bash
# Executar script de deploy
./deploy/deploy.sh
```

O script irá:
1. Validar variáveis críticas
2. Criar backup automático
3. Build das imagens Docker
4. Iniciar serviços
5. Verificar saúde dos serviços

### 4. Deploy Manual

```bash
# Build das imagens
docker build -t gestao-financeira/api:latest ./backend
docker build -t gestao-financeira/frontend:latest ./frontend

# Iniciar serviços
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# Verificar status
docker-compose -f docker-compose.prod.yml ps

# Ver logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 5. Verificação Pós-Deploy

```bash
# Health check da API
curl http://localhost:8080/health

# Health check detalhado
curl http://localhost:8080/health/detailed

# Verificar frontend
curl http://localhost:80

# Verificar Prometheus
curl http://localhost:9090/-/healthy

# Verificar Grafana
curl http://localhost:3001/api/health
```

---

## Deploy em Staging

O deploy em staging segue o mesmo processo de produção, mas com configurações menos restritivas:

```bash
# Usar docker-compose.yml padrão
docker-compose --env-file .env.staging up -d
```

**Diferenças do Staging:**
- SSL pode ser desabilitado
- Logs em formato console (mais legível)
- Recursos menores
- Sem limitações de CORS rígidas

---

## Deploy em Desenvolvimento

Para desenvolvimento local:

```bash
# Iniciar serviços
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar serviços
docker-compose down
```

**Características:**
- Hot reload habilitado
- Logs detalhados
- Sem limitações de recursos
- Banco de dados em memória (opcional)

---

## Configuração de Variáveis

### Variáveis Obrigatórias

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `JWT_SECRET` | Chave secreta JWT | `openssl rand -base64 32` |
| `POSTGRES_PASSWORD` | Senha do PostgreSQL | `StrongPassword123!` |
| `REDIS_PASSWORD` | Senha do Redis | `RedisPassword456!` |
| `ALLOWED_ORIGINS` | Domínios permitidos | `https://app.example.com` |

### Variáveis Opcionais

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `API_PORT` | Porta da API | `8080` |
| `FRONTEND_PORT` | Porta do Frontend | `80` |
| `LOG_LEVEL` | Nível de log | `info` |
| `POSTGRES_MAX_OPEN_CONNS` | Max conexões DB | `50` |

### Geração de Secrets

```bash
# JWT Secret (32 bytes)
openssl rand -base64 32

# Redis Password (24 bytes)
openssl rand -base64 24

# Grafana Secret Key (32 bytes)
openssl rand -base64 32
```

---

## Segurança

### Checklist de Segurança

- [ ] Todas as senhas alteradas dos padrões
- [ ] JWT_SECRET forte e aleatório
- [ ] POSTGRES_SSLMODE=require ou superior
- [ ] Redis protegido com senha
- [ ] CORS configurado apenas para domínios permitidos
- [ ] Firewall configurado
- [ ] Certificados SSL configurados (HTTPS)
- [ ] Logs não contêm informações sensíveis
- [ ] Backups automáticos configurados
- [ ] Acesso SSH restrito

### Configuração de Firewall

```bash
# Permitir apenas portas necessárias
ufw allow 22/tcp    # SSH
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw enable
```

### SSL/TLS

Para produção, configure SSL/TLS usando:

- **Nginx** como reverse proxy
- **Let's Encrypt** para certificados gratuitos
- **Certbot** para gerenciamento automático

Exemplo de configuração Nginx:

```nginx
server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Monitoramento

### Endpoints de Monitoramento

- **API Health:** `http://localhost:8080/health`
- **API Health Detalhado:** `http://localhost:8080/health/detailed`
- **Prometheus:** `http://localhost:9090`
- **Grafana:** `http://localhost:3001`

### Métricas Disponíveis

- Requisições HTTP (total, por status, por endpoint)
- Tempo de resposta
- Uso de memória e CPU
- Conexões do banco de dados
- Uso de cache Redis

### Alertas Recomendados

- API não responde (health check falha)
- Uso de disco > 85%
- Uso de memória > 90%
- Erros HTTP > 5%
- Tempo de resposta > 1s

---

## Troubleshooting

### Problemas Comuns

#### API não inicia

```bash
# Verificar logs
docker-compose -f docker-compose.prod.yml logs api

# Verificar variáveis
docker-compose -f docker-compose.prod.yml config

# Verificar conectividade do banco
docker exec gestao-financeira-db-prod pg_isready -U postgres
```

#### Banco de dados não conecta

```bash
# Verificar se está rodando
docker-compose -f docker-compose.prod.yml ps postgres

# Verificar logs
docker-compose -f docker-compose.prod.yml logs postgres

# Testar conexão
docker exec -it gestao-financeira-db-prod psql -U postgres -d gestao_financeira
```

#### Frontend não carrega

```bash
# Verificar logs
docker-compose -f docker-compose.prod.yml logs frontend

# Verificar se API está acessível
curl http://localhost:8080/health

# Verificar variável VITE_API_URL
docker exec gestao-financeira-frontend-prod env | grep VITE
```

### Comandos Úteis

```bash
# Reiniciar serviço específico
docker-compose -f docker-compose.prod.yml restart api

# Ver uso de recursos
docker stats

# Limpar volumes não utilizados
docker volume prune

# Ver logs em tempo real
docker-compose -f docker-compose.prod.yml logs -f --tail=100
```

---

## Manutenção

### Atualização da Aplicação

```bash
# 1. Fazer backup
./deploy/backup.sh

# 2. Pull das novas imagens
docker-compose -f docker-compose.prod.yml pull

# 3. Recriar containers
docker-compose -f docker-compose.prod.yml up -d --force-recreate

# 4. Verificar saúde
curl http://localhost:8080/health
```

### Backup Regular

Configure backup automático via cron:

```bash
# Adicionar ao crontab
0 2 * * * /path/to/deploy/backup.sh
```

### Limpeza de Logs

```bash
# Limpar logs antigos
docker system prune -f

# Limpar volumes não utilizados
docker volume prune -f
```

### Atualização do Sistema

```bash
# Atualizar Docker
sudo apt update && sudo apt upgrade docker.io docker-compose

# Reiniciar serviços
docker-compose -f docker-compose.prod.yml restart
```

---

## Recursos Adicionais

- [Documentação da API](../configuracao/CONFIG.md)
- [Guia de Backup](../docs/tarefas_concluidas/20251231_123000_OPT-002.md)
- [Health Check Avançado](../docs/tarefas_concluidas/20251231_072100_HEALTH-001.md)
- [Scripts de Deploy](../deploy/README.md)

---

**Última atualização:** 2025-12-31


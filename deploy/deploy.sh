#!/bin/bash

# Script de Deploy em Produção
# Uso: ./deploy/deploy.sh

set -e  # Exit on error

echo "🚀 Iniciando deploy em produção..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar se .env.production existe
if [ ! -f .env.production ]; then
    echo -e "${RED}❌ Arquivo .env.production não encontrado!${NC}"
    echo "Copie .env.production.example para .env.production e configure as variáveis."
    exit 1
fi

# Verificar variáveis críticas
echo -e "${YELLOW}🔍 Verificando variáveis críticas...${NC}"

source .env.production

if [ "$JWT_SECRET" = "CHANGE_ME_GENERATE_STRONG_SECRET_KEY" ] || [ -z "$JWT_SECRET" ]; then
    echo -e "${RED}❌ JWT_SECRET não foi configurado!${NC}"
    exit 1
fi

if [ "$POSTGRES_PASSWORD" = "CHANGE_ME_STRONG_PASSWORD" ] || [ -z "$POSTGRES_PASSWORD" ]; then
    echo -e "${RED}❌ POSTGRES_PASSWORD não foi configurado!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Variáveis críticas verificadas${NC}"

# Criar backup antes do deploy
if [ -f deploy/backup.sh ]; then
    echo -e "${YELLOW}💾 Criando backup antes do deploy...${NC}"
    ./deploy/backup.sh || echo -e "${YELLOW}⚠️  Backup falhou, continuando...${NC}"
fi

# Build das imagens
echo -e "${YELLOW}🔨 Construindo imagens Docker...${NC}"

echo "Building API..."
docker build -t ${DOCKER_REGISTRY:-gestao-financeira}/api:${IMAGE_TAG:-latest} ./backend

echo "Building Frontend..."
docker build -t ${DOCKER_REGISTRY:-gestao-financeira}/frontend:${IMAGE_TAG:-latest} ./frontend

echo -e "${GREEN}✅ Imagens construídas${NC}"

# Parar containers existentes
echo -e "${YELLOW}🛑 Parando containers existentes...${NC}"
docker-compose -f docker-compose.prod.yml --env-file .env.production down || true

# Iniciar serviços
echo -e "${YELLOW}🚀 Iniciando serviços...${NC}"
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# Aguardar serviços iniciarem
echo -e "${YELLOW}⏳ Aguardando serviços iniciarem...${NC}"
sleep 10

# Verificar saúde
echo -e "${YELLOW}🏥 Verificando saúde dos serviços...${NC}"

# Verificar API
if curl -f http://localhost:${API_PORT:-8080}/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ API está saudável${NC}"
else
    echo -e "${RED}❌ API não está respondendo${NC}"
    docker-compose -f docker-compose.prod.yml logs api
    exit 1
fi

# Verificar Frontend
if curl -f http://localhost:${FRONTEND_PORT:-80} > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend está saudável${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend pode não estar pronto ainda${NC}"
fi

echo -e "${GREEN}✅ Deploy concluído com sucesso!${NC}"
echo ""
echo "Serviços disponíveis:"
echo "  - API: http://localhost:${API_PORT:-8080}"
echo "  - Frontend: http://localhost:${FRONTEND_PORT:-80}"
echo "  - Prometheus: http://localhost:${PROMETHEUS_PORT:-9090}"
echo "  - Grafana: http://localhost:${GRAFANA_PORT:-3001}"
echo ""
echo "Para ver logs: docker-compose -f docker-compose.prod.yml logs -f"


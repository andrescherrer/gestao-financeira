#!/bin/bash

# Script de Backup antes do Deploy
# Uso: ./deploy/backup.sh

set -e

echo "💾 Criando backup antes do deploy..."

# Carregar variáveis de ambiente
if [ -f .env.production ]; then
    source .env.production
fi

# Diretório de backups
BACKUP_DIR=${BACKUP_DIR:-./backups}
mkdir -p $BACKUP_DIR

# Timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Backup do banco de dados
if [ -n "$POSTGRES_HOST" ] && [ -n "$POSTGRES_DB" ]; then
    echo "Backing up database..."
    
    BACKUP_FILE="$BACKUP_DIR/db_backup_$TIMESTAMP.sql.gz"
    
    # Verificar se container está rodando
    if docker ps | grep -q gestao-financeira-db-prod; then
        docker exec gestao-financeira-db-prod pg_dump -U ${POSTGRES_USER:-postgres} ${POSTGRES_DB:-gestao_financeira} | gzip > $BACKUP_FILE
        echo "✅ Backup do banco criado: $BACKUP_FILE"
    else
        echo "⚠️  Container do banco não está rodando, pulando backup do banco"
    fi
else
    echo "⚠️  Variáveis do banco não configuradas, pulando backup do banco"
fi

# Backup de volumes (se existirem)
echo "Backing up volumes..."
docker run --rm \
    -v gestao-financeira_postgres_data:/data \
    -v $(pwd)/$BACKUP_DIR:/backup \
    alpine tar czf /backup/volumes_backup_$TIMESTAMP.tar.gz /data 2>/dev/null || echo "⚠️  Não foi possível fazer backup dos volumes"

echo "✅ Backup concluído!"


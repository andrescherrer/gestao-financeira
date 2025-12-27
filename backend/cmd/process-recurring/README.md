# Process Recurring Transactions

Este comando processa transações recorrentes e cria novas instâncias quando necessário.

## Uso

### Via Makefile (Recomendado)

```bash
# Compilar o comando
make build-recurring

# Executar o comando
make run-recurring

# Ou compilar todos os binários
make build-all
```

### Via Go Diretamente

```bash
# Compilar o comando
go build -o bin/process-recurring ./cmd/process-recurring

# Executar o comando
./bin/process-recurring
```

## Variáveis de Ambiente

O comando utiliza as mesmas variáveis de ambiente do servidor principal:

- `DATABASE_URL`: URL de conexão com o banco de dados
- `LOG_LEVEL`: Nível de log (debug, info, warn, error) - padrão: info

## Agendamento com Cron

Para executar diariamente às 00:00, adicione ao crontab:

```bash
# Editar crontab
crontab -e

# Adicionar linha (ajuste o caminho conforme necessário)
0 0 * * * cd /caminho/para/backend && ./bin/process-recurring >> /var/log/recurring-transactions.log 2>&1
```

## Docker

### Via Docker Compose (Recomendado)

O serviço `process-recurring` já está configurado no `docker-compose.yml`:

```bash
# Executar uma vez
docker-compose --profile recurring run process-recurring

# Ou executar no container da API (se o binário estiver disponível)
docker-compose exec api ./bin/process-recurring
```

### Agendamento Automático

Para executar automaticamente, você pode:

1. **Usar cron no host:**
   ```bash
   # Adicionar ao crontab
   0 0 * * * cd /caminho/para/projeto && docker-compose --profile recurring run process-recurring
   ```

2. **Usar Kubernetes CronJob:**
   ```yaml
   apiVersion: batch/v1
   kind: CronJob
   metadata:
     name: process-recurring-transactions
   spec:
     schedule: "0 0 * * *"  # Diariamente às 00:00
     jobTemplate:
       spec:
         template:
           spec:
             containers:
             - name: process-recurring
               image: gestao-financeira-backend:latest
               command: ["./bin/process-recurring"]
             restartPolicy: OnFailure
   ```

3. **Usar um scheduler externo** (ex: GitHub Actions, GitLab CI, etc.)


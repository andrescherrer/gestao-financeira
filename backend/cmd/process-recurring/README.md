# Process Recurring Transactions

Este comando processa transações recorrentes e cria novas instâncias quando necessário.

## Uso

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

Se estiver usando Docker, você pode executar o comando em um container:

```bash
docker-compose exec backend ./bin/process-recurring
```

Ou criar um serviço separado no `docker-compose.yml`:

```yaml
services:
  process-recurring:
    build: ./backend
    command: ./bin/process-recurring
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    depends_on:
      - db
    restart: "no"  # Executa uma vez e sai
```

E agendar com cron no host ou usar um scheduler como Kubernetes CronJob.


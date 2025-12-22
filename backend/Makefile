.PHONY: test test-cover test-race test-identity build run clean help

# Variáveis
GO=go
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## Executa todos os testes
	$(GO) test ./... -v

test-identity: ## Executa testes do Identity Context
	$(GO) test ./internal/identity/... -v

test-cover: ## Executa testes com cobertura
	$(GO) test ./internal/identity/... -cover

test-cover-detail: ## Gera relatório detalhado de cobertura
	$(GO) test ./internal/identity/... -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)
	@echo "\nRelatório HTML gerado em: $(COVERAGE_HTML)"
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

test-race: ## Executa testes com race detector
	$(GO) test ./internal/identity/... -race

test-unit: ## Executa apenas testes unitários (sem integração)
	$(GO) test ./internal/identity/domain/... -v
	$(GO) test ./internal/identity/application/... -v
	$(GO) test ./internal/identity/infrastructure/services/... -v
	$(GO) test ./internal/identity/presentation/handlers/... -v
	$(GO) test ./pkg/middleware/... -v

build: ## Compila a aplicação
	$(GO) build -o bin/api ./cmd/api

run: ## Executa a aplicação
	$(GO) run ./cmd/api/main.go

clean: ## Remove arquivos gerados
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -rf bin/

fmt: ## Formata o código
	$(GO) fmt ./...

vet: ## Executa go vet
	$(GO) vet ./...

lint: fmt vet ## Executa formatação e verificação

mod-tidy: ## Organiza dependências
	$(GO) mod tidy

mod-download: ## Baixa dependências
	$(GO) mod download

deps: mod-download ## Instala todas as dependências

check: lint test ## Verifica código e executa testes


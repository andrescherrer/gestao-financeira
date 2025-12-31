# Planejamento: Upload e Análise Automática de Comprovantes

**Data:** 2025-12-31  
**Feature:** RECEIPT-SCAN-001  
**Status:** 📋 Planejamento  
**Prioridade:** 🟡 Média

---

## 📋 Visão Geral

Implementar funcionalidade para upload de comprovantes (imagens de recibos, notas fiscais, extratos) e análise automática via OCR/IA para extrair informações e criar transações automaticamente, identificando se é receita ou despesa.

---

## 🎯 Objetivos

1. **Upload de Comprovantes**: Permitir upload de imagens (JPG, PNG, PDF)
2. **Análise Automática**: Extrair informações do comprovante via OCR/IA
3. **Criação Automática**: Criar transação automaticamente com dados extraídos
4. **Classificação**: Identificar automaticamente se é receita ou despesa
5. **Validação Manual**: Permitir revisão e correção antes de confirmar

---

## 📊 Análise de Soluções

### Opção 1: Solução Nativa Go (Recomendada para Controle Total)

#### Stack Tecnológico:
- **OCR**: Tesseract OCR (via `gosseract`) ou Google Cloud Vision API
- **ML/IA**: TensorFlow Lite ou integração com APIs de IA
- **Storage**: MinIO (S3-compatible) ou AWS S3
- **Processamento**: Goroutines para processamento assíncrono
- **Queue**: Redis Queue ou RabbitMQ para jobs

#### Vantagens:
- ✅ Controle total sobre o processo
- ✅ Integração nativa com arquitetura DDD existente
- ✅ Sem dependências externas (exceto APIs opcionais)
- ✅ Performance otimizada com Go
- ✅ Custo controlado (pode usar Tesseract gratuito)

#### Desvantagens:
- ⚠️ Desenvolvimento mais complexo
- ⚠️ Manutenção de código de OCR/ML
- ⚠️ Precisão pode ser menor que soluções especializadas

#### Estimativa: 40-60 horas

---

### Opção 2: N8N (Workflow Automation) + APIs Externas

#### Arquitetura:
```
Frontend → Backend Go → N8N Webhook → 
  → OCR API (Google Vision/AWS Textract) → 
  → LLM API (OpenAI/Anthropic) para análise → 
  → N8N Processa → 
  → Webhook de volta para Backend → 
  → Cria Transação
```

#### Stack Tecnológico:
- **N8N**: Workflow automation (self-hosted ou cloud)
- **OCR**: Google Cloud Vision API, AWS Textract, ou Azure Form Recognizer
- **IA/LLM**: OpenAI GPT-4, Anthropic Claude, ou Google Gemini
- **Storage**: MinIO ou S3 (via N8N ou direto do backend)

#### Vantagens:
- ✅ **Visual e Flexível**: Workflows visuais no N8N
- ✅ **Rápido de Implementar**: Menos código, mais configuração
- ✅ **APIs Especializadas**: Maior precisão de OCR/IA
- ✅ **Fácil de Ajustar**: Modificar workflow sem recompilar
- ✅ **Integração com Múltiplas APIs**: Fácil trocar provedores
- ✅ **Processamento Assíncrono Nativo**: N8N gerencia filas
- ✅ **Retry e Error Handling**: N8N tem isso built-in

#### Desvantagens:
- ⚠️ **Dependência Externa**: N8N precisa estar rodando
- ⚠️ **Custo de APIs**: Google Vision, OpenAI, etc. são pagos
- ⚠️ **Complexidade de Deploy**: Mais um serviço para gerenciar
- ⚠️ **Debugging**: Pode ser mais difícil debugar workflows

#### Estimativa: 20-30 horas (mais rápido!)

---

### Opção 3: Híbrida (N8N + Backend Go)

#### Arquitetura:
```
Frontend → Backend Go (upload/storage) → 
  → N8N Webhook (processamento) → 
  → APIs OCR/IA → 
  → N8N analisa e retorna dados → 
  → Backend Go valida e cria transação
```

#### Vantagens:
- ✅ Melhor dos dois mundos
- ✅ Backend controla storage e validação
- ✅ N8N gerencia processamento complexo
- ✅ Fácil escalar processamento

#### Desvantagens:
- ⚠️ Mais complexo de configurar inicialmente

#### Estimativa: 30-40 horas

---

## 🏆 Recomendação: Opção 2 (N8N) ou Opção 3 (Híbrida)

**Por quê?**
1. **Velocidade de Implementação**: N8N acelera muito o desenvolvimento
2. **Precisão**: APIs especializadas (Google Vision, GPT-4) têm melhor precisão
3. **Manutenibilidade**: Workflows visuais são mais fáceis de manter
4. **Flexibilidade**: Fácil ajustar lógica sem recompilar código
5. **Escalabilidade**: N8N gerencia filas e processamento assíncrono

---

## 📐 Arquitetura Detalhada (Opção 2: N8N)

### Fluxo Completo:

```
1. Usuário faz upload no Frontend
   ↓
2. Frontend → POST /api/v1/receipts/upload
   ↓
3. Backend Go:
   - Valida arquivo (tipo, tamanho)
   - Salva em storage (MinIO/S3)
   - Cria Receipt entity (status: PENDING)
   - Retorna receipt_id
   ↓
4. Backend Go → Webhook N8N (receipt_id, file_url)
   ↓
5. N8N Workflow:
   a. Download da imagem
   b. OCR (Google Vision API)
   c. Extração de texto
   d. Análise com LLM (GPT-4):
      - Identifica tipo (receita/despesa)
      - Extrai valor
      - Extrai data
      - Extrai descrição
      - Extrai categoria (se possível)
   e. Retorna dados estruturados
   ↓
6. N8N → Webhook Backend Go (receipt_id, extracted_data)
   ↓
7. Backend Go:
   - Atualiza Receipt (status: PROCESSED, extracted_data)
   - Cria Transaction DRAFT (aguardando confirmação)
   - Notifica usuário (via Notification Context)
   ↓
8. Usuário revisa e confirma/corrige
   ↓
9. Backend Go:
   - Atualiza Transaction (status: CONFIRMED)
   - Aplica transação (atualiza saldo)
```

---

## 🏗️ Estrutura de Implementação

### 1. Receipt Context (Novo Bounded Context)

```
backend/internal/receipt/
├── domain/
│   ├── entities/
│   │   └── receipt.go          # Receipt aggregate root
│   ├── valueobjects/
│   │   ├── receipt_id.go
│   │   ├── receipt_status.go   # PENDING, PROCESSING, PROCESSED, FAILED
│   │   └── file_info.go        # filename, size, mime_type
│   ├── repositories/
│   │   └── receipt_repository.go
│   └── events/
│       ├── receipt_uploaded.go
│       └── receipt_processed.go
├── application/
│   ├── dtos/
│   │   ├── upload_receipt_dto.go
│   │   ├── process_receipt_dto.go
│   │   └── confirm_receipt_dto.go
│   └── usecases/
│       ├── upload_receipt_usecase.go
│       ├── process_receipt_usecase.go
│       └── confirm_receipt_usecase.go
├── infrastructure/
│   ├── persistence/
│   │   ├── receipt_model.go
│   │   └── gorm_receipt_repository.go
│   ├── storage/
│   │   └── file_storage.go      # Interface para MinIO/S3
│   └── services/
│       └── n8n_client.go       # Cliente para N8N webhooks
└── presentation/
    ├── handlers/
    │   └── receipt_handler.go
    └── routes/
        └── receipt_routes.go
```

### 2. Entidade Receipt

```go
type Receipt struct {
    id              ReceiptID
    userID          UserID
    fileName        string
    filePath        string        // Path no storage
    fileSize        int64
    mimeType        string
    status          ReceiptStatus // PENDING, PROCESSING, PROCESSED, FAILED
    extractedData   *ExtractedData // Dados extraídos pelo OCR/IA
    transactionID   *TransactionID // ID da transação criada (se confirmada)
    errorMessage    *string
    createdAt       time.Time
    updatedAt       time.Time
}

type ExtractedData struct {
    Type        string  // INCOME ou EXPENSE
    Amount      float64
    Currency    string
    Date        string
    Description string
    Category    *string
    Account     *string
    Confidence  float64 // 0.0 a 1.0
}
```

### 3. Migration

```sql
CREATE TABLE receipts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL, -- PENDING, PROCESSING, PROCESSED, FAILED
    extracted_data JSONB,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_receipts_user_id ON receipts(user_id);
CREATE INDEX idx_receipts_status ON receipts(status);
CREATE INDEX idx_receipts_created_at ON receipts(created_at);
```

---

## 🔧 Implementação Detalhada

### Fase 1: Storage e Upload (8-12 horas)

#### Backend:
1. **Configurar Storage** (MinIO ou S3)
   - Adicionar MinIO ao docker-compose.yml
   - Criar interface FileStorage
   - Implementar MinIOStorage

2. **Receipt Context - Estrutura Base**
   - Criar entidade Receipt
   - Criar value objects
   - Criar repository interface
   - Implementar GORM repository
   - Criar migration

3. **Upload Endpoint**
   - Handler para upload de arquivo
   - Validação de tipo/tamanho
   - Salvar no storage
   - Criar Receipt entity
   - Retornar receipt_id

#### Frontend:
1. **Componente de Upload**
   - Drag & drop ou file picker
   - Preview da imagem
   - Progress bar
   - Validação client-side

#### Tarefas:
- [ ] RECEIPT-001: Configurar MinIO no docker-compose (2h)
- [ ] RECEIPT-002: Criar interface FileStorage (2h)
- [ ] RECEIPT-003: Implementar MinIOStorage (4h)
- [ ] RECEIPT-004: Criar Receipt Context - estrutura base (4h)
- [ ] RECEIPT-005: Criar migration para receipts (1h)
- [ ] RECEIPT-006: Implementar UploadReceiptUseCase (4h)
- [ ] RECEIPT-007: Criar ReceiptHandler e rotas (2h)
- [ ] RECEIPT-008: Frontend - Componente de upload (4h)
- [ ] RECEIPT-009: Testes de upload (2h)

---

### Fase 2: Integração com N8N (6-8 horas)

#### Backend:
1. **N8N Client**
   - Cliente HTTP para chamar webhooks N8N
   - Retry logic
   - Error handling

2. **Webhook Handler**
   - Endpoint para receber resposta do N8N
   - Validar dados recebidos
   - Atualizar Receipt
   - Criar Transaction DRAFT

#### N8N:
1. **Workflow de Processamento**
   - Trigger: Webhook (recebe receipt_id, file_url)
   - Download da imagem
   - OCR: Google Vision API
   - Análise: OpenAI GPT-4 ou Claude
   - Retorno: Webhook para backend

#### Tarefas:
- [ ] RECEIPT-010: Configurar N8N no docker-compose (2h)
- [ ] RECEIPT-011: Criar N8NClient no backend (2h)
- [ ] RECEIPT-012: Integrar upload com N8N webhook (2h)
- [ ] RECEIPT-013: Criar webhook handler para resposta N8N (2h)
- [ ] RECEIPT-014: Criar workflow N8N (4h)
- [ ] RECEIPT-015: Configurar APIs (Google Vision, OpenAI) (2h)
- [ ] RECEIPT-016: Testes de integração N8N (2h)

---

### Fase 3: Processamento e Análise (8-10 horas)

#### N8N Workflow Detalhado:

**Node 1: Webhook Trigger**
- Recebe: `{ receipt_id, file_url, user_id }`

**Node 2: Download Image**
- Baixa imagem do storage

**Node 3: Google Vision OCR**
- Extrai texto da imagem
- Retorna: `{ text, confidence }`

**Node 4: OpenAI GPT-4 Analysis**
- Prompt:
```
Analise este texto extraído de um comprovante financeiro e retorne JSON:
{
  "type": "INCOME ou EXPENSE",
  "amount": número,
  "currency": "BRL",
  "date": "YYYY-MM-DD",
  "description": "descrição",
  "category": "categoria se identificada",
  "confidence": 0.0 a 1.0
}

Texto: {text}
```

**Node 5: Validação e Formatação**
- Valida dados extraídos
- Formata resposta

**Node 6: Webhook para Backend**
- POST para `/api/v1/receipts/{receipt_id}/process`
- Envia dados extraídos

#### Backend:
1. **ProcessReceiptUseCase**
   - Recebe dados do N8N
   - Valida estrutura
   - Atualiza Receipt (status: PROCESSED)
   - Cria Transaction DRAFT
   - Publica evento ReceiptProcessed

#### Tarefas:
- [ ] RECEIPT-017: Refinar workflow N8N com OCR (2h)
- [ ] RECEIPT-018: Implementar análise com LLM (3h)
- [ ] RECEIPT-019: Criar ProcessReceiptUseCase (3h)
- [ ] RECEIPT-020: Criar Transaction DRAFT (2h)
- [ ] RECEIPT-021: Testes de processamento (2h)

---

### Fase 4: Confirmação e Validação Manual (6-8 horas)

#### Backend:
1. **ConfirmReceiptUseCase**
   - Usuário revisa dados extraídos
   - Pode corrigir valores
   - Confirma transação
   - Transaction DRAFT → CONFIRMED
   - Aplica transação (atualiza saldo)

2. **Endpoints**
   - GET `/api/v1/receipts/{id}` - Ver dados extraídos
   - PUT `/api/v1/receipts/{id}/confirm` - Confirmar transação
   - PUT `/api/v1/receipts/{id}/reject` - Rejeitar

#### Frontend:
1. **Tela de Revisão**
   - Mostra imagem do comprovante
   - Mostra dados extraídos (editáveis)
   - Botões: Confirmar, Rejeitar, Corrigir

#### Tarefas:
- [ ] RECEIPT-022: Criar ConfirmReceiptUseCase (3h)
- [ ] RECEIPT-023: Criar endpoints de confirmação (2h)
- [ ] RECEIPT-024: Frontend - Tela de revisão (4h)
- [ ] RECEIPT-025: Testes de confirmação (2h)

---

### Fase 5: Melhorias e Otimizações (4-6 horas)

1. **Cache de Resultados**
   - Cachear resultados de OCR para imagens similares

2. **Retry Logic**
   - Retry automático em caso de falha

3. **Notificações**
   - Notificar usuário quando processamento completar

4. **Histórico**
   - Listar todos os comprovantes processados

5. **Métricas**
   - Taxa de sucesso de OCR
   - Tempo médio de processamento

#### Tarefas:
- [ ] RECEIPT-026: Implementar cache de OCR (2h)
- [ ] RECEIPT-027: Retry logic no N8N (2h)
- [ ] RECEIPT-028: Notificações de processamento (2h)
- [ ] RECEIPT-029: Listagem de comprovantes (2h)
- [ ] RECEIPT-030: Métricas e monitoramento (2h)

---

## 📦 Dependências e Serviços

### Novos Serviços:

1. **MinIO** (Storage S3-compatible)
   ```yaml
   minio:
     image: minio/minio:latest
     ports:
       - "9000:9000"
       - "9001:9001"
     environment:
       MINIO_ROOT_USER: minioadmin
       MINIO_ROOT_PASSWORD: minioadmin
     volumes:
       - minio_data:/data
     command: server /data --console-address ":9001"
   ```

2. **N8N** (Workflow Automation)
   ```yaml
   n8n:
     image: n8nio/n8n:latest
     ports:
       - "5678:5678"
     environment:
       - N8N_BASIC_AUTH_ACTIVE=true
       - N8N_BASIC_AUTH_USER=admin
       - N8N_BASIC_AUTH_PASSWORD=admin
       - DB_TYPE=postgresdb
       - DB_POSTGRESDB_HOST=postgres
       - DB_POSTGRESDB_DATABASE=n8n
     volumes:
       - n8n_data:/home/node/.n8n
     depends_on:
       - postgres
   ```

### APIs Externas (Configuração):

1. **Google Cloud Vision API**
   - Criar projeto no Google Cloud
   - Habilitar Vision API
   - Criar service account
   - Obter chave JSON

2. **OpenAI API** (ou Anthropic Claude)
   - Criar conta OpenAI
   - Obter API key
   - Configurar no N8N

---

## 💰 Estimativa de Custos

### Opção N8N (Recomendada):

**Infraestrutura:**
- MinIO: Gratuito (self-hosted)
- N8N: Gratuito (self-hosted)
- PostgreSQL: Já existe

**APIs Externas (por 1000 comprovantes/mês):**
- Google Vision API: ~$1.50 (primeiros 1000 são gratuitos)
- OpenAI GPT-4: ~$10-20 (depende do tamanho das imagens)
- **Total: ~$11-21/mês para 1000 comprovantes**

**Alternativa mais barata:**
- Tesseract OCR (gratuito) + GPT-3.5-turbo: ~$2-5/mês

---

## 🔒 Segurança

1. **Validação de Arquivos**
   - Tipos permitidos: JPG, PNG, PDF
   - Tamanho máximo: 10MB
   - Validação de MIME type
   - Scan de vírus (opcional)

2. **Autenticação**
   - Upload requer JWT válido
   - Usuário só acessa seus próprios comprovantes

3. **Storage**
   - Arquivos isolados por usuário
   - URLs com expiração (signed URLs)
   - Não expor caminhos diretos

4. **Dados Sensíveis**
   - Não logar dados extraídos
   - Criptografar dados em repouso (opcional)

---

## 📊 Métricas de Sucesso

1. **Precisão de OCR**: > 90% de extração correta
2. **Precisão de Classificação**: > 85% de identificação correta (receita/despesa)
3. **Tempo de Processamento**: < 30 segundos por comprovante
4. **Taxa de Sucesso**: > 95% de processamentos bem-sucedidos
5. **Adoção**: > 50% dos usuários usando a funcionalidade

---

## 🚀 Roadmap de Implementação

### Sprint 1 (Semana 1): Storage e Upload
- Configurar MinIO
- Implementar upload básico
- Frontend de upload

### Sprint 2 (Semana 2): Integração N8N
- Configurar N8N
- Criar workflow básico
- Integrar com backend

### Sprint 3 (Semana 3): Processamento
- OCR com Google Vision
- Análise com LLM
- Criação de transação DRAFT

### Sprint 4 (Semana 4): Confirmação
- Tela de revisão
- Confirmação/correção
- Aplicação de transação

### Sprint 5 (Semana 5): Melhorias
- Cache
- Retry logic
- Notificações
- Métricas

**Total: 5 semanas (~200 horas)**

---

## 🔄 Alternativas Consideradas

### 1. Tesseract OCR Nativo (Go)
- **Prós**: Gratuito, sem dependências externas
- **Contras**: Precisão menor, mais código para manter
- **Quando usar**: Se custo for crítico e precisão aceitável

### 2. AWS Textract
- **Prós**: Alta precisão, especializado em documentos
- **Contras**: Mais caro que Google Vision
- **Quando usar**: Se já usar AWS

### 3. Azure Form Recognizer
- **Prós**: Boa precisão, especializado em formulários
- **Contras**: Mais caro, menos comum
- **Quando usar**: Se já usar Azure

### 4. Solução Híbrida (Tesseract + LLM)
- **Prós**: Custo baixo, boa precisão
- **Contras**: Mais complexo
- **Quando usar**: Balance entre custo e precisão

---

## 📝 Considerações Finais

### Por que N8N?

1. **Velocidade**: Implementação muito mais rápida
2. **Flexibilidade**: Fácil ajustar workflow sem recompilar
3. **Manutenibilidade**: Workflows visuais são mais fáceis de entender
4. **Escalabilidade**: N8N gerencia filas e processamento assíncrono
5. **Integração**: Fácil integrar com múltiplas APIs
6. **Debugging**: Interface visual facilita debug

### Quando Reconsiderar?

- Se precisar de processamento em tempo real (< 5s)
- Se custo de APIs for proibitivo
- Se precisar de controle total sobre algoritmos
- Se precisar processar offline

---

## ✅ Checklist de Implementação

### Backend:
- [ ] Configurar MinIO
- [ ] Criar Receipt Context
- [ ] Implementar FileStorage
- [ ] Criar upload endpoint
- [ ] Configurar N8N
- [ ] Criar N8NClient
- [ ] Criar webhook handler
- [ ] Implementar ProcessReceiptUseCase
- [ ] Criar Transaction DRAFT
- [ ] Implementar ConfirmReceiptUseCase
- [ ] Adicionar notificações
- [ ] Implementar métricas

### Frontend:
- [ ] Componente de upload
- [ ] Tela de revisão
- [ ] Listagem de comprovantes
- [ ] Notificações de processamento

### N8N:
- [ ] Workflow de processamento
- [ ] Integração Google Vision
- [ ] Integração OpenAI/Claude
- [ ] Retry logic
- [ ] Error handling

### DevOps:
- [ ] Adicionar MinIO ao docker-compose
- [ ] Adicionar N8N ao docker-compose
- [ ] Configurar variáveis de ambiente
- [ ] Documentação de deploy

---

**Próximo Passo:** Revisar este planejamento e decidir entre Opção 2 (N8N) ou Opção 3 (Híbrida).


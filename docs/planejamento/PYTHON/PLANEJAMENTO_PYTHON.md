# Planejamento DDD - Sistema de Gestão Financeira (Python)

## Resumo Executivo

Sistema de gestão financeira pessoal e profissional desenvolvido em **Python 3.11+** com **FastAPI**, seguindo **Domain-Driven Design (DDD)**. Focado em **produtividade**, **type-safety com type hints**, **escalabilidade** e **pronto para produção**, com potencial para evoluir para produto SaaS.

**Stack Principal:**
- **Backend**: Python 3.11+ com FastAPI (framework web assíncrono)
- **Frontend**: Vue 3 com TypeScript (já existente - reutilizar) ⚠️ Ver seção 3.3
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis (cache e rate limiting)
- **ORM**: SQLAlchemy 2.0+ (async) ou Tortoise ORM
- **Observabilidade**: OpenTelemetry, Prometheus + Grafana

> ⚠️ **NOTA:** O projeto já possui um frontend Vue 3 funcional. Não é necessário criar novo frontend. Veja seção 3.3 para detalhes de compatibilidade.

**Diferenciais:**
- Type-safety com type hints nativos
- Performance assíncrona excepcional
- Arquitetura DDD escalável
- Observabilidade completa
- Segurança robusta
- Pronto para produção
- **Feature de Upload e Análise de Comprovantes integrada**

## 1. Visão Geral

Sistema de gestão financeira pessoal e profissional desenvolvido em **Python 3.11+** (backend) seguindo Domain-Driven Design (DDD), aproveitando a produtividade, type hints, ecossistema rico e performance assíncrona do Python moderno.

> ⚠️ **IMPORTANTE:** O projeto já possui um **frontend Vue 3** completamente funcional e independente. Este frontend deve ser **reutilizado** sem modificações, apenas configurando a URL da API. Veja seção 3.3 para detalhes de compatibilidade.

Projetado para **alta produtividade**, **escalabilidade** e **manutenibilidade**, com potencial para evoluir para produto comercial.

**Feature Especial:** O sistema inclui funcionalidade completa de **upload e análise automática de comprovantes** usando OCR/IA para extrair informações e criar transações automaticamente.

## 2. Objetivos

- Controle total de finanças pessoais e profissionais
- Separação clara entre contas pessoais e profissionais
- Análise e relatórios financeiros
- Planejamento orçamentário
- Acompanhamento de metas financeiras
- **Upload e análise automática de comprovantes via OCR/IA**
- Arquitetura escalável e manutenível
- Aproveitamento máximo do Python moderno e FastAPI

## 3. Stack Tecnológico Python

### 3.1. Tecnologias Principais

- **Linguagem**: Python 3.11+ (type hints, performance melhorada)
- **Framework Web**: **FastAPI** (moderno, assíncrono, type-safe)
- **ORM**: SQLAlchemy 2.0+ (async) ou Tortoise ORM
- **Validação**: Pydantic (type-safe, integrado ao FastAPI)
- **Autenticação**: python-jose (JWT) + passlib (hashing)
- **Event Bus**: Celery + Redis ou RabbitMQ
- **Testes**: pytest + pytest-asyncio
- **Migrations**: Alembic (SQLAlchemy) ou Aerich (Tortoise)
- **Config**: pydantic-settings
- **Logging**: structlog ou loguru
- **Banco de Dados**: PostgreSQL (com asyncpg ou psycopg3)
- **Cache**: Redis (redis-py ou aioredis)
- **API Docs**: Swagger/OpenAPI (automático no FastAPI)
- **Rate Limiting**: slowapi ou fastapi-limiter
- **Message Queue**: Celery + Redis/RabbitMQ
- **Monitoring**: Prometheus + Grafana
- **Tracing**: OpenTelemetry + Jaeger
- **Error Tracking**: Sentry (opcional)
- **File Storage**: MinIO (S3-compatible) ou boto3 (AWS S3)
- **OCR/IA**: Google Cloud Vision API, OpenAI GPT-4, ou Tesseract OCR

### 3.2. Framework Web: FastAPI

#### 3.2.1. Sobre o FastAPI

**FastAPI** ⚡ - Modern, fast web framework
- **GitHub**: 70k+ stars
- **Tipo**: Framework web moderno e assíncrono
- **Abordagem**: Baseado em type hints do Python
- **Base**: Starlette (ASGI framework) + Pydantic

#### 3.2.2. Características do FastAPI

**Vantagens:**
- ✅ **Type hints nativos** - Type-safety completo
- ✅ **Performance excepcional** - Uma das opções mais rápidas do Python
- ✅ **Documentação automática** - Swagger/OpenAPI gerado automaticamente
- ✅ **Validação automática** - Pydantic integrado
- ✅ **Assíncrono nativo** - async/await suportado
- ✅ **Fácil de usar** - API intuitiva e moderna
- ✅ **Documentação excelente** - Guias e exemplos completos
- ✅ **WebSocket support** - Suporte nativo a WebSockets
- ✅ **Dependency Injection** - Sistema poderoso de DI
- ✅ **Testes** - Suporte excelente para testes

**Considerações:**
- ⚠️ **Python 3.7+** - Requer Python moderno
- ⚠️ **Type hints** - Beneficia muito de type hints (não obrigatório, mas recomendado)
- ⚠️ **Comunidade** - Menor que Django/Flask, mas crescente e muito ativa

#### 3.2.3. Performance do FastAPI

**Benchmarks (aproximado):**

```
API Simples (Hello World):
- FastAPI:      ~50.000 req/s (async)
- Django:       ~10.000 req/s
- Flask:        ~15.000 req/s

API com JSON (CRUD):
- FastAPI:      ~45.000 req/s
- Django REST:  ~8.000 req/s
```

**Por que FastAPI é rápido?**
- Usa `Starlette` (ASGI framework de alta performance)
- Suporte nativo a async/await
- Validação com Pydantic (muito otimizada)
- Baseado em type hints (menos overhead)

#### 3.2.4. Exemplos de Código com FastAPI

##### Exemplo Básico

```python
# main.py
from fastapi import FastAPI, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import uvicorn

from app.core.config import settings
from app.api.v1 import auth, transactions, accounts
from app.core.database import init_db, close_db

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    await init_db()
    yield
    # Shutdown
    await close_db()

app = FastAPI(
    title="Gestão Financeira API",
    version="1.0.0",
    lifespan=lifespan
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.ALLOWED_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Health checks
@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "gestao-financeira"}

# API v1
app.include_router(auth.router, prefix="/api/v1/auth", tags=["auth"])
app.include_router(
    transactions.router,
    prefix="/api/v1/transactions",
    tags=["transactions"]
)
app.include_router(
    accounts.router,
    prefix="/api/v1/accounts",
    tags=["accounts"]
)

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8000,
        reload=True
    )
```

##### Router com FastAPI

```python
# app/api/v1/auth.py
from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel, EmailStr
from typing import Optional

from app.core.security import verify_password, create_access_token
from app.core.dependencies import get_current_user
from app.domain.identity.entities import User
from app.application.identity.usecases import RegisterUserUseCase, LoginUseCase

router = APIRouter()

class RegisterRequest(BaseModel):
    email: EmailStr
    password: str
    first_name: str
    last_name: str

class RegisterResponse(BaseModel):
    user_id: str
    email: str
    first_name: str
    last_name: str

class LoginResponse(BaseModel):
    access_token: str
    token_type: str = "bearer"
    user: dict

@router.post("/register", response_model=RegisterResponse, status_code=status.HTTP_201_CREATED)
async def register(
    request: RegisterRequest,
    use_case: RegisterUserUseCase = Depends()
):
    """Registrar novo usuário"""
    try:
        result = await use_case.execute(
            email=request.email,
            password=request.password,
            first_name=request.first_name,
            last_name=request.last_name
        )
        return RegisterResponse(**result)
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )

@router.post("/login", response_model=LoginResponse)
async def login(
    form_data: OAuth2PasswordRequestForm = Depends(),
    use_case: LoginUseCase = Depends()
):
    """Login e obter token JWT"""
    try:
        result = await use_case.execute(
            email=form_data.username,
            password=form_data.password
        )
        return LoginResponse(**result)
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail=str(e)
        )

@router.get("/me")
async def get_current_user_info(
    current_user: User = Depends(get_current_user)
):
    """Obter informações do usuário atual"""
    return {
        "user_id": current_user.id.value,
        "email": current_user.email.value,
        "first_name": current_user.first_name,
        "last_name": current_user.last_name
    }
```

##### Dependency Injection

```python
# app/core/dependencies.py
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt

from app.core.config import settings
from app.domain.identity.value_objects import UserID
from app.infrastructure.identity.repositories import UserRepository

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/api/v1/auth/login")

async def get_current_user(
    token: str = Depends(oauth2_scheme),
    user_repo: UserRepository = Depends()
) -> User:
    """Dependency para obter usuário atual do JWT"""
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    
    try:
        payload = jwt.decode(
            token,
            settings.SECRET_KEY,
            algorithms=[settings.ALGORITHM]
        )
        user_id: str = payload.get("sub")
        if user_id is None:
            raise credentials_exception
    except JWTError:
        raise credentials_exception
    
    user = await user_repo.find_by_id(UserID(user_id))
    if user is None:
        raise credentials_exception
    
    return user
```

##### Middleware Customizado

```python
# app/core/middleware.py
from fastapi import Request, Response
from starlette.middleware.base import BaseHTTPMiddleware
import time
import uuid

from app.core.logger import logger

class RequestIDMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        # Gerar request ID
        request_id = str(uuid.uuid4())
        request.state.request_id = request_id
        
        # Adicionar ao header de resposta
        start_time = time.time()
        response = await call_next(request)
        process_time = time.time() - start_time
        
        response.headers["X-Request-ID"] = request_id
        response.headers["X-Process-Time"] = str(process_time)
        
        # Log
        logger.info(
            "Request processed",
            request_id=request_id,
            method=request.method,
            path=request.url.path,
            status_code=response.status_code,
            process_time=process_time
        )
        
        return response
```

#### 3.2.5. Validação com Pydantic

```python
# app/application/transactions/dtos.py
from pydantic import BaseModel, Field, validator
from decimal import Decimal
from datetime import date
from typing import Optional
from enum import Enum

class TransactionType(str, Enum):
    INCOME = "INCOME"
    EXPENSE = "EXPENSE"
    TRANSFER = "TRANSFER"

class CreateTransactionRequest(BaseModel):
    account_id: str = Field(..., description="ID da conta")
    category_id: str = Field(..., description="ID da categoria")
    type: TransactionType = Field(..., description="Tipo da transação")
    amount: Decimal = Field(..., gt=0, description="Valor da transação")
    description: Optional[str] = Field(None, max_length=500)
    date: date = Field(..., description="Data da transação")
    tags: list[str] = Field(default_factory=list)
    
    @validator('amount')
    def validate_amount(cls, v):
        if v <= 0:
            raise ValueError('Amount must be positive')
        return v
    
    @validator('description')
    def validate_description(cls, v):
        if v and len(v.strip()) == 0:
            raise ValueError('Description cannot be empty')
        return v

class TransactionResponse(BaseModel):
    id: str
    account_id: str
    category_id: str
    type: TransactionType
    amount: Decimal
    description: Optional[str]
    date: date
    tags: list[str]
    created_at: str
    updated_at: str
    
    class Config:
        from_attributes = True
```

#### 3.2.6. Vantagens do FastAPI para este Projeto

**Por que FastAPI é uma excelente escolha:**
1. ✅ **Type Safety** - Type hints nativos do Python
2. ✅ **Performance** - Uma das opções mais rápidas do Python
3. ✅ **Documentação Automática** - Swagger gerado automaticamente
4. ✅ **Validação Automática** - Pydantic integrado
5. ✅ **Assíncrono** - Suporte nativo a async/await
6. ✅ **Fácil de Aprender** - API intuitiva e moderna
7. ✅ **Compatível com DDD** - Não impõe estrutura, você controla
8. ✅ **WebSocket Support** - Para features real-time no futuro
9. ✅ **Dependency Injection** - Sistema poderoso e elegante

**Considerações:**
- ⚠️ **Python 3.7+** - Requer Python moderno (3.11+ recomendado)
- ⚠️ **Type Hints** - Beneficia muito de type hints (altamente recomendado)

### 3.3. Compatibilidade com Frontend Vue 3 Existente

**✅ IMPORTANTE:** O projeto já possui um frontend Vue 3 completamente funcional e independente. **NÃO é necessário criar um novo frontend** para Python.

#### 3.3.1. Por que o Frontend Vue 3 é Reutilizável?

O frontend Vue 3 atual foi desenvolvido de forma **desacoplada** do backend, comunicando-se exclusivamente via **API REST**. Isso significa:
- ✅ **Arquitetura independente**: Frontend não depende da tecnologia do backend (Go, Node.js, PHP, Python, etc.)
- ✅ **Interface padronizada**: Comunicação via JSON sobre HTTP
- ✅ **Zero alterações no código**: O frontend Vue 3 funciona sem modificações, apenas com a configuração da URL da API.

#### 3.3.2. Requisitos de Compatibilidade da API

Para que o frontend Vue 3 existente funcione com o backend Python/FastAPI, é necessário implementar a **mesma interface de API REST** que o backend Go atual oferece. Isso inclui:

- **Mesmos Endpoints**:
  ```
  GET    /api/v1/health
  POST   /api/v1/auth/register
  POST   /api/v1/auth/login
  GET    /api/v1/transactions
  POST   /api/v1/transactions
  GET    /api/v1/transactions/{id}
  PUT    /api/v1/transactions/{id}
  DELETE /api/v1/transactions/{id}
  GET    /api/v1/accounts
  POST   /api/v1/accounts
  GET    /api/v1/categories
  POST   /api/v1/categories
  GET    /api/v1/budgets
  POST   /api/v1/budgets
  GET    /api/v1/reports/monthly
  POST   /api/v1/receipts/upload      # Nova feature
  GET    /api/v1/receipts              # Nova feature
  POST   /api/v1/receipts/{id}/confirm # Nova feature
  ... (outros endpoints conforme necessário)
  ```

- **Formato de Request/Response:**
  - Content-Type: `application/json`
  - Autenticação: JWT Bearer Token (header `Authorization: Bearer <token>`)
  - Formato de resposta padronizado:
    ```json
    {
      "data": { ... },
      "message": "Success",
      "status": 200
    }
    ```
  - Códigos de status HTTP padronizados (200, 201, 400, 401, 404, 500, etc.)

- **Autenticação:**
  - JWT tokens com mesma estrutura de payload
  - Refresh tokens (se implementado)
  - Mesmos headers de autenticação

#### 3.3.3. Configuração do Frontend para Python

Para usar o frontend Vue 3 com o backend Python/FastAPI, apenas configure a variável de ambiente:

```bash
# .env ou .env.local no frontend
VITE_API_URL=http://localhost:8000/api/v1
```

Ou no `docker-compose.yml`:
```yaml
frontend:
  environment:
    - VITE_API_URL=http://python-backend:8000/api/v1
```

**Nenhuma alteração no código do frontend é necessária!**

#### 3.3.4. Vantagens dessa Abordagem

1. **Migração facilitada**: Trocar de backend Go para Python é apenas implementar a API
2. **Reutilização total**: Frontend Vue 3 já desenvolvido e testado
3. **Economia de tempo**: Não precisa desenvolver novo frontend
4. **Consistência**: Experiência do usuário consistente, independentemente do backend
5. **Flexibilidade**: Permite experimentar diferentes tecnologias de backend sem impactar o frontend.

---

### 3.4. Por que Python?

**Vantagens:**
- ✅ **Produtividade excepcional**: Código limpo e expressivo
- ✅ **Type hints**: Type-safety com type hints nativos
- ✅ **Ecossistema rico**: Muitas bibliotecas disponíveis
- ✅ **Assíncrono**: Suporte nativo a async/await
- ✅ **Fácil de aprender**: Linguagem intuitiva
- ✅ **ML/IA**: Excelente para integração com ML/IA (OCR, LLMs)
- ✅ **Comunidade**: Grande comunidade e suporte
- ✅ **Versatilidade**: Pode integrar facilmente com serviços de IA

**Desafios:**
- ⚠️ **Performance**: Mais lento que Go/Rust, mas FastAPI é muito rápido
- ⚠️ **GIL**: Global Interpreter Lock (mas async contorna isso)
- ⚠️ **Deploy**: Requer ambiente Python (mas Docker resolve)

## 4. Arquitetura DDD - Estrutura de Pastas

### 4.1. Estrutura Geral

```
backend/
├── app/
│   ├── __init__.py
│   ├── main.py                    # Entry point FastAPI
│   ├── core/                      # Configurações centrais
│   │   ├── __init__.py
│   │   ├── config.py              # Configurações (pydantic-settings)
│   │   ├── database.py            # Conexão DB (SQLAlchemy/Tortoise)
│   │   ├── security.py            # JWT, hashing
│   │   ├── dependencies.py        # Dependency injection
│   │   ├── middleware.py          # Middlewares customizados
│   │   └── logger.py              # Logging (structlog/loguru)
│   ├── api/                       # Camada de apresentação
│   │   ├── __init__.py
│   │   ├── v1/                     # API v1
│   │   │   ├── __init__.py
│   │   │   ├── auth.py             # Rotas de autenticação
│   │   │   ├── transactions.py    # Rotas de transações
│   │   │   ├── accounts.py         # Rotas de contas
│   │   │   ├── categories.py       # Rotas de categorias
│   │   │   ├── budgets.py          # Rotas de orçamentos
│   │   │   ├── reports.py          # Rotas de relatórios
│   │   │   └── receipts.py         # Rotas de comprovantes (nova feature)
│   │   └── dependencies.py         # Dependencies compartilhadas
│   ├── domain/                    # Camada de domínio (DDD)
│   │   ├── identity/               # Bounded Context: Identity
│   │   │   ├── entities/
│   │   │   │   └── user.py
│   │   │   ├── value_objects/
│   │   │   │   ├── user_id.py
│   │   │   │   ├── email.py
│   │   │   │   └── password.py
│   │   │   ├── repositories/
│   │   │   │   └── user_repository.py  # Interface
│   │   │   └── events/
│   │   │       └── user_registered.py
│   │   ├── account/                # Bounded Context: Account
│   │   │   ├── entities/
│   │   │   │   └── account.py
│   │   │   ├── value_objects/
│   │   │   │   ├── account_id.py
│   │   │   │   ├── money.py
│   │   │   │   └── account_context.py
│   │   │   └── repositories/
│   │   │       └── account_repository.py
│   │   ├── transaction/            # Bounded Context: Transaction
│   │   │   ├── entities/
│   │   │   │   ├── transaction.py
│   │   │   │   └── recurring_transaction.py
│   │   │   ├── value_objects/
│   │   │   │   ├── transaction_id.py
│   │   │   │   ├── transaction_type.py
│   │   │   │   └── money.py
│   │   │   └── repositories/
│   │   │       └── transaction_repository.py
│   │   ├── category/               # Bounded Context: Category
│   │   │   ├── entities/
│   │   │   │   └── category.py
│   │   │   └── value_objects/
│   │   │       └── category_id.py
│   │   ├── budget/                 # Bounded Context: Budget
│   │   │   ├── entities/
│   │   │   │   └── budget.py
│   │   │   └── value_objects/
│   │   │       └── budget_id.py
│   │   ├── receipt/                # Bounded Context: Receipt (NOVA)
│   │   │   ├── entities/
│   │   │   │   └── receipt.py
│   │   │   ├── value_objects/
│   │   │   │   ├── receipt_id.py
│   │   │   │   ├── receipt_status.py
│   │   │   │   └── file_info.py
│   │   │   ├── repositories/
│   │   │   │   └── receipt_repository.py
│   │   │   └── events/
│   │   │       ├── receipt_uploaded.py
│   │   │       └── receipt_processed.py
│   │   └── shared/                 # Value objects compartilhados
│   │       └── value_objects/
│   │           └── money.py
│   ├── application/                # Camada de aplicação (Use Cases)
│   │   ├── identity/
│   │   │   └── usecases/
│   │   │       ├── register_user_usecase.py
│   │   │       └── login_usecase.py
│   │   ├── transaction/
│   │   │   └── usecases/
│   │   │       ├── create_transaction_usecase.py
│   │   │       ├── list_transactions_usecase.py
│   │   │       └── update_transaction_usecase.py
│   │   ├── receipt/                # Use cases de comprovantes (NOVA)
│   │   │   └── usecases/
│   │   │       ├── upload_receipt_usecase.py
│   │   │       ├── process_receipt_usecase.py
│   │   │       └── confirm_receipt_usecase.py
│   │   └── shared/
│   │       └── dtos/                # DTOs compartilhados
│   └── infrastructure/             # Camada de infraestrutura
│       ├── identity/
│       │   └── repositories/
│       │       └── sqlalchemy_user_repository.py
│       ├── transaction/
│       │   └── repositories/
│       │       └── sqlalchemy_transaction_repository.py
│       ├── receipt/                # Infra de comprovantes (NOVA)
│       │   ├── repositories/
│       │   │   └── sqlalchemy_receipt_repository.py
│       │   ├── storage/
│       │   │   └── file_storage.py  # MinIO/S3
│       │   └── services/
│       │       ├── ocr_service.py   # Google Vision/Tesseract
│       │       ├── ai_service.py     # OpenAI GPT-4
│       │       └── n8n_client.py     # Cliente N8N (opcional)
│       ├── database/
│       │   ├── models/               # SQLAlchemy models
│       │   │   ├── user_model.py
│       │   │   ├── transaction_model.py
│       │   │   └── receipt_model.py
│       │   └── base.py
│       └── cache/
│           └── redis_cache.py
├── tests/                           # Testes
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── alembic/                         # Migrations (SQLAlchemy)
│   ├── versions/
│   └── env.py
├── requirements.txt                 # Dependências
├── requirements-dev.txt             # Dependências de desenvolvimento
├── pyproject.toml                   # Configuração do projeto
├── Dockerfile
└── docker-compose.yml
```

### 4.2. Exemplo de Entidade de Domínio

```python
# app/domain/transaction/entities/transaction.py
from dataclasses import dataclass
from datetime import date
from decimal import Decimal
from typing import Optional

from app.domain.transaction.value_objects import (
    TransactionID,
    TransactionType,
    Money
)
from app.domain.identity.value_objects import UserID
from app.domain.account.value_objects import AccountID
from app.domain.category.value_objects import CategoryID

@dataclass
class Transaction:
    """Agregado raiz: Transaction"""
    id: TransactionID
    user_id: UserID
    account_id: AccountID
    category_id: CategoryID
    type: TransactionType
    amount: Money
    description: Optional[str]
    date: date
    tags: list[str]
    status: str  # PENDING, APPROVED, CANCELLED
    created_at: date
    updated_at: date
    
    def approve(self) -> None:
        """Aprovar transação"""
        if self.status != "PENDING":
            raise ValueError("Only pending transactions can be approved")
        self.status = "APPROVED"
    
    def cancel(self) -> None:
        """Cancelar transação"""
        if self.status == "CANCELLED":
            raise ValueError("Transaction already cancelled")
        self.status = "CANCELLED"
    
    def update_amount(self, new_amount: Money) -> None:
        """Atualizar valor da transação"""
        if new_amount.value <= 0:
            raise ValueError("Amount must be positive")
        self.amount = new_amount
    
    def _ensure_valid(self) -> None:
        """Validar invariantes"""
        if self.amount.value <= 0:
            raise ValueError("Transaction amount must be positive")
```

### 4.3. Exemplo de Use Case

```python
# app/application/transaction/usecases/create_transaction_usecase.py
from typing import Protocol
from datetime import date
from decimal import Decimal

from app.domain.transaction.entities import Transaction
from app.domain.transaction.value_objects import TransactionID, TransactionType, Money
from app.domain.identity.value_objects import UserID
from app.domain.account.value_objects import AccountID
from app.domain.category.value_objects import CategoryID

class TransactionRepository(Protocol):
    async def save(self, transaction: Transaction) -> None: ...
    async def find_by_id(self, transaction_id: TransactionID) -> Transaction | None: ...

class AccountRepository(Protocol):
    async def find_by_id(self, account_id: AccountID) -> ...: ...

class CreateTransactionUseCase:
    def __init__(
        self,
        transaction_repo: TransactionRepository,
        account_repo: AccountRepository
    ):
        self._transaction_repo = transaction_repo
        self._account_repo = account_repo
    
    async def execute(
        self,
        user_id: str,
        account_id: str,
        category_id: str,
        type: str,
        amount: Decimal,
        description: str | None,
        date: date,
        tags: list[str] = None
    ) -> dict:
        """Criar nova transação"""
        # Validar conta
        account = await self._account_repo.find_by_id(AccountID(account_id))
        if not account:
            raise ValueError("Account not found")
        
        # Criar entidade
        transaction = Transaction(
            id=TransactionID.generate(),
            user_id=UserID(user_id),
            account_id=AccountID(account_id),
            category_id=CategoryID(category_id),
            type=TransactionType(type),
            amount=Money(amount),
            description=description,
            date=date,
            tags=tags or [],
            status="PENDING",
            created_at=date.today(),
            updated_at=date.today()
        )
        
        # Salvar
        await self._transaction_repo.save(transaction)
        
        return {
            "id": transaction.id.value,
            "account_id": transaction.account_id.value,
            "category_id": transaction.category_id.value,
            "type": transaction.type.value,
            "amount": str(transaction.amount.value),
            "description": transaction.description,
            "date": transaction.date.isoformat(),
            "tags": transaction.tags,
            "status": transaction.status
        }
```

## 5. Feature: Upload e Análise Automática de Comprovantes

### 5.1. Visão Geral

Implementar funcionalidade completa para upload de comprovantes (imagens de recibos, notas fiscais, extratos) e análise automática via OCR/IA para extrair informações e criar transações automaticamente.

### 5.2. Arquitetura da Feature

#### 5.2.1. Fluxo Completo

```
1. Usuário faz upload no Frontend
   ↓
2. Frontend → POST /api/v1/receipts/upload
   ↓
3. Backend Python/FastAPI:
   - Valida arquivo (tipo, tamanho)
   - Salva em storage (MinIO/S3)
   - Cria Receipt entity (status: PENDING)
   - Retorna receipt_id
   ↓
4. Backend Python → Webhook N8N (receipt_id, file_url)
   OU
   Backend Python → Processamento direto (Google Vision + OpenAI)
   ↓
5. Processamento (N8N ou direto):
   a. Download imagem
   b. OCR (Google Vision API ou Tesseract)
   c. Extração de texto
   d. Análise com LLM (GPT-4):
      - Identifica tipo (receita/despesa)
      - Extrai valor
      - Extrai data
      - Extrai descrição
      - Extrai categoria (se possível)
   e. Retorna dados estruturados
   ↓
6. Backend Python:
   - Atualiza Receipt (status: PROCESSED, extracted_data)
   - Cria Transaction DRAFT (aguardando confirmação)
   - Notifica usuário
   ↓
7. Usuário revisa e confirma/corrige
   ↓
8. Backend Python:
   - Atualiza Transaction (status: CONFIRMED)
   - Aplica transação (atualiza saldo)
```

#### 5.2.2. Opções de Implementação

**Opção 1: Processamento Direto (Python Nativo)**
- **OCR**: Google Cloud Vision API ou Tesseract OCR (pytesseract)
- **IA**: OpenAI GPT-4 API
- **Vantagens**: Controle total, integração nativa
- **Desvantagens**: Mais código para manter

**Opção 2: N8N (Workflow Automation)**
- **N8N**: Workflow automation (self-hosted ou cloud)
- **OCR**: Google Cloud Vision API, AWS Textract
- **IA**: OpenAI GPT-4, Anthropic Claude
- **Vantagens**: Visual, flexível, rápido de implementar
- **Desvantagens**: Dependência externa

**Recomendação**: Opção 2 (N8N) para MVP, Opção 1 para produção com controle total.

### 5.3. Estrutura de Implementação

#### 5.3.1. Receipt Context (Novo Bounded Context)

```
app/domain/receipt/
├── entities/
│   └── receipt.py          # Receipt aggregate root
├── value_objects/
│   ├── receipt_id.py
│   ├── receipt_status.py  # PENDING, PROCESSING, PROCESSED, FAILED
│   └── file_info.py        # filename, size, mime_type
├── repositories/
│   └── receipt_repository.py
└── events/
    ├── receipt_uploaded.py
    └── receipt_processed.py
```

#### 5.3.2. Entidade Receipt

```python
# app/domain/receipt/entities/receipt.py
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, Any

from app.domain.receipt.value_objects import ReceiptID, ReceiptStatus, FileInfo
from app.domain.identity.value_objects import UserID
from app.domain.transaction.value_objects import TransactionID

@dataclass
class Receipt:
    """Agregado raiz: Receipt"""
    id: ReceiptID
    user_id: UserID
    file_info: FileInfo
    status: ReceiptStatus
    extracted_data: Optional[dict[str, Any]]  # Dados extraídos do OCR/IA
    transaction_id: Optional[TransactionID]  # Transação criada (se houver)
    error_message: Optional[str]
    created_at: datetime
    updated_at: datetime
    
    def mark_as_processing(self) -> None:
        """Marcar como processando"""
        self.status = ReceiptStatus.PROCESSING
        self.updated_at = datetime.now()
    
    def mark_as_processed(self, extracted_data: dict[str, Any]) -> None:
        """Marcar como processado"""
        self.status = ReceiptStatus.PROCESSED
        self.extracted_data = extracted_data
        self.updated_at = datetime.now()
    
    def mark_as_failed(self, error_message: str) -> None:
        """Marcar como falhou"""
        self.status = ReceiptStatus.FAILED
        self.error_message = error_message
        self.updated_at = datetime.now()
    
    def link_transaction(self, transaction_id: TransactionID) -> None:
        """Vincular transação criada"""
        self.transaction_id = transaction_id
        self.updated_at = datetime.now()
```

#### 5.3.3. Use Case: Upload Receipt

```python
# app/application/receipt/usecases/upload_receipt_usecase.py
from fastapi import UploadFile
from typing import Protocol

from app.domain.receipt.entities import Receipt
from app.domain.receipt.value_objects import ReceiptID, ReceiptStatus, FileInfo
from app.domain.identity.value_objects import UserID
from app.infrastructure.receipt.storage import FileStorage

class ReceiptRepository(Protocol):
    async def save(self, receipt: Receipt) -> None: ...

class UploadReceiptUseCase:
    def __init__(
        self,
        receipt_repo: ReceiptRepository,
        file_storage: FileStorage
    ):
        self._receipt_repo = receipt_repo
        self._file_storage = file_storage
    
    async def execute(
        self,
        user_id: str,
        file: UploadFile
    ) -> dict:
        """Upload de comprovante"""
        # Validar arquivo
        if file.content_type not in ["image/jpeg", "image/png", "application/pdf"]:
            raise ValueError("Invalid file type")
        
        if file.size > 10 * 1024 * 1024:  # 10MB
            raise ValueError("File too large")
        
        # Salvar arquivo
        file_path = await self._file_storage.save(file, user_id)
        
        # Criar entidade
        receipt = Receipt(
            id=ReceiptID.generate(),
            user_id=UserID(user_id),
            file_info=FileInfo(
                filename=file.filename,
                size=file.size,
                mime_type=file.content_type,
                path=file_path
            ),
            status=ReceiptStatus.PENDING,
            extracted_data=None,
            transaction_id=None,
            error_message=None,
            created_at=datetime.now(),
            updated_at=datetime.now()
        )
        
        # Salvar
        await self._receipt_repo.save(receipt)
        
        # Trigger processamento assíncrono (Celery ou N8N)
        await self._trigger_processing(receipt.id)
        
        return {
            "receipt_id": receipt.id.value,
            "status": receipt.status.value,
            "file_info": {
                "filename": receipt.file_info.filename,
                "size": receipt.file_info.size,
                "mime_type": receipt.file_info.mime_type
            }
        }
    
    async def _trigger_processing(self, receipt_id: ReceiptID) -> None:
        """Trigger processamento assíncrono"""
        # Opção 1: Celery task
        # process_receipt_task.delay(receipt_id.value)
        
        # Opção 2: N8N webhook
        # await n8n_client.trigger_workflow(receipt_id.value)
        pass
```

#### 5.3.4. Router FastAPI

```python
# app/api/v1/receipts.py
from fastapi import APIRouter, Depends, UploadFile, File, HTTPException, status
from typing import List

from app.core.dependencies import get_current_user
from app.domain.identity.entities import User
from app.application.receipt.usecases import (
    UploadReceiptUseCase,
    ProcessReceiptUseCase,
    ConfirmReceiptUseCase
)

router = APIRouter()

@router.post("/upload", status_code=status.HTTP_201_CREATED)
async def upload_receipt(
    file: UploadFile = File(...),
    current_user: User = Depends(get_current_user),
    use_case: UploadReceiptUseCase = Depends()
):
    """Upload de comprovante"""
    try:
        result = await use_case.execute(
            user_id=current_user.id.value,
            file=file
        )
        return result
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )

@router.get("/")
async def list_receipts(
    current_user: User = Depends(get_current_user),
    receipt_repo: ReceiptRepository = Depends()
):
    """Listar comprovantes do usuário"""
    receipts = await receipt_repo.find_by_user_id(
        UserID(current_user.id.value)
    )
    return [receipt_to_dict(r) for r in receipts]

@router.post("/{receipt_id}/confirm")
async def confirm_receipt(
    receipt_id: str,
    corrections: dict,  # Correções manuais do usuário
    current_user: User = Depends(get_current_user),
    use_case: ConfirmReceiptUseCase = Depends()
):
    """Confirmar e criar transação a partir do comprovante"""
    try:
        result = await use_case.execute(
            receipt_id=receipt_id,
            user_id=current_user.id.value,
            corrections=corrections
        )
        return result
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
```

### 5.4. Migração de Banco de Dados

```python
# alembic/versions/010_create_receipts_table.py
"""create receipts table

Revision ID: 010
Revises: 009
Create Date: 2025-01-01 10:00:00.000000

"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

def upgrade():
    op.create_table(
        'receipts',
        sa.Column('id', postgresql.UUID(as_uuid=True), primary_key=True),
        sa.Column('user_id', postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column('file_name', sa.String(255), nullable=False),
        sa.Column('file_path', sa.String(500), nullable=False),
        sa.Column('file_size', sa.BigInteger(), nullable=False),
        sa.Column('mime_type', sa.String(100), nullable=False),
        sa.Column('status', sa.String(20), nullable=False),
        sa.Column('extracted_data', postgresql.JSONB, nullable=True),
        sa.Column('transaction_id', postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column('error_message', sa.Text(), nullable=True),
        sa.Column('created_at', sa.DateTime(), nullable=False),
        sa.Column('updated_at', sa.DateTime(), nullable=False),
        sa.Column('deleted_at', sa.DateTime(), nullable=True),
        sa.ForeignKeyConstraint(['user_id'], ['users.id'], ondelete='CASCADE'),
        sa.ForeignKeyConstraint(['transaction_id'], ['transactions.id'], ondelete='SET NULL'),
    )
    
    op.create_index('idx_receipts_user_id', 'receipts', ['user_id'])
    op.create_index('idx_receipts_status', 'receipts', ['status'])
    op.create_index('idx_receipts_created_at', 'receipts', ['created_at'])

def downgrade():
    op.drop_table('receipts')
```

## 6. Fases de Desenvolvimento

### Fase 1: Fundação (2-3 semanas)
- Setup do projeto (FastAPI, SQLAlchemy, estrutura DDD)
- Identity Context (autenticação, usuários)
- Configuração de banco de dados e migrations
- Health checks e documentação Swagger

### Fase 2: Core Financeiro (3-4 semanas)
- Account Context (contas e carteiras)
- Transaction Context (transações)
- Category Context (categorias)
- Relacionamentos e validações

### Fase 3: Features Avançadas (3-4 semanas)
- Budget Context (orçamentos)
- Reporting Context (relatórios)
- Recurring Transactions
- **Receipt Context (upload e análise de comprovantes)**

### Fase 4: Observabilidade e Segurança (2 semanas)
- Logging estruturado
- OpenTelemetry
- Prometheus + Grafana
- Rate limiting
- Security headers

### Fase 5: Otimizações e Deploy (2 semanas)
- Cache (Redis)
- Otimizações de queries
- Docker e docker-compose
- CI/CD
- Documentação final

**Total estimado: 12-15 semanas**

## 7. Performance e Otimizações

### 7.1. Async/Await

FastAPI suporta nativamente async/await, permitindo alta concorrência:

```python
@app.get("/transactions")
async def list_transactions(
    current_user: User = Depends(get_current_user)
):
    # Operações assíncronas não bloqueiam
    transactions = await transaction_repo.find_by_user_id(user_id)
    return transactions
```

### 7.2. Connection Pooling

```python
# app/core/database.py
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

engine = create_async_engine(
    settings.DATABASE_URL,
    pool_size=20,
    max_overflow=10,
    pool_pre_ping=True
)

AsyncSessionLocal = sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False
)
```

### 7.3. Cache com Redis

```python
# app/infrastructure/cache/redis_cache.py
from redis import asyncio as aioredis
import json

class RedisCache:
    def __init__(self, redis_url: str):
        self.redis = aioredis.from_url(redis_url)
    
    async def get(self, key: str) -> dict | None:
        value = await self.redis.get(key)
        if value:
            return json.loads(value)
        return None
    
    async def set(self, key: str, value: dict, ttl: int = 3600):
        await self.redis.setex(
            key,
            ttl,
            json.dumps(value)
        )
```

## 8. Testes

### 8.1. Testes Unitários

```python
# tests/unit/application/transaction/test_create_transaction_usecase.py
import pytest
from decimal import Decimal
from datetime import date

from app.application.transaction.usecases import CreateTransactionUseCase
from app.domain.transaction.entities import Transaction

@pytest.mark.asyncio
async def test_create_transaction_success():
    # Arrange
    mock_repo = MockTransactionRepository()
    mock_account_repo = MockAccountRepository()
    use_case = CreateTransactionUseCase(mock_repo, mock_account_repo)
    
    # Act
    result = await use_case.execute(
        user_id="user-123",
        account_id="account-123",
        category_id="category-123",
        type="EXPENSE",
        amount=Decimal("100.50"),
        description="Test transaction",
        date=date.today()
    )
    
    # Assert
    assert result["id"] is not None
    assert result["amount"] == "100.50"
    assert len(mock_repo.saved_transactions) == 1
```

### 8.2. Testes de Integração

```python
# tests/integration/api/test_transactions_api.py
from fastapi.testclient import TestClient

def test_create_transaction(client: TestClient, auth_token: str):
    response = client.post(
        "/api/v1/transactions",
        json={
            "account_id": "account-123",
            "category_id": "category-123",
            "type": "EXPENSE",
            "amount": "100.50",
            "description": "Test",
            "date": "2025-01-01"
        },
        headers={"Authorization": f"Bearer {auth_token}"}
    )
    
    assert response.status_code == 201
    assert response.json()["id"] is not None
```

## 9. Deploy e DevOps

### 9.1. Dockerfile

```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Instalar dependências do sistema
RUN apt-get update && apt-get install -y \
    gcc \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Copiar requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copiar código
COPY . .

# Expor porta
EXPOSE 8000

# Comando
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

### 9.2. docker-compose.yml

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql+asyncpg://user:pass@db:5432/gestao_financeira
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
  
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=gestao_financeira
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
  
  n8n:
    image: n8nio/n8n
    ports:
      - "5678:5678"
    environment:
      - N8N_BASIC_AUTH_ACTIVE=true
      - N8N_BASIC_AUTH_USER=admin
      - N8N_BASIC_AUTH_PASSWORD=admin

volumes:
  postgres_data:
```

## 10. Conclusão

Este planejamento fornece uma base sólida para implementar o sistema de gestão financeira em Python com FastAPI, seguindo DDD e incluindo a feature completa de upload e análise automática de comprovantes.

**Principais vantagens da stack Python/FastAPI:**
- ✅ Produtividade excepcional
- ✅ Type-safety com type hints
- ✅ Performance assíncrona
- ✅ Ecossistema rico (especialmente para ML/IA)
- ✅ Documentação automática
- ✅ Fácil integração com serviços de IA (OCR, LLMs)

**Próximos passos:**
1. Setup inicial do projeto
2. Implementação do Identity Context
3. Desenvolvimento incremental dos outros contexts
4. Integração da feature de comprovantes
5. Testes e otimizações
6. Deploy e monitoramento


# Planejamento DDD - Sistema de Gestão Financeira (Go)

## Resumo Executivo

Sistema de gestão financeira pessoal e profissional desenvolvido em **Go** com **Fiber**, seguindo **Domain-Driven Design (DDD)**. Focado em **alta performance**, **escalabilidade** e **pronto para produção**, com potencial para evoluir para produto SaaS.

**Stack Principal:**
- **Backend**: Go 1.21+ com Fiber (framework web)
- **Frontend**: Next.js 14+ com TypeScript, shadcn/ui e Tailwind CSS
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis (cache e rate limiting)
- **ORM**: GORM
- **Observabilidade**: OpenTelemetry, Prometheus + Grafana

**Diferenciais:**
- Performance excepcional (~200k req/s)
- Arquitetura DDD escalável
- Observabilidade completa
- Segurança robusta
- Pronto para produção

## 1. Visão Geral

Sistema de gestão financeira pessoal e profissional desenvolvido em **Go** (backend) e **Next.js** (frontend) com Domain-Driven Design (DDD), aproveitando a performance excepcional, concorrência nativa e simplicidade da linguagem. O frontend é desenvolvido **incrementalmente**, sempre que um módulo do backend estiver pronto, utilizando **shadcn/ui** para interfaces modernas e acessíveis. Projetado para **alta performance** e **escalabilidade**, com potencial para evoluir para produto comercial.

## 2. Objetivos

- Controle total de finanças pessoais e profissionais
- Separação clara entre contas pessoais e profissionais
- Análise e relatórios financeiros
- Planejamento orçamentário
- Acompanhamento de metas financeiras
- Arquitetura escalável e manutenível
- Aproveitamento máximo da performance e concorrência do Go

## 3. Stack Tecnológico Go

### 3.1. Tecnologias Principais

- **Linguagem**: Go 1.21+
- **Framework Web**: **Fiber** (Express-inspired, alta performance)
- **ORM**: GORM (produtividade) ou ent (type-safety)
- **Validação**: go-playground/validator
- **Autenticação**: golang-jwt/jwt-go
- **Event Bus**: EventBus ou RabbitMQ
- **Testes**: testing package nativo + testify
- **Migrations**: golang-migrate
- **Config**: viper ou env
- **Logging**: zerolog (estruturado, alta performance)
- **Observabilidade**: OpenTelemetry, Prometheus, Grafana
- **Banco de Dados**: PostgreSQL (com connection pooling)
- **Cache**: Redis (go-redis) com cluster support
- **API Docs**: Swagger/OpenAPI (swaggo/swag)
- **Rate Limiting**: go-redis/redis-rate-limiter
- **Message Queue**: RabbitMQ ou NATS (para eventos assíncronos)
- **Monitoring**: Prometheus + Grafana
- **Tracing**: OpenTelemetry + Jaeger
- **Error Tracking**: Sentry (opcional)

### 3.2. Framework Web: Fiber

#### 3.2.1. Sobre o Fiber

**Fiber** ⚡ - Express-inspired web framework
- **GitHub**: 30k+ stars
- **Tipo**: Express-inspired web framework
- **Abordagem**: Familiar para quem vem do Node.js/Express
- **Base**: fasthttp (HTTP engine de alta performance)

#### 3.2.2. Características do Fiber

**Vantagens:**
- ✅ **API similar ao Express.js** - Fácil migração do Node.js
- ✅ **Performance excepcional** - Baseado em fasthttp (~200k req/s)
- ✅ **Middleware compatível** - Muitos middlewares disponíveis
- ✅ **Fácil de usar** - API intuitiva e familiar
- ✅ **Documentação boa** - Guias e exemplos completos
- ✅ **Type-safe** - Go oferece type safety nativo
- ✅ **Zero memory allocations** - Otimizado para performance
- ✅ **WebSocket support** - Suporte nativo a WebSockets
- ✅ **Template engine** - Suporte a templates

**Considerações:**
- ⚠️ **Usa fasthttp** - Não é net/http padrão do Go
- ⚠️ **Compatibilidade** - Algumas libs podem não funcionar (precisam net/http)
- ⚠️ **Comunidade menor** - Menor que alguns frameworks, mas crescente e ativa

#### 3.2.3. Performance do Fiber

**Benchmarks (aproximado):**

```
API Simples (Hello World):
- Fiber:      ~200.000 req/s (fasthttp)
- net/http:    ~50.000 req/s

API com JSON (CRUD):
- Fiber:      ~180.000 req/s
```

**Por que Fiber é mais rápido?**
- Usa `fasthttp` ao invés de `net/http` padrão
- `fasthttp` é otimizado para alta performance
- Menos alocações de memória
- Pooling de objetos reutilizáveis

#### 3.2.4. Exemplos de Código com Fiber

##### Exemplo Básico

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    app := fiber.New(fiber.Config{
        AppName: "Gestão Financeira API",
        ErrorHandler: customErrorHandler,
    })
    
    // Middleware global
    app.Use(logger.New())
    app.Use(recover.New())
    // CORS configurado para produção (não usar * em produção)
    app.Use(cors.New(cors.Config{
        AllowOrigins:     os.Getenv("ALLOWED_ORIGINS"), // Ex: "https://app.example.com"
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: true,
        MaxAge:           86400, // 24 horas
    }))
    
    // Health checks
    app.Get("/health", healthCheck)
    app.Get("/health/ready", readinessCheck) // Verifica DB, Redis, etc.
    app.Get("/health/live", livenessCheck)   // Verifica se app está vivo
    
    // API v1
    api := app.Group("/api/v1")
    {
        // Auth routes
        auth := api.Group("/auth")
        {
            auth.Post("/register", registerHandler)
            auth.Post("/login", loginHandler)
        }
        
        // Protected routes
        protected := api.Group("", authMiddleware)
        {
            // Transaction routes
            transactions := protected.Group("/transactions")
            {
                transactions.Get("", listTransactionsHandler)
                transactions.Post("", createTransactionHandler)
                transactions.Get("/:id", getTransactionHandler)
                transactions.Put("/:id", updateTransactionHandler)
                transactions.Delete("/:id", deleteTransactionHandler)
            }
        }
    }
    
    app.Listen(":8080")
}
```

##### Handler com Fiber

```go
// internal/identity/presentation/handlers/auth_handler.go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "gestao-financeira/internal/identity/application/usecases"
)

type AuthHandler struct {
    registerUserUseCase *usecases.RegisterUserUseCase
    loginUseCase        *usecases.LoginUseCase
}

func NewAuthHandler(
    registerUserUseCase *usecases.RegisterUserUseCase,
    loginUseCase *usecases.LoginUseCase,
) *AuthHandler {
    return &AuthHandler{
        registerUserUseCase: registerUserUseCase,
        loginUseCase:        loginUseCase,
    }
}

type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"firstName" validate:"required"`
    LastName  string `json:"lastName" validate:"required"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    
    // Parse e validação
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Validação com validator
    if err := validateStruct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Executar use case
    input := usecases.RegisterUserInput{
        Email:     req.Email,
        Password:  req.Password,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }
    
    output, err := h.registerUserUseCase.Execute(input)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(output)
}
```

#### 3.2.5. Middleware com Fiber

##### Auth Middleware

```go
// pkg/middleware/auth.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        // Remover "Bearer " prefix se existir
        if len(token) > 7 && token[:7] == "Bearer " {
            token = token[7:]
        }
        
        // Validar token JWT
        claims, err := validateToken(token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
        
        // Adicionar user ID ao locals (context do Fiber)
        c.Locals("userID", claims.UserID)
        
        return c.Next()
    }
}

func validateToken(tokenString string) (*jwt.MapClaims, error) {
    // Implementação da validação JWT
    // ...
    return claims, nil
}
```

##### Logger Middleware Customizado

```go
// pkg/middleware/logger.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func CustomLogger() fiber.Handler {
    return logger.New(logger.Config{
        Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
        TimeFormat: "2006-01-02 15:04:05",
        TimeZone:   "America/Sao_Paulo",
    })
}
```

#### 3.2.6. Estrutura Completa com Fiber

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/recover"
    
    "gestao-financeira/internal/identity/presentation/handlers"
    "gestao-financeira/internal/transaction/presentation/handlers"
    "gestao-financeira/pkg/middleware"
    "gestao-financeira/pkg/database"
)

func main() {
    // Inicializar banco de dados
    db := database.NewDatabase()
    defer db.Close()
    
    // Criar app Fiber
    app := fiber.New(fiber.Config{
        AppName:       "Gestão Financeira API",
        ServerHeader:  "Fiber",
        ErrorHandler:  customErrorHandler,
    })
    
    // Middleware global
    app.Use(logger.New())
    app.Use(recover.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "*",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: true,
    }))
    
    // Health check
    app.Get("/health", healthCheck)
    
    // API v1
    setupRoutes(app, db)
    
    // Graceful shutdown
    go func() {
        if err := app.Listen(":8080"); err != nil {
            log.Fatal("Failed to start server:", err)
        }
    }()
    
    // Aguardar sinal de interrupção
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    
    // Graceful shutdown
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := app.ShutdownWithContext(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    // Fechar conexões
    db.Close()
    log.Println("Server exited")
}

func setupRoutes(app *fiber.App, db *database.DB) {
    v1 := app.Group("/api/v1")
    
    // Auth routes (públicas)
    auth := v1.Group("/auth")
    {
        authHandler := handlers.NewAuthHandler(/* dependencies */)
        auth.Post("/register", authHandler.Register)
        auth.Post("/login", authHandler.Login)
    }
    
    // Rotas protegidas
    protected := v1.Group("", middleware.AuthMiddleware())
    {
        // Transactions
        transactions := protected.Group("/transactions")
        {
            transactionHandler := handlers.NewTransactionHandler(/* dependencies */)
            transactions.Get("", transactionHandler.List)
            transactions.Post("", transactionHandler.Create)
            transactions.Get("/:id", transactionHandler.Get)
            transactions.Put("/:id", transactionHandler.Update)
            transactions.Delete("/:id", transactionHandler.Delete)
        }
        
        // Accounts
        accounts := protected.Group("/accounts")
        {
            accountHandler := handlers.NewAccountHandler(/* dependencies */)
            accounts.Get("", accountHandler.List)
            accounts.Post("", accountHandler.Create)
            accounts.Get("/:id", accountHandler.Get)
        }
    }
}

func healthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status":  "ok",
        "service": "gestao-financeira",
        "version": "1.0.0",
    })
}

func customErrorHandler(c *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError
    message := "Internal server error"
    
    // Log erro interno (não expor detalhes ao cliente)
    if code == fiber.StatusInternalServerError {
        log.Error().
            Err(err).
            Str("path", c.Path()).
            Str("method", c.Method()).
            Msg("Internal server error")
        message = "An unexpected error occurred"
    }
    
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
        message = e.Message
    }
    
    // Não expor stack trace em produção
    response := fiber.Map{
        "error": message,
        "code":  code,
    }
    
    // Adicionar request ID para rastreamento
    if requestID := c.Locals("requestID"); requestID != nil {
        response["request_id"] = requestID
    }
    
    return c.Status(code).JSON(response)
}
```

#### 3.2.7. Validação com Fiber

```go
// pkg/validator/validator.go
package validator

import (
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}

// Middleware de validação
func ValidateRequest(req interface{}) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if err := c.BodyParser(req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid request body",
            })
        }
        
        if err := ValidateStruct(req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        
        c.Locals("validatedRequest", req)
        return c.Next()
    }
}
```

#### 3.2.8. Vantagens do Fiber para este Projeto

**Por que Fiber é uma excelente escolha:**
1. ✅ **API Familiar** - Se você conhece Express.js, Fiber é muito similar
2. ✅ **Performance Máxima** - Uma das opções mais rápidas disponíveis
3. ✅ **Fácil de Aprender** - API intuitiva e bem documentada
4. ✅ **Middleware Rico** - Muitos middlewares disponíveis
5. ✅ **Compatível com DDD** - Não impõe estrutura, você controla
6. ✅ **Type Safety** - Go oferece type safety nativo
7. ✅ **WebSocket Support** - Para features real-time no futuro

**Considerações:**
- ⚠️ **fasthttp vs net/http** - Algumas bibliotecas podem não funcionar
- ⚠️ **Comunidade** - Menor que Gin, mas crescente e ativa

#### 3.2.9. Recursos Adicionais do Fiber

**Middleware Disponíveis:**
- `logger` - Logging de requisições
- `cors` - CORS support
- `recover` - Panic recovery
- `limiter` - Rate limiting
- `compress` - Compression
- `helmet` - Security headers
- `jwt` - JWT authentication
- `session` - Session management
- `cache` - HTTP caching

**Exemplo com Middleware Adicional:**

```go
import (
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/helmet"
)

app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
}))

app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,
}))

app.Use(helmet.New())
```

### 3.3. Por que Go?

**Vantagens:**
- ✅ **Performance excepcional**: Compilado, muito rápido
- ✅ **Concorrência nativa**: Goroutines são incríveis
- ✅ **Baixo consumo de memória**: Eficiente
- ✅ **Type safety**: Forte e estático
- ✅ **Simplicidade**: Linguagem simples e direta
- ✅ **Deploy**: Binário único, fácil deploy
- ✅ **Escalabilidade**: Excelente para alta carga
- ✅ **Aprendizado**: Linguagem moderna e valorizada

**Desafios:**
- ⚠️ **Curva de aprendizado**: Inicial, mas Go é simples
- ⚠️ **Ecossistema menor**: Menos pacotes que PHP/Node
- ⚠️ **DDD menos comum**: Menos exemplos/práticas
- ⚠️ **Error handling**: Explícito (pode ser verboso)

## 4. Stack Tecnológico Frontend

### 4.1. Abordagem de Desenvolvimento

O frontend será desenvolvido **incrementalmente**, sempre que um módulo do backend estiver pronto. Esta abordagem permite:

- ✅ **Validação rápida**: Testar funcionalidades assim que o backend estiver disponível
- ✅ **Feedback contínuo**: Identificar problemas de integração cedo
- ✅ **Desenvolvimento paralelo**: Frontend e backend podem evoluir juntos
- ✅ **MVP funcional**: Ter um produto funcional desde as primeiras fases

**Estratégia:**
- Cada bounded context do backend terá sua interface correspondente no frontend
- Desenvolvimento frontend inicia após a API estar documentada (Swagger)
- Componentes são construídos de forma modular e reutilizável
- Testes de integração frontend-backend são realizados continuamente

### 4.2. Tecnologias Principais

- **Framework**: **Next.js 14+** (React 18+ com App Router)
- **Linguagem**: **TypeScript 5+**
- **UI Library**: **shadcn/ui** (componentes acessíveis e customizáveis)
- **Styling**: **Tailwind CSS** (utility-first)
- **State Management**: **Zustand** (leve) ou **TanStack Query (React Query)** (server state)
- **Formulários**: **React Hook Form** + **Zod** (validação type-safe)
- **Gráficos**: **Recharts** ou **Chart.js** (visualizações financeiras)
- **Roteamento**: Next.js App Router (nativo)
- **Autenticação**: JWT tokens (armazenados em httpOnly cookies)
- **API Client**: **Axios** ou **fetch** nativo com wrappers
- **Testes**: **Vitest** + **React Testing Library**
- **E2E**: **Playwright** ou **Cypress**
- **Build**: Next.js (otimizado para produção)

### 4.3. Por que Next.js?

**Vantagens:**
- ✅ **SSR/SSG**: Performance e SEO excelentes
- ✅ **App Router**: Roteamento moderno e intuitivo
- ✅ **TypeScript nativo**: Type-safety completo
- ✅ **API Routes**: Pode servir como proxy se necessário
- ✅ **Otimizações automáticas**: Image optimization, code splitting
- ✅ **Deploy simples**: Vercel, Netlify, ou qualquer servidor Node.js
- ✅ **Ecossistema**: Grande comunidade e pacotes

### 4.4. Por que shadcn/ui?

**shadcn/ui** é uma coleção de componentes React reutilizáveis construídos com Radix UI e Tailwind CSS.

**Vantagens:**
- ✅ **Acessibilidade**: Baseado em Radix UI (WAI-ARIA compliant)
- ✅ **Customizável**: Código copiado para o projeto (não é uma dependência)
- ✅ **Type-safe**: TypeScript completo
- ✅ **Tailwind CSS**: Estilização com utility classes
- ✅ **Modular**: Use apenas os componentes que precisa
- ✅ **Bem documentado**: Exemplos claros e código limpo
- ✅ **Design system**: Componentes consistentes e profissionais
- ✅ **Dark mode**: Suporte nativo a temas

**Componentes principais que usaremos:**
- `Button`, `Input`, `Select`, `Dialog`, `Dropdown Menu`
- `Table`, `Card`, `Tabs`, `Form`, `Toast`
- `Calendar`, `Date Picker`, `Chart` (com Recharts)
- `Skeleton` (loading states)
- `Alert`, `Badge`, `Avatar`

### 4.5. Estrutura de Pastas Frontend

```
frontend/
├── app/                                    # Next.js App Router
│   ├── (auth)/                             # Grupo de rotas (auth)
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   ├── (dashboard)/                        # Grupo de rotas (dashboard)
│   │   ├── layout.tsx                      # Layout do dashboard
│   │   ├── page.tsx                        # Dashboard home
│   │   ├── accounts/
│   │   │   ├── page.tsx                    # Lista de contas
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx                # Detalhes da conta
│   │   │   └── new/
│   │   │       └── page.tsx                # Criar conta
│   │   ├── transactions/
│   │   │   ├── page.tsx                    # Lista de transações
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx                # Detalhes da transação
│   │   │   └── new/
│   │   │       └── page.tsx                # Criar transação
│   │   ├── categories/
│   │   │   └── page.tsx
│   │   ├── budget/
│   │   │   └── page.tsx
│   │   ├── reports/
│   │   │   └── page.tsx
│   │   └── settings/
│   │       └── page.tsx
│   ├── api/                                # API Routes (se necessário)
│   │   └── auth/
│   │       └── callback/
│   │           └── route.ts
│   ├── layout.tsx                          # Layout raiz
│   ├── page.tsx                            # Home (redirect)
│   └── globals.css                         # Estilos globais
│
├── components/                             # Componentes React
│   ├── ui/                                 # Componentes shadcn/ui
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── card.tsx
│   │   ├── table.tsx
│   │   └── ...                             # Outros componentes
│   ├── features/                           # Componentes por feature
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   └── RegisterForm.tsx
│   │   ├── accounts/
│   │   │   ├── AccountList.tsx
│   │   │   ├── AccountForm.tsx
│   │   │   └── AccountCard.tsx
│   │   ├── transactions/
│   │   │   ├── TransactionList.tsx
│   │   │   ├── TransactionForm.tsx
│   │   │   ├── TransactionTable.tsx
│   │   │   └── TransactionFilters.tsx
│   │   ├── categories/
│   │   │   └── ...
│   │   ├── budget/
│   │   │   └── ...
│   │   └── reports/
│   │       ├── ReportChart.tsx
│   │       └── ReportTable.tsx
│   ├── layout/                             # Componentes de layout
│   │   ├── Header.tsx
│   │   ├── Sidebar.tsx
│   │   ├── Footer.tsx
│   │   └── Navbar.tsx
│   └── shared/                             # Componentes compartilhados
│       ├── LoadingSpinner.tsx
│       ├── ErrorBoundary.tsx
│       └── EmptyState.tsx
│
├── lib/                                    # Utilitários e configurações
│   ├── api/                                # Cliente API
│   │   ├── client.ts                       # Axios instance
│   │   ├── endpoints.ts                    # Endpoints da API
│   │   └── types.ts                        # Types da API
│   ├── hooks/                              # Custom hooks
│   │   ├── useAuth.ts
│   │   ├── useTransactions.ts
│   │   └── useAccounts.ts
│   ├── utils/                              # Funções utilitárias
│   │   ├── formatters.ts                   # Formatação de valores
│   │   ├── validators.ts                   # Validações
│   │   └── constants.ts
│   ├── store/                              # State management (Zustand)
│   │   ├── authStore.ts
│   │   └── uiStore.ts
│   └── config/                             # Configurações
│       ├── env.ts                          # Variáveis de ambiente
│       └── theme.ts                        # Configuração de tema
│
├── types/                                  # TypeScript types
│   ├── api.ts                              # Types da API
│   ├── domain.ts                           # Types do domínio
│   └── index.ts
│
├── styles/                                 # Estilos adicionais
│   └── components.css
│
├── public/                                 # Arquivos estáticos
│   ├── images/
│   └── icons/
│
├── tests/                                  # Testes
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── .env.local                              # Variáveis de ambiente
├── .env.example
├── next.config.js
├── tailwind.config.js
├── tsconfig.json
├── package.json
└── README.md
```

### 4.6. Integração com a API Backend

#### 4.6.1. Cliente API

```typescript
// lib/api/client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Para cookies httpOnly
});

// Interceptor para adicionar token
apiClient.interceptors.request.use((config) => {
  const token = getAuthToken(); // Função para pegar token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor para tratar erros
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirecionar para login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

#### 4.6.2. Custom Hooks

```typescript
// lib/hooks/useTransactions.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/lib/api/client';

export function useTransactions(filters?: TransactionFilters) {
  return useQuery({
    queryKey: ['transactions', filters],
    queryFn: () => apiClient.get('/transactions', { params: filters }),
  });
}

export function useCreateTransaction() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateTransactionDTO) =>
      apiClient.post('/transactions', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] });
    },
  });
}
```

#### 4.6.3. Componente de Formulário com React Hook Form + Zod

```typescript
// components/features/transactions/TransactionForm.tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useCreateTransaction } from '@/lib/hooks/useTransactions';

const transactionSchema = z.object({
  amount: z.number().positive(),
  description: z.string().min(1),
  accountId: z.string().uuid(),
  categoryId: z.string().uuid(),
  date: z.string(),
  type: z.enum(['income', 'expense']),
});

type TransactionFormData = z.infer<typeof transactionSchema>;

export function TransactionForm() {
  const { register, handleSubmit, formState: { errors } } = useForm<TransactionFormData>({
    resolver: zodResolver(transactionSchema),
  });
  
  const createTransaction = useCreateTransaction();
  
  const onSubmit = (data: TransactionFormData) => {
    createTransaction.mutate(data);
  };
  
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* Campos do formulário */}
    </form>
  );
}
```

### 4.7. Desenvolvimento Incremental por Módulo

#### Mapeamento Backend → Frontend

| Backend Context | Frontend Módulo | Prioridade |
|----------------|-----------------|------------|
| **Identity** | Autenticação (Login/Register) | Fase 1 |
| **Account Management** | Gestão de Contas | Fase 1 |
| **Transaction** | Gestão de Transações | Fase 1 |
| **Category** | Gestão de Categorias | Fase 2 |
| **Budget** | Planejamento Orçamentário | Fase 3 |
| **Reporting** | Relatórios e Gráficos | Fase 3 |
| **Investment** | Gestão de Investimentos | Fase 5 |
| **Goal** | Metas Financeiras | Fase 5 |
| **Notification** | Notificações | Fase 5 |

#### Fluxo de Desenvolvimento

1. **Backend disponibiliza API** → Swagger documentado
2. **Frontend gera types** → A partir do Swagger (swagger-typescript-api)
3. **Frontend cria componentes** → Formulários, listas, detalhes
4. **Integração e testes** → Testes E2E da integração
5. **Refinamento** → Ajustes baseados em feedback

### 4.8. Exemplo: Módulo de Transações

#### Estrutura do Módulo

```
components/features/transactions/
├── TransactionList.tsx          # Lista de transações
├── TransactionTable.tsx         # Tabela com paginação
├── TransactionForm.tsx         # Formulário de criação/edição
├── TransactionFilters.tsx       # Filtros (data, tipo, categoria)
├── TransactionCard.tsx          # Card para visualização mobile
└── TransactionDetails.tsx       # Modal/drawer com detalhes
```

#### Integração com Backend

```typescript
// Quando o backend Transaction Context estiver pronto:
// 1. API disponível em /api/v1/transactions
// 2. Swagger documentado
// 3. Frontend consome e exibe dados
// 4. Formulários validam e enviam dados
// 5. Atualização de saldo reflete em tempo real
```

## 5. Arquitetura DDD em Go

### 5.1. Estrutura em Camadas

```
┌─────────────────────────────────────┐
│     Presentation Layer              │  (Handlers, DTOs, HTTP)
├─────────────────────────────────────┤
│     Application Layer                │  (Use Cases, Application Services)
├─────────────────────────────────────┤
│     Domain Layer                    │  (Entities, Value Objects, Domain Services)
├─────────────────────────────────────┤
│     Infrastructure Layer            │  (Repositories, External Services, DB)
└─────────────────────────────────────┘
```

### 5.2. Bounded Contexts

1. **Identity Context** - Autenticação e gestão de usuários
2. **Account Management Context** - Gestão de contas e carteiras
3. **Transaction Context** - Processamento de transações financeiras
4. **Category Context** - Gestão de categorias e taxonomia
5. **Budget Context** - Planejamento e controle orçamentário
6. **Reporting Context** - Análises e relatórios financeiros
7. **Investment Context** - Gestão de investimentos
8. **Goal Context** - Metas e objetivos financeiros
9. **Notification Context** - Notificações e alertas

## 6. Estrutura de Pastas (Go DDD)

**Nota:** O frontend está em um diretório separado (`frontend/`). Ver seção 4.5 para estrutura completa do frontend.

```
gestao-financeira-go/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── shared/                          # Shared Kernel
│   │   ├── domain/
│   │   │   ├── valueobjects/
│   │   │   │   ├── money.go
│   │   │   │   ├── currency.go
│   │   │   │   └── account_context.go
│   │   │   └── events/
│   │   └── infrastructure/
│   │       └── eventbus/
│   │
│   ├── identity/                         # Identity Context
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   └── user.go
│   │   │   ├── valueobjects/
│   │   │   │   ├── email.go
│   │   │   │   ├── password_hash.go
│   │   │   │   └── user_name.go
│   │   │   ├── services/
│   │   │   │   ├── password_service.go
│   │   │   │   └── token_service.go
│   │   │   ├── repositories/
│   │   │   │   └── user_repository.go
│   │   │   └── events/
│   │   │       ├── user_registered.go
│   │   │       └── user_password_changed.go
│   │   ├── application/
│   │   │   ├── usecases/
│   │   │   │   ├── register_user.go
│   │   │   │   ├── login.go
│   │   │   │   └── update_profile.go
│   │   │   └── dtos/
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   └── gorm_user_repository.go
│   │   │   ├── services/
│   │   │   │   └── jwt_token_service.go
│   │   │   └── events/
│   │   │       └── domain_event_publisher.go
│   │   └── presentation/
│   │       ├── handlers/
│   │       │   └── auth_handler.go
│   │       └── dtos/
│   │
│   ├── account/                          # Account Management Context
│   ├── transaction/                      # Transaction Context
│   ├── category/                         # Category Context
│   ├── budget/                           # Budget Context
│   ├── reporting/                        # Reporting Context
│   ├── investment/                      # Investment Context
│   ├── goal/                             # Goal Context
│   └── notification/                     # Notification Context
│
├── pkg/                                  # Pacotes compartilhados
│   ├── database/
│   ├── logger/
│   └── validator/
│
├── migrations/                           # Migrations do banco
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── docs/                                 # Documentação Swagger
│   └── swagger.yaml
│
├── scripts/                              # Scripts utilitários
│   ├── backup.sh
│   ├── migrate.sh
│   └── seed.sh
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── docker-compose.prod.yml
├── .env.example
└── README.md
```

## 7. Detalhamento dos Bounded Contexts (Go)

### 7.1. Identity Context

#### 7.1.1. Entidades (Go)

**User (Agregado Raiz)**
```go
// internal/identity/domain/entities/user.go
package entities

import (
    "time"
    "github.com/google/uuid"
    "gestao-financeira/internal/identity/domain/valueobjects"
    "gestao-financeira/internal/shared/domain/events"
)

type User struct {
    id          valueobjects.UserID
    email       valueobjects.Email
    passwordHash valueobjects.PasswordHash
    name        valueobjects.UserName
    profile     valueobjects.UserProfile
    createdAt   time.Time
    updatedAt   time.Time
    isActive    bool
    events      []events.DomainEvent
}

func NewUser(
    email valueobjects.Email,
    passwordHash valueobjects.PasswordHash,
    name valueobjects.UserName,
) *User {
    return &User{
        id:          valueobjects.NewUserID(uuid.New()),
        email:       email,
        passwordHash: passwordHash,
        name:        name,
        profile:     valueobjects.NewUserProfile(),
        createdAt:   time.Now(),
        updatedAt:   time.Now(),
        isActive:    true,
        events:      []events.DomainEvent{},
    }
}

func (u *User) ChangePassword(oldPassword, newPassword string) error {
    if !u.passwordHash.Verify(oldPassword) {
        return errors.New("invalid old password")
    }
    
    newHash, err := valueobjects.NewPasswordHashFromPlain(newPassword)
    if err != nil {
        return err
    }
    
    u.passwordHash = newHash
    u.updatedAt = time.Now()
    u.addEvent(events.NewUserPasswordChanged(u.id))
    
    return nil
}

func (u *User) UpdateProfile(profile valueobjects.UserProfile) {
    u.profile = profile
    u.updatedAt = time.Now()
    u.addEvent(events.NewUserProfileUpdated(u.id))
}

func (u *User) Deactivate() {
    u.isActive = false
    u.updatedAt = time.Now()
    u.addEvent(events.NewUserDeactivated(u.id))
}

func (u *User) GetID() valueobjects.UserID {
    return u.id
}

func (u *User) GetEmail() valueobjects.Email {
    return u.email
}

func (u *User) GetEvents() []events.DomainEvent {
    return u.events
}

func (u *User) ClearEvents() {
    u.events = []events.DomainEvent{}
}

func (u *User) addEvent(event events.DomainEvent) {
    u.events = append(u.events, event)
}
```

#### 7.1.2. Value Objects (Go)

**Email**
```go
// internal/identity/domain/valueobjects/email.go
package valueobjects

import (
    "strings"
    "regexp"
    "errors"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    email = strings.ToLower(strings.TrimSpace(email))
    
    if !emailRegex.MatchString(email) {
        return Email{}, errors.New("invalid email format")
    }
    
    return Email{value: email}, nil
}

func (e Email) Value() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

func (e Email) String() string {
    return e.value
}
```

**PasswordHash**
```go
// internal/identity/domain/valueobjects/password_hash.go
package valueobjects

import (
    "golang.org/x/crypto/bcrypt"
    "errors"
)

type PasswordHash struct {
    value string
}

func NewPasswordHashFromPlain(password string) (PasswordHash, error) {
    if len(password) < 8 {
        return PasswordHash{}, errors.New("password must be at least 8 characters")
    }
    
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return PasswordHash{}, err
    }
    
    return PasswordHash{value: string(hash)}, nil
}

func NewPasswordHashFromHash(hash string) PasswordHash {
    return PasswordHash{value: hash}
}

func (p PasswordHash) Verify(plainPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(p.value), []byte(plainPassword))
    return err == nil
}

func (p PasswordHash) Value() string {
    return p.value
}
```

#### 7.1.3. Repositórios (Go)

**Interface**
```go
// internal/identity/domain/repositories/user_repository.go
package repositories

import (
    "gestao-financeira/internal/identity/domain/entities"
    "gestao-financeira/internal/identity/domain/valueobjects"
)

type UserRepository interface {
    FindByID(id valueobjects.UserID) (*entities.User, error)
    FindByEmail(email valueobjects.Email) (*entities.User, error)
    Save(user *entities.User) error
    Delete(id valueobjects.UserID) error
    Exists(email valueobjects.Email) (bool, error)
}
```

**Implementação GORM**
```go
// internal/identity/infrastructure/persistence/gorm_user_repository.go
package persistence

import (
    "gorm.io/gorm"
    "gestao-financeira/internal/identity/domain/entities"
    "gestao-financeira/internal/identity/domain/repositories"
    "gestao-financeira/internal/identity/domain/valueobjects"
)

type GormUserRepository struct {
    db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
    return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
    var model UserModel
    if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    
    return r.toDomain(&model)
}

func (r *GormUserRepository) FindByEmail(email valueobjects.Email) (*entities.User, error) {
    var model UserModel
    if err := r.db.Where("email = ?", email.Value()).First(&model).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    
    return r.toDomain(&model)
}

func (r *GormUserRepository) Save(user *entities.User) error {
    model := r.toModel(user)
    return r.db.Save(&model).Error
}

func (r *GormUserRepository) Delete(id valueobjects.UserID) error {
    return r.db.Delete(&UserModel{}, "id = ?", id.Value()).Error
}

func (r *GormUserRepository) Exists(email valueobjects.Email) (bool, error) {
    var count int64
    err := r.db.Model(&UserModel{}).Where("email = ?", email.Value()).Count(&count).Error
    return count > 0, err
}

// UserModel - Modelo de persistência
type UserModel struct {
    ID          string    `gorm:"type:uuid;primary_key"`
    Email       string    `gorm:"type:varchar(255);unique_index;not null"`
    PasswordHash string   `gorm:"type:varchar(255);not null"`
    FirstName   string    `gorm:"type:varchar(100)"`
    LastName    string    `gorm:"type:varchar(100)"`
    Currency    string    `gorm:"type:varchar(3);default:'BRL'"`
    Locale      string    `gorm:"type:varchar(10);default:'pt-BR'"`
    Theme       string    `gorm:"type:varchar(10);default:'light'"`
    IsActive    bool      `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (UserModel) TableName() string {
    return "users"
}

// Mapeamento entre domínio e persistência
func (r *GormUserRepository) toDomain(model *UserModel) (*entities.User, error) {
    email, err := valueobjects.NewEmail(model.Email)
    if err != nil {
        return nil, err
    }
    
    passwordHash := valueobjects.NewPasswordHashFromHash(model.PasswordHash)
    name := valueobjects.NewUserName(model.FirstName, model.LastName)
    
    user := entities.NewUser(email, passwordHash, name)
    // ... resto do mapeamento
    
    return user, nil
}

func (r *GormUserRepository) toModel(user *entities.User) *UserModel {
    // Mapeamento de domínio para persistência
    return &UserModel{
        ID:          user.GetID().Value(),
        Email:       user.GetEmail().Value(),
        // ... resto do mapeamento
    }
}
```

#### 7.1.4. Use Cases (Go)

**RegisterUserUseCase**
```go
// internal/identity/application/usecases/register_user.go
package usecases

import (
    "errors"
    "gestao-financeira/internal/identity/domain/entities"
    "gestao-financeira/internal/identity/domain/repositories"
    "gestao-financeira/internal/identity/domain/valueobjects"
    "gestao-financeira/internal/shared/domain/events"
)

type RegisterUserUseCase struct {
    userRepo repositories.UserRepository
    eventBus events.EventBus
}

func NewRegisterUserUseCase(
    userRepo repositories.UserRepository,
    eventBus events.EventBus,
) *RegisterUserUseCase {
    return &RegisterUserUseCase{
        userRepo: userRepo,
        eventBus: eventBus,
    }
}

type RegisterUserInput struct {
    Email    string
    Password string
    FirstName string
    LastName  string
}

type RegisterUserOutput struct {
    UserID string
    Email  string
}

func (uc *RegisterUserUseCase) Execute(input RegisterUserInput) (*RegisterUserOutput, error) {
    email, err := valueobjects.NewEmail(input.Email)
    if err != nil {
        return nil, err
    }
    
    exists, err := uc.userRepo.Exists(email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("user already exists")
    }
    
    passwordHash, err := valueobjects.NewPasswordHashFromPlain(input.Password)
    if err != nil {
        return nil, err
    }
    
    name := valueobjects.NewUserName(input.FirstName, input.LastName)
    user := entities.NewUser(email, passwordHash, name)
    
    if err := uc.userRepo.Save(user); err != nil {
        return nil, err
    }
    
    // Publicar eventos de domínio
    for _, event := range user.GetEvents() {
        uc.eventBus.Publish(event)
    }
    user.ClearEvents()
    
    return &RegisterUserOutput{
        UserID: user.GetID().Value(),
        Email:  user.GetEmail().Value(),
    }, nil
}
```

#### 7.1.5. Handlers (Fiber)

**AuthHandler**
```go
// internal/identity/presentation/handlers/auth_handler.go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "gestao-financeira/internal/identity/application/usecases"
    "gestao-financeira/pkg/validator"
)

type AuthHandler struct {
    registerUserUseCase *usecases.RegisterUserUseCase
    loginUseCase        *usecases.LoginUseCase
}

func NewAuthHandler(
    registerUserUseCase *usecases.RegisterUserUseCase,
    loginUseCase *usecases.LoginUseCase,
) *AuthHandler {
    return &AuthHandler{
        registerUserUseCase: registerUserUseCase,
        loginUseCase:        loginUseCase,
    }
}

type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"firstName" validate:"required"`
    LastName  string `json:"lastName" validate:"required"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    
    // Parse body
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Validação
    if err := validator.ValidateStruct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    input := usecases.RegisterUserInput{
        Email:     req.Email,
        Password:  req.Password,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }
    
    output, err := h.registerUserUseCase.Execute(input)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(output)
}
```

### 7.2. Transaction Context (Core Domain)

#### 7.2.1. Entidade Transaction (Go)

```go
// internal/transaction/domain/entities/transaction.go
package entities

import (
    "time"
    "errors"
    "gestao-financeira/internal/shared/domain/valueobjects"
    "gestao-financeira/internal/shared/domain/events"
)

type Transaction struct {
    id          valueobjects.TransactionID
    userId      valueobjects.UserID
    accountId   valueobjects.AccountID
    categoryId  valueobjects.CategoryID
    type_       valueobjects.TransactionType
    amount      valueobjects.Money
    description valueobjects.TransactionDescription
    date        time.Time
    tags        []valueobjects.Tag
    status      valueobjects.TransactionStatus
    context     valueobjects.AccountContext
    createdAt   time.Time
    updatedAt   time.Time
    events      []events.DomainEvent
}

func NewTransaction(
    userId valueobjects.UserID,
    accountId valueobjects.AccountID,
    categoryId valueobjects.CategoryID,
    type_ valueobjects.TransactionType,
    amount valueobjects.Money,
    description valueobjects.TransactionDescription,
    date time.Time,
    context valueobjects.AccountContext,
) *Transaction {
    return &Transaction{
        id:          valueobjects.NewTransactionID(),
        userId:      userId,
        accountId:   accountId,
        categoryId:  categoryId,
        type_:       type_,
        amount:      amount,
        description: description,
        date:        date,
        tags:        []valueobjects.Tag{},
        status:      valueobjects.TransactionStatusPending,
        context:     context,
        createdAt:   time.Now(),
        updatedAt:   time.Now(),
        events:      []events.DomainEvent{},
    }
}

func (t *Transaction) Approve() error {
    if !t.status.CanBeApproved() {
        return errors.New("transaction cannot be approved")
    }
    
    t.status = valueobjects.TransactionStatusApproved
    t.updatedAt = time.Now()
    t.addEvent(events.NewTransactionApproved(t.id, t.amount))
    
    return nil
}

func (t *Transaction) Cancel() error {
    if !t.status.CanBeCancelled() {
        return errors.New("transaction cannot be cancelled")
    }
    
    t.status = valueobjects.TransactionStatusCancelled
    t.updatedAt = time.Now()
    t.addEvent(events.NewTransactionCancelled(t.id))
    
    return nil
}

func (t *Transaction) GetID() valueobjects.TransactionID {
    return t.id
}

func (t *Transaction) GetAmount() valueobjects.Money {
    return t.amount
}

func (t *Transaction) GetEvents() []events.DomainEvent {
    return t.events
}

func (t *Transaction) ClearEvents() {
    t.events = []events.DomainEvent{}
}

func (t *Transaction) addEvent(event events.DomainEvent) {
    t.events = append(t.events, event)
}
```

## 8. ORM: GORM vs ent

### 8.1. GORM (Recomendado para Produtividade)

**Vantagens:**
- ✅ Mais popular e maduro
- ✅ Migrations automáticas
- ✅ Hooks (BeforeCreate, AfterUpdate)
- ✅ Associations fáceis
- ✅ Query builder intuitivo
- ✅ Documentação excelente

**Exemplo:**
```go
type AccountModel struct {
    ID        string    `gorm:"type:uuid;primary_key"`
    UserID    string    `gorm:"type:uuid;index"`
    Name      string    `gorm:"type:varchar(255)"`
    Type      string    `gorm:"type:varchar(50)"`
    Balance   float64   `gorm:"type:decimal(15,2)"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (AccountModel) TableName() string {
    return "accounts"
}
```

### 8.2. ent (Recomendado para Type-Safety)

**Vantagens:**
- ✅ Type-safe por geração de código
- ✅ Schema-first (similar ao Prisma)
- ✅ Excelente type-safety
- ✅ Migrations automáticas
- ✅ GraphQL integrado (opcional)

**Exemplo:**
```go
// schema/account.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type Account struct {
    ent.Schema
}

func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.String("user_id"),
        field.String("name"),
        field.String("type"),
        field.Float("balance"),
        field.Time("created_at").Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}
```

## 9. Event Bus em Go

### 9.1. Implementação Simples

```go
// internal/shared/infrastructure/eventbus/simple_event_bus.go
package eventbus

import (
    "sync"
    "gestao-financeira/internal/shared/domain/events"
)

type SimpleEventBus struct {
    handlers map[string][]events.EventHandler
    mu       sync.RWMutex
}

func NewSimpleEventBus() *SimpleEventBus {
    return &SimpleEventBus{
        handlers: make(map[string][]events.EventHandler),
    }
}

func (b *SimpleEventBus) Subscribe(eventType string, handler events.EventHandler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *SimpleEventBus) Publish(event events.DomainEvent) {
    b.mu.RLock()
    handlers := b.handlers[event.EventType()]
    b.mu.RUnlock()
    
    for _, handler := range handlers {
        go handler.Handle(event) // Goroutine para não bloquear
    }
}
```

## 10. Testes em Go

### 10.1. Testes Unitários

```go
// internal/transaction/domain/entities/transaction_test.go
package entities_test

import (
    "testing"
    "time"
    "gestao-financeira/internal/transaction/domain/entities"
    "gestao-financeira/internal/shared/domain/valueobjects"
)

func TestTransaction_Approve(t *testing.T) {
    transaction := entities.NewTransaction(
        valueobjects.NewUserID(),
        valueobjects.NewAccountID(),
        valueobjects.NewCategoryID(),
        valueobjects.TransactionTypeExpense,
        valueobjects.NewMoney(100.0, valueobjects.CurrencyBRL),
        valueobjects.NewTransactionDescription("Test"),
        time.Now(),
        valueobjects.AccountContextPersonal,
    )
    
    err := transaction.Approve()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    // Verificar eventos
    events := transaction.GetEvents()
    if len(events) == 0 {
        t.Fatal("Expected domain event to be added")
    }
}
```

### 10.2. Testes de Integração

```go
// tests/integration/transaction_repository_test.go
package integration_test

import (
    "testing"
    "gestao-financeira/internal/transaction/infrastructure/persistence"
)

func TestGormTransactionRepository_Save(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    repo := persistence.NewGormTransactionRepository(db)
    
    // Test
    transaction := createTestTransaction(t)
    err := repo.Save(transaction)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
}
```

## 11. Fases de Desenvolvimento (Backend + Frontend)

### Fase 1: Fundação e MVP (3-4 semanas)

#### Backend:
- Setup do projeto Go + Fiber
- Shared Kernel (Money, Currency, etc.)
- Identity Context (registro, login, JWT)
- Account Management Context
- Transaction Context (CRUD básico)
- Health checks básicos
- Testes unitários básicos
- Docker setup
- Swagger/OpenAPI básico

#### Frontend (Incremental):
- Setup Next.js + TypeScript + Tailwind CSS
- Instalação e configuração do shadcn/ui
- Layout base (Header, Sidebar, Footer)
- **Módulo de Autenticação:**
  - Página de Login
  - Página de Registro
  - Integração com API de autenticação
  - Proteção de rotas
- **Módulo de Contas:**
  - Lista de contas
  - Formulário de criação/edição
  - Integração com Account Management API
- **Módulo de Transações:**
  - Lista de transações
  - Formulário de criação/edição
  - Integração com Transaction API

**Entregável:** Sistema completo (backend + frontend) onde usuário pode registrar, fazer login, criar contas e transações via interface web

### Fase 2: Core Domain e Integrações (3-4 semanas)

#### Backend:
- Integração Transaction ↔ Account (atualização de saldo)
- Event Bus e Domain Events
- Category Context
- Validações robustas
- Error handling melhorado
- Testes de integração
- Logging estruturado

#### Frontend (Incremental):
- **Módulo de Categorias:**
  - Lista de categorias
  - Formulário de criação/edição
  - Seleção de categoria em transações
  - Integração com Category API
- Melhorias nos módulos existentes:
  - Atualização de saldo em tempo real
  - Filtros avançados em transações
  - Paginação e ordenação
  - Loading states e error handling
- Componentes compartilhados:
  - Toast notifications
  - Confirmation dialogs
  - Empty states

**Entregável:** Sistema funcional com categorias, eventos e interface completa para gestão de categorias

### Fase 3: Funcionalidades Essenciais (4-5 semanas)

#### Backend:
- Budget Context
- Recurring Transactions
- Reporting Context (relatórios básicos)
- Cache com Redis
- Paginação
- Rate limiting
- Swagger/OpenAPI completo
- Testes E2E

#### Frontend (Incremental):
- **Módulo de Orçamento:**
  - Dashboard de orçamentos
  - Formulário de criação/edição de orçamentos
  - Visualização de progresso
  - Alertas de limite
  - Integração com Budget API
- **Módulo de Relatórios:**
  - Dashboard com gráficos (Recharts)
  - Relatórios mensais/anuais
  - Filtros por período, categoria, tipo
  - Exportação de dados (CSV/PDF)
  - Visualizações:
    - Gráfico de receitas vs despesas
    - Gráfico por categoria
    - Tendências temporais
  - Integração com Reporting API
- Melhorias gerais:
  - Dark mode (shadcn/ui)
  - Responsividade mobile
  - Performance (lazy loading, code splitting)
  - Acessibilidade (ARIA labels)

**Entregável:** Sistema completo com orçamentos, relatórios visuais e interface rica para análise financeira

### Fase 4: Produção e Performance (3-4 semanas)

#### Backend:
- Observabilidade (métricas, tracing)
- Monitoramento (Prometheus, Grafana)
- Segurança robusta (headers, validações)
- Graceful shutdown
- CI/CD pipeline
- Backup automático
- Documentação completa
- Otimizações de performance

#### Frontend:
- **Otimizações:**
  - Code splitting avançado
  - Image optimization
  - Bundle size optimization
  - Caching estratégico
- **Testes:**
  - Testes unitários (Vitest)
  - Testes de integração
  - Testes E2E (Playwright)
  - Testes de acessibilidade
- **Deploy:**
  - Configuração de produção
  - Variáveis de ambiente
  - CI/CD para frontend
  - Deploy em Vercel/Netlify
- **Melhorias:**
  - Error tracking (Sentry)
  - Analytics (opcional)
  - PWA (Progressive Web App)
  - Offline support

**Entregável:** Sistema completo pronto para produção (backend + frontend) com monitoramento, testes e deploy automatizado

### Fase 5: Funcionalidades Avançadas (4-5 semanas)

#### Backend:
- Investment Context
- Goal Context
- Notification Context
- Dashboard completo (API)
- Exportação de dados
- Auditoria e compliance
- Multi-tenancy (se necessário)

#### Frontend (Incremental):
- **Módulo de Investimentos:**
  - Lista de investimentos
  - Formulário de criação/edição
  - Acompanhamento de performance
  - Gráficos de evolução
  - Integração com Investment API
- **Módulo de Metas:**
  - Lista de metas financeiras
  - Formulário de criação/edição
  - Progresso visual
  - Alertas de conquista
  - Integração com Goal API
- **Módulo de Notificações:**
  - Centro de notificações
  - Notificações em tempo real (WebSocket)
  - Preferências de notificação
  - Integração com Notification API
- **Dashboard Completo:**
  - Visão geral financeira
  - Cards com métricas principais
  - Gráficos interativos
  - Resumo de contas
  - Últimas transações
  - Metas em destaque
- **Funcionalidades Extras:**
  - Exportação de dados (CSV, PDF, Excel)
  - Importação de dados (CSV)
  - Configurações avançadas
  - Perfil do usuário
  - Temas customizáveis

**Entregável:** Produto completo e escalável com todas as funcionalidades, interface rica e experiência de usuário polida

## 10. Performance e Otimizações Go

### 10.1. Concorrência

**Goroutines para Processamento Assíncrono:**
```go
func (s *TransactionProcessingService) ProcessAsync(transaction *entities.Transaction) {
    go func() {
        // Processar em background
        s.process(transaction)
        
        // Publicar eventos
        for _, event := range transaction.GetEvents() {
            s.eventBus.Publish(event)
        }
    }()
}
```

### 10.2. Connection Pooling

```go
func NewDatabase() (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        PrepareStmt: true, // Prepared statements
    })
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    
    return db, err
}
```

### 10.3. Cache com Redis

```go
type CachedAccountRepository struct {
    repo  repositories.AccountRepository
    cache *redis.Client
}

func (r *CachedAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
    // Tentar cache primeiro
    cacheKey := fmt.Sprintf("account:%s", id.Value())
    cached, err := r.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        // Deserializar do cache
        return deserializeAccount(cached)
    }
    
    // Buscar do repositório
    account, err := r.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Salvar no cache
    serialized := serializeAccount(account)
    r.cache.Set(ctx, cacheKey, serialized, time.Hour)
    
    return account, nil
}
```

## 11. Deploy e DevOps

### 11.1. Dockerfile Multi-stage

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### 11.2. docker-compose.yml

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/gestao_financeira
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
    ports:
      - "6379:6379"
```

## 14. Observabilidade e Monitoramento

### 14.1. Logging Estruturado

```go
// pkg/logger/logger.go
package logger

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "os"
    "time"
)

func InitLogger(level string) {
    zerolog.TimeFieldFormat = time.RFC3339Nano
    zerolog.SetGlobalLevel(parseLevel(level))
    
    log.Logger = zerolog.New(os.Stdout).
        With().
        Timestamp().
        Str("service", "gestao-financeira").
        Logger()
}

func parseLevel(level string) zerolog.Level {
    switch level {
    case "debug":
        return zerolog.DebugLevel
    case "info":
        return zerolog.InfoLevel
    case "warn":
        return zerolog.WarnLevel
    case "error":
        return zerolog.ErrorLevel
    default:
        return zerolog.InfoLevel
    }
}

// Uso
log.Info().
    Str("user_id", userID).
    Str("action", "transaction_created").
    Float64("amount", 100.0).
    Msg("Transaction created successfully")
```

### 14.2. Métricas com Prometheus

```go
// pkg/metrics/metrics.go
package metrics

import (
    "strconv"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
    
    DatabaseQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "database_query_duration_seconds",
            Help: "Database query duration in seconds",
        },
        []string{"operation", "table"},
    )
)

// Middleware de métricas
func MetricsMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Response().StatusCode())
        
        HTTPRequestsTotal.WithLabelValues(
            c.Method(),
            c.Path(),
            status,
        ).Inc()
        
        HTTPRequestDuration.WithLabelValues(
            c.Method(),
            c.Path(),
        ).Observe(duration)
        
        return err
    }
}
```

### 14.3. Tracing com OpenTelemetry

```go
// pkg/tracing/tracing.go
package tracing

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracing(serviceName string) (*trace.TracerProvider, error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    otel.SetTracerProvider(tp)
    return tp, nil
}
```

### 14.4. Health Checks Robustos

```go
// pkg/health/health.go
package health

import (
    "context"
    "database/sql"
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{
        db:    db,
        redis: redis,
    }
}

func (h *HealthChecker) LivenessCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "alive",
    })
}

func (h *HealthChecker) ReadinessCheck(c *fiber.Ctx) error {
    ctx := context.Background()
    checks := make(map[string]string)
    
    // Verificar PostgreSQL
    if err := h.db.PingContext(ctx); err != nil {
        checks["database"] = "unhealthy"
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "status": "not ready",
            "checks": checks,
        })
    }
    checks["database"] = "healthy"
    
    // Verificar Redis
    if err := h.redis.Ping(ctx).Err(); err != nil {
        checks["cache"] = "unhealthy"
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "status": "not ready",
            "checks": checks,
        })
    }
    checks["cache"] = "healthy"
    
    return c.JSON(fiber.Map{
        "status": "ready",
        "checks": checks,
    })
}
```

## 15. Segurança

### 15.1. Headers de Segurança

```go
// pkg/middleware/security.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/limiter"
)

func SecurityMiddleware() fiber.Handler {
    return helmet.New(helmet.Config{
        XSSProtection:             "1; mode=block",
        ContentTypeNosniff:        "nosniff",
        XFrameOptions:             "DENY",
        ReferrerPolicy:            "no-referrer",
        CrossOriginEmbedderPolicy: "require-corp",
        CrossOriginOpenerPolicy:   "same-origin",
        CrossOriginResourcePolicy: "same-origin",
        OriginAgentCluster:        "?1",
        XDNSPrefetchControl:       "off",
        XDownloadOptions:          "noopen",
        XPermittedCrossDomain:     "none",
    })
}

func RateLimitMiddleware() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,              // Máximo de requisições
        Expiration: 1 * time.Minute,  // Por minuto
        KeyGenerator: func(c *fiber.Ctx) string {
            // Usar IP ou user ID
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": "Rate limit exceeded",
            })
        },
    })
}
```

### 15.2. Validação de Entrada Robusta

```go
// pkg/validator/validator.go
package validator

import (
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
    
    // Registrar validações customizadas
    validate.RegisterValidation("money", validateMoney)
    validate.RegisterValidation("date_not_future", validateDateNotFuture)
}

func validateMoney(fl validator.FieldLevel) bool {
    amount := fl.Field().Float()
    return amount > 0 && amount <= 999999999.99
}

func validateDateNotFuture(fl validator.FieldLevel) bool {
    date := fl.Field().Interface().(time.Time)
    return date.Before(time.Now()) || date.Equal(time.Now())
}

func ValidateAndReturnErrors(s interface{}) map[string]string {
    err := validate.Struct(s)
    if err == nil {
        return nil
    }
    
    errors := make(map[string]string)
    for _, err := range err.(validator.ValidationErrors) {
        field := err.Field()
        tag := err.Tag()
        errors[field] = getErrorMessage(field, tag)
    }
    
    return errors
}
```

### 15.3. Proteção contra SQL Injection e XSS

```go
// Sempre usar prepared statements (GORM faz isso automaticamente)
// Mas para queries raw:
func (r *Repository) FindByQuery(query string) ([]Model, error) {
    // NUNCA fazer: fmt.Sprintf("SELECT * FROM table WHERE name = '%s'", query)
    // SEMPRE usar:
    var results []Model
    err := r.db.Where("name = ?", query).Find(&results).Error
    return results, err
}

// Sanitização de entrada
func sanitizeInput(input string) string {
    // Remover caracteres perigosos
    input = strings.TrimSpace(input)
    input = html.EscapeString(input)
    return input
}
```

## 16. Performance e Escalabilidade

### 16.1. Paginação Eficiente

```go
// pkg/pagination/pagination.go
package pagination

type Pagination struct {
    Page  int `json:"page" validate:"min=1"`
    Limit int `json:"limit" validate:"min=1,max=100"`
}

func (p *Pagination) Offset() int {
    return (p.Page - 1) * p.Limit
}

func (p *Pagination) ToSQL() (offset, limit int) {
    return p.Offset(), p.Limit
}

// Uso no handler
func (h *TransactionHandler) List(c *fiber.Ctx) error {
    var pagination Pagination
    if err := c.QueryParser(&pagination); err != nil {
        pagination = Pagination{Page: 1, Limit: 20}
    }
    
    transactions, total, err := h.useCase.List(
        c.Locals("userID").(string),
        pagination,
    )
    
    return c.JSON(fiber.Map{
        "data": transactions,
        "pagination": fiber.Map{
            "page":       pagination.Page,
            "limit":      pagination.Limit,
            "total":      total,
            "totalPages": (total + pagination.Limit - 1) / pagination.Limit,
        },
    })
}
```

### 16.2. Índices de Banco de Dados

```sql
-- Migrations com índices otimizados
CREATE INDEX idx_transactions_user_date ON transactions(user_id, date DESC);
CREATE INDEX idx_transactions_account_date ON transactions(account_id, date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_type_date ON transactions(user_id, type, date DESC);
CREATE INDEX idx_accounts_user_context ON accounts(user_id, context);

-- Índices compostos para queries comuns
CREATE INDEX idx_transactions_reporting ON transactions(user_id, date, type, status);
```

### 16.3. Cache Estratégico

```go
// pkg/cache/cache.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type Cache struct {
    client *redis.Client
    ctx    context.Context
}

func (c *Cache) GetOrSet(key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
    // Tentar buscar do cache
    val, err := c.client.Get(c.ctx, key).Result()
    if err == nil {
        var result interface{}
        json.Unmarshal([]byte(val), &result)
        return result, nil
    }
    
    // Se não encontrou, executar função e cachear
    result, err := fn()
    if err != nil {
        return nil, err
    }
    
    data, _ := json.Marshal(result)
    c.client.Set(c.ctx, key, data, ttl)
    
    return result, nil
}

// Cache de relatórios (TTL maior)
func (s *ReportService) GetMonthlyReport(userID string, month time.Month) (*Report, error) {
    key := fmt.Sprintf("report:%s:%d", userID, month)
    
    return s.cache.GetOrSet(key, 1*time.Hour, func() (interface{}, error) {
        return s.generateReport(userID, month)
    })
}
```

### 16.4. Processamento Assíncrono com Workers

```go
// pkg/worker/worker.go
package worker

import (
    "context"
    "github.com/hibiken/asynq"
)

type Worker struct {
    client *asynq.Client
    server *asynq.Server
}

func NewWorker(redisAddr string) *Worker {
    redisOpt := asynq.RedisClientOpt{Addr: redisAddr}
    
    return &Worker{
        client: asynq.NewClient(redisOpt),
        server: asynq.NewServer(redisOpt, asynq.Config{
            Concurrency: 10, // Número de workers
        }),
    }
}

// Enfileirar tarefa
func (w *Worker) EnqueueProcessRecurringTransactions() error {
    task := asynq.NewTask("process:recurring", nil)
    _, err := w.client.Enqueue(task, asynq.Queue("transactions"))
    return err
}

// Processar tarefa
func (w *Worker) Start() error {
    mux := asynq.NewServeMux()
    mux.HandleFunc("process:recurring", w.handleRecurringTransactions)
    return w.server.Start(mux)
}
```

## 17. Documentação da API

### 17.1. Swagger/OpenAPI

```go
// cmd/api/main.go
import (
    "github.com/gofiber/swagger"
    _ "gestao-financeira/docs" // Gerado pelo swag
)

// @title Gestão Financeira API
// @version 1.0
// @description API para gestão financeira pessoal e profissional
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
    // ...
    app.Get("/swagger/*", swagger.HandlerDefault)
}

// Handler com anotações Swagger
// @Summary Criar transação
// @Description Cria uma nova transação financeira
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body CreateTransactionRequest true "Dados da transação"
// @Success 201 {object} TransactionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /transactions [post]
// @Security Bearer
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
    // ...
}
```

## 18. CI/CD e Deploy

### 18.1. GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
  
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build
        run: |
          CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
      
      - name: Build Docker image
        run: docker build -t gestao-financeira:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin
          docker push gestao-financeira:${{ github.sha }}
```

### 18.2. Deploy em Produção

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  api:
    image: gestao-financeira:latest
    restart: always
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - ENV=production
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  db:
    image: postgres:15
    restart: always
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    deploy:
      resources:
        limits:
          memory: 2G
  
  redis:
    image: redis:7-alpine
    restart: always
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
```

## 19. Backup e Disaster Recovery

### 19.1. Estratégia de Backup

```go
// pkg/backup/backup.go
package backup

import (
    "fmt"
    "os/exec"
    "time"
    "github.com/rs/zerolog/log"
)

type DatabaseConfig struct {
    Host     string
    User     string
    Database string
}

func BackupDatabase(dbConfig DatabaseConfig) error {
    timestamp := time.Now().Format("20060102_150405")
    filename := fmt.Sprintf("backup_%s.sql", timestamp)
    
    cmd := exec.Command("pg_dump",
        "-h", dbConfig.Host,
        "-U", dbConfig.User,
        "-d", dbConfig.Database,
        "-f", filename,
    )
    
    return cmd.Run()
}

// Agendar backups diários
func ScheduleBackups(dbConfig DatabaseConfig) {
    ticker := time.NewTicker(24 * time.Hour)
    go func() {
        for range ticker.C {
            if err := BackupDatabase(dbConfig); err != nil {
                log.Error().Err(err).Msg("Backup failed")
            } else {
                log.Info().Str("file", fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))).Msg("Backup completed")
            }
        }
    }()
}
```

## 20. Testes E2E

### 20.1. Testes End-to-End

```go
// tests/e2e/transaction_test.go
package e2e_test

import (
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

func TestCreateTransactionE2E(t *testing.T) {
    app := setupTestApp()
    
    // 1. Registrar usuário
    registerResp := registerUser(t, app, "test@example.com", "password123")
    token := registerResp["token"].(string)
    
    // 2. Criar conta
    accountResp := createAccount(t, app, token)
    accountID := accountResp["id"].(string)
    
    // 3. Criar transação
    transactionResp := createTransaction(t, app, token, accountID, 100.0)
    
    assert.Equal(t, 201, transactionResp.StatusCode)
    assert.NotEmpty(t, transactionResp["id"])
}

func setupTestApp() *fiber.App {
    // Setup app de teste com banco de teste
    // ...
}
```

## 21. Auditoria e Compliance

### 21.1. Log de Auditoria

```go
// pkg/audit/audit.go
package audit

import (
    "time"
    "github.com/google/uuid"
)

type AuditLog struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`    // CREATE, UPDATE, DELETE, VIEW
    Resource  string    `json:"resource"`  // transaction, account, etc.
    ResourceID string   `json:"resource_id"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
    Timestamp time.Time `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata"`
}

type AuditLogger struct {
    repo AuditRepository
}

func (a *AuditLogger) Log(action, resource, resourceID, userID, ip, userAgent string, metadata map[string]interface{}) {
    log := AuditLog{
        ID:         uuid.New().String(),
        UserID:     userID,
        Action:     action,
        Resource:   resource,
        ResourceID: resourceID,
        IPAddress:  ip,
        UserAgent:  userAgent,
        Timestamp:  time.Now(),
        Metadata:   metadata,
    }
    
    // Salvar de forma assíncrona
    go a.repo.Save(log)
}

// Middleware de auditoria
func AuditMiddleware(auditLogger *AuditLogger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("userID")
        if userID != nil {
            auditLogger.Log(
                c.Method(),
                extractResource(c.Path()),
                c.Params("id"),
                userID.(string),
                c.IP(),
                c.Get("User-Agent"),
                nil,
            )
        }
        return c.Next()
    }
}
```

### 21.2. LGPD/GDPR Compliance

```go
// pkg/compliance/gdpr.go
package compliance

// Direito ao esquecimento (LGPD/GDPR)
func (s *UserService) DeleteUserData(userID string) error {
    // 1. Anonimizar dados pessoais
    // 2. Manter dados financeiros agregados (se necessário para compliance fiscal)
    // 3. Registrar ação de exclusão
    // 4. Notificar usuário
    return nil
}

// Exportação de dados (LGPD/GDPR)
func (s *UserService) ExportUserData(userID string) ([]byte, error) {
    // Exportar todos os dados do usuário em formato JSON
    data := struct {
        User         UserData         `json:"user"`
        Transactions []TransactionData `json:"transactions"`
        Accounts     []AccountData     `json:"accounts"`
        // ...
    }{}
    
    return json.Marshal(data)
}
```

## 22. Escalabilidade e Multi-tenancy (Preparação para Produto)

### 22.1. Estratégia de Escalabilidade

**Horizontal Scaling:**
- Múltiplas instâncias da API (load balancer)
- Database read replicas para relatórios
- Redis cluster para cache distribuído
- Message queue para processamento assíncrono

**Vertical Scaling:**
- Otimização de queries
- Índices estratégicos
- Connection pooling otimizado
- Cache agressivo

### 22.2. Multi-tenancy (Opcional - se virar SaaS)

```go
// Estrutura para suportar múltiplos tenants
type Tenant struct {
    ID       string
    Name     string
    Plan     string // FREE, PREMIUM, ENTERPRISE
    Settings TenantSettings
}

// Isolamento por tenant_id em todas as queries
func (r *Repository) FindByTenant(tenantID string, filters Filters) ([]Entity, error) {
    return r.db.Where("tenant_id = ?", tenantID).Where(filters).Find(&entities).Error
}
```

## 23. Tratamento de Erros Robusto

### 23.1. Erros de Domínio

```go
// pkg/errors/domain_error.go
package errors

type DomainError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

func (e *DomainError) Error() string {
    return e.Message
}

// Erros específicos do domínio
var (
    ErrInsufficientBalance = &DomainError{
        Code:    "INSUFFICIENT_BALANCE",
        Message: "Account balance is insufficient",
    }
    
    ErrTransactionNotFound = &DomainError{
        Code:    "TRANSACTION_NOT_FOUND",
        Message: "Transaction not found",
    }
    
    ErrInvalidAmount = &DomainError{
        Code:    "INVALID_AMOUNT",
        Message: "Transaction amount must be positive",
    }
)
```

### 23.2. Middleware de Request ID

```go
// pkg/middleware/request_id.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

func RequestIDMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        c.Locals("requestID", requestID)
        c.Set("X-Request-ID", requestID)
        
        return c.Next()
    }
}
```

## 24. Testes de Performance e Carga

### 24.1. Benchmarks

```go
// tests/benchmark/transaction_bench_test.go
package benchmark_test

import (
    "testing"
    "github.com/gofiber/fiber/v2"
)

func BenchmarkCreateTransaction(b *testing.B) {
    app := setupTestApp()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Simular criação de transação
        createTransactionRequest(app)
    }
}

func BenchmarkListTransactions(b *testing.B) {
    app := setupTestApp()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Simular listagem
        listTransactionsRequest(app)
    }
}
```

### 24.2. Testes de Carga (k6 ou Vegeta)

```javascript
// tests/load/transactions.js (k6)
import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp up
        { duration: '5m', target: 100 },  // Stay at 100
        { duration: '2m', target: 200 }, // Ramp up to 200
        { duration: '5m', target: 200 }, // Stay at 200
        { duration: '2m', target: 0 },  // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% das requisições < 500ms
        http_req_failed: ['rate<0.01'],  // Taxa de erro < 1%
    },
};

export default function () {
    const res = http.get('http://localhost:8080/api/v1/transactions', {
        headers: { 'Authorization': 'Bearer ' + __ENV.TOKEN },
    });
    
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });
}
```

## 25. Versionamento de API

### 25.1. Estrutura de Versionamento

```go
// API v1 e v2 coexistem
func setupRoutes(app *fiber.App) {
    // API v1 (legacy)
    v1 := app.Group("/api/v1")
    setupV1Routes(v1)
    
    // API v2 (nova)
    v2 := app.Group("/api/v2")
    setupV2Routes(v2)
}

// Deprecation header
func DeprecationMiddleware(version string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("Deprecation", "true")
        c.Set("Sunset", "2025-12-31")
        c.Set("Link", "</api/v2>; rel=\"successor-version\"")
        return c.Next()
    }
}
```

## 14. Considerações Finais

### Vantagens do Go para este Projeto:
1. ✅ **Performance excepcional** - Compilado, muito rápido
2. ✅ **Concorrência nativa** - Goroutines para processamento assíncrono
3. ✅ **Simplicidade** - Código limpo e direto
4. ✅ **Type safety** - Forte e estático
5. ✅ **Deploy simples** - Binário único
6. ✅ **Aprendizado valioso** - Habilidade moderna e valorizada

### Desafios:
1. ⚠️ **Curva de aprendizado inicial** - Mas Go é simples
2. ⚠️ **Menos exemplos DDD** - Mas estrutura é clara
3. ⚠️ **Error handling verboso** - Mas explícito e seguro

### Recomendações para Produção:
- ✅ **Fiber + GORM** - API familiar e performance
- ✅ **Goroutines** - Processamento assíncrono
- ✅ **Redis Cache** - Cache estratégico para relatórios
- ✅ **Testify** - Testes expressivos
- ✅ **Observabilidade** - Logs estruturados, métricas, tracing
- ✅ **Segurança** - Rate limiting, CORS configurado, headers de segurança
- ✅ **Health Checks** - Liveness e readiness robustos
- ✅ **Graceful Shutdown** - Encerramento controlado
- ✅ **Documentação** - Swagger/OpenAPI
- ✅ **CI/CD** - Pipeline automatizado
- ✅ **Backup** - Estratégia de backup automático
- ✅ **Monitoramento** - Prometheus + Grafana
- ✅ **Escalabilidade** - Preparado para horizontal scaling


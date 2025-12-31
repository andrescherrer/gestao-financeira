# Planejamento DDD - Sistema de Gestão Financeira (Node.js)

## Resumo Executivo

Sistema de gestão financeira pessoal e profissional desenvolvido em **Node.js com TypeScript** usando **NestJS** e **Prisma**, seguindo **Domain-Driven Design (DDD)**. Focado em **type-safety**, **escalabilidade** e **pronto para produção**, com potencial para evoluir para produto SaaS.

**Stack Principal:**
- **Backend**: Node.js 20+ LTS com NestJS 10+ (framework web)
- **Frontend**: Vue 3 com TypeScript (já existente - reutilizar) ⚠️ Ver seção 3.3
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis (cache e rate limiting)
- **ORM**: Prisma 5+ (type-safe)
- **Observabilidade**: OpenTelemetry, Prometheus + Grafana

> ⚠️ **NOTA:** O projeto já possui um frontend Vue 3 funcional. Não é necessário criar novo frontend. Veja seção 3.3 para detalhes de compatibilidade.

**Diferenciais:**
- Type-safety excepcional (TypeScript + Prisma)
- Arquitetura DDD nativa do NestJS
- Observabilidade completa
- Segurança robusta
- Pronto para produção

## 1. Visão Geral

Sistema de gestão financeira pessoal e profissional desenvolvido em **Node.js com TypeScript** (backend) seguindo Domain-Driven Design (DDD), aproveitando o type-safety, DDD nativo do NestJS e ecossistema moderno. 

> ⚠️ **IMPORTANTE:** O projeto já possui um **frontend Vue 3** completamente funcional e independente. Este frontend deve ser **reutilizado** sem modificações, apenas configurando a URL da API. Veja seção 3.3 para detalhes de compatibilidade.

Projetado para **alta performance I/O** e **escalabilidade**, com potencial para evoluir para produto comercial.

## 2. Objetivos

- Controle total de finanças pessoais e profissionais
- Separação clara entre contas pessoais e profissionais
- Análise e relatórios financeiros
- Planejamento orçamentário
- Acompanhamento de metas financeiras
- Arquitetura escalável e manutenível
- Aproveitamento máximo do TypeScript e DDD nativo do NestJS

## 3. Stack Tecnológico Node.js

### 3.1. Tecnologias Principais

- **Runtime**: Node.js 20+ LTS
- **Framework**: NestJS 10+
- **Linguagem**: TypeScript 5+
- **ORM**: Prisma 5+
- **Validação**: class-validator + class-transformer
- **Autenticação**: @nestjs/jwt + @nestjs/passport
- **Event Bus**: @nestjs/event-emitter ou RabbitMQ
- **Testes**: Jest + @nestjs/testing
- **Migrations**: Prisma Migrate
- **Config**: @nestjs/config
- **Logging**: nestjs-pino ou winston
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis (ioredis)
- **API Docs**: Swagger/OpenAPI (@nestjs/swagger)

### 3.2. Framework Web: NestJS

#### 3.2.1. Sobre o NestJS

**NestJS** 🚀 - Progressive Node.js framework
- **GitHub**: 70k+ stars
- **Tipo**: Framework progressivo para Node.js
- **Abordagem**: Inspirado em Angular (decorators, DI, modules)
- **Base**: Express.js ou Fastify (HTTP engine)

#### 3.2.2. Características do NestJS

**Vantagens:**
- ✅ **TypeScript nativo** - Type-safety completo
- ✅ **DDD nativo** - Arquitetura modular perfeita para DDD
- ✅ **Dependency Injection** - DI nativa (similar ao Angular)
- ✅ **Decorators** - Código expressivo e elegante
- ✅ **Modules** - Organização modular clara
- ✅ **Guards** - Autenticação e autorização elegantes
- ✅ **Interceptors** - Middleware poderoso
- ✅ **Pipes** - Validação e transformação integradas
- ✅ **Exception Filters** - Tratamento de erros centralizado
- ✅ **Swagger automático** - Documentação API automática
- ✅ **Testes** - Suporte nativo excelente

**Considerações:**
- ⚠️ **Curva de aprendizado** - Se não conhece Angular, pode ser inicial
- ⚠️ **Overhead** - Framework completo tem mais overhead que Express puro
- ⚠️ **Comunidade** - Menor que Express, mas crescente e ativa

#### 3.2.3. Performance do NestJS

**Benchmarks (aproximado):**

```
API Simples (Hello World):
- NestJS (Express):    ~50.000 req/s
- NestJS (Fastify):    ~80.000 req/s
- Express puro:        ~50.000 req/s

API com JSON (CRUD):
- NestJS (Express):    ~45.000 req/s
- NestJS (Fastify):    ~75.000 req/s
```

**Por que NestJS é uma excelente escolha?**
- Type-safety completo com TypeScript
- Arquitetura modular perfeita para DDD
- Dependency Injection nativa
- Decorators expressivos
- Swagger automático
- Testes integrados

#### 3.2.4. Exemplos de Código com NestJS

##### Exemplo Básico

```typescript
// src/main.ts
import { NestFactory } from '@nestjs/core';
import { ValidationPipe } from '@nestjs/common';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Global prefix
  app.setGlobalPrefix('api/v1');
  
  // CORS
  app.enableCors({
    origin: process.env.ALLOWED_ORIGINS?.split(',') || '*',
    credentials: true,
  });
  
  // Validation pipe global
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      forbidNonWhitelisted: true,
      transform: true,
    }),
  );
  
  // Swagger
  const config = new DocumentBuilder()
    .setTitle('Gestão Financeira API')
    .setDescription('API para gestão financeira pessoal e profissional')
    .setVersion('1.0')
    .addBearerAuth()
    .build();
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api/docs', app, document);
  
  // Health checks
  app.getHttpAdapter().get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });
  
  await app.listen(3000);
}
bootstrap();
```

##### Controller com NestJS

```typescript
// src/identity/presentation/controllers/auth.controller.ts
import { Controller, Post, Body, HttpCode, HttpStatus } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { RegisterUserUseCase } from '../../application/use-cases/register-user.use-case';
import { RegisterUserDto } from '../dtos/register-user.dto';

@ApiTags('Authentication')
@Controller('auth')
export class AuthController {
  constructor(
    private readonly registerUserUseCase: RegisterUserUseCase,
  ) {}
  
  @Post('register')
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: 'Register a new user' })
  @ApiResponse({ status: 201, description: 'User successfully registered' })
  @ApiResponse({ status: 400, description: 'Bad request' })
  async register(@Body() dto: RegisterUserDto) {
    return this.registerUserUseCase.execute({
      email: dto.email,
      password: dto.password,
      firstName: dto.firstName,
      lastName: dto.lastName,
    });
  }
}
```

#### 3.2.5. Middleware com NestJS

##### Auth Guard

```typescript
// src/shared/infrastructure/guards/jwt-auth.guard.ts
import { Injectable, CanActivate, ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';

@Injectable()
export class JwtAuthGuard implements CanActivate {
  constructor(private jwtService: JwtService) {}
  
  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest();
    const token = this.extractTokenFromHeader(request);
    
    if (!token) {
      throw new UnauthorizedException();
    }
    
    try {
      const payload = this.jwtService.verify(token);
      request.user = payload;
      return true;
    } catch {
      throw new UnauthorizedException();
    }
  }
  
  private extractTokenFromHeader(request: any): string | undefined {
    const [type, token] = request.headers.authorization?.split(' ') ?? [];
    return type === 'Bearer' ? token : undefined;
  }
}
```

##### Logger Interceptor

```typescript
// src/shared/infrastructure/interceptors/logging.interceptor.ts
import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from '@nestjs/common';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Logger } from '@nestjs/common';

@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  private readonly logger = new Logger(LoggingInterceptor.name);
  
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const request = context.switchToHttp().getRequest();
    const { method, url } = request;
    const now = Date.now();
    
    return next.handle().pipe(
      tap(() => {
        const response = context.switchToHttp().getResponse();
        const { statusCode } = response;
        const delay = Date.now() - now;
        
        this.logger.log(`${method} ${url} ${statusCode} - ${delay}ms`);
      }),
    );
  }
}
```

#### 3.2.6. Estrutura Completa com NestJS

```typescript
// src/app.module.ts
import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { APP_INTERCEPTOR, APP_GUARD } from '@nestjs/core';
import { EventEmitterModule } from '@nestjs/event-emitter';
import { CqrsModule } from '@nestjs/cqrs';
import { ThrottlerModule, ThrottlerGuard } from '@nestjs/throttler';
import { LoggingInterceptor } from './shared/infrastructure/interceptors/logging.interceptor';
import { IdentityModule } from './identity/identity.module';
import { TransactionModule } from './transaction/transaction.module';
// ... outros módulos

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    EventEmitterModule.forRoot(),
    CqrsModule,
    ThrottlerModule.forRoot({
      ttl: 60,
      limit: 100,
    }),
    IdentityModule,
    TransactionModule,
    // ... outros módulos
  ],
  providers: [
    {
      provide: APP_INTERCEPTOR,
      useClass: LoggingInterceptor,
    },
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
    },
  ],
})
export class AppModule {}
```

#### 3.2.7. Validação com NestJS

```typescript
// src/shared/infrastructure/pipes/validation.pipe.ts
import { PipeTransform, Injectable, ArgumentMetadata, BadRequestException } from '@nestjs/common';
import { validate } from 'class-validator';
import { plainToInstance } from 'class-transformer';

@Injectable()
export class ValidationPipe implements PipeTransform<any> {
  async transform(value: any, { metatype }: ArgumentMetadata) {
    if (!metatype || !this.toValidate(metatype)) {
      return value;
    }
    
    const object = plainToInstance(metatype, value);
    const errors = await validate(object);
    
    if (errors.length > 0) {
      throw new BadRequestException('Validation failed');
    }
    
    return value;
  }
  
  private toValidate(metatype: Function): boolean {
    const types: Function[] = [String, Boolean, Number, Array, Object];
    return !types.includes(metatype);
  }
}
```

#### 3.2.8. Vantagens do NestJS para este Projeto

**Por que NestJS é uma excelente escolha:**
1. ✅ **Type-Safety Completo** - TypeScript nativo em tudo
2. ✅ **DDD Nativo** - Arquitetura modular perfeita
3. ✅ **DI Nativa** - Dependency Injection nativa
4. ✅ **Decorators** - Código expressivo e elegante
5. ✅ **Swagger Automático** - Documentação API automática
6. ✅ **Testes Integrados** - Suporte nativo excelente
7. ✅ **Ecossistema** - Muitas ferramentas disponíveis

**Considerações:**
- ⚠️ **Curva de aprendizado** - Se não conhece Angular, pode ser inicial
- ⚠️ **Overhead** - Framework completo tem mais overhead que Express puro

#### 3.2.9. Recursos Adicionais do NestJS

**Módulos Disponíveis:**
- `@nestjs/config` - Configuração
- `@nestjs/jwt` - JWT authentication
- `@nestjs/passport` - Passport integration
- `@nestjs/swagger` - Swagger/OpenAPI
- `@nestjs/throttler` - Rate limiting
- `@nestjs/terminus` - Health checks
- `@nestjs/prometheus` - Prometheus metrics
- `@nestjs/event-emitter` - Event bus
- `@nestjs/cqrs` - CQRS pattern
- `@nestjs/schedule` - Task scheduling

**Exemplo com Módulos Adicionais:**

```typescript
import { ThrottlerModule } from '@nestjs/throttler';
import { TerminusModule } from '@nestjs/terminus';
import { PrometheusModule } from '@nestjs/prometheus';

@Module({
  imports: [
    ThrottlerModule.forRoot({
      ttl: 60,
      limit: 100,
    }),
    TerminusModule,
    PrometheusModule.register(),
  ],
})
export class AppModule {}
```

### 3.3. Por que Node.js + NestJS?

**Vantagens:**
- ✅ **TypeScript nativo** - Type-safety excelente
- ✅ **DDD nativo** - NestJS foi feito pensando em DDD
- ✅ **Ecossistema moderno** - Prisma, TypeORM, etc.
- ✅ **Performance I/O** - Excelente para APIs
- ✅ **Frontend** - Mesma linguagem (TypeScript)
- ✅ **Real-time** - WebSockets nativos
- ✅ **Microservices** - Fácil escalar horizontalmente
- ✅ **Async/Await** - Código limpo e moderno
- ✅ **Decorators** - Código expressivo e elegante

**Desafios:**
- ⚠️ **Single-threaded** - CPU-bound pode ser limitante
- ⚠️ **Runtime overhead** - JavaScript tem overhead
- ⚠️ **Memory** - Pode consumir mais que Go

### 3.3. Compatibilidade com Frontend Vue 3 Existente

**✅ IMPORTANTE:** O projeto já possui um frontend Vue 3 completamente funcional e independente do backend. **NÃO é necessário criar um novo frontend** para Node.js.

#### 3.3.1. Por que o Frontend Vue 3 é Reutilizável?

O frontend Vue 3 atual foi desenvolvido de forma **desacoplada** do backend, comunicando-se exclusivamente via **API REST**. Isso significa:

- ✅ **Arquitetura independente**: Frontend não depende da tecnologia do backend
- ✅ **Comunicação via HTTP/JSON**: Qualquer backend que implemente a mesma API funciona
- ✅ **Configuração via variáveis de ambiente**: Apenas muda a URL da API
- ✅ **Zero alterações no código**: O frontend Vue 3 funciona sem modificações

#### 3.3.2. Requisitos de Compatibilidade da API

Para que o frontend Vue 3 existente funcione com o backend Node.js/NestJS, é necessário implementar a **mesma interface de API REST**:

**Endpoints Principais:**
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/auth/me
POST   /api/v1/transactions
GET    /api/v1/transactions
GET    /api/v1/transactions/:id
PUT    /api/v1/transactions/:id
DELETE /api/v1/transactions/:id
GET    /api/v1/accounts
POST   /api/v1/accounts
GET    /api/v1/categories
POST   /api/v1/categories
GET    /api/v1/budgets
POST   /api/v1/budgets
GET    /api/v1/reports/monthly
... (outros endpoints conforme necessário)
```

**Formato de Request/Response:**
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

**Autenticação:**
- JWT tokens com mesma estrutura de payload
- Refresh tokens (se implementado)
- Mesmos headers de autenticação

#### 3.3.3. Configuração do Frontend para Node.js

Para usar o frontend Vue 3 com o backend Node.js/NestJS, apenas configure a variável de ambiente:

```bash
# .env ou .env.local no frontend
VITE_API_URL=http://localhost:3000/api/v1
```

Ou no `docker-compose.yml`:
```yaml
frontend:
  environment:
    - VITE_API_URL=http://node-backend:3000/api/v1
```

**Nenhuma alteração no código do frontend é necessária!**

#### 3.3.4. Vantagens dessa Abordagem

1. **Migração facilitada**: Trocar de backend Go para Node.js é apenas implementar a API
2. **Reutilização total**: Frontend Vue 3 já desenvolvido e testado
3. **Economia de tempo**: Não precisa desenvolver novo frontend
4. **Consistência**: Mesma experiência de usuário independente do backend
5. **Testes**: Frontend já testado e validado com a API

#### 3.3.5. Checklist de Compatibilidade

Ao implementar o backend Node.js/NestJS, garanta:

- [ ] Todos os endpoints da API implementados
- [ ] Mesmo formato de request/response JSON
- [ ] Mesma estrutura de autenticação JWT
- [ ] Mesmos códigos de status HTTP
- [ ] Mesmas validações e mensagens de erro
- [ ] CORS configurado corretamente
- [ ] Headers de segurança compatíveis
- [ ] Paginação implementada (se aplicável)
- [ ] Filtros e ordenação (se aplicável)

#### 3.3.6. Documentação da API

O projeto Go atual possui documentação Swagger/OpenAPI em `/docs/swagger.json`. Use essa documentação como referência para garantir compatibilidade:

- Endpoints exatos
- Parâmetros de request
- Estrutura de response
- Códigos de erro
- Validações

**Referência:** `backend/docs/swagger.json` ou `http://localhost:8080/docs/swagger/index.html`

#### 3.3.7. Nota sobre a Seção "Stack Tecnológico Frontend"

A seção "4. Stack Tecnológico Frontend" neste documento menciona Next.js como opção de frontend. Isso é apenas uma **referência teórica** para novos projetos. Para este projeto específico, **reutilize o frontend Vue 3 existente** conforme descrito acima.

## 4. Arquitetura DDD em NestJS

### 4.1. Estrutura em Camadas (NestJS)

```
┌─────────────────────────────────────┐
│     Controllers (Presentation)       │  (@Controller, DTOs)
├─────────────────────────────────────┤
│     Use Cases (Application)         │  (@Injectable, Services)
├─────────────────────────────────────┤
│     Domain Layer                    │  (Entities, Value Objects, Domain Services)
├─────────────────────────────────────┤
│     Repositories (Infrastructure)   │  (Prisma, External Services)
└─────────────────────────────────────┘
```

### 4.2. Bounded Contexts

1. **Identity Context** - Autenticação e gestão de usuários
2. **Account Management Context** - Gestão de contas e carteiras
3. **Transaction Context** - Processamento de transações financeiras
4. **Category Context** - Gestão de categorias e taxonomia
5. **Budget Context** - Planejamento e controle orçamentário
6. **Reporting Context** - Análises e relatórios financeiros
7. **Investment Context** - Gestão de investimentos
8. **Goal Context** - Metas e objetivos financeiros
9. **Notification Context** - Notificações e alertas

## 4. Stack Tecnológico Frontend

> ⚠️ **NOTA IMPORTANTE:** Esta seção descreve uma abordagem teórica para novos projetos. Para este projeto específico, **reutilize o frontend Vue 3 existente** conforme descrito na seção 3.3 (Compatibilidade com Frontend Vue 3 Existente). O frontend Vue 3 atual é completamente compatível e não requer desenvolvimento adicional.

### 4.1. Abordagem de Desenvolvimento (Referência Teórica)

> Esta seção é apenas para referência. O frontend Vue 3 existente deve ser reutilizado.

O frontend seria desenvolvido **incrementalmente**, sempre que um módulo do backend estiver pronto. Esta abordagem permite:

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
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000/api/v1',
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

## 5. Arquitetura DDD em NestJS

```
gestao-financeira-node/
├── src/
│   ├── shared/                          # Shared Kernel
│   │   ├── domain/
│   │   │   ├── value-objects/
│   │   │   │   ├── money.ts
│   │   │   │   ├── currency.ts
│   │   │   │   └── account-context.ts
│   │   │   └── events/
│   │   │       └── domain-event.ts
│   │   └── infrastructure/
│   │       └── event-bus/
│   │           └── event-bus.service.ts
│   │
│   ├── identity/                         # Identity Context
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   └── user.entity.ts
│   │   │   ├── value-objects/
│   │   │   │   ├── email.vo.ts
│   │   │   │   ├── password-hash.vo.ts
│   │   │   │   └── user-name.vo.ts
│   │   │   ├── services/
│   │   │   │   ├── password.service.ts
│   │   │   │   └── token.service.ts
│   │   │   ├── repositories/
│   │   │   │   └── user.repository.interface.ts
│   │   │   └── events/
│   │   │       ├── user-registered.event.ts
│   │   │       └── user-password-changed.event.ts
│   │   ├── application/
│   │   │   ├── use-cases/
│   │   │   │   ├── register-user.use-case.ts
│   │   │   │   ├── login.use-case.ts
│   │   │   │   └── update-profile.use-case.ts
│   │   │   └── dtos/
│   │   │       └── register-user.dto.ts
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   └── prisma-user.repository.ts
│   │   │   ├── services/
│   │   │   │   └── jwt-token.service.ts
│   │   │   └── events/
│   │   │       └── domain-event-publisher.service.ts
│   │   └── presentation/
│   │       ├── controllers/
│   │       │   └── auth.controller.ts
│   │       └── dtos/
│   │
│   ├── account-management/              # Account Management Context
│   ├── transaction/                      # Transaction Context
│   ├── category/                        # Category Context
│   ├── budget/                           # Budget Context
│   ├── reporting/                        # Reporting Context
│   ├── investment/                      # Investment Context
│   ├── goal/                             # Goal Context
│   └── notification/                    # Notification Context
│
├── prisma/
│   ├── schema.prisma
│   └── migrations/
│
├── test/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── nest-cli.json
├── package.json
├── tsconfig.json
├── Dockerfile
└── docker-compose.yml
```

## 6. Detalhamento dos Bounded Contexts (NestJS)

### 6.1. Identity Context

#### 6.1.1. Entidades (TypeScript)

**User (Agregado Raiz)**
```typescript
// src/identity/domain/entities/user.entity.ts
import { AggregateRoot } from '@nestjs/cqrs';
import { Email } from '../value-objects/email.vo';
import { PasswordHash } from '../value-objects/password-hash.vo';
import { UserName } from '../value-objects/user-name.vo';
import { UserProfile } from '../value-objects/user-profile.vo';
import { UserRegistered } from '../events/user-registered.event';
import { UserPasswordChanged } from '../events/user-password-changed.event';

export class UserId {
  constructor(private readonly value: string) {}
  
  equals(other: UserId): boolean {
    return this.value === other.value;
  }
  
  toString(): string {
    return this.value;
  }
}

export class User extends AggregateRoot {
  private constructor(
    private readonly id: UserId,
    private email: Email,
    private passwordHash: PasswordHash,
    private name: UserName,
    private profile: UserProfile,
    private readonly createdAt: Date,
    private updatedAt: Date,
    private isActive: boolean,
  ) {
    super();
  }
  
  static create(
    email: Email,
    passwordHash: PasswordHash,
    name: UserName,
  ): User {
    const user = new User(
      new UserId(crypto.randomUUID()),
      email,
      passwordHash,
      name,
      UserProfile.default(),
      new Date(),
      new Date(),
      true,
    );
    
    user.addDomainEvent(new UserRegistered(user.id.toString()));
    return user;
  }
  
  changePassword(oldPassword: string, newPassword: string): void {
    if (!this.passwordHash.verify(oldPassword)) {
      throw new Error('Invalid old password');
    }
    
    this.passwordHash = PasswordHash.fromPlainPassword(newPassword);
    this.updatedAt = new Date();
    this.addDomainEvent(new UserPasswordChanged(this.id.toString()));
  }
  
  updateProfile(profile: UserProfile): void {
    this.profile = profile;
    this.updatedAt = new Date();
  }
  
  deactivate(): void {
    this.isActive = false;
    this.updatedAt = new Date();
  }
  
  getId(): UserId {
    return this.id;
  }
  
  getEmail(): Email {
    return this.email;
  }
  
  getPasswordHash(): PasswordHash {
    return this.passwordHash;
  }
}
```

#### 6.1.2. Value Objects (TypeScript)

**Email**
```typescript
// src/identity/domain/value-objects/email.vo.ts
export class Email {
  private constructor(private readonly value: string) {
    this.validate(value);
  }
  
  static create(email: string): Email {
    return new Email(email.toLowerCase().trim());
  }
  
  private validate(email: string): void {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      throw new Error('Invalid email format');
    }
  }
  
  equals(other: Email): boolean {
    return this.value === other.value;
  }
  
  toString(): string {
    return this.value;
  }
  
  getValue(): string {
    return this.value;
  }
}
```

**PasswordHash**
```typescript
// src/identity/domain/value-objects/password-hash.vo.ts
import * as bcrypt from 'bcrypt';

export class PasswordHash {
  private constructor(private readonly value: string) {}
  
  static fromPlainPassword(password: string): PasswordHash {
    if (password.length < 8) {
      throw new Error('Password must be at least 8 characters');
    }
    
    const hash = bcrypt.hashSync(password, 10);
    return new PasswordHash(hash);
  }
  
  static fromHash(hash: string): PasswordHash {
    return new PasswordHash(hash);
  }
  
  verify(plainPassword: string): boolean {
    return bcrypt.compareSync(plainPassword, this.value);
  }
  
  getValue(): string {
    return this.value;
  }
}
```

#### 6.1.3. Repositórios (TypeScript)

**Interface**
```typescript
// src/identity/domain/repositories/user.repository.interface.ts
import { User } from '../entities/user.entity';
import { UserId } from '../entities/user.entity';
import { Email } from '../value-objects/email.vo';

export const USER_REPOSITORY = Symbol('USER_REPOSITORY');

export interface IUserRepository {
  findById(id: UserId): Promise<User | null>;
  findByEmail(email: Email): Promise<User | null>;
  save(user: User): Promise<void>;
  delete(id: UserId): Promise<void>;
  exists(email: Email): Promise<boolean>;
}
```

**Implementação Prisma**
```typescript
// src/identity/infrastructure/persistence/prisma-user.repository.ts
import { Injectable, Inject } from '@nestjs/common';
import { PrismaService } from '../../../shared/infrastructure/prisma/prisma.service';
import { IUserRepository, USER_REPOSITORY } from '../../domain/repositories/user.repository.interface';
import { User, UserId } from '../../domain/entities/user.entity';
import { Email } from '../../domain/value-objects/email.vo';
import { PasswordHash } from '../../domain/value-objects/password-hash.vo';
import { UserName } from '../../domain/value-objects/user-name.vo';

@Injectable()
export class PrismaUserRepository implements IUserRepository {
  constructor(
    @Inject(PrismaService)
    private readonly prisma: PrismaService,
  ) {}
  
  async findById(id: UserId): Promise<User | null> {
    const data = await this.prisma.user.findUnique({
      where: { id: id.toString() },
    });
    
    if (!data) return null;
    
    return this.toDomain(data);
  }
  
  async findByEmail(email: Email): Promise<User | null> {
    const data = await this.prisma.user.findUnique({
      where: { email: email.toString() },
    });
    
    if (!data) return null;
    
    return this.toDomain(data);
  }
  
  async save(user: User): Promise<void> {
    const data = this.toPersistence(user);
    
    await this.prisma.user.upsert({
      where: { id: data.id },
      create: data,
      update: data,
    });
  }
  
  async delete(id: UserId): Promise<void> {
    await this.prisma.user.delete({
      where: { id: id.toString() },
    });
  }
  
  async exists(email: Email): Promise<boolean> {
    const count = await this.prisma.user.count({
      where: { email: email.toString() },
    });
    
    return count > 0;
  }
  
  private toDomain(data: any): User {
    // Mapear de persistência para domínio
    return User.fromPersistence({
      id: new UserId(data.id),
      email: Email.create(data.email),
      passwordHash: PasswordHash.fromHash(data.passwordHash),
      name: UserName.create(data.firstName, data.lastName),
      // ... resto do mapeamento
    });
  }
  
  private toPersistence(user: User): any {
    return {
      id: user.getId().toString(),
      email: user.getEmail().toString(),
      passwordHash: user.getPasswordHash().getValue(),
      // ... resto do mapeamento
    };
  }
}
```

#### 6.1.4. Use Cases (NestJS)

**RegisterUserUseCase**
```typescript
// src/identity/application/use-cases/register-user.use-case.ts
import { Injectable, Inject } from '@nestjs/common';
import { IUserRepository, USER_REPOSITORY } from '../../domain/repositories/user.repository.interface';
import { User } from '../../domain/entities/user.entity';
import { Email } from '../../domain/value-objects/email.vo';
import { PasswordHash } from '../../domain/value-objects/password-hash.vo';
import { UserName } from '../../domain/value-objects/user-name.vo';
import { EventBus } from '@nestjs/cqrs';

@Injectable()
export class RegisterUserUseCase {
  constructor(
    @Inject(USER_REPOSITORY)
    private readonly userRepository: IUserRepository,
    private readonly eventBus: EventBus,
  ) {}
  
  async execute(input: RegisterUserInput): Promise<RegisterUserOutput> {
    const email = Email.create(input.email);
    
    const exists = await this.userRepository.exists(email);
    if (exists) {
      throw new Error('User already exists');
    }
    
    const passwordHash = PasswordHash.fromPlainPassword(input.password);
    const name = UserName.create(input.firstName, input.lastName);
    
    const user = User.create(email, passwordHash, name);
    
    await this.userRepository.save(user);
    
    // Publicar eventos de domínio
    user.commit();
    user.getUncommittedEvents().forEach(event => {
      this.eventBus.publish(event);
    });
    
    return {
      userId: user.getId().toString(),
      email: user.getEmail().toString(),
    };
  }
}

export interface RegisterUserInput {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface RegisterUserOutput {
  userId: string;
  email: string;
}
```

#### 6.1.5. Controllers (NestJS)

**AuthController**
```typescript
// src/identity/presentation/controllers/auth.controller.ts
import { Controller, Post, Body, HttpCode, HttpStatus } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { RegisterUserUseCase, RegisterUserInput } from '../../application/use-cases/register-user.use-case';
import { RegisterUserDto } from '../dtos/register-user.dto';

@ApiTags('Authentication')
@Controller('auth')
export class AuthController {
  constructor(
    private readonly registerUserUseCase: RegisterUserUseCase,
  ) {}
  
  @Post('register')
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: 'Register a new user' })
  @ApiResponse({ status: 201, description: 'User successfully registered' })
  @ApiResponse({ status: 400, description: 'Bad request' })
  async register(@Body() dto: RegisterUserDto) {
    const input: RegisterUserInput = {
      email: dto.email,
      password: dto.password,
      firstName: dto.firstName,
      lastName: dto.lastName,
    };
    
    return this.registerUserUseCase.execute(input);
  }
}
```

**DTO com Validação**
```typescript
// src/identity/presentation/dtos/register-user.dto.ts
import { IsEmail, IsString, MinLength, IsNotEmpty } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class RegisterUserDto {
  @ApiProperty({ example: 'user@example.com' })
  @IsEmail()
  @IsNotEmpty()
  email: string;
  
  @ApiProperty({ example: 'password123', minLength: 8 })
  @IsString()
  @MinLength(8)
  @IsNotEmpty()
  password: string;
  
  @ApiProperty({ example: 'John' })
  @IsString()
  @IsNotEmpty()
  firstName: string;
  
  @ApiProperty({ example: 'Doe' })
  @IsString()
  @IsNotEmpty()
  lastName: string;
}
```

### 6.2. Transaction Context (Core Domain)

#### 6.2.1. Entidade Transaction (TypeScript)

```typescript
// src/transaction/domain/entities/transaction.entity.ts
import { AggregateRoot } from '@nestjs/cqrs';
import { Money } from '../../../shared/domain/value-objects/money';
import { TransactionType } from '../value-objects/transaction-type.vo';
import { TransactionStatus } from '../value-objects/transaction-status.vo';
import { TransactionDescription } from '../value-objects/transaction-description.vo';
import { TransactionApproved } from '../events/transaction-approved.event';
import { TransactionCancelled } from '../events/transaction-cancelled.event';

export class TransactionId {
  constructor(private readonly value: string) {}
  
  toString(): string {
    return this.value;
  }
}

export class Transaction extends AggregateRoot {
  private constructor(
    private readonly id: TransactionId,
    private readonly userId: string,
    private readonly accountId: string,
    private readonly categoryId: string,
    private type: TransactionType,
    private amount: Money,
    private description: TransactionDescription,
    private readonly date: Date,
    private status: TransactionStatus,
    private readonly context: string,
    private readonly createdAt: Date,
    private updatedAt: Date,
  ) {
    super();
  }
  
  static create(
    userId: string,
    accountId: string,
    categoryId: string,
    type: TransactionType,
    amount: Money,
    description: TransactionDescription,
    date: Date,
    context: string,
  ): Transaction {
    return new Transaction(
      new TransactionId(crypto.randomUUID()),
      userId,
      accountId,
      categoryId,
      type,
      amount,
      description,
      date,
      TransactionStatus.PENDING,
      context,
      new Date(),
      new Date(),
    );
  }
  
  approve(): void {
    if (!this.status.canBeApproved()) {
      throw new Error('Transaction cannot be approved');
    }
    
    this.status = TransactionStatus.APPROVED;
    this.updatedAt = new Date();
    this.addDomainEvent(new TransactionApproved(this.id.toString(), this.amount));
  }
  
  cancel(): void {
    if (!this.status.canBeCancelled()) {
      throw new Error('Transaction cannot be cancelled');
    }
    
    this.status = TransactionStatus.CANCELLED;
    this.updatedAt = new Date();
    this.addDomainEvent(new TransactionCancelled(this.id.toString()));
  }
  
  getId(): TransactionId {
    return this.id;
  }
  
  getAmount(): Money {
    return this.amount;
  }
}
```

## 7. Prisma Schema

### 7.1. Schema Principal

```prisma
// prisma/schema.prisma
generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

model User {
  id           String   @id @default(uuid())
  email        String   @unique
  passwordHash String
  firstName    String
  lastName     String
  currency     String   @default("BRL")
  locale       String   @default("pt-BR")
  theme        String   @default("light")
  isActive     Boolean  @default(true)
  createdAt    DateTime @default(now())
  updatedAt    DateTime @updatedAt
  
  accounts       Account[]
  transactions   Transaction[]
  categories     Category[]
  budgets        Budget[]
  investments    Investment[]
  goals          Goal[]
  
  @@index([email])
}

model Account {
  id        String   @id @default(uuid())
  userId    String
  name      String
  type      String   // BANK, WALLET, INVESTMENT, CREDIT_CARD
  balance   Decimal  @default(0) @db.Decimal(15, 2)
  context   String   // PERSONAL, PROFESSIONAL
  isActive  Boolean  @default(true)
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  
  user         User          @relation(fields: [userId], references: [id], onDelete: Cascade)
  transactions Transaction[]
  investments  Investment[]
  
  @@index([userId])
  @@index([userId, context])
}

model Transaction {
  id          String   @id @default(uuid())
  userId      String
  accountId   String
  categoryId  String
  type        String   // INCOME, EXPENSE, TRANSFER
  amount      Decimal  @db.Decimal(15, 2)
  description String?
  date        DateTime
  status      String   @default("PENDING") // PENDING, APPROVED, CANCELLED
  context     String
  tags        String[]
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt
  
  account  Account  @relation(fields: [accountId], references: [id], onDelete: Cascade)
  category Category @relation(fields: [categoryId], references: [id])
  user     User     @relation(fields: [userId], references: [id], onDelete: Cascade)
  
  @@index([userId, date])
  @@index([accountId])
  @@index([categoryId])
  @@index([userId, type, date])
}

model Category {
  id        String   @id @default(uuid())
  userId    String
  name      String
  type      String   // INCOME, EXPENSE
  parentId  String?
  icon      String?
  color     String?
  context   String
  isSystem  Boolean  @default(false)
  isActive  Boolean  @default(true)
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  
  user         User          @relation(fields: [userId], references: [id], onDelete: Cascade)
  transactions Transaction[]
  budgets      Budget[]
  parent       Category?    @relation("CategoryHierarchy", fields: [parentId], references: [id])
  children     Category[]    @relation("CategoryHierarchy")
  
  @@index([userId])
  @@index([userId, type])
}

model Budget {
  id        String   @id @default(uuid())
  userId    String
  categoryId String
  amount    Decimal  @db.Decimal(15, 2)
  period    String   // MONTHLY, YEARLY
  year      Int
  month     Int?
  context   String
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  
  user     User     @relation(fields: [userId], references: [id], onDelete: Cascade)
  category Category @relation(fields: [categoryId], references: [id])
  
  @@unique([userId, categoryId, period, year, month])
  @@index([userId, year, month])
}
```

## 8. ORM: Prisma (Type-Safety)

### 8.1. Por que Prisma?

**Vantagens:**
- ✅ **Type-safety completo** - Geração automática de tipos
- ✅ **Schema-first** - Schema Prisma como fonte da verdade
- ✅ **Migrations automáticas** - Prisma Migrate
- ✅ **Query builder type-safe** - Autocomplete completo
- ✅ **Performance** - Queries otimizadas
- ✅ **Developer Experience** - Prisma Studio para visualização
- ✅ **Documentação excelente** - Guias completos

**Exemplo:**
```prisma
// prisma/schema.prisma
model Account {
  id        String   @id @default(uuid())
  userId    String
  name      String
  type      String
  balance   Decimal  @default(0) @db.Decimal(15, 2)
  context   String
  isActive  Boolean  @default(true)
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  
  user         User          @relation(fields: [userId], references: [id], onDelete: Cascade)
  transactions Transaction[]
  
  @@index([userId])
  @@index([userId, context])
}
```

### 8.2. Prisma Client Type-Safe

```typescript
// Uso do Prisma Client (type-safe)
const account = await prisma.account.findUnique({
  where: { id: accountId },
  include: {
    transactions: {
      take: 10,
      orderBy: { date: 'desc' },
    },
  },
});

// account é totalmente tipado!
// account.transactions[0].amount // TypeScript sabe que existe
```

## 8. Módulos NestJS

### 8.1. Identity Module

```typescript
// src/identity/identity.module.ts
import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { PrismaModule } from '../shared/infrastructure/prisma/prisma.module';
import { AuthController } from './presentation/controllers/auth.controller';
import { RegisterUserUseCase } from './application/use-cases/register-user.use-case';
import { LoginUseCase } from './application/use-cases/login.use-case';
import { PrismaUserRepository } from './infrastructure/persistence/prisma-user.repository';
import { USER_REPOSITORY } from './domain/repositories/user.repository.interface';
import { JwtTokenService } from './infrastructure/services/jwt-token.service';

@Module({
  imports: [CqrsModule, PrismaModule],
  controllers: [AuthController],
  providers: [
    RegisterUserUseCase,
    LoginUseCase,
    {
      provide: USER_REPOSITORY,
      useClass: PrismaUserRepository,
    },
    JwtTokenService,
  ],
  exports: [USER_REPOSITORY],
})
export class IdentityModule {}
```

### 8.2. Transaction Module

```typescript
// src/transaction/transaction.module.ts
import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { PrismaModule } from '../shared/infrastructure/prisma/prisma.module';
import { TransactionController } from './presentation/controllers/transaction.controller';
import { CreateTransactionUseCase } from './application/use-cases/create-transaction.use-case';
import { ApproveTransactionUseCase } from './application/use-cases/approve-transaction.use-case';
import { PrismaTransactionRepository } from './infrastructure/persistence/prisma-transaction.repository';
import { TRANSACTION_REPOSITORY } from './domain/repositories/transaction.repository.interface';

@Module({
  imports: [CqrsModule, PrismaModule],
  controllers: [TransactionController],
  providers: [
    CreateTransactionUseCase,
    ApproveTransactionUseCase,
    {
      provide: TRANSACTION_REPOSITORY,
      useClass: PrismaTransactionRepository,
    },
  ],
  exports: [TRANSACTION_REPOSITORY],
})
export class TransactionModule {}
```

## 9. Event Bus (NestJS)

### 9.1. Configuração

```typescript
// src/app.module.ts
import { Module } from '@nestjs/common';
import { EventEmitterModule } from '@nestjs/event-emitter';
import { CqrsModule } from '@nestjs/cqrs';

@Module({
  imports: [
    EventEmitterModule.forRoot(),
    CqrsModule,
    // ... outros módulos
  ],
})
export class AppModule {}
```

### 9.2. Event Handlers

```typescript
// src/transaction/infrastructure/events/transaction-approved.handler.ts
import { EventsHandler, IEventHandler } from '@nestjs/cqrs';
import { TransactionApproved } from '../../domain/events/transaction-approved.event';
import { AccountRepository } from '../../../account-management/domain/repositories/account.repository.interface';

@EventsHandler(TransactionApproved)
export class TransactionApprovedHandler implements IEventHandler<TransactionApproved> {
  constructor(
    @Inject(ACCOUNT_REPOSITORY)
    private readonly accountRepository: AccountRepository,
  ) {}
  
  async handle(event: TransactionApproved) {
    // Atualizar saldo da conta
    const account = await this.accountRepository.findById(event.accountId);
    if (account) {
      account.credit(event.amount);
      await this.accountRepository.save(account);
    }
  }
}
```

## 10. Testes em NestJS

### 10.1. Testes Unitários

```typescript
// src/transaction/domain/entities/transaction.entity.spec.ts
import { Transaction } from './transaction.entity';
import { TransactionType } from '../value-objects/transaction-type.vo';
import { Money } from '../../../shared/domain/value-objects/money';

describe('Transaction', () => {
  it('should create a transaction', () => {
    const transaction = Transaction.create(
      'user-id',
      'account-id',
      'category-id',
      TransactionType.EXPENSE,
      Money.create(100, 'BRL'),
      TransactionDescription.create('Test'),
      new Date(),
      'PERSONAL',
    );
    
    expect(transaction).toBeDefined();
    expect(transaction.getAmount().getValue()).toBe(100);
  });
  
  it('should approve a transaction', () => {
    const transaction = Transaction.create(/* ... */);
    
    transaction.approve();
    
    expect(transaction.getStatus()).toBe(TransactionStatus.APPROVED);
  });
});
```

### 10.2. Testes de Integração

```typescript
// test/transaction/transaction.integration.spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from '../../src/app.module';

describe('TransactionController (e2e)', () => {
  let app: INestApplication;
  
  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();
    
    app = moduleFixture.createNestApplication();
    await app.init();
  });
  
  it('/transactions (POST)', () => {
    return request(app.getHttpServer())
      .post('/transactions')
      .send({
        accountId: 'account-id',
        categoryId: 'category-id',
        type: 'EXPENSE',
        amount: 100,
        description: 'Test',
      })
      .expect(201);
  });
});
```

## 11. Fases de Desenvolvimento (Backend + Frontend)

### Fase 1: Fundação e MVP (3-4 semanas)

#### Backend:
- Setup do projeto NestJS + Prisma
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
- Event Bus e Domain Events (@nestjs/event-emitter)
- Category Context
- Validações robustas (class-validator)
- Error handling melhorado
- Testes de integração
- Logging estruturado (nestjs-pino)

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

## 12. Performance e Otimizações

### 12.1. Connection Pooling (Prisma)

```typescript
// src/shared/infrastructure/prisma/prisma.service.ts
import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';

@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit, OnModuleDestroy {
  async onModuleInit() {
    await this.$connect();
  }
  
  async onModuleDestroy() {
    await this.$disconnect();
  }
}
```

### 12.2. Cache com Redis

```typescript
// src/shared/infrastructure/cache/redis-cache.service.ts
import { Injectable, Inject } from '@nestjs/common';
import { Redis } from 'ioredis';

@Injectable()
export class RedisCacheService {
  constructor(
    @Inject('REDIS_CLIENT')
    private readonly redis: Redis,
  ) {}
  
  async get<T>(key: string): Promise<T | null> {
    const value = await this.redis.get(key);
    return value ? JSON.parse(value) : null;
  }
  
  async set(key: string, value: any, ttl: number = 3600): Promise<void> {
    await this.redis.setex(key, ttl, JSON.stringify(value));
  }
}
```

### 12.3. Paginação

```typescript
// src/transaction/application/use-cases/list-transactions.use-case.ts
async execute(input: ListTransactionsInput): Promise<ListTransactionsOutput> {
  const { page = 1, limit = 20 } = input;
  const skip = (page - 1) * limit;
  
  const [transactions, total] = await Promise.all([
    this.transactionRepository.findByUserId(input.userId, { skip, take: limit }),
    this.transactionRepository.count({ userId: input.userId }),
  ]);
  
  return {
    data: transactions,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
}
```

## 13. Observabilidade e Monitoramento

### 13.1. Logging Estruturado

```typescript
// src/shared/infrastructure/logger/pino-logger.service.ts
import { Injectable, LoggerService } from '@nestjs/common';
import { PinoLogger } from 'nestjs-pino';

@Injectable()
export class AppLogger implements LoggerService {
  constructor(private readonly logger: PinoLogger) {}
  
  log(message: string, context?: string) {
    this.logger.info({ context }, message);
  }
  
  error(message: string, trace?: string, context?: string) {
    this.logger.error({ context, trace }, message);
  }
  
  warn(message: string, context?: string) {
    this.logger.warn({ context }, message);
  }
  
  debug(message: string, context?: string) {
    this.logger.debug({ context }, message);
  }
}

// Uso
this.logger.log('Transaction created', {
  userId: user.id,
  transactionId: transaction.id,
  amount: transaction.amount,
});
```

### 13.2. Métricas com Prometheus

```typescript
// src/shared/infrastructure/metrics/prometheus.service.ts
import { Injectable } from '@nestjs/common';
import { Counter, Histogram, register } from 'prom-client';

@Injectable()
export class PrometheusService {
  private readonly httpRequestsTotal = new Counter({
    name: 'http_requests_total',
    help: 'Total number of HTTP requests',
    labelNames: ['method', 'endpoint', 'status'],
  });
  
  private readonly httpRequestDuration = new Histogram({
    name: 'http_request_duration_seconds',
    help: 'HTTP request duration in seconds',
    labelNames: ['method', 'endpoint'],
  });
  
  recordRequest(method: string, endpoint: string, status: number, duration: number) {
    this.httpRequestsTotal.inc({ method, endpoint, status: status.toString() });
    this.httpRequestDuration.observe({ method, endpoint }, duration / 1000);
  }
}
```

### 13.3. Tracing com OpenTelemetry

```typescript
// src/shared/infrastructure/tracing/tracing.module.ts
import { Module } from '@nestjs/common';
import { OpenTelemetryModule } from '@nestjs/opentelemetry';
import { JaegerExporter } from '@opentelemetry/exporter-jaeger';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';

@Module({
  imports: [
    OpenTelemetryModule.forRoot({
      serviceName: 'gestao-financeira',
      traceExporter: new JaegerExporter({
        endpoint: 'http://jaeger:14268/api/traces',
      }),
      resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'gestao-financeira',
        [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
      }),
    }),
  ],
})
export class TracingModule {}
```

### 13.4. Health Checks Robustos

```typescript
// src/shared/infrastructure/health/health.controller.ts
import { Controller, Get } from '@nestjs/common';
import { HealthCheck, HealthCheckService, PrismaHealthIndicator } from '@nestjs/terminus';
import { PrismaService } from '../prisma/prisma.service';
import { RedisHealthIndicator } from './redis-health.indicator';

@Controller('health')
export class HealthController {
  constructor(
    private health: HealthCheckService,
    private prisma: PrismaHealthIndicator,
    private redis: RedisHealthIndicator,
    private prismaService: PrismaService,
  ) {}
  
  @Get()
  @HealthCheck()
  check() {
    return this.health.check([
      () => this.prisma.pingCheck('database', this.prismaService),
      () => this.redis.pingCheck('redis'),
    ]);
  }
  
  @Get('live')
  liveness() {
    return { status: 'alive' };
  }
  
  @Get('ready')
  @HealthCheck()
  readiness() {
    return this.check();
  }
}
```

## 14. Segurança

### 14.1. Headers de Segurança

```typescript
// src/shared/infrastructure/middleware/helmet.middleware.ts
import { Injectable, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import helmet from 'helmet';

@Injectable()
export class HelmetMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    helmet({
      contentSecurityPolicy: {
        directives: {
          defaultSrc: ["'self'"],
          styleSrc: ["'self'", "'unsafe-inline'"],
          scriptSrc: ["'self'"],
          imgSrc: ["'self'", 'data:', 'https:'],
        },
      },
      crossOriginEmbedderPolicy: false,
    })(req, res, next);
  }
}
```

### 14.2. Rate Limiting

```typescript
// src/shared/infrastructure/guards/throttler.guard.ts
import { Injectable } from '@nestjs/common';
import { ThrottlerGuard, ThrottlerOptions } from '@nestjs/throttler';
import { ExecutionContext } from '@nestjs/common';

@Injectable()
export class CustomThrottlerGuard extends ThrottlerGuard {
  protected getTracker(req: Record<string, any>): string {
    // Usar user ID se autenticado, senão IP
    return req.user?.id || req.ip;
  }
}

// Configuração
ThrottlerModule.forRoot({
  ttl: 60,
  limit: 100, // 100 requisições por minuto
})
```

### 14.3. Validação de Entrada Robusta

```typescript
// src/shared/infrastructure/pipes/validation.pipe.ts
import { PipeTransform, Injectable, ArgumentMetadata, BadRequestException } from '@nestjs/common';
import { validate } from 'class-validator';
import { plainToInstance } from 'class-transformer';

@Injectable()
export class CustomValidationPipe implements PipeTransform<any> {
  async transform(value: any, { metatype }: ArgumentMetadata) {
    if (!metatype || !this.toValidate(metatype)) {
      return value;
    }
    
    const object = plainToInstance(metatype, value);
    const errors = await validate(object, {
      whitelist: true,
      forbidNonWhitelisted: true,
      transform: true,
    });
    
    if (errors.length > 0) {
      const messages = errors.map(err => 
        Object.values(err.constraints || {}).join(', ')
      );
      throw new BadRequestException(messages);
    }
    
    return object;
  }
  
  private toValidate(metatype: Function): boolean {
    const types: Function[] = [String, Boolean, Number, Array, Object];
    return !types.includes(metatype);
  }
}
```

### 14.4. Proteção contra SQL Injection e XSS

```typescript
// Prisma usa prepared statements automaticamente
// Mas para queries raw:
const transactions = await prisma.$queryRaw`
  SELECT * FROM transactions 
  WHERE user_id = ${userId} 
  AND date BETWEEN ${startDate} AND ${endDate}
`;

// Sanitização de entrada
import { sanitize } from 'sanitize-html';

function sanitizeInput(input: string): string {
  return sanitize(input, {
    allowedTags: [],
    allowedAttributes: {},
  });
}
```

## 15. Performance e Escalabilidade

### 15.1. Paginação Eficiente

```typescript
// src/shared/domain/pagination/pagination.dto.ts
import { IsOptional, IsInt, Min, Max } from 'class-validator';
import { Type } from 'class-transformer';
import { ApiPropertyOptional } from '@nestjs/swagger';

export class PaginationDto {
  @ApiPropertyOptional({ default: 1, minimum: 1 })
  @IsOptional()
  @Type(() => Number)
  @IsInt()
  @Min(1)
  page?: number = 1;
  
  @ApiPropertyOptional({ default: 20, minimum: 1, maximum: 100 })
  @IsOptional()
  @Type(() => Number)
  @IsInt()
  @Min(1)
  @Max(100)
  limit?: number = 20;
  
  get skip(): number {
    return (this.page - 1) * this.limit;
  }
}

// Uso no use case
async execute(input: ListTransactionsInput): Promise<ListTransactionsOutput> {
  const { page = 1, limit = 20 } = input;
  const skip = (page - 1) * limit;
  
  const [transactions, total] = await Promise.all([
    this.transactionRepository.findByUserId(input.userId, { skip, take: limit }),
    this.transactionRepository.count({ userId: input.userId }),
  ]);
  
  return {
    data: transactions,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
}
```

### 15.2. Índices de Banco de Dados

```prisma
// prisma/schema.prisma
model Transaction {
  // ... campos
  
  @@index([userId, date(sort: Desc)])
  @@index([accountId])
  @@index([categoryId])
  @@index([userId, type, date(sort: Desc)])
  @@index([userId, date, type, status]) // Para relatórios
}

model Account {
  // ... campos
  
  @@index([userId])
  @@index([userId, context])
}
```

### 15.3. Cache Estratégico

```typescript
// src/shared/infrastructure/cache/cache.service.ts
import { Injectable, Inject } from '@nestjs/common';
import { Redis } from 'ioredis';

@Injectable()
export class CacheService {
  constructor(
    @Inject('REDIS_CLIENT')
    private readonly redis: Redis,
  ) {}
  
  async getOrSet<T>(
    key: string,
    fn: () => Promise<T>,
    ttl: number = 3600,
  ): Promise<T> {
    const cached = await this.redis.get(key);
    if (cached) {
      return JSON.parse(cached);
    }
    
    const value = await fn();
    await this.redis.setex(key, ttl, JSON.stringify(value));
    return value;
  }
  
  async invalidate(pattern: string): Promise<void> {
    const keys = await this.redis.keys(pattern);
    if (keys.length > 0) {
      await this.redis.del(...keys);
    }
  }
}

// Uso
async getMonthlyReport(userId: string, month: number): Promise<Report> {
  const key = `report:${userId}:${month}`;
  return this.cacheService.getOrSet(
    key,
    () => this.generateReport(userId, month),
    3600, // 1 hora
  );
}
```

### 15.4. Processamento Assíncrono com Bull

```typescript
// src/shared/infrastructure/queue/queue.module.ts
import { Module } from '@nestjs/common';
import { BullModule } from '@nestjs/bull';
import { ProcessRecurringTransactionsProcessor } from './processors/recurring-transactions.processor';

@Module({
  imports: [
    BullModule.forRoot({
      redis: {
        host: process.env.REDIS_HOST,
        port: parseInt(process.env.REDIS_PORT),
      },
    }),
    BullModule.registerQueue({
      name: 'transactions',
    }),
  ],
  providers: [ProcessRecurringTransactionsProcessor],
})
export class QueueModule {}

// Processor
@Processor('transactions')
export class ProcessRecurringTransactionsProcessor {
  @Process('process-recurring')
  async handle(job: Job) {
    // Processar transações recorrentes
  }
}
```

## 16. Documentação da API

### 16.1. Swagger/OpenAPI

```typescript
// src/main.ts
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';

const config = new DocumentBuilder()
  .setTitle('Gestão Financeira API')
  .setDescription('API para gestão financeira pessoal e profissional')
  .setVersion('1.0')
  .addBearerAuth()
  .addTag('Authentication', 'Endpoints de autenticação')
  .addTag('Transactions', 'Endpoints de transações')
  .addTag('Accounts', 'Endpoints de contas')
  .build();

const document = SwaggerModule.createDocument(app, config);
SwaggerModule.setup('api/docs', app, document);

// Controller com decorators Swagger
@ApiTags('Transactions')
@Controller('transactions')
export class TransactionController {
  @Post()
  @ApiOperation({ summary: 'Criar transação' })
  @ApiResponse({ status: 201, description: 'Transação criada com sucesso' })
  @ApiResponse({ status: 400, description: 'Dados inválidos' })
  @ApiBearerAuth()
  async create(@Body() dto: CreateTransactionDto) {
    // ...
  }
}
```

## 17. CI/CD e Deploy

### 17.1. GitHub Actions

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
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run tests
        run: npm run test
      
      - name: Run E2E tests
        run: npm run test:e2e
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
  
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Build
        run: npm run build
      
      - name: Build Docker image
        run: docker build -t gestao-financeira:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin
          docker push gestao-financeira:${{ github.sha }}
```

### 17.2. Deploy em Produção

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  api:
    image: gestao-financeira:latest
    restart: always
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - NODE_ENV=production
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
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
  
  redis:
    image: redis:7-alpine
    restart: always
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
```

## 18. Backup e Disaster Recovery

### 18.1. Estratégia de Backup

```typescript
// src/shared/infrastructure/backup/backup.service.ts
import { Injectable } from '@nestjs/common';
import { Cron, CronExpression } from '@nestjs/schedule';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

@Injectable()
export class BackupService {
  @Cron(CronExpression.EVERY_DAY_AT_2AM)
  async backupDatabase() {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const filename = `backup_${timestamp}.sql`;
    
    try {
      await execAsync(
        `pg_dump -h ${process.env.DB_HOST} -U ${process.env.DB_USER} -d ${process.env.DB_NAME} -f ${filename}`,
      );
      this.logger.log(`Backup completed: ${filename}`);
    } catch (error) {
      this.logger.error(`Backup failed: ${error.message}`);
    }
  }
}
```

## 19. Testes E2E

### 19.1. Testes End-to-End

```typescript
// test/e2e/transaction.e2e-spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from '../../src/app.module';

describe('TransactionController (e2e)', () => {
  let app: INestApplication;
  let token: string;
  
  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();
    
    app = moduleFixture.createNestApplication();
    await app.init();
    
    // Registrar e fazer login
    await request(app.getHttpServer())
      .post('/api/v1/auth/register')
      .send({
        email: 'test@example.com',
        password: 'password123',
        firstName: 'Test',
        lastName: 'User',
      });
    
    const loginResponse = await request(app.getHttpServer())
      .post('/api/v1/auth/login')
      .send({
        email: 'test@example.com',
        password: 'password123',
      });
    
    token = loginResponse.body.token;
  });
  
  it('/transactions (POST)', () => {
    return request(app.getHttpServer())
      .post('/api/v1/transactions')
      .set('Authorization', `Bearer ${token}`)
      .send({
        accountId: 'account-id',
        categoryId: 'category-id',
        type: 'EXPENSE',
        amount: 100,
        description: 'Test',
      })
      .expect(201);
  });
});
```

## 20. Auditoria e Compliance

### 20.1. Log de Auditoria

```typescript
// src/shared/infrastructure/audit/audit.service.ts
import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';

@Injectable()
export class AuditService {
  constructor(private prisma: PrismaService) {}
  
  async log(
    userId: string,
    action: string,
    resource: string,
    resourceId: string,
    metadata?: Record<string, any>,
  ) {
    await this.prisma.auditLog.create({
      data: {
        userId,
        action,
        resource,
        resourceId,
        metadata,
        ipAddress: metadata?.ipAddress,
        userAgent: metadata?.userAgent,
      },
    });
  }
}

// Interceptor de auditoria
@Injectable()
export class AuditInterceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const request = context.switchToHttp().getRequest();
    const user = request.user;
    
    if (user) {
      this.auditService.log(
        user.id,
        request.method,
        this.extractResource(request.path),
        request.params.id,
        {
          ipAddress: request.ip,
          userAgent: request.get('User-Agent'),
        },
      );
    }
    
    return next.handle();
  }
}
```

### 20.2. LGPD/GDPR Compliance

```typescript
// src/identity/application/use-cases/delete-user-data.use-case.ts
@Injectable()
export class DeleteUserDataUseCase {
  async execute(userId: string): Promise<void> {
    // 1. Anonimizar dados pessoais
    await this.userRepository.anonymize(userId);
    
    // 2. Manter dados financeiros agregados (se necessário para compliance fiscal)
    // 3. Registrar ação de exclusão
    await this.auditService.log(userId, 'DELETE', 'user', userId);
    
    // 4. Notificar usuário
    await this.notificationService.send(userId, 'account_deleted');
  }
}

// Exportação de dados
@Injectable()
export class ExportUserDataUseCase {
  async execute(userId: string): Promise<Buffer> {
    const data = {
      user: await this.userRepository.findById(userId),
      transactions: await this.transactionRepository.findByUserId(userId),
      accounts: await this.accountRepository.findByUserId(userId),
      // ... outros dados
    };
    
    return Buffer.from(JSON.stringify(data, null, 2));
  }
}
```

## 21. Escalabilidade e Multi-tenancy (Preparação para Produto)

### 21.1. Estratégia de Escalabilidade

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

### 21.2. Multi-tenancy (Opcional - se virar SaaS)

```typescript
// Estrutura para suportar múltiplos tenants
export class Tenant {
  constructor(
    private readonly id: string,
    private name: string,
    private plan: 'FREE' | 'PREMIUM' | 'ENTERPRISE',
    private settings: TenantSettings,
  ) {}
}

// Isolamento por tenant_id em todas as queries
async findByTenant(tenantId: string, filters: Filters): Promise<Entity[]> {
  return this.prisma.entity.findMany({
    where: {
      tenantId,
      ...filters,
    },
  });
}
```

## 22. Tratamento de Erros Robusto

### 22.1. Erros de Domínio

```typescript
// src/shared/domain/errors/domain.error.ts
export class DomainError extends Error {
  constructor(
    public readonly code: string,
    message: string,
    public readonly details?: Record<string, any>,
  ) {
    super(message);
    this.name = 'DomainError';
  }
}

// Erros específicos
export class InsufficientBalanceError extends DomainError {
  constructor(accountId: string, required: number, available: number) {
    super(
      'INSUFFICIENT_BALANCE',
      'Account balance is insufficient',
      { accountId, required, available },
    );
  }
}

export class TransactionNotFoundError extends DomainError {
  constructor(transactionId: string) {
    super(
      'TRANSACTION_NOT_FOUND',
      'Transaction not found',
      { transactionId },
    );
  }
}
```

### 22.2. Exception Filter Global

```typescript
// src/shared/infrastructure/filters/http-exception.filter.ts
import { ExceptionFilter, Catch, ArgumentsHost, HttpException } from '@nestjs/common';
import { Request, Response } from 'express';

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
    
    const status = exception instanceof HttpException
      ? exception.getStatus()
      : 500;
    
    const message = exception instanceof HttpException
      ? exception.getResponse()
      : 'Internal server error';
    
    response.status(status).json({
      statusCode: status,
      timestamp: new Date().toISOString(),
      path: request.url,
      requestId: request.headers['x-request-id'],
      message,
    });
  }
}
```

## 23. Testes de Performance e Carga

### 23.1. Benchmarks

```typescript
// test/benchmark/transaction.bench.spec.ts
import { Test } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from '../../src/app.module';

describe('Transaction Performance', () => {
  let app: INestApplication;
  
  beforeAll(async () => {
    const moduleFixture = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();
    
    app = moduleFixture.createNestApplication();
    await app.init();
  });
  
  it('should handle 1000 concurrent requests', async () => {
    const promises = Array.from({ length: 1000 }, () =>
      request(app.getHttpServer())
        .get('/api/v1/transactions')
        .set('Authorization', `Bearer ${token}`),
    );
    
    const start = Date.now();
    await Promise.all(promises);
    const duration = Date.now() - start;
    
    expect(duration).toBeLessThan(5000); // Menos de 5 segundos
  });
});
```

### 23.2. Testes de Carga (k6 ou Artillery)

```javascript
// tests/load/transactions.js (k6)
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 200 },
    { duration: '5m', target: 200 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function () {
  const res = http.get('http://localhost:3000/api/v1/transactions', {
    headers: { 'Authorization': 'Bearer ' + __ENV.TOKEN },
  });
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
```

## 24. Versionamento de API

### 24.1. Estrutura de Versionamento

```typescript
// src/app.module.ts
@Module({
  imports: [
    // API v1
    RouterModule.register([
      {
        path: 'api/v1',
        module: V1Module,
      },
      // API v2 (futuro)
      {
        path: 'api/v2',
        module: V2Module,
      },
    ]),
  ],
})

// Deprecation header
@Injectable()
export class DeprecationInterceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const response = context.switchToHttp().getResponse();
    response.setHeader('Deprecation', 'true');
    response.setHeader('Sunset', '2025-12-31');
    response.setHeader('Link', '</api/v2>; rel="successor-version"');
    return next.handle();
  }
}
```

## 13. Deploy e DevOps

### 13.1. Dockerfile Multi-stage

```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Runtime stage
FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/prisma ./prisma
RUN npx prisma generate
CMD ["node", "dist/main"]
```

### 13.2. docker-compose.yml

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "3000:3000"
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

## 14. Considerações Finais

### Vantagens do Node.js + NestJS:
1. ✅ **TypeScript nativo** - Type-safety excelente
2. ✅ **DDD nativo** - NestJS foi feito para DDD
3. ✅ **Prisma** - ORM moderno e type-safe
4. ✅ **Ecossistema** - Muitas ferramentas disponíveis
5. ✅ **Decorators** - Código expressivo
6. ✅ **Testes** - Suporte nativo excelente

### Recomendações:
- Usar **Prisma** para type-safety máximo
- Aproveitar **Decorators** do NestJS
- Implementar **Event Bus** para desacoplamento
- Usar **Swagger** para documentação automática
- Implementar **cache com Redis** para relatórios
- Aproveitar **async/await** para código limpo


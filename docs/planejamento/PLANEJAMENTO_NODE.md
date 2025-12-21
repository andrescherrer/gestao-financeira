# Planejamento DDD - Sistema de Gestão Financeira (Node.js)

## 1. Visão Geral

Sistema de gestão financeira pessoal e profissional desenvolvido em **Node.js com TypeScript** usando **NestJS** e **Prisma**, aproveitando type-safety, DDD nativo e ecossistema moderno.

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

### 3.2. Por que Node.js + NestJS?

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

## 5. Estrutura de Pastas (NestJS DDD)

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

## 11. Fases de Desenvolvimento

### Fase 1: Fundação (2 semanas)
- Setup NestJS + Prisma
- Shared Kernel
- Identity Context
- Account Management Context

### Fase 2: Core Domain (2 semanas)
- Transaction Context
- Integração com Account
- Event Bus
- Testes

### Fase 3: Expansão (3 semanas)
- Category Context
- Budget Context
- Recurring Transactions
- Relatórios básicos

### Fase 4: Funcionalidades Avançadas (3 semanas)
- Investment Context
- Goal Context
- Notification Context
- Dashboard completo

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


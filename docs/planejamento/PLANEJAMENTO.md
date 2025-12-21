# Planejamento DDD - Sistema de Gestão Financeira

## 1. Visão Geral

Sistema de gestão financeira pessoal e profissional desenvolvido com Domain-Driven Design (DDD), focando em modelagem rica de domínio, separação de responsabilidades e manutenibilidade.

## 2. Objetivos

- Controle total de finanças pessoais e profissionais
- Separação clara entre contas pessoais e profissionais
- Análise e relatórios financeiros
- Planejamento orçamentário
- Acompanhamento de metas financeiras
- Arquitetura escalável e manutenível
- Modelagem de domínio rica e expressiva

## 3. Arquitetura DDD - Visão Geral

### 3.1. Estrutura em Camadas

```
┌─────────────────────────────────────┐
│     Presentation Layer              │  (Controllers, DTOs, API)
├─────────────────────────────────────┤
│     Application Layer               │  (Use Cases, Application Services)
├─────────────────────────────────────┤
│     Domain Layer                    │  (Entities, Value Objects, Domain Services)
├─────────────────────────────────────┤
│     Infrastructure Layer            │  (Repositories, External Services, DB)
└─────────────────────────────────────┘
```

### 3.2. Bounded Contexts Identificados

1. **Identity Context** - Autenticação e gestão de usuários
2. **Account Management Context** - Gestão de contas e carteiras
3. **Transaction Context** - Processamento de transações financeiras
4. **Category Context** - Gestão de categorias e taxonomia
5. **Budget Context** - Planejamento e controle orçamentário
6. **Reporting Context** - Análises e relatórios financeiros
7. **Investment Context** - Gestão de investimentos
8. **Goal Context** - Metas e objetivos financeiros
9. **Notification Context** - Notificações e alertas

## 4. Detalhamento dos Bounded Contexts

---

## 4.1. Identity Context

**Responsabilidade**: Autenticação, autorização e gestão de identidade dos usuários.

### 4.1.1. Entidades

**User (Agregado Raiz)**
```typescript
class User {
  private id: UserId;
  private email: Email;
  private passwordHash: PasswordHash;
  private name: UserName;
  private profile: UserProfile; // Value Object
  private createdAt: Date;
  private updatedAt: Date;
  private isActive: boolean;
  
  // Comportamentos
  changePassword(oldPassword: string, newPassword: string): void
  updateProfile(profile: UserProfile): void
  deactivate(): void
  activate(): void
}
```

### 4.1.2. Value Objects

**Email**
```typescript
class Email {
  private value: string;
  
  constructor(email: string) {
    this.validate(email);
    this.value = email.toLowerCase().trim();
  }
  
  private validate(email: string): void {
    // Validação de formato
  }
  
  equals(other: Email): boolean
  toString(): string
}
```

**PasswordHash**
```typescript
class PasswordHash {
  private value: string;
  
  static fromPlainPassword(password: string): PasswordHash
  verify(plainPassword: string): boolean
}
```

**UserName**
```typescript
class UserName {
  private firstName: string;
  private lastName: string;
  
  getFullName(): string
  getInitials(): string
}
```

**UserProfile**
```typescript
class UserProfile {
  private currency: Currency;
  private locale: Locale;
  private timezone: Timezone;
  private theme: Theme; // light/dark
  private dateFormat: DateFormat;
  
  updateCurrency(currency: Currency): void
  updateLocale(locale: Locale): void
}
```

### 4.1.3. Serviços de Domínio

**PasswordService**
- Hash de senhas
- Verificação de senhas
- Validação de força de senha

**TokenService**
- Geração de tokens JWT
- Refresh tokens
- Validação de tokens

### 4.1.4. Repositórios

**IUserRepository**
```typescript
interface IUserRepository {
  findById(id: UserId): Promise<User | null>;
  findByEmail(email: Email): Promise<User | null>;
  save(user: User): Promise<void>;
  delete(id: UserId): Promise<void>;
  exists(email: Email): Promise<boolean>;
}
```

### 4.1.5. Eventos de Domínio

- `UserRegistered` - Usuário registrado
- `UserPasswordChanged` - Senha alterada
- `UserProfileUpdated` - Perfil atualizado
- `UserDeactivated` - Usuário desativado

---

## 4.2. Account Management Context

**Responsabilidade**: Gestão de contas bancárias, carteiras digitais e contas de investimento.

### 4.2.1. Entidades

**Account (Agregado Raiz)**
```typescript
class Account {
  private id: AccountId;
  private userId: UserId;
  private name: AccountName;
  private type: AccountType; // Value Object
  private balance: Money; // Value Object
  private context: AccountContext; // personal/professional
  private isActive: boolean;
  private createdAt: Date;
  private updatedAt: Date;
  
  // Comportamentos
  credit(amount: Money): void
  debit(amount: Money): void
  transferTo(targetAccount: Account, amount: Money): void
  updateName(name: AccountName): void
  deactivate(): void
  activate(): void
  
  // Invariantes
  private ensurePositiveBalance(): void
  private ensureAccountIsActive(): void
}
```

### 4.2.2. Value Objects

**AccountType**
```typescript
class AccountType {
  private value: 'BANK' | 'WALLET' | 'INVESTMENT' | 'CREDIT_CARD';
  
  isBank(): boolean
  isWallet(): boolean
  isInvestment(): boolean
  isCreditCard(): boolean
}
```

**Money**
```typescript
class Money {
  private amount: number;
  private currency: Currency;
  
  constructor(amount: number, currency: Currency) {
    if (amount < 0) throw new Error('Amount cannot be negative');
    this.amount = amount;
    this.currency = currency;
  }
  
  add(other: Money): Money
  subtract(other: Money): Money
  multiply(factor: number): Money
  isGreaterThan(other: Money): boolean
  isLessThan(other: Money): boolean
  equals(other: Money): boolean
  toNumber(): number
}
```

**AccountName**
```typescript
class AccountName {
  private value: string;
  
  constructor(name: string) {
    if (!name || name.trim().length < 3) {
      throw new Error('Account name must have at least 3 characters');
    }
    this.value = name.trim();
  }
  
  toString(): string
}
```

**AccountContext**
```typescript
class AccountContext {
  private value: 'PERSONAL' | 'PROFESSIONAL';
  
  isPersonal(): boolean
  isProfessional(): boolean
  equals(other: AccountContext): boolean
}
```

### 4.2.3. Serviços de Domínio

**AccountBalanceService**
- Cálculo de saldo consolidado
- Validação de saldo antes de débito
- Sincronização de saldo

### 4.2.4. Repositórios

**IAccountRepository**
```typescript
interface IAccountRepository {
  findById(id: AccountId): Promise<Account | null>;
  findByUserId(userId: UserId): Promise<Account[]>;
  findByUserIdAndContext(userId: UserId, context: AccountContext): Promise<Account[]>;
  save(account: Account): Promise<void>;
  delete(id: AccountId): Promise<void>;
  exists(id: AccountId): Promise<boolean>;
}
```

### 4.2.5. Eventos de Domínio

- `AccountCreated` - Conta criada
- `AccountBalanceUpdated` - Saldo atualizado
- `AccountDeactivated` - Conta desativada
- `AccountNameChanged` - Nome da conta alterado

---

## 4.3. Transaction Context

**Responsabilidade**: Processamento de transações financeiras (receitas, despesas, transferências).

### 4.3.1. Entidades

**Transaction (Agregado Raiz)**
```typescript
class Transaction {
  private id: TransactionId;
  private userId: UserId;
  private accountId: AccountId;
  private categoryId: CategoryId;
  private type: TransactionType; // Value Object
  private amount: Money;
  private description: TransactionDescription;
  private date: TransactionDate;
  private tags: Tag[]; // Value Objects
  private attachment: Attachment | null;
  private context: AccountContext;
  private status: TransactionStatus;
  private createdAt: Date;
  private updatedAt: Date;
  
  // Comportamentos
  approve(): void
  cancel(): void
  updateAmount(amount: Money): void
  updateDescription(description: TransactionDescription): void
  addTag(tag: Tag): void
  removeTag(tag: Tag): void
  attachFile(attachment: Attachment): void
  
  // Invariantes
  private ensureTransactionIsValid(): void
  private ensureAmountIsPositive(): void
}
```

**RecurringTransaction (Agregado)**
```typescript
class RecurringTransaction {
  private id: RecurringTransactionId;
  private userId: UserId;
  private accountId: AccountId;
  private categoryId: CategoryId;
  private type: TransactionType;
  private amount: Money;
  private description: TransactionDescription;
  private frequency: RecurrenceFrequency; // Value Object
  private nextExecutionDate: Date;
  private endDate: Date | null;
  private context: AccountContext;
  private isActive: boolean;
  
  // Comportamentos
  generateNextTransaction(): Transaction
  calculateNextExecutionDate(): Date
  deactivate(): void
  updateFrequency(frequency: RecurrenceFrequency): void
}
```

**TransactionInstallment (Agregado)**
```typescript
class TransactionInstallment {
  private id: InstallmentId;
  private parentTransactionId: TransactionId;
  private installmentNumber: number;
  private totalInstallments: number;
  private amount: Money;
  private dueDate: Date;
  private status: InstallmentStatus;
  
  // Comportamentos
  pay(): void
  markAsOverdue(): void
  updateDueDate(date: Date): void
}
```

### 4.3.2. Value Objects

**TransactionType**
```typescript
class TransactionType {
  private value: 'INCOME' | 'EXPENSE' | 'TRANSFER';
  
  isIncome(): boolean
  isExpense(): boolean
  isTransfer(): boolean
  affectsBalance(): boolean
}
```

**TransactionDescription**
```typescript
class TransactionDescription {
  private value: string;
  
  constructor(description: string) {
    if (description.length > 500) {
      throw new Error('Description too long');
    }
    this.value = description.trim();
  }
  
  toString(): string
}
```

**TransactionDate**
```typescript
class TransactionDate {
  private value: Date;
  
  constructor(date: Date) {
    if (date > new Date()) {
      throw new Error('Transaction date cannot be in the future');
    }
    this.value = date;
  }
  
  isInPeriod(start: Date, end: Date): boolean
  getYear(): number
  getMonth(): number
  toDate(): Date
}
```

**Tag**
```typescript
class Tag {
  private value: string;
  
  constructor(tag: string) {
    if (!tag || tag.length > 50) {
      throw new Error('Invalid tag');
    }
    this.value = tag.toLowerCase().trim();
  }
  
  equals(other: Tag): boolean
  toString(): string
}
```

**RecurrenceFrequency**
```typescript
class RecurrenceFrequency {
  private value: 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'YEARLY';
  private dayOfMonth?: number;
  private dayOfWeek?: number;
  
  calculateNextDate(from: Date): Date
  getDescription(): string
}
```

**TransactionStatus**
```typescript
class TransactionStatus {
  private value: 'PENDING' | 'APPROVED' | 'CANCELLED';
  
  isPending(): boolean
  isApproved(): boolean
  isCancelled(): boolean
  canBeCancelled(): boolean
}
```

**Attachment**
```typescript
class Attachment {
  private url: string;
  private fileName: string;
  private fileType: string;
  private fileSize: number;
  
  isValid(): boolean
  getUrl(): string
}
```

### 4.3.3. Serviços de Domínio

**TransactionProcessingService**
- Processar transação e atualizar saldo da conta
- Validar regras de negócio antes de processar
- Reverter transação cancelada

**RecurringTransactionScheduler**
- Gerar transações recorrentes
- Calcular próxima data de execução
- Validar recorrências ativas

**TransferService**
- Processar transferências entre contas
- Validar contas de origem e destino
- Garantir atomicidade da operação

### 4.3.4. Repositórios

**ITransactionRepository**
```typescript
interface ITransactionRepository {
  findById(id: TransactionId): Promise<Transaction | null>;
  findByAccountId(accountId: AccountId, filters?: TransactionFilters): Promise<Transaction[]>;
  findByUserId(userId: UserId, filters?: TransactionFilters): Promise<Transaction[]>;
  findByDateRange(userId: UserId, start: Date, end: Date): Promise<Transaction[]>;
  findByCategory(categoryId: CategoryId, filters?: TransactionFilters): Promise<Transaction[]>;
  save(transaction: Transaction): Promise<void>;
  delete(id: TransactionId): Promise<void>;
  count(filters: TransactionFilters): Promise<number>;
}
```

**IRecurringTransactionRepository**
```typescript
interface IRecurringTransactionRepository {
  findById(id: RecurringTransactionId): Promise<RecurringTransaction | null>;
  findByUserId(userId: UserId): Promise<RecurringTransaction[]>;
  findActiveByDate(date: Date): Promise<RecurringTransaction[]>;
  save(recurringTransaction: RecurringTransaction): Promise<void>;
  delete(id: RecurringTransactionId): Promise<void>;
}
```

### 4.3.5. Eventos de Domínio

- `TransactionCreated` - Transação criada
- `TransactionApproved` - Transação aprovada
- `TransactionCancelled` - Transação cancelada
- `TransactionAmountUpdated` - Valor atualizado
- `RecurringTransactionGenerated` - Transação recorrente gerada
- `TransferExecuted` - Transferência executada

---

## 4.4. Category Context

**Responsabilidade**: Gestão de categorias e taxonomia financeira.

### 4.4.1. Entidades

**Category (Agregado Raiz)**
```typescript
class Category {
  private id: CategoryId;
  private userId: UserId;
  private name: CategoryName;
  private type: CategoryType; // Value Object
  private parentId: CategoryId | null;
  private icon: Icon; // Value Object
  private color: Color; // Value Object
  private context: AccountContext;
  private isSystem: boolean; // Categoria padrão do sistema
  private isActive: boolean;
  private createdAt: Date;
  private updatedAt: Date;
  
  // Comportamentos
  updateName(name: CategoryName): void
  updateIcon(icon: Icon): void
  updateColor(color: Color): void
  moveToParent(parentId: CategoryId | null): void
  deactivate(): void
  activate(): void
  
  // Invariantes
  private ensureNotSystemCategory(): void
  private ensureValidParent(): void
}
```

### 4.4.2. Value Objects

**CategoryName**
```typescript
class CategoryName {
  private value: string;
  
  constructor(name: string) {
    if (!name || name.length < 2 || name.length > 100) {
      throw new Error('Invalid category name');
    }
    this.value = name.trim();
  }
  
  toString(): string
}
```

**CategoryType**
```typescript
class CategoryType {
  private value: 'INCOME' | 'EXPENSE';
  
  isIncome(): boolean
  isExpense(): boolean
  matchesTransactionType(transactionType: TransactionType): boolean
}
```

**Icon**
```typescript
class Icon {
  private value: string; // Nome do ícone ou URL
  
  isValid(): boolean
  toString(): string
}
```

**Color**
```typescript
class Color {
  private hex: string;
  
  constructor(hex: string) {
    this.validateHex(hex);
    this.hex = hex.toUpperCase();
  }
  
  private validateHex(hex: string): void
  toHex(): string
  toRgb(): { r: number; g: number; b: number }
}
```

### 4.4.3. Serviços de Domínio

**CategoryHierarchyService**
- Validar hierarquia de categorias
- Prevenir ciclos na árvore
- Calcular caminho completo da categoria

**DefaultCategoryService**
- Criar categorias padrão do sistema
- Migrar categorias padrão para novos usuários

### 4.4.4. Repositórios

**ICategoryRepository**
```typescript
interface ICategoryRepository {
  findById(id: CategoryId): Promise<Category | null>;
  findByUserId(userId: UserId, filters?: CategoryFilters): Promise<Category[]>;
  findByType(userId: UserId, type: CategoryType): Promise<Category[]>;
  findByParentId(parentId: CategoryId): Promise<Category[]>;
  findRootCategories(userId: UserId): Promise<Category[]>;
  findSystemCategories(): Promise<Category[]>;
  save(category: Category): Promise<void>;
  delete(id: CategoryId): Promise<void>;
  exists(id: CategoryId): Promise<boolean>;
}
```

### 4.4.5. Eventos de Domínio

- `CategoryCreated` - Categoria criada
- `CategoryUpdated` - Categoria atualizada
- `CategoryDeactivated` - Categoria desativada
- `CategoryMoved` - Categoria movida na hierarquia

---

## 4.5. Budget Context

**Responsabilidade**: Planejamento e controle orçamentário.

### 4.5.1. Entidades

**Budget (Agregado Raiz)**
```typescript
class Budget {
  private id: BudgetId;
  private userId: UserId;
  private categoryId: CategoryId;
  private amount: Money;
  private period: BudgetPeriod; // Value Object
  private context: AccountContext;
  private alerts: BudgetAlert[]; // Value Objects
  private createdAt: Date;
  private updatedAt: Date;
  
  // Comportamentos
  updateAmount(amount: Money): void
  addAlert(alert: BudgetAlert): void
  removeAlert(alertId: string): void
  calculateUsage(transactions: Transaction[]): BudgetUsage
  checkAlerts(usage: BudgetUsage): BudgetAlert[]
  isActive(): boolean
  
  // Invariantes
  private ensureAmountIsPositive(): void
  private ensureValidPeriod(): void
}
```

**BudgetUsage (Value Object Calculado)**
```typescript
class BudgetUsage {
  private budgeted: Money;
  private spent: Money;
  private remaining: Money;
  private percentageUsed: number;
  
  calculatePercentage(): number
  isExceeded(): boolean
  isNearLimit(threshold: number): boolean
  getRemaining(): Money
}
```

### 4.5.2. Value Objects

**BudgetPeriod**
```typescript
class BudgetPeriod {
  private type: 'MONTHLY' | 'YEARLY';
  private year: number;
  private month: number | null;
  
  constructor(type: 'MONTHLY' | 'YEARLY', year: number, month?: number) {
    this.type = type;
    this.year = year;
    this.month = type === 'MONTHLY' ? month : null;
  }
  
  isMonthly(): boolean
  isYearly(): boolean
  getStartDate(): Date
  getEndDate(): Date
  includes(date: Date): boolean
}
```

**BudgetAlert**
```typescript
class BudgetAlert {
  private id: string;
  private threshold: number; // Percentual (0-100)
  private type: 'WARNING' | 'CRITICAL';
  private isActive: boolean;
  
  shouldTrigger(percentageUsed: number): boolean
  getType(): string
}
```

### 4.5.3. Serviços de Domínio

**BudgetCalculationService**
- Calcular uso do orçamento
- Agregar transações por categoria
- Calcular percentual utilizado

**BudgetAlertService**
- Verificar alertas de orçamento
- Gerar notificações quando necessário

### 4.5.4. Repositórios

**IBudgetRepository**
```typescript
interface IBudgetRepository {
  findById(id: BudgetId): Promise<Budget | null>;
  findByUserId(userId: UserId, filters?: BudgetFilters): Promise<Budget[]>;
  findByCategoryAndPeriod(categoryId: CategoryId, period: BudgetPeriod): Promise<Budget | null>;
  findByPeriod(userId: UserId, period: BudgetPeriod): Promise<Budget[]>;
  save(budget: Budget): Promise<void>;
  delete(id: BudgetId): Promise<void>;
  exists(id: BudgetId): Promise<boolean>;
}
```

### 4.5.5. Eventos de Domínio

- `BudgetCreated` - Orçamento criado
- `BudgetUpdated` - Orçamento atualizado
- `BudgetExceeded` - Orçamento excedido
- `BudgetAlertTriggered` - Alerta de orçamento acionado

---

## 4.6. Reporting Context

**Responsabilidade**: Análises, relatórios e visualizações financeiras.

### 4.6.1. Entidades

**Report (Agregado Raiz)**
```typescript
class Report {
  private id: ReportId;
  private userId: UserId;
  private type: ReportType; // Value Object
  private period: ReportPeriod; // Value Object
  private filters: ReportFilters; // Value Object
  private data: ReportData; // Value Object
  private format: ReportFormat; // Value Object
  private generatedAt: Date;
  
  // Comportamentos
  generate(transactions: Transaction[]): void
  export(format: ReportFormat): ExportResult
  updateFilters(filters: ReportFilters): void
}
```

**Dashboard (Agregado)**
```typescript
class Dashboard {
  private userId: UserId;
  private period: ReportPeriod;
  private widgets: DashboardWidget[]; // Value Objects
  
  // Comportamentos
  addWidget(widget: DashboardWidget): void
  removeWidget(widgetId: string): void
  updateWidgetData(widgetId: string, data: any): void
  refresh(transactions: Transaction[], accounts: Account[]): void
}
```

### 4.6.2. Value Objects

**ReportType**
```typescript
class ReportType {
  private value: 'INCOME_EXPENSE' | 'BY_CATEGORY' | 'CASH_FLOW' | 'PATRIMONY' | 'CUSTOM';
  
  getDescription(): string
  requiresGrouping(): boolean
}
```

**ReportPeriod**
```typescript
class ReportPeriod {
  private startDate: Date;
  private endDate: Date;
  private type: 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'YEARLY' | 'CUSTOM';
  
  getDays(): number
  includes(date: Date): boolean
  getDescription(): string
}
```

**ReportFilters**
```typescript
class ReportFilters {
  private accountIds: AccountId[];
  private categoryIds: CategoryId[];
  private transactionTypes: TransactionType[];
  private tags: Tag[];
  private context: AccountContext | null;
  private minAmount: Money | null;
  private maxAmount: Money | null;
  
  matches(transaction: Transaction): boolean
  isEmpty(): boolean
}
```

**ReportData**
```typescript
class ReportData {
  private summary: ReportSummary; // Value Object
  private transactions: Transaction[];
  private groupedData: GroupedData[]; // Por categoria, por período, etc.
  private charts: ChartData[]; // Dados para gráficos
  
  getTotalIncome(): Money
  getTotalExpense(): Money
  getBalance(): Money
  getByCategory(): Map<CategoryId, Money>
}
```

**DashboardWidget**
```typescript
class DashboardWidget {
  private id: string;
  private type: 'BALANCE' | 'INCOME_EXPENSE' | 'CATEGORY_CHART' | 'RECENT_TRANSACTIONS' | 'BUDGET_PROGRESS';
  private position: WidgetPosition;
  private size: WidgetSize;
  private config: WidgetConfig;
  
  render(data: any): WidgetRenderResult
}
```

### 4.6.3. Serviços de Domínio

**ReportGenerationService**
- Agregar dados de transações
- Calcular métricas financeiras
- Gerar visualizações

**DataAggregationService**
- Agrupar transações por categoria
- Agrupar transações por período
- Calcular totais e médias

**ChartDataService**
- Preparar dados para gráficos
- Formatar dados para visualização
- Calcular percentuais e proporções

### 4.6.4. Repositórios

**IReportRepository**
```typescript
interface IReportRepository {
  findById(id: ReportId): Promise<Report | null>;
  findByUserId(userId: UserId): Promise<Report[]>;
  save(report: Report): Promise<void>;
  delete(id: ReportId): Promise<void>;
}
```

**IDashboardRepository**
```typescript
interface IDashboardRepository {
  findByUserId(userId: UserId): Promise<Dashboard | null>;
  save(dashboard: Dashboard): Promise<void>;
}
```

### 4.6.5. Eventos de Domínio

- `ReportGenerated` - Relatório gerado
- `ReportExported` - Relatório exportado
- `DashboardRefreshed` - Dashboard atualizado

---

## 4.7. Investment Context

**Responsabilidade**: Gestão de investimentos e acompanhamento de rentabilidade.

### 4.7.1. Entidades

**Investment (Agregado Raiz)**
```typescript
class Investment {
  private id: InvestmentId;
  private userId: UserId;
  private accountId: AccountId;
  private type: InvestmentType; // Value Object
  private name: InvestmentName;
  private purchaseDate: Date;
  private purchaseAmount: Money;
  private currentValue: Money;
  private quantity: number;
  private context: AccountContext;
  private createdAt: Date;
  private updatedAt: Date;
  
  // Comportamentos
  updateCurrentValue(value: Money): void
  calculateReturn(): InvestmentReturn // Value Object
  calculateReturnPercentage(): number
  addQuantity(quantity: number): void
  removeQuantity(quantity: number): void
  
  // Invariantes
  private ensurePositiveQuantity(): void
  private ensureValidDates(): void
}
```

### 4.7.2. Value Objects

**InvestmentType**
```typescript
class InvestmentType {
  private value: 'STOCK' | 'FUND' | 'CDB' | 'TREASURY' | 'CRYPTO' | 'OTHER';
  
  getDescription(): string
  requiresQuantity(): boolean
  hasVariableValue(): boolean
}
```

**InvestmentName**
```typescript
class InvestmentName {
  private value: string;
  private ticker?: string; // Para ações
  
  constructor(name: string, ticker?: string) {
    if (!name || name.length < 2) {
      throw new Error('Invalid investment name');
    }
    this.value = name;
    this.ticker = ticker;
  }
  
  toString(): string
  getTicker(): string | null
}
```

**InvestmentReturn**
```typescript
class InvestmentReturn {
  private absolute: Money; // Ganho/perda absoluto
  private percentage: number; // Percentual de retorno
  
  constructor(absolute: Money, percentage: number) {
    this.absolute = absolute;
    this.percentage = percentage;
  }
  
  isPositive(): boolean
  isNegative(): boolean
  getAbsolute(): Money
  getPercentage(): number
}
```

### 4.7.3. Serviços de Domínio

**InvestmentValuationService**
- Calcular valor atual do investimento
- Atualizar cotações (se integrado com API externa)
- Calcular rentabilidade

**PortfolioService**
- Calcular valor total do portfólio
- Calcular distribuição por tipo
- Calcular retorno total

### 4.7.4. Repositórios

**IInvestmentRepository**
```typescript
interface IInvestmentRepository {
  findById(id: InvestmentId): Promise<Investment | null>;
  findByUserId(userId: UserId): Promise<Investment[]>;
  findByAccountId(accountId: AccountId): Promise<Investment[]>;
  findByType(userId: UserId, type: InvestmentType): Promise<Investment[]>;
  save(investment: Investment): Promise<void>;
  delete(id: InvestmentId): Promise<void>;
}
```

### 4.7.5. Eventos de Domínio

- `InvestmentCreated` - Investimento criado
- `InvestmentValueUpdated` - Valor atualizado
- `InvestmentReturnCalculated` - Retorno calculado

---

## 4.8. Goal Context

**Responsabilidade**: Metas e objetivos financeiros.

### 4.8.1. Entidades

**Goal (Agregado Raiz)**
```typescript
class Goal {
  private id: GoalId;
  private userId: UserId;
  private name: GoalName;
  private targetAmount: Money;
  private currentAmount: Money;
  private deadline: Date;
  private context: AccountContext;
  private status: GoalStatus; // Value Object
  private milestones: Milestone[]; // Value Objects
  private createdAt: Date;
  private updatedAt: Date;
  
  // Comportamentos
  addContribution(amount: Money): void
  updateProgress(amount: Money): void
  addMilestone(milestone: Milestone): void
  checkStatus(): GoalStatus
  calculateProgress(): number
  calculateRemainingDays(): number
  isCompleted(): boolean
  isOverdue(): boolean
  
  // Invariantes
  private ensureCurrentAmountNotExceedsTarget(): void
  private ensureValidDeadline(): void
}
```

### 4.8.2. Value Objects

**GoalName**
```typescript
class GoalName {
  private value: string;
  
  constructor(name: string) {
    if (!name || name.length < 3 || name.length > 200) {
      throw new Error('Invalid goal name');
    }
    this.value = name.trim();
  }
  
  toString(): string
}
```

**GoalStatus**
```typescript
class GoalStatus {
  private value: 'IN_PROGRESS' | 'COMPLETED' | 'OVERDUE' | 'CANCELLED';
  
  isInProgress(): boolean
  isCompleted(): boolean
  isOverdue(): boolean
  canBeCancelled(): boolean
}
```

**Milestone**
```typescript
class Milestone {
  private id: string;
  private name: string;
  private targetAmount: Money;
  private achievedAt: Date | null;
  
  markAsAchieved(date: Date): void
  isAchieved(): boolean
  getProgress(currentAmount: Money): number
}
```

### 4.8.3. Serviços de Domínio

**GoalProgressService**
- Calcular progresso da meta
- Verificar status (completa, atrasada, etc.)
- Calcular contribuições necessárias

**GoalNotificationService**
- Verificar prazos próximos
- Notificar sobre progresso
- Alertar sobre metas em risco

### 4.8.4. Repositórios

**IGoalRepository**
```typescript
interface IGoalRepository {
  findById(id: GoalId): Promise<Goal | null>;
  findByUserId(userId: UserId): Promise<Goal[]>;
  findByStatus(userId: UserId, status: GoalStatus): Promise<Goal[]>;
  findUpcomingDeadlines(userId: UserId, days: number): Promise<Goal[]>;
  save(goal: Goal): Promise<void>;
  delete(id: GoalId): Promise<void>;
}
```

### 4.8.5. Eventos de Domínio

- `GoalCreated` - Meta criada
- `GoalProgressUpdated` - Progresso atualizado
- `GoalCompleted` - Meta completada
- `GoalOverdue` - Meta vencida
- `MilestoneAchieved` - Marco alcançado

---

## 4.9. Notification Context

**Responsabilidade**: Notificações, alertas e comunicações com o usuário.

### 4.9.1. Entidades

**Notification (Agregado Raiz)**
```typescript
class Notification {
  private id: NotificationId;
  private userId: UserId;
  private type: NotificationType; // Value Object
  private title: string;
  private message: string;
  private priority: NotificationPriority; // Value Object
  private readAt: Date | null;
  private actionUrl: string | null;
  private metadata: NotificationMetadata; // Value Object
  private createdAt: Date;
  
  // Comportamentos
  markAsRead(): void
  isRead(): boolean
  isUnread(): boolean
  getAge(): number // Em horas/dias
}
```

### 4.9.2. Value Objects

**NotificationType**
```typescript
class NotificationType {
  private value: 'BUDGET_ALERT' | 'GOAL_UPDATE' | 'TRANSACTION_REMINDER' | 'SYSTEM' | 'SECURITY';
  
  getIcon(): string
  getColor(): string
}
```

**NotificationPriority**
```typescript
class NotificationPriority {
  private value: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
  
  isUrgent(): boolean
  requiresImmediateAttention(): boolean
}
```

**NotificationMetadata**
```typescript
class NotificationMetadata {
  private data: Map<string, any>;
  
  set(key: string, value: any): void
  get(key: string): any
  has(key: string): boolean
}
```

### 4.9.3. Serviços de Domínio

**NotificationService**
- Criar notificações
- Agrupar notificações similares
- Limpar notificações antigas

**NotificationDeliveryService**
- Enviar notificações por diferentes canais (email, push, in-app)
- Agendar notificações
- Gerenciar preferências de entrega

### 4.9.4. Repositórios

**INotificationRepository**
```typescript
interface INotificationRepository {
  findById(id: NotificationId): Promise<Notification | null>;
  findByUserId(userId: UserId, filters?: NotificationFilters): Promise<Notification[]>;
  findUnreadByUserId(userId: UserId): Promise<Notification[]>;
  countUnread(userId: UserId): Promise<number>;
  save(notification: Notification): Promise<void>;
  markAsRead(id: NotificationId): Promise<void>;
  markAllAsRead(userId: UserId): Promise<void>;
  delete(id: NotificationId): Promise<void>;
  deleteOld(olderThan: Date): Promise<void>;
}
```

### 4.9.5. Eventos de Domínio

- `NotificationCreated` - Notificação criada
- `NotificationRead` - Notificação lida
- `NotificationDeleted` - Notificação deletada

---

## 5. Integração entre Bounded Contexts

### 5.1. Context Mapping

```
Identity Context
    │
    ├──→ Account Management (UserId)
    ├──→ Transaction (UserId)
    ├──→ Category (UserId)
    └──→ ... (todos os contextos)

Account Management ←→ Transaction
    │                      │
    └──────────────────────┘
         (AccountId)

Category ←→ Transaction
    │           │
    └───────────┘
    (CategoryId)

Transaction → Reporting
    │              │
    └──────────────┘
    (dados agregados)

Budget ←→ Transaction
    │          │
    └──────────┘
    (CategoryId)

Goal ←→ Transaction
    │         │
    └─────────┘
    (contribuições)
```

### 5.2. Eventos de Integração

**Event Bus / Message Broker**
- Comunicação assíncrona entre contextos
- Desacoplamento de serviços
- Eventos de domínio publicados

**Eventos Principais de Integração:**

1. **TransactionCreated** → 
   - Reporting: Atualizar agregados
   - Budget: Verificar uso de orçamento
   - Goal: Atualizar progresso de metas
   - Account: Atualizar saldo

2. **AccountBalanceUpdated** →
   - Reporting: Atualizar dashboard
   - Notification: Alertar se saldo baixo

3. **BudgetExceeded** →
   - Notification: Criar alerta

4. **GoalCompleted** →
   - Notification: Notificar usuário

### 5.3. Shared Kernel

**Value Objects Compartilhados:**
- `Money` - Usado em múltiplos contextos
- `Currency` - Moeda
- `AccountContext` - Contexto pessoal/profissional
- `UserId` - Identificador de usuário

## 6. Arquitetura de Implementação

### 6.1. Estrutura de Pastas (DDD)

```
gestao-financeira/
├── src/
│   ├── shared/                          # Shared Kernel
│   │   ├── domain/
│   │   │   ├── value-objects/
│   │   │   │   ├── Money.ts
│   │   │   │   ├── Currency.ts
│   │   │   │   └── AccountContext.ts
│   │   │   └── events/
│   │   └── infrastructure/
│   │       └── event-bus/
│   │
│   ├── identity/                        # Identity Context
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   └── User.ts
│   │   │   ├── value-objects/
│   │   │   │   ├── Email.ts
│   │   │   │   ├── PasswordHash.ts
│   │   │   │   └── UserName.ts
│   │   │   ├── services/
│   │   │   │   ├── PasswordService.ts
│   │   │   │   └── TokenService.ts
│   │   │   ├── repositories/
│   │   │   │   └── IUserRepository.ts
│   │   │   └── events/
│   │   │       ├── UserRegistered.ts
│   │   │       └── UserPasswordChanged.ts
│   │   ├── application/
│   │   │   ├── use-cases/
│   │   │   │   ├── RegisterUserUseCase.ts
│   │   │   │   ├── LoginUseCase.ts
│   │   │   │   └── UpdateProfileUseCase.ts
│   │   │   └── dtos/
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   └── PrismaUserRepository.ts
│   │   │   ├── services/
│   │   │   │   └── JwtTokenService.ts
│   │   │   └── events/
│   │   │       └── DomainEventPublisher.ts
│   │   └── presentation/
│   │       ├── controllers/
│   │       │   └── AuthController.ts
│   │       └── dtos/
│   │
│   ├── account-management/              # Account Management Context
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   └── Account.ts
│   │   │   ├── value-objects/
│   │   │   │   ├── AccountType.ts
│   │   │   │   └── AccountName.ts
│   │   │   ├── services/
│   │   │   │   └── AccountBalanceService.ts
│   │   │   ├── repositories/
│   │   │   │   └── IAccountRepository.ts
│   │   │   └── events/
│   │   │       ├── AccountCreated.ts
│   │   │       └── AccountBalanceUpdated.ts
│   │   ├── application/
│   │   │   ├── use-cases/
│   │   │   │   ├── CreateAccountUseCase.ts
│   │   │   │   ├── UpdateAccountUseCase.ts
│   │   │   │   └── GetAccountsUseCase.ts
│   │   │   └── dtos/
│   │   ├── infrastructure/
│   │   │   └── persistence/
│   │   │       └── PrismaAccountRepository.ts
│   │   └── presentation/
│   │       └── controllers/
│   │           └── AccountController.ts
│   │
│   ├── transaction/                     # Transaction Context
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   ├── Transaction.ts
│   │   │   │   ├── RecurringTransaction.ts
│   │   │   │   └── TransactionInstallment.ts
│   │   │   ├── value-objects/
│   │   │   │   ├── TransactionType.ts
│   │   │   │   ├── TransactionDate.ts
│   │   │   │   ├── Tag.ts
│   │   │   │   └── RecurrenceFrequency.ts
│   │   │   ├── services/
│   │   │   │   ├── TransactionProcessingService.ts
│   │   │   │   ├── RecurringTransactionScheduler.ts
│   │   │   │   └── TransferService.ts
│   │   │   ├── repositories/
│   │   │   │   ├── ITransactionRepository.ts
│   │   │   │   └── IRecurringTransactionRepository.ts
│   │   │   └── events/
│   │   │       ├── TransactionCreated.ts
│   │   │       └── TransactionApproved.ts
│   │   ├── application/
│   │   │   ├── use-cases/
│   │   │   │   ├── CreateTransactionUseCase.ts
│   │   │   │   ├── ApproveTransactionUseCase.ts
│   │   │   │   ├── CancelTransactionUseCase.ts
│   │   │   │   └── CreateRecurringTransactionUseCase.ts
│   │   │   └── dtos/
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   ├── PrismaTransactionRepository.ts
│   │   │   │   └── PrismaRecurringTransactionRepository.ts
│   │   │   └── schedulers/
│   │   │       └── RecurringTransactionJob.ts
│   │   └── presentation/
│   │       └── controllers/
│   │           └── TransactionController.ts
│   │
│   ├── category/                        # Category Context
│   ├── budget/                          # Budget Context
│   ├── reporting/                       # Reporting Context
│   ├── investment/                      # Investment Context
│   ├── goal/                            # Goal Context
│   └── notification/                   # Notification Context
│
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── prisma/
│   └── schema.prisma
│
└── docker-compose.yml
```

### 6.2. Análise Comparativa: PHP vs Node.js vs Go

#### 6.2.1. Visão Geral das Tecnologias

**PHP**
- **Maturidade**: 28+ anos, extremamente maduro
- **Ecossistema**: Laravel, Symfony, PHPUnit
- **Performance**: PHP 8.x com JIT é muito rápido
- **Comunidade**: Enorme, muito material disponível
- **Casos de Uso**: Web apps, APIs REST, CRUDs

**Node.js**
- **Maturidade**: 15+ anos, muito maduro
- **Ecossistema**: NestJS, Express, TypeScript
- **Performance**: Excelente para I/O, single-threaded
- **Comunidade**: Enorme, muito material
- **Casos de Uso**: APIs, real-time, microservices

**Go**
- **Maturidade**: 14+ anos, maduro e estável
- **Ecossistema**: Gin, Echo, GORM, ent
- **Performance**: Excepcional, compilado
- **Comunidade**: Crescendo rapidamente
- **Casos de Uso**: APIs de alta performance, microservices, sistemas distribuídos

#### 6.2.2. Comparação Detalhada

| Aspecto | PHP | Node.js | Go |
|---------|-----|---------|-----|
| **Performance (CPU)** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Performance (I/O)** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Concorrência** | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Curva de Aprendizado** | ⭐⭐⭐⭐⭐ (você já conhece) | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Ecossistema DDD** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Type Safety** | ⭐⭐⭐ (PHP 8+) | ⭐⭐⭐⭐⭐ (TypeScript) | ⭐⭐⭐⭐⭐ |
| **Produtividade** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Deploy/DevOps** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Comunidade** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Documentação** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

#### 6.2.3. Análise por Tecnologia

##### PHP (Laravel/Symfony)

**Vantagens:**
- ✅ **Você já domina**: Produtividade imediata
- ✅ **Ecossistema maduro**: Laravel tem tudo que precisa
- ✅ **ORM excelente**: Eloquent (Laravel) ou Doctrine (Symfony)
- ✅ **Performance PHP 8.x**: JIT compiler, muito rápido
- ✅ **Muitos pacotes**: Composer tem tudo
- ✅ **Documentação excelente**: Laravel docs são ótimas
- ✅ **Validação nativa**: Form Requests, Validators
- ✅ **Jobs/Queues**: Para processar transações recorrentes
- ✅ **Event System**: Laravel Events para Domain Events

**Desvantagens:**
- ❌ **DDD menos comum**: Menos exemplos/práticas DDD
- ❌ **Type safety**: PHP 8+ melhorou, mas não é TypeScript
- ❌ **Performance absoluta**: Ainda abaixo de Go
- ❌ **Concorrência**: Limitada (mas suficiente para este projeto)

**Stack Sugerida (PHP):**
```php
// Framework: Laravel 11 ou Symfony 6
// ORM: Eloquent (Laravel) ou Doctrine (Symfony)
// Validação: Form Requests (Laravel) ou Validator (Symfony)
// Event Bus: Laravel Events ou Symfony EventDispatcher
// Testes: PHPUnit
// API: Laravel Sanctum ou Symfony Security
```

**Exemplo DDD em PHP (Laravel):**
```php
// Domain/Entities/Transaction.php
class Transaction {
    private TransactionId $id;
    private Money $amount;
    private TransactionType $type;
    
    public function approve(): void {
        $this->status = TransactionStatus::APPROVED();
        event(new TransactionApproved($this->id, $this->amount));
    }
}

// Infrastructure/Persistence/EloquentTransactionRepository.php
class EloquentTransactionRepository implements TransactionRepository {
    public function save(Transaction $transaction): void {
        TransactionModel::updateOrCreate(
            ['id' => $transaction->getId()->value()],
            $transaction->toArray()
        );
    }
}
```

**Performance PHP 8.x:**
- JIT compiler ativa em PHP 8.0+
- Performance comparável a Node.js em muitos casos
- Opcache para cache de bytecode
- Para este projeto: **mais que suficiente**

##### Node.js (NestJS + TypeScript)

**Vantagens:**
- ✅ **Type Safety**: TypeScript é excelente
- ✅ **DDD Nativo**: NestJS foi feito pensando em DDD
- ✅ **Ecossistema moderno**: Prisma, TypeORM, etc.
- ✅ **Performance I/O**: Excelente para APIs
- ✅ **Frontend**: Mesma linguagem (TypeScript)
- ✅ **Real-time**: WebSockets nativos
- ✅ **Microservices**: Fácil escalar horizontalmente
- ✅ **Async/Await**: Código limpo e moderno

**Desvantagens:**
- ❌ **Você conhece menos**: Curva de aprendizado
- ❌ **Single-threaded**: CPU-bound pode ser limitante
- ❌ **Runtime overhead**: JavaScript tem overhead
- ❌ **Memory**: Pode consumir mais que Go

**Stack Sugerida (Node.js):**
```typescript
// Framework: NestJS
// ORM: Prisma ou TypeORM
// Validação: class-validator
// Event Bus: NestJS EventEmitter ou RabbitMQ
// Testes: Jest
// API: NestJS Guards + JWT
```

**Exemplo DDD em Node.js (NestJS):**
```typescript
// domain/entities/Transaction.ts
export class Transaction {
  private id: TransactionId;
  private amount: Money;
  
  approve(): void {
    this.status = TransactionStatus.APPROVED;
    this.addDomainEvent(new TransactionApproved(this.id, this.amount));
  }
}

// infrastructure/persistence/PrismaTransactionRepository.ts
@Injectable()
export class PrismaTransactionRepository implements ITransactionRepository {
  constructor(private prisma: PrismaService) {}
  
  async save(transaction: Transaction): Promise<void> {
    await this.prisma.transaction.upsert({
      where: { id: transaction.getId().value },
      create: transaction.toPersistence(),
      update: transaction.toPersistence()
    });
  }
}
```

**Performance Node.js:**
- Excelente para I/O (APIs, banco de dados)
- V8 engine muito otimizado
- Para este projeto: **excelente escolha**

##### Go

**Vantagens:**
- ✅ **Performance excepcional**: Compilado, muito rápido
- ✅ **Concorrência nativa**: Goroutines são incríveis
- ✅ **Baixo consumo de memória**: Eficiente
- ✅ **Type safety**: Forte e estático
- ✅ **Simplicidade**: Linguagem simples e direta
- ✅ **Deploy**: Binário único, fácil deploy
- ✅ **Escalabilidade**: Excelente para alta carga
- ✅ **Aprendizado**: Você quer aprender (motivação!)

**Desvantagens:**
- ❌ **Você não conhece**: Curva de aprendizado inicial
- ❌ **Ecossistema menor**: Menos pacotes que PHP/Node
- ❌ **DDD menos comum**: Menos exemplos/práticas
- ❌ **Produtividade inicial**: Mais lento no começo
- ❌ **ORM limitado**: GORM é bom, mas não é Prisma/Eloquent
- ❌ **Validação**: Mais manual que Laravel/NestJS
- ❌ **Frontend**: Precisa de stack separada (Next.js ainda)

**Stack Sugerida (Go):**
```go
// Framework: Gin ou Echo (ou Fiber para performance)
// ORM: GORM ou ent (do Facebook)
// Validação: go-playground/validator
// Event Bus: EventBus ou RabbitMQ
// Testes: testing package nativo
// API: JWT-Go para autenticação
```

**Exemplo DDD em Go:**
```go
// domain/transaction.go
type Transaction struct {
    id      TransactionID
    amount  Money
    status  TransactionStatus
    events  []DomainEvent
}

func (t *Transaction) Approve() error {
    t.status = TransactionStatusApproved
    t.events = append(t.events, NewTransactionApprovedEvent(t.id, t.amount))
    return nil
}

// infrastructure/persistence/gorm_repository.go
type GormTransactionRepository struct {
    db *gorm.DB
}

func (r *GormTransactionRepository) Save(tx *Transaction) error {
    model := tx.ToPersistence()
    return r.db.Save(&model).Error
}
```

**Performance Go:**
- **Muito superior** a PHP e Node.js
- Compilado nativamente
- Goroutines para concorrência
- Para este projeto: **overkill** (mas excelente para aprender)

#### 6.2.4. Análise Específica para o Projeto

**Requisitos do Projeto:**
- Gestão financeira pessoal/profissional
- CRUD de transações, contas, categorias
- Relatórios e análises
- Orçamentos e metas
- Dashboard com gráficos
- Não é alta escala (milhões de requests)

**Análise por Requisito:**

| Requisito | PHP | Node.js | Go |
|-----------|-----|---------|-----|
| **CRUD Básico** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Relatórios (Agregações)** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Real-time (Dashboard)** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Processamento Assíncrono** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Type Safety** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **DDD** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Produtividade** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |

#### 6.2.5. Recomendações por Cenário

##### Cenário 1: Produtividade e Entrega Rápida
**Recomendação: PHP (Laravel)**
- Você já domina
- MVP em menos tempo
- Ecossistema completo
- Performance PHP 8.x é suficiente

**Tempo estimado MVP:**
- PHP: 2-3 semanas
- Node.js: 3-4 semanas (aprendizado)
- Go: 4-6 semanas (aprendizado + menos tooling)

##### Cenário 2: Aprendizado e Performance
**Recomendação: Go**
- Você quer aprender (motivação!)
- Performance excepcional
- Habilidade valiosa no mercado
- Projeto é bom para aprender (não é crítico)

**Considerações:**
- Mais tempo inicial
- Menos exemplos DDD em Go
- Mas aprendizado valioso

##### Cenário 3: Equilíbrio e Modernidade
**Recomendação: Node.js (NestJS)**
- TypeScript é excelente
- DDD nativo no NestJS
- Ecossistema moderno
- Você conhece um pouco (base JavaScript)

**Considerações:**
- Curva de aprendizado média
- Mas stack moderna e valorizada

#### 6.2.6. Análise de Performance Real

**Benchmarks Típicos (Requests/segundo):**

```
API Simples (CRUD):
- PHP 8.2 (Laravel): ~8.000 req/s
- Node.js (NestJS): ~12.000 req/s
- Go (Gin): ~50.000+ req/s

API com Banco de Dados:
- PHP 8.2 (Laravel): ~3.000 req/s
- Node.js (NestJS): ~4.000 req/s
- Go (Gin + GORM): ~15.000+ req/s

Agregações Complexas (Relatórios):
- PHP 8.2: ~1.500 req/s
- Node.js: ~2.000 req/s
- Go: ~8.000+ req/s
```

**Para este projeto:**
- **PHP**: Mais que suficiente (centenas de usuários)
- **Node.js**: Excelente (milhares de usuários)
- **Go**: Overkill, mas excelente (milhões de usuários)

#### 6.2.7. Considerações de Aprendizado

##### Aprendendo Go
**Vantagens:**
- ✅ Linguagem simples (menos features = menos confusão)
- ✅ Performance excepcional
- ✅ Muito valorizada no mercado
- ✅ Excelente para sistemas distribuídos
- ✅ Concorrência nativa (goroutines)

**Desafios:**
- ❌ Menos material DDD específico
- ❌ Ecossistema menor
- ❌ Mais verboso que PHP/Node
- ❌ Sem generics avançados (até Go 1.18+)
- ❌ Error handling explícito (pode ser verboso)

**Tempo estimado para ficar produtivo:**
- Básico: 1-2 semanas
- Intermediário: 1-2 meses
- Avançado: 3-6 meses

##### Aprendendo Node.js/TypeScript
**Vantagens:**
- ✅ TypeScript é poderoso
- ✅ Muito material disponível
- ✅ Ecossistema enorme
- ✅ Frontend + Backend mesma linguagem

**Desafios:**
- ❌ TypeScript tem curva (generics, tipos avançados)
- ❌ Async/await patterns
- ❌ Node.js event loop (precisa entender)

**Tempo estimado para ficar produtivo:**
- Básico: 1 semana
- Intermediário: 2-3 semanas
- Avançado: 1-2 meses

#### 6.2.8. Recomendação Final

**Minha Recomendação: Go 🚀**

**Por quê?**
1. **Aprendizado valioso**: Você quer aprender, e Go é excelente
2. **Performance**: Mesmo sendo overkill, é um diferencial
3. **Mercado**: Go está em alta, habilidade valiosa
4. **Projeto ideal**: Não é crítico, pode aprender sem pressão
5. **Simplicidade**: Go é simples, curva de aprendizado razoável
6. **Futuro**: Conhecimento que você vai usar em outros projetos

**Plano Sugerido:**
1. **Semana 1-2**: Aprenda Go básico (tour, tutoriais)
2. **Semana 3-4**: Implemente Shared Kernel + Identity Context
3. **Semana 5-6**: Implemente Account + Transaction (Core Domain)
4. **Semana 7+**: Continue evoluindo

**Stack Go Recomendada:**
```go
// Framework: Gin (simples) ou Echo (mais features)
// ORM: GORM (mais popular) ou ent (type-safe)
// Validação: go-playground/validator
// Event Bus: EventBus ou RabbitMQ
// Testes: testing package + testify
// Migrations: golang-migrate
```

**Alternativa (se precisar entregar rápido):**
- **PHP Laravel**: Se o prazo for apertado, use o que você conhece
- Depois refatore partes críticas para Go (se necessário)

#### 6.2.9. Comparação de Código (Mesma Funcionalidade)

**Criar uma Transação - Comparação:**

**PHP (Laravel):**
```php
// app/Domain/UseCases/CreateTransactionUseCase.php
class CreateTransactionUseCase {
    public function execute(CreateTransactionDTO $dto): Transaction {
        $account = $this->accountRepo->findById($dto->accountId);
        $transaction = Transaction::create(
            $account,
            Money::create($dto->amount, Currency::BRL),
            TransactionType::fromString($dto->type)
        );
        $this->transactionRepo->save($transaction);
        return $transaction;
    }
}
```

**Node.js (NestJS):**
```typescript
// application/use-cases/CreateTransactionUseCase.ts
@Injectable()
export class CreateTransactionUseCase {
  async execute(dto: CreateTransactionDTO): Promise<Transaction> {
    const account = await this.accountRepo.findById(dto.accountId);
    const transaction = Transaction.create(
      account,
      Money.create(dto.amount, Currency.BRL),
      TransactionType.fromString(dto.type)
    );
    await this.transactionRepo.save(transaction);
    return transaction;
  }
}
```

**Go:**
```go
// application/usecases/create_transaction.go
type CreateTransactionUseCase struct {
    accountRepo     domain.AccountRepository
    transactionRepo domain.TransactionRepository
}

func (uc *CreateTransactionUseCase) Execute(dto CreateTransactionDTO) (*domain.Transaction, error) {
    account, err := uc.accountRepo.FindByID(dto.AccountID)
    if err != nil {
        return nil, err
    }
    
    transaction := domain.NewTransaction(
        account,
        domain.NewMoney(dto.Amount, domain.CurrencyBRL),
        domain.TransactionTypeFromString(dto.Type),
    )
    
    if err := uc.transactionRepo.Save(transaction); err != nil {
        return nil, err
    }
    
    return transaction, nil
}
```

**Observações:**
- **PHP**: Mais conciso, menos verboso
- **Node.js**: Similar ao PHP, async/await limpo
- **Go**: Mais verboso (error handling explícito), mas mais explícito

### 6.3. Stack Tecnológico (Atualizado)

#### Backend - Opção 1: Go (Recomendado para Aprendizado)
- **Linguagem**: Go 1.21+
- **Framework**: Gin ou Echo
- **ORM**: GORM ou ent
- **Validação**: go-playground/validator
- **Event Bus**: EventBus ou RabbitMQ
- **Testes**: testing package + testify
- **Migrations**: golang-migrate

#### Backend - Opção 2: PHP (Recomendado para Produtividade)
- **Linguagem**: PHP 8.2+
- **Framework**: Laravel 11
- **ORM**: Eloquent
- **Validação**: Form Requests
- **Event Bus**: Laravel Events
- **Testes**: PHPUnit
- **API**: Laravel Sanctum

#### Backend - Opção 3: Node.js (Recomendado para Modernidade)
- **Runtime**: Node.js 20+
- **Framework**: NestJS
- **Linguagem**: TypeScript
- **ORM**: Prisma
- **Validação**: class-validator
- **Event Bus**: NestJS EventEmitter ou RabbitMQ
- **Testes**: Jest

#### Frontend
- **Framework**: Next.js (React)
- **UI Library**: Tailwind CSS + shadcn/ui
- **State Management**: Zustand ou React Query
- **Gráficos**: Recharts
- **Formulários**: React Hook Form + Zod

#### DevOps
- **Containerização**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Deploy**: Railway, Render ou AWS

### 6.3. Análise de ORMs Disponíveis

#### 6.3.1. ORMs Principais para Node.js/TypeScript

**1. Prisma**
- **Tipo**: Query Builder + ORM híbrido
- **Abordagem**: Schema-first (define schema, gera tipos)
- **Características**:
  - Type-safe por padrão
  - Migrations automáticas
  - Prisma Studio (GUI para dados)
  - Excelente DX (Developer Experience)
  - Suporte a múltiplos bancos (PostgreSQL, MySQL, SQLite, MongoDB, etc.)
  - Query builder poderoso e intuitivo
  - Relações bem definidas

**2. TypeORM**
- **Tipo**: ORM completo
- **Abstração**: Active Record + Data Mapper
- **Características**:
  - Decorators para definir entidades
  - Suporte a múltiplos padrões (Active Record, Data Mapper, Repository)
  - Migrations manuais
  - Type-safe com TypeScript
  - Suporte a múltiplos bancos
  - Mais flexível, mas mais complexo
  - Comunidade grande e madura

**3. Sequelize**
- **Tipo**: ORM tradicional
- **Abordagem**: Promise-based
- **Características**:
  - Muito maduro e estável
  - Suporte a múltiplos bancos
  - Migrations incluídas
  - TypeScript suportado (mas não nativo)
  - Comunidade extensa
  - Mais verboso que alternativas modernas

**4. Drizzle ORM**
- **Tipo**: ORM leve e type-safe
- **Abstração**: SQL-like, minimalista
- **Características**:
  - Extremamente leve
  - Type-safe nativo
  - SQL-like syntax
  - Sem runtime overhead
  - Migrations manuais
  - Mais controle, menos "magia"

**5. MikroORM**
- **Tipo**: ORM completo
- **Abstração**: Data Mapper
- **Características**:
  - Unit of Work pattern
  - Identity Map
  - Type-safe
  - Suporte a múltiplos bancos
  - Boa para DDD (Identity Map ajuda)
  - Migrations automáticas

**6. Kysely**
- **Tipo**: Query Builder type-safe
- **Abstração**: SQL puro, type-safe
- **Características**:
  - Type-safe SQL builder
  - Sem abstrações pesadas
  - Controle total sobre queries
  - Excelente para queries complexas
  - Não é um ORM completo

#### 6.3.2. Comparação Rápida

| ORM | Type-Safety | DX | Performance | DDD Friendly | Curva Aprendizado |
|-----|-------------|-----|-------------|--------------|-------------------|
| **Prisma** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **TypeORM** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Sequelize** | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Drizzle** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **MikroORM** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Kysely** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ |

#### 6.3.3. Por que Prisma foi Escolhido?

**1. Type-Safety Excepcional**
```typescript
// Prisma gera tipos automaticamente do schema
const user = await prisma.user.findUnique({
  where: { id: userId }
});
// user é totalmente type-safe, com autocomplete completo
```

**2. Developer Experience Superior**
- Schema declarativo e claro
- Prisma Studio para visualizar dados
- Migrations automáticas e seguras
- Autocomplete excelente no IDE
- Mensagens de erro claras

**3. Compatibilidade com DDD**
- **Repositórios**: Prisma facilita implementação de repositórios
  ```typescript
  class PrismaAccountRepository implements IAccountRepository {
    async findById(id: AccountId): Promise<Account | null> {
      const data = await prisma.account.findUnique({
        where: { id: id.value }
      });
      return data ? Account.fromPersistence(data) : null;
    }
  }
  ```
- **Mapeamento**: Fácil mapear entre modelos de domínio e persistência
- **Transações**: Suporte nativo a transações
- **Queries Complexas**: Query builder poderoso para relatórios

**4. Performance**
- Queries otimizadas
- Connection pooling nativo
- Suporte a índices e constraints
- Preparação de statements

**5. Manutenibilidade**
- Schema como fonte da verdade
- Migrations versionadas
- Fácil refatoração
- Documentação excelente

**6. Ecossistema**
- Integração com NestJS
- Prisma Client extensível
- Middleware para logging, validação, etc.
- Comunidade ativa

#### 6.3.2. ORMs Principais para Go

**1. GORM**
- **Tipo**: ORM completo
- **Abordagem**: Struct-based
- **Características**:
  - Mais popular ORM em Go
  - Migrations automáticas
  - Hooks (BeforeCreate, AfterUpdate, etc.)
  - Associations (has many, belongs to, etc.)
  - Query builder intuitivo
  - Suporte a múltiplos bancos
  - Comunidade grande
  - Documentação boa

**2. ent (Facebook)**
- **Tipo**: Entity framework type-safe
- **Abstração**: Code generation
- **Características**:
  - Type-safe por geração de código
  - Schema-first (define schema, gera código)
  - Similar ao Prisma em conceito
  - Excelente type-safety
  - GraphQL integrado (opcional)
  - Migrations automáticas
  - Mais moderno que GORM

**3. SQLBoiler**
- **Tipo**: Code generator
- **Abordagem**: Database-first
- **Características**:
  - Gera código do schema do banco
  - Type-safe
  - Performance excelente (sem reflection)
  - SQL puro quando necessário
  - Menos "magia", mais controle

**4. sqlx**
- **Tipo**: Extensão do database/sql
- **Abordagem**: SQL direto
- **Características**:
  - Não é um ORM completo
  - Type-safe scanning
  - SQL direto (controle total)
  - Leve e performático
  - Para quem prefere SQL puro

**Comparação Go ORMs:**

| ORM | Type-Safety | DX | Performance | DDD Friendly | Popularidade |
|-----|-------------|-----|-------------|--------------|--------------|
| **GORM** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **ent** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **SQLBoiler** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **sqlx** | ⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ |

**Recomendação para Go:**
- **GORM**: Se quer produtividade e facilidade (similar ao Eloquent)
- **ent**: Se quer type-safety máximo (similar ao Prisma)
- **SQLBoiler**: Se quer performance e controle
- **sqlx**: Se prefere SQL direto

#### 6.3.3. ORMs Principais para PHP

**1. Eloquent (Laravel)**
- **Tipo**: Active Record ORM
- **Abordagem**: Model-based
- **Características**:
  - ORM padrão do Laravel
  - Active Record pattern
  - Migrations integradas
  - Relationships fáceis
  - Query builder poderoso
  - Mutators/Accessors
  - Events e Observers
  - Muito produtivo

**2. Doctrine (Symfony)**
- **Tipo**: Data Mapper ORM
- **Abordagem**: Entity-based
- **Características**:
  - Data Mapper (mais DDD-friendly)
  - Unit of Work pattern
  - Identity Map
  - Migrations automáticas
  - DQL (Doctrine Query Language)
  - Mais complexo, mas mais poderoso
  - Melhor para DDD

**3. Propel**
- **Tipo**: Active Record
- **Abordagem**: Code generation
- **Características**:
  - Gera código do schema
  - Type-safe
  - Performance boa
  - Menos popular

**Comparação PHP ORMs:**

| ORM | Type-Safety | DX | Performance | DDD Friendly | Popularidade |
|-----|-------------|-----|-------------|--------------|--------------|
| **Eloquent** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Doctrine** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Propel** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |

**Recomendação para PHP:**
- **Eloquent**: Se usa Laravel (mais produtivo)
- **Doctrine**: Se usa Symfony ou quer DDD puro (Data Mapper)

#### 6.3.5. Alternativas e Quando Considerar (Node.js)

**TypeORM - Quando usar:**
- Projeto já usa TypeORM
- Precisa de Active Record pattern
- Quer máxima flexibilidade
- Time já conhece bem

**Drizzle ORM - Quando usar:**
- Performance crítica
- Quer controle total sobre SQL
- Projeto pequeno/médio
- Prefere SQL-like syntax

**MikroORM - Quando usar:**
- DDD puro com Identity Map
- Unit of Work é importante
- Projeto complexo com muitos relacionamentos

**Kysely - Quando usar:**
- Queries muito complexas
- Não precisa de ORM completo
- Type-safety é prioridade máxima
- Performance é crítica

#### 6.3.6. Exemplos de Uso com DDD

##### Exemplo 1: Prisma (Node.js)

**Schema Prisma (schema.prisma):**
```prisma
model Account {
  id        String   @id @default(uuid())
  userId    String
  name      String
  type      String   // BANK, WALLET, INVESTMENT
  balance   Decimal  @default(0)
  context   String   // PERSONAL, PROFESSIONAL
  isActive  Boolean  @default(true)
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  
  transactions Transaction[]
  
  @@index([userId])
  @@index([userId, context])
}

model Transaction {
  id          String   @id @default(uuid())
  userId      String
  accountId   String
  categoryId  String
  type        String   // INCOME, EXPENSE, TRANSFER
  amount      Decimal
  description String?
  date        DateTime
  status      String   @default("PENDING") // PENDING, APPROVED, CANCELLED
  context     String
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt
  
  account  Account  @relation(fields: [accountId], references: [id])
  category Category @relation(fields: [categoryId], references: [id])
  
  @@index([userId, date])
  @@index([accountId])
  @@index([categoryId])
}
```

**Implementação de Repositório:**
```typescript
// Domain Layer - Interface
interface IAccountRepository {
  findById(id: AccountId): Promise<Account | null>;
  save(account: Account): Promise<void>;
}

// Infrastructure Layer - Implementação
class PrismaAccountRepository implements IAccountRepository {
  constructor(private prisma: PrismaClient) {}
  
  async findById(id: AccountId): Promise<Account | null> {
    const data = await this.prisma.account.findUnique({
      where: { id: id.value }
    });
    
    if (!data) return null;
    
    // Mapear de persistência para domínio
    return Account.fromPersistence({
      id: AccountId.create(data.id),
      userId: UserId.create(data.userId),
      name: AccountName.create(data.name),
      type: AccountType.fromString(data.type),
      balance: Money.create(Number(data.balance), Currency.BRL),
      context: AccountContext.fromString(data.context),
      isActive: data.isActive,
      createdAt: data.createdAt,
      updatedAt: data.updatedAt
    });
  }
  
  async save(account: Account): Promise<void> {
    const data = account.toPersistence(); // Método no agregado
    
    await this.prisma.account.upsert({
      where: { id: data.id },
      create: data,
      update: data
    });
  }
}
```

**Vantagens desta Abordagem:**
- Domínio permanece puro (sem dependências)
- Prisma fica isolado na camada de infraestrutura
- Fácil trocar ORM no futuro (basta implementar interface)
- Type-safety em todas as camadas
- Testes fáceis (mock do repositório)

##### Exemplo 2: GORM (Go)

**Model GORM:**
```go
// infrastructure/persistence/models/account.go
type AccountModel struct {
    ID        string    `gorm:"type:uuid;primary_key"`
    UserID    string    `gorm:"type:uuid;index"`
    Name      string    `gorm:"type:varchar(255)"`
    Type      string    `gorm:"type:varchar(50)"`
    Balance   float64   `gorm:"type:decimal(15,2)"`
    Context   string    `gorm:"type:varchar(20)"`
    IsActive  bool      `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (AccountModel) TableName() string {
    return "accounts"
}
```

**Implementação de Repositório:**
```go
// Domain Layer - Interface
type AccountRepository interface {
    FindByID(id AccountID) (*Account, error)
    Save(account *Account) error
}

// Infrastructure Layer - Implementação
type GormAccountRepository struct {
    db *gorm.DB
}

func (r *GormAccountRepository) FindByID(id AccountID) (*Account, error) {
    var model AccountModel
    if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    
    // Mapear de persistência para domínio
    return AccountFromPersistence(model), nil
}

func (r *GormAccountRepository) Save(account *Account) error {
    model := account.ToPersistence()
    return r.db.Save(&model).Error
}
```

**Vantagens:**
- Simples e direto
- Migrations automáticas
- Hooks úteis (BeforeCreate, AfterUpdate)
- Comunidade grande

##### Exemplo 3: Eloquent (PHP/Laravel)

**Model Eloquent:**
```php
// app/Infrastructure/Persistence/Models/AccountModel.php
class AccountModel extends Model
{
    protected $table = 'accounts';
    
    protected $fillable = [
        'id', 'user_id', 'name', 'type', 
        'balance', 'context', 'is_active'
    ];
    
    protected $casts = [
        'balance' => 'decimal:2',
        'is_active' => 'boolean',
        'created_at' => 'datetime',
        'updated_at' => 'datetime',
    ];
}
```

**Implementação de Repositório:**
```php
// Domain Layer - Interface
interface AccountRepository
{
    public function findById(AccountId $id): ?Account;
    public function save(Account $account): void;
}

// Infrastructure Layer - Implementação
class EloquentAccountRepository implements AccountRepository
{
    public function findById(AccountId $id): ?Account
    {
        $model = AccountModel::find($id->value());
        
        if (!$model) {
            return null;
        }
        
        // Mapear de persistência para domínio
        return Account::fromPersistence([
            'id' => AccountId::create($model->id),
            'userId' => UserId::create($model->user_id),
            'name' => AccountName::create($model->name),
            'type' => AccountType::fromString($model->type),
            'balance' => Money::create($model->balance, Currency::BRL),
            'context' => AccountContext::fromString($model->context),
            'isActive' => $model->is_active,
            'createdAt' => $model->created_at,
            'updatedAt' => $model->updated_at,
        ]);
    }
    
    public function save(Account $account): void
    {
        $data = $account->toPersistence();
        
        AccountModel::updateOrCreate(
            ['id' => $data['id']],
            $data
        );
    }
}
```

**Vantagens:**
- Muito produtivo
- Migrations integradas
- Relationships fáceis
- Query builder poderoso

#### 6.3.7. Considerações Finais sobre ORM

**Resumo por Tecnologia:**

**Node.js - Prisma:**
1. ✅ Type-safety nativo e excelente
2. ✅ DX superior acelera desenvolvimento
3. ✅ Compatível com DDD (repositórios fáceis)
4. ✅ Performance adequada
5. ✅ Migrations seguras e versionadas

**Go - GORM:**
1. ✅ Simples e direto
2. ✅ Migrations automáticas
3. ✅ Hooks úteis
4. ✅ Comunidade grande
5. ⚠️ Menos type-safe que ent

**PHP - Eloquent:**
1. ✅ Muito produtivo
2. ✅ Migrations integradas
3. ✅ Relationships fáceis
4. ✅ Query builder poderoso
5. ⚠️ Active Record (menos DDD puro)

**Para este projeto específico:**
- **Node.js + Prisma**: Excelente para DDD e type-safety
- **Go + GORM**: Boa escolha, simples e funcional
- **PHP + Eloquent**: Muito produtivo, mas menos DDD puro

**Prisma é uma excelente escolha porque:**
1. ✅ Type-safety nativo e excelente
2. ✅ DX superior acelera desenvolvimento
3. ✅ Compatível com DDD (repositórios fáceis)
4. ✅ Performance adequada
5. ✅ Migrations seguras e versionadas
6. ✅ Comunidade ativa e documentação excelente
7. ✅ Integração perfeita com NestJS

**Limitações do Prisma:**
- Menos flexível que TypeORM para casos muito complexos
- Schema-first pode ser limitante em alguns cenários
- Não é um ORM "puro" (mais próximo de query builder)

**Para este projeto específico:**
- Prisma atende perfeitamente às necessidades
- Facilita implementação de repositórios DDD
- Acelera desenvolvimento do MVP
- Type-safety previne muitos bugs
- Migrations facilitam evolução do schema

## 7. Padrões e Práticas DDD

### 7.1. Agregados

**Regras de Agregados:**
- Cada agregado tem uma raiz (Aggregate Root)
- Acesso a entidades internas apenas através da raiz
- Transações atômicas por agregado
- IDs únicos e imutáveis
- Validação de invariantes na raiz

**Exemplo - Transaction Agregado:**
```typescript
// Transaction é a raiz do agregado
class Transaction {
  // Entidades internas (não acessíveis diretamente)
  private installments: TransactionInstallment[];
  
  // Acesso controlado
  addInstallment(installment: TransactionInstallment): void {
    // Validações e regras de negócio
    this.installments.push(installment);
  }
}
```

### 7.2. Repositórios

**Características:**
- Abstração de persistência
- Interface no domínio, implementação na infraestrutura
- Retorna agregados completos
- Métodos expressivos do domínio

### 7.3. Domain Events

**Uso:**
- Comunicação entre agregados
- Integração entre bounded contexts
- Auditoria e rastreabilidade
- Desacoplamento

**Exemplo:**
```typescript
class Transaction {
  approve(): void {
    this.status = TransactionStatus.APPROVED;
    this.addDomainEvent(new TransactionApproved(this.id, this.amount));
  }
}
```

### 7.4. Value Objects

**Características:**
- Imutáveis
- Comparação por valor
- Validação no construtor
- Sem identidade

### 7.5. Domain Services

**Quando usar:**
- Lógica que não pertence a uma entidade específica
- Operações que envolvem múltiplos agregados
- Cálculos complexos de domínio

## 8. Fases de Desenvolvimento (DDD)

### Fase 1: Fundação e Core Domain
**Objetivo**: Estabelecer a base e os contextos mais críticos

1. **Shared Kernel**
   - Money, Currency, AccountContext
   - Event Bus básico

2. **Identity Context**
   - User agregado
   - Autenticação básica
   - Repositório e persistência

3. **Account Management Context**
   - Account agregado
   - Operações básicas (CRUD)
   - Integração com Identity

4. **Transaction Context (Core Domain)**
   - Transaction agregado
   - Processamento básico
   - Integração com Account
   - Eventos de domínio

**Entregáveis:**
- Usuário pode criar contas
- Usuário pode registrar transações
- Saldo atualizado automaticamente

### Fase 2: Expansão do Domínio
**Objetivo**: Adicionar funcionalidades essenciais

1. **Category Context**
   - Category agregado
   - Hierarquia de categorias
   - Categorias padrão

2. **Transaction Context (Expansão)**
   - RecurringTransaction
   - Transferências
   - Parcelamento

3. **Budget Context**
   - Budget agregado
   - Cálculo de uso
   - Alertas básicos

**Entregáveis:**
- Categorização de transações
- Transações recorrentes
- Orçamentos funcionais

### Fase 3: Análise e Relatórios
**Objetivo**: Fornecer insights financeiros

1. **Reporting Context**
   - Agregação de dados
   - Relatórios básicos
   - Dashboard

2. **Integração Reporting ↔ Transaction**
   - Event handlers
   - Cache de agregados

**Entregáveis:**
- Dashboard funcional
- Relatórios básicos
- Gráficos e visualizações

### Fase 4: Funcionalidades Avançadas
**Objetivo**: Completar o sistema

1. **Investment Context**
   - Investment agregado
   - Cálculo de retorno

2. **Goal Context**
   - Goal agregado
   - Acompanhamento de progresso

3. **Notification Context**
   - Sistema de notificações
   - Integração com outros contextos

**Entregáveis:**
- Gestão de investimentos
- Metas financeiras
- Sistema de notificações completo

### Fase 5: Refinamento e Otimização
**Objetivo**: Melhorar qualidade e performance

1. Otimização de queries
2. Cache estratégico
3. Testes completos
4. Documentação
5. Performance tuning

## 9. Testes em DDD

### 9.1. Testes de Domínio

**Testes Unitários de Entidades:**
```typescript
describe('Transaction', () => {
  it('should not allow negative amounts', () => {
    expect(() => {
      new Transaction(/* ... */, new Money(-100, Currency.BRL));
    }).toThrow();
  });
  
  it('should update account balance when approved', () => {
    // Teste de comportamento
  });
});
```

**Testes de Value Objects:**
```typescript
describe('Money', () => {
  it('should add two money objects with same currency', () => {
    const money1 = new Money(100, Currency.BRL);
    const money2 = new Money(50, Currency.BRL);
    const result = money1.add(money2);
    expect(result.toNumber()).toBe(150);
  });
});
```

### 9.2. Testes de Serviços de Domínio

```typescript
describe('TransactionProcessingService', () => {
  it('should process transaction and update account balance', () => {
    // Teste de serviço de domínio
  });
});
```

### 9.3. Testes de Integração

- Testes de repositórios
- Testes de casos de uso
- Testes de integração entre contextos

### 9.4. Testes E2E

- Fluxos completos do usuário
- Integração frontend-backend

## 10. Métricas e Monitoramento

### 10.1. Métricas de Domínio

- Transações processadas por dia
- Valor total movimentado
- Número de usuários ativos
- Taxa de uso de orçamentos
- Metas completadas

### 10.2. Métricas Técnicas

- Tempo de resposta de queries
- Taxa de erro
- Uptime
- Cobertura de testes

## 11. Próximos Passos

1. **Setup do Projeto**
   - Inicializar NestJS
   - Configurar Prisma
   - Estrutura de pastas DDD

2. **Shared Kernel**
   - Implementar Money, Currency
   - Event Bus básico

3. **Identity Context**
   - Implementar User agregado
   - Casos de uso de autenticação

4. **Account Management Context**
   - Implementar Account agregado
   - Casos de uso básicos

5. **Transaction Context (Core)**
   - Implementar Transaction agregado
   - Processamento de transações
   - Integração com Account

6. **Testes e Validação**
   - Testes unitários
   - Testes de integração
   - Validação de regras de negócio

## 12. Considerações Finais

### 12.1. Vantagens da Abordagem DDD

- **Modelagem Rica**: Domínio expressivo e compreensível
- **Manutenibilidade**: Código organizado e fácil de evoluir
- **Testabilidade**: Domínio isolado e testável
- **Escalabilidade**: Contextos independentes e escaláveis
- **Colaboração**: Linguagem ubíqua facilita comunicação

### 12.2. Desafios

- **Complexidade Inicial**: Curva de aprendizado
- **Over-engineering**: Evitar complexidade desnecessária
- **Performance**: Cuidado com múltiplas camadas
- **Time**: Desenvolvimento pode ser mais lento inicialmente

### 12.3. Boas Práticas

- Começar simples, evoluir gradualmente
- Focar no Core Domain primeiro
- Manter domínio puro (sem dependências de infra)
- Usar eventos para integração
- Documentar decisões de design

# Planejamento DDD - Sistema de Gestão Financeira (PHP)

## 1. Visão Geral

Sistema de gestão financeira pessoal e profissional desenvolvido em **PHP 8.2+** usando **Laravel** ou **Symfony**, aproveitando a produtividade, ecossistema maduro e performance do PHP 8.x com JIT.

## 2. Objetivos

- Controle total de finanças pessoais e profissionais
- Separação clara entre contas pessoais e profissionais
- Análise e relatórios financeiros
- Planejamento orçamentário
- Acompanhamento de metas financeiras
- Arquitetura escalável e manutenível
- Aproveitamento máximo da produtividade do Laravel/Symfony

## 3. Stack Tecnológico PHP

### 3.1. Tecnologias Principais

**Opção 1: Laravel (Recomendado para Produtividade)**
- **Framework**: Laravel 11+
- **Linguagem**: PHP 8.2+
- **ORM**: Eloquent
- **Validação**: Form Requests
- **Autenticação**: Laravel Sanctum
- **Event Bus**: Laravel Events
- **Testes**: PHPUnit
- **Migrations**: Laravel Migrations
- **Queue**: Laravel Queue (Redis/Database)
- **Cache**: Redis/Memcached
- **API Docs**: Laravel API Resources + Swagger

**Opção 2: Symfony (Recomendado para DDD Puro)**
- **Framework**: Symfony 6+
- **Linguagem**: PHP 8.2+
- **ORM**: Doctrine (Data Mapper)
- **Validação**: Symfony Validator
- **Autenticação**: Symfony Security
- **Event Bus**: Symfony EventDispatcher
- **Testes**: PHPUnit
- **Migrations**: Doctrine Migrations
- **Queue**: Symfony Messenger
- **Cache**: Symfony Cache

### 3.2. Por que PHP?

**Vantagens:**
- ✅ **Você já domina** - Produtividade imediata
- ✅ **Ecossistema maduro** - Laravel/Symfony têm tudo
- ✅ **ORM excelente** - Eloquent (Laravel) ou Doctrine (Symfony)
- ✅ **Performance PHP 8.x** - JIT compiler, muito rápido
- ✅ **Muitos pacotes** - Composer tem tudo
- ✅ **Documentação excelente** - Laravel docs são ótimas
- ✅ **Validação nativa** - Form Requests, Validators
- ✅ **Jobs/Queues** - Para processar transações recorrentes
- ✅ **Event System** - Laravel Events/Symfony EventDispatcher

**Desafios:**
- ⚠️ **DDD menos comum** - Menos exemplos/práticas DDD
- ⚠️ **Type safety** - PHP 8+ melhorou, mas não é TypeScript
- ⚠️ **Performance absoluta** - Ainda abaixo de Go

## 4. Arquitetura DDD em PHP

### 4.1. Estrutura em Camadas

```
┌─────────────────────────────────────┐
│     Controllers (Presentation)      │  (Controllers, Requests, Resources)
├─────────────────────────────────────┤
│     Use Cases (Application)         │  (Services, Actions)
├─────────────────────────────────────┤
│     Domain Layer                    │  (Entities, Value Objects, Domain Services)
├─────────────────────────────────────┤
│     Repositories (Infrastructure)  │  (Eloquent, Doctrine, External Services)
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

## 5. Estrutura de Pastas (Laravel DDD)

```
gestao-financeira-laravel/
├── app/
│   ├── Shared/                          # Shared Kernel
│   │   ├── Domain/
│   │   │   ├── ValueObjects/
│   │   │   │   ├── Money.php
│   │   │   │   ├── Currency.php
│   │   │   │   └── AccountContext.php
│   │   │   └── Events/
│   │   │       └── DomainEvent.php
│   │   └── Infrastructure/
│   │       └── EventBus/
│   │           └── EventBusService.php
│   │
│   ├── Identity/                         # Identity Context
│   │   ├── Domain/
│   │   │   ├── Entities/
│   │   │   │   └── User.php
│   │   │   ├── ValueObjects/
│   │   │   │   ├── Email.php
│   │   │   │   ├── PasswordHash.php
│   │   │   │   └── UserName.php
│   │   │   ├── Services/
│   │   │   │   ├── PasswordService.php
│   │   │   │   └── TokenService.php
│   │   │   ├── Repositories/
│   │   │   │   └── UserRepositoryInterface.php
│   │   │   └── Events/
│   │   │       ├── UserRegistered.php
│   │   │       └── UserPasswordChanged.php
│   │   ├── Application/
│   │   │   ├── UseCases/
│   │   │   │   ├── RegisterUserUseCase.php
│   │   │   │   ├── LoginUseCase.php
│   │   │   │   └── UpdateProfileUseCase.php
│   │   │   └── DTOs/
│   │   │       └── RegisterUserDTO.php
│   │   ├── Infrastructure/
│   │   │   ├── Persistence/
│   │   │   │   └── EloquentUserRepository.php
│   │   │   ├── Services/
│   │   │   │   └── JwtTokenService.php
│   │   │   └── Events/
│   │   │       └── DomainEventPublisher.php
│   │   └── Presentation/
│   │       ├── Controllers/
│   │       │   └── AuthController.php
│   │       ├── Requests/
│   │       │   └── RegisterUserRequest.php
│   │       └── Resources/
│   │           └── UserResource.php
│   │
│   ├── AccountManagement/                # Account Management Context
│   ├── Transaction/                       # Transaction Context
│   ├── Category/                         # Category Context
│   ├── Budget/                           # Budget Context
│   ├── Reporting/                        # Reporting Context
│   ├── Investment/                      # Investment Context
│   ├── Goal/                             # Goal Context
│   └── Notification/                    # Notification Context
│
├── database/
│   ├── migrations/
│   └── seeders/
│
├── tests/
│   ├── Unit/
│   ├── Feature/
│   └── Integration/
│
├── composer.json
├── phpunit.xml
├── Dockerfile
└── docker-compose.yml
```

## 6. Detalhamento dos Bounded Contexts (Laravel)

### 6.1. Identity Context

#### 6.1.1. Entidades (PHP)

**User (Agregado Raiz)**
```php
<?php
// app/Identity/Domain/Entities/User.php

namespace App\Identity\Domain\Entities;

use App\Identity\Domain\ValueObjects\Email;
use App\Identity\Domain\ValueObjects\PasswordHash;
use App\Identity\Domain\ValueObjects\UserName;
use App\Identity\Domain\ValueObjects\UserProfile;
use App\Identity\Domain\Events\UserRegistered;
use App\Identity\Domain\Events\UserPasswordChanged;
use App\Shared\Domain\Events\DomainEvent;

class User
{
    private string $id;
    private Email $email;
    private PasswordHash $passwordHash;
    private UserName $name;
    private UserProfile $profile;
    private \DateTimeImmutable $createdAt;
    private \DateTimeImmutable $updatedAt;
    private bool $isActive;
    
    /** @var DomainEvent[] */
    private array $events = [];
    
    private function __construct(
        string $id,
        Email $email,
        PasswordHash $passwordHash,
        UserName $name
    ) {
        $this->id = $id;
        $this->email = $email;
        $this->passwordHash = $passwordHash;
        $this->name = $name;
        $this->profile = UserProfile::default();
        $this->createdAt = new \DateTimeImmutable();
        $this->updatedAt = new \DateTimeImmutable();
        $this->isActive = true;
    }
    
    public static function create(
        Email $email,
        PasswordHash $passwordHash,
        UserName $name
    ): self {
        $user = new self(
            \Ramsey\Uuid\Uuid::uuid4()->toString(),
            $email,
            $passwordHash,
            $name
        );
        
        $user->addEvent(new UserRegistered($user->id));
        return $user;
    }
    
    public function changePassword(string $oldPassword, string $newPassword): void
    {
        if (!$this->passwordHash->verify($oldPassword)) {
            throw new \InvalidArgumentException('Invalid old password');
        }
        
        $this->passwordHash = PasswordHash::fromPlainPassword($newPassword);
        $this->updatedAt = new \DateTimeImmutable();
        $this->addEvent(new UserPasswordChanged($this->id));
    }
    
    public function updateProfile(UserProfile $profile): void
    {
        $this->profile = $profile;
        $this->updatedAt = new \DateTimeImmutable();
    }
    
    public function deactivate(): void
    {
        $this->isActive = false;
        $this->updatedAt = new \DateTimeImmutable();
    }
    
    public function getId(): string
    {
        return $this->id;
    }
    
    public function getEmail(): Email
    {
        return $this->email;
    }
    
    public function getPasswordHash(): PasswordHash
    {
        return $this->passwordHash;
    }
    
    public function getEvents(): array
    {
        return $this->events;
    }
    
    public function clearEvents(): void
    {
        $this->events = [];
    }
    
    private function addEvent(DomainEvent $event): void
    {
        $this->events[] = $event;
    }
}
```

#### 6.1.2. Value Objects (PHP)

**Email**
```php
<?php
// app/Identity/Domain/ValueObjects/Email.php

namespace App\Identity\Domain\ValueObjects;

class Email
{
    private string $value;
    
    public function __construct(string $email)
    {
        $email = strtolower(trim($email));
        
        if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
            throw new \InvalidArgumentException('Invalid email format');
        }
        
        $this->value = $email;
    }
    
    public function equals(Email $other): bool
    {
        return $this->value === $other->value;
    }
    
    public function __toString(): string
    {
        return $this->value;
    }
    
    public function getValue(): string
    {
        return $this->value;
    }
}
```

**PasswordHash**
```php
<?php
// app/Identity/Domain/ValueObjects/PasswordHash.php

namespace App\Identity\Domain\ValueObjects;

class PasswordHash
{
    private string $value;
    
    private function __construct(string $hash)
    {
        $this->value = $hash;
    }
    
    public static function fromPlainPassword(string $password): self
    {
        if (strlen($password) < 8) {
            throw new \InvalidArgumentException('Password must be at least 8 characters');
        }
        
        return new self(password_hash($password, PASSWORD_BCRYPT));
    }
    
    public static function fromHash(string $hash): self
    {
        return new self($hash);
    }
    
    public function verify(string $plainPassword): bool
    {
        return password_verify($plainPassword, $this->value);
    }
    
    public function getValue(): string
    {
        return $this->value;
    }
}
```

#### 6.1.3. Repositórios (PHP)

**Interface**
```php
<?php
// app/Identity/Domain/Repositories/UserRepositoryInterface.php

namespace App\Identity\Domain\Repositories;

use App\Identity\Domain\Entities\User;
use App\Identity\Domain\ValueObjects\Email;

interface UserRepositoryInterface
{
    public function findById(string $id): ?User;
    public function findByEmail(Email $email): ?User;
    public function save(User $user): void;
    public function delete(string $id): void;
    public function exists(Email $email): bool;
}
```

**Implementação Eloquent**
```php
<?php
// app/Identity/Infrastructure/Persistence/EloquentUserRepository.php

namespace App\Identity\Infrastructure\Persistence;

use App\Identity\Domain\Entities\User;
use App\Identity\Domain\Repositories\UserRepositoryInterface;
use App\Identity\Domain\ValueObjects\Email;
use App\Identity\Domain\ValueObjects\PasswordHash;
use App\Identity\Domain\ValueObjects\UserName;
use App\Models\User as UserModel;

class EloquentUserRepository implements UserRepositoryInterface
{
    public function findById(string $id): ?User
    {
        $model = UserModel::find($id);
        
        if (!$model) {
            return null;
        }
        
        return $this->toDomain($model);
    }
    
    public function findByEmail(Email $email): ?User
    {
        $model = UserModel::where('email', $email->getValue())->first();
        
        if (!$model) {
            return null;
        }
        
        return $this->toDomain($model);
    }
    
    public function save(User $user): void
    {
        $data = $this->toPersistence($user);
        
        UserModel::updateOrCreate(
            ['id' => $data['id']],
            $data
        );
    }
    
    public function delete(string $id): void
    {
        UserModel::destroy($id);
    }
    
    public function exists(Email $email): bool
    {
        return UserModel::where('email', $email->getValue())->exists();
    }
    
    private function toDomain(UserModel $model): User
    {
        $email = new Email($model->email);
        $passwordHash = PasswordHash::fromHash($model->password_hash);
        $name = new UserName($model->first_name, $model->last_name);
        
        // Usar reflection ou método estático para reconstruir
        $user = User::fromPersistence([
            'id' => $model->id,
            'email' => $email,
            'passwordHash' => $passwordHash,
            'name' => $name,
            'createdAt' => $model->created_at,
            'updatedAt' => $model->updated_at,
            'isActive' => $model->is_active,
        ]);
        
        return $user;
    }
    
    private function toPersistence(User $user): array
    {
        return [
            'id' => $user->getId(),
            'email' => $user->getEmail()->getValue(),
            'password_hash' => $user->getPasswordHash()->getValue(),
            'first_name' => $user->getName()->getFirstName(),
            'last_name' => $user->getName()->getLastName(),
            'is_active' => $user->isActive(),
            'created_at' => $user->getCreatedAt(),
            'updated_at' => $user->getUpdatedAt(),
        ];
    }
}
```

#### 6.1.4. Use Cases (PHP)

**RegisterUserUseCase**
```php
<?php
// app/Identity/Application/UseCases/RegisterUserUseCase.php

namespace App\Identity\Application\UseCases;

use App\Identity\Domain\Entities\User;
use App\Identity\Domain\Repositories\UserRepositoryInterface;
use App\Identity\Domain\ValueObjects\Email;
use App\Identity\Domain\ValueObjects\PasswordHash;
use App\Identity\Domain\ValueObjects\UserName;
use Illuminate\Events\Dispatcher;

class RegisterUserUseCase
{
    public function __construct(
        private UserRepositoryInterface $userRepository,
        private Dispatcher $eventDispatcher
    ) {}
    
    public function execute(RegisterUserDTO $dto): RegisterUserOutput
    {
        $email = new Email($dto->email);
        
        if ($this->userRepository->exists($email)) {
            throw new \InvalidArgumentException('User already exists');
        }
        
        $passwordHash = PasswordHash::fromPlainPassword($dto->password);
        $name = new UserName($dto->firstName, $dto->lastName);
        
        $user = User::create($email, $passwordHash, $name);
        
        $this->userRepository->save($user);
        
        // Publicar eventos de domínio
        foreach ($user->getEvents() as $event) {
            $this->eventDispatcher->dispatch($event);
        }
        $user->clearEvents();
        
        return new RegisterUserOutput(
            userId: $user->getId(),
            email: $user->getEmail()->getValue()
        );
    }
}

class RegisterUserDTO
{
    public function __construct(
        public string $email,
        public string $password,
        public string $firstName,
        public string $lastName
    ) {}
}

class RegisterUserOutput
{
    public function __construct(
        public string $userId,
        public string $email
    ) {}
}
```

#### 6.1.5. Controllers (Laravel)

**AuthController**
```php
<?php
// app/Identity/Presentation/Controllers/AuthController.php

namespace App\Identity\Presentation\Controllers;

use App\Identity\Application\UseCases\RegisterUserUseCase;
use App\Identity\Presentation\Requests\RegisterUserRequest;
use App\Identity\Presentation\Resources\UserResource;
use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;

class AuthController extends Controller
{
    public function __construct(
        private RegisterUserUseCase $registerUserUseCase
    ) {}
    
    public function register(RegisterUserRequest $request): JsonResponse
    {
        $dto = new \App\Identity\Application\UseCases\RegisterUserDTO(
            email: $request->email,
            password: $request->password,
            firstName: $request->firstName,
            lastName: $request->lastName
        );
        
        $output = $this->registerUserUseCase->execute($dto);
        
        return response()->json($output, 201);
    }
}
```

**Form Request (Validação)**
```php
<?php
// app/Identity/Presentation/Requests/RegisterUserRequest.php

namespace App\Identity\Presentation\Requests;

use Illuminate\Foundation\Http\FormRequest;

class RegisterUserRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'email' => ['required', 'email', 'unique:users,email'],
            'password' => ['required', 'string', 'min:8'],
            'firstName' => ['required', 'string', 'max:100'],
            'lastName' => ['required', 'string', 'max:100'],
        ];
    }
    
    public function messages(): array
    {
        return [
            'email.required' => 'Email é obrigatório',
            'email.email' => 'Email inválido',
            'email.unique' => 'Email já cadastrado',
            'password.required' => 'Senha é obrigatória',
            'password.min' => 'Senha deve ter no mínimo 8 caracteres',
        ];
    }
}
```

### 6.2. Transaction Context (Core Domain)

#### 6.2.1. Entidade Transaction (PHP)

```php
<?php
// app/Transaction/Domain/Entities/Transaction.php

namespace App\Transaction\Domain\Entities;

use App\Shared\Domain\ValueObjects\Money;
use App\Transaction\Domain\ValueObjects\TransactionType;
use App\Transaction\Domain\ValueObjects\TransactionStatus;
use App\Transaction\Domain\ValueObjects\TransactionDescription;
use App\Transaction\Domain\Events\TransactionApproved;
use App\Transaction\Domain\Events\TransactionCancelled;

class Transaction
{
    private string $id;
    private string $userId;
    private string $accountId;
    private string $categoryId;
    private TransactionType $type;
    private Money $amount;
    private TransactionDescription $description;
    private \DateTimeImmutable $date;
    private TransactionStatus $status;
    private string $context;
    private \DateTimeImmutable $createdAt;
    private \DateTimeImmutable $updatedAt;
    
    /** @var \App\Shared\Domain\Events\DomainEvent[] */
    private array $events = [];
    
    private function __construct(
        string $id,
        string $userId,
        string $accountId,
        string $categoryId,
        TransactionType $type,
        Money $amount,
        TransactionDescription $description,
        \DateTimeImmutable $date,
        string $context
    ) {
        $this->id = $id;
        $this->userId = $userId;
        $this->accountId = $accountId;
        $this->categoryId = $categoryId;
        $this->type = $type;
        $this->amount = $amount;
        $this->description = $description;
        $this->date = $date;
        $this->status = TransactionStatus::pending();
        $this->context = $context;
        $this->createdAt = new \DateTimeImmutable();
        $this->updatedAt = new \DateTimeImmutable();
    }
    
    public static function create(
        string $userId,
        string $accountId,
        string $categoryId,
        TransactionType $type,
        Money $amount,
        TransactionDescription $description,
        \DateTimeImmutable $date,
        string $context
    ): self {
        return new self(
            \Ramsey\Uuid\Uuid::uuid4()->toString(),
            $userId,
            $accountId,
            $categoryId,
            $type,
            $amount,
            $description,
            $date,
            $context
        );
    }
    
    public function approve(): void
    {
        if (!$this->status->canBeApproved()) {
            throw new \InvalidArgumentException('Transaction cannot be approved');
        }
        
        $this->status = TransactionStatus::approved();
        $this->updatedAt = new \DateTimeImmutable();
        $this->addEvent(new TransactionApproved($this->id, $this->amount));
    }
    
    public function cancel(): void
    {
        if (!$this->status->canBeCancelled()) {
            throw new \InvalidArgumentException('Transaction cannot be cancelled');
        }
        
        $this->status = TransactionStatus::cancelled();
        $this->updatedAt = new \DateTimeImmutable();
        $this->addEvent(new TransactionCancelled($this->id));
    }
    
    public function getId(): string
    {
        return $this->id;
    }
    
    public function getAmount(): Money
    {
        return $this->amount;
    }
    
    public function getEvents(): array
    {
        return $this->events;
    }
    
    public function clearEvents(): void
    {
        $this->events = [];
    }
    
    private function addEvent(\App\Shared\Domain\Events\DomainEvent $event): void
    {
        $this->events[] = $event;
    }
}
```

## 7. Migrations (Laravel)

### 7.1. Migration Users

```php
<?php
// database/migrations/xxxx_create_users_table.php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('users', function (Blueprint $table) {
            $table->uuid('id')->primary();
            $table->string('email')->unique();
            $table->string('password_hash');
            $table->string('first_name');
            $table->string('last_name');
            $table->string('currency', 3)->default('BRL');
            $table->string('locale', 10)->default('pt-BR');
            $table->string('theme', 10)->default('light');
            $table->boolean('is_active')->default(true);
            $table->timestamps();
            
            $table->index('email');
        });
    }
    
    public function down(): void
    {
        Schema::dropIfExists('users');
    }
};
```

### 7.2. Migration Transactions

```php
<?php
// database/migrations/xxxx_create_transactions_table.php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('transactions', function (Blueprint $table) {
            $table->uuid('id')->primary();
            $table->uuid('user_id');
            $table->uuid('account_id');
            $table->uuid('category_id');
            $table->string('type'); // INCOME, EXPENSE, TRANSFER
            $table->decimal('amount', 15, 2);
            $table->text('description')->nullable();
            $table->dateTime('date');
            $table->string('status')->default('PENDING');
            $table->string('context'); // PERSONAL, PROFESSIONAL
            $table->json('tags')->nullable();
            $table->timestamps();
            
            $table->foreign('user_id')->references('id')->on('users')->onDelete('cascade');
            $table->foreign('account_id')->references('id')->on('accounts')->onDelete('cascade');
            $table->foreign('category_id')->references('id')->on('categories');
            
            $table->index(['user_id', 'date']);
            $table->index(['account_id']);
            $table->index(['category_id']);
            $table->index(['user_id', 'type', 'date']);
        });
    }
    
    public function down(): void
    {
        Schema::dropIfExists('transactions');
    }
};
```

## 8. Event Bus (Laravel)

### 8.1. Event Listener

```php
<?php
// app/Transaction/Infrastructure/Events/TransactionApprovedListener.php

namespace App\Transaction\Infrastructure\Events;

use App\Transaction\Domain\Events\TransactionApproved;
use App\AccountManagement\Domain\Repositories\AccountRepositoryInterface;
use Illuminate\Contracts\Queue\ShouldQueue;

class TransactionApprovedListener implements ShouldQueue
{
    public function __construct(
        private AccountRepositoryInterface $accountRepository
    ) {}
    
    public function handle(TransactionApproved $event): void
    {
        // Atualizar saldo da conta
        $account = $this->accountRepository->findById($event->getAccountId());
        if ($account) {
            $account->credit($event->getAmount());
            $this->accountRepository->save($account);
        }
    }
}
```

### 8.2. Event Service Provider

```php
<?php
// app/Providers/EventServiceProvider.php

namespace App\Providers;

use Illuminate\Foundation\Support\Providers\EventServiceProvider as ServiceProvider;
use App\Transaction\Domain\Events\TransactionApproved;
use App\Transaction\Infrastructure\Events\TransactionApprovedListener;

class EventServiceProvider extends ServiceProvider
{
    protected $listen = [
        TransactionApproved::class => [
            TransactionApprovedListener::class,
        ],
    ];
}
```

## 9. Jobs/Queues (Laravel)

### 9.1. Processar Transações Recorrentes

```php
<?php
// app/Transaction/Application/Jobs/ProcessRecurringTransactionsJob.php

namespace App\Transaction\Application\Jobs;

use App\Transaction\Domain\Repositories\RecurringTransactionRepositoryInterface;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;

class ProcessRecurringTransactionsJob implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;
    
    public function __construct(
        private RecurringTransactionRepositoryInterface $recurringTransactionRepository
    ) {}
    
    public function handle(): void
    {
        $today = new \DateTimeImmutable();
        $recurringTransactions = $this->recurringTransactionRepository
            ->findActiveByDate($today);
        
        foreach ($recurringTransactions as $recurring) {
            $transaction = $recurring->generateNextTransaction();
            // Salvar transação e atualizar próxima data
        }
    }
}
```

### 9.2. Agendar Job

```php
<?php
// app/Console/Kernel.php

namespace App\Console;

use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Console\Kernel as ConsoleKernel;
use App\Transaction\Application\Jobs\ProcessRecurringTransactionsJob;

class Kernel extends ConsoleKernel
{
    protected function schedule(Schedule $schedule): void
    {
        // Executar diariamente às 00:00
        $schedule->job(new ProcessRecurringTransactionsJob)
            ->daily();
    }
}
```

## 10. Testes em PHP

### 10.1. Testes Unitários

```php
<?php
// tests/Unit/Transaction/TransactionTest.php

namespace Tests\Unit\Transaction;

use PHPUnit\Framework\TestCase;
use App\Transaction\Domain\Entities\Transaction;
use App\Transaction\Domain\ValueObjects\TransactionType;
use App\Shared\Domain\ValueObjects\Money;

class TransactionTest extends TestCase
{
    public function test_can_create_transaction(): void
    {
        $transaction = Transaction::create(
            'user-id',
            'account-id',
            'category-id',
            TransactionType::expense(),
            Money::create(100, 'BRL'),
            TransactionDescription::create('Test'),
            new \DateTimeImmutable(),
            'PERSONAL'
        );
        
        $this->assertInstanceOf(Transaction::class, $transaction);
        $this->assertEquals(100, $transaction->getAmount()->getValue());
    }
    
    public function test_can_approve_transaction(): void
    {
        $transaction = Transaction::create(/* ... */);
        
        $transaction->approve();
        
        $this->assertTrue($transaction->getStatus()->isApproved());
        $this->assertCount(1, $transaction->getEvents());
    }
}
```

### 10.2. Testes de Integração

```php
<?php
// tests/Feature/Transaction/CreateTransactionTest.php

namespace Tests\Feature\Transaction;

use Tests\TestCase;
use Illuminate\Foundation\Testing\RefreshDatabase;

class CreateTransactionTest extends TestCase
{
    use RefreshDatabase;
    
    public function test_can_create_transaction(): void
    {
        $user = User::factory()->create();
        
        $response = $this->actingAs($user)
            ->postJson('/api/transactions', [
                'accountId' => 'account-id',
                'categoryId' => 'category-id',
                'type' => 'EXPENSE',
                'amount' => 100,
                'description' => 'Test',
            ]);
        
        $response->assertStatus(201)
            ->assertJsonStructure([
                'id',
                'amount',
                'status',
            ]);
    }
}
```

## 11. Fases de Desenvolvimento

### Fase 1: Fundação (1-2 semanas)
- Setup Laravel
- Shared Kernel
- Identity Context
- Account Management Context

### Fase 2: Core Domain (2 semanas)
- Transaction Context
- Integração com Account
- Event Bus
- Jobs/Queues

### Fase 3: Expansão (2-3 semanas)
- Category Context
- Budget Context
- Recurring Transactions
- Relatórios básicos

### Fase 4: Funcionalidades Avançadas (2-3 semanas)
- Investment Context
- Goal Context
- Notification Context
- Dashboard completo

## 12. Performance e Otimizações

### 12.1. Cache (Redis)

```php
<?php
// app/Reporting/Application/Services/ReportCacheService.php

namespace App\Reporting\Application\Services;

use Illuminate\Support\Facades\Cache;

class ReportCacheService
{
    public function getOrGenerate(string $key, callable $generator, int $ttl = 3600): mixed
    {
        return Cache::remember($key, $ttl, $generator);
    }
    
    public function invalidate(string $pattern): void
    {
        Cache::forget($pattern);
    }
}
```

### 12.2. Eager Loading (Eloquent)

```php
<?php
// app/Transaction/Infrastructure/Persistence/EloquentTransactionRepository.php

public function findByUserId(string $userId): array
{
    return TransactionModel::where('user_id', $userId)
        ->with(['account', 'category']) // Eager loading
        ->get()
        ->map(fn($model) => $this->toDomain($model))
        ->toArray();
}
```

### 12.3. Query Optimization

```php
<?php
// Usar índices e queries otimizadas
public function findByDateRange(string $userId, \DateTimeImmutable $start, \DateTimeImmutable $end): array
{
    return TransactionModel::where('user_id', $userId)
        ->whereBetween('date', [$start, $end])
        ->select(['id', 'amount', 'type', 'date']) // Selecionar apenas campos necessários
        ->get()
        ->map(fn($model) => $this->toDomain($model))
        ->toArray();
}
```

## 13. Deploy e DevOps

### 13.1. Dockerfile

```dockerfile
FROM php:8.2-fpm

# Instalar dependências
RUN apt-get update && apt-get install -y \
    git \
    curl \
    libpng-dev \
    libonig-dev \
    libxml2-dev \
    zip \
    unzip

# Instalar extensões PHP
RUN docker-php-ext-install pdo_mysql pdo_pgsql mbstring exif pcntl bcmath gd

# Instalar Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www

COPY . .

RUN composer install --no-dev --optimize-autoloader

CMD php artisan serve --host=0.0.0.0 --port=8000
```

### 13.2. docker-compose.yml

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8000:8000"
    volumes:
      - .:/var/www
    environment:
      - DB_CONNECTION=pgsql
      - DB_HOST=db
      - DB_DATABASE=gestao_financeira
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

### 14.1. Logging Estruturado (Monolog)

```php
<?php
// config/logging.php
return [
    'channels' => [
        'stack' => [
            'driver' => 'stack',
            'channels' => ['daily', 'sentry'],
        ],
        'daily' => [
            'driver' => 'daily',
            'path' => storage_path('logs/laravel.log'),
            'level' => env('LOG_LEVEL', 'debug'),
            'days' => 14,
        ],
        'sentry' => [
            'driver' => 'sentry',
            'level' => 'error',
        ],
    ],
];

// Uso
use Illuminate\Support\Facades\Log;

Log::info('Transaction created', [
    'userId' => $user->getId(),
    'transactionId' => $transaction->getId(),
    'amount' => $transaction->getAmount()->getValue(),
    'requestId' => request()->header('X-Request-ID'),
]);
```

### 14.2. Métricas com Prometheus

```php
<?php
// app/Shared/Infrastructure/Metrics/PrometheusService.php

namespace App\Shared\Infrastructure\Metrics;

use Prometheus\CollectorRegistry;
use Prometheus\Storage\Redis;
use Prometheus\Counter;
use Prometheus\Histogram;

class PrometheusService
{
    private CollectorRegistry $registry;
    private Counter $httpRequestsTotal;
    private Histogram $httpRequestDuration;
    
    public function __construct()
    {
        $adapter = new Redis(['host' => env('REDIS_HOST')]);
        $this->registry = new CollectorRegistry($adapter);
        
        $this->httpRequestsTotal = $this->registry->getOrRegisterCounter(
            'app',
            'http_requests_total',
            'Total number of HTTP requests',
            ['method', 'endpoint', 'status']
        );
        
        $this->httpRequestDuration = $this->registry->getOrRegisterHistogram(
            'app',
            'http_request_duration_seconds',
            'HTTP request duration in seconds',
            ['method', 'endpoint']
        );
    }
    
    public function recordRequest(string $method, string $endpoint, int $status, float $duration): void
    {
        $this->httpRequestsTotal->inc([$method, $endpoint, (string)$status]);
        $this->httpRequestDuration->observe([$method, $endpoint], $duration);
    }
    
    public function getMetrics(): string
    {
        $renderer = new \Prometheus\RenderTextFormat();
        return $renderer->render($this->registry->getMetricFamilySamples());
    }
}

// Middleware
class PrometheusMiddleware
{
    public function handle($request, \Closure $next)
    {
        $start = microtime(true);
        $response = $next($request);
        $duration = microtime(true) - $start;
        
        app(PrometheusService::class)->recordRequest(
            $request->method(),
            $request->path(),
            $response->getStatusCode(),
            $duration
        );
        
        return $response;
    }
}
```

### 14.3. Tracing com OpenTelemetry

```php
<?php
// app/Shared/Infrastructure/Tracing/TracingService.php

namespace App\Shared\Infrastructure\Tracing;

use OpenTelemetry\SDK\Trace\TracerProvider;
use OpenTelemetry\SDK\Trace\SpanExporter\ConsoleSpanExporter;
use OpenTelemetry\SDK\Trace\SpanProcessor\BatchSpanProcessor;
use OpenTelemetry\SDK\Resource\ResourceInfo;
use OpenTelemetry\SemConv\ResourceAttributes;

class TracingService
{
    private TracerProvider $tracerProvider;
    
    public function __construct()
    {
        $exporter = new ConsoleSpanExporter();
        $processor = new BatchSpanProcessor($exporter);
        
        $resource = ResourceInfo::create([
            ResourceAttributes::SERVICE_NAME => 'gestao-financeira',
            ResourceAttributes::SERVICE_VERSION => '1.0.0',
        ]);
        
        $this->tracerProvider = new TracerProvider($processor, null, $resource);
    }
    
    public function getTracer(string $name = 'gestao-financeira'): \OpenTelemetry\API\Trace\TracerInterface
    {
        return $this->tracerProvider->getTracer($name);
    }
}

// Uso
$tracer = app(TracingService::class)->getTracer();
$span = $tracer->spanBuilder('create-transaction')->startSpan();
$scope = $span->activate();

try {
    // Lógica do use case
} finally {
    $span->end();
    $scope->detach();
}
```

### 14.4. Health Checks Robustos

```php
<?php
// app/Shared/Infrastructure/Health/HealthController.php

namespace App\Shared\Infrastructure\Health;

use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Redis;
use Illuminate\Support\Facades\Cache;

class HealthController
{
    public function liveness(): JsonResponse
    {
        return response()->json(['status' => 'alive']);
    }
    
    public function readiness(): JsonResponse
    {
        $checks = [
            'database' => $this->checkDatabase(),
            'redis' => $this->checkRedis(),
            'cache' => $this->checkCache(),
        ];
        
        $allHealthy = !in_array(false, $checks);
        
        return response()->json([
            'status' => $allHealthy ? 'ready' : 'not_ready',
            'checks' => $checks,
        ], $allHealthy ? 200 : 503);
    }
    
    private function checkDatabase(): bool
    {
        try {
            DB::connection()->getPdo();
            return true;
        } catch (\Exception $e) {
            return false;
        }
    }
    
    private function checkRedis(): bool
    {
        try {
            Redis::ping();
            return true;
        } catch (\Exception $e) {
            return false;
        }
    }
    
    private function checkCache(): bool
    {
        try {
            Cache::put('health_check', 'ok', 1);
            return Cache::get('health_check') === 'ok';
        } catch (\Exception $e) {
            return false;
        }
    }
}
```

## 15. Segurança

### 15.1. Headers de Segurança

```php
<?php
// app/Http/Middleware/SecurityHeaders.php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;

class SecurityHeaders
{
    public function handle(Request $request, Closure $next)
    {
        $response = $next($request);
        
        $response->headers->set('X-Content-Type-Options', 'nosniff');
        $response->headers->set('X-Frame-Options', 'DENY');
        $response->headers->set('X-XSS-Protection', '1; mode=block');
        $response->headers->set('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
        $response->headers->set('Content-Security-Policy', "default-src 'self'");
        $response->headers->set('Referrer-Policy', 'strict-origin-when-cross-origin');
        
        return $response;
    }
}
```

### 15.2. Rate Limiting

```php
<?php
// app/Http/Kernel.php
protected $middlewareGroups = [
    'api' => [
        'throttle:api',
        // ...
    ],
];

// app/Providers/RouteServiceProvider.php
protected function configureRateLimiting()
{
    RateLimiter::for('api', function (Request $request) {
        return Limit::perMinute(100)->by($request->user()?->id ?: $request->ip());
    });
    
    RateLimiter::for('auth', function (Request $request) {
        return Limit::perMinute(5)->by($request->ip());
    });
}

// Uso em rotas
Route::middleware(['throttle:auth'])->group(function () {
    Route::post('/login', [AuthController::class, 'login']);
});
```

### 15.3. Validação de Entrada Robusta

```php
<?php
// app/Identity/Presentation/Requests/RegisterUserRequest.php

namespace App\Identity\Presentation\Requests;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rules\Password;

class RegisterUserRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'email' => ['required', 'email', 'unique:users,email', 'max:255'],
            'password' => ['required', 'string', Password::min(8)->mixedCase()->numbers()->symbols()],
            'firstName' => ['required', 'string', 'max:100', 'regex:/^[a-zA-Z\s]+$/'],
            'lastName' => ['required', 'string', 'max:100', 'regex:/^[a-zA-Z\s]+$/'],
        ];
    }
    
    public function messages(): array
    {
        return [
            'email.required' => 'Email é obrigatório',
            'email.email' => 'Email inválido',
            'email.unique' => 'Email já cadastrado',
            'password.required' => 'Senha é obrigatória',
            'password.min' => 'Senha deve ter no mínimo 8 caracteres',
            'firstName.required' => 'Primeiro nome é obrigatório',
            'firstName.regex' => 'Primeiro nome contém caracteres inválidos',
        ];
    }
    
    protected function prepareForValidation(): void
    {
        $this->merge([
            'email' => strtolower(trim($this->email)),
            'firstName' => trim($this->firstName),
            'lastName' => trim($this->lastName),
        ]);
    }
}
```

### 15.4. Proteção contra SQL Injection e XSS

```php
<?php
// Eloquent usa prepared statements automaticamente
// Mas para queries raw:
$transactions = DB::select(
    'SELECT * FROM transactions WHERE user_id = ? AND date BETWEEN ? AND ?',
    [$userId, $startDate, $endDate]
);

// Sanitização de entrada
use Illuminate\Support\Str;

function sanitizeInput(string $input): string
{
    return strip_tags(trim($input));
}

// Para HTML (se necessário)
use HTMLPurifier;

function sanitizeHtml(string $html): string
{
    $purifier = new HTMLPurifier();
    return $purifier->purify($html);
}
```

## 16. Documentação da API

### 16.1. Swagger/OpenAPI (L5 Swagger)

```php
<?php
// composer.json
{
    "require": {
        "darkaonline/l5-swagger": "^10.0"
    }
}

// config/l5-swagger.php
return [
    'default' => 'default',
    'documentations' => [
        'default' => [
            'api' => [
                'title' => 'Gestão Financeira API',
                'version' => '1.0.0',
            ],
            'routes' => [
                'api' => 'api/documentation',
            ],
        ],
    ],
];

// Controller com anotações Swagger
/**
 * @OA\Tag(
 *     name="Transactions",
 *     description="Endpoints de transações"
 * )
 */
class TransactionController extends Controller
{
    /**
     * @OA\Post(
     *     path="/api/v1/transactions",
     *     summary="Criar transação",
     *     tags={"Transactions"},
     *     security={{"bearerAuth":{}}},
     *     @OA\RequestBody(
     *         required=true,
     *         @OA\JsonContent(ref="#/components/schemas/CreateTransactionRequest")
     *     ),
     *     @OA\Response(
     *         response=201,
     *         description="Transação criada com sucesso",
     *         @OA\JsonContent(ref="#/components/schemas/TransactionResponse")
     *     ),
     *     @OA\Response(
     *         response=400,
     *         description="Dados inválidos"
     *     )
     * )
     */
    public function store(CreateTransactionRequest $request)
    {
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
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testing
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.2'
          extensions: pdo, pdo_pgsql, redis
      
      - name: Copy .env
        run: php -r "file_exists('.env') || copy('.env.example', '.env');"
      
      - name: Install Dependencies
        run: composer install -q --no-ansi --no-interaction --no-scripts --no-progress --prefer-dist
      
      - name: Generate key
        run: php artisan key:generate
      
      - name: Directory Permissions
        run: chmod -R 777 storage bootstrap/cache
      
      - name: Create Database
        run: php artisan migrate --force
      
      - name: Execute tests (Unit and Feature tests) via PHPUnit
        env:
          DB_CONNECTION: pgsql
          DB_HOST: postgres
          DB_PORT: 5432
          DB_DATABASE: testing
          DB_USERNAME: postgres
          DB_PASSWORD: postgres
        run: vendor/bin/phpunit
  
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: gestao-financeira:${{ github.sha }},gestao-financeira:latest
```

### 17.2. Deploy em Produção

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  app:
    image: gestao-financeira:latest
    restart: always
    ports:
      - "8000:8000"
    environment:
      - APP_ENV=production
      - APP_DEBUG=false
      - DB_CONNECTION=pgsql
      - DB_HOST=db
      - REDIS_HOST=redis
      - QUEUE_CONNECTION=redis
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health/live"]
      interval: 30s
      timeout: 10s
      retries: 3
    depends_on:
      - db
      - redis
  
  queue:
    image: gestao-financeira:latest
    command: php artisan queue:work --tries=3
    restart: always
    environment:
      - APP_ENV=production
      - DB_CONNECTION=pgsql
      - DB_HOST=db
      - REDIS_HOST=redis
      - QUEUE_CONNECTION=redis
    depends_on:
      - db
      - redis
  
  scheduler:
    image: gestao-financeira:latest
    command: php artisan schedule:work
    restart: always
    environment:
      - APP_ENV=production
      - DB_CONNECTION=pgsql
      - DB_HOST=db
    depends_on:
      - db
  
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

```php
<?php
// app/Console/Commands/BackupDatabase.php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\DB;

class BackupDatabase extends Command
{
    protected $signature = 'db:backup';
    protected $description = 'Backup do banco de dados';
    
    public function handle(): int
    {
        $this->info('Iniciando backup do banco de dados...');
        
        $timestamp = now()->format('Y-m-d_H-i-s');
        $filename = "backup_{$timestamp}.sql";
        $path = storage_path("app/backups/{$filename}");
        
        $config = config('database.connections.pgsql');
        $command = sprintf(
            'PGPASSWORD=%s pg_dump -h %s -U %s -d %s -F c -f %s',
            $config['password'],
            $config['host'],
            $config['username'],
            $config['database'],
            $path
        );
        
        exec($command, $output, $returnVar);
        
        if ($returnVar === 0) {
            $this->info("Backup criado: {$filename}");
            
            // Upload para S3 (opcional)
            if (config('backup.s3.enabled')) {
                Storage::disk('s3')->put("backups/{$filename}", file_get_contents($path));
                $this->info("Backup enviado para S3");
            }
            
            // Limpar backups antigos (manter últimos 30 dias)
            $this->cleanOldBackups();
            
            return Command::SUCCESS;
        }
        
        $this->error('Falha ao criar backup');
        return Command::FAILURE;
    }
    
    private function cleanOldBackups(): void
    {
        $files = Storage::disk('local')->files('backups');
        $cutoff = now()->subDays(30);
        
        foreach ($files as $file) {
            $lastModified = Storage::disk('local')->lastModified($file);
            if ($lastModified < $cutoff->timestamp) {
                Storage::disk('local')->delete($file);
                $this->info("Backup antigo removido: {$file}");
            }
        }
    }
}

// Agendar no Kernel
protected function schedule(Schedule $schedule): void
{
    $schedule->command('db:backup')->dailyAt('02:00');
}
```

## 19. Testes E2E

### 19.1. Testes End-to-End (Laravel Dusk)

```php
<?php
// tests/Browser/TransactionTest.php

namespace Tests\Browser;

use Laravel\Dusk\Browser;
use Tests\DuskTestCase;
use App\Models\User;
use App\Models\Account;
use App\Models\Category;

class TransactionTest extends DuskTestCase
{
    public function test_user_can_create_transaction(): void
    {
        $user = User::factory()->create();
        $account = Account::factory()->create(['user_id' => $user->id]);
        $category = Category::factory()->create(['user_id' => $user->id]);
        
        $this->browse(function (Browser $browser) use ($user, $account, $category) {
            $browser->loginAs($user)
                    ->visit('/transactions/create')
                    ->type('amount', '100.00')
                    ->type('description', 'Test transaction')
                    ->select('account_id', $account->id)
                    ->select('category_id', $category->id)
                    ->select('type', 'EXPENSE')
                    ->press('Salvar')
                    ->assertPathIs('/transactions')
                    ->assertSee('Transação criada com sucesso');
        });
    }
}
```

## 20. Auditoria e Compliance

### 20.1. Log de Auditoria

```php
<?php
// app/Shared/Infrastructure/Audit/AuditService.php

namespace App\Shared\Infrastructure\Audit;

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class AuditService
{
    public function log(
        string $userId,
        string $action,
        string $resource,
        string $resourceId,
        ?array $metadata = null
    ): void {
        DB::table('audit_logs')->insert([
            'user_id' => $userId,
            'action' => $action,
            'resource' => $resource,
            'resource_id' => $resourceId,
            'metadata' => json_encode($metadata),
            'ip_address' => request()->ip(),
            'user_agent' => request()->userAgent(),
            'created_at' => now(),
        ]);
        
        Log::info('Audit log created', [
            'userId' => $userId,
            'action' => $action,
            'resource' => $resource,
            'resourceId' => $resourceId,
        ]);
    }
}

// Middleware de auditoria
class AuditMiddleware
{
    public function handle($request, \Closure $next)
    {
        $response = $next($request);
        
        if ($request->user() && in_array($request->method(), ['POST', 'PUT', 'DELETE'])) {
            app(AuditService::class)->log(
                $request->user()->id,
                $request->method(),
                $this->extractResource($request->path()),
                $request->route('id') ?? 'N/A',
                [
                    'path' => $request->path(),
                    'ip' => $request->ip(),
                ]
            );
        }
        
        return $response;
    }
    
    private function extractResource(string $path): string
    {
        $parts = explode('/', $path);
        return $parts[2] ?? 'unknown';
    }
}
```

### 20.2. LGPD/GDPR Compliance

```php
<?php
// app/Identity/Application/UseCases/DeleteUserDataUseCase.php

namespace App\Identity\Application\UseCases;

use App\Identity\Domain\Repositories\UserRepositoryInterface;
use App\Shared\Infrastructure\Audit\AuditService;

class DeleteUserDataUseCase
{
    public function __construct(
        private UserRepositoryInterface $userRepository,
        private AuditService $auditService
    ) {}
    
    public function execute(string $userId): void
    {
        // 1. Anonimizar dados pessoais
        $user = $this->userRepository->findById($userId);
        if ($user) {
            $user->anonymize();
            $this->userRepository->save($user);
        }
        
        // 2. Manter dados financeiros agregados (se necessário para compliance fiscal)
        // 3. Registrar ação de exclusão
        $this->auditService->log($userId, 'DELETE', 'user', $userId);
        
        // 4. Notificar usuário
        // ...
    }
}

// Exportação de dados
class ExportUserDataUseCase
{
    public function execute(string $userId): array
    {
        return [
            'user' => $this->userRepository->findById($userId),
            'transactions' => $this->transactionRepository->findByUserId($userId),
            'accounts' => $this->accountRepository->findByUserId($userId),
            // ... outros dados
        ];
    }
}
```

## 21. Escalabilidade e Multi-tenancy

### 21.1. Estratégia de Escalabilidade

**Horizontal Scaling:**
- Múltiplas instâncias da API (load balancer)
- Database read replicas para relatórios
- Redis cluster para cache distribuído
- Queue workers distribuídos

**Vertical Scaling:**
- Otimização de queries
- Índices estratégicos
- Connection pooling otimizado
- Cache agressivo
- OPcache para PHP

### 21.2. Multi-tenancy (Opcional - se virar SaaS)

```php
<?php
// Estrutura para suportar múltiplos tenants
class Tenant
{
    public function __construct(
        private string $id,
        private string $name,
        private string $plan, // FREE, PREMIUM, ENTERPRISE
        private array $settings
    ) {}
}

// Middleware para isolar dados por tenant
class TenantMiddleware
{
    public function handle($request, \Closure $next)
    {
        $tenantId = $request->header('X-Tenant-ID') ?? $request->user()?->tenant_id;
        
        if ($tenantId) {
            app()->instance('tenant_id', $tenantId);
        }
        
        return $next($request);
    }
}

// Scope global para filtrar por tenant
class TenantScope implements Scope
{
    public function apply(Builder $builder, Model $model)
    {
        $tenantId = app('tenant_id');
        if ($tenantId) {
            $builder->where('tenant_id', $tenantId);
        }
    }
}
```

## 22. Tratamento de Erros Robusto

### 22.1. Erros de Domínio

```php
<?php
// app/Shared/Domain/Exceptions/DomainException.php

namespace App\Shared\Domain\Exceptions;

class DomainException extends \Exception
{
    public function __construct(
        private string $code,
        string $message,
        private ?array $details = null
    ) {
        parent::__construct($message);
    }
    
    public function getCode(): string
    {
        return $this->code;
    }
    
    public function getDetails(): ?array
    {
        return $this->details;
    }
}

// Erros específicos
class InsufficientBalanceException extends DomainException
{
    public function __construct(string $accountId, float $required, float $available)
    {
        parent::__construct(
            'INSUFFICIENT_BALANCE',
            'Account balance is insufficient',
            [
                'accountId' => $accountId,
                'required' => $required,
                'available' => $available,
            ]
        );
    }
}

class TransactionNotFoundException extends DomainException
{
    public function __construct(string $transactionId)
    {
        parent::__construct(
            'TRANSACTION_NOT_FOUND',
            'Transaction not found',
            ['transactionId' => $transactionId]
        );
    }
}
```

### 22.2. Exception Handler Global

```php
<?php
// app/Exceptions/Handler.php

namespace App\Exceptions;

use Illuminate\Foundation\Exceptions\Handler as ExceptionHandler;
use App\Shared\Domain\Exceptions\DomainException;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Throwable;

class Handler extends ExceptionHandler
{
    public function render($request, Throwable $exception): JsonResponse
    {
        if ($exception instanceof DomainException) {
            return response()->json([
                'error' => [
                    'code' => $exception->getCode(),
                    'message' => $exception->getMessage(),
                    'details' => $exception->getDetails(),
                ],
                'requestId' => $request->header('X-Request-ID'),
            ], 400);
        }
        
        if ($exception instanceof \Illuminate\Validation\ValidationException) {
            return response()->json([
                'error' => [
                    'code' => 'VALIDATION_ERROR',
                    'message' => 'Validation failed',
                    'details' => $exception->errors(),
                ],
                'requestId' => $request->header('X-Request-ID'),
            ], 422);
        }
        
        return parent::render($request, $exception);
    }
}
```

## 23. Testes de Performance e Carga

### 23.1. Benchmarks

```php
<?php
// tests/Benchmark/TransactionBenchTest.php

namespace Tests\Benchmark;

use Tests\TestCase;
use App\Models\User;
use App\Models\Transaction;
use Illuminate\Foundation\Testing\RefreshDatabase;

class TransactionBenchTest extends TestCase
{
    use RefreshDatabase;
    
    public function test_list_transactions_performance(): void
    {
        $user = User::factory()->create();
        Transaction::factory()->count(1000)->create(['user_id' => $user->id]);
        
        $start = microtime(true);
        
        for ($i = 0; $i < 100; $i++) {
            $this->getJson("/api/v1/transactions", [
                'Authorization' => "Bearer {$user->createToken('test')->plainTextToken}",
            ]);
        }
        
        $duration = microtime(true) - $start;
        $avgDuration = ($duration / 100) * 1000; // em milissegundos
        
        $this->assertLessThan(500, $avgDuration, "Average response time should be less than 500ms");
    }
}
```

### 23.2. Testes de Carga (Artillery ou Apache Bench)

```yaml
# artillery-config.yml
config:
  target: "http://localhost:8000"
  phases:
    - duration: 60
      arrivalRate: 10
      name: "Warm up"
    - duration: 300
      arrivalRate: 50
      name: "Sustained load"
    - duration: 60
      arrivalRate: 100
      name: "Spike"
  defaults:
    headers:
      Authorization: "Bearer {{ token }}"
scenarios:
  - name: "List transactions"
    flow:
      - get:
          url: "/api/v1/transactions"
```

## 24. Versionamento de API

### 24.1. Estrutura de Versionamento

```php
<?php
// routes/api.php
Route::prefix('v1')->group(function () {
    Route::apiResource('transactions', V1\TransactionController::class);
    // ...
});

Route::prefix('v2')->group(function () {
    Route::apiResource('transactions', V2\TransactionController::class);
    // ...
});

// Middleware de deprecation
class DeprecationMiddleware
{
    public function handle($request, \Closure $next)
    {
        $response = $next($request);
        
        if (str_contains($request->path(), '/v1/')) {
            $response->headers->set('Deprecation', 'true');
            $response->headers->set('Sunset', '2025-12-31');
            $response->headers->set('Link', '</api/v2>; rel="successor-version"');
        }
        
        return $response;
    }
}
```

## 14. Considerações Finais

### Vantagens do PHP + Laravel:
1. ✅ **Produtividade máxima** - Você já domina
2. ✅ **Ecossistema maduro** - Muitas ferramentas
3. ✅ **Eloquent** - ORM produtivo
4. ✅ **Form Requests** - Validação integrada
5. ✅ **Events/Jobs** - Processamento assíncrono
6. ✅ **Performance PHP 8.x** - JIT compiler

### Recomendações:
- Usar **Laravel** para produtividade máxima
- Aproveitar **Form Requests** para validação
- Usar **Laravel Events** para Domain Events
- Implementar **Jobs/Queues** para processamento assíncrono
- Usar **Redis** para cache de relatórios
- Aproveitar **Eloquent** para queries expressivas
- Usar **API Resources** para serialização

### Alternativa: Symfony + Doctrine
- Se preferir **DDD mais puro** (Data Mapper)
- **Doctrine** tem Identity Map e Unit of Work
- Mais complexo, mas mais alinhado com DDD


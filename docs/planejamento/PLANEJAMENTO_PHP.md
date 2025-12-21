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


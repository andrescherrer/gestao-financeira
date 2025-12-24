# Status da Migração Completa para shadcn-vue

## ✅ Componentes Migrados

### Layout
- [x] Header.vue - Migrado para Lucide Vue e shadcn-vue Button
- [x] Sidebar.vue - Migrado para Lucide Vue
- [x] Layout.vue - Sem mudanças necessárias
- [x] Breadcrumbs.vue - Migrado para Lucide Vue

### Componentes
- [x] AccountCard.vue - Migrado para shadcn-vue Card, Badge e Lucide Vue
- [x] TransactionTable.vue - Migrado para shadcn-vue Table, Badge, Button e Lucide Vue
- [x] AccountForm.vue - Migrado para shadcn-vue Input, Label, Button, Card e Lucide Vue
- [x] TransactionForm.vue - Migrado para shadcn-vue Input, Textarea, Label, Button, Card e Lucide Vue

### Configuração
- [x] main.ts - Removido PrimeVue
- [x] package.json - Removidas dependências PrimeVue

## ⏳ Views Pendentes (9 arquivos)

Ainda há 9 views que usam PrimeIcons e precisam ser migradas:

1. LoginView.vue - 6 ocorrências
2. HomeView.vue - 13 ocorrências
3. AccountsView.vue - 8 ocorrências
4. TransactionsView.vue - 9 ocorrências
5. AccountDetailsView.vue - 5 ocorrências
6. TransactionDetailsView.vue - 4 ocorrências
7. NewAccountView.vue - 4 ocorrências
8. NewTransactionView.vue - 4 ocorrências
9. EditAccountView.vue - 6 ocorrências

## Mapeamento de Ícones

| PrimeIcons | Lucide Vue |
|-----------|------------|
| pi-wallet | Wallet |
| pi-user | User |
| pi-briefcase | Briefcase |
| pi-list | List |
| pi-arrow-down | ArrowDown |
| pi-arrow-up | ArrowUp |
| pi-chart-line | TrendingUp |
| pi-plus | Plus |
| pi-chevron-right | ChevronRight |
| pi-spinner pi-spin | Loader2 (com animate-spin) |
| pi-exclamation-circle | AlertCircle |
| pi-inbox | Inbox |
| pi-credit-card | CreditCard |
| pi-home | Home |
| pi-sign-out | LogOut |
| pi-bell | Bell |
| pi-chevron-down | ChevronDown |
| pi-pencil | Pencil |
| pi-arrow-left | ArrowLeft |
| pi-envelope | Mail |
| pi-lock | Lock |
| pi-eye | Eye |
| pi-eye-slash | EyeOff |
| pi-check | Check |

## Próximos Passos

1. Migrar todas as views restantes
2. Testar a aplicação
3. Verificar se há algum componente quebrado
4. Atualizar documentação


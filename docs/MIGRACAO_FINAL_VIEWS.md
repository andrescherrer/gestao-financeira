# Migração Final das Views Restantes

## Views Pendentes (5 arquivos, 39 ocorrências)

1. **HomeView.vue** - 13 ocorrências
   - pi-wallet, pi-user, pi-briefcase, pi-list
   - pi-arrow-down, pi-arrow-up, pi-chart-line
   - pi-plus, pi-chevron-right

2. **AccountsView.vue** - 8 ocorrências
   - pi-plus, pi-spinner, pi-exclamation-circle
   - pi-wallet, pi-user, pi-briefcase

3. **TransactionsView.vue** - 9 ocorrências
   - pi-plus, pi-spinner, pi-exclamation-circle
   - pi-credit-card, pi-list, pi-arrow-down, pi-arrow-up, pi-chart-line

4. **AccountDetailsView.vue** - 5 ocorrências
   - pi-spinner, pi-exclamation-circle, pi-pencil, pi-arrow-left

5. **TransactionDetailsView.vue** - 4 ocorrências
   - pi-spinner, pi-exclamation-circle, pi-arrow-left

## Mapeamento de Ícones

| PrimeIcons | Lucide Vue | Uso |
|-----------|------------|-----|
| pi-wallet | Wallet | Ícone de carteira/conta |
| pi-user | User | Ícone de usuário/pessoal |
| pi-briefcase | Briefcase | Ícone de negócio |
| pi-list | List | Ícone de lista/transações |
| pi-arrow-down | ArrowDown | Receita |
| pi-arrow-up | ArrowUp | Despesa |
| pi-chart-line | TrendingUp | Gráfico/saldo |
| pi-plus | Plus | Adicionar/criar |
| pi-chevron-right | ChevronRight | Navegação |
| pi-spinner pi-spin | Loader2 (animate-spin) | Loading |
| pi-exclamation-circle | AlertCircle | Erro/alerta |
| pi-credit-card | CreditCard | Cartão/transação |
| pi-pencil | Pencil | Editar |
| pi-arrow-left | ArrowLeft | Voltar |

## Estratégia de Migração

1. Substituir todos os `<i class="pi pi-*">` por componentes Lucide Vue
2. Atualizar botões HTML para usar shadcn-vue Button
3. Atualizar cards HTML para usar shadcn-vue Card
4. Atualizar estados de loading/error para usar shadcn-vue componentes
5. Manter a mesma estrutura e funcionalidade


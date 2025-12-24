# Migração para shadcn-vue

## Por que migrar?

- **Design mais moderno**: shadcn-vue oferece componentes com design mais próximo do shadcn/ui original
- **Mais personalizável**: Componentes copiáveis e editáveis diretamente no projeto
- **Melhor integração com Tailwind**: Usa Tailwind CSS nativamente
- **Acessibilidade**: Baseado em Radix Vue (headless UI)
- **TypeScript**: Suporte completo a TypeScript

## Opções disponíveis

### 1. shadcn-vue (Recomendado)
- Porte oficial do shadcn/ui para Vue 3
- Componentes copiáveis (não é uma biblioteca npm)
- Baseado em Radix Vue
- Tailwind CSS

### 2. Radix Vue + Tailwind
- Biblioteca headless (sem estilos)
- Máxima flexibilidade
- Mais trabalho manual

### 3. Manter PrimeVue
- Já está funcionando
- Menos trabalho de migração
- Design diferente do shadcn

## Plano de Migração

### Fase 1: Instalação
1. Instalar shadcn-vue CLI
2. Configurar Tailwind CSS (já temos)
3. Instalar dependências base (Radix Vue, etc)

### Fase 2: Substituição Gradual
1. Substituir ícones (PrimeIcons → Lucide Vue)
2. Substituir componentes um por um
3. Manter PrimeVue durante transição

### Fase 3: Limpeza
1. Remover PrimeVue
2. Remover PrimeIcons
3. Atualizar todos os componentes

## Componentes a Migrar

- [ ] Ícones (PrimeIcons → Lucide Vue)
- [ ] Botões
- [ ] Inputs/Forms
- [ ] Cards
- [ ] Tabelas
- [ ] Dropdowns/Menus
- [ ] Modals/Dialogs
- [ ] Badges/Tags


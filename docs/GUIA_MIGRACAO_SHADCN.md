# Guia de Migração: PrimeVue → shadcn-vue

## Status da Migração

### ✅ Instalado e Configurado
- [x] shadcn-vue CLI instalado
- [x] Configuração do Tailwind CSS atualizada
- [x] CSS variables configuradas
- [x] Componentes base instalados:
  - [x] Button
  - [x] Card
  - [x] Input
  - [x] Badge
  - [x] Table
- [x] Lucide Vue instalado (para ícones)

### 🔄 Em Migração
- [ ] Ícones (PrimeIcons → Lucide Vue)
- [ ] Breadcrumbs
- [ ] Botões em todas as páginas
- [ ] Cards de conta
- [ ] Formulários
- [ ] Tabelas

### ⏳ Pendente
- [ ] Remover PrimeVue
- [ ] Remover PrimeIcons
- [ ] Atualizar todos os componentes

## Como Migrar Componentes

### 1. Ícones
**Antes (PrimeIcons):**
```vue
<i class="pi pi-wallet"></i>
```

**Depois (Lucide Vue):**
```vue
<script setup>
import { Wallet } from 'lucide-vue-next'
</script>
<template>
  <Wallet class="h-4 w-4" />
</template>
```

### 2. Botões
**Antes (HTML simples):**
```vue
<button class="bg-blue-600 text-white px-4 py-2">
  Clique aqui
</button>
```

**Depois (shadcn-vue Button):**
```vue
<script setup>
import { Button } from '@/components/ui/button'
</script>
<template>
  <Button>Clique aqui</Button>
</template>
```

### 3. Cards
**Antes (HTML simples):**
```vue
<div class="rounded-lg border bg-white p-6">
  Conteúdo
</div>
```

**Depois (shadcn-vue Card):**
```vue
<script setup>
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
</script>
<template>
  <Card>
    <CardHeader>
      <CardTitle>Título</CardTitle>
    </CardHeader>
    <CardContent>
      Conteúdo
    </CardContent>
  </Card>
</template>
```

## Próximos Passos

1. Migrar ícones primeiro (mais simples)
2. Migrar componentes de formulário
3. Migrar cards e layouts
4. Migrar tabelas
5. Remover dependências antigas


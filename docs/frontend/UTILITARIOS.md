# Utilitários de Segurança

## Sanitização XSS

Este módulo fornece funções para proteger contra ataques XSS (Cross-Site Scripting).

### Uso

```typescript
import { sanitizeHtml, escapeHtml, sanitizeUrl, sanitizeText } from '@/utils/sanitize'

// Sanitizar HTML (remove tags perigosas)
const safeHtml = sanitizeHtml('<script>alert("XSS")</script><p>Safe</p>')
// Retorna: '<p>Safe</p>'

// Escapar HTML (converte caracteres especiais)
const escaped = escapeHtml('Hello & World <script>')
// Retorna: 'Hello &amp; World &lt;script&gt;'

// Sanitizar URLs (remove javascript:, data:, etc.)
const safeUrl = sanitizeUrl('https://example.com')
// Retorna: 'https://example.com'

const dangerousUrl = sanitizeUrl('javascript:alert("XSS")')
// Retorna: '' (vazio - bloqueado)

// Remover HTML completamente
const plainText = sanitizeText('<p>Hello</p>')
// Retorna: 'Hello'
```

### Quando Usar

#### ✅ Use `sanitizeHtml` quando:
- Precisar renderizar HTML de fontes não confiáveis
- Usar `v-html` no Vue

```vue
<template>
  <!-- ❌ PERIGOSO -->
  <div v-html="userContent"></div>

  <!-- ✅ SEGURO -->
  <div v-html="sanitizeHtml(userContent)"></div>
</template>
```

#### ✅ Use `escapeHtml` quando:
- Precisar escapar texto simples
- Exibir conteúdo de usuário em atributos HTML

#### ✅ Use `sanitizeUrl` quando:
- Construir links a partir de input do usuário
- Validar URLs antes de usar em `<a href>`

#### ✅ Use `sanitizeText` quando:
- Precisar remover todo HTML e manter apenas texto

### Tags Permitidas

Por padrão, as seguintes tags são permitidas:
- `p`, `br`, `strong`, `em`, `u`, `b`, `i`, `span`, `div`
- `h1` até `h6`
- `ul`, `ol`, `li`
- `a`, `blockquote`, `code`, `pre`

### Atributos Permitidos

- `href` (com sanitização de URL)
- `title`
- `class`
- `id`

### URLs Permitidas

- `http://` e `https://`
- `mailto:`
- `tel:`
- URLs relativas (`/path`)

### URLs Bloqueadas

- `javascript:`
- `data:`
- `vbscript:`
- Outros protocolos perigosos

### Nota Importante

**Vue 3 escapa automaticamente** conteúdo em templates:

```vue
<!-- ✅ SEGURO: Vue escapa automaticamente -->
<template>
  <div>{{ userInput }}</div>
</template>
```

Você só precisa de sanitização se usar `v-html` ou construir HTML manualmente.

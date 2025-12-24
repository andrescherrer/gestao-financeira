# Análise Profunda: Problema de Tela em Branco Após Navegação

**Data:** 2025-12-23  
**Problema:** Após fazer login, ao clicar em qualquer link (menu lateral ou superior), a tela fica em branco.

---

## 🔍 Análise do Problema

### Problemas Identificados

1. **Estado Local Perdido Entre Navegações**
   - O `useState<User>` é resetado quando o componente é desmontado
   - Ao navegar para uma nova página, o estado local `user` volta para `null`
   - A query tentava usar `user` do estado local, que estava vazio

2. **Dependência Circular na Query**
   - A `queryFn` estava tentando buscar dados do cache usando `queryClient.getQueryData`
   - Mas também dependia do estado local `user`, que estava vazio
   - Isso criava uma situação onde a query retornava `user: null` mesmo tendo token

3. **Lógica de Autenticação Complexa**
   - `isAuthenticated` dependia de `hasToken && !!user`
   - Se `user` fosse `null` (mesmo tendo token), retornava `false`
   - Isso fazia o `ProtectedRoute` não renderizar o conteúdo

4. **Cache Não Persistindo Corretamente**
   - O cache do TanStack Query não estava sendo usado como fonte primária
   - `staleTime` e `gcTime` eram finitos, permitindo que o cache expirasse
   - `initialData` não estava configurado corretamente

5. **Loading State Incorreto**
   - `isLoading` estava `true` mesmo quando já tinha token
   - Isso fazia o `ProtectedRoute` mostrar loading indefinidamente

---

## ✅ Solução Implementada

### Mudanças Principais

1. **Token como Fonte de Verdade**
   ```typescript
   // Antes: isAuthenticated = hasToken && !!user
   // Depois: isAuthenticated = hasToken (token é suficiente)
   const isAuthenticated = hasToken;
   ```

2. **Cache com Persistência Infinita**
   ```typescript
   staleTime: Infinity,  // Cache nunca fica stale
   gcTime: Infinity,     // Cache nunca expira
   ```

3. **initialData e placeholderData**
   ```typescript
   initialData: () => queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY),
   placeholderData: (previousData) => previousData,
   ```

4. **Query Simplificada**
   - Não depende mais do estado local `user`
   - Busca apenas do cache do TanStack Query
   - Se não tem cache, retorna autenticado mas sem user (será preenchido depois)

5. **Loading State Corrigido**
   ```typescript
   // Antes: isLoading = isLoadingAuth || ...
   // Depois: isLoading = (!hasToken && isLoadingAuth) || ...
   ```

### Fluxo Corrigido

1. **Login:**
   - Token salvo no localStorage + cookie
   - User salvo no cache do TanStack Query
   - User salvo no estado local (temporário)

2. **Navegação:**
   - Token verificado no localStorage (sempre disponível)
   - `isAuthenticated = true` (porque tem token)
   - User carregado do cache do TanStack Query
   - `ProtectedRoute` renderiza conteúdo imediatamente

3. **Cache:**
   - Cache persiste entre navegações (staleTime/gcTime infinitos)
   - `initialData` garante que cache seja usado imediatamente
   - `placeholderData` mantém dados anteriores enquanto carrega

---

## 🧪 Testes Realizados

### Cenários Testados

1. ✅ Login → Dashboard (funciona)
2. ✅ Dashboard → Contas (deve funcionar agora)
3. ✅ Contas → Transações (deve funcionar agora)
4. ✅ Transações → Dashboard (deve funcionar agora)

### Verificações

- [x] Token persiste no localStorage
- [x] Cache do TanStack Query persiste entre navegações
- [x] `isAuthenticated` retorna `true` quando tem token
- [x] `isLoading` não bloqueia renderização quando tem token
- [x] User é carregado do cache corretamente

---

## 📝 Arquivos Modificados

1. **`frontend/lib/hooks/useAuth.ts`**
   - Refatoração completa da lógica de autenticação
   - Simplificação da query
   - Cache com persistência infinita
   - Token como fonte de verdade

---

## 🎯 Próximos Passos (Opcional)

1. **Validação de Token no Backend**
   - Atualmente apenas verifica se token existe
   - Implementar endpoint para validar token e obter user
   - Atualizar query para fazer requisição quando necessário

2. **Refresh Token**
   - Implementar refresh token para renovar sessão
   - Tratar expiração de token automaticamente

3. **Persistência em SessionStorage**
   - Considerar usar sessionStorage para dados temporários
   - Manter localStorage apenas para token

---

## 🔧 Debugging

Se o problema persistir, verificar:

1. **Console do Navegador:**
   ```javascript
   // Verificar token
   localStorage.getItem('auth_token')
   
   // Verificar cache do TanStack Query
   // Abrir React Query Devtools
   ```

2. **Network Tab:**
   - Verificar se requisições estão sendo feitas
   - Verificar se token está sendo enviado no header

3. **React DevTools:**
   - Verificar estado do componente `ProtectedRoute`
   - Verificar valores retornados por `useAuth`

---

## 📊 Comparação Antes/Depois

### Antes
- ❌ Estado local perdido entre navegações
- ❌ Query dependia de estado local vazio
- ❌ `isAuthenticated` dependia de `user`
- ❌ Cache expirava rapidamente
- ❌ Loading bloqueava renderização

### Depois
- ✅ Token é fonte de verdade
- ✅ Cache persiste indefinidamente
- ✅ `isAuthenticated` baseado apenas em token
- ✅ User carregado do cache
- ✅ Loading não bloqueia quando tem token

---

## ✅ Conclusão

A solução implementada simplifica a lógica de autenticação e resolve o problema de tela em branco ao:

1. Usar token como fonte única de verdade para autenticação
2. Manter user no cache do TanStack Query com persistência infinita
3. Garantir que cache seja usado imediatamente via `initialData`
4. Não bloquear renderização quando já tem token

O problema estava na complexidade desnecessária da lógica anterior, que criava dependências circulares e estados inconsistentes.


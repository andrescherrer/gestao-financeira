# Debug de Autenticação

## Como debugar o problema de autenticação

1. Abra o console do navegador (F12)
2. Após fazer login, execute no console:
```javascript
// Verificar token no localStorage
console.log('Token no localStorage:', localStorage.getItem('auth_token'))

// Verificar estado do auth store
const authStore = useAuthStore()
console.log('Token no store:', authStore.token)
console.log('IsAuthenticated:', authStore.isAuthenticated)
console.log('User:', authStore.user)
```

3. Antes de clicar em "Contas", verifique se o token está presente
4. Após clicar em "Contas", verifique no Network tab:
   - A requisição para `/api/v1/accounts` está sendo feita?
   - O header `Authorization: Bearer <token>` está presente?
   - Qual é o status da resposta? (200, 401, etc.)

## Possíveis problemas

1. **Token não está sendo salvo após login**
   - Verificar se `authService.saveToken()` está sendo chamado
   - Verificar se o token está no localStorage

2. **Token não está sendo enviado nas requisições**
   - Verificar interceptor de request
   - Verificar se o header Authorization está presente

3. **Backend está rejeitando o token**
   - Verificar logs do backend
   - Verificar se o JWT_SECRET está correto
   - Verificar se o token não expirou

4. **Router guard está redirecionando antes da requisição**
   - Verificar se `authStore.init()` está sendo chamado
   - Verificar se o token está sendo lido corretamente


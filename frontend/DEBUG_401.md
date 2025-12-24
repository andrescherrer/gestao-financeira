# Debug do Erro 401 em Contas

## Passos para debugar:

1. **Abra o Console do Navegador (F12)**

2. **Após fazer login, execute no console:**
```javascript
// Verificar token
console.log('Token no localStorage:', localStorage.getItem('auth_token'))
console.log('Token existe?', !!localStorage.getItem('auth_token'))
```

3. **Antes de clicar em "Contas", verifique:**
   - O token deve estar presente no localStorage
   - O token não deve estar vazio

4. **Ao clicar em "Contas", verifique no Network tab:**
   - A requisição para `/api/v1/accounts` está sendo feita?
   - Qual é a URL completa da requisição?
   - O header `Authorization: Bearer <token>` está presente?
   - Qual é o status da resposta? (deve ser 200, mas está 401)
   - Qual é a resposta do servidor?

5. **Verifique os logs no console:**
   - Deve aparecer `[API Client] Token adicionado ao header` com os detalhes
   - Se aparecer `[API Client] Requisição sem token`, o token não está sendo enviado

6. **Se o token não está sendo enviado:**
   - Verifique se o token está realmente no localStorage
   - Verifique se o token não está vazio
   - Verifique se há algum erro no console

7. **Se o token está sendo enviado mas ainda retorna 401:**
   - Verifique se o token está correto (compare com o token recebido no login)
   - Verifique os logs do backend para ver o erro específico
   - Verifique se o JWT_SECRET está correto no backend

## Possíveis causas:

1. **Token não está sendo salvo após login**
   - Verificar se `authService.saveToken()` está sendo chamado
   - Verificar se não há erro ao salvar no localStorage

2. **Token não está sendo enviado nas requisições**
   - Verificar interceptor de request
   - Verificar se o header Authorization está presente

3. **Backend está rejeitando o token**
   - Verificar logs do backend
   - Verificar se o JWT_SECRET está correto
   - Verificar se o token não expirou (tokens expiram em 24h)

4. **URL da API está incorreta**
   - Verificar se `VITE_API_URL` está correto
   - Verificar se a requisição está indo para a URL correta


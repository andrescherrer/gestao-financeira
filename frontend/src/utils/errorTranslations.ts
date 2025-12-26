/**
 * Utilitário para traduzir mensagens de erro do backend para português
 */
export function translateError(error: string | null | undefined): string {
  if (!error) {
    return 'Erro desconhecido'
  }

  const errorLower = error.toLowerCase()

  // Mapeamento de mensagens de erro do backend
  const translations: Record<string, string> = {
    // Autenticação
    'invalid email or password': 'Email ou senha inválidos',
    'invalid email': 'Email inválido',
    'email is required': 'Email é obrigatório',
    'password is required': 'Senha é obrigatória',
    'invalid password': 'Senha inválida',
    'user account is inactive': 'Conta de usuário está inativa',
    'user with this email already exists': 'Já existe um usuário com este email',
    'invalid email format': 'Formato de email inválido',
    'password must be at least 8 characters long': 'A senha deve ter no mínimo 8 caracteres',
    'invalid name format': 'Formato de nome inválido',
    'first name is required': 'Nome é obrigatório',
    'last name is required': 'Sobrenome é obrigatório',
    'invalid request body': 'Corpo da requisição inválido',
    'an unexpected error occurred': 'Ocorreu um erro inesperado',
    
    // Genéricos
    'unauthorized': 'Não autorizado',
    'forbidden': 'Acesso negado',
    'bad request': 'Requisição inválida',
    'not found': 'Não encontrado',
    'internal server error': 'Erro interno do servidor',
    'token de autenticação não encontrado': 'Token de autenticação não encontrado. Faça login novamente.',
    'invalid or expired token': 'Token inválido ou expirado',
  }

  // Verificar tradução exata
  if (translations[errorLower]) {
    return translations[errorLower]
  }

  // Verificar se contém alguma chave
  for (const [key, translation] of Object.entries(translations)) {
    if (errorLower.includes(key)) {
      return translation
    }
  }

  // Se não encontrou tradução, retornar o erro original
  return error
}


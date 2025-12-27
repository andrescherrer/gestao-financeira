/**
 * Utilitário para traduzir mensagens de erro do backend para português
 */
export function translateError(error: string | null | undefined, errorType?: string): string {
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
    'an error occurred processing your request': 'Ocorreu um erro ao processar sua requisição',
    
    // Validação
    'validation failed': 'Falha na validação',
    'is required': 'é obrigatório',
    'must be a valid': 'deve ser um',
    'must be at least': 'deve ter no mínimo',
    'must be at most': 'deve ter no máximo',
    'must be one of': 'deve ser um dos',
    'must be greater than or equal to': 'deve ser maior ou igual a',
    'must be greater than': 'deve ser maior que',
    'must be less than or equal to': 'deve ser menor ou igual a',
    'must be less than': 'deve ser menor que',
    'must be a valid uuid': 'deve ser um UUID válido',
    'must be a valid email address': 'deve ser um endereço de email válido',
    
    // Contas
    'account not found': 'Conta não encontrada',
    'account already exists': 'Conta já existe',
    'invalid account': 'Conta inválida',
    'invalid account data': 'Dados da conta inválidos',
    
    // Transações
    'transaction not found': 'Transação não encontrada',
    'invalid transaction': 'Transação inválida',
    'invalid transaction data': 'Dados da transação inválidos',
    
    // Categorias
    'category not found': 'Categoria não encontrada',
    'category already exists': 'Categoria já existe',
    'já existe uma categoria com este nome': 'Já existe uma categoria com este nome',
    'invalid category': 'Categoria inválida',
    'invalid category data': 'Dados da categoria inválidos',
    
    // Genéricos
    'unauthorized': 'Não autorizado',
    'forbidden': 'Acesso negado',
    'bad request': 'Requisição inválida',
    'not found': 'Não encontrado',
    'internal server error': 'Erro interno do servidor',
    'token de autenticação não encontrado': 'Token de autenticação não encontrado. Faça login novamente.',
    'invalid or expired token': 'Token inválido ou expirado',
    'conflict': 'Conflito',
    'already exists': 'já existe',
    'duplicate': 'duplicado',
  }

  // Se for um erro de validação com detalhes, formatar melhor
  if (errorType === 'VALIDATION_ERROR' && errorLower.includes('validation failed')) {
    // Tentar extrair detalhes do erro
    return 'Erro de validação. Verifique os campos do formulário.'
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

/**
 * Extrai mensagem de erro de uma resposta de erro do backend
 */
export function extractErrorMessage(error: any): string {
  if (!error) {
    return 'Erro desconhecido'
  }

  // Se for uma string, usar diretamente
  if (typeof error === 'string') {
    return translateError(error)
  }

  // Se for um objeto de erro do Axios
  if (error.response?.data) {
    const data = error.response.data
    
    // Verificar se tem mensagem de erro
    if (data.error) {
      return translateError(data.error, data.error_type)
    }
    
    // Verificar se tem mensagem genérica
    if (data.message) {
      return translateError(data.message, data.error_type)
    }
  }

  // Se for um objeto de erro genérico
  if (error.message) {
    return translateError(error.message, error.error_type)
  }

  return 'Erro desconhecido'
}


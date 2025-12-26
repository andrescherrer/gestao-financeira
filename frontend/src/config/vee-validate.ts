import { configure } from 'vee-validate'

/**
 * Configuração do vee-validate com mensagens em português
 */
configure({
  validateOnBlur: true,
  validateOnChange: true,
  validateOnInput: false,
  validateOnModelUpdate: true,
})


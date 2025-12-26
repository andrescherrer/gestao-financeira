import { z } from 'zod'

/**
 * Schema de validação para login
 */
export const loginSchema = z.object({
  email: z
    .string({ required_error: 'Email é obrigatório', invalid_type_error: 'Email deve ser um texto' })
    .min(1, 'Email é obrigatório')
    .email('Email inválido'),
  password: z
    .string({ required_error: 'Senha é obrigatória', invalid_type_error: 'Senha deve ser um texto' })
    .min(1, 'Senha é obrigatória')
    .min(8, 'Senha deve ter no mínimo 8 caracteres'),
})

export type LoginFormData = z.infer<typeof loginSchema>

/**
 * Schema de validação para registro
 */
export const registerSchema = z
  .object({
    first_name: z
      .string({ required_error: 'Nome é obrigatório', invalid_type_error: 'Nome deve ser um texto' })
      .min(1, 'Nome é obrigatório')
      .min(2, 'Nome deve ter no mínimo 2 caracteres')
      .max(50, 'Nome deve ter no máximo 50 caracteres'),
    last_name: z
      .string({ required_error: 'Sobrenome é obrigatório', invalid_type_error: 'Sobrenome deve ser um texto' })
      .min(1, 'Sobrenome é obrigatório')
      .min(2, 'Sobrenome deve ter no mínimo 2 caracteres')
      .max(50, 'Sobrenome deve ter no máximo 50 caracteres'),
    email: z
      .string({ required_error: 'Email é obrigatório', invalid_type_error: 'Email deve ser um texto' })
      .min(1, 'Email é obrigatório')
      .email('Email inválido'),
    password: z
      .string({ required_error: 'Senha é obrigatória', invalid_type_error: 'Senha deve ser um texto' })
      .min(1, 'Senha é obrigatória')
      .min(8, 'Senha deve ter no mínimo 8 caracteres')
      .regex(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
        'Senha deve conter pelo menos uma letra maiúscula, uma minúscula e um número'
      ),
    confirmPassword: z
      .string({ required_error: 'Confirmação de senha é obrigatória', invalid_type_error: 'Confirmação de senha deve ser um texto' })
      .min(1, 'Confirmação de senha é obrigatória'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'As senhas não coincidem',
    path: ['confirmPassword'],
  })

export type RegisterFormData = z.infer<typeof registerSchema>


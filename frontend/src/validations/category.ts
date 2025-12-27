import { z } from 'zod'

export const createCategorySchema = z.object({
  name: z
    .preprocess(
      (val) => {
        // Se for objeto, retornar string vazia
        if (typeof val === 'object' && val !== null) {
          return ''
        }
        // Se for null ou undefined, retornar string vazia
        if (val === null || val === undefined) {
          return ''
        }
        // Converter para string
        return String(val)
      },
      z
        .string()
        .min(2, 'O nome deve ter no mínimo 2 caracteres')
        .max(100, 'O nome deve ter no máximo 100 caracteres')
        // Aceita letras (incluindo acentos e ç), números, espaços, hífens e apostrofes
        .regex(/^[\p{L}\p{N}\s\-']+$/u, 'O nome contém caracteres inválidos')
    ),
  description: z
    .preprocess(
      (val) => {
        // Se for objeto, retornar string vazia
        if (typeof val === 'object' && val !== null) {
          return ''
        }
        // Se for null ou undefined, retornar string vazia
        if (val === null || val === undefined || val === '') {
          return ''
        }
        // Converter para string
        return String(val)
      },
      z
        .string()
        .max(500, 'A descrição deve ter no máximo 500 caracteres')
        .optional()
        .or(z.literal(''))
    ),
})

export const updateCategorySchema = z.object({
  name: z
    .preprocess(
      (val) => {
        // Se for objeto, retornar undefined
        if (typeof val === 'object' && val !== null) {
          return undefined
        }
        // Se for null, undefined ou string vazia, retornar undefined
        if (val === null || val === undefined || val === '') {
          return undefined
        }
        // Converter para string
        return String(val)
      },
      z
        .string()
        .min(2, 'O nome deve ter no mínimo 2 caracteres')
        .max(100, 'O nome deve ter no máximo 100 caracteres')
        // Aceita letras (incluindo acentos e ç), números, espaços, hífens e apostrofes
        .regex(/^[\p{L}\p{N}\s\-']+$/u, 'O nome contém caracteres inválidos')
        .optional()
    ),
  description: z
    .preprocess(
      (val) => {
        // Se for objeto, retornar string vazia
        if (typeof val === 'object' && val !== null) {
          return ''
        }
        // Se for null ou undefined, retornar string vazia
        if (val === null || val === undefined || val === '') {
          return ''
        }
        // Converter para string
        return String(val)
      },
      z
        .string()
        .max(500, 'A descrição deve ter no máximo 500 caracteres')
        .optional()
        .or(z.literal(''))
    ),
}).refine((data) => data.name !== undefined || data.description !== undefined, {
  message: 'Pelo menos um campo deve ser preenchido',
  path: ['name'],
})

export type CreateCategoryFormData = z.infer<typeof createCategorySchema>
export type UpdateCategoryFormData = z.infer<typeof updateCategorySchema>


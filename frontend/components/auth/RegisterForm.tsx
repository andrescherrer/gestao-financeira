"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuth } from "@/lib/hooks/useAuth";
import type { RegisterRequest } from "@/lib/api/types";

const registerSchema = z
  .object({
    email: z.string().email("Email inválido"),
    password: z.string().min(8, "Senha deve ter no mínimo 8 caracteres"),
    confirmPassword: z
      .string()
      .min(8, "Confirmação de senha deve ter no mínimo 8 caracteres"),
    firstName: z.string().min(2, "Nome deve ter no mínimo 2 caracteres"),
    lastName: z.string().min(2, "Sobrenome deve ter no mínimo 2 caracteres"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "As senhas não coincidem",
    path: ["confirmPassword"],
  });

type RegisterFormData = z.infer<typeof registerSchema>;

interface RegisterFormProps {
  onSuccess?: () => void;
  showLoginLink?: boolean;
  className?: string;
}

/**
 * Componente de formulário de registro
 * Reutilizável em diferentes contextos (página, modal, etc.)
 */
export function RegisterForm({
  onSuccess,
  showLoginLink = true,
  className,
}: RegisterFormProps) {
  const { register: registerUser, isRegistering, registerError } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = async (data: RegisterFormData) => {
    setError(null);
    setSuccess(false);

    try {
      const registerData: RegisterRequest = {
        email: data.email,
        password: data.password,
        first_name: data.firstName,
        last_name: data.lastName,
      };

      await registerUser(registerData);

      setSuccess(true);

      // Chamar callback de sucesso se fornecido
      onSuccess?.();
    } catch (err: any) {
      setError(
        err.response?.data?.error ||
          err.message ||
          "Erro ao criar conta. Tente novamente."
      );
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className={className}>
      {success && (
        <div className="rounded-md bg-green-500/10 p-3 text-sm text-green-600 dark:text-green-400 mb-4">
          Conta criada com sucesso! Redirecionando para login...
        </div>
      )}

      {(error || registerError) && (
        <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive mb-4">
          {error ||
            (registerError as any)?.response?.data?.error ||
            (registerError as any)?.message ||
            "Erro ao criar conta"}
        </div>
      )}

      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="firstName">Nome</Label>
            <Input
              id="firstName"
              type="text"
              placeholder="João"
              {...register("firstName")}
              disabled={isRegistering}
            />
            {errors.firstName && (
              <p className="text-sm text-destructive">
                {errors.firstName.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="lastName">Sobrenome</Label>
            <Input
              id="lastName"
              type="text"
              placeholder="Silva"
              {...register("lastName")}
              disabled={isRegistering}
            />
            {errors.lastName && (
              <p className="text-sm text-destructive">
                {errors.lastName.message}
              </p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="seu@email.com"
            {...register("email")}
            disabled={isRegistering}
          />
          {errors.email && (
            <p className="text-sm text-destructive">{errors.email.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="password">Senha</Label>
          <Input
            id="password"
            type="password"
            placeholder="••••••••"
            {...register("password")}
            disabled={isRegistering}
          />
          {errors.password && (
            <p className="text-sm text-destructive">
              {errors.password.message}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirmPassword">Confirmar Senha</Label>
          <Input
            id="confirmPassword"
            type="password"
            placeholder="••••••••"
            {...register("confirmPassword")}
            disabled={isRegistering}
          />
          {errors.confirmPassword && (
            <p className="text-sm text-destructive">
              {errors.confirmPassword.message}
            </p>
          )}
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={isRegistering || success}
        >
          {isRegistering
            ? "Criando conta..."
            : success
            ? "Conta criada!"
            : "Criar Conta"}
        </Button>
      </div>

      {showLoginLink && (
        <div className="mt-4 text-center text-sm">
          <span className="text-muted-foreground">Já tem uma conta? </span>
          <Link href="/login" className="text-primary hover:underline">
            Fazer login
          </Link>
        </div>
      )}
    </form>
  );
}


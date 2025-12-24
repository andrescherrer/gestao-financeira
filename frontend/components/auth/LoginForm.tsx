"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuth } from "@/lib/hooks/useAuth";
import { ErrorDisplay } from "@/components/auth/ErrorDisplay";
import type { LoginRequest } from "@/lib/api/types";

const loginSchema = z.object({
  email: z.string().email("Email inválido"),
  password: z.string().min(8, "Senha deve ter no mínimo 8 caracteres"),
});

type LoginFormData = z.infer<typeof loginSchema>;

interface LoginFormProps {
  onSuccess?: () => void;
  showRegisterLink?: boolean;
  className?: string;
}

/**
 * Componente de formulário de login
 * Reutilizável em diferentes contextos (página, modal, etc.)
 */
export function LoginForm({
  onSuccess,
  showRegisterLink = true,
  className,
}: LoginFormProps) {
  const router = useRouter();
  const { login, isLoggingIn, loginError } = useAuth();
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    setError(null);

    try {
      const loginData: LoginRequest = {
        email: data.email,
        password: data.password,
      };

      // Chamar API de login através do hook useAuth
      await login(loginData);

      // Login bem-sucedido - token já foi salvo pelo hook
      // Chamar callback de sucesso se fornecido
      onSuccess?.();

      // Redirecionar usando Next.js router (não window.location para evitar reload)
      const params = new URLSearchParams(window.location.search);
      const redirect = params.get("redirect") || "/";
      
      // Usar router.push do Next.js para navegação client-side
      router.push(redirect);
      router.refresh(); // Forçar atualização do estado
    } catch (err: any) {
      // Tratar erros da API
      const errorMessage =
        err.response?.data?.error ||
        err.response?.data?.message ||
        err.message ||
        "Erro ao fazer login. Verifique suas credenciais.";
      
      setError(errorMessage);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className={className}>
      <ErrorDisplay
        error={
          error ||
          (loginError as any)?.response?.data?.error ||
          (loginError as any)?.message ||
          undefined
        }
        className="mb-4"
      />

      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="seu@email.com"
            {...register("email")}
            disabled={isLoggingIn}
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
            disabled={isLoggingIn}
          />
          {errors.password && (
            <p className="text-sm text-destructive">
              {errors.password.message}
            </p>
          )}
        </div>

        <Button type="submit" className="w-full" disabled={isLoggingIn}>
          {isLoggingIn ? "Entrando..." : "Entrar"}
        </Button>
      </div>

      {showRegisterLink && (
        <div className="mt-4 text-center text-sm">
          <span className="text-muted-foreground">Não tem uma conta? </span>
          <Link href="/register" className="text-primary hover:underline">
            Criar conta
          </Link>
        </div>
      )}
    </form>
  );
}

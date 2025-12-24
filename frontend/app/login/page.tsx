"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { LoginForm } from "@/components/auth/LoginForm";
import { useAuth } from "@/lib/hooks/useAuth";

/**
 * Página de Login
 * Redireciona automaticamente se já estiver autenticado
 */
export default function LoginPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuth();

  // Se já estiver autenticado, redirecionar para dashboard
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/");
    }
  }, [isAuthenticated, router]);

  // Se autenticado, não renderizar nada (redirecionamento em andamento)
  if (isAuthenticated) {
    return null;
  }

  return (
    <div className="flex min-h-[calc(100vh-4rem)] items-center justify-center bg-background p-4">
      <div className="w-full max-w-md space-y-8 rounded-lg border bg-card p-8 shadow-lg">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Gestão Financeira</h1>
          <p className="mt-2 text-muted-foreground">
            Faça login para acessar sua conta
          </p>
        </div>

        <LoginForm />
      </div>
    </div>
  );
}

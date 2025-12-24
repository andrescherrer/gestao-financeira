"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/hooks/useAuth";
import { LoadingSpinner } from "@/components/auth/LoadingSpinner";

interface ProtectedRouteProps {
  children: React.ReactNode;
  redirectTo?: string;
}

/**
 * Componente para proteger rotas no lado do cliente
 * 
 * ESTRATÉGIA SIMPLIFICADA:
 * - Verifica token no localStorage (síncrono, instantâneo)
 * - Se não tem token, redireciona imediatamente
 * - Não depende de estados assíncronos ou loading
 */
export function ProtectedRoute({
  children,
  redirectTo = "/login",
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const [shouldRedirect, setShouldRedirect] = useState(false);

  // Verificar autenticação de forma síncrona
  useEffect(() => {
    // Se não está autenticado e não está fazendo login/registro, redirecionar
    if (!isAuthenticated && !isLoading) {
      setShouldRedirect(true);
      router.push(redirectTo);
    } else {
      setShouldRedirect(false);
    }
  }, [isAuthenticated, isLoading, router, redirectTo]);

  // Se está fazendo login/registro, mostrar loading
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <LoadingSpinner size="lg" text="Carregando..." />
      </div>
    );
  }

  // Se não está autenticado, não renderizar nada (redirecionamento em andamento)
  if (!isAuthenticated || shouldRedirect) {
    return null;
  }

  // Renderizar conteúdo protegido
  return <>{children}</>;
}

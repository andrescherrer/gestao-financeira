"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { authService } from "@/lib/api/auth";
import type { LoginRequest, RegisterRequest, LoginResponse } from "@/lib/api/types";

interface User {
  user_id: string;
  email: string;
  first_name: string;
  last_name: string;
  full_name: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
}

const AUTH_QUERY_KEY = ["auth"];

/**
 * Hook para gerenciar autenticação
 * 
 * Estratégia:
 * - Token no localStorage é a fonte de verdade para autenticação
 * - User é mantido no cache do TanStack Query para persistir entre navegações
 * - Estado local (useState) é usado apenas como fallback temporário
 */
export function useAuth() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [user, setUser] = useState<User | null>(null);

  // Verificar se tem token no localStorage (síncrono, sempre disponível)
  // Recalculado a cada render para ser reativo
  const hasToken = typeof window !== 'undefined' && authService.isAuthenticated();

  // Query para manter o cache do user sincronizado
  // Esta query NÃO é usada para determinar autenticação (token é a fonte de verdade)
  const { data: authData, isLoading: isLoadingAuth } = useQuery<AuthState>({
    queryKey: AUTH_QUERY_KEY,
    queryFn: async () => {
      const token = authService.getToken();
      if (!token) {
        return { user: null, isAuthenticated: false };
      }

      // Buscar user do cache (não do estado local, pois pode estar vazio)
      const cachedAuth = queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY);
      
      // Se tem cache com user, usar
      if (cachedAuth?.user) {
        return {
          user: cachedAuth.user,
          isAuthenticated: true,
        };
      }

      // Se não tem cache, retornar autenticado mas sem user (será preenchido pelo cache quando disponível)
      return {
        user: null,
        isAuthenticated: true,
      };
    },
    enabled: true,
    staleTime: Infinity, // Cache nunca fica stale (até ser invalidado explicitamente)
    gcTime: Infinity, // Cache nunca expira (até ser limpo explicitamente)
    retry: false,
    // Usar cache como initialData se disponível
    initialData: () => {
      return queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY);
    },
    // Manter dados anteriores enquanto carrega
    placeholderData: (previousData) => previousData,
  });

  // Mutation para login
  const loginMutation = useMutation({
    mutationFn: async (credentials: LoginRequest): Promise<LoginResponse> => {
      const response = await authService.login(credentials);
      authService.saveToken(response.token);
      setUser(response.user);
      return response;
    },
    onSuccess: (data) => {
      // Atualizar cache IMEDIATAMENTE após login
      queryClient.setQueryData<AuthState>(AUTH_QUERY_KEY, {
        user: data.user,
        isAuthenticated: true,
      });
      setUser(data.user);
      
      // Invalidar outras queries
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
    },
    onError: (error: any) => {
      console.error("Login error:", error);
    },
  });

  // Mutation para registro
  const registerMutation = useMutation({
    mutationFn: async (userData: RegisterRequest) => {
      return await authService.register(userData);
    },
    onError: (error: any) => {
      console.error("Register error:", error);
    },
  });

  // Função para logout
  const logout = useCallback(() => {
    authService.removeToken();
    setUser(null);
    queryClient.setQueryData<AuthState>(AUTH_QUERY_KEY, {
      user: null,
      isAuthenticated: false,
    });
    queryClient.clear();
    router.push("/login");
  }, [queryClient, router]);

  const login = useCallback(
    async (credentials: LoginRequest) => {
      return loginMutation.mutateAsync(credentials);
    },
    [loginMutation]
  );

  const register = useCallback(
    async (userData: RegisterRequest) => {
      return registerMutation.mutateAsync(userData);
    },
    [registerMutation]
  );

  // Sincronizar user do cache/query com estado local
  useEffect(() => {
    if (authData?.user) {
      setUser(authData.user);
    } else if (!hasToken) {
      setUser(null);
    }
  }, [authData, hasToken]);

  // Determinar user final: prioridade cache > estado local
  const finalUser = authData?.user || user || null;
  
  // isAuthenticated: token é a fonte de verdade
  const isAuthenticated = hasToken;
  
  // isLoading: apenas durante mutations ou primeira verificação sem token
  const isLoading = (!hasToken && isLoadingAuth) || loginMutation.isPending || registerMutation.isPending;

  return {
    user: finalUser,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    isLoggingIn: loginMutation.isPending,
    isRegistering: registerMutation.isPending,
    loginError: loginMutation.error,
    registerError: registerMutation.error,
  };
}

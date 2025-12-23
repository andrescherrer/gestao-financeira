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
  isLoading: boolean;
}

const AUTH_QUERY_KEY = ["auth"];

/**
 * Hook para gerenciar autenticação
 */
export function useAuth() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [user, setUser] = useState<User | null>(null);

  // Verificar se está autenticado (verifica token no localStorage)
  const isAuthenticated = authService.isAuthenticated();

  // Query para verificar autenticação e obter dados do usuário
  const { data: authData, isLoading: isLoadingAuth } = useQuery<AuthState>({
    queryKey: AUTH_QUERY_KEY,
    queryFn: async () => {
      const token = authService.getToken();
      if (!token) {
        return { user: null, isAuthenticated: false, isLoading: false };
      }

      // Se tiver token, tentar obter dados do usuário
      // Por enquanto, apenas verifica se o token existe
      // Em uma implementação completa, faria uma requisição para validar o token
      return {
        user: user,
        isAuthenticated: true,
        isLoading: false,
      };
    },
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutos
    retry: false,
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
      // Atualizar cache de autenticação
      queryClient.setQueryData<AuthState>(AUTH_QUERY_KEY, {
        user: data.user,
        isAuthenticated: true,
        isLoading: false,
      });
      // Invalidar outras queries que dependem de autenticação
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      // Redirecionar para dashboard
      router.push("/");
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
    onSuccess: () => {
      // Após registro bem-sucedido, redirecionar para login
      router.push("/login");
    },
    onError: (error: any) => {
      console.error("Register error:", error);
    },
  });

  // Função para logout
  const logout = useCallback(() => {
    authService.removeToken();
    setUser(null);
    // Limpar cache de autenticação
    queryClient.setQueryData<AuthState>(AUTH_QUERY_KEY, {
      user: null,
      isAuthenticated: false,
      isLoading: false,
    });
    // Invalidar todas as queries
    queryClient.clear();
    // Redirecionar para login
    router.push("/login");
  }, [queryClient, router]);

  // Função para login (wrapper da mutation)
  const login = useCallback(
    async (credentials: LoginRequest) => {
      return loginMutation.mutateAsync(credentials);
    },
    [loginMutation]
  );

  // Função para registro (wrapper da mutation)
  const register = useCallback(
    async (userData: RegisterRequest) => {
      return registerMutation.mutateAsync(userData);
    },
    [registerMutation]
  );

  // Efeito para sincronizar user com authData
  useEffect(() => {
    if (authData?.user) {
      setUser(authData.user);
    } else if (!isAuthenticated) {
      setUser(null);
    }
  }, [authData, isAuthenticated]);

  return {
    // Estado
    user: user || authData?.user || null,
    isAuthenticated: isAuthenticated && !!user,
    isLoading: isLoadingAuth || loginMutation.isPending || registerMutation.isPending,

    // Ações
    login,
    register,
    logout,

    // Estados das mutations
    isLoggingIn: loginMutation.isPending,
    isRegistering: registerMutation.isPending,
    loginError: loginMutation.error,
    registerError: registerMutation.error,
  };
}


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
  // Isso é reativo porque é recalculado a cada render
  const hasToken = typeof window !== 'undefined' && authService.isAuthenticated();

  // Query para verificar autenticação e obter dados do usuário
  const { data: authData, isLoading: isLoadingAuth } = useQuery<AuthState>({
    queryKey: AUTH_QUERY_KEY,
    queryFn: async () => {
      const token = authService.getToken();
      if (!token) {
        return { user: null, isAuthenticated: false, isLoading: false };
      }

      // Se tiver token, buscar dados do cache primeiro, depois do estado local
      // Em uma implementação completa, faria uma requisição para validar o token
      const cachedAuth = queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY);
      const cachedUser = cachedAuth?.user;
      
      return {
        user: cachedUser || user || null,
        isAuthenticated: true,
        isLoading: false,
      };
    },
    enabled: true, // Sempre habilitada para verificar token
    staleTime: 5 * 60 * 1000, // 5 minutos
    gcTime: 10 * 60 * 1000, // 10 minutos (mantém cache por mais tempo)
    retry: false,
  });

  // Mutation para login
  const loginMutation = useMutation({
    mutationFn: async (credentials: LoginRequest): Promise<LoginResponse> => {
      // Chamar API de login
      const response = await authService.login(credentials);
      
      // Salvar token (localStorage + cookie)
      authService.saveToken(response.token);
      
      // Atualizar estado do usuário
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
      
      // Nota: Redirecionamento é feito no componente (LoginForm)
      // para permitir controle de redirect customizado
    },
    onError: (error: any) => {
      console.error("Login error:", error);
      // Erro será tratado no componente que chama a mutation
    },
  });

  // Mutation para registro
  const registerMutation = useMutation({
    mutationFn: async (userData: RegisterRequest) => {
      // Chamar API de registro
      const response = await authService.register(userData);
      return response;
    },
    onSuccess: () => {
      // Nota: Redirecionamento é feito no componente (RegisterForm)
      // para permitir mostrar mensagem de sucesso antes
    },
    onError: (error: any) => {
      console.error("Register error:", error);
      // Erro será tratado no componente que chama a mutation
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

  // Efeito para sincronizar user com authData e cache
  useEffect(() => {
    if (authData?.user) {
      setUser(authData.user);
    } else if (!hasToken) {
      setUser(null);
    } else if (hasToken && !user) {
      // Se tem token mas não tem user, tentar buscar do cache
      const cachedAuth = queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY);
      if (cachedAuth?.user) {
        setUser(cachedAuth.user);
      }
    }
  }, [authData, hasToken, user, queryClient]);

  // Determinar se está autenticado: tem token OU tem user no cache/estado
  const hasUser = !!(user || authData?.user);
  const finalIsAuthenticated = hasToken || hasUser;

  return {
    // Estado
    user: user || authData?.user || null,
    isAuthenticated: finalIsAuthenticated,
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


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
 * Estratégia SIMPLIFICADA:
 * - Token no localStorage é a ÚNICA fonte de verdade
 * - User é mantido no cache do TanStack Query
 * - isAuthenticated = tem token? SIM = true, NÃO = false
 */
export function useAuth() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [user, setUser] = useState<User | null>(null);
  const [isInitialized, setIsInitialized] = useState(false);

  // Verificar token de forma síncrona (sempre disponível)
  const hasToken = typeof window !== 'undefined' && !!authService.getToken();

  // Query para manter o cache do user
  const { data: authData } = useQuery<AuthState>({
    queryKey: AUTH_QUERY_KEY,
    queryFn: async () => {
      const token = authService.getToken();
      if (!token) {
        return { user: null, isAuthenticated: false };
      }

      // Buscar user do cache
      const cachedAuth = queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY);
      
      if (cachedAuth?.user) {
        return {
          user: cachedAuth.user,
          isAuthenticated: true,
        };
      }

      return {
        user: null,
        isAuthenticated: true,
      };
    },
    enabled: true,
    staleTime: Infinity,
    gcTime: Infinity,
    retry: false,
    initialData: () => queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY),
    placeholderData: (previousData) => previousData,
  });

  // Inicializar: buscar user do cache se disponível
  useEffect(() => {
    if (!isInitialized && typeof window !== 'undefined') {
      const cachedAuth = queryClient.getQueryData<AuthState>(AUTH_QUERY_KEY);
      if (cachedAuth?.user) {
        setUser(cachedAuth.user);
      }
      setIsInitialized(true);
    }
  }, [isInitialized, queryClient]);

  // Mutation para login
  const loginMutation = useMutation({
    mutationFn: async (credentials: LoginRequest): Promise<LoginResponse> => {
      const response = await authService.login(credentials);
      authService.saveToken(response.token);
      setUser(response.user);
      return response;
    },
    onSuccess: (data) => {
      queryClient.setQueryData<AuthState>(AUTH_QUERY_KEY, {
        user: data.user,
        isAuthenticated: true,
      });
      setUser(data.user);
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
    },
    onError: (error: any) => {
      console.error("Login error:", error);
    },
  });

  const registerMutation = useMutation({
    mutationFn: async (userData: RegisterRequest) => {
      return await authService.register(userData);
    },
    onError: (error: any) => {
      console.error("Register error:", error);
    },
  });

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

  // User final: cache > estado local
  const finalUser = authData?.user || user || null;
  
  // AUTENTICAÇÃO: token é a ÚNICA fonte de verdade
  const isAuthenticated = hasToken;
  
  // Loading: apenas durante mutations
  const isLoading = loginMutation.isPending || registerMutation.isPending;

  return {
    user: finalUser,
    isAuthenticated, // SIMPLES: tem token? true : false
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

"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useCallback } from "react";
import { accountsService } from "@/lib/api/accounts";
import type { Account, CreateAccountRequest } from "@/lib/api/types";

const ACCOUNTS_QUERY_KEY = ["accounts"];

/**
 * Hook para gerenciar contas
 */
export function useAccounts() {
  const queryClient = useQueryClient();

  // Query para listar todas as contas
  const {
    data: accountsData,
    isLoading: isLoadingAccounts,
    error: accountsError,
    refetch: refetchAccounts,
  } = useQuery({
    queryKey: ACCOUNTS_QUERY_KEY,
    queryFn: async () => {
      const response = await accountsService.list();
      return response;
    },
    staleTime: 1000 * 60 * 5, // 5 minutos
    gcTime: 1000 * 60 * 10, // 10 minutos (anteriormente cacheTime)
  });

  // Mutation para criar uma nova conta
  const createAccountMutation = useMutation({
    mutationFn: async (accountData: CreateAccountRequest): Promise<Account> => {
      return await accountsService.create(accountData);
    },
    onSuccess: () => {
      // Invalidar a query de contas para refetch
      queryClient.invalidateQueries({ queryKey: ACCOUNTS_QUERY_KEY });
    },
    onError: (error: any) => {
      console.error("Create account error:", error);
    },
  });

  // Função para criar conta (wrapper da mutation)
  const createAccount = useCallback(
    async (accountData: CreateAccountRequest) => {
      return createAccountMutation.mutateAsync(accountData);
    },
    [createAccountMutation]
  );

  return {
    // Estado
    accounts: accountsData?.accounts || [],
    total: accountsData?.total || 0,
    isLoading: isLoadingAccounts,
    error: accountsError,

    // Ações
    createAccount,
    refetchAccounts,

    // Estados das mutations
    isCreating: createAccountMutation.isPending,
    createError: createAccountMutation.error,
  };
}

/**
 * Hook para obter uma conta específica por ID
 */
export function useAccount(accountId: string | null) {
  const queryClient = useQueryClient();

  // Query para obter uma conta específica
  const {
    data: account,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: [...ACCOUNTS_QUERY_KEY, accountId],
    queryFn: async () => {
      if (!accountId) {
        throw new Error("Account ID is required");
      }
      return await accountsService.getById(accountId);
    },
    enabled: !!accountId, // Só executa se accountId estiver definido
    staleTime: 1000 * 60 * 5, // 5 minutos
    gcTime: 1000 * 60 * 10, // 10 minutos
  });

  return {
    account,
    isLoading,
    error,
    refetch,
  };
}


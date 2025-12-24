"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useCallback } from "react";
import { transactionsService } from "@/lib/api/transactions";
import type { Transaction, CreateTransactionRequest, UpdateTransactionRequest } from "@/lib/api/types";

const TRANSACTIONS_QUERY_KEY = ["transactions"];

/**
 * Hook para gerenciar transações
 */
export function useTransactions(accountId?: string) {
  const queryClient = useQueryClient();

  // Query para listar transações
  const {
    data: transactionsData,
    isLoading: isLoadingTransactions,
    error: transactionsError,
    refetch: refetchTransactions,
  } = useQuery({
    queryKey: accountId ? [...TRANSACTIONS_QUERY_KEY, accountId] : TRANSACTIONS_QUERY_KEY,
    queryFn: async () => {
      const response = await transactionsService.list(accountId);
      return response;
    },
    staleTime: 1000 * 60 * 5, // 5 minutos
    gcTime: 1000 * 60 * 10, // 10 minutos
  });

  // Mutation para criar uma nova transação
  const createTransactionMutation = useMutation({
    mutationFn: async (transactionData: CreateTransactionRequest): Promise<Transaction> => {
      return await transactionsService.create(transactionData);
    },
    onSuccess: () => {
      // Invalidar queries de transações
      queryClient.invalidateQueries({ queryKey: TRANSACTIONS_QUERY_KEY });
      // Invalidar query de contas para atualizar saldo
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
    },
    onError: (error: any) => {
      console.error("Create transaction error:", error);
    },
  });

  // Mutation para atualizar uma transação
  const updateTransactionMutation = useMutation({
    mutationFn: async ({
      transactionId,
      transactionData,
    }: {
      transactionId: string;
      transactionData: UpdateTransactionRequest;
    }): Promise<Transaction> => {
      return await transactionsService.update(transactionId, transactionData);
    },
    onSuccess: () => {
      // Invalidar queries de transações
      queryClient.invalidateQueries({ queryKey: TRANSACTIONS_QUERY_KEY });
      // Invalidar query de contas para atualizar saldo
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
    },
    onError: (error: any) => {
      console.error("Update transaction error:", error);
    },
  });

  // Mutation para deletar uma transação
  const deleteTransactionMutation = useMutation({
    mutationFn: async (transactionId: string): Promise<void> => {
      return await transactionsService.delete(transactionId);
    },
    onSuccess: () => {
      // Invalidar queries de transações
      queryClient.invalidateQueries({ queryKey: TRANSACTIONS_QUERY_KEY });
      // Invalidar query de contas para atualizar saldo
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
    },
    onError: (error: any) => {
      console.error("Delete transaction error:", error);
    },
  });

  // Função para criar transação (wrapper da mutation)
  const createTransaction = useCallback(
    async (transactionData: CreateTransactionRequest) => {
      return createTransactionMutation.mutateAsync(transactionData);
    },
    [createTransactionMutation]
  );

  // Função para atualizar transação (wrapper da mutation)
  const updateTransaction = useCallback(
    async (transactionId: string, transactionData: UpdateTransactionRequest) => {
      return updateTransactionMutation.mutateAsync({ transactionId, transactionData });
    },
    [updateTransactionMutation]
  );

  // Função para deletar transação (wrapper da mutation)
  const deleteTransaction = useCallback(
    async (transactionId: string) => {
      return deleteTransactionMutation.mutateAsync(transactionId);
    },
    [deleteTransactionMutation]
  );

  return {
    // Estado
    transactions: transactionsData?.transactions || [],
    total: transactionsData?.count || 0,
    isLoading: isLoadingTransactions,
    error: transactionsError,

    // Ações
    createTransaction,
    updateTransaction,
    deleteTransaction,
    refetchTransactions,

    // Estados das mutations
    isCreating: createTransactionMutation.isPending,
    isUpdating: updateTransactionMutation.isPending,
    isDeleting: deleteTransactionMutation.isPending,
    createError: createTransactionMutation.error,
    updateError: updateTransactionMutation.error,
    deleteError: deleteTransactionMutation.error,
  };
}

/**
 * Hook para obter uma transação específica por ID
 */
export function useTransaction(transactionId: string | null) {
  const queryClient = useQueryClient();

  // Query para obter uma transação específica
  const {
    data: transaction,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: [...TRANSACTIONS_QUERY_KEY, transactionId],
    queryFn: async () => {
      if (!transactionId) {
        throw new Error("Transaction ID is required");
      }
      return await transactionsService.getById(transactionId);
    },
    enabled: !!transactionId, // Só executa se transactionId estiver definido
    staleTime: 1000 * 60 * 5, // 5 minutos
    gcTime: 1000 * 60 * 10, // 10 minutos
  });

  return {
    transaction,
    isLoading,
    error,
    refetch,
  };
}


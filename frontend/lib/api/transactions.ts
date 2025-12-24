import apiClient from './client';
import type {
  Transaction,
  CreateTransactionRequest,
  UpdateTransactionRequest,
  ListTransactionsResponse,
} from './types';

/**
 * Serviço de transações
 */
export const transactionsService = {
  /**
   * Lista todas as transações do usuário autenticado
   */
  async list(accountId?: string): Promise<ListTransactionsResponse> {
    const params = accountId ? { account_id: accountId } : {};
    const response = await apiClient.get<{ message: string; data: ListTransactionsResponse }>('/transactions', {
      params,
    });
    // O backend retorna { message: "...", data: { transactions: [...], count: ... } }
    return response.data.data;
  },

  /**
   * Obtém uma transação específica por ID
   */
  async getById(transactionId: string): Promise<Transaction> {
    const response = await apiClient.get<{ message: string; data: Transaction }>(`/transactions/${transactionId}`);
    // O backend retorna { message: "...", data: {...} }
    return response.data.data;
  },

  /**
   * Cria uma nova transação
   */
  async create(transactionData: CreateTransactionRequest): Promise<Transaction> {
    const response = await apiClient.post<{ message: string; data: Transaction }>('/transactions', transactionData);
    // O backend retorna { message: "...", data: {...} }
    return response.data.data;
  },

  /**
   * Atualiza uma transação existente
   */
  async update(
    transactionId: string,
    transactionData: UpdateTransactionRequest
  ): Promise<Transaction> {
    const response = await apiClient.put<{ message: string; data: Transaction }>(
      `/transactions/${transactionId}`,
      transactionData
    );
    // O backend retorna { message: "...", data: {...} }
    return response.data.data;
  },

  /**
   * Deleta uma transação
   */
  async delete(transactionId: string): Promise<void> {
    await apiClient.delete(`/transactions/${transactionId}`);
  },
};


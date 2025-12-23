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
    const response = await apiClient.get<ListTransactionsResponse>('/transactions', {
      params,
    });
    return response.data;
  },

  /**
   * Obtém uma transação específica por ID
   */
  async getById(transactionId: string): Promise<Transaction> {
    const response = await apiClient.get<Transaction>(`/transactions/${transactionId}`);
    return response.data;
  },

  /**
   * Cria uma nova transação
   */
  async create(transactionData: CreateTransactionRequest): Promise<Transaction> {
    const response = await apiClient.post<Transaction>('/transactions', transactionData);
    return response.data;
  },

  /**
   * Atualiza uma transação existente
   */
  async update(
    transactionId: string,
    transactionData: UpdateTransactionRequest
  ): Promise<Transaction> {
    const response = await apiClient.put<Transaction>(
      `/transactions/${transactionId}`,
      transactionData
    );
    return response.data;
  },

  /**
   * Deleta uma transação
   */
  async delete(transactionId: string): Promise<void> {
    await apiClient.delete(`/transactions/${transactionId}`);
  },
};


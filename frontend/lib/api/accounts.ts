import apiClient from './client';
import type {
  Account,
  CreateAccountRequest,
  ListAccountsResponse,
} from './types';

/**
 * Serviço de contas
 */
export const accountsService = {
  /**
   * Lista todas as contas do usuário autenticado
   */
  async list(): Promise<ListAccountsResponse> {
    const response = await apiClient.get<{ message: string; data: ListAccountsResponse }>('/accounts');
    // O backend retorna { message: "...", data: { accounts: [...], total: ... } }
    return response.data.data;
  },

  /**
   * Obtém uma conta específica por ID
   */
  async getById(accountId: string): Promise<Account> {
    const response = await apiClient.get<{ message: string; data: Account }>(`/accounts/${accountId}`);
    // O backend retorna { message: "...", data: {...} }
    return response.data.data;
  },

  /**
   * Cria uma nova conta
   */
  async create(accountData: CreateAccountRequest): Promise<Account> {
    const response = await apiClient.post<{ message: string; data: Account }>('/accounts', accountData);
    // O backend retorna { message: "...", data: {...} }
    return response.data.data;
  },
};


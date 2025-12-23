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
    const response = await apiClient.get<ListAccountsResponse>('/accounts');
    return response.data;
  },

  /**
   * Obtém uma conta específica por ID
   */
  async getById(accountId: string): Promise<Account> {
    const response = await apiClient.get<Account>(`/accounts/${accountId}`);
    return response.data;
  },

  /**
   * Cria uma nova conta
   */
  async create(accountData: CreateAccountRequest): Promise<Account> {
    const response = await apiClient.post<Account>('/accounts', accountData);
    return response.data;
  },
};


"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Sidebar } from "@/components/layout";
import { ProtectedRoute } from "@/components/auth";
import { TransactionForm } from "@/components/transactions";
import { useTransactions } from "@/lib/hooks/useTransactions";
import type { CreateTransactionRequest } from "@/lib/api/types";

export default function NewTransactionPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const accountId = searchParams.get("account_id"); // Permite pré-selecionar conta via query param
  
  const { createTransaction, isCreating, createError } = useTransactions();

  const handleSubmit = async (data: CreateTransactionRequest) => {
    await createTransaction(data);
    // Redirecionar para a lista de transações após criar
    router.push("/transactions");
  };

  const handleCancel = () => {
    router.push("/transactions");
  };

  return (
    <ProtectedRoute>
      <div className="flex">
        <Sidebar />
        <div className="ml-64 flex-1 p-8">
          <div className="container max-w-2xl">
            <div className="space-y-6">
              <div>
                <h1 className="text-4xl font-bold mb-2">Nova Transação</h1>
                <p className="text-muted-foreground">
                  Preencha os dados para criar uma nova transação
                </p>
              </div>

              <TransactionForm
                onSubmit={handleSubmit}
                onCancel={handleCancel}
                isLoading={isCreating}
                error={
                  createError
                    ? (createError as any)?.response?.data?.error ||
                      (createError as any)?.message ||
                      "Erro ao criar transação"
                    : null
                }
                defaultAccountId={accountId || undefined}
              />
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}


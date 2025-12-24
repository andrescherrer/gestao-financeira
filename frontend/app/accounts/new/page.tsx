"use client";

import { useRouter } from "next/navigation";
import { Sidebar } from "@/components/layout";
import { ProtectedRoute } from "@/components/auth";
import { AccountForm } from "@/components/accounts";
import { useAccounts } from "@/lib/hooks/useAccounts";
import type { CreateAccountRequest } from "@/lib/api/types";

export default function NewAccountPage() {
  const router = useRouter();
  const { createAccount, isCreating, createError } = useAccounts();

  const handleSubmit = async (data: CreateAccountRequest) => {
    await createAccount(data);
    // Redirecionar para a lista de contas após criar
    router.push("/accounts");
  };

  const handleCancel = () => {
    router.push("/accounts");
  };

  return (
    <ProtectedRoute>
      <div className="flex">
        <Sidebar />
        <div className="ml-64 flex-1 p-8">
          <div className="container max-w-2xl">
            <div className="space-y-6">
              <div>
                <h1 className="text-4xl font-bold mb-2">Nova Conta</h1>
                <p className="text-muted-foreground">
                  Preencha os dados para criar uma nova conta
                </p>
              </div>

              <AccountForm
                onSubmit={handleSubmit}
                onCancel={handleCancel}
                isLoading={isCreating}
                error={
                  createError
                    ? (createError as any)?.response?.data?.error ||
                      (createError as any)?.message ||
                      "Erro ao criar conta"
                    : null
                }
              />
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}


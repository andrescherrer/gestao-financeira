"use client";

import { Sidebar } from "@/components/layout";
import { ProtectedRoute } from "@/components/auth";
import { AccountCard } from "@/components/accounts";
import { useAccounts } from "@/lib/hooks/useAccounts";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { LoadingSpinner } from "@/components/auth";
import { ErrorDisplay } from "@/components/auth";

export default function AccountsPage() {
  const { accounts, isLoading, error, refetchAccounts } = useAccounts();

  return (
    <ProtectedRoute>
      <div className="flex">
        <Sidebar />
        <div className="ml-64 flex-1 p-8">
          <div className="container space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-4xl font-bold mb-2">Contas</h1>
                <p className="text-muted-foreground">
                  Gerencie suas contas bancárias e financeiras
                </p>
              </div>
              <Link href="/accounts/new">
                <Button>Nova Conta</Button>
              </Link>
            </div>

            {/* Loading State */}
            {isLoading && (
              <div className="flex items-center justify-center py-12">
                <div className="text-center">
                  <LoadingSpinner className="mx-auto mb-4" />
                  <p className="text-muted-foreground">Carregando contas...</p>
                </div>
              </div>
            )}

            {/* Error State */}
            {error && !isLoading && (
              <ErrorDisplay
                error={
                  (error as any)?.response?.data?.error ||
                  (error as any)?.message ||
                  "Erro ao carregar contas. Tente novamente."
                }
              />
            )}

            {/* Empty State */}
            {!isLoading && !error && accounts.length === 0 && (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <p className="text-lg text-muted-foreground mb-4">
                  Você ainda não possui contas cadastradas
                </p>
                <Link href="/accounts/new">
                  <Button>Criar Primeira Conta</Button>
                </Link>
              </div>
            )}

            {/* Accounts Grid */}
            {!isLoading && !error && accounts.length > 0 && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {accounts.map((account) => (
                  <AccountCard key={account.account_id} account={account} />
                ))}
              </div>
            )}

            {/* Total Summary */}
            {!isLoading && !error && accounts.length > 0 && (
              <div className="mt-8 p-4 bg-muted rounded-lg">
                <p className="text-sm text-muted-foreground">
                  Total de contas: <span className="font-semibold">{accounts.length}</span>
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}


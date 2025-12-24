"use client";

import { Sidebar } from "@/components/layout";
import { ProtectedRoute } from "@/components/auth";
import { TransactionTable } from "@/components/transactions";
import { useTransactions } from "@/lib/hooks/useTransactions";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { LoadingSpinner } from "@/components/auth";
import { ErrorDisplay } from "@/components/auth";

export default function TransactionsPage() {
  const { transactions, total, isLoading, error, refetchTransactions } = useTransactions();

  return (
    <ProtectedRoute>
      <div className="flex">
        <Sidebar />
        <div className="ml-64 flex-1 p-8">
          <div className="container space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-4xl font-bold mb-2">Transações</h1>
                <p className="text-muted-foreground">
                  Gerencie suas receitas e despesas
                </p>
              </div>
              <Link href="/transactions/new">
                <Button>Nova Transação</Button>
              </Link>
            </div>

            {/* Loading State */}
            {isLoading && (
              <div className="flex items-center justify-center py-12">
                <div className="text-center">
                  <LoadingSpinner className="mx-auto mb-4" />
                  <p className="text-muted-foreground">Carregando transações...</p>
                </div>
              </div>
            )}

            {/* Error State */}
            {error && !isLoading && (
              <ErrorDisplay
                error={
                  (error as any)?.response?.data?.error ||
                  (error as any)?.message ||
                  "Erro ao carregar transações. Tente novamente."
                }
              />
            )}

            {/* Transactions Table */}
            {!isLoading && !error && (
              <>
                <TransactionTable transactions={transactions} />
                
                {/* Summary */}
                {transactions.length > 0 && (
                  <div className="mt-4 p-4 bg-muted rounded-lg">
                    <p className="text-sm text-muted-foreground">
                      Total de transações: <span className="font-semibold">{total}</span>
                    </p>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}


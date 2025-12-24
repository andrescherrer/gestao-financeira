"use client";

import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Sidebar } from "@/components/layout";
import { ProtectedRoute } from "@/components/auth";
import { useTransaction } from "@/lib/hooks/useTransactions";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LoadingSpinner } from "@/components/auth";
import { ErrorDisplay } from "@/components/auth";
import { cn } from "@/lib/utils";

/**
 * Formata um valor monetário
 */
function formatCurrencyValue(value: string, currency: string = "BRL"): string {
  const numValue = parseFloat(value);
  if (isNaN(numValue)) {
    return `R$ 0,00`;
  }

  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: currency === "BRL" ? "BRL" : currency === "USD" ? "USD" : "EUR",
  }).format(numValue);
}

/**
 * Mapeia o tipo de transação para um label
 */
const transactionTypeLabels: Record<string, string> = {
  INCOME: "Receita",
  EXPENSE: "Despesa",
};

export default function TransactionDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const transactionId = params.id as string;

  const { transaction, isLoading, error } = useTransaction(transactionId);

  if (isLoading) {
    return (
      <ProtectedRoute>
        <div className="flex">
          <Sidebar />
          <div className="ml-64 flex-1 p-8">
            <div className="container">
              <div className="flex items-center justify-center py-12">
                <div className="text-center">
                  <LoadingSpinner className="mx-auto mb-4" />
                  <p className="text-muted-foreground">Carregando transação...</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  if (error || !transaction) {
    return (
      <ProtectedRoute>
        <div className="flex">
          <Sidebar />
          <div className="ml-64 flex-1 p-8">
            <div className="container">
              <ErrorDisplay
                error={
                  error
                    ? (error as any)?.response?.data?.error ||
                      (error as any)?.message ||
                      "Erro ao carregar transação"
                    : "Transação não encontrada"
                }
              />
              <div className="mt-4">
                <Link href="/transactions">
                  <Button variant="outline">Voltar para Transações</Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  const amount = formatCurrencyValue(transaction.amount, transaction.currency);
  const isIncome = transaction.type === "INCOME";
  const typeColor = isIncome
    ? "text-green-600 dark:text-green-400"
    : "text-red-600 dark:text-red-400";

  return (
    <ProtectedRoute>
      <div className="flex">
        <Sidebar />
        <div className="ml-64 flex-1 p-8">
          <div className="container max-w-4xl space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-4xl font-bold mb-2">Detalhes da Transação</h1>
                <p className="text-muted-foreground">
                  Informações completas da transação
                </p>
              </div>
              <Link href="/transactions">
                <Button variant="outline">Voltar</Button>
              </Link>
            </div>

            {/* Transaction Details Card */}
            <Card>
              <CardHeader>
                <CardTitle>Informações da Transação</CardTitle>
                <CardDescription>Detalhes completos da transação</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">Descrição</p>
                    <p className="text-lg font-semibold">{transaction.description}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Tipo</p>
                    <p className={cn("text-lg font-semibold", typeColor)}>
                      {transactionTypeLabels[transaction.type] || transaction.type}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Data</p>
                    <p className="text-lg font-semibold">
                      {new Date(transaction.date).toLocaleDateString("pt-BR", {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                      })}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Moeda</p>
                    <p className="text-lg font-semibold">{transaction.currency}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">ID da Conta</p>
                    <p className="text-lg font-semibold font-mono text-sm">
                      {transaction.account_id}
                    </p>
                  </div>
                </div>

                <div className="pt-4 border-t">
                  <p className="text-sm text-muted-foreground mb-2">Valor</p>
                  <p className={cn("text-3xl font-bold", typeColor)}>
                    {isIncome ? "+" : "-"} {amount}
                  </p>
                </div>

                <div className="pt-4 border-t grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-muted-foreground">
                  <div>
                    <p>Criada em</p>
                    <p className="font-semibold text-foreground">
                      {new Date(transaction.created_at).toLocaleDateString("pt-BR", {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </p>
                  </div>
                  <div>
                    <p>Última atualização</p>
                    <p className="font-semibold text-foreground">
                      {new Date(transaction.updated_at).toLocaleDateString("pt-BR", {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}


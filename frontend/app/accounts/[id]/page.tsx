"use client";

import { use } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Sidebar } from "@/components/layout";
import { ProtectedRoute } from "@/components/auth";
import { useAccount } from "@/lib/hooks/useAccounts";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LoadingSpinner } from "@/components/auth";
import { ErrorDisplay } from "@/components/auth";

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
 * Mapeia o tipo de conta para um label
 */
const accountTypeLabels: Record<string, string> = {
  BANK: "Banco",
  WALLET: "Carteira Digital",
  INVESTMENT: "Investimento",
  CREDIT_CARD: "Cartão de Crédito",
};

/**
 * Mapeia o contexto para um label
 */
const contextLabels: Record<string, string> = {
  PERSONAL: "Pessoal",
  BUSINESS: "Empresarial",
};

export default function AccountDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const accountId = params.id as string;

  const { account, isLoading, error } = useAccount(accountId);

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
                  <p className="text-muted-foreground">Carregando conta...</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  if (error || !account) {
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
                      "Erro ao carregar conta"
                    : "Conta não encontrada"
                }
              />
              <div className="mt-4">
                <Link href="/accounts">
                  <Button variant="outline">Voltar para Contas</Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  const balance = formatCurrencyValue(account.balance, account.currency);
  const isPositive = parseFloat(account.balance) >= 0;

  return (
    <ProtectedRoute>
      <div className="flex">
        <Sidebar />
        <div className="ml-64 flex-1 p-8">
          <div className="container max-w-4xl space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-4xl font-bold mb-2">{account.name}</h1>
                <p className="text-muted-foreground">
                  Detalhes da conta
                </p>
              </div>
              <Link href="/accounts">
                <Button variant="outline">Voltar</Button>
              </Link>
            </div>

            {/* Account Details Card */}
            <Card>
              <CardHeader>
                <CardTitle>Informações da Conta</CardTitle>
                <CardDescription>Detalhes completos da conta</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">Nome</p>
                    <p className="text-lg font-semibold">{account.name}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Tipo</p>
                    <p className="text-lg font-semibold">
                      {accountTypeLabels[account.type] || account.type}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Contexto</p>
                    <p className="text-lg font-semibold">
                      {contextLabels[account.context] || account.context}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Moeda</p>
                    <p className="text-lg font-semibold">{account.currency}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Status</p>
                    <p className="text-lg font-semibold">
                      {account.is_active ? (
                        <span className="text-green-600 dark:text-green-400">Ativa</span>
                      ) : (
                        <span className="text-red-600 dark:text-red-400">Inativa</span>
                      )}
                    </p>
                  </div>
                </div>

                <div className="pt-4 border-t">
                  <p className="text-sm text-muted-foreground mb-2">Saldo Atual</p>
                  <p
                    className={cn(
                      "text-3xl font-bold",
                      isPositive
                        ? "text-green-600 dark:text-green-400"
                        : "text-red-600 dark:text-red-400"
                    )}
                  >
                    {balance}
                  </p>
                </div>

                <div className="pt-4 border-t grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-muted-foreground">
                  <div>
                    <p>Criada em</p>
                    <p className="font-semibold text-foreground">
                      {new Date(account.created_at).toLocaleDateString("pt-BR", {
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
                      {new Date(account.updated_at).toLocaleDateString("pt-BR", {
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

// Helper function para cn (se não estiver importado)
import { cn } from "@/lib/utils";


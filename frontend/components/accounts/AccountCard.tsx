"use client";

import Link from "next/link";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import type { Account } from "@/lib/api/types";
import { cn } from "@/lib/utils";

interface AccountCardProps {
  account: Account;
  className?: string;
}

/**
 * Mapeia o tipo de conta para um ícone e label
 */
const accountTypeConfig: Record<Account["type"], { icon: string; label: string }> = {
  BANK: { icon: "🏦", label: "Banco" },
  WALLET: { icon: "💳", label: "Carteira Digital" },
  INVESTMENT: { icon: "📈", label: "Investimento" },
  CREDIT_CARD: { icon: "💳", label: "Cartão de Crédito" },
};

/**
 * Formata um valor monetário
 */
function formatCurrency(value: string, currency: string = "BRL"): string {
  const numValue = parseFloat(value);
  if (isNaN(numValue)) {
    return `R$ 0,00`;
  }

  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: currency === "BRL" ? "BRL" : "USD",
  }).format(numValue);
}

/**
 * Componente de card para exibir uma conta
 */
export function AccountCard({ account, className }: AccountCardProps) {
  const typeConfig = accountTypeConfig[account.type];
  const balance = formatCurrency(account.balance, account.currency);
  const isPositive = parseFloat(account.balance) >= 0;

  return (
    <Card className={cn("hover:shadow-md transition-shadow", className)}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-2xl">{typeConfig.icon}</span>
            <div>
              <CardTitle className="text-lg">{account.name}</CardTitle>
              <CardDescription>{typeConfig.label}</CardDescription>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div>
            <p className="text-sm text-muted-foreground">Saldo</p>
            <p
              className={cn(
                "text-2xl font-bold",
                isPositive ? "text-green-600 dark:text-green-400" : "text-red-600 dark:text-red-400"
              )}
            >
              {balance}
            </p>
          </div>
          <div className="flex gap-2">
            <Link href={`/accounts/${account.account_id}`} className="flex-1">
              <Button variant="outline" className="w-full" size="sm">
                Ver Detalhes
              </Button>
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}


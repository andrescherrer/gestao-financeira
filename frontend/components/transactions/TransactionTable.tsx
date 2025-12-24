"use client";

import Link from "next/link";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import type { Transaction } from "@/lib/api/types";
import { cn } from "@/lib/utils";

interface TransactionTableProps {
  transactions: Transaction[];
  className?: string;
}

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
    currency: currency === "BRL" ? "BRL" : currency === "USD" ? "USD" : "EUR",
  }).format(numValue);
}

/**
 * Formata uma data
 */
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("pt-BR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

/**
 * Mapeia o tipo de transação para um label e cor
 */
const transactionTypeConfig: Record<Transaction["type"], { label: string; color: string }> = {
  INCOME: { label: "Receita", color: "text-green-600 dark:text-green-400" },
  EXPENSE: { label: "Despesa", color: "text-red-600 dark:text-red-400" },
};

/**
 * Componente de tabela para exibir transações
 */
export function TransactionTable({ transactions, className }: TransactionTableProps) {
  if (transactions.length === 0) {
    return (
      <div className={cn("text-center py-12 text-muted-foreground", className)}>
        <p>Nenhuma transação encontrada</p>
      </div>
    );
  }

  return (
    <div className={cn("rounded-md border", className)}>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Data</TableHead>
            <TableHead>Descrição</TableHead>
            <TableHead>Tipo</TableHead>
            <TableHead className="text-right">Valor</TableHead>
            <TableHead className="text-right">Ações</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {transactions.map((transaction) => {
            const typeConfig = transactionTypeConfig[transaction.type];
            const amount = formatCurrency(transaction.amount, transaction.currency);
            const isIncome = transaction.type === "INCOME";

            return (
              <TableRow key={transaction.transaction_id}>
                <TableCell className="font-medium">
                  {formatDate(transaction.date)}
                </TableCell>
                <TableCell>{transaction.description}</TableCell>
                <TableCell>
                  <span className={typeConfig.color}>{typeConfig.label}</span>
                </TableCell>
                <TableCell className={cn("text-right font-semibold", typeConfig.color)}>
                  {isIncome ? "+" : "-"} {amount}
                </TableCell>
                <TableCell className="text-right">
                  <Link href={`/transactions/${transaction.transaction_id}`}>
                    <Button variant="ghost" size="sm">
                      Ver Detalhes
                    </Button>
                  </Link>
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}


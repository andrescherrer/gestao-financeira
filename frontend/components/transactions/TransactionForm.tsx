"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ErrorDisplay } from "@/components/auth";
import { LoadingSpinner } from "@/components/auth";
import { useAccounts } from "@/lib/hooks/useAccounts";
import type { CreateTransactionRequest } from "@/lib/api/types";
import { cn } from "@/lib/utils";

const transactionFormSchema = z.object({
  account_id: z.string().min(1, "Selecione uma conta"),
  type: z.enum(["INCOME", "EXPENSE"]).refine((val) => val !== undefined, {
    message: "Selecione um tipo de transação",
  }),
  amount: z
    .string()
    .min(1, "Valor é obrigatório")
    .refine(
      (val) => {
        const num = parseFloat(val);
        return !isNaN(num) && num > 0;
      },
      {
        message: "Valor deve ser um número positivo",
      }
    ),
  currency: z.enum(["BRL", "USD", "EUR"]).refine((val) => val !== undefined, {
    message: "Selecione uma moeda",
  }),
  description: z.string().min(3, "Descrição deve ter no mínimo 3 caracteres").max(500, "Descrição deve ter no máximo 500 caracteres"),
  date: z.string().min(1, "Data é obrigatória"),
});

type TransactionFormData = z.infer<typeof transactionFormSchema>;

interface TransactionFormProps {
  onSubmit: (data: CreateTransactionRequest) => Promise<void>;
  onCancel?: () => void;
  isLoading?: boolean;
  error?: string | null;
  className?: string;
  defaultAccountId?: string;
}

/**
 * Componente de formulário para criar/editar transação
 * Reutilizável em diferentes contextos (página, modal, etc.)
 */
export function TransactionForm({
  onSubmit,
  onCancel,
  isLoading = false,
  error,
  className,
  defaultAccountId,
}: TransactionFormProps) {
  const { accounts, isLoading: isLoadingAccounts } = useAccounts();
  const [formError, setFormError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
  } = useForm<TransactionFormData>({
    resolver: zodResolver(transactionFormSchema),
    defaultValues: {
      currency: "BRL",
      date: new Date().toISOString().split("T")[0], // Data atual no formato YYYY-MM-DD
      account_id: defaultAccountId || "",
    },
  });

  const typeValue = watch("type");
  const currencyValue = watch("currency");
  const accountIdValue = watch("account_id");

  // Atualizar account_id quando defaultAccountId mudar
  useEffect(() => {
    if (defaultAccountId && !accountIdValue) {
      setValue("account_id", defaultAccountId);
    }
  }, [defaultAccountId, accountIdValue, setValue]);

  const onFormSubmit = async (data: TransactionFormData) => {
    setFormError(null);

    try {
      const transactionData: CreateTransactionRequest = {
        account_id: data.account_id,
        type: data.type,
        amount: data.amount,
        currency: data.currency,
        description: data.description,
        date: data.date,
      };

      await onSubmit(transactionData);
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error ||
        err.response?.data?.message ||
        err.message ||
        "Erro ao criar transação. Tente novamente.";

      setFormError(errorMessage);
    }
  };

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className={cn("space-y-6", className)}>
      {(error || formError) && (
        <ErrorDisplay
          error={error || formError || undefined}
        />
      )}

      <div className="space-y-4">
        {/* Conta */}
        <div className="space-y-2">
          <Label htmlFor="account_id">Conta *</Label>
          {isLoadingAccounts ? (
            <div className="flex items-center gap-2">
              <LoadingSpinner />
              <span className="text-sm text-muted-foreground">Carregando contas...</span>
            </div>
          ) : (
            <Select
              value={accountIdValue}
              onValueChange={(value) => setValue("account_id", value)}
              disabled={isLoading || isLoadingAccounts}
            >
              <SelectTrigger id="account_id">
                <SelectValue placeholder="Selecione uma conta" />
              </SelectTrigger>
              <SelectContent>
                {accounts.map((account) => (
                  <SelectItem key={account.account_id} value={account.account_id}>
                    {account.name} ({accountTypeLabels[account.type] || account.type})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
          {errors.account_id && (
            <p className="text-sm text-destructive">{errors.account_id.message}</p>
          )}
          {accounts.length === 0 && !isLoadingAccounts && (
            <p className="text-sm text-muted-foreground">
              Nenhuma conta disponível.{" "}
              <a href="/accounts/new" className="text-primary hover:underline">
                Criar conta
              </a>
            </p>
          )}
        </div>

        {/* Tipo de Transação */}
        <div className="space-y-2">
          <Label htmlFor="type">Tipo de Transação *</Label>
          <Select
            value={typeValue}
            onValueChange={(value) => setValue("type", value as TransactionFormData["type"])}
            disabled={isLoading}
          >
            <SelectTrigger id="type">
              <SelectValue placeholder="Selecione o tipo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="INCOME">Receita</SelectItem>
              <SelectItem value="EXPENSE">Despesa</SelectItem>
            </SelectContent>
          </Select>
          {errors.type && (
            <p className="text-sm text-destructive">{errors.type.message}</p>
          )}
        </div>

        {/* Valor e Moeda */}
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="amount">Valor *</Label>
            <Input
              id="amount"
              type="number"
              step="0.01"
              min="0.01"
              placeholder="0.00"
              {...register("amount")}
              disabled={isLoading}
            />
            {errors.amount && (
              <p className="text-sm text-destructive">{errors.amount.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="currency">Moeda *</Label>
            <Select
              value={currencyValue}
              onValueChange={(value) => setValue("currency", value as TransactionFormData["currency"])}
              disabled={isLoading}
            >
              <SelectTrigger id="currency">
                <SelectValue placeholder="Selecione a moeda" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="BRL">Real (BRL)</SelectItem>
                <SelectItem value="USD">Dólar (USD)</SelectItem>
                <SelectItem value="EUR">Euro (EUR)</SelectItem>
              </SelectContent>
            </Select>
            {errors.currency && (
              <p className="text-sm text-destructive">{errors.currency.message}</p>
            )}
          </div>
        </div>

        {/* Descrição */}
        <div className="space-y-2">
          <Label htmlFor="description">Descrição *</Label>
          <Textarea
            id="description"
            placeholder="Descreva a transação..."
            {...register("description")}
            disabled={isLoading}
            rows={3}
          />
          {errors.description && (
            <p className="text-sm text-destructive">{errors.description.message}</p>
          )}
          <p className="text-xs text-muted-foreground">
            {watch("description")?.length || 0}/500 caracteres
          </p>
        </div>

        {/* Data */}
        <div className="space-y-2">
          <Label htmlFor="date">Data *</Label>
          <Input
            id="date"
            type="date"
            {...register("date")}
            disabled={isLoading}
          />
          {errors.date && (
            <p className="text-sm text-destructive">{errors.date.message}</p>
          )}
        </div>
      </div>

      {/* Botões */}
      <div className="flex gap-4">
        {onCancel && (
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isLoading}
            className="flex-1"
          >
            Cancelar
          </Button>
        )}
        <Button type="submit" disabled={isLoading || isLoadingAccounts || accounts.length === 0} className="flex-1">
          {isLoading ? <LoadingSpinner /> : "Criar Transação"}
        </Button>
      </div>
    </form>
  );
}

// Helper para labels de tipo de conta
const accountTypeLabels: Record<string, string> = {
  BANK: "Banco",
  WALLET: "Carteira",
  INVESTMENT: "Investimento",
  CREDIT_CARD: "Cartão",
};


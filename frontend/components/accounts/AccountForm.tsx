"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ErrorDisplay } from "@/components/auth";
import { LoadingSpinner } from "@/components/auth";
import type { CreateAccountRequest } from "@/lib/api/types";
import { cn } from "@/lib/utils";

const accountFormSchema = z.object({
  name: z.string().min(3, "Nome deve ter no mínimo 3 caracteres").max(100, "Nome deve ter no máximo 100 caracteres"),
  type: z.enum(["BANK", "WALLET", "INVESTMENT", "CREDIT_CARD"]).refine((val) => val !== undefined, {
    message: "Selecione um tipo de conta",
  }),
  initial_balance: z
    .string()
    .optional()
    .refine(
      (val) => {
        if (!val || val === "") return true; // Opcional
        const num = parseFloat(val);
        return !isNaN(num) && num >= 0;
      },
      {
        message: "Saldo inicial deve ser um número positivo ou zero",
      }
    ),
  currency: z.enum(["BRL", "USD", "EUR"]).refine((val) => val !== undefined, {
    message: "Selecione uma moeda",
  }),
  context: z.enum(["PERSONAL", "BUSINESS"]).refine((val) => val !== undefined, {
    message: "Selecione um contexto",
  }),
});

type AccountFormData = z.infer<typeof accountFormSchema>;

interface AccountFormProps {
  onSubmit: (data: CreateAccountRequest) => Promise<void>;
  onCancel?: () => void;
  isLoading?: boolean;
  error?: string | null;
  className?: string;
}

/**
 * Componente de formulário para criar/editar conta
 * Reutilizável em diferentes contextos (página, modal, etc.)
 */
export function AccountForm({
  onSubmit,
  onCancel,
  isLoading = false,
  error,
  className,
}: AccountFormProps) {
  const [formError, setFormError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
  } = useForm<AccountFormData>({
    resolver: zodResolver(accountFormSchema),
    defaultValues: {
      currency: "BRL",
      context: "PERSONAL",
    },
  });

  const typeValue = watch("type");
  const currencyValue = watch("currency");
  const contextValue = watch("context");

  const onFormSubmit = async (data: AccountFormData) => {
    setFormError(null);

    try {
      const accountData: CreateAccountRequest = {
        name: data.name,
        type: data.type,
        currency: data.currency,
        context: data.context,
        initial_balance: data.initial_balance
          ? parseFloat(data.initial_balance)
          : undefined,
      };

      await onSubmit(accountData);
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error ||
        err.response?.data?.message ||
        err.message ||
        "Erro ao criar conta. Tente novamente.";

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
        {/* Nome da Conta */}
        <div className="space-y-2">
          <Label htmlFor="name">Nome da Conta *</Label>
          <Input
            id="name"
            type="text"
            placeholder="Ex: Conta Corrente Principal"
            {...register("name")}
            disabled={isLoading}
          />
          {errors.name && (
            <p className="text-sm text-destructive">{errors.name.message}</p>
          )}
        </div>

        {/* Tipo de Conta */}
        <div className="space-y-2">
          <Label htmlFor="type">Tipo de Conta *</Label>
          <Select
            value={typeValue}
            onValueChange={(value) => setValue("type", value as AccountFormData["type"])}
            disabled={isLoading}
          >
            <SelectTrigger id="type">
              <SelectValue placeholder="Selecione o tipo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="BANK">Banco</SelectItem>
              <SelectItem value="WALLET">Carteira Digital</SelectItem>
              <SelectItem value="INVESTMENT">Investimento</SelectItem>
              <SelectItem value="CREDIT_CARD">Cartão de Crédito</SelectItem>
            </SelectContent>
          </Select>
          {errors.type && (
            <p className="text-sm text-destructive">{errors.type.message}</p>
          )}
        </div>

        {/* Contexto */}
        <div className="space-y-2">
          <Label htmlFor="context">Contexto *</Label>
          <Select
            value={contextValue}
            onValueChange={(value) => setValue("context", value as AccountFormData["context"])}
            disabled={isLoading}
          >
            <SelectTrigger id="context">
              <SelectValue placeholder="Selecione o contexto" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="PERSONAL">Pessoal</SelectItem>
              <SelectItem value="BUSINESS">Empresarial</SelectItem>
            </SelectContent>
          </Select>
          {errors.context && (
            <p className="text-sm text-destructive">{errors.context.message}</p>
          )}
        </div>

        {/* Moeda */}
        <div className="space-y-2">
          <Label htmlFor="currency">Moeda *</Label>
          <Select
            value={currencyValue}
            onValueChange={(value) => setValue("currency", value as AccountFormData["currency"])}
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

        {/* Saldo Inicial */}
        <div className="space-y-2">
          <Label htmlFor="initial_balance">Saldo Inicial (opcional)</Label>
          <Input
            id="initial_balance"
            type="number"
            step="0.01"
            min="0"
            placeholder="0.00"
            {...register("initial_balance")}
            disabled={isLoading}
          />
          {errors.initial_balance && (
            <p className="text-sm text-destructive">
              {errors.initial_balance.message}
            </p>
          )}
          <p className="text-xs text-muted-foreground">
            Deixe em branco para começar com saldo zero
          </p>
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
        <Button type="submit" disabled={isLoading} className="flex-1">
          {isLoading ? <LoadingSpinner /> : "Criar Conta"}
        </Button>
      </div>
    </form>
  );
}


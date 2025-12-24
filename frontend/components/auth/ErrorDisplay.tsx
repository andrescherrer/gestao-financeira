"use client";

import { AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";

interface ErrorDisplayProps {
  error: string | null | undefined;
  className?: string;
}

/**
 * Componente para exibir erros de forma consistente
 */
export function ErrorDisplay({ error, className }: ErrorDisplayProps) {
  if (!error) return null;

  return (
    <div
      className={cn(
        "rounded-md bg-destructive/10 p-3 text-sm text-destructive flex items-start gap-2",
        className
      )}
    >
      <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
      <span>{error}</span>
    </div>
  );
}


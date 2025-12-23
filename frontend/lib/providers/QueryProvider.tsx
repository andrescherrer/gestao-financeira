"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { useState } from "react";
import { env } from "@/lib/config/env";

interface QueryProviderProps {
  children: React.ReactNode;
}

export function QueryProvider({ children }: QueryProviderProps) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            // Tempo de cache padrão: 5 minutos
            staleTime: 5 * 60 * 1000,
            // Tempo de garbage collection: 10 minutos
            gcTime: 10 * 60 * 1000,
            // Retry automático em caso de erro
            retry: 1,
            // Refetch quando a janela recebe foco
            refetchOnWindowFocus: true,
            // Refetch quando reconecta à rede
            refetchOnReconnect: true,
          },
          mutations: {
            // Retry automático em caso de erro
            retry: 1,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {/* React Query Devtools apenas em desenvolvimento */}
      {env.isDevelopment && <ReactQueryDevtools initialIsOpen={false} />}
    </QueryClientProvider>
  );
}


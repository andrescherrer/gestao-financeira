import { LoginForm } from "@/components/auth/LoginForm";

/**
 * Página de Login
 * Usa o componente LoginForm para renderizar o formulário
 */
export default function LoginPage() {
  return (
    <div className="flex min-h-[calc(100vh-4rem)] items-center justify-center bg-background p-4">
      <div className="w-full max-w-md space-y-8 rounded-lg border bg-card p-8 shadow-lg">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Gestão Financeira</h1>
          <p className="mt-2 text-muted-foreground">
            Faça login para acessar sua conta
          </p>
        </div>

        <LoginForm />
      </div>
    </div>
  );
}


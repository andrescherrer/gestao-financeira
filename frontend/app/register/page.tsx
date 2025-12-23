import { RegisterForm } from "@/components/auth/RegisterForm";

/**
 * Página de Registro
 * Usa o componente RegisterForm para renderizar o formulário
 */
export default function RegisterPage() {
  return (
    <div className="flex min-h-[calc(100vh-4rem)] items-center justify-center bg-background p-4">
      <div className="w-full max-w-md space-y-8 rounded-lg border bg-card p-8 shadow-lg">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Criar Conta</h1>
          <p className="mt-2 text-muted-foreground">
            Preencha os dados para criar sua conta
          </p>
        </div>

        <RegisterForm />
      </div>
    </div>
  );
}

"use client";

import { Separator } from "@/components/ui/separator";

export function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="border-t bg-background">
      <div className="container flex h-16 items-center justify-between px-4">
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          <span>© {currentYear} Gestão Financeira</span>
          <Separator orientation="vertical" className="h-4" />
          <span>Todos os direitos reservados</span>
        </div>
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          <a href="#" className="hover:text-foreground">
            Sobre
          </a>
          <Separator orientation="vertical" className="h-4" />
          <a href="#" className="hover:text-foreground">
            Contato
          </a>
        </div>
      </div>
    </footer>
  );
}


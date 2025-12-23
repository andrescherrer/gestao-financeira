"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-xl font-bold">Gestão Financeira</span>
          </Link>
        </div>
        <nav className="flex items-center gap-4">
          <Link href="/accounts">
            <Button variant="ghost" size="sm">
              Contas
            </Button>
          </Link>
          <Link href="/transactions">
            <Button variant="ghost" size="sm">
              Transações
            </Button>
          </Link>
          <Separator orientation="vertical" className="h-6" />
          <Button variant="outline" size="sm">
            Login
          </Button>
        </nav>
      </div>
    </header>
  );
}


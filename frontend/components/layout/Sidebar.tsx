"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface SidebarProps {
  className?: string;
}

const navigation = [
  { name: "Dashboard", href: "/", icon: "📊" },
  { name: "Contas", href: "/accounts", icon: "💳" },
  { name: "Transações", href: "/transactions", icon: "💰" },
  { name: "Relatórios", href: "/reports", icon: "📈" },
];

export function Sidebar({ className }: SidebarProps) {
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        "fixed left-0 top-16 z-40 h-[calc(100vh-4rem)] w-64 border-r bg-background",
        className
      )}
    >
      <nav className="flex flex-col gap-1 p-4">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link key={item.href} href={item.href}>
              <Button
                variant={isActive ? "secondary" : "ghost"}
                className={cn(
                  "w-full justify-start gap-2",
                  isActive && "bg-secondary"
                )}
              >
                <span>{item.icon}</span>
                <span>{item.name}</span>
              </Button>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}


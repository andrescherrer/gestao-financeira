import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { Header, Footer } from '@/components/layout'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Gestão Financeira',
  description: 'Sistema de gestão financeira pessoal e profissional',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="pt-BR">
      <body className={inter.className}>
        <div className="flex min-h-screen flex-col">
          <Header />
          <main className="flex-1">{children}</main>
          <Footer />
        </div>
      </body>
    </html>
  )
}


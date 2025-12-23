import { Sidebar } from '@/components/layout'

export default function Home() {
  return (
    <div className="flex">
      <Sidebar />
      <div className="ml-64 flex-1 p-8">
        <div className="container">
          <h1 className="text-4xl font-bold mb-4">Dashboard</h1>
          <p className="text-muted-foreground">
            Bem-vindo ao sistema de gestão financeira
          </p>
        </div>
      </div>
    </div>
  )
}


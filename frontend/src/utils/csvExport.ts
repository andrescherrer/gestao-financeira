/**
 * Utilitário para exportar dados para CSV
 */

/**
 * Converte um array de objetos para CSV
 */
export function arrayToCSV(data: Record<string, any>[]): string {
  if (data.length === 0) return ''

  // Obter cabeçalhos
  const firstRow = data[0]
  if (!firstRow) return ''
  const headers = Object.keys(firstRow)

  // Criar linha de cabeçalho
  const headerRow = headers.map((h) => `"${h}"`).join(',')

  // Criar linhas de dados
  const dataRows = data.map((row) => {
    return headers
      .map((header) => {
        const value = row?.[header]
        // Escapar aspas e quebras de linha
        const escaped = String(value ?? '')
          .replace(/"/g, '""')
          .replace(/\n/g, ' ')
          .replace(/\r/g, ' ')
        return `"${escaped}"`
      })
      .join(',')
  })

  // Combinar tudo
  return [headerRow, ...dataRows].join('\n')
}

/**
 * Faz download de um arquivo CSV
 */
export function downloadCSV(csvContent: string, filename: string): void {
  // Adicionar BOM para Excel reconhecer UTF-8
  const BOM = '\uFEFF'
  const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)

  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Exporta relatório mensal para CSV
 */
export function exportMonthlyReportToCSV(report: any, year: number, month?: number): void {
  const monthNames = [
    'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
    'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
  ]

  const monthIndex = month ? month - 1 : 0
  const data = [
    {
      'Período': `${monthNames[monthIndex]} ${year}`,
      'Moeda': report.currency,
      'Total Receitas': report.total_income,
      'Total Despesas': report.total_expense,
      'Saldo': report.balance,
      'Quantidade Receitas': report.income_count,
      'Quantidade Despesas': report.expense_count,
      'Total Transações': report.total_count,
    },
  ]

  // Adicionar breakdown por categoria se disponível
  if (report.category_breakdown && report.category_breakdown.length > 0) {
    report.category_breakdown.forEach((cat: any) => {
      data.push({
        'Período': '',
        'Moeda': '',
        'Total Receitas': cat.type === 'INCOME' ? cat.total_amount : '',
        'Total Despesas': cat.type === 'EXPENSE' ? cat.total_amount : '',
        'Saldo': '',
        'Quantidade Receitas': cat.type === 'INCOME' ? cat.count : '',
        'Quantidade Despesas': cat.type === 'EXPENSE' ? cat.count : '',
        'Total Transações': cat.count,
      })
    })
  }

  const csv = arrayToCSV(data)
  const monthStr = month ? month.toString().padStart(2, '0') : '00'
  const filename = `relatorio_mensal_${year}_${monthStr}.csv`
  downloadCSV(csv, filename)
}

/**
 * Exporta relatório anual para CSV
 */
export function exportAnnualReportToCSV(report: any, year: number): void {
  const data = [
    {
      'Ano': year,
      'Moeda': report.currency,
      'Total Receitas': report.total_income,
      'Total Despesas': report.total_expense,
      'Saldo': report.balance,
      'Quantidade Receitas': report.income_count,
      'Quantidade Despesas': report.expense_count,
      'Total Transações': report.total_count,
    },
  ]

  // Adicionar breakdown mensal
  if (report.monthly_breakdown && report.monthly_breakdown.length > 0) {
    const monthNames = [
      'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
      'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
    ]

    report.monthly_breakdown.forEach((month: any) => {
      const monthIndex = month.month ? month.month - 1 : 0
      data.push({
        'Ano': monthNames[monthIndex],
        'Moeda': report.currency,
        'Total Receitas': month.total_income,
        'Total Despesas': month.total_expense,
        'Saldo': month.balance,
        'Quantidade Receitas': month.income_count,
        'Quantidade Despesas': month.expense_count,
        'Total Transações': month.income_count + month.expense_count,
      })
    })
  }

  const csv = arrayToCSV(data)
  const filename = `relatorio_anual_${year}.csv`
  downloadCSV(csv, filename)
}

/**
 * Exporta relatório por categoria para CSV
 */
export function exportCategoryReportToCSV(report: any): void {
  const data = report.category_breakdown.map((cat: any) => ({
    'Categoria': cat.category_name || 'Sem categoria',
    'Tipo': cat.type === 'INCOME' ? 'Receita' : 'Despesa',
    'Valor Total': cat.total_amount,
    'Quantidade': cat.count,
    'Percentual': `${cat.percentage.toFixed(2)}%`,
    'Moeda': report.currency,
  }))

  // Adicionar totais
  data.push({
    'Categoria': 'TOTAL',
    'Tipo': '',
    'Valor Total': report.balance,
    'Quantidade': report.total_count,
    'Percentual': '100%',
    'Moeda': report.currency,
  })

  const csv = arrayToCSV(data)
  const filename = `relatorio_categorias_${new Date().toISOString().split('T')[0]}.csv`
  downloadCSV(csv, filename)
}


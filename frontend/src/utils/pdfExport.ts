/**
 * Utilitário para exportar dados para PDF
 */

import jsPDF from 'jspdf'

/**
 * Exporta relatório mensal para PDF
 */
export function exportMonthlyReportToPDF(report: any, year: number, month?: number): void {
  const doc = new jsPDF()
  const monthNames = [
    'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
    'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
  ]

  const monthName = month ? monthNames[month - 1] : 'N/A'
  const period = `${monthName} ${year}`

  // Título
  doc.setFontSize(18)
  doc.text('Relatório Mensal', 14, 20)
  doc.setFontSize(12)
  doc.text(`Período: ${period}`, 14, 30)
  doc.text(`Moeda: ${report.currency}`, 14, 36)

  // Resumo
  let yPos = 50
  doc.setFontSize(14)
  doc.text('Resumo', 14, yPos)
  yPos += 10

  doc.setFontSize(11)
  doc.text(`Total de Receitas: ${formatCurrency(report.total_income, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Total de Despesas: ${formatCurrency(report.total_expense, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Saldo: ${formatCurrency(report.balance, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Quantidade de Transações: ${report.total_count}`, 14, yPos)

  // Breakdown por categoria
  if (report.category_breakdown && report.category_breakdown.length > 0) {
    yPos += 15
    doc.setFontSize(14)
    doc.text('Breakdown por Categoria', 14, yPos)
    yPos += 10

    doc.setFontSize(10)
    report.category_breakdown.forEach((cat: any) => {
      if (yPos > 270) {
        doc.addPage()
        yPos = 20
      }
      doc.text(
        `${cat.category_name || 'Sem categoria'} (${cat.type}): ${formatCurrency(cat.total_amount, report.currency)}`,
        14,
        yPos
      )
      yPos += 7
    })
  }

  // Salvar PDF
  const filename = `relatorio_mensal_${year}_${month ? month.toString().padStart(2, '0') : '00'}.pdf`
  doc.save(filename)
}

/**
 * Exporta relatório anual para PDF
 */
export function exportAnnualReportToPDF(report: any, year: number): void {
  const doc = new jsPDF()
  const monthNames = [
    'Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun',
    'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'
  ]

  // Título
  doc.setFontSize(18)
  doc.text('Relatório Anual', 14, 20)
  doc.setFontSize(12)
  doc.text(`Ano: ${year}`, 14, 30)
  doc.text(`Moeda: ${report.currency}`, 14, 36)

  // Resumo
  let yPos = 50
  doc.setFontSize(14)
  doc.text('Resumo Anual', 14, yPos)
  yPos += 10

  doc.setFontSize(11)
  doc.text(`Total de Receitas: ${formatCurrency(report.total_income, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Total de Despesas: ${formatCurrency(report.total_expense, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Saldo: ${formatCurrency(report.balance, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Quantidade de Transações: ${report.total_count}`, 14, yPos)

  // Breakdown mensal
  if (report.monthly_breakdown && report.monthly_breakdown.length > 0) {
    yPos += 15
    doc.setFontSize(14)
    doc.text('Breakdown Mensal', 14, yPos)
    yPos += 10

    doc.setFontSize(10)
    report.monthly_breakdown.forEach((month: any) => {
      if (yPos > 270) {
        doc.addPage()
        yPos = 20
      }
      const monthName = monthNames[month.month - 1] || 'N/A'
      doc.text(`${monthName}:`, 14, yPos)
      doc.text(`  Receitas: ${formatCurrency(month.total_income, report.currency)}`, 20, yPos + 5)
      doc.text(`  Despesas: ${formatCurrency(month.total_expense, report.currency)}`, 20, yPos + 10)
      doc.text(`  Saldo: ${formatCurrency(month.balance, report.currency)}`, 20, yPos + 15)
      yPos += 20
    })
  }

  // Salvar PDF
  const filename = `relatorio_anual_${year}.pdf`
  doc.save(filename)
}

/**
 * Exporta relatório por categoria para PDF
 */
export function exportCategoryReportToPDF(report: any): void {
  const doc = new jsPDF()

  // Título
  doc.setFontSize(18)
  doc.text('Relatório por Categoria', 14, 20)
  doc.setFontSize(12)
  doc.text(`Moeda: ${report.currency}`, 14, 30)

  // Resumo
  let yPos = 45
  doc.setFontSize(11)
  doc.text(`Total de Receitas: ${formatCurrency(report.total_income, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Total de Despesas: ${formatCurrency(report.total_expense, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Saldo: ${formatCurrency(report.balance, report.currency)}`, 14, yPos)
  yPos += 7
  doc.text(`Total de Transações: ${report.total_count}`, 14, yPos)

  // Breakdown por categoria
  if (report.category_breakdown && report.category_breakdown.length > 0) {
    yPos += 15
    doc.setFontSize(14)
    doc.text('Breakdown por Categoria', 14, yPos)
    yPos += 10

    doc.setFontSize(10)
    report.category_breakdown.forEach((cat: any) => {
      if (yPos > 270) {
        doc.addPage()
        yPos = 20
      }
      doc.text(
        `${cat.category_name || 'Sem categoria'} (${cat.type === 'INCOME' ? 'Receita' : 'Despesa'}):`,
        14,
        yPos
      )
      doc.text(`  Valor: ${formatCurrency(cat.total_amount, report.currency)}`, 20, yPos + 5)
      doc.text(`  Quantidade: ${cat.count}`, 20, yPos + 10)
      doc.text(`  Percentual: ${cat.percentage.toFixed(2)}%`, 20, yPos + 15)
      yPos += 20
    })
  }

  // Salvar PDF
  const filename = `relatorio_categorias_${new Date().toISOString().split('T')[0]}.pdf`
  doc.save(filename)
}

/**
 * Formata valor como moeda
 */
function formatCurrency(value: number, currency: string): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  }).format(value)
}


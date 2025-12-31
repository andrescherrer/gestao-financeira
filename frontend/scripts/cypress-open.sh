#!/bin/bash
# Script para abrir a interface gráfica do Cypress
# Requer que o servidor de desenvolvimento esteja rodando

echo "🚀 Abrindo Cypress Test Runner..."
echo ""
echo "⚠️  Certifique-se de que o servidor está rodando:"
echo "   npm run preview"
echo "   ou"
echo "   npm run dev"
echo ""
echo "📝 A interface gráfica do Cypress será aberta em breve..."
echo ""

# Verificar se o servidor está rodando
if ! curl -s http://localhost:4173 > /dev/null 2>&1 && ! curl -s http://localhost:3000 > /dev/null 2>&1; then
  echo "⚠️  AVISO: Nenhum servidor detectado em localhost:4173 ou localhost:3000"
  echo "   Inicie o servidor antes de executar os testes!"
  echo ""
fi

# Abrir Cypress
npm run test:e2e:open


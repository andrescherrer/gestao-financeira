#!/bin/bash
# Script para executar testes E2E no Docker com Xvfb

# Exportar DISPLAY antes de iniciar Xvfb
export DISPLAY=:99

# Iniciar Xvfb em background
Xvfb :99 -screen 0 1280x720x24 > /dev/null 2>&1 &
XVFB_PID=$!

# Aguardar Xvfb iniciar
sleep 3

# Verificar se Xvfb está rodando
if ! kill -0 $XVFB_PID 2>/dev/null; then
  echo "Erro: Xvfb não iniciou corretamente (PID: $XVFB_PID)"
  exit 1
fi

# Executar testes E2E
npm run test:e2e
EXIT_CODE=$?

# Matar Xvfb ao finalizar
kill $XVFB_PID 2>/dev/null || true

exit $EXIT_CODE


# Plano de Reorganização de Arquivos Markdown

**Data:** 2025-12-31  
**Objetivo:** Organizar todos os arquivos `.md` (exceto README.md) em categorias dentro de `docs/`

---

## 📊 Análise dos Arquivos

### Arquivos na Raiz do Projeto
- `TAREFAS.md` - Documentação principal de tarefas
- `VERIFICACAO_IMPLEMENTACAO.md` - Verificação histórica (desatualizado)
- `VERIFICACAO_TAREFAS.md` - Verificação histórica (desatualizado)
- `QUICK_START_TESTING.md` - Guia rápido de testes

### Arquivos no Backend
- `backend/CONFIG.md` - Configuração da API
- `backend/docs/SECURITY.md` - Segurança do backend
- `backend/TEST_COVERAGE_ANALYSIS.md` - Análise de cobertura de testes
- `backend/TESTING_GUIDE.md` - Guia completo de testes
- `backend/cmd/process-recurring/README.md` - Documentação de comando

### Arquivos no Frontend
- `frontend/CYPRESS.md` - Guia de testes E2E com Cypress
- `frontend/src/utils/README.md` - Documentação de utilitários

### Arquivos em Deploy
- `deploy/README.md` - Scripts de deploy

### Arquivos já em docs/ (serão reorganizados)
- `docs/DEPLOY.md` - Guia completo de deploy
- `docs/TESTES_E2E_DOCKER.md` - Testes E2E no Docker
- `docs/ANALISE_ARQUIVOS_MD.md` - Análise de arquivos MD
- `docs/ANALISE_MELHORIAS_SENIOR.md` - Análise de melhorias
- `docs/PLANO_MELHORIAS_PENDENTES.md` - Plano de melhorias
- `docs/RELATORIO_ANALISE_ENGENHARIA_SOFTWARE.md` - Relatório de análise
- `docs/api/README.md` - Documentação da API (Postman)

---

## 📁 Estrutura Proposta

```
docs/
├── README.md (novo - índice principal)
│
├── guia/                          # Guias e tutoriais
│   ├── QUICK_START_TESTING.md    # ← QUICK_START_TESTING.md
│   └── README.md                  # Índice dos guias
│
├── configuracao/                  # Configurações e setup
│   ├── CONFIG.md                  # ← backend/CONFIG.md
│   └── README.md                  # Índice de configurações
│
├── testes/                        # Documentação de testes
│   ├── TESTING_GUIDE.md           # ← backend/TESTING_GUIDE.md
│   ├── TEST_COVERAGE_ANALYSIS.md  # ← backend/TEST_COVERAGE_ANALYSIS.md
│   ├── CYPRESS.md                 # ← frontend/CYPRESS.md
│   ├── TESTES_E2E_DOCKER.md       # ← docs/TESTES_E2E_DOCKER.md
│   └── README.md                  # Índice de testes
│
├── seguranca/                     # Documentação de segurança
│   ├── SECURITY.md                # ← backend/docs/SECURITY.md
│   └── README.md                  # Índice de segurança
│
├── analise/                       # Análises e relatórios
│   ├── ANALISE_MELHORIAS_SENIOR.md           # ← docs/ANALISE_MELHORIAS_SENIOR.md
│   ├── RELATORIO_ANALISE_ENGENHARIA_SOFTWARE.md  # ← docs/RELATORIO_ANALISE_ENGENHARIA_SOFTWARE.md
│   ├── ANALISE_ARQUIVOS_MD.md     # ← docs/ANALISE_ARQUIVOS_MD.md
│   ├── PLANO_MELHORIAS_PENDENTES.md  # ← docs/PLANO_MELHORIAS_PENDENTES.md
│   └── README.md                  # Índice de análises
│
├── verificacao/                   # Verificações históricas
│   ├── VERIFICACAO_IMPLEMENTACAO.md  # ← VERIFICACAO_IMPLEMENTACAO.md
│   ├── VERIFICACAO_TAREFAS.md    # ← VERIFICACAO_TAREFAS.md
│   └── README.md                  # Índice de verificações
│
├── deploy/                        # Deploy e scripts
│   ├── DEPLOY.md                  # ← docs/DEPLOY.md
│   ├── DEPLOY_SCRIPTS.md          # ← deploy/README.md (renomeado)
│   └── README.md                  # Índice de deploy
│
├── backend/                       # Documentação específica do backend
│   ├── COMANDOS.md                # ← backend/cmd/process-recurring/README.md
│   └── README.md                  # Índice do backend
│
├── frontend/                      # Documentação específica do frontend
│   ├── UTILITARIOS.md             # ← frontend/src/utils/README.md
│   └── README.md                  # Índice do frontend
│
├── tarefas/                       # Tarefas e planejamento
│   ├── TAREFAS.md                 # ← TAREFAS.md (raiz)
│   └── README.md                  # Índice de tarefas
│
├── api/                           # Documentação da API (mantido)
│   └── README.md                  # ← docs/api/README.md (mantido)
│
├── planejamento/                  # Planejamento (mantido)
│   └── [arquivos existentes]     # Mantidos como estão
│
└── tarefas_concluidas/            # Tarefas concluídas (mantido)
    └── [arquivos existentes]      # Mantidos como estão
```

---

## 📋 Mapeamento Detalhado

### 1. Guias (`docs/guia/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `QUICK_START_TESTING.md` | `docs/guia/QUICK_START_TESTING.md` | Guia rápido de testes |

### 2. Configuração (`docs/configuracao/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `backend/CONFIG.md` | `docs/configuracao/CONFIG.md` | Configuração da API |

### 3. Testes (`docs/testes/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `backend/TESTING_GUIDE.md` | `docs/testes/TESTING_GUIDE.md` | Guia completo de testes |
| `backend/TEST_COVERAGE_ANALYSIS.md` | `docs/testes/TEST_COVERAGE_ANALYSIS.md` | Análise de cobertura |
| `frontend/CYPRESS.md` | `docs/testes/CYPRESS.md` | Testes E2E |
| `docs/TESTES_E2E_DOCKER.md` | `docs/testes/TESTES_E2E_DOCKER.md` | Testes E2E no Docker |

### 4. Segurança (`docs/seguranca/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `backend/docs/SECURITY.md` | `docs/seguranca/SECURITY.md` | Segurança do backend |

### 5. Análise (`docs/analise/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `docs/ANALISE_MELHORIAS_SENIOR.md` | `docs/analise/ANALISE_MELHORIAS_SENIOR.md` | Análise de melhorias |
| `docs/RELATORIO_ANALISE_ENGENHARIA_SOFTWARE.md` | `docs/analise/RELATORIO_ANALISE_ENGENHARIA_SOFTWARE.md` | Relatório de análise |
| `docs/ANALISE_ARQUIVOS_MD.md` | `docs/analise/ANALISE_ARQUIVOS_MD.md` | Análise de arquivos MD |
| `docs/PLANO_MELHORIAS_PENDENTES.md` | `docs/analise/PLANO_MELHORIAS_PENDENTES.md` | Plano de melhorias |

### 6. Verificação (`docs/verificacao/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `VERIFICACAO_IMPLEMENTACAO.md` | `docs/verificacao/VERIFICACAO_IMPLEMENTACAO.md` | Verificação histórica |
| `VERIFICACAO_TAREFAS.md` | `docs/verificacao/VERIFICACAO_TAREFAS.md` | Verificação histórica |

### 7. Deploy (`docs/deploy/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `docs/DEPLOY.md` | `docs/deploy/DEPLOY.md` | Guia completo de deploy |
| `deploy/README.md` | `docs/deploy/DEPLOY_SCRIPTS.md` | Scripts de deploy (renomeado) |

### 8. Backend (`docs/backend/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `backend/cmd/process-recurring/README.md` | `docs/backend/COMANDOS.md` | Documentação de comandos |

### 9. Frontend (`docs/frontend/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `frontend/src/utils/README.md` | `docs/frontend/UTILITARIOS.md` | Utilitários do frontend |

### 10. Tarefas (`docs/tarefas/`)
| Origem | Destino | Motivo |
|--------|---------|--------|
| `TAREFAS.md` | `docs/tarefas/TAREFAS.md` | Documentação principal de tarefas |

### 11. Mantidos como estão
- `docs/api/README.md` - Mantido em `docs/api/`
- `docs/planejamento/*` - Mantidos em `docs/planejamento/`
- `docs/tarefas_concluidas/*` - Mantidos em `docs/tarefas_concluidas/`

---

## 🔄 Atualizações Necessárias

### Arquivos que Referenciam os Movidos

Após a reorganização, será necessário atualizar referências em:

1. **README.md (raiz)**
   - `[Tarefas do Projeto](./TAREFAS.md)` → `[Tarefas do Projeto](./docs/tarefas/TAREFAS.md)`
   - `[Guia de Deploy](./docs/DEPLOY.md)` → `[Guia de Deploy](./docs/deploy/DEPLOY.md)`
   - `[Configuração da API](./backend/CONFIG.md)` → `[Configuração da API](./docs/configuracao/CONFIG.md)`

2. **Outros arquivos .md**
   - Verificar e atualizar links internos que referenciam arquivos movidos

---

## ✅ Resumo

### Estatísticas
- **Total de arquivos a mover:** 15 arquivos
- **Novas pastas a criar:** 10 pastas
- **Arquivos README.md a criar:** 10 arquivos (índices)

### Categorias Criadas
1. `docs/guia/` - 1 arquivo
2. `docs/configuracao/` - 1 arquivo
3. `docs/testes/` - 4 arquivos
4. `docs/seguranca/` - 1 arquivo
5. `docs/analise/` - 4 arquivos
6. `docs/verificacao/` - 2 arquivos
7. `docs/deploy/` - 2 arquivos
8. `docs/backend/` - 1 arquivo
9. `docs/frontend/` - 1 arquivo
10. `docs/tarefas/` - 1 arquivo

### Arquivos Mantidos
- `docs/api/` - Mantido
- `docs/planejamento/` - Mantido
- `docs/tarefas_concluidas/` - Mantido

---

## 🎯 Resultado Final Esperado

Após a reorganização:
- ✅ Todos os arquivos `.md` (exceto README.md) estarão em `docs/`
- ✅ Organização clara por categoria
- ✅ Fácil navegação com README.md em cada categoria
- ✅ Estrutura escalável para futuras documentações
- ✅ Links atualizados nos arquivos principais

---

**Status:** ⏳ Aguardando aprovação para execução


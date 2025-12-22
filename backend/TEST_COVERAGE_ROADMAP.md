# Roadmap de Cobertura de Testes - Identity Context

## ✅ Status Atual (Sprint 1.3)

### Cobertura: 75.2%

| Componente | Cobertura | Status |
|------------|-----------|--------|
| Value Objects | 89.4% | ✅ Excelente |
| Services (JWT) | 88.9% | ✅ Excelente |
| Handlers (HTTP) | 87.5% | ✅ Muito Bom |
| Use Cases | 86.7% | ✅ Muito Bom |
| Entities | 81.2% | ✅ Bom |
| **Persistence (Repository)** | **0.0%** | ⚠️ **Pendente** |
| Routes | 0.0% | ⚠️ Baixa Prioridade |

---

## 📅 Planejamento de Cobertura

### ✅ Sprint 1.3 (Concluída)
- **ID-013**: Testes unitários para Identity Context ✅
- **Cobertura alcançada**: 75.2%
- **Foco**: Testes unitários de domínio, use cases, handlers e services

### ⏳ Sprint 2.7: Testes de Integração (Semana 8)

**Tarefa:** `TEST-INT-001 - Criar testes de integração para Identity Context`
- **Dependências**: ID-013 ✅ (já concluída)
- **Esforço**: 4h
- **Prioridade**: 🟡 Média
- **Status**: ⏳ Pendente

**O que será coberto:**
1. ✅ **GormUserRepository** (0% → 80%+)
   - Testes de integração com banco de dados
   - SQLite em memória ou testcontainers
   - Todos os métodos: FindByID, FindByEmail, Save, Delete, Exists, Count
   - Conversões Domain ↔ Persistence (toDomain, toModel)
   - Casos de erro do banco de dados

2. ✅ **Fluxos End-to-End**
   - Register → Login → Acesso Protegido
   - Validação de integração entre camadas

3. ✅ **Routes** (opcional)
   - Verificação de registro de rotas
   - Middleware aplicado corretamente

**Meta de cobertura após TEST-INT-001: 85%+**

---

## 🎯 Garantias do Planejamento

### ✅ SIM, os pontos de atenção SERÃO cobertos

**Evidências:**

1. **Tarefa específica planejada:**
   - `TEST-INT-001` está no Sprint 2.7 (Semana 8)
   - Dependência já satisfeita (ID-013 ✅)
   - Esforço estimado: 4h

2. **Padrão estabelecido:**
   - Cada contexto terá sua tarefa de testes de integração:
     - TEST-INT-001: Identity Context
     - TEST-INT-002: Account Context
     - TEST-INT-003: Transaction Context
     - TEST-INT-004: Category Context

3. **Sprint dedicada:**
   - Sprint 2.7 é inteiramente dedicada a testes de integração
   - Entregável: "Suite de testes de integração completa"

---

## 📊 Projeção de Cobertura

### Após TEST-INT-001 (Sprint 2.7):

| Componente | Atual | Projetado | Melhoria |
|------------|-------|-----------|----------|
| Value Objects | 89.4% | 90%+ | +0.6% |
| Services | 88.9% | 90%+ | +1.1% |
| Handlers | 87.5% | 90%+ | +2.5% |
| Use Cases | 86.7% | 90%+ | +3.3% |
| Entities | 81.2% | 90%+ | +8.8% |
| **Persistence** | **0.0%** | **80%+** | **+80%** |
| Routes | 0.0% | 70%+ | +70% |

### Cobertura Total Projetada: **85-90%**

---

## ⏰ Timeline

```
Sprint 1.3 (Semana 2) ✅
├── Testes unitários básicos
└── Cobertura: 75.2%

Sprint 2.7 (Semana 8) ⏳
├── TEST-INT-001: Testes de integração Identity
└── Cobertura projetada: 85%+

Sprint 4.4 (Semana 15-16) ⏳
├── Testes E2E completos
└── Cobertura final: 90%+
```

---

## ✅ Conclusão

**SIM, os pontos de atenção serão cobertos até o final do desenvolvimento.**

### Garantias:

1. ✅ **Tarefa específica planejada** (TEST-INT-001)
2. ✅ **Dependências satisfeitas** (ID-013 concluída)
3. ✅ **Sprint dedicada** (Sprint 2.7)
4. ✅ **Padrão estabelecido** (todos os contextos terão testes de integração)
5. ✅ **Meta clara** (85%+ de cobertura)

### Recomendação:

A cobertura atual (75.2%) é **suficiente para desenvolvimento ativo**, e os testes de integração estão **planejados e garantidos** para a Sprint 2.7.

**Não é necessário antecipar** os testes de integração agora, pois:
- O código está em desenvolvimento ativo
- Os testes unitários cobrem a lógica de negócio
- Os testes de integração serão feitos quando o sistema estiver mais estável
- O planejamento já contempla essa necessidade

---

## 📝 Nota sobre Priorização

A estratégia atual (testes unitários primeiro, integração depois) é **correta** porque:

1. **Desenvolvimento rápido**: Permite iterar rápido sem depender de banco
2. **Isolamento**: Testes unitários são mais rápidos e isolados
3. **Custo-benefício**: Testes de integração são mais caros e demorados
4. **Padrão da indústria**: Testes unitários → Integração → E2E


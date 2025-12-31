# Documentação de Planejamento

Esta pasta contém toda a documentação de planejamento do sistema de gestão financeira, incluindo planejamentos específicos para diferentes stacks tecnológicos.

## Estrutura

```
docs/planejamento/
├── README.md (este arquivo)
├── PLANEJAMENTO.md (planejamento geral agnóstico de tecnologia)
├── ANALISE_ARQUIVOS_PLANEJAMENTO.md (análise dos arquivos de planejamento)
├── GO/
│   ├── README.md
│   ├── PLANEJAMENTO_GO.md
│   ├── EXPLICACAO_GO.md
│   ├── RECEIPT_SCANNING.md
│   └── RECEIPT_SCANNING_COMPARISON.md
├── NODE/
│   ├── README.md
│   ├── PLANEJAMENTO_NODE.md
│   └── EXPLICACAO_NODE.md
└── PHP/
    ├── README.md
    ├── PLANEJAMENTO_PHP.md
    └── EXPLICACAO_PHP.md
```

## Arquivos Principais

### Planejamento Geral

- **PLANEJAMENTO.md** - Planejamento geral agnóstico de tecnologia, definindo a arquitetura DDD, bounded contexts, entidades, value objects e estrutura de camadas de forma genérica.

### Análise

- **ANALISE_ARQUIVOS_PLANEJAMENTO.md** - Análise dos arquivos de planejamento existentes, determinando a necessidade de criar novos arquivos e comparando o nível de detalhe entre as diferentes stacks.

### Planejamentos Específicos por Stack

Cada pasta (GO, NODE, PHP) contém:

1. **PLANEJAMENTO_[STACK].md** - Planejamento completo e detalhado para a stack específica, incluindo:
   - Stack tecnológico completo
   - Arquitetura DDD adaptada para a stack
   - Estrutura de pastas detalhada
   - Exemplos de código práticos
   - Fases de desenvolvimento
   - Performance e otimizações
   - Observabilidade
   - Segurança
   - Deploy e DevOps
   - Testes
   - Versionamento de API
   - Auditoria e compliance
   - Multi-tenancy

2. **EXPLICACAO_[STACK].md** - Explicação e resumo do planejamento, destacando:
   - Visão geral
   - Principais seções
   - Stack tecnológico
   - Por que escolher a stack
   - Arquitetura DDD
   - Estrutura de pastas
   - Exemplos de código
   - Fases de desenvolvimento
   - Destaques do documento

## Comparação entre Stacks

| Aspecto | Go | Node.js | PHP |
|---------|----|---------|-----|
| **Framework** | Fiber | NestJS | Laravel/Symfony |
| **ORM** | GORM | Prisma | Eloquent/Doctrine |
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Produtividade** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Type Safety** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Ecossistema** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Curva de Aprendizado** | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **DDD Nativo** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |

## Como Usar

1. **Para entender a arquitetura geral**: Leia `PLANEJAMENTO.md`
2. **Para escolher uma stack**: Leia os `EXPLICACAO_[STACK].md` de cada pasta
3. **Para implementar**: Use o `PLANEJAMENTO_[STACK].md` correspondente como guia completo


## Status do Projeto

O projeto atual está implementado em **Go** (conforme código existente). Os planejamentos para Node.js e PHP foram criados para:
- Comparação entre stacks
- Possível migração futura
- Aprendizado e referência
- Decisões arquiteturais


# Documentação de Planejamento - Python

Esta pasta contém a documentação de planejamento específica para a implementação do sistema em **Python**.

## Arquivos

- **PLANEJAMENTO_PYTHON.md** - Planejamento completo e detalhado para implementação em Python com FastAPI, SQLAlchemy, seguindo DDD
- **EXPLICACAO_PYTHON.md** - Explicação e resumo do PLANEJAMENTO_PYTHON.md, destacando os principais pontos e estrutura

## Stack Tecnológico

- **Python 3.11+** com **FastAPI** (framework web assíncrono)
- **PostgreSQL** (banco de dados)
- **Redis** (cache e rate limiting)
- **SQLAlchemy 2.0+** (ORM async) ou Tortoise ORM
- **Pydantic** (validação type-safe)
- **OpenTelemetry** (observabilidade)
- **Prometheus + Grafana** (monitoramento)
- **MinIO/S3** (armazenamento de arquivos)
- **Google Cloud Vision API / Tesseract** (OCR)
- **OpenAI GPT-4** (análise de IA)
- **N8N** (workflow automation - opcional)

## Estrutura

O planejamento inclui:
- Arquitetura DDD completa
- Estrutura de pastas detalhada
- Exemplos de código práticos
- Fases de desenvolvimento (5 fases, 12-15 semanas)
- Performance e otimizações (async/await, connection pooling, cache)
- Observabilidade completa
- Segurança robusta
- Deploy e DevOps (Docker, docker-compose)
- Testes (unitários, integração, E2E)
- **Feature completa de upload e análise automática de comprovantes**

## Funcionalidades Especiais

### Upload e Análise Automática de Comprovantes

**Feature integrada no planejamento:** Upload de comprovantes (imagens de recibos, notas fiscais, extratos) e análise automática via OCR/IA para extrair informações e criar transações automaticamente.

**Fluxo:**
1. Upload de arquivo (JPG, PNG, PDF)
2. Validação e armazenamento (MinIO/S3)
3. Processamento assíncrono (OCR + IA)
4. Extração de dados (valor, data, descrição, tipo)
5. Criação de Transaction DRAFT
6. Revisão e confirmação do usuário
7. Aplicação da transação

**Opções de Implementação:**
- **Opção 1**: Processamento direto (Python nativo com Google Vision + OpenAI)
- **Opção 2**: N8N (workflow automation visual)
- **Recomendação**: Opção 2 (N8N) para MVP, Opção 1 para produção

## Compatibilidade com Frontend

**✅ IMPORTANTE:** O projeto já possui um frontend Vue 3 completamente funcional e independente. **NÃO é necessário criar novo frontend** para Python.

O frontend Vue 3 é **reutilizável** sem modificações, apenas configurando a URL da API. Veja seção 3.3 do PLANEJAMENTO_PYTHON.md para detalhes de compatibilidade.

## Diferenciais da Stack Python/FastAPI

- ✅ **Produtividade excepcional**: Código limpo e expressivo
- ✅ **Type-safety**: Type hints nativos do Python
- ✅ **Performance assíncrona**: FastAPI é uma das opções mais rápidas do Python
- ✅ **Documentação automática**: Swagger gerado automaticamente
- ✅ **Validação automática**: Pydantic integrado
- ✅ **Ecossistema rico**: Muitas bibliotecas disponíveis (especialmente para ML/IA)
- ✅ **ML/IA**: Excelente para integração com serviços de IA (OCR, LLMs)
- ✅ **Fácil integração**: Integração simples com Google Vision, OpenAI, etc.

## Próximos Passos

1. Setup inicial do projeto (FastAPI, SQLAlchemy, estrutura DDD)
2. Implementação do Identity Context
3. Desenvolvimento incremental dos outros contexts
4. Integração da feature de comprovantes
5. Testes e otimizações
6. Deploy e monitoramento


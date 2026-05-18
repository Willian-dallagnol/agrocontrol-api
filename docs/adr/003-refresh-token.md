# ADR 003 — Access Token + Refresh Token

**Status:** Aceito  
**Data:** 2025

## Contexto

O sistema original usava tokens JWT com expiração de 24h. Qualquer token vazado daria acesso por até 24h sem possibilidade de revogação.

## Decisão

Adotamos o padrão de dois tokens:
- **Access token:** curta duração (1h configurável), usado nas requisições protegidas
- **Refresh token:** longa duração (7 dias), usado apenas no endpoint `POST /auth/refresh`

Ambos são JWT assinados com HS256. O refresh token tem o campo `is_refresh: true` nas claims, impedindo seu uso como access token.

## Consequências

- ✅ Janela de exposição reduzida de 24h para 1h
- ✅ Sem impacto na experiência — clientes renovam automaticamente
- ✅ Sem estado no servidor (sem tabela de refresh tokens)
- ⚠️ Sem revogação imediata (limitação do JWT stateless — aceitável para este contexto)

# ADR 005 — Security Headers HTTP

**Status:** Aceito  
**Data:** 2026

## Contexto

A API não retornava headers de segurança HTTP nas respostas, deixando clientes vulneráveis a ataques comuns como clickjacking, XSS, sniffing de MIME type e exposição de informações sobre a stack tecnológica.

## Decisão

Criamos o middleware `SecurityHeaders()` aplicado globalmente em todas as respostas da API, com os seguintes headers:

| Header | Valor | Proteção |
|--------|-------|----------|
| `X-Frame-Options` | `DENY` | Clickjacking |
| `X-XSS-Protection` | `1; mode=block` | XSS em browsers legados |
| `X-Content-Type-Options` | `nosniff` | MIME sniffing |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Vazamento de URLs internas |
| `Permissions-Policy` | `geolocation=(), microphone=(), camera=()` | Acesso a recursos do browser |
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains` | Força HTTPS por 1 ano |
| `Content-Security-Policy` | `default-src 'self'` | Carregamento de recursos |
| `Server` | `""` | Oculta stack tecnológica |

## Consequências

- ✅ Proteção contra os ataques mais comuns de segurança web
- ✅ Score A+ em ferramentas de análise de headers (securityheaders.com)
- ✅ Zero impacto em performance — headers adicionados antes do handler
- ⚠️ CSP com `unsafe-inline` para scripts e estilos — necessário para o Swagger UI funcionar
- ⚠️ HSTS deve ser usado apenas com HTTPS — em desenvolvimento local pode causar problemas se acessado via HTTP

## Alternativa considerada

Usar biblioteca externa como `secure` do unrolled. Descartado para manter zero dependências adicionais — a implementação manual é simples e totalmente controlada.
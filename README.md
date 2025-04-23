# 🧠 RESTful API em Golang

Este projeto é uma API RESTful escrita em Go, onde decidi aprofundar meus conhecimentos no padrão REST proposto por Roy Fielding bem como aprofundar meu conhecimento em Golang. A API Implementa até o nível 3 da Escala de Maturidade de Richardson e inclui caching com ETags para otimizar o tráfego entre cliente e servidor.

---

## ✅ Recursos já implementados

### 📊 Escala de Maturidade de Richardson

| Nível | Descrição | Status |
|------:|-----------|--------|
|  Level 1 | Recursos com URIs distintas e orientadas a recursos (`/categories`, `/products`) | ✅ |
|  Level 2 | Uso correto de verbos HTTP (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`) e status codes (`200`, `201`, `304`, etc), além de outras boas práticas no uso do HTTP como headers, content-negotiation etc. | ✅ |
|  Level 3 | HATEOAS (Hypermedia links nas respostas) | ✅ |

### 🗃️ Caching via ETag

- Geração de ETag a partir do hash SHA256 do corpo da resposta e codificado depois em hexadecimal.
- Verificação automática com `If-None-Match`
- Retorno `304 Not Modified` quando aplicável
- Redução de consumo de banda e carga de processamento

### 🪵 Logs estruturados utilizando [slog](https://github.com/sirupsen/logrus) com rotação automática usando [Lumberjack](https://github.com/natefinch/lumberjack)

### 🔐 Autenticação com JWT

- Implementado fluxo de autenticação via JWT 
- Atualmente implementado
  - Cadastro de usuários
  - Rotas para login e atualização do par de tokens
  - Middleware para validação do token nas rotas administrativas

---

## 📌 Em andamento / Próximos passos

- 🚀 **Deploy como serviço (Windows/Linux)**  
  Utilizar o [Kardianos/service](https://github.com/kardianos/service) para rodar a API como serviço nativo do sistema operacional.

---

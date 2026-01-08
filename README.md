# Network Toolkit 🔧

Canivete suíço para atividades de gerenciamento de redes desenvolvido em Go.

## 📋 Descrição

Network Toolkit é uma aplicação de linha de comando que fornece ferramentas avançadas para administradores de sistemas e profissionais de segurança gerenciarem, monitorarem e auditarem conexões de rede. A aplicação oferece uma interface interativa e fácil de usar, com funcionalidades equivalentes ao nmap, netstat e outras ferramentas de rede essenciais.

## ✨ Funcionalidades Implementadas

### 1. Listar Portas em Escuta
Alternativa ao comando `netstat -tuln` (Linux) ou `Get-NetTCPConnection -State Listen` (PowerShell).

Exibe todas as portas TCP em estado de escuta com:
- ✅ Endereço local
- ✅ Porta
- ✅ Estado da conexão
- ✅ PID do processo
- ✅ Nome do processo

**Funções Auxiliares:**
- `GetListeningPortsCount()` - Retorna o número de portas em escuta
- `IsPortListening(port)` - Verifica se uma porta específica está em escuta
- `GetProcessByPort(port)` - Retorna o processo que está usando uma porta

### 2. Scanner de Rede (nmap -sS -sV -p-)
Scanner de rede completo para múltiplos hosts em notação CIDR.

Funcionalidades:
- ✅ Parse de redes CIDR (ex: 192.168.1.0/24)
- ✅ Detecção automática de hosts ativos
- ✅ Scan paralelo de portas TCP
- ✅ Identificação de 20+ serviços comuns
- ✅ Banner grabbing para detecção avançada
- ✅ Configuração de threads (1-100)
- ✅ Múltiplas opções de range de portas
- ✅ Relatório detalhado com estatísticas

**Opções de Portas:**
- Portas comuns (~20 portas principais)
- Range específico (ex: 1-1024)
- Portas customizadas (ex: 80,443,8080)

### 3. Scanner Stealth de Host Único (nmap -sS -sV -p- -T4 --reason)
Scanner agressivo focado em um único alvo com máximo desempenho.

Funcionalidades:
- ✅ TCP SYN Scan (stealth mode)
- ✅ Detecção de versão de serviços (-sV)
- ✅ Timing agressivo T4 (até 200 threads)
- ✅ Análise de motivos (--reason): syn-ack, conn-refused, timeout
- ✅ Estados de porta: open, closed, filtered
- ✅ Banner grabbing com extração de versão
- ✅ Progresso em tempo real
- ✅ Estimativa de tempo antes do scan

**Modos de Scan:**
- **Rápido**: Portas 1-1024 (~20 segundos)
- **Completo**: Todas as 65535 portas (~5-10 minutos)
- **Personalizado**: Range definido pelo usuário

## 🚀 Instalação

### Pré-requisitos
- Go 1.21 ou superior
- Privilégios de administrador (recomendado para visualizar todos os processos)

### Compilar

```bash
# Navegue até o diretório do projeto
cd network-toolkit

# Baixe as dependências
go mod download

# Compile o executável
go build -o network-toolkit.exe
```

## 💻 Uso

### Executar a Aplicação

```bash
# Windows (recomendado: executar como Administrador)
.\network-toolkit.exe
```

### Menu Interativo
A aplicação apresenta um menu interativo:

```
============================================================
  Network Toolkit 🔧 - v1.2.0
  Canivete suíço para atividades de gerenciamento de redes
============================================================

------------------------------------------------------------
MENU PRINCIPAL
------------------------------------------------------------
[1] Listar Portas em Escuta (netstat -tuln)
[2] Scanner de Rede (nmap -sS -sV -p-)
[3] Scanner Stealth de Host Único (nmap -sS -sV -p- -T4)
[0] Sair
------------------------------------------------------------
```

### Exemplo de Saída - Portas em Escuta

```
=== PORTAS EM ESCUTA ===
ENDEREÇO             PORTA      ESTADO          PID        PROCESSO
--------------------------------------------------------------------------------------------
0.0.0.0              80         LISTEN          1234       nginx.exe
0.0.0.0              443        LISTEN          1234       nginx.exe
127.0.0.1            3306       LISTEN          5678       mysqld.exe
0.0.0.0              8080       LISTEN          9012       java.exe

Total: 4 porta(s) em escuta
```

### Exemplo de Saída - Scanner de Rede

```
🔍 Iniciando scan de rede: 192.168.1.0/24
📊 Hosts a escanear: 254
🔌 Portas por host: 20
⚙️  Threads: 10

✅ 192.168.1.1 - 4 porta(s) aberta(s)
✅ 192.168.1.20 - 6 porta(s) aberta(s)

================================================================================
📊 RELATÓRIO DE SCAN DE REDE
================================================================================

🖥️  HOST: 192.168.1.1 (router.local)
   Tempo de scan: 2.3s
   🔓 Portas abertas: 4

   PORTA      SERVIÇO              BANNER
   ----------------------------------------------------------------------
   80         HTTP                 nginx/1.18.0
   443        HTTPS                
   22         SSH                  OpenSSH_8.2p1
   8080       HTTP-Proxy           
```

### Exemplo de Saída - Scanner Stealth

```
🎯 TARGET: 192.168.1.20 (server.local)
🔍 Scanning 65535 ports (range: 1-65535)
⚙️  Threads: 100 | Timeout: 1s | Timing: Aggressive (T4)

✅ Port 22/tcp      open    SSH
✅ Port 80/tcp      open    HTTP
✅ Port 443/tcp     open    HTTPS
⏳ Progresso: 25% (16384/65535 portas escaneadas)

================================================================================
🎯 RELATÓRIO DE SCAN STEALTH (NMAP-LIKE)
================================================================================

📍 TARGET: 192.168.1.20 (server.local)
⏱️  Duração: 5m 23s

📊 ESTATÍSTICAS
   🟢 Abertas:   8
   🔴 Fechadas:  65520
   🟡 Filtradas: 7

🔓 PORTAS ABERTAS DETECTADAS
PORTA      ESTADO     SERVIÇO         RAZÃO                VERSÃO/BANNER
----------------------------------------------------------------------------------
22         open       SSH             syn-ack              OpenSSH_8.2p1 Ubuntu
80         open       HTTP            syn-ack              nginx/1.18.0
443        open       HTTPS           syn-ack              nginx/1.18.0
3306       open       MySQL           syn-ack              MySQL 8.0.28
```

## 📁 Estrutura do Projeto

```
network-toolkit/
├── main.go                          # Entrada da aplicação e menu interativo
├── network/
│   ├── listening_ports.go           # Módulo de portas em escuta
│   ├── port_scanner.go              # Scanner de rede CIDR
│   └── port_scanner_stealthy.go     # Scanner stealth de host único
├── go.mod                           # Gerenciamento de dependências
├── go.sum                           # Checksums das dependências
├── .gitignore                       # Arquivos ignorados pelo Git
├── network-toolkit.exe              # Executável compilado
└── README.md                        # Este arquivo
```

## 📦 Dependências

- [`github.com/shirou/gopsutil/v3`](https://github.com/shirou/gopsutil) - Biblioteca para obter informações de sistema, processos e rede de forma multiplataforma

## 📝 Notas Importantes

### Windows
- **Privilégios de Administrador**: Execute o programa como Administrador para visualizar informações completas de todos os processos
- **Windows Defender/Antivírus**: Algumas soluções de segurança podem alertar sobre o executável. Isso é normal para ferramentas de rede.

### Compatibilidade
- ✅ Windows 10/11
- ✅ Windows Server 2016+
- ⚠️ Linux (funcionalidade básica - necessita testes)
- ⚠️ macOS (funcionalidade básica - necessita testes)

### ⚠️ Avisos de Segurança e Uso Ético

**IMPORTANTE**: As funcionalidades de scan de rede devem ser utilizadas apenas:
- Em redes e sistemas que você possui ou tem autorização explícita
- Para fins de auditoria de segurança legítima
- Em ambientes de teste e desenvolvimento próprios

**Uso não autorizado pode:**
- Violar leis de crimes cibernéticos
- Resultar em ações legais
- Ser detectado por sistemas IDS/IPS
- Gerar alertas de segurança

**Recomendações:**
- Sempre obtenha autorização por escrito antes de escanear redes
- Use em horários de baixo movimento quando possível
- Configure threads e timeouts apropriados
- Mantenha logs de atividades de scan
- Respeite políticas de segurança da informação

### Limitações Conhecidas
- Processos do sistema protegidos podem aparecer como "Unknown" sem privilégios administrativos
- A performance pode variar dependendo do número de conexões ativas no sistema
- Scanner stealth usa TCP connect scan (não SYN real) devido a limitações do Go
- Detecção de OS é limitada (não implementada completamente)
- Suporte apenas para IPv4 no momento
- Firewalls podem bloquear ou limitar scans de rede

## 🗺️ Roadmap

### ✅ Versão 1.1.0 (Concluída)
- [x] Scanner de rede com suporte a CIDR
- [x] Detecção de hosts ativos
- [x] Scan paralelo de portas TCP
- [x] Identificação de serviços comuns
- [x] Banner grabbing básico

### ✅ Versão 1.2.0 (Concluída)
- [x] Scanner stealth de host único
- [x] Timing agressivo (T4)
- [x] Detecção de versão de serviços
- [x] Análise de motivos (--reason)
- [x] Estados de porta (open/closed/filtered)
- [x] Progresso em tempo real

### Versão 1.3.0 (Em Planejamento)
- [ ] Adicionar suporte para portas UDP
- [ ] Implementar filtros (por porta, por processo, por endereço)
- [ ] Adicionar opção de exportar resultados para CSV/JSON
- [ ] Melhorar tratamento de erros e mensagens ao usuário
- [ ] Listar todas as conexões ativas (não apenas LISTEN)

### Versão 2.0.0
- [ ] Teste de conectividade (ping, traceroute)
- [ ] Análise de latência e jitter
- [ ] Interface web opcional (modo servidor)
- [ ] Suporte a IPv6 completo
- [ ] Detecção de OS (fingerprinting)
- [ ] Modo de monitoramento contínuo

### Futuras Funcionalidades
- [ ] Monitoramento de largura de banda por processo
- [ ] Alertas e notificações
- [ ] Histórico de conexões
- [ ] Detecção de conexões suspeitas
- [ ] Integração com ferramentas de logging
- [ ] API REST para integração com outras ferramentas
- [ ] Modo daemon/serviço para monitoramento contínuo

## 🐛 Problemas Conhecidos

Nenhum problema crítico identificado até o momento.

## 🤝 Contribuindo

Sugestões e melhorias são bem-vindas! Este projeto está em desenvolvimento ativo.

### Como Contribuir
1. Identifique um bug ou funcionalidade desejada
2. Implemente a solução
3. Teste em diferentes cenários
4. Documente as mudanças

## 📄 Licença

Este projeto é de uso interno e educacional.

## 👨‍💻 Desenvolvimento

### Tecnologias Utilizadas
- **Linguagem**: Go 1.21+
- **Bibliotecas**: gopsutil v3
- **Plataforma**: Windows (primário)

### Status do Projeto
🟢 Em desenvolvimento ativo - v1.2.0

### Última Atualização
8 de Janeiro de 2026

### Histórico de Versões
- **v1.2.0** (08/01/2026) - Scanner Stealth de Host Único
- **v1.1.0** (07/01/2026) - Scanner de Rede CIDR
- **v1.0.1** (07/01/2026) - Ajustes intermediários
- **v1.0.0** (07/01/2026) - Release inicial

---

**Network Toolkit** - Simplificando o gerenciamento de redes 🚀

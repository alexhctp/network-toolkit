# Network Toolkit 🔧

Canivete suíço para atividades de gerenciamento de redes desenvolvido em Go.

## 📋 Descrição

Network Toolkit é uma aplicação de linha de comando que fornece ferramentas úteis para administradores de sistemas e desenvolvedores gerenciarem e monitorarem conexões de rede. A aplicação oferece uma interface interativa e fácil de usar.

## ✨ Funcionalidades Implementadas

### 1. Listar Portas em Escuta
Alternativa ao comando `netstat -tuln` (Linux) ou `Get-NetTCPConnection -State Listen` (PowerShell).

Exibe todas as portas TCP em estado de escuta com:
- ✅ Endereço local
- ✅ Porta
- ✅ Estado da conexão
- ✅ PID do processo
- ✅ Nome do processo

### Funções Auxiliares Disponíveis
- `GetListeningPortsCount()` - Retorna o número de portas em escuta
- `IsPortListening(port)` - Verifica se uma porta específica está em escuta
- `GetProcessByPort(port)` - Retorna o processo que está usando uma porta

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
================================================================
  Network Toolkit 🔧 - v1.0.0
  Canivete suíço para atividades de gerenciamento de redes
================================================================

------------------------------------------------------------
MENU PRINCIPAL
------------------------------------------------------------
[1] Listar Portas em Escuta (netstat -tuln)
[0] Sair
------------------------------------------------------------
```

### Exemplo de Saída

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

## 📁 Estrutura do Projeto

```
network-toolkit/
├── main.go                      # Entrada da aplicação e menu interativo
├── network/
│   └── listening_ports.go       # Módulo de portas em escuta
├── go.mod                       # Gerenciamento de dependências
├── go.sum                       # Checksums das dependências
├── network-toolkit.exe          # Executável compilado
└── README.md                    # Este arquivo
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

### Limitações Conhecidas
- Processos do sistema protegidos podem aparecer como "Unknown" sem privilégios administrativos
- A performance pode variar dependendo do número de conexões ativas no sistema

## 🗺️ Roadmap

### Versão 1.1.0 (Próxima Release)
- [ ] Adicionar suporte para portas UDP
- [ ] Implementar filtros (por porta, por processo, por endereço)
- [ ] Adicionar opção de exportar resultados para CSV/JSON
- [ ] Melhorar tratamento de erros e mensagens ao usuário

### Versão 1.2.0
- [ ] Listar todas as conexões ativas (não apenas LISTEN)
- [ ] Adicionar estatísticas de rede (bytes enviados/recebidos)
- [ ] Implementar modo de monitoramento contínuo (refresh automático)
- [ ] Adicionar gráficos ASCII de uso de rede

### Versão 2.0.0
- [ ] Scanner de portas (verificar se portas remotas estão abertas)
- [ ] Teste de conectividade (ping, traceroute)
- [ ] Análise de latência e jitter
- [ ] Interface web opcional (modo servidor)
- [ ] Suporte a IPv6 completo

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
🟢 Em desenvolvimento ativo - v1.0.0

### Última Atualização
7 de Janeiro de 2026

---

**Network Toolkit** - Simplificando o gerenciamento de redes 🚀
# network-toolkit

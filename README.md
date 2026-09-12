# scaffolder-go

**scaffolder-go** é uma ferramenta de linha de comando escrita em Go que permite criar novos projetos a partir de templates de frameworks e ferramentas já configurados, sem precisar lembrar comandos, flags ou passos de instalação manuais.

A ideia central é simples: em vez de manter dezenas de templates à mão que ficam desatualizados com o tempo, cada framework é adicionado a este repositório executando **o comando oficial de criação do próprio framework** (o mesmo recomendado pelos seus mantenedores). O scaffolder-go se encarrega de copiar essa estrutura para o projeto do usuário e instalar as dependências automaticamente.

## ✨ Funcionalidades

- 📦 Scaffolding de projetos para múltiplos frameworks/linguagens/ferramentas a partir de uma única CLI.
- 🧱 Cada template vive na sua própria pasta dentro de `structure/`, gerada com o comando oficial do framework — não são templates "feitos à mão".
- ⚙️ Instalação automática de dependências ao criar o projeto (via `setup.go`).
- 🧹 Limpeza automática: os arquivos internos usados no setup não ficam no projeto final do usuário.
- 🤝 Pensado para crescer com contribuições da comunidade: adicionar um framework novo não exige alterar o código da CLI.

## 🚀 Instalação

### Opção 1: Baixar o binário (recomendado)

Baixe o binário pré-compilado para o seu sistema operacional na seção [Releases](https://github.com/<seu-usuario>/scaffolder-go/releases). Há builds disponíveis para:

- 🪟 Windows
- 🐧 Linux
- 🍎 macOS

Baixe o arquivo correspondente à sua plataforma, descompacte e adicione o binário ao seu `PATH`.

### Opção 2: Com `go install`

```bash
go install github.com/<seu-usuario>/scaffolder-go@latest
```

### Opção 3: Compilando a partir do código-fonte

```bash
git clone https://github.com/<seu-usuario>/scaffolder-go.git
cd scaffolder-go
go build -o scaffolder .
```

## 🧑‍💻 Uso

```bash
scaffolder create <framework> <nome-do-projeto>
```

Exemplo:

```bash
scaffolder create deno-fresh minha-app
```

Isso vai:
1. Procurar a pasta correspondente dentro de `structure/`.
2. Copiar o conteúdo dela para `./minha-app`.
3. Executar o `setup.go` interno para instalar as dependências do framework escolhido.
4. Remover os arquivos internos de scaffolding (`setup.go`) para que o projeto fique limpo, como se tivesse sido criado manualmente com o comando oficial.

## 📁 Estrutura do repositório

```
scaffolder-go/
├── structure/
│   ├── <framework-1>/
│   │   ├── setup.go
│   │   └── ... (arquivos gerados pelo comando oficial do framework)
│   ├── <framework-2>/
│   │   ├── setup.go
│   │   └── ...
│   └── ...
├── main.go
└── README.md
```

Cada pasta dentro de `structure/` representa **um framework ou ferramenta suportado**. O nome da pasta é o identificador usado ao chamar `scaffolder create <nome>`.

### Como funciona o `setup.go`?

Cada template inclui, na raiz da sua pasta, um arquivo `setup.go` responsável por instalar as dependências do projeto gerado, usando o comando de instalação próprio daquele framework (`npm install`, `deno install`, `cargo build`, `pip install -r requirements.txt`, etc.).

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Setting up project dependencies...")

	cmd := exec.Command("deno", "install") // <-- substituir pelo comando real de instalação
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: Failed to install dependencies.")
		os.Exit(1)
	}

	// Se autodeleta para não deixar rastro na pasta final do usuário
	_ = os.Remove("setup.go")

	fmt.Println("Setup completed successfully!")
}
```

Ao final, o `setup.go` se autodeleta: o projeto que fica para o usuário final é indistinguível de um criado rodando o comando oficial do framework manualmente.

## 🧩 Frameworks suportados

| Framework / Ferramenta | Comando de instalação usado |
|---|---|
| _(adicionar aqui conforme novos templates forem incluídos)_ | |

## ⚙️ CI/CD — Workflows

O repositório tem dois workflows do GitHub Actions em `.github/workflows/`, cada um com um gatilho (`trigger`) diferente e totalmente independentes entre si.

### `release.yml` — Compilar e publicar binários

**É ativado quando:** é feito push de uma **tag** que começa com `v` (`on: push: tags: ['v*']`).

**Qual o fluxo tradicional para disparar esse workflow?**

```bash
# 1. Adicionar os arquivos alterados
git add .

# 2. Commitar seguindo Conventional Commits
git commit -m "feat: primeiro release do scaffolder-go"

# 3. Enviar o commit para a branch main
git push origin main

# 4. Criar a tag da versão (segue SemVer: vMAJOR.MINOR.PATCH)
git tag v0.1.0

# 5. Enviar a tag — é esse push que dispara o release.yml
git push origin v0.1.0
```

Regra prática para escolher o número da versão, baseada em Conventional Commits:

| Prefixo do commit | Tipo de mudança | Exemplo de bump |
|---|---|---|
| `fix: ...` | correção | `v1.0.0` → `v1.0.1` |
| `feat: ...` | nova funcionalidade | `v1.0.0` → `v1.1.0` |
| `feat!: ...` ou commit com `BREAKING CHANGE` | mudança incompatível | `v1.0.0` → `v2.0.0` |

**O que o workflow faz:**
1. Faz checkout do repositório e configura o Go 1.22.
2. Compila o binário para 5 combinações de sistema operacional/arquitetura:
   - `windows/amd64`
   - `linux/amd64`
   - `linux/arm64`
   - `darwin/amd64` (macOS Intel)
   - `darwin/arm64` (macOS Apple Silicon)
3. Cria um GitHub Release com o mesmo nome da tag e sobe os 5 binários como assets, usando `softprops/action-gh-release`.
4. Gera as release notes automaticamente (`generate_release_notes: true`) a partir dos commits e PRs desde a tag anterior.

**Requisitos:** nenhum especial — usa o `GITHUB_TOKEN` implícito do repositório (por isso tem `permissions: contents: write`).

### `sync-structure.yml` — Sincronizar templates com o R2

**É ativado quando:** é feito push para a branch `main` **e** o push modifica algo dentro de `structure/**` (`on: push: branches: [main]: paths: ['structure/**']`). Um push que altere apenas o README, o CLI em `main.go` ou qualquer outro arquivo fora de `structure/` **não** dispara esse workflow.

**Fluxo tradicional para disparar esse workflow:**

```bash
# Depois de adicionar/editar algo dentro de structure/<framework>/
git add structure/nestjs
git commit -m "feat: adicionar template do NestJS"
git push origin main
```

Basta esse `git push origin main` — não é necessário criar tag nenhuma para este workflow rodar.

**O que o workflow faz:**
1. Instala o `rclone` e o configura para falar com o Cloudflare R2 usando o provedor S3.
2. Empacota **cada pasta** de `structure/<framework>/` como um `.tar.gz` (`structure/nestjs/` → `dist/templates/nestjs/latest.tar.gz`).
3. Gera um `manifest.json` com a data de atualização e a lista de todos os templates disponíveis.
4. Sincroniza a pasta `dist/` inteira com o bucket do R2 (`rclone sync dist/ r2:scaffolder-dist/`).

É isso que permite à CLI do scaffolder-go baixar templates sob demanda em vez de tê-los empacotados dentro do binário.

**Requisitos:** três secrets configurados no repositório (Settings → Secrets and variables → Actions):

| Secret | Para que serve |
|---|---|
| `R2_ACCESS_KEY_ID` | Access key do bucket do Cloudflare R2 |
| `R2_SECRET_ACCESS_KEY` | Secret key do bucket |
| `R2_ENDPOINT` | Endpoint S3 do bucket do R2 |

### Resumo: o que dispara o quê

| Ação | `release.yml` | `sync-structure.yml` |
|---|:---:|:---:|
| Push para `main` sem alterar `structure/` | ❌ | ❌ |
| Push para `main` alterando `structure/**` | ❌ | ✅ |
| Push de uma tag `vX.Y.Z` | ✅ | ❌ |
| Commit com mudanças em `structure/` + push para `main` + tag depois | ✅ (pela tag) | ✅ (pelo commit em main) |

Os dois workflows são independentes e podem rodar em paralelo no mesmo ciclo de mudanças. Por exemplo, ao adicionar um template novo e depois lançar uma versão:

```bash
git add structure/nestjs
git commit -m "feat: adicionar template do NestJS"
git push origin main          # dispara sync-structure.yml

git tag v0.2.0
git push origin v0.2.0        # dispara release.yml
```

## 🤝 Contribuir

Quer adicionar um framework novo? Veja o guia completo em [`CONTRIBUTING.md`](./CONTRIBUTING.md). Em resumo:

1. Criar uma pasta em `structure/<nome-do-framework>` (kebab-case).
2. Dentro dela, executar o comando oficial de criação do framework (o mesmo recomendado na documentação oficial para começar um projeto do zero).
3. Adicionar um `setup.go` na raiz dessa pasta com o comando real de instalação de dependências.
4. Abrir um Pull Request.

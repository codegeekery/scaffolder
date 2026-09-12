# Contribuir com o scaffolder-go

Obrigado por querer adicionar um framework ou ferramenta nova! Adicionar suporte a algo novo no scaffolder-go **não exige alterar o código da CLI**, basta seguir a convenção de pastas descrita abaixo.

## Regra de ouro

> A pasta que você adicionar em `structure/` precisa ser o resultado real de rodar o comando oficial de criação do framework — não um template montado manualmente. Assim garantimos que o que os usuários do scaffolder-go recebem é idêntico ao que obteriam seguindo a documentação oficial do framework.

## Passos para adicionar um framework novo

### 1. Criar a pasta

Dentro de `structure/`, crie uma pasta com o nome do framework ou ferramenta, em **kebab-case** e em minúsculas:

```
structure/deno-fresh/
structure/vite-react/
structure/nestjs/
```

### 2. Gerar o projeto com o comando oficial

Dentro dessa pasta, execute o comando que o próprio framework recomenda para iniciar um projeto do zero. Por exemplo:

```bash
cd structure/deno-fresh
deno run -A -r https://fresh.deno.dev .
```

```bash
cd structure/vite-react
npm create vite@latest . -- --template react
```

O resultado desse comando (todos os arquivos gerados) é o que vai ficar versionado no repositório dentro dessa pasta.

> ⚠️ Não monte a estrutura de arquivos manualmente nem copie de outro lugar. Ela precisa vir da execução do comando oficial, para garantir que fique sempre alinhada com o que o framework espera.

### 3. Adicionar o `setup.go`

Na **raiz da pasta do framework** (mesmo nível dos demais arquivos gerados), crie um arquivo `setup.go` com este template, substituindo o comando pelo que aquele framework usa para instalar dependências:

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Setting up project dependencies...")

	cmd := exec.Command("npm", "install") // substituir pelo comando real
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

Você só precisa alterar a linha do `exec.Command(...)` para o comando de instalação real do seu framework. Alguns exemplos:

| Framework | Comando de instalação |
|---|---|
| Deno (Fresh, etc.) | `exec.Command("deno", "install")` |
| Node / npm | `exec.Command("npm", "install")` |
| Node / pnpm | `exec.Command("pnpm", "install")` |
| Rust | `exec.Command("cargo", "build")` |
| Python | `exec.Command("pip", "install", "-r", "requirements.txt")` |

### 4. Testar localmente antes de abrir o PR

Antes de enviar seu template, teste você mesmo:

```bash
cd structure/<seu-framework>
go run setup.go
```

Verifique se:
- As dependências são instaladas corretamente.
- O `setup.go` se autodeleta ao final.
- Não sobrou nenhum arquivo extra (node_modules, .git herdado do comando de criação, lockfiles de outro gerenciador de pacotes, etc.).

### 5. Abrir o Pull Request

Na descrição do PR, inclua:
- Nome do framework/ferramenta.
- Link para a documentação oficial de instalação usada como referência.
- Comando exato que você executou para gerar a pasta (passo 2).
- Confirmação de que testou `go run setup.go` localmente.

## Boas práticas e segurança

Como o `setup.go` executa um comando do sistema (`exec.Command`), todo PR que adicionar ou modificar um `setup.go` é revisado manualmente antes do merge. Por favor:

- Não inclua comandos que baixem ou executem scripts remotos (`curl ... | sh`, etc.) dentro do `setup.go`.
- Use apenas o comando de instalação de dependências — nada de passos adicionais, telemetria ou configuração de credenciais.
- Se o framework exigir passos extras além de instalar dependências (por exemplo, gerar um `.env`), documente isso na descrição do PR para discutirmos antes de adicioná-lo ao `setup.go`.
- Não suba `node_modules/`, binários compilados nem artefatos gerados pela instalação — a pasta em `structure/` deve conter apenas o que o comando de criação gera antes de instalar as dependências.

## Nomenclatura

- Nome da pasta: kebab-case, em minúsculas, sem espaços ou acentos (`nest-js`, `vite-vue`, `spring-boot`).
- Se houver variantes do mesmo framework (por exemplo, com templates diferentes), use um sufixo descritivo: `vite-react`, `vite-vue`, `vite-svelte`.

Obrigado por contribuir! 🙌

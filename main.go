package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var baseURL = "https://scaffolder.codegeekery.com"

type Manifest struct {
	UpdatedAt string `json:"updated_at"`
	Templates []struct {
		Name string `json:"name"`
	} `json:"templates"`
}


func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "-h", "--help":
		printUsage()
		return
	case "--available":
		if err := listAvailableTemplates(); err != nil {
			log.Fatalf("Error listing templates: %v", err)
		}
		return
	}

	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	templateName := os.Args[1]
	projectName := os.Args[2]

	if err := os.MkdirAll(projectName, 0755); err != nil {
		log.Fatalf("Error creating project directory: %v", err)
	}

	if err := fetchAndExtractTemplate(templateName, projectName); err != nil {
		log.Fatalf("Error fetching template: %v", err)
	}

	if err := runPostCreateHook(projectName); err != nil {
		log.Fatalf("Error running post-create setup: %v", err)
	}

	fmt.Printf("Project '%s' successfully created and configured!\n", filepath.Base(projectName))
}

func printUsage() {
	fmt.Println(`scaffolder - crea proyectos a partir de templates remotos

Uso:
  scaffolder <template> <project_name>   Crea un nuevo proyecto usando el template indicado
  scaffolder --available                 Lista los templates disponibles
  scaffolder -h, --help                  Muestra esta ayuda

Ejemplos:
  scaffolder go-api mi-proyecto
  scaffolder --available`)
}

func listAvailableTemplates() error {
	url := fmt.Sprintf("%s/manifest.json", baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("no se pudo conectar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta inesperada del servidor: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var m Manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return fmt.Errorf("no se pudo leer el manifiesto: %w", err)
	}

	if len(m.Templates) == 0 {
		fmt.Println("No hay templates disponibles.")
		return nil
	}

	fmt.Printf("Templates disponibles (actualizado: %s):\n\n", m.UpdatedAt)
	for _, t := range m.Templates {
		fmt.Printf("  - %s\n", t.Name)
	}
	return nil
}

func fetchAndExtractTemplate(templateName, destDir string) error {
	url := fmt.Sprintf("%s/templates/%s/latest.tar.gz", baseURL, templateName)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("no se pudo conectar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("el template '%s' no existe", templateName)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta inesperada del servidor: %s", resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}

func runPostCreateHook(targetDir string) error {
	setupPath := filepath.Join(targetDir, "setup.go")
	if _, err := os.Stat(setupPath); err == nil {
		fmt.Println("Running framework setup script...")
		cmd := exec.Command("go", "run", "setup.go")
		cmd.Dir = targetDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

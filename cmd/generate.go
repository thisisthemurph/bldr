package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/thisisthemurph/bldr/internal/builder"
	"github.com/thisisthemurph/bldr/internal/config"
	"github.com/thisisthemurph/bldr/internal/parser"
	"os"
	"strings"
)

type snippet struct {
	Data []byte
	Path string
}

var generateFromYamlCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the builder pattern for structs as detailed in your bldr.yml file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Read()
		if err != nil {
			return err
		}

		module, err := readGoModule()
		if err != nil {
			return err
		}

		snippets := make([]snippet, 0, len(cfg.Structs))

		for _, s := range cfg.Structs {
			structDetail, err := parser.ParseStruct(module, s.Path, s.Name)
			if err != nil {
				return err
			}

			code, err := builder.Generate(structDetail, s.Go.Package)
			if err != nil {
				return err
			}

			snippets = append(snippets, snippet{
				Data: []byte(code),
				Path: s.Go.Output,
			})
		}

		cmd.Println("Generating code:")
		for _, code := range snippets {
			if err := os.WriteFile(code.Path, code.Data, 0644); err != nil {
				return err
			}
			cmd.Printf("  ✅ %s generated successfully\n", code.Path)
		}

		return nil
	},
}

func readGoModule() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", errors.New("module not found")
}

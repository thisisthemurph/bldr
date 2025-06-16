package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thisisthemurph/bldr/internal/builder"
	"github.com/thisisthemurph/bldr/internal/config"
	"github.com/thisisthemurph/bldr/internal/parser"
	"os"
	"strings"
)

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

		for _, s := range cfg.Structs {
			structDetail, err := parser.ParseStruct(s.Path, s.Name)
			if err != nil {
				return err
			}

			structImport := fmt.Sprintf("%s/%s", module, structDetail.PackageDir)
			code, err := builder.Generate(structDetail, structImport, s.Go.Package)
			if err != nil {
				return err
			}

			fmt.Println(code)
			os.WriteFile(s.Go.Output, []byte(code), 0644)
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

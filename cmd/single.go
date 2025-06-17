package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thisisthemurph/bldr/internal/builder"
	"github.com/thisisthemurph/bldr/internal/config"
	"github.com/thisisthemurph/bldr/internal/parser"
	"os"
)

var structInfo config.StructInfo

var generateSingleCmd = &cobra.Command{
	Use:   "single",
	Short: "Generate the builder pattern for a given struct",
	RunE: func(cmd *cobra.Command, args []string) error {
		module, err := readGoModule()
		if err != nil {
			return err
		}

		structDetail, err := parser.ParseStruct(module, structInfo.Path, structInfo.Name)
		if err != nil {
			return err
		}

		code, err := builder.Generate(structDetail, structInfo.Go.Package)
		if err != nil {
			return err
		}

		fmt.Println(code)
		if err := os.WriteFile(structInfo.Go.Output, []byte(code), 0644); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	generateSingleCmd.Flags().StringVarP(&structInfo.Path, "file", "f", "", "Path to Go file containing struct")
	generateSingleCmd.Flags().StringVarP(&structInfo.Name, "struct", "s", "", "Name of the struct")
	generateSingleCmd.Flags().StringVarP(&structInfo.Go.Package, "package", "p", "", "Name of the output package")
	generateSingleCmd.Flags().StringVarP(&structInfo.Go.Output, "out", "o", "", "Path to output file")

	panicif(generateSingleCmd.MarkFlagRequired("file"))
	panicif(generateSingleCmd.MarkFlagRequired("struct"))
	panicif(generateSingleCmd.MarkFlagRequired("package"))
	panicif(generateSingleCmd.MarkFlagRequired("out"))
}

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}

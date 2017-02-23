package main

import "github.com/spf13/cobra"
import "fmt"
import "leaf/pkg"

var (
	// host string
	file string
)

func NewLeafCommand() *cobra.Command {
	var leafCommand = &cobra.Command{
		Use:   "api",
		Short: "-。 -",
		Long:  `-。 -`,
		Run:   readFile,
	}

	// leafCommand.Flags().StringVarP(&host, "localhost", "h", "h", "set host")
	leafCommand.Flags().StringVarP(&file, "file", "f", "api.txt", "set file path")
	return leafCommand
}

//read api doc
func readFile(cmd *cobra.Command, args []string) {
	f, err := File(file, false)
	if err != nil {
		fmt.Printf("open %s file:%s", file, err)
		return
	}

	ag, err := pkg.NewApiList(f)
	if err != nil {
		fmt.Printf("resolve template fail:%s", err)
		return
	}
	pkg.Test(ag)
}

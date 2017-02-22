package main

import ("github.com/spf13/cobra"
    "fmt"
    "leaf/pkg"
)

var Host string
var File string

func NewLeafCommand() *cobra.Command {
	var leafCommand = &cobra.Command{
		Use:   "pull",
		Short: "-。 -",
		Long:  `-。 -`,
		Run:   readFile,
	}

	leafCommand.Flags().StringVarP(&Host, "localhost", "h", "", "set host")
	leafCommand.Flags().StringVarP(&File, "api.txt", "f", "", "set file path")
	return leafCommand
}

func readFile(cmd *cobra.Command, args []string) {
    f,err:=file(File,false);err!=nil{
        fmt.Errorf("read file %s fail:%s ",File,err)
    }  
    
}

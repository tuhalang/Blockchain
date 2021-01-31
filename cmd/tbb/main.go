package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main(){
	var tbbCmd = &cobra.Command{
		Use: "tbb",
		Short: "The Blockchain CLI",
		Run: func(cmd *cobra.Command, args []string) {
			
		},
	}

	err := tbbCmd.Execute()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import(
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "term",
	Short: "simple command for words / terminology management",
	Run: func(cmd *cobra.Command, args []string){
		//somthing here
	},
}

func main() {
	rootCmd.Execute()
}

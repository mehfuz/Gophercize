package cmd

import (
	"github.com/spf13/cobra"
	
)

var RootCmd = &cobra.Command{
	Use:   "todoz",
	Short: "TODOZ is a command line tool to store activities.",
	
}

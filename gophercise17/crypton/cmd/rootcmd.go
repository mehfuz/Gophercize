package cmd

import (
	"github.com/spf13/cobra"
)

var Rootcmd = &cobra.Command{
	Use:   "crypton",
	Short: "crypton makes your system secure by providing API to enccypt and decrypt ",
}

var nkey string

func init() {
	Rootcmd.PersistentFlags().StringVarP(&nkey, "key", "k", "", "key which is used to encode data")
}

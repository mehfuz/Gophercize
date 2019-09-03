package cmd

import (
	"fmt"
	"strings"

	"github.com/gophercises/gophercise17/encrypt"
	"github.com/spf13/cobra"
)

var filename = ".Vault"
var Setcmd = &cobra.Command{
	Use:   "set",
	Short: "set makes your system secure by providing API to enccypt and decrypt ",
	Run: func(cmd *cobra.Command, args []string) {
		vault := encrypt.NewVault(nkey, filename)
		//code to include blank spaces
		key := args[0]
		value := strings.Join(args[1:], " ")
		if key == "" || value == "" {
			fmt.Println("Provide valid key or value")
			return
		}
		err := vault.Set(key, value)
		if err != nil {
			fmt.Println(err)
			//panic(err)
			return
		}
		fmt.Println("Data successfully encoded.")
	},
}

func init() {
	Rootcmd.AddCommand(Setcmd)
}

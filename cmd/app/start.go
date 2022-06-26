package cmder

import (
	"fmt"

	"github.com/spf13/cobra"
)

var start = &cobra.Command{
	Use:   "start",
	Short: "start ast",
	Long:  `start ast`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("helloast start~")
	},
}

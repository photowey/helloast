package cmder

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/photowey/helloast/internal/app"
)

var (
	conf string

	root = &cobra.Command{
		Use:   "helloast",
		Short: "helloast project",
		Long:  "A study project of ast",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing
		},
	}
)

func init() {
	cobra.OnInitialize(appRun)
	root.PersistentFlags().StringVarP(&conf, "conf", "f", "", "config file")
	root.AddCommand(start)
}

func Run() {
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func appRun() {
	if err := app.Run(); err != nil {
		cobra.CheckErr(err)
	}
}

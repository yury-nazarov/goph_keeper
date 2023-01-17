package handler

import "C"
import (
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "secret service",
	Long:  ``,
}

func init() {
	App.Cmd.AddCommand(secretCmd)
}

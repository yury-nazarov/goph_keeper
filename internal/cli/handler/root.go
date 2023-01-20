package handler

import "C"

// Переменные для мапинга из флагов

//var Cmd = &cobra.Command{
//	Use:   "gkc",
//	Short: "Goph Keeper cli",
//	Long:  `Goph Keeper command line interface`,
//}


var App = New()
func Executor() error {

	return App.Cmd.Execute()
}

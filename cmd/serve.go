package cmd

import (
	"github.com/AlexKLWS/lws-blog-server/config"
	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Runs the server",
		Run: func(cmd *cobra.Command, args []string) {
			runServer()
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

}

func runServer() {
	e := echo.New()
	e.Static("/", "output")
	e.Logger.Fatal(e.Start(viper.GetString(config.Port)))
}

package cmd

import (
	"fmt"

	"github.com/AlexKLWS/lws-blog-server/config"
	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/directories"
	"github.com/AlexKLWS/youtube-audio-stream/handlers"
	"github.com/AlexKLWS/youtube-audio-stream/router"
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

// SetupServeCommand binds command line flags and arguments to viper config instance
func SetupServeCommand() {
	serveCmd.Flags().StringP(consts.SourceDir, "s", viper.GetString(consts.SourceDir), "specify transmuxer source directory name")
	serveCmd.Flags().StringP(consts.OutputDir, "o", viper.GetString(consts.OutputDir), "specify transmuxer output directory name")
	serveCmd.Flags().StringP(consts.Port, "p", viper.GetString(consts.Port), "specify server port")
	serveCmd.Flags().String("socks-proxy", "", "The Socks 5 proxy, e.g. 10.10.10.10:7878")

	viper.BindPFlags(serveCmd.Flags())
}

func runServer() {
	directories.PrepareDirectories()

	httpTransport := client.GetHTTPTransport()
	client.New(httpTransport)

	if viper.GetBool(consts.Debug) {
		fmt.Print("Running in debug mode\n")
	}

	r := router.New()
	handlers.RegisterHandlers(r)
	r.Server.Logger.Fatal(r.Server.Start(viper.GetString(config.Port)))
}

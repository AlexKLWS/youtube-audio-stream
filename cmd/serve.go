package cmd

import (
	"github.com/AlexKLWS/lws-blog-server/config"
	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/directories"
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
	serveCmd.Flags().StringP("source", "ts", viper.GetString(consts.SourceDir), "specify transmuxer source directory name")
	serveCmd.Flags().StringP("output", "to", viper.GetString(consts.OutputDir), "specify transmuxer output directory name")
	serveCmd.Flags().StringP("port", "p", viper.GetString(consts.Port), "specify server port")
	serveCmd.Flags().StringP("socks-proxy", "sp", "", "The Socks 5 proxy, e.g. 10.10.10.10:7878")

	viper.BindPFlags(serveCmd.Flags())
}

func runServer() {
	directories.PrepareDirectories()

	httpTransport := client.GetHTTPTransport()
	client.New(httpTransport)

	r := router.New()

	r.Server.Logger.Fatal(r.Server.Start(viper.GetString(config.Port)))
}

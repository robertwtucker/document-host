//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package cmd

import (
	"context"

	"github.com/robertwtucker/document-host/internal/api"
	"github.com/robertwtucker/document-host/internal/config"
	zlog "github.com/robertwtucker/document-host/pkg/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the API server",
	Long: `Starts the HTTP(S) server on the configured port and exposes the API endpoints
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("server", RootApp.Config.Server).Debug("starting API server")
		logger := zlog.NewDebug().With(context.Background(), "version", config.AppVersion().String())
		app, err := api.NewApp(&RootApp.Config, logger)
		if err != nil {
			return
		}
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

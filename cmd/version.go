//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package cmd

import (
	"fmt"

	"github.com/robertwtucker/document-host/internal/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the " + config.AppName + " version",
	Long:  "Prints the " + config.AppName + " version",
	Run: func(cmd *cobra.Command, args []string) {
		version := config.AppVersion()
		fmt.Printf("%s %s\n", config.AppName, version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

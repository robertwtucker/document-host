//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/robertwtucker/document-host/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootApp struct {
	Config *config.Configuration
}

// RootApp represents the root application object
var RootApp = &rootApp{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   config.AppName,
	Short: "provides temporary hosting of demo documents",
	Long: `The Document Host (Docuhost) service provides a REST API endpoint to upload
demo-generated documents for temporary storage. Documents can be retrieved via the
short link returned in the upload response.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		RootApp.Config = &config.Configuration{}
		if err := viper.UnmarshalExact(&RootApp.Config); err != nil {
			return errors.Wrapf(err, "failed to unmarshal config")
		}
		if err := initLog(RootApp.Config); err != nil {
			return errors.Wrapf(err, "failed to initialize logging")
		}
		logrus.WithField("version", config.AppVersion().String()).Debug("initialized")
		return nil
	},
}

// rootCmdArgs holds the flags configured in the root Cmd
var rootCmdArgs struct {
	ConfigFile string
	LogFormat  string
	LogDebug   bool
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Process the PersistentFlags
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.ConfigFile, "config", "c",
		"", "specify the config file (default is ./config/"+config.AppName+".yaml)")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.LogFormat, "log-format", "f",
		"text", "set the logging format [text|json]")
	rootCmd.PersistentFlags().BoolVarP(&rootCmdArgs.LogDebug, "verbose", "v",
		false, "set verbose logging")

	// Hide the completions options
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// initConfig reads in config file and ENV variables, if set.
func initConfig() {
	if rootCmdArgs.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(rootCmdArgs.ConfigFile)
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
		viper.SetConfigName(config.AppName)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Have Viper check the environment for matching keys
	viper.AutomaticEnv()

	// WORKAROUND: Viper doesn't seem to be overriding the config file with values
	// from the environment. See: https://github.com/spf13/viper/issues/584
	_ = viper.BindEnv(config.AppURLKVey, config.AppURLEnv)
	_ = viper.BindEnv(config.DBPrefixKey, config.DBPrefixEnv)
	_ = viper.BindEnv(config.DBUserKey, config.DBUserEnv)
	_ = viper.BindEnv(config.DBPasswordKey, config.DBPasswordEnv)
	_ = viper.BindEnv(config.DBHostKey, config.DBHostEnv)
	_ = viper.BindEnv(config.DBPortKey, config.DBPortEnv)
	_ = viper.BindEnv(config.DBNameKey, config.DBNameEnv)
	_ = viper.BindEnv(config.DBTimeoutKey, config.DBTimeoutEnv)
	_ = viper.BindEnv(config.ServerPortKey, config.ServerPortEnv)
	_ = viper.BindEnv(config.ServerTimeoutKey, config.ServerTimeoutEnv)
	_ = viper.BindEnv(config.LogDebugKey, config.LogDebugEnv)
	_ = viper.BindEnv(config.LogFormatKey, config.LogFormatEnv)
	_ = viper.BindEnv(config.ShortLinkAPIKey, config.ShortLinkAPIEnv)
	_ = viper.BindEnv(config.ShortLinkDomainKey, config.ShortLinkDomainEnv)

	// Command-line pflags replace environment
	if rootCmdArgs.LogFormat != "" {
		viper.Set(config.LogFormatKey, rootCmdArgs.LogFormat)
	}
	if rootCmdArgs.LogDebug {
		viper.Set(config.LogDebugKey, rootCmdArgs.LogDebug)
	}
}

func initLog(cfg *config.Configuration) error {
	if "json" == strings.ToLower(cfg.Log.Format) {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	if cfg.Log.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	log.SetOutput(logrus.New().Writer())
	log.SetFlags(0)

	return nil
}

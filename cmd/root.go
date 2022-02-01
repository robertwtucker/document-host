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

var cfgFile string

type rootApp struct {
	Config config.Configuration
}

var RootApp = &rootApp{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docuhost",
	Short: "provides temporary hosting of demo documents",
	Long: `The Document Host (Docuhost) service provides a REST API endpoint to upload
demo-generated documents for temporary storage. Documents can be retrieved via the
short link returned in the upload response.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//r := &rootApp{}
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is ./config/"+config.AppName+".yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
		viper.SetConfigName(config.AppName)
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// WORKAROUND: Viper doesn't seem to be overriding the config file with values
	// from the environment. See: https://github.com/spf13/viper/issues/584
	_ = viper.BindEnv("app.url", "APP_URL")
	_ = viper.BindEnv("db.prefix", "DB_PREFIX")
	_ = viper.BindEnv("db.user", "DB_USER")
	_ = viper.BindEnv("db.password", "DB_PASSWORD")
	_ = viper.BindEnv("db.host", "DB_HOST")
	_ = viper.BindEnv("db.port", "DB_PORT")
	_ = viper.BindEnv("db.name", "DB_NAME")
	_ = viper.BindEnv("db.timeout", "DB_TIMEOUT")
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.timeout", "SERVER_TIMEOUT")
	_ = viper.BindEnv("log.debug", "LOG_DEBUG")
	_ = viper.BindEnv("log.format", "LOG_FORMAT")
	_ = viper.BindEnv("shortlink.apikey", "SHORTLINK_APIKEY")
	_ = viper.BindEnv("shortlink.domain", "SHORTLINK_DOMAIN")
}

func initLog(cfg config.Configuration) error {
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

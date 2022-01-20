//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package config

import (
	"strings"

	"github.com/robertwtucker/document-host/pkg/log"
	"github.com/spf13/viper"
)

// Init sets up the Viper configuration
func Init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if env := viper.GetString("env"); strings.ToUpper(env) != "PROD" {
		viper.Set("log.debug", true)
	}
}

// Load attempts to read the app configuration file
func Load(logger log.Logger) error {
	logger.Debug("attempting to load config file")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No config file, will use env settings
		} else {
			logger.Errorf("error loading config file: %v \n", err)
			return err
		}
	}
	logger.Infof("config file '%s' used", viper.ConfigFileUsed())

	return nil
}

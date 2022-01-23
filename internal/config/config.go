//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package config

import (
	"encoding/json"
	"fmt"
	"github.com/robertwtucker/document-host/pkg/log"
	"github.com/spf13/viper"
)

// Configuration represents the application configuration settings
type Configuration struct {
	App struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"app"`
	DB struct {
		Prefix   string `mapstructure:"prefix"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		Timeout  int64  `mapstructure:"timeout"`
	} `mapstructure:"db"`
	Server struct {
		Port    string `mapstructure:"port"`
		Timeout int64  `mapstructure:"timeout"`
	} `mapstructure:"server"`
	Log struct {
		Debug bool `mapstructure:"debug"`
	} `mapstructure:"log"`
	ShortLink struct {
		APIKey string `mapstructure:"apiKey"`
		Domain string `mapstructure:"domain"`
	} `mapstructure:"shortlink"`
}

// PrettyPrint outputs a formatted listing of the configuration settings
func (c Configuration) PrettyPrint() {
	p := fmt.Println
	p("Configuration Settings:")
	p("App:")
	p("  URL:     ", c.App.URL)
	p("DB:")
	p("  Prefix:  ", c.DB.Prefix)
	p("  User:    ", c.DB.User)
	p("  Password:", c.DB.Password)
	p("  Host:    ", c.DB.Host)
	p("  Port:    ", c.DB.Port)
	p("  Name:    ", c.DB.Name)
	p("  Timeout: ", c.DB.Timeout)
	p("Server:")
	p("  Port:    ", c.Server.Port)
	p("  Timeout: ", c.Server.Timeout)
	p("Log:")
	p("  Debug:   ", c.Log.Debug)
	p("ShortLink: ")
	p("  APIKey:  ", c.ShortLink.APIKey)
	p("  Domain:  ", c.ShortLink.Domain)
}

// String displays the configuration settings
func (c Configuration) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(out)
}

// Init sets up the Viper configuration
func Init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// WORKAROUND: Viper doesn't seem to be overriding the config file with values
	// from the environment. See: https://github.com/spf13/viper/issues/584
	viper.BindEnv("app.url", "APP_URL")
	viper.BindEnv("db.prefix", "DB_PREFIX")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.name", "DB_NAME")
	viper.BindEnv("db.timeout", "DB_TIMEOUT")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.timeout", "SERVER_TIMEOUT")
	viper.BindEnv("log.debug", "LOG_DEBUG")
	viper.BindEnv("shortlink.apikey", "SHORTLINK_APIKEY")
	viper.BindEnv("shortlink.domain", "SHORTLINK_DOMAIN")
}

// Load attempts to read the app configuration file
func Load(logger log.Logger) (*Configuration, error) {
	logger.Debug("attempting to load config file")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No config file, will use env settings
		} else {
			logger.Errorf("error loading config file: %v \n", err)
			return nil, err
		}
	}
	logger.Infof("config file '%s' used", viper.ConfigFileUsed())

	configuration := new(Configuration)
	err = viper.Unmarshal(&configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

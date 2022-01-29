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

const AppName = "docuhost"

// Configuration represents the application configuration settings
type Configuration struct {
	App struct {
		URL     string `mapstructure:"url"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`
	DB struct {
		Prefix   string `mapstructure:"prefix"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int64  `mapstructure:"port"`
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

type VersionInfo struct {
	Version  string `mapstructure:"version"`
	Revision string `mapstructure:"revision"`
}

func AppVersion() VersionInfo { return VersionInfo{Version: appVersion, Revision: revision} }

var (
	appVersion = "development"
	revision   = "unknown"
)

// PrettyPrint outputs a formatted listing of the configuration settings
func (c Configuration) PrettyPrint() {
	p := fmt.Println
	_, _ = p("Configuration Settings:")
	_, _ = p("App:")
	_, _ = p("  URL:     ", c.App.URL)
	_, _ = p("DB:")
	_, _ = p("  Prefix:  ", c.DB.Prefix)
	_, _ = p("  User:    ", c.DB.User)
	_, _ = p("  Password:", c.DB.Password)
	_, _ = p("  Host:    ", c.DB.Host)
	_, _ = p("  Port:    ", c.DB.Port)
	_, _ = p("  Name:    ", c.DB.Name)
	_, _ = p("  Timeout: ", c.DB.Timeout)
	_, _ = p("Server:")
	_, _ = p("  Port:    ", c.Server.Port)
	_, _ = p("  Timeout: ", c.Server.Timeout)
	_, _ = p("Log:")
	_, _ = p("  Debug:   ", c.Log.Debug)
	_, _ = p("ShortLink: ")
	_, _ = p("  APIKey:  ", c.ShortLink.APIKey)
	_, _ = p("  Domain:  ", c.ShortLink.Domain)
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
	_ = viper.BindEnv("shortlink.apikey", "SHORTLINK_APIKEY")
	_ = viper.BindEnv("shortlink.domain", "SHORTLINK_DOMAIN")
}

// Load attempts to read the app configuration file
func Load(logger log.Logger) (*Configuration, error) {
	// OK if No config file, will use env settings
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("config file '%s' used", viper.ConfigFileUsed())
	}

	configuration := new(Configuration)
	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Error("failed to get configuration:", err)
		return nil, err
	}

	return configuration, nil
}

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
)

// AppName represents the name of the application
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
		Debug  bool   `mapstructure:"debug"`
		Format string `mapstructure:"format"`
	} `mapstructure:"log"`
	ShortLink struct {
		APIKey string `mapstructure:"apiKey"`
		Domain string `mapstructure:"domain"`
	} `mapstructure:"shortlink"`
}

// VersionInfo represents the application's latest version tag and Git revision
type VersionInfo struct {
	Version  string `mapstructure:"version"`
	Revision string `mapstructure:"revision"`
}

// AppVersion returns the application's latest version and Git revision
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
	_, _ = p("  Format:   ", c.Log.Format)
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

// String returns a formatted form of the version and revision
func (v VersionInfo) String() string {
	return fmt.Sprintf("v%s-%s", v.Version, v.Revision)
}

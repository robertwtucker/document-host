//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package config

import (
	"errors"
	"io/ioutil"

	"github.com/robertwtucker/document-host/pkg/log"
	"gopkg.in/yaml.v2"
)

// ErrFileLoadFailed occurs when the specified config file cannot be read
var ErrFileLoadFailed = errors.New("error loading configuration file")

// ErrFileDecodeFailed occurs when the config file is not formatted correctly
var ErrFileDecodeFailed = errors.New("error decoding configuration file")

// Configuration holds the app configuration settings
type Configuration struct {
	Server *Server   `yaml:"server,omitempty"`
	DB     *Database `yaml:"database,omitempty"`
}

// Database holds the database configuration settings
type Database struct {
	URI string `yaml:"uri,omitempty" env:"DB_URI"`
}

// Server holds the HTTP server configuration settings
type Server struct {
	Addr                   string `yaml:"addr,omitempty" env:"SERVER_ADDR"`
	ReadTimeoutSeconds     int    `yaml:"read_timeout_seconds,omitempty"`
	ShutdownTimeoutSeconds int    `yaml:"shutdown_timeout_seconds,omitempty"`
	WriteTimeoutSeconds    int    `yaml:"write_timeout_seconds,omitempty"`
}

// Load creates a Configuration given a properly formatted file
func Load(file string, logger log.Logger) (*Configuration, error) {
	logger.Infof("loading config file: %s", file)

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Errorf("unable to load config file: %s", err)
		return nil, ErrFileLoadFailed
	}

	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		logger.Errorf("unable to unmarshal config file: %s", err)
		return nil, ErrFileDecodeFailed
	}

	// TODO: Look for ENV variable overrides

	return cfg, nil
}

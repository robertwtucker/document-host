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

// default values for configuration settings
const (
	defaultPort = "8080"
)

// ErrFileLoadFailed occurs when the specified config file cannot be read
var ErrFileLoadFailed = errors.New("error loading configuration file")

// ErrFileDecodeFailed occurs when the config file is not formatted correctly
var ErrFileDecodeFailed = errors.New("error decoding configuration file")

// Config holds the app configuration settings
type Config struct {
	Port string `yaml:"port,omitempty" env:"PORT"`
}

// Load creates a Config given a properly formatted file
func Load(file string, logger log.Logger) (*Config, error) {
	logger.Infof("loading config file: %s", file)

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Errorf("unable to load config file: %s", err)
		return nil, ErrFileLoadFailed
	}

	c := Config{
		Port: defaultPort,
	}
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		logger.Errorf("unable to unmarshal config file: %s", err)
		return nil, ErrFileDecodeFailed
	}

	// TODO: Look for ENV variable overrides

	return &c, nil
}

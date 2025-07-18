package config

import (
	"fmt"
)

type Controller struct {
	Host            string `config:"host" json:"host" toml:"host" yaml:"host"`
	Port            int    `config:"port" json:"port" toml:"host" yaml:"port"`
	ReadTimeout     int    `config:"read_timeout" json:"read_timeout" toml:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    int    `config:"write_timeout" json:"write_timeout" toml:"write_timeout" yaml:"write_timeout"`
	IdleTimeout     int    `config:"idle_timeout" json:"idle_timeout" toml:"idle_timeout" yaml:"idle_timeout"`
	ShutdownTimeout int    `config:"shutdown_timeout" json:"shutdown_timeout" toml:"shutdown_timeout" yaml:"shutdown_timeout"`
	CertFile        string `config:"cert_file" json:"cert_file" toml:"cert_file" yaml:"cert_file"`
	KeyFile         string `config:"key_file" json:"key_file" toml:"key_file" yaml:"key_file"`
	UseTLS          bool   `config:"use_tls" json:"use_tls" toml:"use_tls" yaml:"use_tls"`
	Enabled         bool   `config:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
}

func (c *Controller) BindAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Controller) copy() *Controller {
	if c == nil {
		return nil
	}

	return &Controller{
		Host:            c.Host,
		Port:            c.Port,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		IdleTimeout:     c.IdleTimeout,
		ShutdownTimeout: c.ShutdownTimeout,
		CertFile:        c.CertFile,
		KeyFile:         c.KeyFile,
		UseTLS:          c.UseTLS,
		Enabled:         c.Enabled,
	}
}

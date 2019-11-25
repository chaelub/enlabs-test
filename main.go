package main

import (
	"flag"
	"github.com/BurntSushi/toml"
)

const defaultConfigPath = "config/dev.toml"

var configPath *string

type Config struct {
	ServiceConfig *struct {
		HttpInterface       string `toml:"http"`
		CancellationTimeout int64  `toml:"cancellationTimeout"`
	} `toml:"service"`
}

func init() {
	configPath = flag.String("config", defaultConfigPath, "Full path to config file")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	/*
		read config
		connect to DB
		create user, transaction repo (pass DB connection to encapsulate it)
		create a user cache (pass user repo)
		start finance worker: any balance updates only via finance worker. receive user and transaction model
		create a user-transaction service (handler) to handle transaction requests POST /user/:userid/transaction
		run cancellation worker
		run a web server in separate routine?
		handle system signals for graceful shutdown
	*/

	conf := new(Config)
	_, err := toml.DecodeFile(*configPath, conf)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"enlabs-test/app/cache/user"
	"enlabs-test/app/repo/transaction"
	user_repo "enlabs-test/app/repo/user"
	"enlabs-test/server/router"
	"enlabs-test/store/pg"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/valyala/fasthttp"
)

const defaultConfigPath = "config/dev.toml"

var configPath *string

type Config struct {
	ServiceConfig *struct {
		HttpInterface       string `toml:"http"`
		CancellationTimeout int64  `toml:"cancellationTimeout"`
	} `toml:"service"`
	PGConfig *pg.PGConfig `toml:"postgres"`
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

	pgStore, err := pg.InitPGStore(*conf.PGConfig)
	if err != nil {
		return err
	}

	ur, err := user_repo.NewUserRepo(pgStore)
	if err != nil {
		return err
	}

	userCache, err := user.Cache(ur)
	if err != nil {
		return err
	}

	transactionRepo, err := transaction.NewTransactionRepo(pgStore)
	if err != nil {
		return nil
	}

	r := router.Router(userCache, transactionRepo)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}

	if err := server.ListenAndServe(":" + conf.ServiceConfig.HttpInterface); err != nil {
		return err
	}

	return nil
}

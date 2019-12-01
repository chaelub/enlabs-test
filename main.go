package main

import (
	"enlabs-test/cache/user"
	finance_manager "enlabs-test/finance-manager"
	"enlabs-test/repo/transaction"
	user_repo "enlabs-test/repo/user"
	user_transaction "enlabs-test/repo/user-transaction"
	"enlabs-test/server/router"
	"enlabs-test/store/pg"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

const defaultConfigPath = "config/dev.toml"

var configPath *string

type Config struct {
	ServiceConfig *struct {
		HttpInterface string `toml:"http"`
	} `toml:"service"`
	FinanceManagerConf *finance_manager.Config `toml:"finance-manager"`
	PGConfig           *pg.PGConfig            `toml:"postgres"`
}

func init() {
	configPath = flag.String("config", defaultConfigPath, "Full path to config file")
}

func main() {
	/*
		read config
		connect to DB
		create user, transaction repo (pass DB connection to encapsulate it)
		create a user cache (pass user repo)
		start finance worker: any balance updates only via finance worker.
		create a user-transaction service (handler) to handle transaction requests POST /user/:userid/transaction
		run cancellation worker
		run a web server
		handle system signals for graceful shutdown ?
	*/

	log := log.New(os.Stdout, "Finance-service : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	conf := new(Config)
	_, err := toml.DecodeFile(*configPath, conf)
	if err != nil {
		panic(err)
	}

	pgStore, err := pg.InitPGStore(*conf.PGConfig)
	if err != nil {
		panic(err)
	}

	ur, err := user_repo.NewUserRepo(pgStore)
	if err != nil {
		panic(err)
	}

	userCache, err := user.Cache(ur)
	if err != nil {
		panic(err)
	}

	transactionRepo, err := transaction.NewTransactionRepo(pgStore)
	if err != nil {
		panic(err)
	}

	userTransactionRepo := user_transaction.NewUserTransactionRepo(pgStore, ur, transactionRepo)

	financeManager := finance_manager.NewFinanceManager(log, conf.FinanceManagerConf, userCache, userTransactionRepo)

	go financeManager.ScheduleCancellation()

	r := router.Router(log, userCache, financeManager)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}

	if err := server.ListenAndServe(":" + conf.ServiceConfig.HttpInterface); err != nil {
		panic(err)
	}
}

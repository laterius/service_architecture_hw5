package main

import (
	"fmt"
	"github.com/ezn-go/mixture"
	"github.com/laterius/service_architecture_hw3/app/internal/transport/client/dbrepo"
	dblogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jinzhu/configor"

	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	_ "github.com/laterius/service_architecture_hw3/app/migrations"
)

func main() {
	var cfg domain.Config
	err := configor.New(&configor.Config{Silent: true}).Load(&cfg, "config/config.yaml", "./config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbrepo.Dsn(cfg.Db),
	}), &gorm.Config{
		Logger: dblogger.Default.LogMode(dblogger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = mixture.Apply(db, cfg.App.Env)
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations applied")
}

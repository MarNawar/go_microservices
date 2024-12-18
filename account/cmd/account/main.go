package main

import (
	"log"
	"time"

	"github.com/MarNawar/microservices/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	Database string `envconfig:"DATABASR_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil{
		log.Fatal(err);
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int)(err error){
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)

		if err != nil{
			log.Panicln(err)
		}
		return
	})

	defer r.Close()
	log.Panicln("Listening on port 8080...",)
	s:= account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}

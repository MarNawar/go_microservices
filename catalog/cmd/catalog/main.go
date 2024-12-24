package main

import (
	"log"
	"time"

	"github.com/MarNawar/microservices/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil{
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int)(err error){
		r, err = catalog.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil{
			log.Panicln(err)
		}
		return
	})

	defer r.Close()

	log.Panicln("Listening on port 8080 ...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
package app

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/AlexNov03/UrlShortener/internal/adapters"
	"github.com/AlexNov03/UrlShortener/internal/bootstrap"
	"github.com/AlexNov03/UrlShortener/internal/delivery"
	localrepo "github.com/AlexNov03/UrlShortener/internal/repository/local"
	"github.com/AlexNov03/UrlShortener/internal/repository/pg"
	"github.com/AlexNov03/UrlShortener/internal/server"
	"github.com/AlexNov03/UrlShortener/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type ApiEntryPoint struct {
	cfg    *bootstrap.Config
	server *server.Server
}

func NewApiEntryPoint() *ApiEntryPoint {
	return &ApiEntryPoint{}
}

func (ae *ApiEntryPoint) Init() error {
	config, err := bootstrap.ReadConfig()
	if err != nil {
		return err
	}

	ae.cfg = config

	validator := validator.New(validator.WithRequiredStructEnabled())

	inMemory := flag.Bool("in-memory", true, "defines, whether app use in-memory or postgres db")

	flag.Parse()

	var repo usecase.UrlRepository
	if *inMemory == true {
		repo = localrepo.NewUrlRepository()
		log.Printf("app is using in-memory db")
	} else {
		db, err := adapters.GetDB(ae.cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		repo = pg.NewUrlRepository(db)
		log.Printf("app is using postgres db")
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	uc := usecase.NewUrlUsecase(repo, rnd, ae.cfg)
	deliv := delivery.NewUrlDelivery(uc, validator)

	ae.server = server.NewServer(ae.cfg, deliv)
	ae.server.Init()

	return nil
}

func (ae *ApiEntryPoint) Run() error {
	return ae.server.Run()

}

func (ae *ApiEntryPoint) Stop() error {
	return ae.server.Stop()
}

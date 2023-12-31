package main

import (
	"github.com/hamza72x/islamqa-scrapper/content"
	"github.com/hamza72x/islamqa-scrapper/log"
	"github.com/hamza72x/islamqa-scrapper/scrapper"
	"github.com/hamza72x/islamqa-scrapper/sitemap"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	log.Initialize()

	// open database
	db, err := gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: " + err.Error())
	}

	// migrate
	if err := db.AutoMigrate(
		&sitemap.URL{},
		&content.Content{},
		&content.ContentV2{},
	); err != nil {
		log.Fatal("failed to migrate database: " + err.Error())
	}

	// create scrapper
	s := scrapper.New(db)

	// sync sitemaps
	errs := s.SyncSitemaps(siteMaps)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Err(err)
		}
	}

	// // sync conents
	// if err := s.SyncContents(); err != nil {
	// 	log.Err(err)
	// }

	//  sync conents v2
	if err := s.SyncContentsV2(); err != nil {
		log.Err(err)
	}
}

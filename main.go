package main

import (
	"github.com/hamza72x/islamqa-scrapper/content"
	"github.com/hamza72x/islamqa-scrapper/log"
	scapper "github.com/hamza72x/islamqa-scrapper/scrapper"
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
	if err := db.AutoMigrate(&sitemap.URL{}, &content.Content{}); err != nil {
		log.Fatal("failed to migrate database: " + err.Error())
	}

	// create scrapper
	s := scapper.New(db)

	// sync sitemaps
	errs := s.SyncSitemaps(siteMaps)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Err(err)
		}
	}

	// sync conents
	if err := s.SyncContents(); err != nil {
		log.Err(err)
	}
}

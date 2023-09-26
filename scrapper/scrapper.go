package scrapper

import (
	"errors"
	"sync"
	"time"

	"github.com/hamza72x/islamqa-scrapper/content"
	"github.com/hamza72x/islamqa-scrapper/log"
	"github.com/hamza72x/islamqa-scrapper/sitemap"

	"gorm.io/gorm"
)

const (
	THREADS = 20
)

type Scapper struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Scapper {
	return &Scapper{
		db: db,
	}
}

func (s *Scapper) SyncContents() error {
	ch := make(chan int, THREADS)
	wg := &sync.WaitGroup{}

	urls := []*sitemap.URL{}

	if err := s.db.
		Find(&urls).
		Error; err != nil {
		return err
	}

	count := len(urls)
	completed := 0

	log.Ok("total urls to sync:", count)
	time.Sleep(1 * time.Second)

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url *sitemap.URL) {
			defer wg.Done()
			ch <- i

			if err := s.syncContent(url); err != nil {
				log.Err(err)
			}

			completed++

			if completed%100 == 0 {
				log.Ok("completed", completed, "out of", count)
			}

			<-ch
		}(i, url)
	}

	wg.Wait()

	return nil
}

func (s *Scapper) SyncSitemaps(sitemaps []string) []error {
	ch := make(chan int, THREADS)
	wg := &sync.WaitGroup{}

	errs := []error{}

	// sync all site maps
	for i, url := range sitemaps {
		wg.Add(1)

		go func(i int, url string) {
			defer wg.Done()
			ch <- i

			if err := s.syncSitemap(url); err != nil {
				errs = append(errs, err)
			}

			<-ch
		}(i, url)
	}

	wg.Wait()

	return errs
}

func (s *Scapper) syncSitemap(sitemapURL string) error {

	log.Info("syncing sitemap", sitemapURL)

	smap, err := sitemap.Get(sitemapURL)

	if err != nil {
		return err
	}

	for _, url := range smap.URLS {
		existingCounter := int64(0)

		url.SitemapUrl = sitemapURL

		if err := s.db.
			Model(&sitemap.URL{}).
			Where("loc = ?", url.Loc).
			Count(&existingCounter).
			Error; err != nil {
			return err
		}

		if existingCounter > 0 {
			continue
		}

		if err := s.db.
			Create(url).
			Error; err != nil {
			return err
		}
	}

	log.Ok("completed syncing", sitemapURL)

	return nil
}

func (s *Scapper) syncContent(url *sitemap.URL) error {
	existingContent := &content.Content{}

	if err := s.db.
		Model(&content.Content{}).
		Where("url = ?", url.Loc).
		First(existingContent).
		Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if existingContent.ID > 0 {
		return nil
	}

	newContent, err := content.New(url)
	if err != nil {
		return err
	}

	// update the content if it exists
	if existingContent.ID > 0 {

		if existingContent.Title == newContent.Title &&
			existingContent.Content == newContent.Content &&
			existingContent.Summary == newContent.Summary &&
			existingContent.Body == newContent.Body {
			return nil
		}

		existingContent.Title = newContent.Title
		existingContent.Content = newContent.Content
		existingContent.Summary = newContent.Summary
		existingContent.Body = newContent.Body

		if err := s.db.Save(existingContent).Error; err != nil {
			return err
		}

		return nil
	}

	// otherwise create a new one
	if err := s.db.Create(newContent).Error; err != nil {
		return err
	}

	return nil
}

// for ContentV2
func (s *Scapper) SyncContentsV2() error {
	ch := make(chan int, THREADS)
	wg := &sync.WaitGroup{}

	urls := []*sitemap.URL{}

	if err := s.db.
		Limit(1000).
		Order("last_mod desc").
		Find(&urls).
		Error; err != nil {
		return err
	}

	count := len(urls)
	completed := 0

	log.Ok("total urls to sync:", count)
	time.Sleep(1 * time.Second)

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url *sitemap.URL) {
			defer wg.Done()
			ch <- i

			if err := s.syncContentV2(url); err != nil {
				log.Err(err)
			}

			completed++

			if completed%100 == 0 {
				log.Ok("completed", completed, "out of", count)
			}

			<-ch
		}(i, url)
	}

	wg.Wait()

	return nil
}

func (s *Scapper) syncContentV2(url *sitemap.URL) error {
	existingContent := &content.ContentV2{}

	if err := s.db.
		Model(&content.ContentV2{}).
		Where("url = ?", url.Loc).
		First(existingContent).
		Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if existingContent.ID > 0 && existingContent.LastModified == url.LastMod {
		return nil
	}

	newContent, err := content.NewV2(url)
	if err != nil {
		return err
	}

	// update the content if it exists
	if existingContent.ID > 0 {

		if existingContent.Title == newContent.Title &&
			existingContent.Content == newContent.Content &&
			existingContent.LastModified == newContent.LastModified {
			return nil
		}

		existingContent.Title = newContent.Title
		existingContent.Content = newContent.Content
		existingContent.LastModified = newContent.LastModified

		if err := s.db.Save(existingContent).Error; err != nil {
			return err
		}

		return nil
	}

	// otherwise create a new one
	if err := s.db.Create(newContent).Error; err != nil {
		return err
	}

	return nil
}

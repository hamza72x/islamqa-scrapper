package content

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hamza72x/islamqa-scrapper/helper"
	"github.com/hamza72x/islamqa-scrapper/log"
	"github.com/hamza72x/islamqa-scrapper/sitemap"

	"github.com/PuerkitoBio/goquery"
)

type ContentV2 struct {
	ID    uint    `gorm:"primarykey;column:id"`
	Title *string `gorm:"column:title"`

	// Content is either question body or seo description
	// for fatwa, it's "question body"
	// for article, it's "seo description"
	Content *string `gorm:"column:content"`

	URL string `gorm:"column:url;uniqueIndex"`

	// LastModified is the last modified date of the content
	// populated from sitemap
	LastModified time.Time `gorm:"column:last_modified"`
}

func (ContentV2) TableName() string {
	return "contents_v2"
}

func NewV2(url *sitemap.URL) (*ContentV2, error) {
	resp, err := helper.GetURLResponse(url.Loc, "")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not ok, but: " + resp.Status)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c := &ContentV2{
		URL:          url.Loc,
		LastModified: url.LastMod,
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		log.Fatal(err)
	}

	// title
	/*
		<div class="single-layout__title has-text-centered">
				<h1 class="title is-4 is-size-5-touch">
					If he gets married with a dowry, part of which is deferred until the time of death or separation
				</h1>
		</div>
	*/
	title := doc.Find("div.single-layout__title").Find("h1").Text()
	title = strings.TrimSpace(title)
	if len(title) > 0 {
		c.Title = &title
	} else {
		log.Warn("failed to parse title for", c.URL)
	}

	// question
	/*
		<section class="single_fatwa__question text-justified">
			<h2 class="has-text-weight-bold subtitle">প্রশ্ন</h2>
			<div>প্রশ্ন: আমি আমেরিকাতে প্রবাসী। স্বর্ণের নিসাব আমেরিকান ডলারে কত আসবে?</div>
		</section>
	*/

	question := doc.Find("section.single_fatwa__question").Find("div").Text()
	question = strings.TrimSpace(question)
	if len(question) > 0 {
		c.Content = &question
	} else {
		// article, 'seo description'
		/*
			<meta name="description" content="Giving a name such as this" />
		*/
		seoDescription, exists := doc.Find("meta[name='description']").Attr("content")
		if exists {
			seoDescription = strings.TrimSpace(seoDescription)
			c.Content = &seoDescription
		} else {
			log.Warn("failed to parse question or description for", c.URL)
		}
	}

	return c, nil
}

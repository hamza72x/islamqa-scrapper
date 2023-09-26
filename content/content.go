package content

import (
	"bytes"
	"io"
	"islam-qa-scrapper/log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
)

type Content struct {
	gorm.Model

	// URL is the url of the content
	URL string `gorm:"column:url;uniqueIndex"`

	// Title is parsed question/title (fatwa/article)
	Title *string `gorm:"column:title"`

	// Content is the parsed answer/content (fatwa/article)
	Content *string `gorm:"column:content"`

	// Summary is the parsed "summarized" part
	// not all of the answer has that
	Summary *string `gorm:"column:summary"`

	// Body is the whole html body
	Body string `gorm:"column:body"`
}

func New(resp *http.Response) (*Content, error) {
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c := &Content{
		URL:  resp.Request.URL.String(),
		Body: string(htmlBytes),
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		log.Fatal(err)
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
		c.Title = &question
	} else {
		// single-layout__title
		/*
			<div class="single-layout__title has-text-centered">
				<h1 class="title is-4 is-size-5-touch" itemprop="name">
					هكذا بشر رسول الله أصحابه بقدوم رمضان
				</h1>
			</div>
		*/
		title := doc.Find("div.single-layout__title").Find("h1").Text()
		title = strings.TrimSpace(title)
		c.Title = &title
	}

	// answer
	/*
		<section class="single_fatwa__answer__body text-justified _pa--0">
			<div class="content">
				<p>আলহামদু লিল্লাহ।.</p>
			</div>
		</section>
	*/

	answer, err := doc.Find("section.single_fatwa__answer__body").Find("div.content").Html()
	if err != nil {
		log.Warn("failed to parse answer for", c.URL, "err", err)
	} else {
		answer = strings.TrimSpace(answer)
		c.Content = &answer
	}

	// summary
	/*
		<div class="single_fatwa__summary__body">
			<div>
				Giving a name such as this
			</div>
		</div>
	*/
	summary := doc.Find("div.single_fatwa__summary__body").Find("div").Text()
	summary = strings.TrimSpace(summary)
	c.Summary = &summary

	return c, nil
}

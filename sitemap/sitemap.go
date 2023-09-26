package sitemap

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"islam-qa-scrapper/helper"
	"log"
	"regexp"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	// Time interval to be used in Index.get
	interval = time.Second
)

// Index is a structure of <sitemapindex>
type Index struct {
	XMLName xml.Name `xml:"sitemapindex"`
	Sitemap []parts  `xml:"sitemap"`
}

// parts is a structure of <sitemap> in <sitemapindex>
type parts struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// Sitemap is a structure of <sitemap>
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLS    []*URL   `xml:"url"`
}

// URL is a structure of <url> in <sitemap>
type URL struct {
	gorm.Model
	SitemapUrl string  `xml:"sitemap" gorm:"column:sitemap_url;index"`
	Loc        string  `xml:"loc" gorm:"column:loc;uniqueIndex"`
	LastMod    string  `xml:"lastmod" gorm:"column:last_mod"`
	ChangeFreq string  `xml:"changefreq" gorm:"column:change_freq"`
	Priority   float32 `xml:"priority" gorm:"column:priority"`
}

func (URL) TableName() string {
	return "urls"
}

// GetTime get time from url.LastMod
func (url URL) GetTime() time.Time {
	t, err := time.Parse(time.RFC3339, url.LastMod)
	if err != nil {
		log.Println("Error parsing time", url.LastMod, url)
		return time.Now()
	}
	return t
}

// GetByWPJSON make urls from wp-json api
// WpJsonUrl: https://www.muslimmedia.info/wp-json/wp/v2/posts?per_page=50
// WpJsonUrl: https://www.muslimmedia.info/wp-json/wp/v2/posts?per_page=50&post_type=post
func GetByWPJSON(wpJSONURL string, userAgent string) Sitemap {
	var siteMap = Sitemap{}

	type WpPost struct {
		Date string `json:"date"`
		Link string `json:"link"`
	}

	var page = 1
	timeRegex := regexp.MustCompile(`\d+-\d+-\d+`)
	for {
		var posts []WpPost
		u := wpJSONURL + "&page=" + strconv.Itoa(page)

		log.Println("Getting urls from =>", u)

		bytes, err := helper.GetURLBytes(u, userAgent)

		if err != nil {
			log.Println("Error getting up json", err)
			break
		}

		if err := json.Unmarshal(bytes, &posts); err != nil {
			log.Println("Error getting up json", err)
			break
		}

		log.Println("Found posts count =>", len(posts))

		for _, p := range posts {
			t, err := time.Parse("2006-01-02", timeRegex.FindString(p.Date))

			if err != nil {
				log.Println(
					"Error parsing Time",
					"Regex found Time =>", timeRegex.FindString(p.Date),
					"Original Time =>", p.Date,
					"Sitemap input Time =>", t.Format(time.RFC3339),
				)
				continue
			}

			siteMap.URLS = append(siteMap.URLS, &URL{
				SitemapUrl: wpJSONURL,
				Loc:        p.Link,
				LastMod:    t.Format(time.RFC3339),
			})
		}
		page++
	}

	return siteMap
}

// Get sitemap data from URL
func Get(URL string) (Sitemap, error) {
	data, err := helper.GetURLBytes(URL, "")
	if err != nil {
		return Sitemap{}, err
	}

	idx, idxErr := ParseIndex(data)
	smap, smapErr := Parse(data)

	if idxErr != nil && smapErr != nil {
		return Sitemap{}, errors.New("URL is not a sitemap or sitemapindex")
	} else if idxErr != nil {
		return smap, nil
	}

	smap, err = idx.get(data)
	if err != nil {
		return Sitemap{}, err
	}

	return smap, nil
}

// Get Sitemap data from sitemapindex file
func (s *Index) get(data []byte) (Sitemap, error) {
	idx, err := ParseIndex(data)
	if err != nil {
		return Sitemap{}, err
	}

	var smap Sitemap
	for _, s := range idx.Sitemap {
		time.Sleep(interval)

		data, err := helper.GetURLBytes(s.Loc, "")
		if err != nil {
			return smap, err
		}

		err = xml.Unmarshal(data, &smap)
		if err != nil {
			return smap, err
		}
	}

	return smap, err
}

// Parse create Sitemap data from text
func Parse(data []byte) (smap Sitemap, err error) {
	err = xml.Unmarshal(data, &smap)
	return
}

// ParseIndex create Index data from text
func ParseIndex(data []byte) (idx Index, err error) {
	err = xml.Unmarshal(data, &idx)
	return
}

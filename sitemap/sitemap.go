package sitemap

import (
	"encoding/xml"
	"errors"
	"time"

	"github.com/hamza72x/islamqa-scrapper/helper"
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
	ID         uint      `gorm:"primarykey;column:id"`
	SitemapUrl string    `xml:"sitemap" gorm:"column:sitemap_url;index"`
	Loc        string    `xml:"loc" gorm:"column:loc;uniqueIndex"`
	LastMod    time.Time `xml:"lastmod" gorm:"column:last_mod"`
	ChangeFreq string    `xml:"changefreq" gorm:"-"`
	Priority   float32   `xml:"priority" gorm:"-"`
}

func (URL) TableName() string {
	return "urls"
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

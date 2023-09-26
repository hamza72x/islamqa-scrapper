package scapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetContent(t *testing.T) {
	tt := []struct {
		url        string
		hasSummary bool
	}{
		{
			url:        "https://islamqa.info/en/answers/314410/is-it-permissible-to-call-something-haleem",
			hasSummary: true,
		},
		{
			url:        "https://islamqa.info/fr/articles/19/le-ramadan-et-la-maison",
			hasSummary: false,
		},
		{
			url:        "https://islamqa.info/ar/articles/16/%D9%87%D9%83%D8%B0%D8%A7-%D8%A8%D8%B4%D8%B1-%D8%B1%D8%B3%D9%88%D9%84-%D8%A7%D9%84%D9%84%D9%87-%D8%A7%D8%B5%D8%AD%D8%A7%D8%A8%D9%87-%D8%A8%D9%82%D8%AF%D9%88%D9%85-%D8%B1%D9%85%D8%B6%D8%A7%D9%86",
			hasSummary: false,
		},
		{

			url:        "https://islamqa.info/bn/answers/47761/%E0%A6%AF%E0%A6%B0-%E0%A6%AC%E0%A6%AF%E0%A6%AC%E0%A6%B8%E0%A7%9F%E0%A6%95-%E0%A6%AA%E0%A6%A3%E0%A6%AF-%E0%A6%AF%E0%A6%95%E0%A6%A4-%E0%A6%93%E0%A7%9F%E0%A6%9C%E0%A6%AC-%E0%A6%B9%E0%A7%9F%E0%A6%9B-%E0%A6%95%E0%A6%A8%E0%A6%A4-%E0%A6%A4%E0%A6%B0-%E0%A6%95%E0%A6%9B-%E0%A6%A8%E0%A6%97%E0%A6%A6-%E0%A6%85%E0%A6%B0%E0%A6%A5-%E0%A6%A8%E0%A6%87",
			hasSummary: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.url, func(t *testing.T) {
			c, err := getContent(tc.url)
			require.NoError(t, err)

			require.NotEmpty(t, c.URL)
			require.NotEmpty(t, c.Body)
			require.NotEmpty(t, c.Title)
			require.NotEmpty(t, c.Content)

			if tc.hasSummary {
				require.NotEmpty(t, c.Summary)
			}
		})
	}
}

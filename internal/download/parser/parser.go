package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetDoc(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func Title(doc *goquery.Document) string {
	h1Elem := doc.Find("h1").First()
	h1Title, found := h1Elem.Attr("title")
	if !found {
		h1Title = h1Elem.Text()
	}
	title := strings.ReplaceAll(strings.TrimSpace(h1Title), "\n", "")
	if title == "" {
		// Bilibili: Some movie page got no h1 tag
		title, _ = doc.Find("meta[property=\"og:title\"]").Attr("content")
	}
	if title == "" {
		title = doc.Find("title").Text()
	}
	return title
}

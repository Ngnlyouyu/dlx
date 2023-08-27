package extractors

import (
	"dlx/internal/download/datafile"
	"dlx/utils"
	"net/url"
	"strings"
	"sync"
)

var lock sync.RWMutex
var extractorMap = make(map[string]datafile.Extractor)

// Register registers an Extractor.
func Register(domain string, e datafile.Extractor) {
	lock.Lock()
	extractorMap[domain] = e
	lock.Unlock()
}

// Extract is the main function to extract the data.
func Extract(u string, option datafile.Options) ([]*datafile.Data, error) {
	u = strings.TrimSpace(u)
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, err
		}
		if u.Host == "haokan.baidu.com" {
			domain = "haokan"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor := extractorMap[domain]
	if extractor == nil {
		extractor = extractorMap[""]
	}
	videos, err := extractor.Extract(u, option)
	if err != nil {
		return nil, err
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}

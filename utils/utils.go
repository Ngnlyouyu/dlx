package utils

import (
	"os"
	"regexp"

	"github.com/cheggaaa/pb/v3"
)

// MatchOneOf match one of the patterns
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// MatchAll return all matching results
func MatchAll(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	value := re.FindAllStringSubmatch(text, -1)
	return value
}

func FileSize(filePath string) (int64, bool, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return fi.Size(), true, nil
}

func ProgressBar(size int64) *pb.ProgressBar {
	tmpl := `{{counters .}} {{bar . "[" "=" ">" "-" "]"}} {{speed .}} {{percent . | green}} {{rtime .}}`
	return pb.New64(size).
		Set(pb.Bytes, true).
		SetMaxWidth(1000).
		SetTemplate(pb.ProgressBarTemplate(tmpl))
}

// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return ""
}

// Reverse Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Range generate a sequence of numbers by range
func Range(min, max int) []int {
	items := make([]int, max-min+1)
	for index := range items {
		items[index] = min + index
	}
	return items
}

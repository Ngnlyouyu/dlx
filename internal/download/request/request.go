package request

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"dlx/internal/download/config"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	cookiemonster "github.com/MercuryEngineering/CookieMonster"
)

var DefaultRequest *Req

var ErrDefaultReqIsNil = errors.New("default request is nil")

type Req struct {
	RetryTimes int
	Cookie     string
	UserAgent  string
	Refer      string
	Debug      bool
}

func NewRequest(retryTimes int, cookie, userAgent, refer string, debug bool) *Req {
	return &Req{
		RetryTimes: retryTimes,
		Cookie:     cookie,
		UserAgent:  userAgent,
		Refer:      refer,
		Debug:      debug,
	}
}

func Request(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	if DefaultRequest == nil {
		return nil, ErrDefaultReqIsNil
	}
	return DefaultRequest.Call(http.MethodGet, url, nil, headers)
}

func (r *Req) Call(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range config.DefaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["Referer"]; !ok {
		req.Header.Set("Referer", url)
	}
	if r.Cookie != "" {
		cookies, _ := cookiemonster.ParseString(r.Cookie)
		if len(cookies) > 0 {
			for _, c := range cookies {
				req.AddCookie(c)
			}
		} else {
			req.Header.Set("Cookie", r.Cookie)
		}
	}

	if r.UserAgent != "" {
		req.Header.Set("User-Agent", r.UserAgent)
	}
	if r.Refer != "" {
		req.Header.Set("Referer", r.Refer)
	}

	var (
		res    *http.Response
		reqErr error
	)
	for i := 0; ; i++ {
		res, reqErr = client.Do(req)
		if reqErr == nil && res.StatusCode < 400 {
			break
		} else if i+1 >= r.RetryTimes {
			var err error
			if reqErr != nil {
				err = fmt.Errorf("request error: %v", reqErr)
			} else {
				err = fmt.Errorf("%s request error: HTTP %d", url, res.StatusCode)
			}
			return nil, err
		}
		time.Sleep(time.Second)
	}

	if r.Debug {
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("Method: %s\n", method)
		fmt.Printf("Headers: %v\n", req.Header)
		fmt.Printf("Status Code: %d\n", res.StatusCode)
	}
	return res, nil
}

func Get(url, refer string, headers map[string]string) (string, error) {
	body, err := GetByte(url, refer, headers)
	return string(body), err
}

func GetByte(url, refer string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if refer != "" {
		headers["Referer"] = refer
	}

	res, err := Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

func Headers(url, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
	}

	res, err := Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return res.Header, nil
}

func Size(url, refer string) (int64, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return 0, err
	}
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not present")
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func ContentType(url, refer string) (string, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return "", err
	}
	s := h.Get("Content-Type")
	return strings.Split(s, ";")[0], nil
}

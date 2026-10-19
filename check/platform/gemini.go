package platform

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

// The availability flag used to appear as the literal "45631641,null,true".
// Google now embeds it in an escaped JSON map as \"45631641\":true/false.
var geminiFlagRe = regexp.MustCompile(`45631641\\?"\s*:\s*(true|false)`)

// https://github.com/clash-verge-rev/clash-verge-rev/blob/c894a15d13d5bcce518f8412cc393b56272a9afa/src-tauri/src/cmd/media_unlock_checker.rs#L241
func CheckGemini(httpClient *http.Client) (bool, error) {
	req, err := http.NewRequest("GET", "https://gemini.google.com/", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return geminiAvailable(body), nil
}

func geminiAvailable(body []byte) bool {
	if match := geminiFlagRe.FindSubmatch(body); match != nil {
		return string(match[1]) == "true"
	}
	return strings.Contains(string(body), "45631641,null,true")
}

package helper

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	// UserAgentCrawler generic crawler user agent
	UserAgentCrawler = "Crawler"

	// UserAgentChrome79Windows Chrome 79 Windows
	UserAgentChrome79Windows = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"
)

// URLContentMust return []bytes
// panics if failed
func GetURLBytesMust(urlStr string, userAgent string) []byte {

	bytes, err := GetURLBytes(urlStr, userAgent)

	if err != nil {
		panic("[GetURLBytesMust] Error getting data - " + err.Error())
	}

	return bytes
}

func GetURLBytes(urlStr string, userAgent string) ([]byte, error) {

	// Make request
	resp, err := GetURLResponse(urlStr, userAgent)

	if err != nil {
		return nil, err
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return htmlBytes, nil
}

// GetURLResponse get full response of a url
// make sure to call `defer response.Body.Close()` in your caller function
func GetURLResponse(urlStr string, userAgent string) (*http.Response, error) {
	// fmt.Printf("HTML code of %s ...\n", urlStr)
	if len(userAgent) == 0 {
		userAgent = UserAgentCrawler
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create and modify HTTP req before sending
	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return &http.Response{}, err
	}

	// set user agent
	req.Header.Set("User-Agent", userAgent)

	// Make request
	resp, err := client.Do(req)

	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}

// Exec does a os command
// and return stdout as string
// and stderr as error
func Exec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

// ExecStd does a os command
// uses  stdout, and stderr
func ExecStd(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RemoveFileIfExists remove file if exists
func RemoveFileIfExists(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}

	return nil
}

package main

import (
	"bytes"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"log"
	"os"
	"strings"
)

var dlSession string

var UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

func init() {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_105),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	tokenURL := "https://clearance.deepl.com/token"
	req, _ := http.NewRequest("GET", tokenURL, nil)

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", UA)
	req.Header.Set("Referer", "https://www.deepl.com/")
	req.Header.Set("Origin", "https://www.deepl.com")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request failed: ", err)
	}

	setCookies := resp.Header.Get("Set-Cookie")
	// 判断dl_clearance
	if !strings.Contains(setCookies, "dl_clearance") {
		log.Fatal("dl_clearance not found")
	}

	email := os.Getenv("DL_EMAIL")
	password := os.Getenv("DL_PASSWORD")

	jsonStr := []byte(fmt.Sprintf(
		`{"id":33490001,"jsonrpc":"2.0","method":"login","params":{"clearanceInfo":{"status":200,"duration":819},"referrer":"https://www.deepl.com/zh/login","email":"%s","password":"%s","version":"44","loginDomain":"default"}}`,
		email, password))
	loginURL := "https://w.deepl.com/account?request_type=jsonrpc&il=zh&method=login"
	loginReq, _ := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonStr))

	loginReq.Header.Set("Accept", "*/*")
	loginReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	loginReq.Header.Set("User-Agent", UA)
	loginReq.Header.Set("Referer", "https://www.deepl.com/")
	loginReq.Header.Set("Origin", "https://www.deepl.com")
	loginReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Fatal("Request failed: ", err)
	}

	setCookies = loginResp.Header.Get("Set-Cookie")

	// 判断dl_session
	if !strings.Contains(setCookies, "dl_session") {
		log.Fatal("dl_session not found")
	}

	cookies := strings.Split(setCookies, ";")
	for _, cookie := range cookies {
		parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(parts) == 2 && parts[0] == "dl_session" {
			dlSession = parts[1]
			break
		}
	}

}

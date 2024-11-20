package http

import (
	"context"
	"net/http"
	"net/url"

	"efaturas-xtreme/pkg/errors"
)

func Get(ctx context.Context, baseURL string, cookies map[string]string) (*http.Response, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.New("failed to parse url:", err)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, errors.New("failed to create http request:", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Referer", "https://faturas.portaldasfinancas.gov.pt/")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("Priority", "u=0, i")

	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to make request:", err)
	}

	return res, nil
}

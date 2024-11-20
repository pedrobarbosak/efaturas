package http

import (
	"context"
	"io"
	"net/http"

	"efaturas-xtreme/pkg/errors"
)

func Post(ctx context.Context, baseURL string, body io.Reader, cookies map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, body)
	if err != nil {
		return nil, errors.New("failed to create http request:", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd") // disabled to be able to parse the response without having to decode
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.acesso.gov.pt")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.acesso.gov.pt/jsp/loginRedirectForm.jsp?path=painelAdquirente.action&partID=EFPF")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")

	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to make request:", err)
	}

	return res, nil
}

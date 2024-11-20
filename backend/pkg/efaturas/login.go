package efaturas

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	httpClient "efaturas-xtreme/pkg/efaturas/http"
	"efaturas-xtreme/pkg/efaturas/parser"

	"efaturas-xtreme/pkg/errors"
)

type loginURLs struct {
	csrf         string
	login        string
	authenticate string
}

func (s *service) Login(ctx context.Context, uname string, pword string) (map[string]string, error) {
	csrf, cookies, err := s.getCSRFLogin(ctx)
	if err != nil {
		return nil, err
	}

	fields, cookies, err := s.login(ctx, cookies, csrf, uname, pword)
	if err != nil {
		return nil, err
	}

	authCookies, err := s.authenticateAndRedirect(ctx, cookies, fields)
	if err != nil {
		return nil, err
	}

	return authCookies, nil
}

func (s *service) getCSRFLogin(ctx context.Context) (string, map[string]string, error) {
	resp, err := httpClient.Get(ctx, s.urls.login.csrf, nil)
	if err != nil {
		return "", nil, errors.New("failed to make login request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("failed to make csrf request:", resp.StatusCode, resp.Status)
	}

	token, err := parser.ParseCSRF(ctx, resp.Body)
	if err != nil {
		return "", nil, err
	}

	cookies := make(map[string]string)
	for _, cookie := range resp.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	return token, cookies, nil
}

func (s *service) login(ctx context.Context, cookies map[string]string, csrf string, uname string, pword string) (map[string]string, map[string]string, error) {
	body := fmt.Sprintf("path=painelAdquirente.action&partID=EFPF&authVersion=1&_csrf=%s&selectedAuthMethod=N&username=%s&password=%s", csrf, uname, pword)
	resp, err := httpClient.Post(ctx, s.urls.login.login, strings.NewReader(body), cookies)
	if err != nil {
		return nil, nil, errors.New("failed to make login request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("failed to make login request:", resp.StatusCode, resp.Status)
	}

	fields, err := parser.ParseLoginForm(ctx, resp.Body)
	if err != nil {
		return nil, nil, err
	}

	for _, cookie := range resp.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	return fields, cookies, nil
}

func (s *service) authenticateAndRedirect(ctx context.Context, cookies map[string]string, fields map[string]string) (map[string]string, error) {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}

	resp, err := httpClient.Post(ctx, s.urls.login.authenticate, bytes.NewBufferString(form.Encode()), cookies)
	if err != nil {
		return nil, errors.New("failed to make authenticate request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to make authenticate request:", resp.StatusCode, resp.Status)
	}

	authCookies := make(map[string]string)
	for _, cookie := range resp.Cookies() {
		authCookies[cookie.Name] = cookie.Value
	}

	return authCookies, nil
}

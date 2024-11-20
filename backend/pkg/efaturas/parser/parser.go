package parser

import (
	"context"
	"io"
	"strings"

	"efaturas-xtreme/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)

func ParseCSRF(ctx context.Context, data io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return "", errors.New("failed to read data:", err)
	}

	token := doc.Find(`input[name="_csrf"]`).AttrOr("value", "")
	if token == "" {
		return "", errors.New("failed to get token: empty")
	}

	return token, nil
}

func ParseLoginForm(ctx context.Context, data io.Reader) (map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil, errors.New("failed to read data:", err)
	}

	forms := doc.Find(`form#forwardParticipantForm`)
	if forms.Length() != 1 {
		return nil, errors.New("something as changed: forms length", forms.Length())
	}

	fields := make(map[string]string)
	forms.Find("input").Each(func(_ int, s *goquery.Selection) {
		name := s.AttrOr("name", "")
		if name == "" {
			return
		}

		value := s.AttrOr("value", "")
		if value == "" {
			return
		}

		fields[name] = value
	})

	if len(fields) == 0 {
		return nil, errors.New("something went wrong: no fields found ...")
	}

	return fields, nil
}

func ParseCheckInvoice(ctx context.Context, data io.Reader) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return false, errors.New("failed to read data:", err)
	}

	errorDiv := doc.Find("div.alert-error")
	if errorDiv.Length() != 1 {
		return false, errors.New("found 0 or many error divs:", errorDiv.Length())
	}

	successDiv := doc.Find("div.alert-success")
	if successDiv.Length() != 1 {
		return false, errors.New("found 0 or many success divs:", successDiv.Length())
	}

	errorMsg := strings.TrimSpace(errorDiv.Text())
	if len(errorMsg) != 0 {
		return false, nil
	}

	successMsg := strings.TrimSpace(successDiv.Text())
	if len(successMsg) == 0 {
		return false, nil
	}

	return true, nil
}

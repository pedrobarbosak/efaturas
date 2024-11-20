package efaturas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"efaturas-xtreme/internal/service/domain"
	httpClient "efaturas-xtreme/pkg/efaturas/http"
	"efaturas-xtreme/pkg/efaturas/parser"

	"efaturas-xtreme/pkg/errors"
)

func (s *service) CheckInvoice(ctx context.Context, cookies map[string]string, invoice *domain.Invoice, category domain.Category) (bool, error) {
	body := url.Values{}
	body.Set("idDocumento", fmt.Sprintf("%d", invoice.ID))
	body.Set("hashDocumento", invoice.Document.Hash)
	body.Set("ambitoAquisicaoPend", category.String())

	resp, err := httpClient.Post(ctx, s.urls.invoices.invoice, strings.NewReader(body.Encode()), cookies)
	if err != nil {
		return false, errors.New("failed to make checkInvoice request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to make checkInvoice request:", resp.StatusCode, resp.Status)
	}

	success, err := parser.ParseCheckInvoice(ctx, resp.Body)
	if err != nil {
		return false, err
	}

	return success, nil
}

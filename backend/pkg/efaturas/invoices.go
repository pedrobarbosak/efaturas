package efaturas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"efaturas-xtreme/internal/service/domain"
	httpClient "efaturas-xtreme/pkg/efaturas/http"
	"efaturas-xtreme/pkg/efaturas/responses"

	"efaturas-xtreme/pkg/errors"
)

type invoicesURLs struct {
	list    string
	invoice string
}

func (s *service) GetInvoices(ctx context.Context, cookies map[string]string, dates ...Dates) ([]*domain.Invoice, int, error) {
	return s.GetInvoicesWithCategory(ctx, cookies, "TODOS", dates...)
}

func (s *service) GetInvoicesWithCategory(ctx context.Context, cookies map[string]string, category domain.Category, dates ...Dates) ([]*domain.Invoice, int, error) {
	date := defaultDate()
	if len(dates) != 0 {
		date = dates[0]
		if err := date.parse(); err != nil {
			return nil, 0, err
		}
	}

	url := fmt.Sprintf("%s&dataInicioFilter=%s&dataFimFilter=+%s&ambitoAquisicaoFilter=%s", s.urls.invoices.list, date.From, date.To, category)

	resp, err := httpClient.Get(ctx, url, cookies)
	if err != nil {
		return nil, 0, errors.New("failed to make getInvoices request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, 0, errors.New("failed to make getInvoices request:", resp.StatusCode, resp.Status)
	}

	var response responses.ListInvoices
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, errors.New("failed to read body:", err)
	}

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, 0, errors.New("failed to unmarshal to response:", err)
	}

	return response.ToDomain(), response.Total, nil
}

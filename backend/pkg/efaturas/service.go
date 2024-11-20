package efaturas

import (
	"context"
	"time"

	"efaturas-xtreme/internal/service/domain"
	"efaturas-xtreme/pkg/validator"

	"efaturas-xtreme/pkg/errors"
)

type Service interface {
	Login(ctx context.Context, uname string, pword string) (map[string]string, error)
	GetInvoices(ctx context.Context, cookies map[string]string, dates ...Dates) ([]*domain.Invoice, int, error)
	GetInvoicesWithCategory(ctx context.Context, cookies map[string]string, category domain.Category, dates ...Dates) ([]*domain.Invoice, int, error)
	CheckInvoice(ctx context.Context, cookies map[string]string, invoice *domain.Invoice, category domain.Category) (bool, error)
}

type service struct {
	urls urls
}

type urls struct {
	login    loginURLs
	invoices invoicesURLs
}

func New() Service {
	return &service{
		urls: urls{
			login: loginURLs{
				csrf:         "https://www.acesso.gov.pt/jsp/loginRedirectForm.jsp?path=painelAdquirente.action&partID=EFPF",
				login:        "https://www.acesso.gov.pt/jsp/submissaoFormularioLogin",
				authenticate: "https://faturas.portaldasfinancas.gov.pt/painelAdquirente.action",
			},
			invoices: invoicesURLs{
				list:    "https://faturas.portaldasfinancas.gov.pt/json/obterDocumentosIRSAdquirente.action?ambitoAquisicaoFilter=TODOS",
				invoice: "https://faturas.portaldasfinancas.gov.pt/resolverPendenciaAdquirente.action",
			},
		},
	}
}

type Dates struct {
	From string `validate:"omitempty,datetime=2006-01-02"`
	To   string `validate:"omitempty,datetime=2006-01-02"`
}

func (d *Dates) parse() error {
	if err := validator.Validate(d); err != nil {
		return errors.New("invalid input:", err)
	}

	dDate := defaultDate()

	if len(d.From) == 0 {
		d.From = dDate.From
	}

	if len(d.To) == 0 {
		d.To = dDate.To
	}

	return nil
}

func defaultDate() Dates {
	const layout = "2006-01-02"
	n := time.Now()

	return Dates{
		//	From: time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC).Format(layout),
		From: time.Date(n.Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format(layout),
		To:   time.Date(n.Year(), 12, 31, 0, 0, 0, 0, time.UTC).Format(layout),
	}
}

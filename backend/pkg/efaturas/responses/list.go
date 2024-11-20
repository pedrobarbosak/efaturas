package responses

import "efaturas-xtreme/internal/service/domain"

type ListInvoices struct {
	Success  bool       `json:"success"`
	Invoices []*Invoice `json:"linhas"`
	Count    uint       `json:"numElementos"`
	Total    int        `json:"totalElementos"`
}

func (i ListInvoices) ToDomain() []*domain.Invoice {
	invoices := make([]*domain.Invoice, 0, len(i.Invoices))
	for _, inv := range i.Invoices {
		invoices = append(invoices, inv.ToDomain())
	}
	return invoices
}

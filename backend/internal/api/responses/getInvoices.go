package responses

import (
	"efaturas-xtreme/internal/api/models"
	"efaturas-xtreme/internal/service/domain"
)

type ListInvoices []*models.Invoice

func NewListInvoices(invoices []*domain.Invoice) ListInvoices {
	data := make([]*models.Invoice, 0, len(invoices))
	for _, inv := range invoices {
		data = append(data, models.InvoiceFromDomain(inv))
	}
	return data
}

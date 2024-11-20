package responses

import (
	"efaturas-xtreme/internal/api/models"
	"efaturas-xtreme/internal/service/domain"
)

type InvoiceUpdate struct {
	Done    bool            `json:"done"`
	Invoice *models.Invoice `json:"invoice"`
}

func NewInvoiceUpdate(invoice *domain.Invoice, done bool) *InvoiceUpdate {
	return &InvoiceUpdate{Done: done, Invoice: models.InvoiceFromDomain(invoice)}
}

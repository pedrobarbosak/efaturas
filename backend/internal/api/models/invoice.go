package models

import "efaturas-xtreme/internal/service/domain"

type Invoice struct {
	ID         int64                       `json:"id"`
	Origin     Origin                      `json:"origin"`
	Issuer     Issuer                      `json:"issuer"`
	Buyer      Buyer                       `json:"buyer"`
	Document   Document                    `json:"document"`
	Total      Total                       `json:"total"`
	Activity   Activity                    `json:"activity"`
	ATCud      string                      `json:"atCud"`
	Tested     bool                        `json:"tested"`
	Categories map[domain.Category]*Values `json:"categories"`
}

func InvoiceFromDomain(invoice *domain.Invoice) *Invoice {
	categories := make(map[domain.Category]*Values, len(invoice.Categories))
	for cat, v := range invoice.Categories {
		categories[cat] = &Values{Success: v.Success, Benefit: v.Benefit, Others: v.Others}
	}

	return &Invoice{
		ID: invoice.ID,
		Origin: Origin{
			Value:       invoice.Origin.Value,
			Description: invoice.Origin.Description,
		},
		Issuer: Issuer{
			NIF:  invoice.Issuer.NIF,
			Name: invoice.Issuer.Name,
		},
		Buyer: Buyer{
			NIF:              invoice.Buyer.NIF,
			Name:             invoice.Buyer.Name,
			Country:          invoice.Buyer.Country,
			NIFInternational: invoice.Buyer.NIFInternational,
		},
		Document: Document{
			Type:        invoice.Document.Type,
			Description: invoice.Document.Description,
			Number:      invoice.Document.Number,
			Hash:        invoice.Document.Hash,
			Date:        invoice.Document.Date,
		},
		Total: Total{
			Value:         invoice.Total.Value,
			Taxable:       invoice.Total.Taxable,
			VAT:           invoice.Total.VAT,
			Benefit:       invoice.Total.Benefit,
			OthersBenefit: invoice.Total.OthersBenefit,
		},
		Activity: Activity{
			Category:    invoice.Activity.Category,
			Description: invoice.Activity.Description,
		},
		ATCud:      invoice.ATCud,
		Tested:     invoice.Tested,
		Categories: categories,
	}
}

type Activity struct {
	Category    domain.Category `json:"category"`
	Description string          `json:"description"`
}

type Total struct {
	Value         int `json:"value"`
	Taxable       int `json:"taxable"`
	VAT           int `json:"vat"`
	Benefit       int `json:"benefit"`
	OthersBenefit int `json:"othersBenefit"`
}

type Document struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Number      string `json:"number"`
	Hash        string `json:"hash"`
	Date        string `json:"date"`
}

type Buyer struct {
	NIF              uint   `json:"nif"`
	Name             string `json:"name"`
	Country          string `json:"country"`
	NIFInternational string `json:"nifInternational"`
}

type Issuer struct {
	NIF  uint   `json:"nif"`
	Name string `json:"name"`
}

type Origin struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

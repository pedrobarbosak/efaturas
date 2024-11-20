package responses

import (
	"efaturas-xtreme/internal/service/domain"
	"efaturas-xtreme/pkg/entity"
)

type Invoice struct {
	IdDocumento                 int64       `json:"idDocumento"`
	OrigemRegisto               string      `json:"origemRegisto"`
	OrigemRegistoDesc           string      `json:"origemRegistoDesc"`
	NifEmitente                 uint        `json:"nifEmitente"`
	NomeEmitente                string      `json:"nomeEmitente"`
	NifAdquirente               uint        `json:"nifAdquirente"`
	NomeAdquirente              string      `json:"nomeAdquirente"`
	PaisAdquirente              string      `json:"paisAdquirente"`
	NifAdquirenteInternac       string      `json:"nifAdquirenteInternac"`
	TipoDocumento               string      `json:"tipoDocumento"`
	TipoDocumentoDesc           string      `json:"tipoDocumentoDesc"`
	Numerodocumento             string      `json:"numerodocumento"`
	HashDocumento               string      `json:"hashDocumento"`
	DataEmissaoDocumento        string      `json:"dataEmissaoDocumento"`
	ValorTotal                  int         `json:"valorTotal"`
	ValorTotalBaseTributavel    int         `json:"valorTotalBaseTributavel"`
	ValorTotalIva               int         `json:"valorTotalIva"`
	ValorTotalBeneficioProv     int         `json:"valorTotalBeneficioProv"`
	ValorTotalSetorBeneficio    interface{} `json:"valorTotalSetorBeneficio"`
	ValorTotalDespesasGerais    int         `json:"valorTotalDespesasGerais"`
	EstadoBeneficio             string      `json:"estadoBeneficio"`
	EstadoBeneficioDesc         string      `json:"estadoBeneficioDesc"`
	EstadoBeneficioEmitente     string      `json:"estadoBeneficioEmitente"`
	EstadoBeneficioDescEmitente string      `json:"estadoBeneficioDescEmitente"`
	ExisteTaxaNormal            interface{} `json:"existeTaxaNormal"`
	ActividadeEmitente          string      `json:"actividadeEmitente"`
	ActividadeEmitenteDesc      string      `json:"actividadeEmitenteDesc"`
	ActividadeProf              interface{} `json:"actividadeProf"`
	ActividadeProfDesc          interface{} `json:"actividadeProfDesc"`
	ComunicacaoComerciante      bool        `json:"comunicacaoComerciante"`
	ComunicacaoConsumidor       bool        `json:"comunicacaoConsumidor"`
	IsDocumentoEstrangeiro      bool        `json:"isDocumentoEstrangeiro"`
	Atcud                       string      `json:"atcud"`
	Autofaturacao               bool        `json:"autofaturacao"`
}

func (inv *Invoice) ToDomain() *domain.Invoice {
	return &domain.Invoice{
		Entity: &entity.Entity{ID: inv.IdDocumento},
		Origin: domain.Origin{
			Value:       inv.OrigemRegisto,
			Description: inv.OrigemRegistoDesc,
		},
		Issuer: domain.Issuer{
			NIF:  inv.NifEmitente,
			Name: inv.NomeEmitente,
		},
		Buyer: domain.Buyer{
			NIF:              inv.NifAdquirente,
			Name:             inv.NomeAdquirente,
			Country:          inv.PaisAdquirente,
			NIFInternational: inv.NifAdquirenteInternac,
		},
		Document: domain.Document{
			Type:        inv.TipoDocumento,
			Description: inv.TipoDocumentoDesc,
			Number:      inv.Numerodocumento,
			Hash:        inv.HashDocumento,
			Date:        inv.DataEmissaoDocumento,
		},
		Total: domain.Total{
			Value:         inv.ValorTotal,
			Taxable:       inv.ValorTotalBaseTributavel,
			VAT:           inv.ValorTotalIva,
			Benefit:       inv.ValorTotalBeneficioProv,
			OthersBenefit: inv.ValorTotalDespesasGerais,
		},
		Activity: domain.Activity{
			Category:    domain.Category(inv.ActividadeEmitente),
			Description: inv.ActividadeEmitenteDesc,
		},
		ATCud: inv.Atcud,
	}
}

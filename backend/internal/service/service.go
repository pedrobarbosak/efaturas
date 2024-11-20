package service

import (
	"context"
	"log"
	"os"
	"slices"

	"efaturas-xtreme/internal/service/domain"
	"efaturas-xtreme/internal/service/repository"
	"efaturas-xtreme/pkg/efaturas"
	"efaturas-xtreme/pkg/errors"
	"efaturas-xtreme/pkg/sse"
)

type Service interface {
	GetInvoices(ctx context.Context, userID string) ([]*domain.Invoice, error)
	CreateOrUpdate(ctx context.Context, userID string, uname string, pword string) ([]*domain.Invoice, error)
	ScanInvoices(ctx context.Context, userID string, uname string, pword string) error
	UpdateInvoiceCategories(ctx context.Context, userID string, uname string, pword string, invoicesCategories map[int64]domain.Category) ([]*domain.Invoice, error)
	SetOnUpdateInvoice(onUpdate func(invoice *domain.Invoice, done bool))
}

type service struct {
	log      *log.Logger
	repo     repository.Repository
	efaturas efaturas.Service
	sse      sse.Publisher

	onUpdateInvoice []func(*domain.Invoice, bool)
}

func (s *service) GetInvoices(ctx context.Context, userID string) ([]*domain.Invoice, error) {
	invoices, err := s.repo.GetInvoicesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(invoices, func(a, b *domain.Invoice) int {
		switch {
		case a.Document.Date < b.Document.Date:
			return 1
		case a.Document.Date > b.Document.Date:
			return -1
		default:
			return 0
		}
	})

	return invoices, nil
}

func (s *service) CreateOrUpdate(ctx context.Context, userID string, uname string, pword string) ([]*domain.Invoice, error) {
	cookies, err := s.efaturas.Login(ctx, uname, pword)
	if err != nil {
		return nil, errors.NewUnauthorized("failed to login:", err)
	}

	invoices, _, err := s.efaturas.GetInvoices(ctx, cookies)
	if err != nil {
		return nil, errors.NewUnauthorized("failed to get new invoices:", err)
	}

	for _, inv := range invoices {
		if err = inv.Prepare(userID); err != nil {
			return nil, errors.NewUnauthorized("failed to prepare invoice to insert:", err)
		}
	}

	if err = s.repo.CreateOrUpdate(ctx, invoices); err != nil {
		return nil, errors.NewUnauthorized("failed to insert invoices:", err)
	}

	return s.GetInvoices(ctx, userID)
}

func (s *service) ScanInvoices(ctx context.Context, userID string, uname string, pword string) error {
	cookies, err := s.efaturas.Login(ctx, uname, pword)
	if err != nil {
		return err
	}

	invoices, err := s.repo.GetInvoicesByUserID(ctx, userID)
	if err != nil {
		return errors.New(err)
	}

	originalCategories := make(map[*domain.Invoice]domain.Category, len(invoices))
	for _, inv := range invoices {
		originalCategories[inv] = inv.Activity.Category
	}

	for _, category := range domain.CategoryList {
		s.log.Println("checking for category:", category)

		// Set all the invoices to a certain category
		for _, inv := range invoices {
			if inv.Tested {
				continue
			}

			if inv.TestedCategory(category) {
				continue
			}

			success, err := s.efaturas.CheckInvoice(ctx, cookies, inv, category)
			if err != nil {
				s.log.Println("failed to scan:", inv.ID, category, err)
				s.publish(inv, false)
				continue
			}

			if !success {
				inv.SetCategory(category, success, 0, 0)
				s.publish(inv, false)
			}
		}

		// Check the results ...
		updatedInvoices, _, err := s.efaturas.GetInvoices(ctx, cookies)
		if err != nil {
			s.log.Println("failed to check results from category:", category, err)
			continue
		}

		for _, updated := range updatedInvoices {
			for _, existing := range invoices {
				if updated.ID != existing.ID {
					continue
				}

				existing.SetCategory(updated.Activity.Category, true, int64(updated.Total.Benefit), int64(updated.Total.OthersBenefit))
				s.publish(existing, false)
			}

		}

		// Update database
		if err = s.repo.Update(ctx, invoices); err != nil {
			s.log.Println("failed to save updated invoices:", err)
			continue
		}
	}

	for inv, cat := range originalCategories {
		_, _ = s.efaturas.CheckInvoice(ctx, cookies, inv, cat)
	}

	return nil
}

func (s *service) UpdateInvoiceCategories(ctx context.Context, userID string, uname string, pword string, invoicesCategories map[int64]domain.Category) ([]*domain.Invoice, error) {
	cookies, err := s.efaturas.Login(ctx, uname, pword)
	if err != nil {
		return nil, err
	}

	invoices, err := s.repo.GetInvoicesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, inv := range invoices {
		cat, ok := invoicesCategories[inv.ID]
		if !ok {
			continue
		}

		_, err = s.efaturas.CheckInvoice(ctx, cookies, inv, cat)
		if err != nil {
			return nil, err
		}
	}

	return s.CreateOrUpdate(ctx, userID, uname, pword)
}

func (s *service) publish(inv *domain.Invoice, done bool) {
	for _, fn := range s.onUpdateInvoice {
		fn(inv, done)
	}
}

func (s *service) SetOnUpdateInvoice(fn func(invoice *domain.Invoice, done bool)) {
	s.onUpdateInvoice = append(s.onUpdateInvoice, fn)
}

func New(repo repository.Repository, efaturas efaturas.Service, sse sse.Publisher) Service {
	return &service{
		log:             log.New(os.Stdout, "<service> ", log.Flags()),
		repo:            repo,
		efaturas:        efaturas,
		sse:             sse,
		onUpdateInvoice: []func(*domain.Invoice, bool){},
	}
}

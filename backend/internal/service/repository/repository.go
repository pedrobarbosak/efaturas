package repository

import (
	"context"

	"efaturas-xtreme/internal/service/domain"
	"efaturas-xtreme/pkg/db"
	"efaturas-xtreme/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, invoices []*domain.Invoice) error
	Update(ctx context.Context, invoices []*domain.Invoice) error
	GetInvoicesByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error)
}

type repository struct {
	db         *db.DB
	collection string
}

func (r repository) CreateOrUpdate(ctx context.Context, invoices []*domain.Invoice) error {
	collection := r.db.Database.Collection(r.collection)

	for _, inv := range invoices {
		exists, err := r.exists(ctx, inv)
		if err != nil {
			return errors.New(err)
		}

		if !exists {
			if _, err = collection.InsertOne(ctx, inv); err != nil {
				return errors.New(err)
			}

			continue
		}

		filter := bson.M{"_id": inv.ID}
		update := bson.M{"$set": bson.M{"activity": inv.Activity, "total": inv.Total}}

		if _, err = collection.UpdateOne(ctx, filter, update); err != nil {
			return errors.New("failed to update:", err)
		}

		continue
	}

	return nil
}

func (r repository) exists(ctx context.Context, invoice *domain.Invoice) (bool, error) {
	count, err := r.db.Database.Collection(r.collection).CountDocuments(ctx, bson.M{"_id": invoice.ID})
	if err != nil {
		return false, errors.New(err)
	}

	return count != 0, nil
}

func (r repository) Update(ctx context.Context, invoices []*domain.Invoice) error {
	for _, inv := range invoices {
		filter := bson.M{"_id": inv.ID}
		update := bson.M{"$set": bson.M{
			"activity":   inv.Activity,
			"total":      inv.Total,
			"categories": inv.Categories,
			"tested":     inv.Tested,
			"testedat":   inv.TestedAt,
		}}

		if _, err := r.db.Database.Collection(r.collection).UpdateOne(ctx, filter, update); err != nil {
			return errors.New("failed to update:", err)
		}
	}

	return nil
}

func (r repository) GetInvoicesByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error) {
	filter := bson.M{"userid": userID}
	cursor, err := r.db.Database.Collection(r.collection).Find(ctx, filter)
	if err != nil {
		return nil, errors.New("failed to get:", err)
	}

	invoices := make([]*domain.Invoice, 0)
	for cursor.Next(ctx) {
		var inv domain.Invoice
		if err = cursor.Decode(&inv); err != nil {
			return nil, errors.New("failed to decode:", err)
		}

		invoices = append(invoices, &inv)
	}

	if err = cursor.Err(); err != nil {
		return nil, errors.New("failed:", err)
	}
	_ = cursor.Close(ctx)

	return invoices, nil
}

func New(db *db.DB) (Repository, error) {
	return &repository{db: db, collection: "invoices"}, nil
}

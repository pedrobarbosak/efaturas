package domain

import (
	"fmt"
	"time"

	"efaturas-xtreme/pkg/entity"
	"efaturas-xtreme/pkg/validator"

	"efaturas-xtreme/pkg/errors"
)

type Invoice struct {
	*entity.Entity `bson:",inline"`
	Origin         Origin
	Issuer         Issuer
	Buyer          Buyer
	Document       Document
	Total          Total
	Activity       Activity
	ATCud          string

	UserID     string `validate:"required"`
	Tested     bool
	TestedAt   int64
	Categories map[Category]Values
}

func (i *Invoice) Prepare(userID string) error {
	i.Entity = entity.New(i.ID)
	i.UserID = userID
	i.Tested = false
	i.TestedAt = 0
	i.Categories = nil

	if err := validator.Validate(i); err != nil {
		return errors.New(err)
	}

	return nil
}

func (i *Invoice) TestedCategory(category Category) bool {
	if i.Categories == nil {
		return false
	}

	_, ok := i.Categories[category]
	return ok
}

func (i *Invoice) SetCategory(category Category, success bool, benefit int64, others int64) {
	if i.Categories == nil {
		i.Categories = map[Category]Values{}
	}

	i.Categories[category] = Values{Success: success, Benefit: benefit, Others: others}

	if len(i.Categories) == len(CategoryList) {
		i.Tested = true
		i.TestedAt = time.Now().Unix()
	}
}

func (i *Invoice) GetID() string {
	return fmt.Sprintf("%d", i.ID)
}

func (i *Invoice) HasMultipleCategories() bool {
	count := 0
	for _, v := range i.Categories {
		if v.Success {
			count += 1
		}
	}
	return count > 1
}

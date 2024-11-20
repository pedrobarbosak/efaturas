package responses

import "efaturas-xtreme/internal/api/models"

type GetCategories []*models.CategoryInfo

func NewGetCategories() GetCategories {
	return models.CategoriesInfo
}

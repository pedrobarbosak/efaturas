package models

import "efaturas-xtreme/internal/service/domain"

type CategoryInfo struct {
	Name    domain.Category `json:"name"`
	Color   string          `json:"color"`
	Unicode string          `json:"unicode"`
}

var CategoriesInfo = []*CategoryInfo{
	{Name: domain.AutoMaintenance, Color: "#357db0", Unicode: "E600"},
	{Name: domain.MotoMaintenance, Color: "#55a3db", Unicode: "E609"},
	{Name: domain.LodgingAndRestaurants, Color: "#febc17", Unicode: "E60A"},
	{Name: domain.BeautySalon, Color: "#a98ddf", Unicode: "E601"},
	{Name: domain.Veterinary, Color: "#d07361", Unicode: "E900"},
	{Name: domain.PublicTransport, Color: "#309ffa", Unicode: "E624"},
	{Name: domain.Health, Color: "#ff5a59", Unicode: "E60B"},
	{Name: domain.Education, Color: "#ff893a", Unicode: "E602"},
	{Name: domain.RealEstate, Color: "#95d655", Unicode: "E606"},
	{Name: domain.NursingHomes, Color: "#6cb664", Unicode: "E608"},
	{Name: domain.Other, Color: "#48c1da", Unicode: "E604"},
	{Name: domain.Gym, Color: "#942192", Unicode: "E901"},
	{Name: domain.NewspapersAndMagazines, Color: "#c86c44", Unicode: "E904"},
}

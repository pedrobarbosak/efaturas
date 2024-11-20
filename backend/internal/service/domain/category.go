package domain

type Category string

func (c Category) String() string {
	return string(c)
}

const (
	Noop                   Category = ""
	AutoMaintenance        Category = "C01"
	MotoMaintenance        Category = "C02"
	LodgingAndRestaurants  Category = "C03"
	BeautySalon            Category = "C04"
	Veterinary             Category = "C09"
	PublicTransport        Category = "C10"
	Health                 Category = "C05"
	Education              Category = "C06"
	RealEstate             Category = "C07"
	NursingHomes           Category = "C08"
	Other                  Category = "C99"
	Gym                    Category = "C11"
	NewspapersAndMagazines Category = "C12"
)

var CategoryList = []Category{
	AutoMaintenance,
	MotoMaintenance,
	LodgingAndRestaurants,
	BeautySalon,
	Veterinary,
	PublicTransport,
	Health,
	Education,
	RealEstate,
	NursingHomes,
	Other,
	Gym,
	NewspapersAndMagazines,
}

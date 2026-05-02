package models

var AllowedCategories = map[string]bool{
	"coffee":    true,
	"nightlife": true,
	"art":       true,
	"food":      true,
}

func IsValidCategory(cat string) bool {
	return AllowedCategories[cat]
}

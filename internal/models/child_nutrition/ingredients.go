package models

type Ingredient struct {
	ID       string  `json:"ingredient_id" db:"id"`
	Name     string  `json:"name" db:"name"`
	ImageURL *string `json:"image_url" db:"image_url"`
	Category *string `json:"category" db:"category"`
	Price    *string `json:"price" db:"price"`
}

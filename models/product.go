package models

type Product struct {
	Index                   string   `graphql:"-"`
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	Price                   int64    `json:"price"`
	Stock                   int64    `json:"stock"`
	ProductSpecificDiscount int64    `json:"productSpecificDiscount"`
	FullImages              []string `json:"fullImages"`
}

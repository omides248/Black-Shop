package dto

import ozzo "github.com/go-ozzo/ozzo-validation/v4"

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type ProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r CreateProductRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Name, ozzo.Required, ozzo.Length(3, 100)),
		ozzo.Field(&r.Price, ozzo.Required),
	)
}

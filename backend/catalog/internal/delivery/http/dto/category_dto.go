package dto

import (
	"mime/multipart"
	"pkg/validation"
	"time"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateCategoryRequest struct {
	Name     string                `form:"name"`
	ParentID *string               `form:"parentId,omitempty"`
	Image    *multipart.FileHeader `form:"image,omitempty"`
}

type CategoryResponse struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Image         *string             `json:"image,omitempty"`
	ParentID      *string             `json:"parentId,omitempty"`
	Depth         int                 `json:"depth"`
	Subcategories []*CategoryResponse `json:"subcategory,omitempty"`
	CreatedAt     time.Time           `json:"createdAt"`
}

func (r *CreateCategoryRequest) Validate() error {
	return ozzo.ValidateStruct(r,
		ozzo.Field(&r.Name, ozzo.Required, ozzo.Length(3, 100)),
		ozzo.Field(&r.ParentID, ozzo.NilOrNotEmpty),
		ozzo.Field(&r.Image,
			ozzo.When(r.Image != nil,
				ozzo.By(validation.ImageRule(
					5*1024*1024, // 5 MB
					[]string{".jpg", ".jpeg", ".png"},
				)),
			),
		),
	)
}

type UpdateCategoryRequest struct {
	Name     *string               `form:"name"`
	ParentID *string               `form:"parent_id,omitempty"`
	Image    *multipart.FileHeader `form:"image,omitempty"`
}

func (r *UpdateCategoryRequest) Validate() error {
	return ozzo.ValidateStruct(r,
		ozzo.Field(&r.Name, ozzo.NilOrNotEmpty, ozzo.Length(3, 100)),
		ozzo.Field(&r.ParentID, ozzo.NilOrNotEmpty),
		ozzo.Field(&r.Image,
			ozzo.When(r.Image != nil,
				ozzo.By(validation.ImageRule(
					5*1024*1024, // 5 MB
					[]string{".jpg", ".jpeg", ".png"},
				)),
			),
		),
	)
}

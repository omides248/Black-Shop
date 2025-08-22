package error_mapping

import (
	"catalog/internal/domain"
	"net/http"
	"pkg/echo/error_handler"
)

func GetDomainErrorMappings() map[error]error_handler.DomainErrorMapping {
	domainErrorMappings := map[error]error_handler.DomainErrorMapping{
		domain.ErrProductNotFound:            {StatusCode: http.StatusNotFound, Message: domain.ErrProductNotFound.Error()},
		domain.ErrCategoryNotFound:           {StatusCode: http.StatusNotFound, Message: domain.ErrCategoryNotFound.Error()},
		domain.ErrCategoryAlreadyExists:      {StatusCode: http.StatusConflict, Message: domain.ErrCategoryAlreadyExists.Error()},
		domain.ErrCategoryDepthLimitExceeded: {StatusCode: http.StatusBadRequest, Message: domain.ErrCategoryDepthLimitExceeded.Error()},
		domain.ErrCategoryHasProducts:        {StatusCode: http.StatusBadRequest, Message: domain.ErrCategoryHasProducts.Error()},
	}

	return domainErrorMappings
}

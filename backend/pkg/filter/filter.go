package filter

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"strings"
)

type FilterSet struct {
	FilterFields   map[string]string
	SearchFields   []string
	OrderingFields []string
}

type QueryBuilderResult struct {
	FilterQuery bson.M
	SortOptions bson.D
}

func (fs *FilterSet) BuildMongoQuery(c echo.Context) QueryBuilderResult {
	queryParams := c.QueryParams()

	result := QueryBuilderResult{
		FilterQuery: bson.M{},
		SortOptions: bson.D{},
	}

	// Ordering
	ordering := queryParams.Get("ordering")
	if ordering != "" {
		sortFiled := ordering
		sortDirection := 1 // 1 for ASC, -1 for DESC

		if strings.HasPrefix(sortFiled, "-") {
			sortDirection = -1
			sortFiled = strings.TrimPrefix(sortFiled, "-")
		}

		for _, allowedField := range fs.OrderingFields {
			if sortFiled == allowedField {
				result.SortOptions = append(result.SortOptions, bson.E{Key: allowedField, Value: sortDirection})
				break
			}
		}
	} else {
		result.SortOptions = append(result.SortOptions, bson.E{Key: "_id", Value: 1}) // Default sort
	}

	// Search & Filter
	var searchClauses bson.A

	for key, values := range queryParams {
		if len(values) == 0 || values[0] == "" {
			continue
		}
		value := values[0]

		// Search Logic
		if (key == "p" || key == "search") && len(fs.SearchFields) > 0 {
			for _, field := range fs.SearchFields {
				searchClauses = append(searchClauses, bson.M{field: bson.M{"$regex": value, "$options": "i"}})
			}
			continue
		}

		// Filter (gte, lte, in, ...)
		parts := strings.Split(key, "__") // price__gte, price__lte
		if len(parts) == 2 {
			fieldName := parts[0] // price
			operator := parts[1]  //gte

			if dbField, ok := fs.FilterFields[fieldName]; ok {
				mongoOperator := "$" + operator // gte -> $gte

				if existingFilter, ok := result.FilterQuery[dbField].(bson.M); ok {
					existingFilter[mongoOperator] = value
				} else {
					result.FilterQuery[dbField] = bson.M{mongoOperator: value}
				}
				continue
			}
		}

		// Filter (equal)
		if dbField, ok := fs.FilterFields[key]; ok {
			result.FilterQuery[dbField] = value
		}
	}

	if len(searchClauses) > 0 {
		if len(result.FilterQuery) > 0 {
			result.FilterQuery = bson.M{"$and": []bson.M{{"$or": searchClauses}, result.FilterQuery}}
		} else {
			result.FilterQuery["$or"] = searchClauses
		}
	}

	return result
}

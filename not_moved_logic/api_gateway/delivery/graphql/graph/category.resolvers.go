package graph

import (
	pbcatalog "black-shop/api/proto/catalog/v1"
	"black-shop/internal/api_gateway/delivery/graphql/graph/model"
	"context"
	"fmt"
)

type categoryResolver struct{ *Resolver }

func (r *Resolver) Category() CategoryResolver { return &categoryResolver{r} }

func (r *categoryResolver) Children(ctx context.Context, obj *model.Category) ([]*model.Category, error) {
	if obj == nil {
		return nil, nil
	}
	res, err := r.Resolver.CatalogClient.ListCategories(ctx, &pbcatalog.ListCategoriesRequest{
		ParentId: obj.ID,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("res", res)
	categories := make([]*model.Category, 0, len(res.Categories))

	for _, cat := range res.Categories {
		if cat.ParentId != nil && *cat.ParentId == obj.ID {
			categories = append(categories, &model.Category{
				ID:       cat.Id,
				Name:     cat.Name,
				ImageURL: cat.ImageUrl,
				ParentID: cat.ParentId,
				Depth:    cat.Depth,
			})
		}
	}
	return categories, nil
}

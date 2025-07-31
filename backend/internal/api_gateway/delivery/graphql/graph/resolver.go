package graph

import pbcatalog "black-shop/api/proto/catalog/v1"

type Resolver struct {
	CatalogClient pbcatalog.CatalogServiceClient
}

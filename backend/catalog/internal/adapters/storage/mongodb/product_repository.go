package mongodb

import (
	"catalog/internal/adapters/storage/mongodb/model"
	"catalog/internal/domain"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type productRepo struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewProductRepository(db *mongo.Database, logger *zap.Logger) (domain.ProductRepository, error) {
	col := db.Collection(model.ProductsCollection)

	indexModel := mongo.IndexModel{
		Keys: bson.M{"categoryId": 1},
	}
	_, err := col.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create index on categoryId: %w", err)
	}

	return &productRepo{
		collection: col,
		logger:     logger.Named("mongodb_product_repo"),
	}, nil
}

func (r *productRepo) FindByID(ctx context.Context, id domain.ProductID) (*domain.Product, error) {

	r.logger.Info("finding product by id", zap.String("product_id", string(id)))

	objectID, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		r.logger.Error("failed to convert objectID from hex", zap.String("product_id", string(id)), zap.Error(err))
		return nil, domain.ErrProductNotFound
	}

	var p model.MongoProduct
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("product not found", zap.String("product_id", string(id)), zap.Error(err))
			return nil, domain.ErrProductNotFound
		}
		r.logger.Error("failed to find product by id", zap.String("product_id", string(id)), zap.Error(err))
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}

	return &domain.Product{
		ID:   domain.ProductID(p.ID.Hex()),
		Name: p.Name,
	}, nil
}

func (r *productRepo) FindAll(ctx context.Context, filterQuery bson.M, sortOptions bson.D, page, limit int) ([]*domain.Product, int64, error) {
	r.logger.Info("finding all products")

	offset := (page - 1) * limit

	matchStage := bson.D{{"$match", filterQuery}}
	dataStage := bson.A{
		bson.D{{"$sort", sortOptions}},
		bson.D{{"$skip", int64(offset)}},
		bson.D{{"$limit", int64(limit)}},
	}
	countStage := bson.A{
		bson.D{{"$count", "count"}},
	}
	facetStage := bson.D{{"$facet", bson.M{
		"data":       dataStage,
		"totalCount": countStage,
	}}}

	pipeline := mongo.Pipeline{matchStage, facetStage}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error("failed to execute aggregate query for products", zap.Error(err))
		return nil, 0, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	var results struct {
		Data       []model.MongoProduct `bson:"data"`
		TotalCount []struct {
			Count int64 `bson:"count"`
		} `bson:"totalCount"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&results); err != nil {
			r.logger.Error("failed to decode product aggregation result", zap.Error(err))
			return nil, 0, err
		}
	}

	if len(results.Data) == 0 {
		return []*domain.Product{}, 0, nil
	}

	var total int64
	if len(results.TotalCount) > 0 {
		total = results.TotalCount[0].Count
	}

	products := make([]*domain.Product, 0, len(results.Data))
	for _, mongoProd := range results.Data {
		products = append(products, &domain.Product{
			ID:   domain.ProductID(mongoProd.ID.Hex()),
			Name: mongoProd.Name,
		})
	}

	return products, total, nil
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	r.logger.Info("creating a new product", zap.String("product_name", product.Name))

	p := model.MongoProduct{
		Name: product.Name,
	}

	res, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		r.logger.Error("failed to save product", zap.Error(err))
		return fmt.Errorf("failed to insert product into db: %w", err)
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		product.ID = domain.ProductID(oid.Hex())
		r.logger.Info("successfully created product", zap.String("product_id", string(product.ID)))
	}

	return nil
}

func (r *productRepo) CategoryHasProducts(ctx context.Context, id domain.CategoryID) (bool, error) {
	//TODO Add Index for performance
	filter := bson.M{"categoryId": id}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("failed to count products by category", zap.String("categoryId", string(id)), zap.Error(err))
		return false, fmt.Errorf("failed to count products by category: %w", err)

	}
	return count > 0, nil
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

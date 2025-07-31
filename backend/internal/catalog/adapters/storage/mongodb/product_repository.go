package mongodb

import (
	"black-shop/internal/catalog/domain"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type mongoProduct struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}
type productRepo struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewProductRepository(db *mongo.Database, logger *zap.Logger) (domain.ProductRepository, error) {
	collection := db.Collection("products")

	indexModel := mongo.IndexModel{
		Keys: bson.M{"categoryId": 1},
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create index on categoryId: %w", err)
	}

	return &productRepo{
		collection: collection,
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

	var p mongoProduct
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

func (r *productRepo) FindAll(ctx context.Context) ([]*domain.Product, error) {
	r.logger.Info("finding all products")

	var domainProducts []*domain.Product

	filter := bson.M{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("failed to execute find all query", zap.Error(err))
		return nil, fmt.Errorf("failed to execute find all query: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var p mongoProduct
		if err := cursor.Decode(&p); err != nil {
			r.logger.Error("failed to decode a product document", zap.Error(err))
			continue
		}
		domainProducts = append(domainProducts, &domain.Product{
			ID:   domain.ProductID(p.ID.Hex()),
			Name: p.Name,
		})
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("cursor error after iterating", zap.Error(err))
		return nil, fmt.Errorf("cursor iteration failed: %w", err)
	}

	r.logger.Info("successfully found all products", zap.Int("count", len(domainProducts)))
	return domainProducts, nil
}

func (r *productRepo) Save(ctx context.Context, product *domain.Product) error {
	r.logger.Info("creating a new product", zap.String("product_name", product.Name))

	p := mongoProduct{
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

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

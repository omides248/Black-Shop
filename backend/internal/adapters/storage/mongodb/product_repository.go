package mongodb

import (
	"black-shop-service/internal/domain/catalog"
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

func NewProductRepository(db *mongo.Database, logger *zap.Logger) catalog.ProductRepository {
	collection := db.Collection("products")
	return &productRepo{
		collection: collection,
		logger:     logger.Named("mongodb_repo"),
	}
}

func (r *productRepo) FindByID(ctx context.Context, id catalog.ProductID) (*catalog.Product, error) {

	r.logger.Info("finding product by id", zap.String("product_id", string(id)))

	objectID, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		r.logger.Error("failed to convert objectID from hex", zap.String("product_id", string(id)), zap.Error(err))
		return nil, catalog.ErrProductNotFound
	}

	var p mongoProduct
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("product not found", zap.String("product_id", string(id)), zap.Error(err))
			return nil, catalog.ErrProductNotFound
		}
		r.logger.Error("failed to find product by id", zap.String("product_id", string(id)), zap.Error(err))
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}

	return &catalog.Product{
		ID:   catalog.ProductID(p.ID.Hex()),
		Name: p.Name,
	}, nil
}

func (r *productRepo) FindAll(ctx context.Context) ([]*catalog.Product, error) {
	r.logger.Info("finding all products")

	var domainProducts []*catalog.Product

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
		domainProducts = append(domainProducts, &catalog.Product{
			ID:   catalog.ProductID(p.ID.Hex()),
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

func (r *productRepo) Save(ctx context.Context, product *catalog.Product) error {
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
		product.ID = catalog.ProductID(oid.Hex())
		r.logger.Info("successfully created product", zap.String("product_id", string(product.ID)))
	}

	return nil
}

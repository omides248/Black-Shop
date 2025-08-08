package mongodb

import (
	"catalog/internal/adapters/storage/mongodb/model"
	"catalog/internal/domain"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *productRepo) FindAll(ctx context.Context, page, limit int) ([]*domain.Product, int64, error) {
	r.logger.Info("finding all products")

	offset := (page - 1) * limit

	// TODO Use one aggregate query instead collection.CountDocuments() and  r.collection.Find(ctx, filter, opts)

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	var domainProducts []*domain.Product

	filter := bson.M{}
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("failed to execute find all query", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to execute find all query: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var p model.MongoProduct
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
		return nil, 0, fmt.Errorf("cursor iteration failed: %w", err)
	}

	r.logger.Info("successfully found all products", zap.Int("count", len(domainProducts)))
	return domainProducts, total, nil
}

func (r *productRepo) Save(ctx context.Context, product *domain.Product) error {
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

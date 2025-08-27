package mongodb

import (
	"catalog/internal/adapters/storage/mongodb/model"
	"catalog/internal/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type categoryRepo struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewCategoryRepository(db *mongo.Database, logger *zap.Logger) (domain.CategoryRepository, error) {
	col := db.Collection(model.CategoriesCollection)

	if err := ensureCategoryIndexes(context.Background(), col); err != nil {
		return nil, err
	}

	return &categoryRepo{
		collection: col,
		logger:     logger.Named("mongodb_category_repo"),
	}, nil

}

func ensureCategoryIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "parentId", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	slugIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "slug", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetSparse(true),
	}

	_, err = collection.Indexes().CreateOne(ctx, slugIndexModel)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepo) Create(ctx context.Context, category *domain.Category) error {
	r.logger.Info("creating a new category", zap.String("category_name", category.Name))

	mc := model.MongoCategory{
		Name:      category.Name,
		Slug:      category.Slug,
		Image:     category.Image,
		ParentID:  category.ParentID,
		Depth:     category.Depth,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := r.collection.InsertOne(ctx, mc)
	if err != nil {
		r.logger.Error("failed to create category", zap.Error(err))
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		category.ID = domain.CategoryID(oid.Hex())
		r.logger.Info("successfully created category", zap.String("category_id", string(category.ID)))
	}

	return nil
}

func (r *categoryRepo) Update(ctx context.Context, category *domain.Category) error {
	r.logger.Info("updating category", zap.String("category_id", string(category.ID)))

	oid, err := toObjectID(string(category.ID))
	if err != nil {
		return domain.ErrCategoryNotFound
	}

	mc := model.MongoCategory{
		ID:        oid,
		Name:      category.Name,
		Image:     category.Image,
		ParentID:  category.ParentID,
		Depth:     category.Depth,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": mc}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed to update category", zap.Error(err))
		return err
	}

	if res.MatchedCount == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepo) FindByID(ctx context.Context, id domain.CategoryID) (*domain.Category, error) {
	r.logger.Info("finding category by id", zap.String("category_id", string(id)))

	oid, err := toObjectID(string(id))
	if err != nil {
		return nil, domain.ErrCategoryNotFound
	}

	var mc model.MongoCategory
	filter := bson.M{"_id": oid}

	if err := r.collection.FindOne(ctx, filter).Decode(&mc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("category not found", zap.String("category_id", string(id)))
			return nil, domain.ErrCategoryNotFound
		}
		r.logger.Error("failed to find category by id", zap.Error(err))
		return nil, err
	}

	return &domain.Category{
		ID:        domain.CategoryID(mc.ID.Hex()),
		Name:      mc.Name,
		Image:     mc.Image,
		ParentID:  mc.ParentID,
		Depth:     mc.Depth,
		CreatedAt: mc.CreatedAt,
		UpdatedAt: mc.UpdatedAt,
	}, nil
}

func (r *categoryRepo) FindAll(ctx context.Context) ([]*domain.Category, error) {
	r.logger.Info("finding all categories")
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		r.logger.Error("failed to execute find all categories query", zap.Error(err))
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	var categories []*domain.Category
	for cursor.Next(ctx) {
		var mc model.MongoCategory
		if err := cursor.Decode(&mc); err != nil {
			r.logger.Error("failed to decode a category document", zap.Error(err))
			continue
		}
		categories = append(categories, &domain.Category{
			ID:        domain.CategoryID(mc.ID.Hex()),
			Name:      mc.Name,
			Image:     mc.Image,
			ParentID:  mc.ParentID,
			Depth:     mc.Depth,
			CreatedAt: mc.CreatedAt,
			UpdatedAt: mc.UpdatedAt,
		})
	}

	return categories, cursor.Err()
}

func (r *categoryRepo) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {

	r.logger.Info("finding category by slug", zap.String("slug", slug))

	var mc model.MongoCategory
	filter := bson.M{"slug": slug}

	if err := r.collection.FindOne(ctx, filter).Decode(&mc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("category not found by slug", zap.String("slug", slug))
			return nil, domain.ErrCategoryNotFound
		}
		r.logger.Error("failed to find category by slug", zap.Error(err))
	}

	return &domain.Category{
		ID:        domain.CategoryID(mc.ID.Hex()),
		Name:      mc.Name,
		Image:     mc.Image,
		Slug:      mc.Slug,
		ParentID:  mc.ParentID,
		Depth:     mc.Depth,
		CreatedAt: mc.CreatedAt,
		UpdatedAt: mc.UpdatedAt,
	}, nil
}

func (r *categoryRepo) HasChildren(ctx context.Context, id domain.CategoryID) (bool, error) {
	filter := bson.M{"parentId": id}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("failed to count child categories", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *categoryRepo) FindByNameAndParentID(ctx context.Context, name string, parentID *domain.CategoryID) (*domain.Category, error) {
	filter := bson.M{
		"name":     name,
		"parentId": parentID,
	}

	var mc model.MongoCategory
	if err := r.collection.FindOne(ctx, filter).Decode(&mc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return &domain.Category{
		ID:        domain.CategoryID(mc.ID.Hex()),
		Name:      mc.Name,
		Image:     mc.Image,
		ParentID:  mc.ParentID,
		Depth:     mc.Depth,
		CreatedAt: mc.CreatedAt,
		UpdatedAt: mc.UpdatedAt,
	}, nil
}

func toObjectID(id string) (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(id)
}

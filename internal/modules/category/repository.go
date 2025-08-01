package category

import (
	"context"
	"gorm.io/gorm"
)

// repository implements the Repository interface using GORM
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new category repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new category and returns the created category
func (r *repository) Create(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByID retrieves a category by its ID
func (r *repository) GetByID(ctx context.Context, id uint) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByName retrieves a category by its name
func (r *repository) GetByName(ctx context.Context, name string) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// List retrieves a paginated list of categories
func (r *repository) List(ctx context.Context, offset, limit int) ([]*CategoryList, error) {
	var categories []*CategoryList
	err := r.db.WithContext(ctx).
		Model(&Category{}).
		Select("id, name, description, created_at").
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetAll retrieves all categories (for dropdowns)
func (r *repository) GetAll(ctx context.Context) ([]*CategoryList, error) {
	var categories []*CategoryList
	err := r.db.WithContext(ctx).
		Model(&Category{}).
		Select("id, name, description, created_at").
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Count returns the total number of categories
func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Category{}).Count(&count).Error
	return count, err
}

// Update updates a category by ID
func (r *repository) Update(ctx context.Context, id uint, updates *CategoryUpdate) error {
	updateMap := make(map[string]interface{})

	if updates.Name != nil {
		updateMap["name"] = *updates.Name
	}
	if updates.Description != nil {
		updateMap["description"] = *updates.Description
	}

	if len(updateMap) == 0 {
		return nil // No updates to perform
	}

	return r.db.WithContext(ctx).Model(&Category{}).Where("id = ?", id).Updates(updateMap).Error
}

// Delete soft deletes a category by ID
func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Category{}, id).Error
}

// Exists checks if a category name already exists
func (r *repository) Exists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Category{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

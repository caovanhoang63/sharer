package page

import (
	"context"
	"gorm.io/gorm"
)

// repository implements the Repository interface using GORM
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new page repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new page and returns the created page
func (r *repository) Create(ctx context.Context, page *Page) error {
	return r.db.WithContext(ctx).Create(page).Error
}

// GetBySlug retrieves a page by its slug
func (r *repository) GetBySlug(ctx context.Context, slug string) (*Page, error) {
	var page Page
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetByID retrieves a page by its ID
func (r *repository) GetByID(ctx context.Context, id uint) (*Page, error) {
	var page Page
	err := r.db.WithContext(ctx).First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// List retrieves a paginated list of pages
func (r *repository) List(ctx context.Context, offset, limit int) ([]*PageList, error) {
	var pages []*PageList
	err := r.db.WithContext(ctx).
		Table("shared_content p").
		Select("p.id, p.slug, p.title, p.category_id, c.name as category_name, p.created_at").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Order("p.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&pages).Error

	if err != nil {
		return nil, err
	}
	return pages, nil
}

// ListByCategory retrieves a paginated list of pages filtered by category
func (r *repository) ListByCategory(ctx context.Context, categoryID uint, offset, limit int) ([]*PageList, error) {
	var pages []*PageList
	err := r.db.WithContext(ctx).
		Table("shared_content p").
		Select("p.id, p.slug, p.title, p.category_id, c.name as category_name, p.created_at").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Where("p.category_id = ?", categoryID).
		Order("p.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&pages).Error

	if err != nil {
		return nil, err
	}
	return pages, nil
}

// Count returns the total number of pages
func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Page{}).Count(&count).Error
	return count, err
}

// CountByCategory returns the total number of pages in a category
func (r *repository) CountByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Page{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// Update updates a page by ID
func (r *repository) Update(ctx context.Context, id uint, updates *PageUpdate) error {
	updateMap := make(map[string]interface{})

	if updates.HTMLContent != nil {
		updateMap["html_content"] = *updates.HTMLContent
	}
	if updates.Title != nil {
		updateMap["title"] = *updates.Title
	}

	if len(updateMap) == 0 {
		return nil // No updates to perform
	}

	return r.db.WithContext(ctx).Model(&Page{}).Where("id = ?", id).Updates(updateMap).Error
}

// Delete soft deletes a page by ID
func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Page{}, id).Error
}

// Exists checks if a slug already exists
func (r *repository) Exists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Page{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

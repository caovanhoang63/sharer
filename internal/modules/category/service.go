package category

import (
	"context"
	"strings"
)

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new category service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateCategory creates a new category
func (s *service) CreateCategory(ctx context.Context, req *CategoryCreate) (*CategoryResponse, error) {
	// Validate category name
	if strings.TrimSpace(req.Name) == "" {
		return &CategoryResponse{Error: "Category name is required"}, nil
	}

	// Check if category name already exists
	exists, err := s.repo.Exists(ctx, req.Name)
	if err != nil {
		return &CategoryResponse{Error: "Error checking category existence"}, err
	}
	if exists {
		return &CategoryResponse{Error: "Category name already exists"}, nil
	}

	// Create category model
	category := &Category{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	}

	// Save to repository
	if err := s.repo.Create(ctx, category); err != nil {
		return &CategoryResponse{Error: "Error creating category"}, err
	}

	// Return created category
	categoryDetail := &CategoryDetail{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	return &CategoryResponse{Category: categoryDetail}, nil
}

// GetCategoryByID retrieves a category by its ID
func (s *service) GetCategoryByID(ctx context.Context, id uint) (*CategoryDetail, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &CategoryDetail{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// GetCategoriesList retrieves a paginated list of categories
func (s *service) GetCategoriesList(ctx context.Context, page, pageSize int) ([]*CategoryList, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	categories, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// GetAllCategories retrieves all categories for dropdowns
func (s *service) GetAllCategories(ctx context.Context) ([]*CategoryList, error) {
	return s.repo.GetAll(ctx)
}

// UpdateCategory updates a category
func (s *service) UpdateCategory(ctx context.Context, id uint, req *CategoryUpdate) (*CategoryResponse, error) {
	// Check if category exists
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &CategoryResponse{Error: "Category not found"}, err
	}

	// Validate name if provided
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return &CategoryResponse{Error: "Category name cannot be empty"}, nil
		}

		// Check if new name already exists (excluding current category)
		if name != category.Name {
			exists, err := s.repo.Exists(ctx, name)
			if err != nil {
				return &CategoryResponse{Error: "Error checking category existence"}, err
			}
			if exists {
				return &CategoryResponse{Error: "Category name already exists"}, nil
			}
		}
		*req.Name = name
	}

	// Trim description if provided
	if req.Description != nil {
		*req.Description = strings.TrimSpace(*req.Description)
	}

	// Update category
	if err := s.repo.Update(ctx, id, req); err != nil {
		return &CategoryResponse{Error: "Error updating category"}, err
	}

	// Get updated category
	updatedCategory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &CategoryResponse{Error: "Error retrieving updated category"}, err
	}

	categoryDetail := &CategoryDetail{
		ID:          updatedCategory.ID,
		Name:        updatedCategory.Name,
		Description: updatedCategory.Description,
		CreatedAt:   updatedCategory.CreatedAt,
		UpdatedAt:   updatedCategory.UpdatedAt,
	}

	return &CategoryResponse{Category: categoryDetail}, nil
}

// DeleteCategory deletes a category
func (s *service) DeleteCategory(ctx context.Context, id uint) error {
	// Check if category exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if category is being used by any pages
	// This should be implemented after updating the page module

	return s.repo.Delete(ctx, id)
}

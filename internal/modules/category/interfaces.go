package category

import "context"

// Repository defines the interface for category data access operations
type Repository interface {
	// Create creates a new category and returns the created category
	Create(ctx context.Context, category *Category) error

	// GetByID retrieves a category by its ID
	GetByID(ctx context.Context, id uint) (*Category, error)

	// GetByName retrieves a category by its name
	GetByName(ctx context.Context, name string) (*Category, error)

	// List retrieves a paginated list of categories
	List(ctx context.Context, offset, limit int) ([]*CategoryList, error)

	// GetAll retrieves all categories (for dropdowns)
	GetAll(ctx context.Context) ([]*CategoryList, error)

	// Count returns the total number of categories
	Count(ctx context.Context) (int64, error)

	// Update updates a category by ID
	Update(ctx context.Context, id uint, updates *CategoryUpdate) error

	// Delete soft deletes a category by ID
	Delete(ctx context.Context, id uint) error

	// Exists checks if a category name already exists
	Exists(ctx context.Context, name string) (bool, error)
}

// Service defines the interface for category business logic operations
type Service interface {
	// CreateCategory creates a new category
	CreateCategory(ctx context.Context, req *CategoryCreate) (*CategoryResponse, error)

	// GetCategoryByID retrieves a category by its ID
	GetCategoryByID(ctx context.Context, id uint) (*CategoryDetail, error)

	// GetCategoriesList retrieves a paginated list of categories
	GetCategoriesList(ctx context.Context, page, pageSize int) ([]*CategoryList, int64, error)

	// GetAllCategories retrieves all categories for dropdowns
	GetAllCategories(ctx context.Context) ([]*CategoryList, error)

	// UpdateCategory updates a category
	UpdateCategory(ctx context.Context, id uint, req *CategoryUpdate) (*CategoryResponse, error)

	// DeleteCategory deletes a category
	DeleteCategory(ctx context.Context, id uint) error
}

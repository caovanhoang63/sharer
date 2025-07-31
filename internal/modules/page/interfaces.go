package page

import "context"

// Repository defines the interface for page data access operations
type Repository interface {
	// Create creates a new page and returns the created page
	Create(ctx context.Context, page *Page) error

	// GetBySlug retrieves a page by its slug
	GetBySlug(ctx context.Context, slug string) (*Page, error)

	// GetByID retrieves a page by its ID
	GetByID(ctx context.Context, id uint) (*Page, error)

	// List retrieves a paginated list of pages
	List(ctx context.Context, offset, limit int) ([]*PageList, error)

	// Count returns the total number of pages
	Count(ctx context.Context) (int64, error)

	// Update updates a page by ID
	Update(ctx context.Context, id uint, updates *PageUpdate) error

	// Delete soft deletes a page by ID
	Delete(ctx context.Context, id uint) error

	// Exists checks if a slug already exists
	Exists(ctx context.Context, slug string) (bool, error)
}

// Service defines the interface for page business logic operations
type Service interface {
	// CreatePage creates a new shared page
	CreatePage(ctx context.Context, req *PageCreate) (*PageResponse, error)

	// GetPageBySlug retrieves a page by its slug for viewing
	GetPageBySlug(ctx context.Context, slug string) (*PageDetail, error)

	// GetPagesList retrieves a paginated list of pages
	GetPagesList(ctx context.Context, page, pageSize int) ([]*PageList, int64, error)

	// GenerateUniqueSlug generates a unique slug for a new page
	GenerateUniqueSlug(ctx context.Context) (string, error)

	// ExtractTitle attempts to extract title from HTML content
	ExtractTitle(htmlContent string) string
}

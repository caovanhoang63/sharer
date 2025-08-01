package page

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new page service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreatePage creates a new shared page
func (s *service) CreatePage(ctx context.Context, req *PageCreate) (*PageResponse, error) {
	// Validate HTML content
	if strings.TrimSpace(req.HTMLContent) == "" {
		return &PageResponse{Error: "No HTML content provided"}, nil
	}

	// Generate unique slug
	slug, err := s.GenerateUniqueSlug(ctx)
	if err != nil {
		return &PageResponse{Error: "Error generating unique slug"}, err
	}

	// Extract title if not provided
	title := req.Title
	if title == "" {
		title = s.ExtractTitle(req.HTMLContent)
	}

	// Create page model
	page := &Page{
		Slug:        slug,
		HTMLContent: req.HTMLContent,
		Title:       title,
		CategoryID:  req.CategoryID,
	}

	// Save to repository
	if err := s.repo.Create(ctx, page); err != nil {
		return &PageResponse{Error: "Error saving content"}, err
	}

	return &PageResponse{URL: "/shared/" + slug}, nil
}

// GetPageBySlug retrieves a page by its slug for viewing
func (s *service) GetPageBySlug(ctx context.Context, slug string) (*PageDetail, error) {
	page, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return &PageDetail{
		ID:          page.ID,
		Slug:        page.Slug,
		HTMLContent: page.HTMLContent,
		Title:       page.Title,
		CreatedAt:   page.CreatedAt,
		UpdatedAt:   page.UpdatedAt,
	}, nil
}

// GetPagesList retrieves a paginated list of pages
func (s *service) GetPagesList(ctx context.Context, page, pageSize int) ([]*PageList, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	pages, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// GetPagesByCategory retrieves a paginated list of pages filtered by category
func (s *service) GetPagesByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]*PageList, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	pages, err := s.repo.ListByCategory(ctx, categoryID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByCategory(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// GenerateUniqueSlug generates a unique slug for a new page
func (s *service) GenerateUniqueSlug(ctx context.Context) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	const maxAttempts = 10

	rand.Seed(time.Now().UnixNano())

	for attempt := 0; attempt < maxAttempts; attempt++ {
		slug := make([]byte, length)
		for i := range slug {
			slug[i] = charset[rand.Intn(len(charset))]
		}

		slugStr := string(slug)
		exists, err := s.repo.Exists(ctx, slugStr)
		if err != nil {
			return "", err
		}

		if !exists {
			return slugStr, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique slug after %d attempts", maxAttempts)
}

// ExtractTitle attempts to extract title from HTML content
func (s *service) ExtractTitle(htmlContent string) string {
	// Try to extract title from <title> tag
	titleRegex := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	matches := titleRegex.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		if title != "" {
			return title
		}
	}

	// Try to extract from first <h1> tag
	h1Regex := regexp.MustCompile(`<h1[^>]*>([^<]+)</h1>`)
	matches = h1Regex.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		if title != "" {
			return title
		}
	}

	// Default title
	return "Shared HTML Page"
}

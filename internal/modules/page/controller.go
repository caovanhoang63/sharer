package page

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sharer/views/components"
	"sharer/views/pages"
)

// Controller handles HTTP requests for page operations
type Controller struct {
	service Service
}

// NewController creates a new page controller
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// Home handles the home page display
func (c *Controller) Home(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	pages.Home().Render(ctx.Request.Context(), ctx.Writer)
}

// Index handles the index page showing list of shared pages
func (c *Controller) Index(ctx *gin.Context) {
	page := 1
	if p := ctx.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := ctx.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Check for category filter
	var pagesList []*pages.PageData
	var total int64
	var err error

	if categoryIDStr := ctx.Query("category"); categoryIDStr != "" {
		if categoryID, parseErr := strconv.ParseUint(categoryIDStr, 10, 32); parseErr == nil {
			pagesListRaw, totalRaw, serviceErr := c.service.GetPagesByCategory(ctx.Request.Context(), uint(categoryID), page, pageSize)
			pagesList = make([]*pages.PageData, len(pagesListRaw))
			for i, p := range pagesListRaw {
				pagesList[i] = &pages.PageData{
					ID:           p.ID,
					Slug:         p.Slug,
					Title:        p.Title,
					CategoryID:   p.CategoryID,
					CategoryName: p.CategoryName,
					CreatedAt:    p.CreatedAt,
				}
			}
			total = totalRaw
			err = serviceErr
		} else {
			ctx.Status(http.StatusBadRequest)
			return
		}
	} else {
		pagesListRaw, totalRaw, serviceErr := c.service.GetPagesList(ctx.Request.Context(), page, pageSize)
		pagesList = make([]*pages.PageData, len(pagesListRaw))
		for i, p := range pagesListRaw {
			pagesList[i] = &pages.PageData{
				ID:           p.ID,
				Slug:         p.Slug,
				Title:        p.Title,
				CategoryID:   p.CategoryID,
				CategoryName: p.CategoryName,
				CreatedAt:    p.CreatedAt,
			}
		}
		total = totalRaw
		err = serviceErr
	}

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Convert PageList to PageData
	pagesData := make([]*pages.PageData, len(pagesList))
	for i, p := range pagesList {
		pagesData[i] = &pages.PageData{
			ID:        p.ID,
			Slug:      p.Slug,
			Title:     p.Title,
			CreatedAt: p.CreatedAt,
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	hasNext := page < int(totalPages)
	hasPrev := page > 1

	ctx.Header("Content-Type", "text/html")
	pages.Index(pagesData, page, totalPages, total, hasNext, hasPrev).Render(ctx.Request.Context(), ctx.Writer)
}

// CreateFromForm handles form submission for creating pages
func (c *Controller) CreateFromForm(ctx *gin.Context) {
	var htmlContent string

	// Priority: textarea content over file
	textareaContent := ctx.PostForm("htmlContent")
	if strings.TrimSpace(textareaContent) != "" {
		htmlContent = textareaContent
	} else {
		// Try to read from uploaded file
		file, err := ctx.FormFile("htmlFile")
		if err == nil {
			// Validate file extension
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".html" && ext != ".htm" {
				ctx.String(http.StatusBadRequest, "Please upload an HTML file")
				return
			}

			src, err := file.Open()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Error reading file")
				return
			}
			defer src.Close()

			content, err := io.ReadAll(src)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Error reading file")
				return
			}
			htmlContent = string(content)
		}
	}

	if strings.TrimSpace(htmlContent) == "" {
		ctx.String(http.StatusBadRequest, "No HTML content provided")
		return
	}

	// Get category ID from form if provided
	var categoryID *uint
	if categoryIDStr := ctx.PostForm("category_id"); categoryIDStr != "" {
		if parsed, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			id := uint(parsed)
			categoryID = &id
		}
	}

	// Create page using service
	req := &PageCreate{
		HTMLContent: htmlContent,
		CategoryID:  categoryID,
	}
	response, err := c.service.CreatePage(ctx.Request.Context(), req)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error creating page")
		return
	}

	if response.Error != "" {
		ctx.String(http.StatusBadRequest, response.Error)
		return
	}

	// Return success component for htmx or redirect for regular form
	if ctx.GetHeader("HX-Request") == "true" {
		fullURL := "http://" + ctx.Request.Host + response.URL
		ctx.Header("Content-Type", "text/html")
		components.Success(fullURL).Render(ctx.Request.Context(), ctx.Writer)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/?success="+strings.TrimPrefix(response.URL, "/shared/"))
	}
}

// CreateFromAPI handles API requests for creating pages
func (c *Controller) CreateFromAPI(ctx *gin.Context) {
	var req PageCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	response, err := c.service.CreatePage(ctx.Request.Context(), &req)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	if response.Error != "" {
		ctx.String(http.StatusBadRequest, response.Error)
		return
	}

	// Return success component for htmx
	fullURL := ctx.Request.Host + response.URL
	ctx.Header("Content-Type", "text/html")
	components.Success(fullURL).Render(ctx.Request.Context(), ctx.Writer)
}

// GetSharedContent handles requests to view shared content
func (c *Controller) GetSharedContent(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	page, err := c.service.GetPageBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.serve404(ctx)
		} else {
			ctx.String(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page.HTMLContent))
}

// serve404 renders a 404 error page
func (c *Controller) serve404(ctx *gin.Context) {
	ctx.Status(http.StatusNotFound)
	ctx.Header("Content-Type", "text/html")
	pages.NotFound().Render(ctx.Request.Context(), ctx.Writer)
}

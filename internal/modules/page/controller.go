package page

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	ctx.HTML(http.StatusOK, "home.html", gin.H{
		"title": "HTML Sharer",
	})
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

	pages, total, err := c.service.GetPagesList(ctx.Request.Context(), page, pageSize)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Error",
			"error": "Failed to load pages",
		})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":       "Shared Pages",
		"pages":       pages,
		"currentPage": page,
		"totalPages":  totalPages,
		"total":       total,
		"hasNext":     page < int(totalPages),
		"hasPrev":     page > 1,
		"nextPage":    page + 1,
		"prevPage":    page - 1,
	})
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

	// Create page using service
	req := &PageCreate{HTMLContent: htmlContent}
	response, err := c.service.CreatePage(ctx.Request.Context(), req)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error creating page")
		return
	}

	if response.Error != "" {
		ctx.String(http.StatusBadRequest, response.Error)
		return
	}

	// Redirect to success page or return JSON
	if ctx.GetHeader("Accept") == "application/json" {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/?success="+strings.TrimPrefix(response.URL, "/shared/"))
	}
}

// CreateFromAPI handles API requests for creating pages
func (c *Controller) CreateFromAPI(ctx *gin.Context) {
	var req PageCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, PageResponse{Error: "Invalid JSON"})
		return
	}

	response, err := c.service.CreatePage(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, PageResponse{Error: "Internal server error"})
		return
	}

	if response.Error != "" {
		ctx.JSON(http.StatusBadRequest, *response)
		return
	}

	ctx.JSON(http.StatusOK, *response)
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
	ctx.HTML(http.StatusNotFound, "404.html", gin.H{
		"title": "Page Not Found",
	})
}

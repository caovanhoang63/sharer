package category

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sharer/views/components"
	"sharer/views/pages"
)

// Controller handles HTTP requests for category operations
type Controller struct {
	service Service
}

// NewController creates a new category controller
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// Index handles the category list page
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

	categoriesList, total, err := c.service.GetCategoriesList(ctx.Request.Context(), page, pageSize)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Convert CategoryList to CategoryData
	categoriesData := make([]*pages.CategoryData, len(categoriesList))
	for i, cat := range categoriesList {
		categoriesData[i] = &pages.CategoryData{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	hasNext := page < int(totalPages)
	hasPrev := page > 1

	ctx.Header("Content-Type", "text/html")
	pages.Categories(categoriesData, page, totalPages, total, hasNext, hasPrev).Render(ctx.Request.Context(), ctx.Writer)
}

// Create handles category creation form display
func (c *Controller) Create(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	pages.CategoryCreateForm().Render(ctx.Request.Context(), ctx.Writer)
}

// Store handles category creation from form submission
func (c *Controller) Store(ctx *gin.Context) {
	var req CategoryCreate
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	response, err := c.service.CreateCategory(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if response.Error != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
		return
	}

	// Return success response or redirect
	if ctx.GetHeader("HX-Request") == "true" {
		// For modal requests, return HTML that triggers modal close and page refresh
		ctx.Header("Content-Type", "text/html")
		ctx.Header("HX-Trigger", "closeModal")
		ctx.String(http.StatusOK, `<script>document.getElementById('create_category_modal').close(); window.location.reload();</script>`)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/categories")
	}
}

// Show handles displaying a single category
func (c *Controller) Show(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	_, err = c.service.GetCategoryByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	// For now, redirect to categories list since we don't have a show template
	ctx.Redirect(http.StatusSeeOther, "/categories")
}

// Edit handles category edit form display
func (c *Controller) Edit(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := c.service.GetCategoryByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	ctx.Header("Content-Type", "text/html")
	pages.CategoryEditForm(category.ID, category.Name, category.Description).Render(ctx.Request.Context(), ctx.Writer)
}

// EditModal handles category edit modal content display
func (c *Controller) EditModal(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := c.service.GetCategoryByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	ctx.Header("Content-Type", "text/html")
	components.CategoryEditModalContent(category.ID, category.Name, category.Description).Render(ctx.Request.Context(), ctx.Writer)
}

// Update handles category update from form submission
func (c *Controller) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req CategoryUpdate
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	response, err := c.service.UpdateCategory(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if response.Error != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
		return
	}

	// Return success response or redirect
	if ctx.GetHeader("HX-Request") == "true" {
		// For modal requests, return HTML that triggers modal close and page refresh
		ctx.Header("Content-Type", "text/html")
		ctx.Header("HX-Trigger", "closeModal")
		ctx.String(http.StatusOK, `<script>document.getElementById('edit_category_modal').close(); window.location.reload();</script>`)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/categories")
	}
}

// Delete handles category deletion
func (c *Controller) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = c.service.DeleteCategory(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Return success response or redirect
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	} else {
		ctx.Redirect(http.StatusSeeOther, "/categories")
	}
}

// GetAllForDropdown handles API requests for category dropdown data
func (c *Controller) GetAllForDropdown(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories(ctx.Request.Context())
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error loading categories")
		return
	}

	// Return HTML options for the dropdown
	html := `<option value="">Select a category...</option>`
	for _, cat := range categories {
		html += fmt.Sprintf(`<option value="%d">%s</option>`, cat.ID, cat.Name)
	}

	ctx.Data(http.StatusOK, "text/html", []byte(html))
}

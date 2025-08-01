package category

import (
	"gorm.io/gorm"
	"time"
)

// Category represents the database table for categories
type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CategoryCreate represents the data needed to create a new category
type CategoryCreate struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description,omitempty" form:"description"`
}

// CategoryUpdate represents the data that can be updated for a category
type CategoryUpdate struct {
	Name        *string `json:"name,omitempty" form:"name"`
	Description *string `json:"description,omitempty" form:"description"`
}

// CategoryList represents a simplified category for listing purposes
type CategoryList struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CategoryDetail represents detailed category information
type CategoryDetail struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryResponse represents the API response for category operations
type CategoryResponse struct {
	Category *CategoryDetail `json:"category,omitempty"`
	Error    string          `json:"error,omitempty"`
}

// TableName returns the table name for the Category model
func (Category) TableName() string {
	return "categories"
}

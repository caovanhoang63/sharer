package page

import (
	"gorm.io/gorm"
	"time"
)

// Page represents the database table for shared HTML pages
type Page struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	HTMLContent string         `gorm:"type:text;not null" json:"html_content"`
	Title       string         `gorm:"size:255" json:"title,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// PageCreate represents the data needed to create a new page
type PageCreate struct {
	HTMLContent string `json:"html_content" binding:"required"`
	Title       string `json:"title,omitempty"`
}

// PageUpdate represents the data that can be updated for a page
type PageUpdate struct {
	HTMLContent *string `json:"html_content,omitempty"`
	Title       *string `json:"title,omitempty"`
}

// PageList represents a simplified page for listing purposes
type PageList struct {
	ID        uint      `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// PageDetail represents detailed page information
type PageDetail struct {
	ID          uint      `json:"id"`
	Slug        string    `json:"slug"`
	HTMLContent string    `json:"html_content"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PageResponse represents the API response for page operations
type PageResponse struct {
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

// TableName returns the table name for the Page model
func (Page) TableName() string {
	return "shared_content"
}

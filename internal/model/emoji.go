package model

import (
	"time"

	"gorm.io/datatypes"
)

type Emoji struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	AuthorUUID       string         `gorm:"size:36;not null" json:"author_uuid"`
	Name             string         `gorm:"size:100;" json:"name"`
	Tags             datatypes.JSON `gorm:"type:json" json:"tags"`
	URL              string         `gorm:"size:255" json:"url"`
	View_count       int            `gorm:"default:0" json:"view_count"`
	Collection_count int            `gorm:"default:0" json:"collection_count"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

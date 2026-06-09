package model

import (
	"time"

	"gorm.io/datatypes"
)

type EmojiPack struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:100;" json:"name"`
	IconURL          string         `gorm:"size:255" json:"iconUrl"`
	View_count       int            `gorm:"default:0" json:"viewCount"`
	IsDefault        bool           `gorm:"default:false" json:"isDefault"`
	Tags             datatypes.JSON `gorm:"type:json" json:"tags"`
	Collection_count int            `gorm:"default:0" json:"collectionCount"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	AuthorUUID       string         `gorm:"size:36;not null" json:"authorUuid"`
}

type EmojiPack_Emoji struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EmojiPackID uint      `gorm:"index" json:"emojiPackId"`
	EmojiID     uint      `gorm:"index" json:"emojiId"`
	CreatedAt   time.Time `json:"createdAt"`
}

package models

import "time"

type Document struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Filename    string    `json:"filename"`
	URL         string    `json:"url"`
	ContentType string    `json:"content_type"`
	VehicleID   int       `json:"vehicle_id"`
	UploadedAt  time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
}

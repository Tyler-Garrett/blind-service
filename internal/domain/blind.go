package domain

import "time"

type Blind struct {
	Id   string     `json:"id"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type BlindsInRadiusRequest struct {
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"`
	RadiusMeters float64 `json:"radiusMeters" binding:"required,gt=0"`
}

type CreateBlindRequest struct {
	BlindID   string  `json:"blindId" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

type UpdateBlindRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

type BlindStats struct {
	TotalBlinds int     `json:"totalBlinds"`
	MinLatitude float64 `json:"minLatitude"`
	MaxLatitude float64 `json:"maxLatitude"`
	MinLongitude float64 `json:"minLongitude"`
	MaxLongitude float64 `json:"maxLongitude"`
}
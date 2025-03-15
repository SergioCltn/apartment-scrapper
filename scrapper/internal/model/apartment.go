package model

import "time"

type Apartment struct {
	ID                string    `json:"id"`
	Title             string    `json:"title,omitempty"`
	PropertyPrice     string    `json:"propertyPrice,omitempty"`
	PricePerSqm       string    `json:"pricePerSqm,omitempty"`
	CommunityFees     string    `json:"communityFees,omitempty"`
	Details           string    `json:"details,omitempty"`
	Location          string    `json:"location,omitempty"`
	Description       string    `json:"description,omitempty"`
	ApartmentStatus   string    `json:"apartmentStatus,omitempty"`
	BasicFeatures     string    `json:"basicFeatures,omitempty"`
	Building          string    `json:"building,omitempty"`
	EnergyCertificate string    `json:"energyCertificate,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
}

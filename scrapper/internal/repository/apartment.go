package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sergiocltn/apartment-scrapper/internal/model"
	"github.com/sergiocltn/apartment-scrapper/internal/provider"
)

type ApartmentRepository struct {
	db *sql.DB
}

func NewApartmentRepository(db *sql.DB) *ApartmentRepository {
	return &ApartmentRepository{db: db}
}

func (r *ApartmentRepository) Initialize() error {
	query := `
        CREATE TABLE IF NOT EXISTS apartments (
            id TEXT PRIMARY KEY,
            title TEXT,
            propertyPrice TEXT,
            pricePerSqm TEXT,
            communityFees TEXT,
            location TEXT,
            description TEXT,
            basicFeatures TEXT,
            building TEXT,
            energyCertificate TEXT,
            apartmentStatus TEXT,
            createdAt TEXT NOT NULL,
            updatedAt TEXT
        )
    `
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	provider.InfoLogger.Println("Initializing database connection")
	return nil
}

func (r *ApartmentRepository) Save(apartment model.Apartment) error {
	query := `
        INSERT INTO apartments (
            id, title, location, description, createdAt, updatedAt, propertyPrice, pricePerSqm, communityFees, apartmentStatus, basicFeatures, building, energyCertificate
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := r.db.Exec(query,
		apartment.ID,
		apartment.Title,
		apartment.Location,
		apartment.Description,
		apartment.CreatedAt.Format(time.RFC3339),
		nullTime(apartment.UpdatedAt),
		apartment.PropertyPrice,
		apartment.PricePerSqm,
		apartment.CommunityFees,
		apartment.ApartmentStatus,
		apartment.BasicFeatures,
		apartment.Building,
		apartment.EnergyCertificate,
	)

	if err != nil {
		return fmt.Errorf("failed to save apartment %s: %v", apartment.ID, err)
	}

	provider.InfoLogger.Printf("Saved apartment with id: %s", apartment.ID)
	return nil
}

func (r *ApartmentRepository) Exists(id string) (bool, error) {
	var exists int
	err := r.db.QueryRow("SELECT 1 FROM apartments WHERE id = ? LIMIT 1", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ApartmentRepository) FindByID(id string) (*model.Apartment, error) {
	var row struct {
		ID            string
		Title         sql.NullString
		PropertyPrice sql.NullString
		PricePerSqm   sql.NullString
		CommunityFees sql.NullString
		Details       sql.NullString
		Location      sql.NullString
		Description   sql.NullString
		CreatedAt     string
		UpdatedAt     sql.NullString
	}

	err := r.db.QueryRow("SELECT * FROM apartments WHERE id = ?", id).Scan(
		&row.ID,
		&row.Title,
		&row.Details,
		&row.PropertyPrice,
		&row.PricePerSqm,
		&row.CommunityFees,
		&row.Location,
		&row.Description,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse createdAt: %v", err)
	}

	var updatedAt time.Time
	if row.UpdatedAt.Valid {
		updatedAt, err = time.Parse(time.RFC3339, row.UpdatedAt.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updatedAt: %v", err)
		}
	}

	return &model.Apartment{
		ID:            row.ID,
		Title:         row.Title.String,
		PropertyPrice: row.PropertyPrice.String,
		PricePerSqm:   row.PricePerSqm.String,
		CommunityFees: row.CommunityFees.String,
		Details:       row.Details.String,
		Location:      row.Location.String,
		Description:   row.Description.String,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func (r *ApartmentRepository) Update(id string, data model.Apartment) error {
	fields := []string{}
	values := []any{}

	if data.Title != "" {
		fields = append(fields, "title = ?")
		values = append(values, data.Title)
	}
	if data.PricePerSqm != "" {
		fields = append(fields, "pricePerSqm = ?")
		values = append(values, data.PricePerSqm)
	}
	if data.PropertyPrice != "" {
		fields = append(fields, "propertyPrice = ?")
		values = append(values, data.PropertyPrice)
	}
	if data.CommunityFees != "" {
		fields = append(fields, "communityFees = ?")
		values = append(values, data.CommunityFees)
	}
	if data.Location != "" {
		fields = append(fields, "location = ?")
		values = append(values, data.Location)
	}
	if data.Description != "" {
		fields = append(fields, "description = ?")
		values = append(values, data.Description)
	}
	if data.Details != "" {
		fields = append(fields, "details = ?")
		values = append(values, data.Details)
	}

	fields = append(fields, "updatedAt = ?")
	values = append(values, time.Now().Format(time.RFC3339))
	values = append(values, id)

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE apartments SET %s WHERE id = ?", strings.Join(fields, ", "))
	_, err := r.db.Exec(query, values...)
	if err != nil {
		return err
	}
	provider.InfoLogger.Printf("Updated apartment with id: %s", id)
	return nil
}

func (r *ApartmentRepository) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	provider.InfoLogger.Println("Closing database connection")
	return nil
}

func nullTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t.Format(time.RFC3339)
}

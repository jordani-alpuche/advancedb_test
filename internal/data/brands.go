package data

import (
	"context"
	"database/sql"
	"github/jordani-alpuche/test2/internal/validator"
	"time"
)

// BrandData holds the information of a journal entry
type BrandData struct {
	ID        int64     `json:"id"`
	BrandName	 string    `json:"brand_name"`
	BrandDescription   string    `json:"brand_description"`
	BrandCreatedAt time.Time `json:"brand_created_at"`
	BrandUpdatedAt time.Time `json:"brand_updated_at"`
}

// BrandDataModel represents the database model for journal entries
type BrandDataModel struct {
	DB *sql.DB
}

func ValidateBrands(v *validator.Validator, brand *BrandData) {
	v.Check(validator.NotBlank(brand.BrandName), "BrandName", "Brand Name must be provided")
	v.Check(validator.MaxLength(brand.BrandName, 50), "BrandName", "Brand Name must not be more than 50 bytes long")
	v.Check(validator.NotBlank(brand.BrandDescription), "BrandDescription", "Brand Description must be provided")
	v.Check(validator.MaxLength(brand.BrandDescription, 1500), "BrandDescription", "Brand Description must not be more than 1500 bytes long")
}


// This method is used to count the number of brands in the database
// It returns a slice of BrandData and an error if any
// It uses a context with a timeout of 3 seconds for the query
func (m *BrandDataModel) CountAllBrands() (int, error) {
	var count int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM brand`

	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Select method fetches all brand entries from the database
func (m *BrandDataModel) GET(id int) ([]BrandData, error) {
	var query string
	var rows *sql.Rows
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if id == 0 {
		query = `SELECT id, brand_name,brand_description,brand_created_at FROM brand ORDER BY brand_created_at DESC`
		rows, err = m.DB.QueryContext(ctx, query)
	} else {
		query = `SELECT id, brand_name,brand_description,brand_created_at FROM brand WHERE id = $1 ORDER BY brand_created_at DESC`
		rows, err = m.DB.QueryContext(ctx, query, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []BrandData

	for rows.Next() {
		var brand BrandData		
		var BrandCreatedAt sql.NullTime

		err := rows.Scan(
			&brand.ID,
			&brand.BrandName,
			&brand.BrandDescription,			
			&BrandCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if BrandCreatedAt.Valid {
			brand.BrandCreatedAt = BrandCreatedAt.Time
		}

		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

// Create method inserts a new journal entry into the database
func (m *BrandDataModel) POST(brand *BrandData) error {
	// SQL query to insert a new brand entry

	// fmt.Printf("passed data: %v",brand)
	query := `
		INSERT INTO brand (brand_name,brand_description,brand_created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	// Create a context with a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	
	// Insert the brand into the database and get the generated ID
	err := m.DB.QueryRowContext(ctx, query, brand.BrandName, brand.BrandDescription, time.Now()).Scan(&brand.ID)
	if err != nil {
		return err
	}

	// Return nil if the creation is successful
	return nil
}

// PUT updates an existing brand record
func (m *BrandDataModel) PUT(id int,brand *BrandData) error {
	query := `
		UPDATE brand
		SET
			brand_name = $1,
			brand_description = $2,
			brand_updated_at = $3
		WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the update
	_, err := m.DB.ExecContext(ctx, query,
		brand.BrandName,
		brand.BrandDescription,
		time.Now(),     // Update the timestamp
		id,       // Use ID in the WHERE clause
	)

	return err
}


// Delete method removes a journal entry by its ID from the database
func (m *BrandDataModel) DELETE(id int) error {
	// SQL query to delete a journal entry by its ID
	query := `DELETE FROM brand WHERE id = $1`

	// Create a context with a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the delete query
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Return nil if the deletion is successful
	return nil
}

package data

import (
	"context"
	"database/sql"
	"github/jordani-alpuche/test2/internal/validator"
	"time"
)

// CategoryData holds the information of a journal entry
type CategoryData struct {
	ID        int64     `json:"id"`
	CategoryName     string    `json:"title"`
    CategoryDescription   string    `json:"category_description"`
	CategoryCode   string    `json:"content"`	
	CategoryCreatedAt time.Time `json:"category_created_at"`
	CategoryUpdatedAt time.Time `json:"category_updated_at"`
}

// CategoryDataModel represents the database model for journal entries
type CategoryDataModel struct {
	DB *sql.DB
}

func ValidateCategory(v *validator.Validator, category *CategoryData) {
	v.Check(validator.NotBlank(category.CategoryName), "CategoryName", "Category Name must be provided")
	v.Check(validator.MaxLength(category.CategoryName, 50), "CategoryName", "Category Name must not be more than 50 bytes long")
    v.Check(validator.NotBlank(category.CategoryDescription), "CategoryDescription", "Category Description must be provided")
	v.Check(validator.MaxLength(category.CategoryDescription, 1500), "CategoryDescription", "Category Description must not be more than 1500 bytes long")
	v.Check(validator.NotBlank(category.CategoryCode), "CategoryCode", "Category Code must be provided")
	v.Check(validator.MaxLength(category.CategoryCode, 150), "CategoryCode", "Category Code must not be more than 150 bytes long")
}



func (m *CategoryDataModel) CountAllCategories() (int, error) {
	var count int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM category`

	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Select method fetches all category entries from the database
func (m *CategoryDataModel) GET(id int) ([]CategoryData, error) {
	var query string
	var rows *sql.Rows
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if id == 0 {
		query = `SELECT id, category_name, category_description, category_code,category_created_at FROM category ORDER BY category_created_at DESC`
		rows, err = m.DB.QueryContext(ctx, query)
	} else {
		query = `SELECT id, category_name, category_description, category_code,category_created_at FROM category WHERE id = $1 ORDER BY category_created_at DESC`
		rows, err = m.DB.QueryContext(ctx, query, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryData

	for rows.Next() {
		var category CategoryData		
		var CategoryCreatedAt sql.NullTime

		err := rows.Scan(
			&category.ID,
			&category.CategoryName,
			&category.CategoryDescription,
            &category.CategoryCode,
			
			&CategoryCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if CategoryCreatedAt.Valid {
			category.CategoryCreatedAt = CategoryCreatedAt.Time
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// Create method inserts a new journal entry into the database
func (m *CategoryDataModel) POST(category *CategoryData) error {
	// SQL query to insert a new category entry

	// fmt.Printf("passed data: %v",category)
	query := `
		INSERT INTO category (category_name, category_description, category_code,category_created_at)
		VALUES ($1, $2, $3,$4)
		RETURNING id
	`

	// Create a context with a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	
	// Insert the category into the database and get the generated ID
	err := m.DB.QueryRowContext(ctx, query, category.CategoryName, category.CategoryDescription,category.CategoryCode ,time.Now()).Scan(&category.ID)
	if err != nil {
		return err
	}

	// Return nil if the creation is successful
	return nil
}

func (m *CategoryDataModel) PUT(id int,category *CategoryData) error {
	query := `
		UPDATE category
		SET
			category_name = $1,
			category_description = $2,
			category_code = $3,
			category_updated_at = $4
		WHERE id = $5
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query,
		category.CategoryName,
		category.CategoryDescription,
		category.CategoryCode,
		time.Now(),          // updating `category_updated_at`
		id,       // WHERE clause
	)

	return err
}

// Delete method removes a journal entry by its ID from the database
func (m *CategoryDataModel) DELETE(id int) error {
	// SQL query to delete a journal entry by its ID
	query := `DELETE FROM category WHERE id = $1`

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

package data

import (
	"context"
	"database/sql"
	"github/jordani-alpuche/test1/internal/validator"
	"time"
)

// UsersData holds the information of a journal entry
type UsersData struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`	
	CreatedAt time.Time `json:"created_at"`
	Image_Url string   `json:"image_url"`
}

// UsersDataModel represents the database model for journal entries
type UsersDataModel struct {
	DB *sql.DB
}

func ValidateUsers(v *validator.Validator, users *UsersData) {
	v.Check(validator.NotBlank(users.Title), "title", "must be provided")
	v.Check(validator.MaxLength(users.Title, 50), "title", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(users.Content), "content", "must be provided")
	v.Check(validator.MaxLength(users.Content, 5000), "content", "must not be more than 500 bytes long")
}

// Select method fetches all journal entries from the database
func (m *UsersDataModel) Select(id int) ([]UsersData, error) {
	var query string
	var rows *sql.Rows
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if id == 0 {
		query = `SELECT id, title, content, created_at FROM journals ORDER BY created_at DESC`
		rows, err = m.DB.QueryContext(ctx, query)
	} else {
		query = `SELECT id, title, content, created_at FROM journals WHERE id = $1 ORDER BY created_at DESC`
		rows, err = m.DB.QueryContext(ctx, query, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UsersData

	for rows.Next() {
		var user UsersData		
		var createdAt sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.Title,
			&user.Content,
			
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Create method inserts a new journal entry into the database
func (m *UsersDataModel) Insert(journal *UsersData) error {
	// SQL query to insert a new journal entry

	// fmt.Printf("passed data: %v",journal)
	query := `
		INSERT INTO journals (title, content, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	// Create a context with a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	
	// Insert the journal into the database and get the generated ID
	err := m.DB.QueryRowContext(ctx, query, journal.Title, journal.Content, time.Now()).Scan(&journal.ID)
	if err != nil {
		return err
	}

	// Return nil if the creation is successful
	return nil
}

// Delete method removes a journal entry by its ID from the database
func (m *UsersDataModel) Delete(id int64) error {
	// SQL query to delete a journal entry by its ID
	query := `DELETE FROM journals WHERE id = $1`

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

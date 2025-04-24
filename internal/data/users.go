package data

import (
	"context"
	"database/sql"
	"github/jordani-alpuche/test2/internal/validator"
	"time"
)

// UsersData holds the information of a journal entry
type UsersData struct {
	ID        int64     `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Username	 string    `json:"username"`
	Password 	 string    `json:"password"`
	Email      string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role	   string    `json:"role"`
	Status     string    `json:"status"`

}

// UsersDataModel represents the database model for journal entries
type UsersDataModel struct {
	DB *sql.DB
}

func ValidateUsers(v *validator.Validator, users *UsersData) {
	v.Check(validator.NotBlank(users.FirstName), "FirstName", "First Name must be provided")
	v.Check(validator.MaxLength(users.FirstName, 150), "FirstName", "must not be more than 150 bytes long")
	v.Check(validator.NotBlank(users.LastName), "LastName", "Last Name must be provided")
	v.Check(validator.MaxLength(users.LastName, 150), "LastName", "must not be more than 150 bytes long")
	v.Check(validator.NotBlank(users.Username), "Username", "Username must be provided")
	v.Check(validator.MaxLength(users.Username, 150), "Username", "must not be more than 150 bytes long")
	v.Check(validator.NotBlank(users.Password), "Password", "Password must be provided")
	v.Check(len(users.Password) >= 8, "Password", "Password must be at least 8 characters long")
	v.Check(validator.NotBlank(users.Email), "Email", "Email must be provided")
	v.Check(validator.IsValidEmail(users.Email), "Email", "Must be a valid Email address")
	v.Check(validator.NotBlank(users.PhoneNumber), "PhoneNumber", "Phone Number must be provided")
	v.Check(validator.NotBlank(users.Role), "Role", "Role must be provided")
	v.Check(validator.NotBlank(users.Status), "Status", "Status must be provided")
}
func ValidateLogin(v *validator.Validator, users *UsersData) {
	v.Check(validator.NotBlank(users.Username), "Username", "Username must be provided")
	v.Check(validator.NotBlank(users.Password), "Password", "Password must be provided")
}

func ValidateSignup(v *validator.Validator, users *UsersData) {
	v.Check(validator.NotBlank(users.FirstName), "FirstName", "First Name must be provided")
	v.Check(validator.MaxLength(users.FirstName, 150), "FirstName", "must not be more than 150 bytes long")
	v.Check(validator.NotBlank(users.LastName), "LastName", "Last Name must be provided")
	v.Check(validator.MaxLength(users.LastName, 150), "LastName", "must not be more than 150 bytes long")
	v.Check(validator.NotBlank(users.Email), "Email", "Email must be provided")
	v.Check(validator.IsValidEmail(users.Email), "Email", "Must be a valid Email address")
	v.Check(validator.NotBlank(users.Password), "Password", "Password must be provided")
}

func (m *UsersDataModel) CountAllActiveUsers() (int, error) {
	var count int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM users WHERE status = 'active'`

	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByUsername checks if a user with the given username exists
// func (m *UsersDataModel) FindByUsername(username string) ([]UsersData, error) {
func (m *UsersDataModel) FindByUsername(username string,formType string) ([]UsersData, error){
	var query string
	var rows *sql.Rows
	var err error

	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if formType == "users" {

	query = `
		SELECT id, first_name, last_name, username, password, email, phone_number, role, status
		FROM users
		WHERE username = $1
		LIMIT 1
	`
	rows, err = m.DB.QueryContext(ctx, query, username)
	} else if formType == "login" {	
		query = `
		SELECT id, first_name, last_name, username, password, email, phone_number, role, status
		FROM users
		WHERE username = $1 and status = 'active'
		LIMIT 1
	`
	rows, err = m.DB.QueryContext(ctx, query, username)
	}


	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UsersData

	for rows.Next() {
		var user UsersData		

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.PhoneNumber,
			&user.Role,
			&user.Status,		
			
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}


// Select method fetches all journal entries from the database
func (m *UsersDataModel) GET(id int) ([]UsersData, error) {
	var query string
	var rows *sql.Rows
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if id == 0 {
		query = `SELECT id, first_name, last_name, username,password,email,phone_number,role,status FROM users ORDER BY first_name DESC`
		rows, err = m.DB.QueryContext(ctx, query)
	} else {
		query = `SELECT id, first_name, last_name, username,password,email,phone_number,role,status FROM users WHERE id = $1 ORDER BY first_name DESC`
		rows, err = m.DB.QueryContext(ctx, query, id)
	}


	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UsersData

	for rows.Next() {
		var user UsersData		

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.PhoneNumber,
			&user.Role,
			&user.Status,		
			
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Create method inserts a new user entry into the database
func (m *UsersDataModel) POST(user *UsersData) error {
	// SQL query to insert a new user entry

	// fmt.Printf("passed data: %v",user)
	query := `
		INSERT INTO users (first_name, last_name, username,password,email,phone_number,role,status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	// Create a context with a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	
	// Insert the user into the database and get the generated ID
	err := m.DB.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username,
		user.Password,user.Email,user.PhoneNumber,user.Role,user.Status).Scan(&user.ID)
	if err != nil {
		return err
	}

	// Return nil if the creation is successful
	return nil
}

// Update method updates an existing user entry in the database
func (m *UsersDataModel) PUT(id int,user *UsersData) error {
	query := `
		UPDATE users
		SET 
			first_name = $1,
			last_name = $2,
			username = $3,
			password = $4,
			email = $5,
			phone_number = $6,
			role = $7,
			status = $8		
		WHERE id = $9
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Email,
		user.PhoneNumber,
		user.Role,
		user.Status,		
		id,                 // user ID
	)

	return err
}

// Delete method removes a journal entry by its ID from the database
func (m *UsersDataModel) DELETE(id int) error {
	// SQL query to delete a journal entry by its ID
	query := `DELETE FROM users WHERE id = $1`

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

package data

import (
	"context"
	"database/sql"
	"github/jordani-alpuche/test2/internal/validator"
	"time"
)

// ProductData holds the information of a feedback entry
type ProductData struct {
	ID        int64     `json:"id"`
	ProductName     string    `json:"product_name"`
	ProductDescription   string    `json:"product_description"`
	ProductPrice  float64    `json:"product_price"`
	ProductCategoryID  int    `json:"product_category_id"`
	ProductBrandID int `json:"product_brand_id"`
	ProductQTY int   `json:"product_qty"`
	ProductStatus string   `json:"product_status"`
	ProductCreateTime time.Time `json:"product_create_time"`
	ProductUpdateTime time.Time `json:"product_update_time"`
	ProductPurchasedFrom string   `json:"product_purchased_from"`
	ProductTag string `json:"product_tag"`

	// --- Add fields for related names ---
    CategoryName       string `db:"category_name"` 
    BrandName          string `db:"brand_name"`    
}

// ProductDataModel represents the database model for feedback entries
type ProductDataModel struct {
	DB *sql.DB
}

func ValidateProduct(v *validator.Validator, product *ProductData) {
	v.Check(validator.NotBlank(product.ProductName), "ProductName", "Product Name must be provided")
	v.Check(validator.MaxLength(product.ProductName, 50), "ProductName", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(product.ProductDescription), "ProductDescription", "Product Description must be provided")
	v.Check(validator.MaxLength(product.ProductDescription, 150), "ProductDescription", "must not be more than 150 bytes long")
	v.Check(validator.NotBlankFloat(product.ProductPrice), "ProductPrice", "Product Price must be provided")
	v.Check(product.ProductPrice > 0, "ProductPrice", "Product Price must be greater than 0")
	v.Check(validator.NotBlankInt(product.ProductCategoryID), "ProductCategoryID", "Product Category must be provided")
	v.Check(validator.NotBlankInt(product.ProductBrandID), "ProductBrandID", "Product Brand must be provided")
	v.Check(validator.NotBlankInt(product.ProductQTY), "ProductQTY", "Product QTY must be provided")
	v.Check(validator.NotBlank(product.ProductStatus), "ProductStatus", "Product Status must be provided")
	v.Check(validator.MaxLength(product.ProductStatus, 50), "ProductStatus", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(product.ProductPurchasedFrom), "ProductPurchasedFrom", "Product Purchase From must be provided")
	v.Check(validator.MaxLength(product.ProductPurchasedFrom, 150), "ProductPurchasedFrom", "must not be more than 150 bytes long")
}




func (m *ProductDataModel) CountAllProducts() (int, error) {
	var count int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM products`

	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Select method fetches all product entries from the database
func (m *ProductDataModel) GET(id int) ([]ProductData, error) {
	var query string
	var rows *sql.Rows
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

    const baseQuery = `
        SELECT
            p.id, p.product_name, p.product_description, p.product_price,
            p.product_category_id, p.product_brand_id, p.product_qty,
            p.product_status, p.product_purchased_from, p.product_create_time,
            COALESCE(c.category_name, '') AS category_name, -- Use COALESCE for LEFT JOIN robustness
            COALESCE(b.brand_name, '') AS brand_name
        FROM
            products p
        LEFT JOIN category c ON p.product_category_id = c.id -- LEFT JOIN is safer if category/brand might be missing
        LEFT JOIN brand b ON p.product_brand_id = b.id`

    if id == 0 {
        // Query for all products
        query = baseQuery + ` ORDER BY p.product_create_time DESC`
        rows, err = m.DB.QueryContext(ctx, query)
    } else {
        // Query for a single product by ID
        query = baseQuery + ` WHERE p.id=$1 ORDER BY p.product_create_time DESC` // Order might be less relevant for single item
        rows, err = m.DB.QueryContext(ctx, query, id)
    }
    // --- End Updated Queries ---

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductData

	for rows.Next() {
		var product ProductData	
		var ProductCreateTime sql.NullTime

		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.ProductDescription,
			&product.ProductPrice,
			&product.ProductCategoryID,
			&product.ProductBrandID,
			&product.ProductQTY,
			&product.ProductStatus,
			&product.ProductPurchasedFrom,
			&product.ProductCreateTime,
			&product.CategoryName,         // Scan the category name
            &product.BrandName,            // Scan the brand name
		)
		if err != nil {
			return nil, err
		}

		if ProductCreateTime.Valid {
			product.ProductCreateTime = ProductCreateTime.Time
		}


		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Create method inserts a new products entry into the database
func (m *ProductDataModel) POST(products *ProductData) error {
	// SQL query to insert a new products entry

	query := `
		INSERT INTO products(product_name,product_description,product_price,product_category_id,product_brand_id,product_qty,product_status,product_purchased_from,product_create_time,product_tag)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	// Create a context with a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	

	// Insert the products into the database and get the generated ID
	err := m.DB.QueryRowContext(ctx, query, products.ProductName, products.ProductDescription,products.ProductPrice,products.ProductCategoryID,products.ProductBrandID,products.ProductQTY,products.ProductStatus,products.ProductPurchasedFrom ,time.Now(),products.ProductTag).Scan(&products.ID)
	if err != nil {
		return err
	}

	// Return nil if the creation is successful
	return nil
}

// PUT updates an existing product in the database
func (m *ProductDataModel) PUT(id int,product *ProductData) error {
	query := `
		UPDATE products
		SET 
			product_name = $1,
			product_description = $2,
			product_price = $3,
			product_category_id = $4,
			product_brand_id = $5,
			product_qty = $6,
			product_status = $7,
			product_purchased_from = $8,
			product_update_time = $9
		WHERE id = $10
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query,
		product.ProductName,
		product.ProductDescription,
		product.ProductPrice,
		product.ProductCategoryID,
		product.ProductBrandID,
		product.ProductQTY,
		product.ProductStatus,
		product.ProductPurchasedFrom,
		time.Now(),          // update timestamp
		id,                 // product ID
	)

	return err
}

// Delete method removes a products entry by its ID from the database
func (m *ProductDataModel) DELETE(id int) error {
	// SQL query to delete a products entry by its ID
	query := `DELETE FROM products WHERE id = $1`

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

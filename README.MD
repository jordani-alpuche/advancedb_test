# Inventory Web App

A Go-based web application that manages **products**, **categories**, **brands**, and **users** with a PostgreSQL backend. Built using Go standard packages and HTML templating, this app serves as a lightweight inventory or catalog management system with a clear MVC separation.

---

# Note

There is two branch namely main and test2, main correspond to test1 and test2 correspond to test2

---

## 🚀 Features

- Full CRUD functionality for:
  - Products
  - Categories
  - Brands
- User management and role-based access
- PostgreSQL-backed storage with SQL migrations
- Clean HTML templating with `html/template`
- Logging with `slog`
- Organized file structure following Go best practices

---

## 🧱 Project Structure

```
cmd                         # Main application entry point
└── web                     # Web server layer and HTTP-related handlers
    ├── handlers_brands.go     # HTTP handlers for brand-related routes
    ├── handlers_category.go   # HTTP handlers for category-related routes
    ├── handlers_home.go       # Handler for homepage/dashboard
    ├── handlers_product.go    # HTTP handlers for product-related routes
    ├── main.go                # Entry point of the web server application
    ├── middleware.go          # Middleware functions (e.g., logging, auth)
    ├── render.go              # Template rendering helpers
    ├── routes.go              # Route definitions and route setup
    ├── server.go              # HTTP server configuration and startup logic
    ├── template_data.go       # Structs for passing data to templates
    └── templates.go           # Template parsing and loading utilities
internal                    # Application-specific internal packages
├── data                    # Business logic and database access layer
│   ├── brands.go               # Brand model and data access functions
│   ├── category.go             # Category model and data access functions
│   ├── products.go             # Product model and data access functions
│   └── users.go                # User model and data access functions
└── validator              # Custom validation logic
    └── validator.go           # Custom input validation utilities
migrations                 # Database schema migrations (SQL scripts)
├── 000001_create_products_table.down.sql   # Rollback: drop products table
├── 000001_create_products_table.up.sql     # Migration: create products table
├── 000002_create_category_table.down.sql   # Rollback: drop category table
├── 000002_create_category_table.up.sql     # Migration: create category table
├── 000003_create_brand_table.down.sql      # Rollback: drop brand table
├── 000003_create_brand_table.up.sql        # Migration: create brand table
├── 000004_create_users_table.down.sql      # Rollback: drop users table
└── 000004_create_users_table.up.sql        # Migration: create users table
ui                          # Frontend resources
├── html                    # HTML templates for rendering views
│   ├── addbrand.tmpl           # Template for adding a new brand
│   ├── addcategory.tmpl        # Template for adding a new category
│   ├── addproduct.tmpl         # Template for adding a new product
│   ├── adduser.tmpl            # Template for adding a new user
│   ├── brand-details.tmpl      # Template for viewing brand details
│   ├── brandlist.tmpl          # Template for listing brands
│   ├── category-details.tmpl   # Template for viewing category details
│   ├── categorylist.tmpl       # Template for listing categories
│   ├── editbrand.tmpl          # Template for editing brand
│   ├── editcategory.tmpl       # Template for editing category
│   ├── editproduct.tmpl        # Template for editing product
│   ├── edituser.tmpl           # Template for editing user
│   ├── error-404.tmpl          # Template for 404 error page
│   ├── error-500.tmpl          # Template for 500 error page
│   ├── index.tmpl              # Template for homepage
│   ├── layout.tmpl             # Base layout template used by others
│   ├── newuseredit.tmpl        # Template for new user edit form
│   ├── newuser.tmpl            # Template for new user creation
│   ├── product-details.tmpl    # Template for viewing product details
│   ├── productlist.tmpl        # Template for listing products
│   └── userlists.tmpl          # Template for listing users
└── static                  # Static assets (served as-is)
    ├── css                    # CSS stylesheets
    ├── js                     # JavaScript files
    ├── images                 # Image assets
    └── plugins                # Third-party JS/CSS plugins

```

## 🛠️ Technologies Used

- **Go** 1.23+
- **PostgreSQL**
- **HTML Templates**
- **SQL Migrations** (manual)
- **Standard Logging** (`slog`)

---

## 🗃️ Database Schema (via SQL Migrations)

The schema includes the following tables:

- `products`
- `category`
- `brand`
- `users`

Each table is created using plain SQL migrations located in the `/migrations` folder. Migrations are versioned and include `up` and `down` scripts.

---

## 🧪 Getting Started

### 1. Clone the Repo

```bash
git clone https://github.com/jordani-alpuche/advancedb_test.git
cd test1/tapir

createdb tapir_db
psql -d tapir_db -f migrations/000001_create_products_table.up.sql
psql -d tapir_db -f migrations/000002_create_category_table.up.sql
psql -d tapir_db -f migrations/000003_create_brand_table.up.sql
psql -d tapir_db -f migrations/000004_create_users_table.up.sql

make run

🌐 Endpoints (Examples)
HTTP Method	Route	Description
GET	/products	List all products
GET	/categories	List all categories
GET	/brands	List all brands
POST	/products/new	Create a new product
...	...	(Additional CRUD routes)
All rendered using Go’s built-in HTML templating engine.

🙌 Author
Jordaní Alpuche
Project module: github/jordani-alpuche/advancedb_test

📝 License
This project is licensed under the MIT License.
```

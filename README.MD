# Tapir Web App

A Go-based web application that manages **products**, **categories**, **brands**, and **users** with a PostgreSQL backend. Built using Go standard packages and HTML templating, this app serves as a lightweight inventory or catalog management system with a clear MVC separation.

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
tapir/
├── cmd/web/ # Web server entry point and handlers
│ ├── main.go # Main application entry
│ ├── handlers_*.go # Route handlers (products, categories, brands, home)
│ ├── routes.go # HTTP routes
│ ├── templates.go # HTML template loading
│ └── middleware.go # HTTP middleware
├── internal/
│ ├── data/ # Data models and DB logic
│ │ ├── products.go
│ │ ├── category.go
│ │ ├── brands.go
│ │ └── users.go
│ └── validator/ # Input validation logic
│ └── validator.go
├── migrations/ # SQL migration files
├── go.mod / go.sum # Go module definitions
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
git clone https://github.com/jordani-alpuche/test1.git
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
Project module: github/jordani-alpuche/test1

📝 License
This project is licensed under the MIT License.
```

package repository

import (
	"database/sql"

	"golang-crud-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT id, name, description, price, created_at, updated_at
	          FROM products ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT id, name, description, price, created_at, updated_at
	          FROM products WHERE id = $1`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(req models.CreateProductRequest) (*models.Product, error) {
	query := `INSERT INTO products (name, description, price)
	          VALUES ($1, $2, $3)
	          RETURNING id, name, description, price, created_at, updated_at`

	var p models.Product
	err := r.db.QueryRow(query, req.Name, req.Description, req.Price).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(id int, req models.UpdateProductRequest) (*models.Product, error) {
	query := `UPDATE products
	          SET name = $1, description = $2, price = $3, updated_at = NOW()
	          WHERE id = $4
	          RETURNING id, name, description, price, created_at, updated_at`

	var p models.Product
	err := r.db.QueryRow(query, req.Name, req.Description, req.Price, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Delete(id int) (bool, error) {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

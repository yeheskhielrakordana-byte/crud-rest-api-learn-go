package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	_ "github.com/lib/pq"
)

var db *sql.DB

type product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func connectDB() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "root")
	dbname := getEnv("DB_NAME", "golang_crud")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal membuka koneksi database:", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}
	return conn
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func toJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func args(req mcp.CallToolRequest) map[string]any {
	if m, ok := req.Params.Arguments.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func getInt(m map[string]any, key string) (int64, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key %q tidak ditemukan", key)
	}
	switch n := v.(type) {
	case json.Number:
		return n.Int64()
	case float64:
		return int64(n), nil
	}
	return 0, fmt.Errorf("key %q bukan angka", key)
}

func getFloat(m map[string]any, key string) float64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case json.Number:
		f, _ := n.Float64()
		return f
	case float64:
		return n
	}
	return 0
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// tool handlers

func handleGetAllProducts(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rows, err := db.Query(`SELECT id, name, description, price, created_at, updated_at FROM products ORDER BY id`)
	if err != nil {
		return mcp.NewToolResultText("Error: " + err.Error()), nil
	}
	defer rows.Close()

	var products []product
	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return mcp.NewToolResultText("Error scan: " + err.Error()), nil
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		return mcp.NewToolResultText("Tidak ada produk ditemukan."), nil
	}
	return mcp.NewToolResultText(toJSON(products)), nil
}

func handleGetProductByID(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args(req)
	id, err := getInt(a, "id")
	if err != nil {
		return mcp.NewToolResultText("Error: id tidak valid"), nil
	}

	var p product
	err = db.QueryRow(
		`SELECT id, name, description, price, created_at, updated_at FROM products WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return mcp.NewToolResultText(fmt.Sprintf("Produk dengan id %d tidak ditemukan.", id)), nil
	}
	if err != nil {
		return mcp.NewToolResultText("Error: " + err.Error()), nil
	}
	return mcp.NewToolResultText(toJSON(p)), nil
}

func handleCreateProduct(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args(req)
	name := getString(a, "name")
	description := getString(a, "description")
	price := getFloat(a, "price")

	if name == "" {
		return mcp.NewToolResultText("Error: name wajib diisi"), nil
	}

	var p product
	err := db.QueryRow(
		`INSERT INTO products (name, description, price) VALUES ($1, $2, $3)
		 RETURNING id, name, description, price, created_at, updated_at`,
		name, description, price,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return mcp.NewToolResultText("Error: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Produk berhasil dibuat:\n" + toJSON(p)), nil
}

func handleUpdateProduct(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args(req)
	id, err := getInt(a, "id")
	if err != nil {
		return mcp.NewToolResultText("Error: id tidak valid"), nil
	}
	name := getString(a, "name")
	description := getString(a, "description")
	price := getFloat(a, "price")

	if name == "" {
		return mcp.NewToolResultText("Error: name wajib diisi"), nil
	}

	var p product
	err = db.QueryRow(
		`UPDATE products SET name=$1, description=$2, price=$3, updated_at=NOW()
		 WHERE id=$4
		 RETURNING id, name, description, price, created_at, updated_at`,
		name, description, price, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return mcp.NewToolResultText(fmt.Sprintf("Produk dengan id %d tidak ditemukan.", id)), nil
	}
	if err != nil {
		return mcp.NewToolResultText("Error: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Produk berhasil diupdate:\n" + toJSON(p)), nil
}

func handleDeleteProduct(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args(req)
	id, err := getInt(a, "id")
	if err != nil {
		return mcp.NewToolResultText("Error: id tidak valid"), nil
	}

	result, err := db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return mcp.NewToolResultText("Error: " + err.Error()), nil
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Produk dengan id %d tidak ditemukan.", id)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Produk dengan id %d berhasil dihapus.", id)), nil
}

func main() {
	db = connectDB()
	defer db.Close()

	s := server.NewMCPServer("golang-crud-api-mcp", "1.0.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("get_all_products",
			mcp.WithDescription("Ambil semua data produk dari database"),
		),
		handleGetAllProducts,
	)

	s.AddTool(
		mcp.NewTool("get_product_by_id",
			mcp.WithDescription("Ambil satu produk berdasarkan ID"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("ID produk")),
		),
		handleGetProductByID,
	)

	s.AddTool(
		mcp.NewTool("create_product",
			mcp.WithDescription("Buat produk baru"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Nama produk")),
			mcp.WithString("description", mcp.Description("Deskripsi produk")),
			mcp.WithNumber("price", mcp.Description("Harga produk")),
		),
		handleCreateProduct,
	)

	s.AddTool(
		mcp.NewTool("update_product",
			mcp.WithDescription("Update produk berdasarkan ID"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("ID produk")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Nama produk")),
			mcp.WithString("description", mcp.Description("Deskripsi produk")),
			mcp.WithNumber("price", mcp.Description("Harga produk")),
		),
		handleUpdateProduct,
	)

	s.AddTool(
		mcp.NewTool("delete_product",
			mcp.WithDescription("Hapus produk berdasarkan ID"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("ID produk")),
		),
		handleDeleteProduct,
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

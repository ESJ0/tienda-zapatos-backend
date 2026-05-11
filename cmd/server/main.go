package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ESJ0/tienda-zapatos-backend/internal/config"
	"github.com/ESJ0/tienda-zapatos-backend/internal/db"
	"github.com/ESJ0/tienda-zapatos-backend/internal/handlers"
	"github.com/ESJ0/tienda-zapatos-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	pool, err := db.New(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to database")

	// Handlers
	authH := &handlers.AuthHandler{DB: pool, Cfg: cfg}
	productH := &handlers.ProductHandler{DB: pool}
	categoryH := &handlers.CategoryHandler{DB: pool}
	supplierH := &handlers.SupplierHandler{DB: pool}
	customerH := &handlers.CustomerHandler{DB: pool}
	employeeH := &handlers.EmployeeHandler{DB: pool}
	saleH := &handlers.SaleHandler{DB: pool}
	reportH := &handlers.ReportHandler{DB: pool}

	// Router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"zapatos-api"}`))
	})

	r.Route("/api", func(r chi.Router) {
		// Público
		r.Post("/auth/login", authH.Login)

		// Protegido
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg))

			r.Post("/auth/logout", authH.Logout)

			// Products
			r.Get("/products", productH.List)
			r.Post("/products", productH.Create)
			r.Get("/products/{id}", productH.Get)
			r.Put("/products/{id}", productH.Update)
			r.Delete("/products/{id}", productH.Delete)

			// Categories
			r.Get("/categories", categoryH.List)
			r.Post("/categories", categoryH.Create)
			r.Get("/categories/{id}", categoryH.Get)
			r.Put("/categories/{id}", categoryH.Update)
			r.Delete("/categories/{id}", categoryH.Delete)

			// Suppliers
			r.Get("/suppliers", supplierH.List)
			r.Post("/suppliers", supplierH.Create)
			r.Get("/suppliers/{id}", supplierH.Get)
			r.Put("/suppliers/{id}", supplierH.Update)
			r.Delete("/suppliers/{id}", supplierH.Delete)

			// Customers
			r.Get("/customers", customerH.List)
			r.Post("/customers", customerH.Create)
			r.Get("/customers/{id}", customerH.Get)
			r.Put("/customers/{id}", customerH.Update)
			r.Delete("/customers/{id}", customerH.Delete)

			// Employees
			r.Get("/employees", employeeH.List)
			r.Post("/employees", employeeH.Create)
			r.Get("/employees/{id}", employeeH.Get)
			r.Put("/employees/{id}", employeeH.Update)
			r.Delete("/employees/{id}", employeeH.Delete)

			// Sales
			r.Get("/sales", saleH.List)
			r.Post("/sales", saleH.Create)
			r.Get("/sales/{id}", saleH.Get)
			r.Put("/sales/{id}", saleH.Update)
			r.Delete("/sales/{id}", saleH.Delete)

			// Reports
			r.Get("/reports/sales-total", reportH.SalesTotal)
			r.Get("/reports/stock", reportH.Stock)
			r.Get("/reports/top-products", reportH.TopProducts)
			r.Get("/reports/sales-by-employee", reportH.SalesByEmployee)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

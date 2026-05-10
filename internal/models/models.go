package models

import "time"

// ─── Category ────────────────────────────────────────────────────────────────

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ─── Supplier ────────────────────────────────────────────────────────────────

type Supplier struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Country   string    `json:"country,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSupplierRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Country string `json:"country"`
}

// ─── Product ─────────────────────────────────────────────────────────────────

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Size        float64   `json:"size,omitempty"`
	Color       string    `json:"color,omitempty"`
	Brand       string    `json:"brand,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	CategoryID  *int      `json:"category_id,omitempty"`
	SupplierID  *int      `json:"supplier_id,omitempty"`
	Category    *Category `json:"category,omitempty"`
	Supplier    *Supplier `json:"supplier,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Size        float64 `json:"size"`
	Color       string  `json:"color"`
	Brand       string  `json:"brand"`
	ImageURL    string  `json:"image_url"`
	CategoryID  *int    `json:"category_id"`
	SupplierID  *int    `json:"supplier_id"`
}

// ─── Employee ────────────────────────────────────────────────────────────────

type Employee struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	HireDate  string    `json:"hire_date"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEmployeeRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// ─── User / Auth ─────────────────────────────────────────────────────────────

type User struct {
	ID           int       `json:"id"`
	EmployeeID   *int      `json:"employee_id,omitempty"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string   `json:"token"`
	User     User     `json:"user"`
	Employee Employee `json:"employee"`
}

// ─── Customer ────────────────────────────────────────────────────────────────

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

// ─── Sale ────────────────────────────────────────────────────────────────────

type SaleItem struct {
	ID        int      `json:"id"`
	SaleID    int      `json:"sale_id"`
	ProductID *int     `json:"product_id"`
	Quantity  int      `json:"quantity"`
	UnitPrice float64  `json:"unit_price"`
	Subtotal  float64  `json:"subtotal"`
	Product   *Product `json:"product,omitempty"`
}

type Sale struct {
	ID         int        `json:"id"`
	CustomerID *int       `json:"customer_id,omitempty"`
	EmployeeID *int       `json:"employee_id,omitempty"`
	SaleDate   time.Time  `json:"sale_date"`
	Total      float64    `json:"total"`
	Status     string     `json:"status"`
	Notes      string     `json:"notes,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	Customer   *Customer  `json:"customer,omitempty"`
	Employee   *Employee  `json:"employee,omitempty"`
	Items      []SaleItem `json:"items,omitempty"`
}

type CreateSaleItemRequest struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type CreateSaleRequest struct {
	CustomerID *int                    `json:"customer_id"`
	EmployeeID *int                    `json:"employee_id"`
	Notes      string                  `json:"notes"`
	Items      []CreateSaleItemRequest `json:"items"`
}

// ─── Reports ─────────────────────────────────────────────────────────────────

type SalesTotalReport struct {
	TotalSales  int     `json:"total_sales"`
	TotalAmount float64 `json:"total_amount"`
	DateFrom    string  `json:"date_from,omitempty"`
	DateTo      string  `json:"date_to,omitempty"`
}

type StockReport struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Brand       string  `json:"brand"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url,omitempty"`
}

type TopProductReport struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Brand        string  `json:"brand"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

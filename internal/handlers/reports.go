package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportHandler struct {
	DB *pgxpool.Pool
}

// GET /api/reports/sales-total?from=2026-01-01&to=2026-12-31
func (h *ReportHandler) SalesTotal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")

	query := `SELECT COUNT(*), COALESCE(SUM(total),0) FROM sales WHERE status != 'anulada'`
	args := []any{}
	n := 1

	if from != "" {
		query += " AND sale_date >= $" + strconv.Itoa(n)
		args = append(args, from)
		n++
	}
	if to != "" {
		query += " AND sale_date <= $" + strconv.Itoa(n)
		args = append(args, to)
		n++
	}

	var report models.SalesTotalReport
	report.DateFrom = from
	report.DateTo = to
	h.DB.QueryRow(context.Background(), query, args...).
		Scan(&report.TotalSales, &report.TotalAmount)
	respond(w, http.StatusOK, report)
}

// GET /api/reports/stock?low=true
func (h *ReportHandler) Stock(w http.ResponseWriter, r *http.Request) {
	lowOnly := r.URL.Query().Get("low") == "true"

	query := `
		SELECT p.id, p.name, COALESCE(p.brand,''), p.stock, p.price, COALESCE(p.image_url,'')
		FROM products p`
	if lowOnly {
		query += " WHERE p.stock <= 5"
	}
	query += " ORDER BY p.stock ASC"

	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch stock report")
		return
	}
	defer rows.Close()

	report := []models.StockReport{}
	for rows.Next() {
		var s models.StockReport
		rows.Scan(&s.ProductID, &s.ProductName, &s.Brand, &s.Stock, &s.Price, &s.ImageURL)
		report = append(report, s)
	}
	respond(w, http.StatusOK, report)
}

// GET /api/reports/top-products?limit=10
func (h *ReportHandler) TopProducts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}

	rows, err := h.DB.Query(context.Background(), `
		SELECT p.id, p.name, COALESCE(p.brand,''),
		       COALESCE(SUM(si.quantity),0)::int,
		       COALESCE(SUM(si.subtotal),0)
		FROM products p
		LEFT JOIN sale_items si ON p.id = si.product_id
		LEFT JOIN sales s       ON si.sale_id = s.id AND s.status != 'anulada'
		GROUP BY p.id, p.name, p.brand
		ORDER BY SUM(si.quantity) DESC NULLS LAST
		LIMIT $1`, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch top products")
		return
	}
	defer rows.Close()

	report := []models.TopProductReport{}
	for rows.Next() {
		var t models.TopProductReport
		rows.Scan(&t.ProductID, &t.ProductName, &t.Brand, &t.TotalSold, &t.TotalRevenue)
		report = append(report, t)
	}
	respond(w, http.StatusOK, report)
}

// GET /api/reports/sales-by-employee
func (h *ReportHandler) SalesByEmployee(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(), `
		SELECT e.id, e.first_name || ' ' || e.last_name,
		       COUNT(s.id)::int, COALESCE(SUM(s.total),0)
		FROM employees e
		LEFT JOIN sales s ON e.id = s.employee_id AND s.status != 'anulada'
		GROUP BY e.id, e.first_name, e.last_name
		ORDER BY SUM(s.total) DESC NULLS LAST`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch report")
		return
	}
	defer rows.Close()

	type Row struct {
		EmployeeID   int     `json:"employee_id"`
		EmployeeName string  `json:"employee_name"`
		TotalSales   int     `json:"total_sales"`
		TotalAmount  float64 `json:"total_amount"`
	}
	result := []Row{}
	for rows.Next() {
		var row Row
		rows.Scan(&row.EmployeeID, &row.EmployeeName, &row.TotalSales, &row.TotalAmount)
		result = append(result, row)
	}
	respond(w, http.StatusOK, result)
}

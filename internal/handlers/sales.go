package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleHandler struct {
	DB *pgxpool.Pool
}

func (h *SaleHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	status := q.Get("status")

	query := `
		SELECT s.id, s.customer_id, s.employee_id, s.sale_date, s.total,
		       s.status, COALESCE(s.notes,''), s.created_at,
		       COALESCE(c.first_name,''), COALESCE(c.last_name,''),
		       COALESCE(e.first_name,''), COALESCE(e.last_name,'')
		FROM sales s
		LEFT JOIN customers c ON s.customer_id = c.id
		LEFT JOIN employees e ON s.employee_id  = e.id
		WHERE 1=1`

	args := []any{}
	n := 1

	if from != "" {
		query += " AND s.sale_date >= $" + strconv.Itoa(n)
		args = append(args, from)
		n++
	}
	if to != "" {
		query += " AND s.sale_date <= $" + strconv.Itoa(n)
		args = append(args, to)
		n++
	}
	if status != "" {
		query += " AND s.status = $" + strconv.Itoa(n)
		args = append(args, status)
		n++
	}
	query += " ORDER BY s.id DESC"

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch sales")
		return
	}
	defer rows.Close()

	sales := []models.Sale{}
	for rows.Next() {
		var s models.Sale
		var cust models.Customer
		var emp models.Employee
		rows.Scan(
			&s.ID, &s.CustomerID, &s.EmployeeID, &s.SaleDate, &s.Total,
			&s.Status, &s.Notes, &s.CreatedAt,
			&cust.FirstName, &cust.LastName,
			&emp.FirstName, &emp.LastName,
		)
		if s.CustomerID != nil {
			s.Customer = &cust
		}
		if s.EmployeeID != nil {
			s.Employee = &emp
		}
		sales = append(sales, s)
	}
	respond(w, http.StatusOK, sales)
}

func (h *SaleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var s models.Sale
	var cust models.Customer
	var emp models.Employee

	err = h.DB.QueryRow(context.Background(), `
		SELECT s.id, s.customer_id, s.employee_id, s.sale_date, s.total,
		       s.status, COALESCE(s.notes,''), s.created_at,
		       COALESCE(c.first_name,''), COALESCE(c.last_name,''),
		       COALESCE(e.first_name,''), COALESCE(e.last_name,'')
		FROM sales s
		LEFT JOIN customers c ON s.customer_id = c.id
		LEFT JOIN employees e ON s.employee_id  = e.id
		WHERE s.id=$1`, id,
	).Scan(
		&s.ID, &s.CustomerID, &s.EmployeeID, &s.SaleDate, &s.Total,
		&s.Status, &s.Notes, &s.CreatedAt,
		&cust.FirstName, &cust.LastName,
		&emp.FirstName, &emp.LastName,
	)
	if err != nil {
		respondError(w, http.StatusNotFound, "sale not found")
		return
	}
	if s.CustomerID != nil {
		s.Customer = &cust
	}
	if s.EmployeeID != nil {
		s.Employee = &emp
	}

	// Cargar items
	itemRows, err := h.DB.Query(context.Background(), `
		SELECT si.id, si.sale_id, si.product_id, si.quantity, si.unit_price, si.subtotal,
		       COALESCE(p.name,''), COALESCE(p.image_url,'')
		FROM sale_items si
		LEFT JOIN products p ON si.product_id = p.id
		WHERE si.sale_id=$1`, id)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var item models.SaleItem
			var prod models.Product
			itemRows.Scan(
				&item.ID, &item.SaleID, &item.ProductID,
				&item.Quantity, &item.UnitPrice, &item.Subtotal,
				&prod.Name, &prod.ImageURL,
			)
			if item.ProductID != nil {
				item.Product = &prod
			}
			s.Items = append(s.Items, item)
		}
	}

	respond(w, http.StatusOK, s)
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSaleRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Items) == 0 {
		respondError(w, http.StatusBadRequest, "at least one item is required")
		return
	}

	tx, err := h.DB.Begin(context.Background())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not start transaction")
		return
	}
	defer tx.Rollback(context.Background())

	var saleID int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO sales (customer_id, employee_id, notes) VALUES ($1,$2,$3) RETURNING id`,
		req.CustomerID, req.EmployeeID, req.Notes,
	).Scan(&saleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not create sale")
		return
	}

	var total float64
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			respondError(w, http.StatusBadRequest, "quantity must be > 0")
			return
		}

		var currentStock int
		err := tx.QueryRow(context.Background(),
			`SELECT stock FROM products WHERE id=$1 FOR UPDATE`, item.ProductID,
		).Scan(&currentStock)
		if err != nil {
			respondError(w, http.StatusBadRequest, "product not found: "+strconv.Itoa(item.ProductID))
			return
		}
		if currentStock < item.Quantity {
			respondError(w, http.StatusConflict, "insufficient stock for product "+strconv.Itoa(item.ProductID))
			return
		}

		_, err = tx.Exec(context.Background(),
			`UPDATE products SET stock = stock - $1, updated_at=NOW() WHERE id=$2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "could not update stock")
			return
		}

		_, err = tx.Exec(context.Background(),
			`INSERT INTO sale_items (sale_id, product_id, quantity, unit_price)
			 VALUES ($1,$2,$3,$4)`,
			saleID, item.ProductID, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "could not insert sale item")
			return
		}
		total += float64(item.Quantity) * item.UnitPrice
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE sales SET total=$1 WHERE id=$2`, total, saleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not update total")
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		respondError(w, http.StatusInternalServerError, "could not commit transaction")
		return
	}

	// Retornar la venta completa
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(saleID))
	r2 := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	h.Get(w, r2)
}

func (h *SaleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := decode(r, &body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	_, err = h.DB.Exec(context.Background(),
		`UPDATE sales
		 SET status = COALESCE(NULLIF($1,''), status),
		     notes  = COALESCE(NULLIF($2,''), notes)
		 WHERE id=$3`,
		body.Status, body.Notes, id,
	)
	if err != nil {
		respondError(w, http.StatusNotFound, "sale not found")
		return
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(id))
	r2 := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	h.Get(w, r2)
}

func (h *SaleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.DB.Exec(context.Background(),
		`UPDATE sales SET status='anulada' WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "sale not found")
		return
	}
	respond(w, http.StatusOK, map[string]string{"message": "sale cancelled"})
}

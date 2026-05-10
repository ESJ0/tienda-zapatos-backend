package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductHandler struct {
	DB *pgxpool.Pool
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	categoryID := q.Get("category_id")
	search := q.Get("search")

	query := `
		SELECT p.id, p.name, COALESCE(p.description,''), p.price, p.stock,
		       COALESCE(p.size,0), COALESCE(p.color,''), COALESCE(p.brand,''),
		       COALESCE(p.image_url,''),
		       p.category_id, p.supplier_id, p.created_at, p.updated_at,
		       COALESCE(c.id,0), COALESCE(c.name,''),
		       COALESCE(s.id,0), COALESCE(s.name,'')
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN suppliers  s ON p.supplier_id  = s.id
		WHERE 1=1`

	args := []any{}
	argN := 1

	if categoryID != "" {
		query += " AND p.category_id = $" + strconv.Itoa(argN)
		args = append(args, categoryID)
		argN++
	}
	if search != "" {
		query += " AND (p.name ILIKE $" + strconv.Itoa(argN) +
			" OR p.brand ILIKE $" + strconv.Itoa(argN) + ")"
		args = append(args, "%"+search+"%")
		argN++
	}
	query += " ORDER BY p.id DESC"

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch products")
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		var cat models.Category
		var sup models.Supplier
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock,
			&p.Size, &p.Color, &p.Brand, &p.ImageURL,
			&p.CategoryID, &p.SupplierID, &p.CreatedAt, &p.UpdatedAt,
			&cat.ID, &cat.Name, &sup.ID, &sup.Name,
		)
		if err != nil {
			continue
		}
		if cat.ID != 0 {
			p.Category = &cat
		}
		if sup.ID != 0 {
			p.Supplier = &sup
		}
		products = append(products, p)
	}
	respond(w, http.StatusOK, products)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p models.Product
	var cat models.Category
	var sup models.Supplier

	err = h.DB.QueryRow(context.Background(), `
		SELECT p.id, p.name, COALESCE(p.description,''), p.price, p.stock,
		       COALESCE(p.size,0), COALESCE(p.color,''), COALESCE(p.brand,''),
		       COALESCE(p.image_url,''),
		       p.category_id, p.supplier_id, p.created_at, p.updated_at,
		       COALESCE(c.id,0), COALESCE(c.name,''),
		       COALESCE(s.id,0), COALESCE(s.name,'')
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN suppliers  s ON p.supplier_id  = s.id
		WHERE p.id=$1`, id,
	).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock,
		&p.Size, &p.Color, &p.Brand, &p.ImageURL,
		&p.CategoryID, &p.SupplierID, &p.CreatedAt, &p.UpdatedAt,
		&cat.ID, &cat.Name, &sup.ID, &sup.Name,
	)
	if err != nil {
		respondError(w, http.StatusNotFound, "product not found")
		return
	}
	if cat.ID != 0 {
		p.Category = &cat
	}
	if sup.ID != 0 {
		p.Supplier = &sup
	}
	respond(w, http.StatusOK, p)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Price < 0 {
		respondError(w, http.StatusBadRequest, "price must be >= 0")
		return
	}

	var p models.Product
	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO products
		  (name, description, price, stock, size, color, brand, image_url, category_id, supplier_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, name, COALESCE(description,''), price, stock,
		          COALESCE(size,0), COALESCE(color,''), COALESCE(brand,''),
		          COALESCE(image_url,''), category_id, supplier_id, created_at, updated_at`,
		req.Name, req.Description, req.Price, req.Stock, req.Size,
		req.Color, req.Brand, req.ImageURL, req.CategoryID, req.SupplierID,
	).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock,
		&p.Size, &p.Color, &p.Brand, &p.ImageURL,
		&p.CategoryID, &p.SupplierID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not create product: "+err.Error())
		return
	}
	respond(w, http.StatusCreated, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req models.CreateProductRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var p models.Product
	err = h.DB.QueryRow(context.Background(), `
		UPDATE products
		SET name=$1, description=$2, price=$3, stock=$4, size=$5,
		    color=$6, brand=$7, image_url=$8, category_id=$9, supplier_id=$10,
		    updated_at=NOW()
		WHERE id=$11
		RETURNING id, name, COALESCE(description,''), price, stock,
		          COALESCE(size,0), COALESCE(color,''), COALESCE(brand,''),
		          COALESCE(image_url,''), category_id, supplier_id, created_at, updated_at`,
		req.Name, req.Description, req.Price, req.Stock, req.Size,
		req.Color, req.Brand, req.ImageURL, req.CategoryID, req.SupplierID, id,
	).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock,
		&p.Size, &p.Color, &p.Brand, &p.ImageURL,
		&p.CategoryID, &p.SupplierID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		respondError(w, http.StatusNotFound, "product not found")
		return
	}
	respond(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.DB.Exec(context.Background(),
		`DELETE FROM products WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "product not found")
		return
	}
	respond(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryHandler struct {
	DB *pgxpool.Pool
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, name, COALESCE(description,''), created_at
		 FROM categories ORDER BY name`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch categories")
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
		categories = append(categories, c)
	}
	respond(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var c models.Category
	err = h.DB.QueryRow(context.Background(),
		`SELECT id, name, COALESCE(description,''), created_at
		 FROM categories WHERE id=$1`, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "category not found")
		return
	}
	respond(w, http.StatusOK, c)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	var c models.Category
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO categories (name, description) VALUES ($1,$2)
		 RETURNING id, name, COALESCE(description,''), created_at`,
		req.Name, req.Description,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		respondError(w, http.StatusConflict, "category already exists or db error")
		return
	}
	respond(w, http.StatusCreated, c)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req models.CreateCategoryRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var c models.Category
	err = h.DB.QueryRow(context.Background(),
		`UPDATE categories SET name=$1, description=$2 WHERE id=$3
		 RETURNING id, name, COALESCE(description,''), created_at`,
		req.Name, req.Description, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "category not found")
		return
	}
	respond(w, http.StatusOK, c)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.DB.Exec(context.Background(),
		`DELETE FROM categories WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "category not found")
		return
	}
	respond(w, http.StatusOK, map[string]string{"message": "category deleted"})
}

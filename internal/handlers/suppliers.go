package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierHandler struct {
	DB *pgxpool.Pool
}

func (h *SupplierHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, name, COALESCE(contact,''), COALESCE(email,''),
		        COALESCE(phone,''), COALESCE(country,''), created_at
		 FROM suppliers ORDER BY name`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch suppliers")
		return
	}
	defer rows.Close()

	suppliers := []models.Supplier{}
	for rows.Next() {
		var s models.Supplier
		rows.Scan(&s.ID, &s.Name, &s.Contact, &s.Email, &s.Phone, &s.Country, &s.CreatedAt)
		suppliers = append(suppliers, s)
	}
	respond(w, http.StatusOK, suppliers)
}

func (h *SupplierHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var s models.Supplier
	err = h.DB.QueryRow(context.Background(),
		`SELECT id, name, COALESCE(contact,''), COALESCE(email,''),
		        COALESCE(phone,''), COALESCE(country,''), created_at
		 FROM suppliers WHERE id=$1`, id,
	).Scan(&s.ID, &s.Name, &s.Contact, &s.Email, &s.Phone, &s.Country, &s.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "supplier not found")
		return
	}
	respond(w, http.StatusOK, s)
}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSupplierRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	var s models.Supplier
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO suppliers (name, contact, email, phone, country)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, name, COALESCE(contact,''), COALESCE(email,''),
		           COALESCE(phone,''), COALESCE(country,''), created_at`,
		req.Name, req.Contact, req.Email, req.Phone, req.Country,
	).Scan(&s.ID, &s.Name, &s.Contact, &s.Email, &s.Phone, &s.Country, &s.CreatedAt)
	if err != nil {
		respondError(w, http.StatusConflict, "supplier already exists or db error")
		return
	}
	respond(w, http.StatusCreated, s)
}

func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req models.CreateSupplierRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var s models.Supplier
	err = h.DB.QueryRow(context.Background(),
		`UPDATE suppliers SET name=$1, contact=$2, email=$3, phone=$4, country=$5
		 WHERE id=$6
		 RETURNING id, name, COALESCE(contact,''), COALESCE(email,''),
		           COALESCE(phone,''), COALESCE(country,''), created_at`,
		req.Name, req.Contact, req.Email, req.Phone, req.Country, id,
	).Scan(&s.ID, &s.Name, &s.Contact, &s.Email, &s.Phone, &s.Country, &s.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "supplier not found")
		return
	}
	respond(w, http.StatusOK, s)
}

func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.DB.Exec(context.Background(),
		`DELETE FROM suppliers WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "supplier not found")
		return
	}
	respond(w, http.StatusOK, map[string]string{"message": "supplier deleted"})
}

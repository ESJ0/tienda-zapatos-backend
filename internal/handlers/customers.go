package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tu-usuario/zapatos-api/internal/models"
)

type CustomerHandler struct {
	DB *pgxpool.Pool
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, first_name, last_name, COALESCE(email,''),
		        COALESCE(phone,''), COALESCE(address,''), created_at
		 FROM customers ORDER BY last_name, first_name`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch customers")
		return
	}
	defer rows.Close()

	customers := []models.Customer{}
	for rows.Next() {
		var c models.Customer
		rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Address, &c.CreatedAt)
		customers = append(customers, c)
	}
	respond(w, http.StatusOK, customers)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var c models.Customer
	err = h.DB.QueryRow(context.Background(),
		`SELECT id, first_name, last_name, COALESCE(email,''),
		        COALESCE(phone,''), COALESCE(address,''), created_at
		 FROM customers WHERE id=$1`, id,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Address, &c.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "customer not found")
		return
	}
	respond(w, http.StatusOK, c)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.FirstName == "" || req.LastName == "" {
		respondError(w, http.StatusBadRequest, "first_name and last_name are required")
		return
	}
	var c models.Customer
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO customers (first_name, last_name, email, phone, address)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, first_name, last_name, COALESCE(email,''),
		           COALESCE(phone,''), COALESCE(address,''), created_at`,
		req.FirstName, req.LastName, req.Email, req.Phone, req.Address,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Address, &c.CreatedAt)
	if err != nil {
		respondError(w, http.StatusConflict, "customer already exists or db error")
		return
	}
	respond(w, http.StatusCreated, c)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req models.CreateCustomerRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var c models.Customer
	err = h.DB.QueryRow(context.Background(),
		`UPDATE customers SET first_name=$1, last_name=$2, email=$3, phone=$4, address=$5
		 WHERE id=$6
		 RETURNING id, first_name, last_name, COALESCE(email,''),
		           COALESCE(phone,''), COALESCE(address,''), created_at`,
		req.FirstName, req.LastName, req.Email, req.Phone, req.Address, id,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Address, &c.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "customer not found")
		return
	}
	respond(w, http.StatusOK, c)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.DB.Exec(context.Background(),
		`DELETE FROM customers WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "customer not found")
		return
	}
	respond(w, http.StatusOK, map[string]string{"message": "customer deleted"})
}

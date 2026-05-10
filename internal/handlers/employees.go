package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeHandler struct {
	DB *pgxpool.Pool
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, first_name, last_name, email, COALESCE(role,''), hire_date::text, created_at
		 FROM employees ORDER BY last_name, first_name`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not fetch employees")
		return
	}
	defer rows.Close()

	employees := []models.Employee{}
	for rows.Next() {
		var e models.Employee
		rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Role, &e.HireDate, &e.CreatedAt)
		employees = append(employees, e)
	}
	respond(w, http.StatusOK, employees)
}

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var e models.Employee
	err = h.DB.QueryRow(context.Background(),
		`SELECT id, first_name, last_name, email, COALESCE(role,''), hire_date::text, created_at
		 FROM employees WHERE id=$1`, id,
	).Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Role, &e.HireDate, &e.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "employee not found")
		return
	}
	respond(w, http.StatusOK, e)
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		respondError(w, http.StatusBadRequest, "first_name, last_name and email are required")
		return
	}
	if req.Role == "" {
		req.Role = "vendedor"
	}
	var e models.Employee
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO employees (first_name, last_name, email, role)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, first_name, last_name, email, COALESCE(role,''), hire_date::text, created_at`,
		req.FirstName, req.LastName, req.Email, req.Role,
	).Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Role, &e.HireDate, &e.CreatedAt)
	if err != nil {
		respondError(w, http.StatusConflict, "employee already exists or db error")
		return
	}
	respond(w, http.StatusCreated, e)
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req models.CreateEmployeeRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var e models.Employee
	err = h.DB.QueryRow(context.Background(),
		`UPDATE employees SET first_name=$1, last_name=$2, email=$3, role=$4
		 WHERE id=$5
		 RETURNING id, first_name, last_name, email, COALESCE(role,''), hire_date::text, created_at`,
		req.FirstName, req.LastName, req.Email, req.Role, id,
	).Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Role, &e.HireDate, &e.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "employee not found")
		return
	}
	respond(w, http.StatusOK, e)
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.DB.Exec(context.Background(),
		`DELETE FROM employees WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "employee not found")
		return
	}
	respond(w, http.StatusOK, map[string]string{"message": "employee deleted"})
}

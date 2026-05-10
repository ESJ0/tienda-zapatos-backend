package handlers

import (
	"context"
	"net/http"

	"github.com/ESJ0/tienda-zapatos-backend/internal/config"
	"github.com/ESJ0/tienda-zapatos-backend/internal/middleware"
	"github.com/ESJ0/tienda-zapatos-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB  *pgxpool.Pool
	Cfg *config.Config
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	var user models.User
	var hash string
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, employee_id, username, password_hash, created_at
		 FROM users WHERE username = $1`, req.Username,
	).Scan(&user.ID, &user.EmployeeID, &user.Username, &hash, &user.CreatedAt)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	var emp models.Employee
	if user.EmployeeID != nil {
		h.DB.QueryRow(context.Background(),
			`SELECT id, first_name, last_name, email, role, hire_date::text, created_at
			 FROM employees WHERE id = $1`, *user.EmployeeID,
		).Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Role, &emp.HireDate, &emp.CreatedAt)
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, h.Cfg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	respond(w, http.StatusOK, models.LoginResponse{
		Token:    token,
		User:     user,
		Employee: emp,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

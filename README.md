# 👟 ESJ0 Footwear — Backend

API REST en Go para la tienda de zapatos ESJ0. Construida con **Chi**, **pgx** y **JWT**.

## Stack
- **Go 1.22** + Chi router
- **PostgreSQL 16**
- **Docker / Docker Compose**

---

## 🚀 Levantar con Docker (recomendado)

```bash
# 1. Clonar el repositorio
git clone <url-del-repo>
cd tienda-zapatos-backend

# 2. Copiar variables de entorno
cp .env.example .env

# 3. Levantar todo con un solo comando
docker compose up --build
```

La API estará disponible en `http://localhost:8080`.
La base de datos se inicializa automáticamente con el schema y datos de prueba.

---

## 🛠️ Desarrollo local (sin Docker)

```bash
# Requisitos: Go 1.22+, PostgreSQL corriendo localmente

# 1. Instalar dependencias
go mod download

# 2. Editar .env apuntando a localhost
cp .env.example .env
# Cambiar DB_HOST=localhost

# 3. Levantar solo la DB
docker compose up -d db

# 4. Correr el servidor
go run ./cmd/server/main.go
```

---

## 🔑 Credenciales de prueba

| Campo    | Valor      |
|----------|------------|
| Usuario  | `admin`    |
| Password | `admin123` |

---

## 🗄️ Credenciales de base de datos

| Campo    | Valor     |
|----------|-----------|
| Usuario  | `proy2`   |
| Password | `secret`  |
| DB       | `zapatos` |
| Puerto   | `5432`    |

---

## 📡 Endpoints

### Auth (público)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/auth/login` | Login → devuelve JWT |
| POST | `/api/auth/logout` | Logout |

### Productos (requiere `Authorization: Bearer <token>`)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/products` | Listar productos |
| GET | `/api/products?search=nike` | Buscar por nombre o marca |
| GET | `/api/products?category_id=1` | Filtrar por categoría |
| POST | `/api/products` | Crear producto |
| GET | `/api/products/{id}` | Obtener producto |
| PUT | `/api/products/{id}` | Actualizar producto |
| DELETE | `/api/products/{id}` | Eliminar producto |

### Categorías
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/categories` | Listar |
| POST | `/api/categories` | Crear |
| GET | `/api/categories/{id}` | Obtener |
| PUT | `/api/categories/{id}` | Actualizar |
| DELETE | `/api/categories/{id}` | Eliminar |

### Proveedores
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/suppliers` | Listar |
| POST | `/api/suppliers` | Crear |
| GET | `/api/suppliers/{id}` | Obtener |
| PUT | `/api/suppliers/{id}` | Actualizar |
| DELETE | `/api/suppliers/{id}` | Eliminar |

### Clientes
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/customers` | Listar |
| POST | `/api/customers` | Crear |
| GET | `/api/customers/{id}` | Obtener |
| PUT | `/api/customers/{id}` | Actualizar |
| DELETE | `/api/customers/{id}` | Eliminar |

### Empleados
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/employees` | Listar |
| POST | `/api/employees` | Crear |
| GET | `/api/employees/{id}` | Obtener |
| PUT | `/api/employees/{id}` | Actualizar |
| DELETE | `/api/employees/{id}` | Eliminar |

### Ventas
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/sales` | Listar ventas |
| GET | `/api/sales?from=2026-01-01&to=2026-12-31` | Filtrar por fecha |
| GET | `/api/sales?status=completada` | Filtrar por estado |
| POST | `/api/sales` | Crear venta |
| GET | `/api/sales/{id}` | Obtener venta con items |
| PUT | `/api/sales/{id}` | Actualizar estado/notas |
| DELETE | `/api/sales/{id}` | Anular venta |

### Reportes
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/reports/sales-total` | Total ventas e ingresos |
| GET | `/api/reports/sales-total?from=2026-01-01&to=2026-12-31` | Con rango de fechas |
| GET | `/api/reports/stock` | Stock por producto |
| GET | `/api/reports/stock?low=true` | Solo stock crítico (≤5) |
| GET | `/api/reports/top-products` | Productos más vendidos |
| GET | `/api/reports/top-products?limit=5` | Top N productos |
| GET | `/api/reports/sales-by-employee` | Ventas por empleado |

---

## 📦 Ejemplo — Crear una venta

```json
POST /api/sales
Authorization: Bearer <token>
Content-Type: application/json

{
  "customer_id": 1,
  "employee_id": 2,
  "notes": "Pago en efectivo",
  "items": [
    { "product_id": 1, "quantity": 2, "unit_price": 950.00 },
    { "product_id": 5, "quantity": 1, "unit_price": 299.99 }
  ]
}
```

---

## 🐳 Variables de entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `DB_HOST` | Host de la base de datos | `db` |
| `DB_PORT` | Puerto PostgreSQL | `5432` |
| `DB_USER` | Usuario | `proy2` |
| `DB_PASSWORD` | Contraseña | `secret` |
| `DB_NAME` | Nombre de la DB | `zapatos` |
| `JWT_SECRET` | Secreto para firmar JWT | — |
| `JWT_EXPIRY_HOURS` | Duración del token | `24` |
| `PORT` | Puerto del servidor | `8080` |

---

## 📁 Estructura del proyecto

```
tienda-zapatos-backend/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── db/db.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── products.go
│   │   ├── categories.go
│   │   ├── suppliers.go
│   │   ├── customers.go
│   │   ├── employees.go
│   │   ├── sales.go
│   │   ├── reports.go
│   │   └── helpers.go
│   └── models/models.go
├── migrations/
│   ├── 001_init.sql
│── seeds/
│   ├── 001_seed.sql
├── .env.example
├── docker-compose.yml
├── Dockerfile
└── README.md
```
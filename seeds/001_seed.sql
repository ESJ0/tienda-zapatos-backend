-- ─── Categorías ───────────────────────────────────────────────────────────────

INSERT INTO categories (name, description) VALUES
    ('Deportivo', 'Tenis y zapatillas para deporte y gym'),
    ('Casual',    'Zapatos para uso diario informal'),
    ('Formal',    'Zapatos elegantes para ocasiones formales'),
    ('Sandalia',  'Sandalias y chanclas para clima cálido'),
    ('Bota',      'Botas y botines para todo terreno')
ON CONFLICT (name) DO NOTHING;

-- ─── Proveedores ──────────────────────────────────────────────────────────────

INSERT INTO suppliers (name, contact, email, phone, country) VALUES
    ('Nike Guatemala',    'Carlos López',   'carlos@nikegt.com',     '2222-1111', 'Guatemala'),
    ('Adidas Imports',    'María Pérez',    'maria@adidas-imp.com',  '2333-2222', 'México'),
    ('New Balance GT',    'Luis Morales',   'luis@nbgt.com',         '2444-3333', 'Guatemala'),
    ('EuroShoe SA',       'Pierre Martin',  'pierre@euroshoe.eu',    '+34600000', 'España')
ON CONFLICT (email) DO NOTHING;

-- ─── Empleados ────────────────────────────────────────────────────────────────

INSERT INTO employees (first_name, last_name, email, role) VALUES
    ('Admin',    'Sistema',    'admin@zapatos.gt',    'administrador'),
    ('Laura',    'Monterroso', 'laura@zapatos.gt',    'vendedor'),
    ('Diego',    'Cifuentes',  'diego@zapatos.gt',    'vendedor'),
    ('Patricia', 'Alvarado',   'patricia@zapatos.gt', 'supervisor')
ON CONFLICT (email) DO NOTHING;

-- ─── Usuario admin ────────────────────────────────────────────────────────────
-- password: admin123

INSERT INTO users (employee_id, username, password_hash) VALUES
    (1, 'admin', '$2a$10$P4nflM86wGUJMg2OLACHU.vZ5oCivTeofEwq4bWW3F2jQ/VhH5fLa')
ON CONFLICT (username) DO NOTHING;

-- ─── Clientes ─────────────────────────────────────────────────────────────────

INSERT INTO customers (first_name, last_name, email, phone, address) VALUES
    ('Juan',    'Rodríguez', 'juan@gmail.com',    '5111-2222', 'Zona 1, Guatemala City'),
    ('Sofía',   'Herrera',   'sofia@gmail.com',   '5222-3333', 'Zona 10, Guatemala City'),
    ('Roberto', 'Méndez',    'roberto@gmail.com', '5333-4444', 'Mixco, Guatemala'),
    ('Lucía',   'Castillo',  'lucia@gmail.com',   '5444-5555', 'Villa Nueva, Guatemala')
ON CONFLICT (email) DO NOTHING;

-- ─── Productos ────────────────────────────────────────────────────────────────

INSERT INTO products (name, description, price, stock, size, color, brand, image_url, category_id, supplier_id) VALUES
    ('Nike Air Force 1',
     'El clásico de clásicos. Suela de aire, cuero blanco icónico.',
     950.00, 30, 42, 'Blanco', 'Nike', '/shoes/nike-air-force-1.png', 2, 1),

    ('Nike Air Max 90',
     'Amortiguación visible Air Max con estilo retro de los 90s.',
     1150.00, 20, 41, 'Negro/Rojo', 'Nike', '/shoes/nike-air-max-90.png', 1, 1),

    ('Adidas Stan Smith',
     'Tenis de tenis reconvertido en ícono del streetwear mundial.',
     899.00, 25, 40, 'Blanco/Verde', 'Adidas', '/shoes/adidas-stan-smith.png', 2, 2),

    ('Adidas Ultraboost 22',
     'Máxima amortiguación con tecnología Boost para corredores.',
     1399.00, 15, 43, 'Gris', 'Adidas', '/shoes/adidas-ultraboost-22.png', 1, 2),

    ('Converse Chuck Taylor All Star',
     'La bota de lona más vendida de la historia.',
     699.00, 35, 39, 'Negro', 'Converse', '/shoes/converse-chuck-taylor-all-star.png', 2, 1),

    ('New Balance 574',
     'Perfil retro con amortiguación ENCAP. Cómodo todo el día.',
     1050.00, 18, 42, 'Azul marino', 'New Balance', '/shoes/new-balance-574.png', 2, 3),

    ('Vans Old Skool',
     'La primera Vans con la franja lateral. Skate y street.',
     799.00, 22, 41, 'Negro/Blanco', 'Vans', '/shoes/vans-old-skool.jpg', 2, 3),

    ('Nike Air Jordan 1 Retro',
     'El zapato que cambió la historia del basketball y la moda.',
     2199.00, 10, 44, 'Rojo/Negro', 'Nike', '/shoes/nike-air-jordan-1-retro.png', 1, 1),

    ('Adidas Samba OG',
     'Originalmente para fútbol en interiores, ahora ícono de moda.',
     949.00, 20, 40, 'Blanco/Negro', 'Adidas', '/shoes/adidas-samba-og.png', 2, 2),

    ('Timberland 6-Inch Boot',
     'Bota impermeable de cuero premium, icónica en color trigo.',
     1899.00, 12, 43, 'Trigo', 'Timberland', '/shoes/timberland-6-inch-boot.png', 5, 4)
ON CONFLICT DO NOTHING;
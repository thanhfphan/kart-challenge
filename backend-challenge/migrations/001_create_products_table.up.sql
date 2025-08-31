-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sku VARCHAR(255) NOT NULL DEFAULT '',
    name VARCHAR(255) NOT NULL DEFAULT '',
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    category VARCHAR(255) NOT NULL DEFAULT '',
    thumbnail_url TEXT DEFAULT '',
    mobile_url TEXT DEFAULT '',
    tablet_url TEXT DEFAULT '',
    desktop_url TEXT DEFAULT '',
    description TEXT DEFAULT '',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    
    INDEX idx_products_sku (sku),
    INDEX idx_products_category (category),
    INDEX idx_products_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed product data
INSERT INTO products (id, sku, name, price, category, thumbnail_url, mobile_url, tablet_url, desktop_url, description, created_at, updated_at) VALUES
(1, 'WAFFLE-001', 'Chicken Waffle', 13.30, 'Waffle', 'https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg', 'https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg', 'https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg', 'https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg', 'Delicious chicken waffle with crispy exterior and fluffy interior', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW())),
(2, 'BURGER-001', 'Classic Beef Burger', 15.50, 'Burger', 'https://orderfoodonline.deno.dev/public/images/image-burger-thumbnail.jpg', 'https://orderfoodonline.deno.dev/public/images/image-burger-mobile.jpg', 'https://orderfoodonline.deno.dev/public/images/image-burger-tablet.jpg', 'https://orderfoodonline.deno.dev/public/images/image-burger-desktop.jpg', 'Juicy beef patty with fresh lettuce, tomato, and special sauce', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW())),
(3, 'PIZZA-001', 'Margherita Pizza', 18.75, 'Pizza', 'https://orderfoodonline.deno.dev/public/images/image-pizza-thumbnail.jpg', 'https://orderfoodonline.deno.dev/public/images/image-pizza-mobile.jpg', 'https://orderfoodonline.deno.dev/public/images/image-pizza-tablet.jpg', 'https://orderfoodonline.deno.dev/public/images/image-pizza-desktop.jpg', 'Traditional Italian pizza with fresh mozzarella and basil', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW())),
(4, 'PASTA-001', 'Spaghetti Carbonara', 16.25, 'Pasta', 'https://orderfoodonline.deno.dev/public/images/image-pasta-thumbnail.jpg', 'https://orderfoodonline.deno.dev/public/images/image-pasta-mobile.jpg', 'https://orderfoodonline.deno.dev/public/images/image-pasta-tablet.jpg', 'https://orderfoodonline.deno.dev/public/images/image-pasta-desktop.jpg', 'Creamy pasta with bacon, eggs, and parmesan cheese', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW())),
(5, 'SALAD-001', 'Caesar Salad', 12.90, 'Salad', 'https://orderfoodonline.deno.dev/public/images/image-salad-thumbnail.jpg', 'https://orderfoodonline.deno.dev/public/images/image-salad-mobile.jpg', 'https://orderfoodonline.deno.dev/public/images/image-salad-tablet.jpg', 'https://orderfoodonline.deno.dev/public/images/image-salad-desktop.jpg', 'Fresh romaine lettuce with caesar dressing and croutons', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW()));



-- Membuat tabel categories terlebih dahulu karena akan direferensikan oleh blogs
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE, -- soft delete
);

-- Membuat tabel blogs
CREATE TABLE IF NOT EXISTS blogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    excerpt TEXT, -- Cuplikan pendek untuk card di frontend
    cover_image_url VARCHAR(255),
    
    -- Status publikasi: draft, published, atau archived
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    
    -- Foreign Keys
    author_id UUID NOT NULL,
    category_id UUID, -- Bisa null jika Anda memperbolehkan blog tanpa kategori (Uncategorized)
    
    -- Timestamps
    published_at TIMESTAMP WITH TIME ZONE, -- Terpisah dari created_at untuk fitur schedule publish
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE, -- Soft delete implementation

    -- Constraints
    CONSTRAINT fk_blogs_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_blogs_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- BEST PRACTICE: Membuat Index untuk optimasi performa pencarian
-- Next.js akan sangat sering melakukan fetch berdasarkan slug untuk ISR/SSG
CREATE INDEX idx_blogs_slug ON blogs(slug);
CREATE INDEX idx_categories_slug ON categories(slug);

-- Index untuk filter status dan mempercepat load halaman blog public (hanya yang published)
CREATE INDEX idx_blogs_status_published ON blogs(status) WHERE status = 'published';
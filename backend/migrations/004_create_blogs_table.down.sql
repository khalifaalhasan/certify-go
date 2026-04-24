-- Hapus index terlebih dahulu (opsional, karena drop table otomatis menghapus index, 
-- tapi lebih clean code jika ditulis eksplisit)
DROP INDEX IF EXISTS idx_blogs_status_published;
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_blogs_slug;

-- Hapus tabel blogs terlebih dahulu karena memiliki foreign key ke categories
DROP TABLE IF EXISTS blogs;

-- Setelah blogs terhapus, aman untuk menghapus categories
DROP TABLE IF EXISTS categories;
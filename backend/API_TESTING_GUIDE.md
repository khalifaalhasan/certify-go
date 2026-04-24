# Panduan Alur Pengujian API (API Testing Flow)
> Dokumen ini menjelaskan urutan yang benar dalam menguji sistem *Next-Go Starter Kit* menggunakan Postman Collection yang disediakan, untuk mencegah pengujian membentur error relasi / Auth. 

---

## 1. Fase Inisiasi & Autentikasi (WAJIB)
Karena 90% manipulasi dikunci oleh akses Role dan Authorization, Anda diwajibkan untuk *Login* terlebih dahulu.

1. Buka folder **Public Namespace > Auth**.
2. Jalankan request **`Login`**.
   *(Pastikan email dan password Admin valid, sesuaikan dengan apa yang disemai di Seeders database Anda)*
3. Dari *Response Body*-nya, salin nilai `access_token`.
4. Buka **Postman Root Collection** (klik nama koleksi "Next-Go Starter Kit API" di sidebar kiri).
5. Pergi ke tab **Variables**, lalu _paste_ di baris variabel `token`.
6. Simpan konfigurasi koleksi tersebut (`Cmd+S` atau `Ctrl+S`). Semua akses di folder *Admin* kini akan terotorisasi.

---

## 2. Fase Profil & Setup Role (Opsional)
Pastikan status profil dan Role akun yang sedang dipakai memiliki hak akses tertinggi jika ingin mentest keseluruhan fitur.
1. Masuk ke **Admin Namespace > Auth / Profile**.
2. Jalankan **`Get Profile`** untuk melihat data akun beserta list `roles` yang terhubung. Salin `id` *User* Anda bila nanti Anda ingin memodifikasi profile ini.
3. Jalankan **Admin Namespace > RBAC > Roles > `Get All Roles`** untuk mendaftar Role apa saja yang tersedia di sistem saat ini. 

---

## 3. Fase Manipulasi Kategori (Prerequisite/Syarat)
Syarat utama dari membuat sebuah *Blog* adalah adanya sebuah *Kategori*. Oleh karena itu, kita **tidak bisa meng-create blog jika kita belum membuat sebuah Kategori**.

1. Buka folder **Admin Namespace > Categories**.
2. Jalankan request **`Create Category`**. (Ganti body sesuai selera jika perlu).
3. Anda akan mendapatkan response berisikan profil struktur `Kategori` baru. **PENTING: Salin UUID yang ada di object `id` dari response ini!**
4. Tes operasi perubahannya menggunakan **`Update Category`**. Jangan lupa replace param url `:id` dengan UUID barusan.

---

## 4. Fase Publikasi Artikel (Blog)
Setelah memiliki Kategori yang sah, barulah sistem bisa menerima request pembuatan blog post.

1. Buka folder **Admin Namespace > Blogs**.
2. Klik request **`Create Blog`**. 
3. Pergi ke bagian payload (Body `raw`). Ganti form properti `"category_id"` dengan *UUID kategori* yang sudah kita copy pada Fase ke-3 tadi.
4. *Send Request*. Anda akan mendapatkan kembalian data sebuah blog. Salin **UUID** blog tersebut, serta catat bahwa status bawaannya adalah `"draft"`.
5. Klik request **`Update Blog`** (pastikan URL param dinamis `:id` terisi).
6. Di bagian Body, kirimkan form json `{"status": "published"}` untuk mengubah artikel ke keadaan siap rilis komersial.

---

## 5. Fase Validasi Output (Testing Public Namespace)
Ini adalah fase pembuktian di mana *Role rules* dan *Status Filtration* bekerja sempurna. Semua namespace ini diakses secara terbuka tanpa perlu otorisasi.

1. Tutup folder Admin, dan jelajahi folder **Public Namespace > Blogs**.
2. Jalankan request **`Get All Blogs (Published)`**. 
   > **Validasi:** Anda seharusnya bisa melihat Blog yang baru saja rilis. Artikel lain yang tidak berstatus `"published"` tidak boleh/tidak akan dimunculkan oleh server.
3. Coba panggil spesifik via slug menggunakan **`Get Blog by Slug`**.

---

## 6. Fase Destruksi & Pembersihan (Teardown)
Pengujian fase ini berfokus pada tes Delete dan validasi hilangnya data. Karena modul menggunakan *Soft-Delete*, kita mengamati fungsionalitas sembunyi datanya.

1. Buka ulang folder **Admin Namespace > Blogs**.
2. Jalankan **`Delete Blog`** menggunakan *UUID* blog tadi. Kembalian haruslah pesan Sukses.
3. Buka **Public Namespace > Blogs**, jalankan ulang **`Get All Blogs (Published)`**. Artikel terkait mestinya sudah hilang dari radar.
4. Lakukan pembersihan serupa untuk Kategori dengan mengeksekusi **`Delete Category`**.
5. Tutup koneksi anda dan hapuskan riwayat dengan menjajal **`Logout`** pada Admin Namespace.

> **Catatan Endpoint Pencarian:**
> Untuk API `GetAll` (baik itu Blogs, Roles, atau Categories), tambahkan saja atribut kueri di ujung *Request URL* seperti `?search=kucing&sort_by=created_at&sort_order=desc`. Ini sesuai dengan standar *best practice REST*.

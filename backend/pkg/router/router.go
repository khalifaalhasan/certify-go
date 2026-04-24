package router

import (
	"compress/zlib"
	"net/http"

	"github.com/base-go/backend/internal/auth"
	"github.com/base-go/backend/internal/blog"
	"github.com/base-go/backend/internal/category"
	"github.com/base-go/backend/internal/rbac"
	"github.com/base-go/backend/pkg/middleware"
	"github.com/base-go/backend/pkg/response"
	"github.com/go-chi/chi/v5"
	cmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

// SetupRoutes this function for centralize setup all route in this app.
// why wee need to centralize?, it's for easies debugging if any issue
//
// parameters: all interface handlers we need to expose with rest
func SetupRoutes(
	authHandler auth.Handler,
	rbacHandler rbac.Handler,
	rbacRepo rbac.Repository,
	categoryHandler category.Handler,
	blogHandler blog.Handler,
) *chi.Mux {
	mux := chi.NewRouter()

	// chi middleware
	mux.Use(cmiddleware.Logger)
	mux.Use(cmiddleware.Recoverer)
	mux.Use(cmiddleware.RealIP)
	mux.Use(cmiddleware.NoCache)
	mux.Use(cmiddleware.GetHead)
	mux.Use(cmiddleware.Compress(zlib.BestCompression))
	mux.Use(cmiddleware.AllowContentType("application/json"))
	mux.Use(secure.New(secure.Options{
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		STSSeconds:           900,
	}).Handler)

	mux.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.JSON{Code: http.StatusMethodNotAllowed, Message: "Route method not allowed"}
		response.ResponseJSON(w, res.Code, res)
	}))

	mux.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.JSON{Code: http.StatusNotFound, Message: "Route not found"}
		response.ResponseJSON(w, res.Code, res)
	}))

	// set cors middleware
	mux.Use(middleware.Cors())
	// set middleware rate limiter
	mux.Use(middleware.RateLimit(1000, 10))

	// set prefix v1
	mux.Route("/v1", func(r chi.Router) {

		// ==========================================
		// PUBLIC NAMESPACE
		// ==========================================
		r.Route("/public", func(r chi.Router) {
			
			// Auth Routes (Public)
			r.Route("/auth", func(r chi.Router) {
				// No registration point in public namespace (closed system)
				r.Post("/login", authHandler.Login)
				r.Post("/oauth/google", authHandler.LoginWithGoogle)
				r.Post("/refresh", authHandler.RefreshToken)
			})

			// Categories (Read-Only Public)
			r.Route("/categories", func(r chi.Router) {
				// Supports Search & Sort via Query Parameters
				r.Get("/", categoryHandler.GetActive)
				r.Get("/{id}", categoryHandler.GetByID)
				r.Get("/slug/{slug}", categoryHandler.GetBySlug)
			})

			// Blogs (Read-Only Public)
			r.Route("/blogs", func(r chi.Router) {
				// Supports Search & Sort via Query Parameters
				r.Get("/", blogHandler.GetPublishedBlogs)
				r.Get("/{id}", blogHandler.GetByID)
				r.Get("/slug/{slug}", blogHandler.GetPublishedBySlug)
			})
		})

		// ==========================================
		// ADMIN NAMESPACE (Protected by Global JWT)
		// ==========================================
		r.Route("/admin", func(r chi.Router) {
			// Apply JWT protection for all /admin routes
			r.Use(middleware.JWTAuthMiddleware)

			// Auth Routes (Protected)
			r.Route("/auth", func(r chi.Router) {
				r.Post("/logout", authHandler.Logout)
			})

			// Profile Management
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", authHandler.GetProfile)
				r.Put("/", authHandler.UpdateProfile)
				r.Post("/change-password", authHandler.ChangePassword)
			})

			// Categories CRUD (Admin/Editor)
			r.Route("/categories", func(r chi.Router) {
				r.Use(middleware.RequireRole("Super Admin", "Admin", "Editor"))
				
				// Supports Search & Pagination via Query Parameters
				r.Get("/", categoryHandler.GetAll)
				r.Get("/{id}", categoryHandler.GetByID)
				r.Post("/", categoryHandler.Create)
				r.Put("/{id}", categoryHandler.Update)
				r.Delete("/{id}", categoryHandler.Delete)
			})

			// Blogs CRUD (Author/Editor/Admin)
			r.Route("/blogs", func(r chi.Router) {
				r.Use(middleware.RequireRole("Super Admin", "Admin", "Editor", "Author"))
				
				// Supports Search & Pagination via Query Parameters
				r.Get("/", blogHandler.GetAll)
				r.Get("/{id}", blogHandler.GetByID)
				r.Post("/", blogHandler.Create)
				r.Put("/{id}", blogHandler.Update)
				r.Delete("/{id}", blogHandler.Delete)
			})

			// User Management (Super Admin only)
			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.RequireRole("Super Admin"))
				
				r.Get("/", authHandler.ListUsers)
				r.Post("/", authHandler.CreateUser)
				r.Get("/deleted", authHandler.ListDeletedUsers)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", authHandler.GetUserByID)
					r.Put("/", authHandler.UpdateUser)
					r.Delete("/", authHandler.DeleteUser)
					r.Post("/toggle-status", authHandler.ToggleUserStatus)
					r.Post("/restore", authHandler.RestoreUser)
				})
			})

			// RBAC / Roles & Permissions Management (Super Admin & Admin)
			r.Route("/rbac", func(r chi.Router) {
				r.Use(middleware.RequireRole("Super Admin", "Admin"))

				// Roles
				r.Route("/roles", func(r chi.Router) {
					r.Post("/", rbacHandler.CreateRole)
					r.Get("/", rbacHandler.GetAllRoles)

					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", rbacHandler.GetRoleByID)
						r.Put("/", rbacHandler.UpdateRole)
						r.Delete("/", rbacHandler.DeleteRole)

						// Role permissions
						r.Post("/permissions", rbacHandler.AssignPermissionsToRole)
						r.Get("/permissions", rbacHandler.GetRolePermissions)

						// Module access
						r.Post("/module-access", rbacHandler.UpdateModuleAccess)
						r.Get("/module-access", rbacHandler.GetModuleAccessByRole)
					})
				})

				// Permissions
				r.Route("/permissions", func(r chi.Router) {
					r.Post("/", rbacHandler.CreatePermission)
					r.Get("/", rbacHandler.GetAllPermissions)
					r.Get("/by-module", rbacHandler.GetPermissionsByModule)

					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", rbacHandler.GetPermissionByID)
						r.Put("/", rbacHandler.UpdatePermission)
						r.Delete("/", rbacHandler.DeletePermission)
					})
				})

				// User roles mapping
				r.Route("/users/{userId}/roles", func(r chi.Router) {
					r.Post("/", rbacHandler.AssignRolesToUser)
					r.Get("/", rbacHandler.GetUserRoles)
				})

				// Permission checking endpoints
				r.Group(func(r chi.Router) {
					r.Post("/check-permission", rbacHandler.CheckPermission)
					r.Post("/check-module-access", rbacHandler.CheckModuleAccess)
				})
			})
		})
	})

	return mux
}

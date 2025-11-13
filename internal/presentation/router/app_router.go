package router

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	authadmincontroller "mqfm_backend/internal/handler/auth/admin"
	authadminusecase "mqfm_backend/internal/domain/usecases/auth/admin"
	authadminservice "mqfm_backend/internal/application/services/auth/admin"
	authadminrepository "mqfm_backend/internal/infrastructure/repository/auth/admin"
	authusercontroller "mqfm_backend/internal/handler/auth/user"
	authuserusecase "mqfm_backend/internal/domain/usecases/auth/user"
	authuserservice "mqfm_backend/internal/application/services/auth/user"
	authuserrepository "mqfm_backend/internal/infrastructure/repository/auth/user"
	categoriesadmincontroller "mqfm_backend/internal/handler/podcast/categories/admin"
	categoriesadminusecase "mqfm_backend/internal/domain/usecases/podcast/categories/admin"
	categoriesadminservice "mqfm_backend/internal/application/services/podcast/categories/admin"
	categoriesadminrepository "mqfm_backend/internal/infrastructure/repository/podcast/categories/admin"
	categoriesusercontroller "mqfm_backend/internal/handler/podcast/categories/user"
	categoriesuserservice "mqfm_backend/internal/application/services/podcast/categories/user"
	categoriesuserrepository "mqfm_backend/internal/infrastructure/repository/podcast/categories/user"
	middleware "mqfm_backend/internal/middleware"
	"mqfm_backend/internal/infrastructure/database"

)

func SetupAppRouter() http.Handler {
	if err := godotenv.Load("env_staging.env"); err != nil {
		log.Fatalf("Gagal load env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := database.NewMySQLConnection(dbUser, dbPass, dbHost, dbName)
	if err != nil {
		log.Fatalf("Gagal konek database: %v", err)
	}

	adminRepo := authadminrepository.NewMySQLAuthAdminRepository(db)
	adminService := authadminservice.NewAuthAdminService(adminRepo, jwtSecret)
	adminUsecase := authadminusecase.NewAuthAdminUsecase(adminService)
	adminController := authadmincontroller.NewAuthAdminController(adminUsecase, adminRepo)

	userRepo := authuserrepository.NewMySQLAuthUserRepository(db)
	userService := authuserservice.NewAuthUserService(userRepo, jwtSecret)
	userUsecase := authuserusecase.NewAuthUserUsecase(userService)
	userController := authusercontroller.NewAuthUserController(userUsecase, userRepo)

	catRepo := categoriesadminrepository.NewMySQLCategoriesAdminRepository(db)
	catAdminService := categoriesadminservice.NewCategoriesAdminService(catRepo)
	catAdminUsecase := categoriesadminusecase.NewCategoriesAdminUsecase(catAdminService)
	catAdminController := categoriesadmincontroller.NewCategoriesAdminController(catAdminUsecase)

	catUserRepo := categoriesuserrepository.NewMySQLCategoriesUserRepository(db)
	catUserService := categoriesuserservice.NewCategoriesUserService(catUserRepo)
	catUserController := categoriesusercontroller.NewCategoriesUserController(catUserService)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MQFM Backend API (staging) is running"))
	})

	mux.Handle("/uploads/", http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("./storage/uploads"))))

	mux.HandleFunc("/api/v1/auth/admin/register", adminController.Register)
	mux.HandleFunc("/api/v1/auth/admin/login", adminController.Login)
	mux.HandleFunc("/api/v1/auth/admin/me", adminController.Me)
	mux.HandleFunc("/api/v1/auth/admin/update", adminController.Update)
	mux.HandleFunc("/api/v1/auth/admin/delete", adminController.DeleteAccount)
	mux.HandleFunc("/api/v1/auth/admin/logout", adminController.Logout)

	mux.HandleFunc("/api/v1/auth/user/register", userController.Register)
	mux.HandleFunc("/api/v1/auth/user/login", userController.Login)
	mux.HandleFunc("/api/v1/auth/user/me", userController.Me)
	mux.HandleFunc("/api/v1/auth/user/update", userController.Update)
	mux.HandleFunc("/api/v1/auth/user/delete", userController.DeleteAccount)
	mux.HandleFunc("/api/v1/auth/user/logout", userController.Logout)

	mux.HandleFunc("/api/v1/podcast/admin/categories/create",
		middleware.AdminAuthMiddleware(adminService, catAdminController.Create))
	mux.HandleFunc("/api/v1/podcast/admin/categories/update",
		middleware.AdminAuthMiddleware(adminService, catAdminController.Update))
	mux.HandleFunc("/api/v1/podcast/admin/categories/delete",
		middleware.AdminAuthMiddleware(adminService, catAdminController.Delete))
	mux.HandleFunc("/api/v1/podcast/admin/categories/all",
		middleware.AdminAuthMiddleware(adminService, catAdminController.GetAll))
	mux.HandleFunc("/api/v1/podcast/admin/categories/detail",
		middleware.AdminAuthMiddleware(adminService, catAdminController.GetByID))

	mux.HandleFunc("/api/v1/podcast/user/categories/all",
		middleware.UserAuthMiddleware(userService, catUserController.GetAll))
	mux.HandleFunc("/api/v1/podcast/user/categories/detail",
		middleware.UserAuthMiddleware(userService, catUserController.GetByID))

	return middleware.WrapWithCORS(mux)
}

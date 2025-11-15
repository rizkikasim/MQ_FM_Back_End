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
	audioadmincontroller "mqfm_backend/internal/handler/podcast/audio/admin"
audioadminusecase "mqfm_backend/internal/domain/usecases/podcast/audio/admin"
audioadminservice "mqfm_backend/internal/application/services/podcast/audio/admin"
audioadminrepository "mqfm_backend/internal/infrastructure/repository/podcast/audio/admin"
audiousercontroller "mqfm_backend/internal/handler/podcast/audio/user"
audiouserusecase "mqfm_backend/internal/domain/usecases/podcast/audio/user"
audiouserservice "mqfm_backend/internal/application/services/podcast/audio/user"
audiouserrepository "mqfm_backend/internal/infrastructure/repository/podcast/audio/user"
	middleware "mqfm_backend/internal/middleware"
	"mqfm_backend/internal/infrastructure/database"

)

func SetupAppRouter() http.Handler {

	// Load env
	if err := godotenv.Load("env_staging.env"); err != nil {
		log.Fatalf("Gagal load env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Database
	db, err := database.NewMySQLConnection(dbUser, dbPass, dbHost, dbName)
	if err != nil {
		log.Fatalf("Gagal konek database: %v", err)
	}

	// Admin
	adminRepo := authadminrepository.NewMySQLAuthAdminRepository(db)
	adminService := authadminservice.NewAuthAdminService(adminRepo, jwtSecret)
	adminUsecase := authadminusecase.NewAuthAdminUsecase(adminService)
	adminController := authadmincontroller.NewAuthAdminController(adminUsecase, adminRepo)

	// User
	userRepo := authuserrepository.NewMySQLAuthUserRepository(db)
	userService := authuserservice.NewAuthUserService(userRepo, jwtSecret)
	userUsecase := authuserusecase.NewAuthUserUsecase(userService)
	userController := authusercontroller.NewAuthUserController(userUsecase, userRepo)

	// Categories Admin
	catRepo := categoriesadminrepository.NewMySQLCategoriesAdminRepository(db)
	catAdminService := categoriesadminservice.NewCategoriesAdminService(catRepo)
	catAdminUsecase := categoriesadminusecase.NewCategoriesAdminUsecase(catAdminService)
	catAdminController := categoriesadmincontroller.NewCategoriesAdminController(catAdminUsecase)

	// Categories User
	catUserRepo := categoriesuserrepository.NewMySQLCategoriesUserRepository(db)
	catUserService := categoriesuserservice.NewCategoriesUserService(catUserRepo)
	catUserController := categoriesusercontroller.NewCategoriesUserController(catUserService)

	// Audio Podcast Admin
audioRepo := audioadminrepository.NewMySQLAudioAdminPodcastRepository(db)
audioService := audioadminservice.NewAudioAdminPodcastService(audioRepo)
audioUsecase := audioadminusecase.NewAudioAdminPodcastUsecase(audioService)
audioController := audioadmincontroller.NewAudioAdminPodcastController(audioUsecase)

// Audio Podcast User
audioUserRepo := audiouserrepository.NewMySQLAudioUserPodcastRepository(db)
audioUserService := audiouserservice.NewAudioUserPodcastService(audioUserRepo)
audioUserUsecase := audiouserusecase.NewAudioUserPodcastUsecase(audioUserService)
audioUserController := audiousercontroller.NewAudioUserPodcastController(audioUserUsecase)



	// Router
	mux := http.NewServeMux()

	// Root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Write([]byte("MQFM Backend API (staging) is running"))
	})

	// Static uploads
	mux.Handle("/uploads/", http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("./storage/uploads"))))

	// =====================================
	// 🔥 ADMIN AUTH (PUBLIC)
	// =====================================
	mux.Handle("/api/v1/auth/admin/register",
		middleware.CORSMiddleware(http.HandlerFunc(adminController.Register)),
	)
	mux.Handle("/api/v1/auth/admin/login",
		middleware.CORSMiddleware(http.HandlerFunc(adminController.Login)),
	)
	mux.Handle("/api/v1/auth/admin/me",
		middleware.CORSMiddleware(http.HandlerFunc(adminController.Me)),
	)
	mux.Handle("/api/v1/auth/admin/update",
		middleware.CORSMiddleware(http.HandlerFunc(adminController.Update)),
	)
	mux.Handle("/api/v1/auth/admin/delete",
		middleware.CORSMiddleware(http.HandlerFunc(adminController.DeleteAccount)),
	)
	mux.Handle("/api/v1/auth/admin/logout",
		middleware.CORSMiddleware(http.HandlerFunc(adminController.Logout)),
	)

	// =====================================
	// 🔥 USER AUTH (PUBLIC)
	// =====================================
	mux.Handle("/api/v1/auth/user/register",
		middleware.UserCORSMiddleware(http.HandlerFunc(userController.Register)),
	)
	mux.Handle("/api/v1/auth/user/login",
		middleware.UserCORSMiddleware(http.HandlerFunc(userController.Login)),
	)
	mux.Handle("/api/v1/auth/user/me",
		middleware.UserCORSMiddleware(http.HandlerFunc(userController.Me)),
	)
	mux.Handle("/api/v1/auth/user/update",
		middleware.UserCORSMiddleware(http.HandlerFunc(userController.Update)),
	)
	mux.Handle("/api/v1/auth/user/delete",
		middleware.UserCORSMiddleware(http.HandlerFunc(userController.DeleteAccount)),
	)
	mux.Handle("/api/v1/auth/user/logout",
		middleware.UserCORSMiddleware(http.HandlerFunc(userController.Logout)),
	)

	// =====================================
	// 🔥 ADMIN CATEGORIES (PROTECTED + CORS)
	// =====================================
	mux.Handle("/api/v1/podcast/admin/categories/create",
		middleware.CORSMiddleware(
			http.HandlerFunc(
				middleware.AdminAuthMiddleware(adminService, catAdminController.Create),
			),
		),
	)
	mux.Handle("/api/v1/podcast/admin/categories/update",
		middleware.CORSMiddleware(
			http.HandlerFunc(
				middleware.AdminAuthMiddleware(adminService, catAdminController.Update),
			),
		),
	)
	mux.Handle("/api/v1/podcast/admin/categories/delete",
		middleware.CORSMiddleware(
			http.HandlerFunc(
				middleware.AdminAuthMiddleware(adminService, catAdminController.Delete),
			),
		),
	)
	mux.Handle("/api/v1/podcast/admin/categories/all",
		middleware.CORSMiddleware(
			http.HandlerFunc(
				middleware.AdminAuthMiddleware(adminService, catAdminController.GetAll),
			),
		),
	)
	mux.Handle("/api/v1/podcast/admin/categories/detail",
		middleware.CORSMiddleware(
			http.HandlerFunc(
				middleware.AdminAuthMiddleware(adminService, catAdminController.GetByID),
			),
		),
	)

	// =====================================
	// 🔥 USER CATEGORIES (Protected + CORS)
	// =====================================
	mux.Handle("/api/v1/podcast/user/categories/all",
		middleware.UserCORSMiddleware(
			http.HandlerFunc(
				middleware.UserAuthMiddleware(userService, catUserController.GetAll),
			),
		),
	)

	mux.Handle("/api/v1/podcast/user/categories/detail",
		middleware.UserCORSMiddleware(
			http.HandlerFunc(
				middleware.UserAuthMiddleware(userService, catUserController.GetByID),
			),
		),
	)

	// =====================================
// 🔥 ADMIN AUDIO PODCAST (PROTECTED + CORS)
// =====================================

mux.Handle("/api/v1/podcast/admin/audio/create",
    middleware.CORSMiddleware(
        http.HandlerFunc(
            middleware.AdminAuthMiddleware(adminService, audioController.Create),
        ),
    ),
)

mux.Handle("/api/v1/podcast/admin/audio/update",
    middleware.CORSMiddleware(
        http.HandlerFunc(
            middleware.AdminAuthMiddleware(adminService, audioController.Update),
        ),
    ),
)

mux.Handle("/api/v1/podcast/admin/audio/delete",
    middleware.CORSMiddleware(
        http.HandlerFunc(
            middleware.AdminAuthMiddleware(adminService, audioController.Delete),
        ),
    ),
)

mux.Handle("/api/v1/podcast/admin/audio/all",
    middleware.CORSMiddleware(
        http.HandlerFunc(
            middleware.AdminAuthMiddleware(adminService, audioController.GetAll),
        ),
    ),
)

mux.Handle("/api/v1/podcast/admin/audio/detail",
    middleware.CORSMiddleware(
        http.HandlerFunc(
            middleware.AdminAuthMiddleware(adminService, audioController.GetByID),
        ),
    ),
)

// =====================================
// 🔥 USER AUDIO PODCAST (Protected + CORS)
// =====================================

mux.Handle("/api/v1/podcast/user/audio/all",
    middleware.UserCORSMiddleware(
        http.HandlerFunc(
            middleware.UserAuthMiddleware(userService, audioUserController.GetAll),
        ),
    ),
)

mux.Handle("/api/v1/podcast/user/audio/detail",
    middleware.UserCORSMiddleware(
        http.HandlerFunc(
            middleware.UserAuthMiddleware(userService, audioUserController.GetByID),
        ),
    ),
)



	return mux
}

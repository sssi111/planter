package api

import (
	"net/http"

	"github.com/anpanovv/planter/internal/config"
	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// API represents the API server
type API struct {
	router          *mux.Router
	authService     *services.AuthService
	userService     *services.UserService
	plantService    *services.PlantService
	shopService     *services.ShopService
	recommendationService *services.RecommendationService
	notificationService *services.NotificationService
	auth            *middleware.Auth
}

// New creates a new API server
func New(
	authService *services.AuthService,
	userService *services.UserService,
	plantService *services.PlantService,
	shopService *services.ShopService,
	recommendationService *services.RecommendationService,
	notificationService *services.NotificationService,
	auth *middleware.Auth,
) *API {
	api := &API{
		router:          mux.NewRouter(),
		authService:     authService,
		userService:     userService,
		plantService:    plantService,
		shopService:     shopService,
		recommendationService: recommendationService,
		notificationService: notificationService,
		auth:            auth,
	}

	api.setupRoutes()
	return api
}

// setupRoutes sets up the API routes
func (a *API) setupRoutes() {
	// Auth routes
	a.router.HandleFunc("/auth/login", a.handleLogin).Methods(http.MethodPost)
	a.router.HandleFunc("/auth/register", a.handleRegister).Methods(http.MethodPost)

	// User routes
	userRouter := a.router.PathPrefix("/users").Subrouter()
	userRouter.Use(a.auth.RequireAuth)
	userRouter.HandleFunc("/{userId}", a.handleGetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/{userId}", a.handleUpdateUser).Methods(http.MethodPut)

	// Plant routes
	a.router.HandleFunc("/plants", a.handleGetAllPlants).Methods(http.MethodGet)
	a.router.HandleFunc("/plants/search", a.handleSearchPlants).Methods(http.MethodGet)
	a.router.HandleFunc("/plants/{plantId}", a.handleGetPlant).Methods(http.MethodGet)

	// Plant routes that require authentication
	plantRouter := a.router.PathPrefix("/plants").Subrouter()
	plantRouter.Use(a.auth.RequireAuth)
	userRouter.HandleFunc("/me/favorites", a.handleGetFavoritePlants).Methods(http.MethodGet)
	plantRouter.HandleFunc("/{plantId}/favorite", a.handleAddToFavorites).Methods(http.MethodPost)
	plantRouter.HandleFunc("/{plantId}/favorite", a.handleRemoveFromFavorites).Methods(http.MethodDelete)
	plantRouter.HandleFunc("/{plantId}/water", a.handleMarkAsWatered).Methods(http.MethodPost)
	plantRouter.HandleFunc("/user", a.handleGetUserPlants).Methods(http.MethodGet)
	plantRouter.HandleFunc("/user/{plantId}", a.handleAddUserPlant).Methods(http.MethodPost)
	plantRouter.HandleFunc("/user/{plantId}", a.handleUpdateUserPlant).Methods(http.MethodPut)
	plantRouter.HandleFunc("/user/{plantId}", a.handleRemoveUserPlant).Methods(http.MethodDelete)

	// Shop routes
	a.router.HandleFunc("/shops", a.handleGetAllShops).Methods(http.MethodGet)
	a.router.HandleFunc("/shops/{shopId}", a.handleGetShop).Methods(http.MethodGet)
	a.router.HandleFunc("/shops/{shopId}/plants", a.handleGetShopPlants).Methods(http.MethodGet)

	// Recommendation routes
	recommendationRouter := a.router.PathPrefix("/recommendations").Subrouter()
	recommendationRouter.HandleFunc("/questionnaire", a.handleSaveQuestionnaire).Methods(http.MethodPost)
	recommendationRouter.HandleFunc("/questionnaire/detailed", a.handleSaveDetailedQuestionnaire).Methods(http.MethodPost)
	recommendationRouter.HandleFunc("/questionnaire/{questionnaireId}", a.handleGetRecommendations).Methods(http.MethodGet)
	
	// Admin routes
	adminRouter := a.router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/plants", a.handleAdminCreatePlant).Methods(http.MethodPost)
	
	// Chat routes (require authentication)
	chatRouter := a.router.PathPrefix("/chat").Subrouter()
	chatRouter.Use(a.auth.RequireAuth)
	chatRouter.HandleFunc("/sessions", a.handleCreateChatSession).Methods(http.MethodPost)
	chatRouter.HandleFunc("/sessions", a.handleGetChatSessions).Methods(http.MethodGet)
	chatRouter.HandleFunc("/sessions/{sessionId}", a.handleGetChatSession).Methods(http.MethodGet)
	chatRouter.HandleFunc("/sessions/{sessionId}/messages", a.handleGetChatMessages).Methods(http.MethodGet)
	chatRouter.HandleFunc("/sessions/{sessionId}/messages", a.handleSendChatMessage).Methods(http.MethodPost)

	// Notification routes
	a.router.Handle("/notifications", a.auth.RequireAuth(http.HandlerFunc(a.handleGetUserNotifications))).Methods(http.MethodGet)
	a.router.Handle("/notifications/{notificationId}/read", a.auth.RequireAuth(http.HandlerFunc(a.handleMarkNotificationAsRead))).Methods(http.MethodPost)
}

// Handler returns the HTTP handler for the API
func (a *API) Handler() http.Handler {
	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap router with logging middleware and CORS
	return c.Handler(middleware.LoggingMiddleware(a.router))
}

// Start starts the API server
func (a *API) Start(cfg *config.Config) error {
	addr := ":" + cfg.Server.Port
	return http.ListenAndServe(addr, a.Handler())
}
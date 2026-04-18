package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	authpb "github.com/fercho/school-tracking/proto/gen/auth/v1"
	_ "github.com/fercho/school-tracking/services/gateway/docs/api"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/handlers"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/middlewares"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	cfg *env.Config,
	log *zap.Logger,
	authClient authpb.AuthServiceClient,
	fleetHandler *handlers.FleetHandler,
	tripHandler *handlers.TripHandler,
	notificationHandler *handlers.NotificationHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("gateway is healthy"))
		})

		// Public Auth routes proxy (Since gRPC doesn't have login/register yet)
		authProxy, err := handlers.NewReverseProxy(cfg.AuthServiceURL, "")
		if err == nil {
			r.Mount("/auth", authProxy)
		} else {
			log.Error("failed to create auth service proxy", zap.Error(err))
		}

		// Protected Fleet routes
		r.Route("/fleet", func(r chi.Router) {
			r.Use(middlewares.RequireAuth(cfg, authClient))

			// Vehicles
			r.Post("/vehicles", fleetHandler.CreateVehicle)
			r.Get("/vehicles", fleetHandler.ListVehicles)
			r.Get("/vehicles/{id}", fleetHandler.GetVehicle)

			// Schools
			r.Post("/schools", fleetHandler.CreateSchool)
			r.Get("/schools", fleetHandler.ListSchools)
			r.Get("/schools/{id}", fleetHandler.GetSchool)
			r.Put("/schools/{id}", fleetHandler.UpdateSchool)

			// Drivers
			r.Post("/drivers", fleetHandler.RegisterDriver)
			r.Get("/drivers", fleetHandler.ListDrivers)
			r.Get("/drivers/{id}", fleetHandler.GetDriver)
			r.Put("/drivers/{id}", fleetHandler.UpdateDriver)

			// Students
			r.Post("/students", fleetHandler.RegisterStudent)
			r.Get("/students", fleetHandler.ListStudents)
			r.Get("/students/{id}", fleetHandler.GetStudent)
			r.Put("/students/{id}", fleetHandler.UpdateStudent)
			r.Delete("/students/{id}", fleetHandler.DeactivateStudent)
			r.Get("/students/{student_id}/guardians", fleetHandler.GetGuardiansByStudent)

			// Guardians
			r.Post("/guardians", fleetHandler.LinkGuardian)
			r.Delete("/guardians/{id}", fleetHandler.UnlinkGuardian)

			// Routes
			r.Post("/routes", fleetHandler.CreateRoute)
			r.Get("/routes", fleetHandler.ListRoutes)
			r.Get("/routes/{id}", fleetHandler.GetRoute)
			r.Put("/routes/{id}", fleetHandler.UpdateRoute)
			r.Post("/routes/{id}/stops", fleetHandler.AddStop)
			r.Get("/routes/{id}/stops", fleetHandler.GetRouteStops)
			r.Delete("/routes/{id}/stops/{stop_id}", fleetHandler.RemoveStop)
		})

		// Protected Trip routes
		r.Route("/trips", func(r chi.Router) {
			r.Use(middlewares.RequireAuth(cfg, authClient))

			r.Post("/", tripHandler.StartTrip)
			r.Put("/{id}/end", tripHandler.EndTrip)
			r.Get("/{id}", tripHandler.GetTrip)
			r.Get("/", tripHandler.ListTrips)
			r.Get("/active", tripHandler.ListActiveTrips)

			r.Post("/{id}/checkins", tripHandler.CheckinStudent)
			r.Get("/{id}/checkins", tripHandler.GetTripCheckins)
		})

		// Protected Notification routes
		r.Route("/notifications", func(r chi.Router) {
			r.Use(middlewares.RequireAuth(cfg, authClient))

			r.Post("/push", notificationHandler.SendPush)
			r.Post("/sms", notificationHandler.SendSMS)
			r.Get("/{id}", notificationHandler.GetNotification)
			r.Get("/", notificationHandler.ListNotifications)
			r.Post("/retry", notificationHandler.RetryFailed)
		})
	})

	return r
}

func StartHTTPServer(lc fx.Lifecycle, r *chi.Mux, cfg *env.Config, log *zap.Logger) {
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server", zap.String("port", cfg.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})
}

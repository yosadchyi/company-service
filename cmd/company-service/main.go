// @title Company Service API
// @version 1.0
// @description This service allows CRUD operations on companies
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v2"

	_ "company-service/docs"

	"company-service/internal/auth"
	"company-service/internal/company"
	"company-service/internal/middleware"
	"company-service/internal/postgresrepo"
	"company-service/pkg/jwt"
	"company-service/pkg/kafka"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
	} `yaml:"kafka"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

// runMigrationsWithMigrate uses golang-migrate to run DB migrations from /migrations
func runMigrationsWithMigrate(dbURL string) error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	cfgBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("failed to read config:", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		log.Fatal("failed to parse config:", err)
	}

	// DSN for pgx and migrate
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName)

	// Run DB migrations
	if err := runMigrationsWithMigrate(dsn); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Connect to DB
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer dbpool.Close()

	// Init Kafka + JWT middleware
	kafka.InitKafka(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	middleware.SetSecret([]byte(cfg.JWT.Secret))
	jwt.SetSecret([]byte(cfg.JWT.Secret))

	// DI wiring
	companyRepo := postgresrepo.NewCompanyRepo(dbpool)
	companySvc := company.NewService(companyRepo)
	companyHandler := company.NewHandler(companySvc)

	authRepo := postgresrepo.NewUserRepo(dbpool)
	authSvc := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authSvc)

	// Routes
	r := chi.NewRouter()

	// Public routes
	r.Group(func(p chi.Router) {
		p.Get("/swagger/*", httpSwagger.WrapHandler)
		p.Post("/auth/login", authHandler.Login)
		p.Post("/auth/refresh", authHandler.Refresh)
	})

	// Secured routes
	r.Group(func(priv chi.Router) {
		priv.Use(middleware.JWTMiddleware)
		companyHandler.RegisterRoutes(priv)
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

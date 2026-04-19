package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"dashboard-cs-be/config"
	deliveryHTTP "dashboard-cs-be/delivery/http"
	"dashboard-cs-be/repository"
	"dashboard-cs-be/usecase"
)

func main() {
	cfg := config.Load()

	// Connect to MySQL 
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("db.Ping: %v — check DB_HOST, DB_USER, DB_PASSWORD, DB_NAME", err)
	}
	log.Printf("Connected to MySQL  →  %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Dependency Injection

	// Dashboard
	dashRepo    := repository.NewMySQLDashboardRepository(db)
	dashUC      := usecase.NewDashboardUsecase(dashRepo)
	dashHandler := deliveryHTTP.NewDashboardHandler(dashUC)

	// Import
	importRepo    := repository.NewMySQLImportRepository(db)
	importUC      := usecase.NewImportUsecase(importRepo)
	importHandler := deliveryHTTP.NewImportHandler(importUC)

	// Export
	exportRepo    := repository.NewMySQLExportRepository(db)
	exportUC      := usecase.NewExportUsecase(exportRepo)
	exportHandler := deliveryHTTP.NewExportHandler(exportUC)

	// Incidents
	incidentRepo    := repository.NewMySQLIncidentRepository(db)
	incidentUC      := usecase.NewIncidentUsecase(incidentRepo, db)
	incidentHandler := deliveryHTTP.NewIncidentHandler(incidentUC)

	// Router
	router := deliveryHTTP.NewRouter(dashHandler, importHandler, exportHandler, incidentHandler)

	// Start server
	addr := ":" + cfg.AppPort
	log.Printf("Dashboard API  →  http://localhost%s", addr)
	log.Println("──────────────────────────────────────────────────────────────────")
	log.Printf("  GET    /health")
	log.Printf("  GET    /api/v1/dashboard?from=YYYY-MM-DD&to=YYYY-MM-DD")
	log.Printf("  GET    /api/v1/realtime")
	log.Printf("  POST   /api/v1/import                        (multipart, field: file)")
	log.Printf("  GET    /api/v1/export?from=&to=&channel=     (all|whatsapp|email,...)")
	log.Printf("  GET    /api/v1/incidents/active")
	log.Printf("  GET    /api/v1/incidents/history?from=&to=")
	log.Printf("  POST   /api/v1/incidents")
	log.Printf("  PATCH  /api/v1/incidents/{id}/resolve")
	log.Println("──────────────────────────────────────────────────────────────────")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
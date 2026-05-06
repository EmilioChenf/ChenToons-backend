package main

import (
	"log"
	"os"

	"chentoons-backend/internal/config"
	"chentoons-backend/internal/database"
	"chentoons-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadCheninConfig()

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	db, err := database.ConectarPostgresChen(cfg)
	if err != nil {
		log.Fatal("no se pudo conectar a postgres: ", err)
	}
	defer db.Close()

	if err := database.MigrarTablasChenin(db); err != nil {
		log.Fatal("fallo la migracion: ", err)
	}
	if err := database.SembrarDatosChenin(db); err != nil {
		log.Fatal("fallo el seed: ", err)
	}

	app := fiber.New(fiber.Config{
		AppName:   "ChenToons API",
		BodyLimit: 2 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.RegistrarRutasChen(app, db, cfg.UploadDir)

	log.Printf("ChenToons escuchando en 0.0.0.0:%s (%s)", cfg.AppPort, cfg.AppEnv)
	log.Fatal(app.Listen("0.0.0.0:" + cfg.AppPort))
}

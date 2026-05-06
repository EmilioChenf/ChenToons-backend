package routes

import (
	"chentoons-backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegistrarRutasChen(app *fiber.App, db *pgxpool.Pool, uploadDir string) {
	series := handlers.SeriesHandler{DB: db}
	personajes := handlers.PersonajesHandler{DB: db}
	episodios := handlers.EpisodiosHandler{DB: db}
	ratings := handlers.RatingsHandler{DB: db}
	uploads := handlers.UploadsHandler{Folder: uploadDir}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"project": "ChenToons",
		})
	})

	app.Get("/series", series.CargarSeriesChenin)
	app.Get("/series/:id", series.ObtenerSerieChen)
	app.Post("/series", series.CrearSerieChen)
	app.Put("/series/:id", series.ActualizarSerieChen)
	app.Delete("/series/:id", series.EliminarSerieChen)

	app.Get("/personajes", personajes.CargarPersonajesChen)
	app.Get("/personajes/:id", personajes.ObtenerPersonajeChen)
	app.Post("/personajes", personajes.CrearPersonajeChen)
	app.Put("/personajes/:id", personajes.ActualizarPersonajeJosuc)
	app.Delete("/personajes/:id", personajes.EliminarPersonajeChen)

	app.Get("/episodios", episodios.CargarEpisodiosEmilio)
	app.Get("/episodios/:id", episodios.ObtenerEpisodioChen)
	app.Post("/episodios", episodios.CrearEpisodioChen)
	app.Put("/episodios/:id", episodios.ActualizarEpisodioChen)
	app.Delete("/episodios/:id", episodios.EliminarEpisodioChen)

	app.Get("/series/:id/ratings", ratings.CargarRatingsSerieChen)
	app.Post("/series/:id/ratings", ratings.CrearRatingChen)
	app.Get("/series/:id/promedio-rating", ratings.PromedioRatingToon)

	app.Post("/uploads", uploads.SubirImagenChenin)
	app.Static("/uploads", uploadDir)

	app.Get("/export/series.csv", series.ExportarCSVChenin)

	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.html")
	})
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/docs")
	})
	app.Get("/swagger/index.html", func(c *fiber.Ctx) error {
		return c.Redirect("/docs")
	})
	app.Get("/docs/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/openapi.yaml")
	})
}

package handlers

import (
	"context"
	"strconv"
	"strings"

	"chentoons-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatingsHandler struct {
	DB *pgxpool.Pool
}

func (h RatingsHandler) CargarRatingsSerieChen(c *fiber.Ctx) error {
	serieID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	rows, err := h.DB.Query(context.Background(), `SELECT id, serie_id, puntuacion, comentario,
		created_at FROM ratings WHERE serie_id=$1 ORDER BY created_at DESC`, serieID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	ratings := []models.Rating{}
	for rows.Next() {
		var r models.Rating
		if err := rows.Scan(&r.ID, &r.SerieID, &r.Puntuacion, &r.Comentario, &r.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		ratings = append(ratings, r)
	}

	return c.JSON(ratings)
}

func (h RatingsHandler) CrearRatingChen(c *fiber.Ctx) error {
	serieID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var r models.Rating
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if r.Puntuacion < 1 || r.Puntuacion > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "la puntuacion debe estar entre 1 y 5"})
	}
	r.SerieID = serieID
	r.Comentario = strings.TrimSpace(r.Comentario)

	err = h.DB.QueryRow(context.Background(), `INSERT INTO ratings (serie_id, puntuacion, comentario)
		VALUES ($1,$2,$3) RETURNING id, created_at`,
		r.SerieID, r.Puntuacion, r.Comentario,
	).Scan(&r.ID, &r.CreatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h RatingsHandler) PromedioRatingToon(c *fiber.Ctx) error {
	serieID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var promedio float64
	var total int
	err = h.DB.QueryRow(context.Background(), `SELECT COALESCE(AVG(puntuacion), 0), COUNT(*)
		FROM ratings WHERE serie_id=$1`, serieID).Scan(&promedio, &total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"serie_id": serieID,
		"promedio": promedio,
		"total":    total,
		"maximo":   5,
	})
}

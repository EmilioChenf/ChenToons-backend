package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"

	"chentoons-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatingsHandler struct {
	DB *pgxpool.Pool
}

type ratingInputChen struct {
	Puntuacion int    `json:"puntuacion"`
	Rating     int    `json:"rating"`
	Estrellas  int    `json:"estrellas"`
	Comentario string `json:"comentario"`
}

func (h RatingsHandler) CargarRatingsSerieChen(c *fiber.Ctx) error {
	serieID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	rows, err := h.DB.Query(context.Background(), `SELECT id, serie_id, puntuacion, comentario,
		created_at FROM ratings WHERE serie_id=$1 ORDER BY created_at DESC, id DESC`, serieID)
	if err != nil {
		log.Printf("chen ratings GET fallo serie_id=%d error=%v", serieID, err)
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
	if err := rows.Err(); err != nil {
		log.Printf("chen ratings rows fallo serie_id=%d error=%v", serieID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("chen ratings GET serie_id=%d total=%d", serieID, len(ratings))

	return c.JSON(ratings)
}

func (h RatingsHandler) CrearRatingChen(c *fiber.Ctx) error {
	serieID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var input ratingInputChen
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}

	puntuacion := input.Puntuacion
	if puntuacion == 0 {
		puntuacion = input.Rating
	}
	if puntuacion == 0 {
		puntuacion = input.Estrellas
	}

	r := models.Rating{
		SerieID:    serieID,
		Puntuacion: puntuacion,
		Comentario: strings.TrimSpace(input.Comentario),
	}

	if r.Puntuacion < 1 || r.Puntuacion > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "la puntuacion debe estar entre 1 y 5"})
	}

	var existeSerie bool
	err = h.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM series WHERE id=$1)", serieID).Scan(&existeSerie)
	if err != nil {
		log.Printf("chen rating validacion serie fallo serie_id=%d error=%v", serieID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !existeSerie {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "serie no encontrada"})
	}

	err = h.DB.QueryRow(context.Background(), `INSERT INTO ratings (serie_id, puntuacion, comentario)
		VALUES ($1,$2,$3)
		RETURNING id, serie_id, puntuacion, comentario, created_at`,
		r.SerieID, r.Puntuacion, r.Comentario,
	).Scan(&r.ID, &r.SerieID, &r.Puntuacion, &r.Comentario, &r.CreatedAt)
	if err != nil {
		log.Printf("chen rating INSERT fallo serie_id=%d puntuacion=%d error=%v", serieID, r.Puntuacion, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("chen rating creado id=%d serie_id=%d puntuacion=%d", r.ID, r.SerieID, r.Puntuacion)

	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h RatingsHandler) PromedioRatingToon(c *fiber.Ctx) error {
	serieID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var promedio float64
	var total int
	err = h.DB.QueryRow(context.Background(), `SELECT COALESCE(AVG(puntuacion)::float8, 0), COUNT(*)
		FROM ratings WHERE serie_id=$1`, serieID).Scan(&promedio, &total)
	if err != nil {
		log.Printf("chen promedio rating fallo serie_id=%d error=%v", serieID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("chen promedio rating serie_id=%d promedio=%.2f total=%d", serieID, promedio, total)

	return c.JSON(fiber.Map{
		"serie_id": serieID,
		"promedio": promedio,
		"total":    total,
		"maximo":   5,
	})
}

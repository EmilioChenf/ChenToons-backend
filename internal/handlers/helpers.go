package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func errorChen(c *fiber.Ctx, status int, mensaje string) error {
	return c.Status(status).JSON(fiber.Map{"error": mensaje})
}

func existeSerieChenin(db *pgxpool.Pool, serieID int) (bool, error) {
	var existe bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM series WHERE id=$1)", serieID).Scan(&existe)
	return existe, err
}

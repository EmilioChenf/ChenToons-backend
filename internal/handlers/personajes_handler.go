package handlers

import (
	"context"
	"strconv"
	"strings"

	"chentoons-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonajesHandler struct {
	DB *pgxpool.Pool
}

func (h PersonajesHandler) CargarPersonajesChen(c *fiber.Ctx) error {
	query := `SELECT id, serie_id, nombre, descripcion, rol, personalidad, imagen, created_at
		FROM personajes`
	args := []any{}
	if c.Query("serie_id") != "" {
		serieID, err := strconv.Atoi(c.Query("serie_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "serie_id invalido"})
		}
		query += " WHERE serie_id=$1"
		args = append(args, serieID)
	}
	query += " ORDER BY id DESC"

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	personajes := []models.Personaje{}
	for rows.Next() {
		var p models.Personaje
		if err := rows.Scan(&p.ID, &p.SerieID, &p.Nombre, &p.Descripcion, &p.Rol, &p.Personalidad, &p.Imagen, &p.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		personajes = append(personajes, p)
	}

	return c.JSON(personajes)
}

func (h PersonajesHandler) ObtenerPersonajeChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var p models.Personaje
	err = h.DB.QueryRow(context.Background(), `SELECT id, serie_id, nombre, descripcion, rol,
		personalidad, imagen, created_at FROM personajes WHERE id=$1`, id).
		Scan(&p.ID, &p.SerieID, &p.Nombre, &p.Descripcion, &p.Rol, &p.Personalidad, &p.Imagen, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "personaje no encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(p)
}

func (h PersonajesHandler) CrearPersonajeChen(c *fiber.Ctx) error {
	var p models.Personaje
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if strings.TrimSpace(p.Nombre) == "" || p.SerieID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nombre y serie_id son obligatorios"})
	}

	err := h.DB.QueryRow(context.Background(), `INSERT INTO personajes
		(serie_id, nombre, descripcion, rol, personalidad, imagen)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`,
		p.SerieID, p.Nombre, p.Descripcion, p.Rol, p.Personalidad, p.Imagen,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h PersonajesHandler) ActualizarPersonajeJosuc(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var p models.Personaje
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if strings.TrimSpace(p.Nombre) == "" || p.SerieID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nombre y serie_id son obligatorios"})
	}

	err = h.DB.QueryRow(context.Background(), `UPDATE personajes SET
		serie_id=$1, nombre=$2, descripcion=$3, rol=$4, personalidad=$5, imagen=$6
		WHERE id=$7 RETURNING id, serie_id, nombre, descripcion, rol, personalidad, imagen, created_at`,
		p.SerieID, p.Nombre, p.Descripcion, p.Rol, p.Personalidad, p.Imagen, id,
	).Scan(&p.ID, &p.SerieID, &p.Nombre, &p.Descripcion, &p.Rol, &p.Personalidad, &p.Imagen, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "personaje no encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(p)
}

func (h PersonajesHandler) EliminarPersonajeChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	result, err := h.DB.Exec(context.Background(), "DELETE FROM personajes WHERE id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "personaje no encontrado"})
	}

	return c.JSON(fiber.Map{"mensaje": "personaje eliminado"})
}

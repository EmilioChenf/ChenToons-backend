package handlers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"chentoons-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EpisodiosHandler struct {
	DB *pgxpool.Pool
}

func (h EpisodiosHandler) CargarEpisodiosEmilio(c *fiber.Ctx) error {
	query := `SELECT id, serie_id, titulo, temporada, numero_episodio, descripcion,
		duracion_minutos, fecha_estreno, created_at FROM episodios`
	args := []any{}
	if c.Query("serie_id") != "" {
		serieID, err := strconv.Atoi(c.Query("serie_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "serie_id invalido"})
		}
		query += " WHERE serie_id=$1"
		args = append(args, serieID)
	}
	query += " ORDER BY serie_id ASC, temporada ASC, numero_episodio ASC"

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	episodios := []models.Episodio{}
	for rows.Next() {
		e, err := scanEpisodioChen(rows)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		episodios = append(episodios, e)
	}
	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(episodios)
}

func (h EpisodiosHandler) ObtenerEpisodioChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	row := h.DB.QueryRow(context.Background(), `SELECT id, serie_id, titulo, temporada,
		numero_episodio, descripcion, duracion_minutos, fecha_estreno, created_at
		FROM episodios WHERE id=$1`, id)

	e, err := scanEpisodioChen(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "episodio no encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(e)
}

func (h EpisodiosHandler) CrearEpisodioChen(c *fiber.Ctx) error {
	var e models.Episodio
	if err := c.BodyParser(&e); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if strings.TrimSpace(e.Titulo) == "" || e.SerieID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "titulo y serie_id son obligatorios"})
	}
	existe, err := existeSerieChenin(h.DB, e.SerieID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !existe {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "serie no encontrada"})
	}
	fecha, err := fechaEpisodioJosuc(e.FechaEstreno)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_estreno debe usar formato YYYY-MM-DD"})
	}

	err = h.DB.QueryRow(context.Background(), `INSERT INTO episodios
		(serie_id, titulo, temporada, numero_episodio, descripcion, duracion_minutos, fecha_estreno)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at`,
		e.SerieID, e.Titulo, e.Temporada, e.NumeroEpisodio, e.Descripcion, e.DuracionMinutos, fecha,
	).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(e)
}

func (h EpisodiosHandler) ActualizarEpisodioChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var e models.Episodio
	if err := c.BodyParser(&e); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if strings.TrimSpace(e.Titulo) == "" || e.SerieID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "titulo y serie_id son obligatorios"})
	}
	existe, err := existeSerieChenin(h.DB, e.SerieID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !existe {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "serie no encontrada"})
	}
	fecha, err := fechaEpisodioJosuc(e.FechaEstreno)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_estreno debe usar formato YYYY-MM-DD"})
	}

	row := h.DB.QueryRow(context.Background(), `UPDATE episodios SET
		serie_id=$1, titulo=$2, temporada=$3, numero_episodio=$4, descripcion=$5,
		duracion_minutos=$6, fecha_estreno=$7
		WHERE id=$8 RETURNING id, serie_id, titulo, temporada, numero_episodio,
		descripcion, duracion_minutos, fecha_estreno, created_at`,
		e.SerieID, e.Titulo, e.Temporada, e.NumeroEpisodio, e.Descripcion, e.DuracionMinutos, fecha, id,
	)
	e, err = scanEpisodioChen(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "episodio no encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(e)
}

func (h EpisodiosHandler) EliminarEpisodioChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	result, err := h.DB.Exec(context.Background(), "DELETE FROM episodios WHERE id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "episodio no encontrado"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type filaEpisodio interface {
	Scan(dest ...any) error
}

func scanEpisodioChen(row filaEpisodio) (models.Episodio, error) {
	var e models.Episodio
	var fecha pgtype.Date
	err := row.Scan(&e.ID, &e.SerieID, &e.Titulo, &e.Temporada, &e.NumeroEpisodio, &e.Descripcion, &e.DuracionMinutos, &fecha, &e.CreatedAt)
	if err != nil {
		return e, err
	}
	if fecha.Valid {
		e.FechaEstreno = fecha.Time.Format("2006-01-02")
	}
	return e, nil
}

func fechaEpisodioJosuc(valor string) (any, error) {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return nil, nil
	}
	fecha, err := time.Parse("2006-01-02", valor)
	if err != nil {
		return nil, err
	}
	return fecha, nil
}

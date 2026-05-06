package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"chentoons-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeriesHandler struct {
	DB *pgxpool.Pool
}

type respuestaListaChen struct {
	Data       []models.Serie `json:"data"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

func (h SeriesHandler) CargarSeriesChenin(c *fiber.Ctx) error {
	page := leerIntQueryChen(c, "page", 1)
	limit := leerIntQueryChen(c, "limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	sort := ordenarCampoSerie(c.Query("sort", "id"))
	order := strings.ToUpper(c.Query("order", "DESC"))
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	where, args := filtrosSeriesChen(c)
	offset := (page - 1) * limit

	var total int
	countSQL := "SELECT COUNT(*) FROM series " + where
	if err := h.DB.QueryRow(context.Background(), countSQL, args...).Scan(&total); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	argsLista := append(args, limit, offset)
	query := fmt.Sprintf(`SELECT id, nombre, descripcion, categoria, genero, anio_lanzamiento,
		temporadas, estado, plataforma, creador, pais_origen, imagen, color_tema,
		destacada, created_at, updated_at
		FROM series %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		where, sort, order, len(args)+1, len(args)+2,
	)

	rows, err := h.DB.Query(context.Background(), query, argsLista...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	series := []models.Serie{}
	for rows.Next() {
		serie, err := scanSerieJosuc(rows)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		series = append(series, serie)
	}
	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return c.JSON(respuestaListaChen{
		Data:       series,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h SeriesHandler) ObtenerSerieChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	row := h.DB.QueryRow(context.Background(), `SELECT id, nombre, descripcion, categoria, genero,
		anio_lanzamiento, temporadas, estado, plataforma, creador, pais_origen,
		imagen, color_tema, destacada, created_at, updated_at FROM series WHERE id=$1`, id)

	serie, err := scanSerieJosuc(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "serie no encontrada"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(serie)
}

func (h SeriesHandler) CrearSerieChen(c *fiber.Ctx) error {
	var serie models.Serie
	if err := c.BodyParser(&serie); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if strings.TrimSpace(serie.Nombre) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nombre es obligatorio"})
	}

	err := h.DB.QueryRow(context.Background(), `INSERT INTO series
		(nombre, descripcion, categoria, genero, anio_lanzamiento, temporadas, estado,
		plataforma, creador, pais_origen, imagen, color_tema, destacada)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, created_at, updated_at`,
		serie.Nombre, serie.Descripcion, serie.Categoria, serie.Genero, serie.AnioLanzamiento,
		serie.Temporadas, serie.Estado, serie.Plataforma, serie.Creador, serie.PaisOrigen,
		serie.Imagen, serie.ColorTema, serie.Destacada,
	).Scan(&serie.ID, &serie.CreatedAt, &serie.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(serie)
}

func (h SeriesHandler) ActualizarSerieChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	var serie models.Serie
	if err := c.BodyParser(&serie); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "json invalido"})
	}
	if strings.TrimSpace(serie.Nombre) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nombre es obligatorio"})
	}

	row := h.DB.QueryRow(context.Background(), `UPDATE series SET
		nombre=$1, descripcion=$2, categoria=$3, genero=$4, anio_lanzamiento=$5,
		temporadas=$6, estado=$7, plataforma=$8, creador=$9, pais_origen=$10,
		imagen=$11, color_tema=$12, destacada=$13, updated_at=NOW()
		WHERE id=$14
		RETURNING id, nombre, descripcion, categoria, genero, anio_lanzamiento,
		temporadas, estado, plataforma, creador, pais_origen, imagen, color_tema,
		destacada, created_at, updated_at`,
		serie.Nombre, serie.Descripcion, serie.Categoria, serie.Genero, serie.AnioLanzamiento,
		serie.Temporadas, serie.Estado, serie.Plataforma, serie.Creador, serie.PaisOrigen,
		serie.Imagen, serie.ColorTema, serie.Destacada, id,
	)

	actualizada, err := scanSerieJosuc(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "serie no encontrada"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(actualizada)
}

func (h SeriesHandler) EliminarSerieChen(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalido"})
	}

	result, err := h.DB.Exec(context.Background(), "DELETE FROM series WHERE id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "serie no encontrada"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h SeriesHandler) ExportarCSVChenin(c *fiber.Ctx) error {
	rows, err := h.DB.Query(context.Background(), `SELECT id, nombre, categoria, genero,
		anio_lanzamiento, temporadas, estado, plataforma, creador, pais_origen,
		imagen, destacada FROM series ORDER BY id ASC`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="series_chentoons.csv"`)

	writer := csv.NewWriter(c.Response().BodyWriter())
	_ = writer.Write([]string{"id", "nombre", "categoria", "genero", "anio_lanzamiento", "temporadas", "estado", "plataforma", "creador", "pais_origen", "imagen", "destacada"})

	for rows.Next() {
		var s models.Serie
		if err := rows.Scan(&s.ID, &s.Nombre, &s.Categoria, &s.Genero, &s.AnioLanzamiento, &s.Temporadas, &s.Estado, &s.Plataforma, &s.Creador, &s.PaisOrigen, &s.Imagen, &s.Destacada); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		_ = writer.Write([]string{
			strconv.Itoa(s.ID), s.Nombre, s.Categoria, s.Genero, strconv.Itoa(s.AnioLanzamiento),
			strconv.Itoa(s.Temporadas), s.Estado, s.Plataforma, s.Creador, s.PaisOrigen,
			s.Imagen, strconv.FormatBool(s.Destacada),
		})
	}
	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	writer.Flush()

	return nil
}

func filtrosSeriesChen(c *fiber.Ctx) (string, []any) {
	partes := []string{}
	args := []any{}

	agregar := func(campo string, valor string, exacto bool) {
		if strings.TrimSpace(valor) == "" {
			return
		}
		args = append(args, valor)
		if exacto {
			args[len(args)-1] = strings.ToLower(valor)
			partes = append(partes, fmt.Sprintf("LOWER(%s) = $%d", campo, len(args)))
			return
		}
		args[len(args)-1] = "%" + strings.ToLower(valor) + "%"
		partes = append(partes, fmt.Sprintf("LOWER(%s) LIKE $%d", campo, len(args)))
	}

	search := c.Query("search")
	if strings.TrimSpace(search) == "" {
		search = c.Query("q")
	}
	agregar("nombre", search, false)
	agregar("genero", c.Query("genero"), true)
	agregar("categoria", c.Query("categoria"), true)
	agregar("estado", c.Query("estado"), true)

	if len(partes) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(partes, " AND "), args
}

func ordenarCampoSerie(campo string) string {
	permitidos := map[string]string{
		"id":               "id",
		"nombre":           "nombre",
		"genero":           "genero",
		"categoria":        "categoria",
		"anio_lanzamiento": "anio_lanzamiento",
		"temporadas":       "temporadas",
		"estado":           "estado",
		"created_at":       "created_at",
	}
	if valor, ok := permitidos[campo]; ok {
		return valor
	}
	return "id"
}

type filaSerie interface {
	Scan(dest ...any) error
}

func scanSerieJosuc(row filaSerie) (models.Serie, error) {
	var s models.Serie
	err := row.Scan(&s.ID, &s.Nombre, &s.Descripcion, &s.Categoria, &s.Genero,
		&s.AnioLanzamiento, &s.Temporadas, &s.Estado, &s.Plataforma, &s.Creador,
		&s.PaisOrigen, &s.Imagen, &s.ColorTema, &s.Destacada, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func leerIntQueryChen(c *fiber.Ctx, nombre string, fallback int) int {
	valor, err := strconv.Atoi(c.Query(nombre))
	if err != nil {
		return fallback
	}
	return valor
}

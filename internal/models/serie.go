package models

import "time"

type Serie struct {
	ID              int       `json:"id"`
	Nombre          string    `json:"nombre"`
	Descripcion     string    `json:"descripcion"`
	Categoria       string    `json:"categoria"`
	Genero          string    `json:"genero"`
	AnioLanzamiento int       `json:"anio_lanzamiento"`
	Temporadas      int       `json:"temporadas"`
	Estado          string    `json:"estado"`
	Plataforma      string    `json:"plataforma"`
	Creador         string    `json:"creador"`
	PaisOrigen      string    `json:"pais_origen"`
	Imagen          string    `json:"imagen"`
	ColorTema       string    `json:"color_tema"`
	Destacada       bool      `json:"destacada"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

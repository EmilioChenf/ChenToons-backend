package models

import "time"

type Personaje struct {
	ID           int       `json:"id"`
	SerieID      int       `json:"serie_id"`
	Nombre       string    `json:"nombre"`
	Descripcion  string    `json:"descripcion"`
	Rol          string    `json:"rol"`
	Personalidad string    `json:"personalidad"`
	CreatedAt    time.Time `json:"created_at"`
}

package models

import "time"

type Rating struct {
	ID         int       `json:"id"`
	SerieID    int       `json:"serie_id"`
	Puntuacion int       `json:"puntuacion"`
	Comentario string    `json:"comentario"`
	CreatedAt  time.Time `json:"created_at"`
}

package models

import "time"

type Episodio struct {
	ID              int       `json:"id"`
	SerieID         int       `json:"serie_id"`
	Titulo          string    `json:"titulo"`
	Temporada       int       `json:"temporada"`
	NumeroEpisodio  int       `json:"numero_episodio"`
	Descripcion     string    `json:"descripcion"`
	DuracionMinutos int       `json:"duracion_minutos"`
	FechaEstreno    string    `json:"fecha_estreno"`
	CreatedAt       time.Time `json:"created_at"`
}

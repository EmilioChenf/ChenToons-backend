package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MigrarTablasChenin(db *pgxpool.Pool) error {
	consultas := []string{
		`CREATE TABLE IF NOT EXISTS series (
			id SERIAL PRIMARY KEY,
			nombre VARCHAR(120) NOT NULL,
			descripcion TEXT,
			categoria VARCHAR(80),
			genero VARCHAR(80),
			anio_lanzamiento INT,
			temporadas INT DEFAULT 1,
			estado VARCHAR(50),
			plataforma VARCHAR(100),
			creador VARCHAR(120),
			pais_origen VARCHAR(80),
			imagen TEXT,
			color_tema VARCHAR(30),
			destacada BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS personajes (
			id SERIAL PRIMARY KEY,
			serie_id INT NOT NULL REFERENCES series(id) ON DELETE CASCADE,
			nombre VARCHAR(120) NOT NULL,
			descripcion TEXT,
			rol VARCHAR(80),
			personalidad VARCHAR(120),
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS episodios (
			id SERIAL PRIMARY KEY,
			serie_id INT NOT NULL REFERENCES series(id) ON DELETE CASCADE,
			titulo VARCHAR(150) NOT NULL,
			temporada INT DEFAULT 1,
			numero_episodio INT DEFAULT 1,
			descripcion TEXT,
			duracion_minutos INT,
			fecha_estreno DATE,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS ratings (
			id SERIAL PRIMARY KEY,
			serie_id INT NOT NULL REFERENCES series(id) ON DELETE CASCADE,
			puntuacion INT NOT NULL CHECK (puntuacion BETWEEN 1 AND 5),
			comentario TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
	}

	for _, consulta := range consultas {
		if _, err := db.Exec(context.Background(), consulta); err != nil {
			return err
		}
	}

	limpiezasChen := []string{
		`ALTER TABLE personajes DROP COLUMN IF EXISTS imagen;`,
		`ALTER TABLE episodios DROP COLUMN IF EXISTS imagen;`,
	}
	for _, consulta := range limpiezasChen {
		if _, err := db.Exec(context.Background(), consulta); err != nil {
			return err
		}
	}

	return nil
}

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SembrarDatosChenin(db *pgxpool.Pool) error {
	var total int
	if err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM series").Scan(&total); err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	series := []struct {
		nombre, descripcion, categoria, genero, estado, plataforma, creador, pais, imagen, color string
		anio, temporadas                                                                         int
		destacada                                                                                bool
	}{
		{"Pocoyo", "Un pequeno explorador azul aprende jugando con sus amigos.", "Preescolar", "Infantil", "En emision", "YouTube", "Guillermo Garcia Carsi", "Espana", "/uploads/pocoyo.jpg", "#4aa3df", 2005, 5, true},
		{"Escandalosos", "Tres hermanos osos intentan encajar en la ciudad moderna.", "Comedia", "Aventura", "Finalizada", "Max", "Daniel Chong", "Estados Unidos", "/uploads/escandalosos.jpg", "#70b77e", 2015, 4, true},
		{"Snoopy y Charlie Brown", "Snoopy, Charlie Brown y la pandilla viven pequenas historias con mucho corazon.", "Clasica", "Comedia", "Clasica", "Apple TV+", "Charles M. Schulz", "Estados Unidos", "/uploads/snoopy.jpg", "#f6d04d", 1965, 3, true},
		{"Bluey", "Una cachorrita usa la imaginacion para convertir cada dia en una aventura.", "Familiar", "Infantil", "En emision", "Disney+", "Joe Brumm", "Australia", "/uploads/bluey.jpg", "#5b8fd9", 2018, 3, true},
		{"Peppa Pig", "Peppa comparte juegos, familia y aprendizajes sencillos.", "Preescolar", "Infantil", "En emision", "Netflix", "Neville Astley", "Reino Unido", "/uploads/peppa.jpg", "#f28bb2", 2004, 8, false},
		{"Doraemon", "Un gato robot del futuro ayuda a Nobita con inventos increibles.", "Anime", "Comedia", "En emision", "Crunchyroll", "Fujiko F. Fujio", "Japon", "/uploads/doraemon.jpg", "#3d9bd8", 1979, 10, true},
		{"Hora de Aventura", "Finn y Jake recorren Ooo entre dulces, magia y caos divertido.", "Fantasia", "Aventura", "Finalizada", "Max", "Pendleton Ward", "Estados Unidos", "/uploads/hora-aventura.jpg", "#6ec6b8", 2010, 10, false},
		{"El Laboratorio de Dexter", "Dexter esconde un laboratorio enorme mientras Dee Dee lo complica todo.", "Clasica", "Comedia", "Finalizada", "Cartoon Network", "Genndy Tartakovsky", "Estados Unidos", "/uploads/dexter.jpg", "#e45b4f", 1996, 4, false},
		{"Las Chicas Superpoderosas", "Bombon, Burbuja y Bellota defienden Saltadilla antes de dormir.", "Superheroes", "Accion", "Finalizada", "Max", "Craig McCracken", "Estados Unidos", "/uploads/chicas-superpoderosas.jpg", "#f277a8", 1998, 6, false},
		{"Tom y Jerry", "Un gato y un raton convierten cada persecucion en comedia visual.", "Clasica", "Comedia", "Clasica", "Max", "William Hanna", "Estados Unidos", "/uploads/tom-jerry.jpg", "#a96f4c", 1940, 5, true},
		{"Masha y el Oso", "Masha visita a su amigo Oso y siempre aparece una travesura.", "Preescolar", "Comedia", "En emision", "YouTube", "Oleg Kuzovkov", "Rusia", "/uploads/masha.jpg", "#9ac65b", 2009, 5, false},
		{"Craig del Arroyo", "Craig y sus amigos exploran un arroyo lleno de pequenos mundos.", "Aventura", "Infantil", "Finalizada", "Max", "Matt Burnett", "Estados Unidos", "/uploads/craig.jpg", "#8fcb9b", 2018, 5, false},
		{"Gravity Falls", "Dipper y Mabel descubren misterios raros durante sus vacaciones.", "Misterio", "Aventura", "Finalizada", "Disney+", "Alex Hirsch", "Estados Unidos", "/uploads/gravity-falls.jpg", "#8061a8", 2012, 2, true},
		{"Steven Universe", "Steven aprende sobre amistad, familia y gemas magicas.", "Fantasia", "Aventura", "Finalizada", "Max", "Rebecca Sugar", "Estados Unidos", "/uploads/steven.jpg", "#ef8fa8", 2013, 5, false},
		{"Oggy y las Cucarachas", "Oggy solo quiere paz, pero las cucarachas tienen otros planes.", "Comedia", "Slapstick", "Finalizada", "YouTube", "Jean-Yves Raimbaud", "Francia", "/uploads/oggy.jpg", "#4d8ec9", 1998, 7, false},
		{"Los Rugrats", "Bebes curiosos interpretan el mundo adulto a su manera.", "Clasica", "Comedia", "Finalizada", "Paramount+", "Arlene Klasky", "Estados Unidos", "/uploads/rugrats.jpg", "#f0a84f", 1991, 9, false},
		{"Hey Arnold!", "Arnold vive historias de barrio con amigos muy distintos.", "Clasica", "Drama ligero", "Finalizada", "Paramount+", "Craig Bartlett", "Estados Unidos", "/uploads/hey-arnold.jpg", "#5db7a6", 1996, 5, false},
		{"El Increible Mundo de Gumball", "Gumball y Darwin mezclan humor absurdo con vida escolar.", "Comedia", "Aventura", "Finalizada", "Max", "Ben Bocquelet", "Reino Unido", "/uploads/gumball.jpg", "#59b8e8", 2011, 6, true},
	}

	ids := map[string]int{}
	for _, s := range series {
		var id int
		err := db.QueryRow(context.Background(), `INSERT INTO series
			(nombre, descripcion, categoria, genero, anio_lanzamiento, temporadas, estado, plataforma, creador, pais_origen, imagen, color_tema, destacada)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
			s.nombre, s.descripcion, s.categoria, s.genero, s.anio, s.temporadas, s.estado, s.plataforma, s.creador, s.pais, s.imagen, s.color, s.destacada,
		).Scan(&id)
		if err != nil {
			return err
		}
		ids[s.nombre] = id
	}

	personajes := []struct {
		serie, nombre, descripcion, rol, personalidad, imagen string
	}{
		{"Pocoyo", "Pocoyo", "Nino curioso vestido de azul.", "Protagonista", "Curioso y alegre", "/uploads/pocoyo-personaje.jpg"},
		{"Pocoyo", "Pato", "Amigo amarillo que intenta mantener el orden.", "Amigo", "Serio pero noble", "/uploads/pato.jpg"},
		{"Escandalosos", "Pardo", "Oso grizzly que lidera muchas ideas del grupo.", "Protagonista", "Sociable", "/uploads/pardo.jpg"},
		{"Escandalosos", "Panda", "Oso sensible que ama internet y el romance.", "Protagonista", "Timido", "/uploads/panda.jpg"},
		{"Escandalosos", "Polar", "Oso callado con habilidades sorprendentes.", "Protagonista", "Misterioso", "/uploads/polar.jpg"},
		{"Snoopy y Charlie Brown", "Snoopy", "Beagle imaginativo que suena con ser piloto.", "Protagonista", "Creativo", "/uploads/snoopy-personaje.jpg"},
		{"Snoopy y Charlie Brown", "Charlie Brown", "Nino amable que nunca deja de intentarlo.", "Protagonista", "Tierno", "/uploads/charlie-brown.jpg"},
		{"Bluey", "Bluey", "Cachorrita con energia para inventar juegos.", "Protagonista", "Imaginativa", "/uploads/bluey-personaje.jpg"},
		{"Doraemon", "Doraemon", "Gato robot con bolsillo magico.", "Protagonista", "Leal", "/uploads/doraemon-personaje.jpg"},
		{"Gravity Falls", "Mabel Pines", "Hermana alegre con sueteres inolvidables.", "Protagonista", "Optimista", "/uploads/mabel.jpg"},
	}
	for _, p := range personajes {
		if _, err := db.Exec(context.Background(), `INSERT INTO personajes
			(serie_id, nombre, descripcion, rol, personalidad, imagen)
			VALUES ($1,$2,$3,$4,$5,$6)`,
			ids[p.serie], p.nombre, p.descripcion, p.rol, p.personalidad, p.imagen,
		); err != nil {
			return err
		}
	}

	episodios := []struct {
		serie, titulo, descripcion, fecha, imagen string
		temporada, numero, duracion               int
	}{
		{"Pocoyo", "El baile de Pocoyo", "Pocoyo descubre el ritmo con sus amigos.", "2005-01-07", "/uploads/ep-pocoyo-baile.jpg", 1, 1, 7},
		{"Escandalosos", "Nuestra cueva", "Los osos intentan ordenar su hogar.", "2015-07-27", "/uploads/ep-cueva.jpg", 1, 1, 11},
		{"Snoopy y Charlie Brown", "Una pequena victoria", "Charlie Brown intenta levantar el animo.", "1965-12-09", "/uploads/ep-snoopy.jpg", 1, 1, 25},
		{"Bluey", "Magia de juegos", "Bluey convierte la casa en una aventura.", "2018-10-01", "/uploads/ep-bluey.jpg", 1, 1, 7},
		{"Gravity Falls", "Trampa turistica", "Dipper y Mabel llegan a Gravity Falls.", "2012-06-15", "/uploads/ep-gravity.jpg", 1, 1, 22},
		{"El Increible Mundo de Gumball", "El DVD", "Gumball y Darwin deben devolver una pelicula.", "2011-05-03", "/uploads/ep-gumball.jpg", 1, 1, 11},
	}
	for _, e := range episodios {
		if _, err := db.Exec(context.Background(), `INSERT INTO episodios
			(serie_id, titulo, temporada, numero_episodio, descripcion, duracion_minutos, fecha_estreno, imagen)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			ids[e.serie], e.titulo, e.temporada, e.numero, e.descripcion, e.duracion, e.fecha, e.imagen,
		); err != nil {
			return err
		}
	}

	ratings := []struct {
		serie, comentario string
		puntuacion        int
	}{
		{"Pocoyo", "Perfecta para explicar curiosidad y amistad.", 5},
		{"Escandalosos", "Muy divertida para ver en familia.", 5},
		{"Snoopy y Charlie Brown", "Clasico bonito y tranquilo.", 5},
		{"Bluey", "Historias cortas pero con mucho corazon.", 5},
		{"Gravity Falls", "Misterio y comedia bien mezclados.", 5},
		{"Tom y Jerry", "Comedia visual que no envejece.", 4},
	}
	for _, r := range ratings {
		if _, err := db.Exec(context.Background(), `INSERT INTO ratings
			(serie_id, puntuacion, comentario) VALUES ($1,$2,$3)`,
			ids[r.serie], r.puntuacion, r.comentario,
		); err != nil {
			return err
		}
	}

	return nil
}

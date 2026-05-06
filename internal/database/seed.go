package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serieSeedChen struct {
	nombre, descripcion, categoria, genero, estado, plataforma, creador, pais, imagen, color string
	anio, temporadas                                                                         int
	destacada                                                                                bool
}

type personajeSeedChen struct {
	nombre, descripcion, rol, personalidad string
}

type episodioSeedChen struct {
	titulo, descripcion, fecha  string
	temporada, numero, duracion int
}

type ratingSeedChen struct {
	puntuacion int
	comentario string
	fecha      string
}

type detalleSerieSeedChen struct {
	personajes []personajeSeedChen
	episodios  []episodioSeedChen
	ratings    []ratingSeedChen
}

func SembrarDatosChenin(db *pgxpool.Pool) error {
	ids, err := asegurarSeriesChenin(db)
	if err != nil {
		return err
	}

	detalles := detallesToonChenin()
	for nombreSerie, detalle := range detalles {
		serieID, ok := ids[nombreSerie]
		if !ok {
			continue
		}
		if err := asegurarPersonajesChenin(db, serieID, detalle.personajes); err != nil {
			return err
		}
		if err := asegurarEpisodiosChenin(db, serieID, detalle.episodios); err != nil {
			return err
		}
		if err := asegurarRatingsChenin(db, serieID, detalle.ratings); err != nil {
			return err
		}
	}

	return nil
}

func asegurarSeriesChenin(db *pgxpool.Pool) (map[string]int, error) {
	ids := map[string]int{}

	for _, s := range seriesBaseChenin() {
		var id int
		err := db.QueryRow(context.Background(), `SELECT id FROM series
			WHERE LOWER(nombre)=LOWER($1) ORDER BY id LIMIT 1`, s.nombre).Scan(&id)
		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		}

		if err == pgx.ErrNoRows {
			err = db.QueryRow(context.Background(), `INSERT INTO series
				(nombre, descripcion, categoria, genero, anio_lanzamiento, temporadas, estado,
				plataforma, creador, pais_origen, imagen, color_tema, destacada)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
				s.nombre, s.descripcion, s.categoria, s.genero, s.anio, s.temporadas, s.estado,
				s.plataforma, s.creador, s.pais, s.imagen, s.color, s.destacada,
			).Scan(&id)
			if err != nil {
				return nil, err
			}
		}

		ids[s.nombre] = id
	}

	return ids, nil
}

func asegurarPersonajesChenin(db *pgxpool.Pool, serieID int, personajes []personajeSeedChen) error {
	for _, p := range personajes {
		var existe bool
		err := db.QueryRow(context.Background(), `SELECT EXISTS(
			SELECT 1 FROM personajes WHERE serie_id=$1 AND LOWER(nombre)=LOWER($2)
		)`, serieID, p.nombre).Scan(&existe)
		if err != nil {
			return err
		}
		if existe {
			continue
		}

		_, err = db.Exec(context.Background(), `INSERT INTO personajes
			(serie_id, nombre, descripcion, rol, personalidad)
			VALUES ($1,$2,$3,$4,$5)`,
			serieID, p.nombre, p.descripcion, p.rol, p.personalidad,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func asegurarEpisodiosChenin(db *pgxpool.Pool, serieID int, episodios []episodioSeedChen) error {
	for _, e := range episodios {
		var existe bool
		err := db.QueryRow(context.Background(), `SELECT EXISTS(
			SELECT 1 FROM episodios WHERE serie_id=$1 AND LOWER(titulo)=LOWER($2)
		)`, serieID, e.titulo).Scan(&existe)
		if err != nil {
			return err
		}
		if existe {
			continue
		}

		_, err = db.Exec(context.Background(), `INSERT INTO episodios
			(serie_id, titulo, temporada, numero_episodio, descripcion, duracion_minutos, fecha_estreno)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			serieID, e.titulo, e.temporada, e.numero, e.descripcion, e.duracion, e.fecha,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func asegurarRatingsChenin(db *pgxpool.Pool, serieID int, ratings []ratingSeedChen) error {
	for _, r := range ratings {
		var existe bool
		err := db.QueryRow(context.Background(), `SELECT EXISTS(
			SELECT 1 FROM ratings WHERE serie_id=$1 AND LOWER(comentario)=LOWER($2)
		)`, serieID, r.comentario).Scan(&existe)
		if err != nil {
			return err
		}
		if existe {
			continue
		}

		_, err = db.Exec(context.Background(), `INSERT INTO ratings
			(serie_id, puntuacion, comentario, created_at) VALUES ($1,$2,$3,$4)`,
			serieID, r.puntuacion, r.comentario, r.fecha,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func seriesBaseChenin() []serieSeedChen {
	return []serieSeedChen{
		{"Pocoyo", "Un pequeno explorador azul aprende jugando con sus amigos.", "Preescolar", "Infantil", "En emision", "YouTube", "Guillermo Garcia Carsi", "Espana", "/uploads/pocoyo.png", "#4aa3df", 2005, 5, true},
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
}

func detallesToonChenin() map[string]detalleSerieSeedChen {
	datos := map[string]detalleSerieSeedChen{
		"Pocoyo": {
			personajes: []personajeSeedChen{
				{"Pocoyo", "Nino curioso vestido de azul.", "Protagonista", "Curioso y alegre"},
				{"Pato", "Amigo amarillo que intenta mantener el orden.", "Amigo", "Serio pero noble"},
				{"Elly", "Elefanta rosada que cuida a sus amigos.", "Amiga", "Dulce y paciente"},
				{"Loula", "Perrita fiel que acompana las aventuras.", "Mascota", "Juguetona"},
			},
			episodios: episodiosChen("pocoyo", 2005, []string{
				"El baile de Pocoyo", "La carrera de Pato", "El regalo de Elly", "Loula se esconde", "Musica en el parque", "La gran torre azul",
			}),
			ratings: ratingsChen("Perfecta para explicar curiosidad y amistad.", "Muy bonita para verla en familia.", "Tiene capitulos simples pero entretenidos.", "La volveria a ver por sus colores y musica."),
		},
		"Escandalosos": {
			personajes: []personajeSeedChen{
				{"Pardo", "Oso grizzly que lidera muchas ideas del grupo.", "Protagonista", "Sociable"},
				{"Panda", "Oso sensible que ama internet y el romance.", "Protagonista", "Timido"},
				{"Polar", "Oso callado con habilidades sorprendentes.", "Protagonista", "Misterioso"},
				{"Chloe", "Nina inteligente que se vuelve amiga de los osos.", "Amiga", "Aplicada"},
				{"Nom Nom", "Koala famoso que complica varias situaciones.", "Rival", "Egocentrico"},
			},
			episodios: episodiosChen("escandalosos", 2015, []string{
				"Nuestra cueva", "La selfie perfecta", "El dia de Chloe", "Polar cocina", "Panda en linea", "Nom Nom visita",
			}),
			ratings: ratingsChen("Muy divertida para ver en familia.", "Me gusto porque los osos tienen personalidades muy distintas.", "Los capitulos son cortos y se sienten frescos.", "La volveria a ver por Pardo, Panda y Polar."),
		},
		"Snoopy y Charlie Brown": {
			personajes: []personajeSeedChen{
				{"Snoopy", "Beagle imaginativo que suena con ser piloto.", "Protagonista", "Creativo"},
				{"Charlie Brown", "Nino amable que nunca deja de intentarlo.", "Protagonista", "Tierno"},
				{"Lucy", "Amiga directa que siempre opina fuerte.", "Amiga", "Segura"},
				{"Linus", "Nino tranquilo con su manta inseparable.", "Amigo", "Reflexivo"},
			},
			episodios: episodiosChen("snoopy", 1965, []string{
				"Una pequena victoria", "La casita de Snoopy", "El partido de beisbol", "La manta de Linus", "Carta para el piloto", "El arbol de Navidad",
			}),
			ratings: ratingsChen("Clasico bonito y tranquilo.", "Me gusto porque los personajes son memorables.", "Tiene humor sencillo pero con mucho corazon.", "La volveria a ver por Snoopy y Charlie Brown."),
		},
		"Bluey": {
			personajes: []personajeSeedChen{
				{"Bluey", "Cachorrita con energia para inventar juegos.", "Protagonista", "Imaginativa"},
				{"Bingo", "Hermana menor que juega con mucha ternura.", "Protagonista", "Dulce"},
				{"Bandit", "Papa de Bluey que se presta para todos los juegos.", "Padre", "Divertido"},
				{"Chilli", "Mama tranquila que acompana los aprendizajes.", "Madre", "Paciente"},
			},
			episodios: episodiosChen("bluey", 2018, []string{
				"Magia de juegos", "El hospital de Bingo", "Carrera en el patio", "La tienda imaginaria", "Camping familiar", "La leccion de papa",
			}),
			ratings: ratingsChen("Historias cortas pero con mucho corazon.", "Muy bonita para verla en familia.", "Cada capitulo deja una idea facil de entender.", "La volveria a ver por la relacion familiar."),
		},
		"Peppa Pig": {
			personajes: []personajeSeedChen{
				{"Peppa", "Cerdita curiosa que disfruta jugar con su familia.", "Protagonista", "Alegre"},
				{"George", "Hermano menor que ama los dinosaurios.", "Hermano", "Tierno"},
				{"Mama Pig", "Mama paciente que guia a Peppa.", "Madre", "Amable"},
				{"Papa Pig", "Papa bromista que siempre intenta ayudar.", "Padre", "Divertido"},
			},
			episodios: episodiosChen("peppa", 2004, []string{
				"Charcos de lodo", "El dinosaurio de George", "La visita de los abuelos", "Dia de escuela", "Picnic en familia", "El globo rojo",
			}),
			ratings: ratingsChen("Buena para ninos pequenos.", "Tiene capitulos simples pero entretenidos.", "Me gusto porque muestra momentos familiares.", "La volveria a ver por su humor tranquilo."),
		},
		"Doraemon": {
			personajes: []personajeSeedChen{
				{"Doraemon", "Gato robot con bolsillo magico.", "Protagonista", "Leal"},
				{"Nobita", "Nino distraido que aprende con ayuda de Doraemon.", "Protagonista", "Sonador"},
				{"Shizuka", "Amiga amable de Nobita.", "Amiga", "Inteligente"},
				{"Gigante", "Companero fuerte que suele intimidar.", "Companero", "Impulsivo"},
				{"Suneo", "Nino presumido que crea conflictos.", "Companero", "Vanidoso"},
			},
			episodios: episodiosChen("doraemon", 1979, []string{
				"El bolsillo magico", "La maquina del tiempo", "Nobita invisible", "El examen imposible", "El invento perdido", "Regreso al futuro",
			}),
			ratings: ratingsChen("Los inventos hacen que cada capitulo sea distinto.", "Me gusto por la amistad entre Doraemon y Nobita.", "Tiene ideas creativas y graciosas.", "La volveria a ver por sus aventuras."),
		},
		"Hora de Aventura": {
			personajes: []personajeSeedChen{
				{"Finn", "Heroe humano que busca aventuras en Ooo.", "Protagonista", "Valiente"},
				{"Jake", "Perro magico que cambia de forma.", "Protagonista", "Relajado"},
				{"Dulce Princesa", "Gobernante inteligente del Dulce Reino.", "Aliada", "Analitica"},
				{"Marceline", "Vampira musica con historia misteriosa.", "Aliada", "Libre"},
			},
			episodios: episodiosChen("hora-aventura", 2010, []string{
				"El heroe de Ooo", "La espada perdida", "Dulce Reino en peligro", "Cancion de Marceline", "Jake se estira", "Mazmorra de amigos",
			}),
			ratings: ratingsChen("Tiene mucha imaginacion y humor raro.", "Me gusto porque mezcla aventura con momentos emotivos.", "Los personajes cambian bastante con el tiempo.", "La volveria a ver por el mundo de Ooo."),
		},
		"El Laboratorio de Dexter": {
			personajes: []personajeSeedChen{
				{"Dexter", "Nino cientifico con laboratorio secreto.", "Protagonista", "Genial y terco"},
				{"Dee Dee", "Hermana que entra al laboratorio sin permiso.", "Hermana", "Traviesa"},
				{"Mandark", "Rival cientifico de Dexter.", "Rival", "Presumido"},
				{"Mama", "Mama de Dexter que mantiene la casa en orden.", "Familia", "Protectora"},
			},
			episodios: episodiosChen("dexter", 1996, []string{
				"El laboratorio secreto", "Dee Dee toca botones", "Mandark reta a Dexter", "Robot fuera de control", "Formula para crecer", "Experimento en casa",
			}),
			ratings: ratingsChen("Muy divertida por el contraste entre Dexter y Dee Dee.", "Me gusto el humor de ciencia exagerada.", "Tiene capitulos rapidos y faciles de recordar.", "La volveria a ver por nostalgia."),
		},
		"Las Chicas Superpoderosas": {
			personajes: []personajeSeedChen{
				{"Bombon", "Lider del equipo y estratega.", "Heroina", "Responsable"},
				{"Burbuja", "Heroina tierna con gran sensibilidad.", "Heroina", "Dulce"},
				{"Bellota", "Heroina fuerte que enfrenta el peligro directo.", "Heroina", "Valiente"},
				{"Mojo Jojo", "Villano que siempre prepara planes complicados.", "Villano", "Dramatico"},
			},
			episodios: episodiosChen("chicas-superpoderosas", 1998, []string{
				"Salvando Saltadilla", "El plan de Mojo", "Burbuja al rescate", "Bellota no se rinde", "Bombon lidera", "Dia de escuela heroico",
			}),
			ratings: ratingsChen("Tiene accion y comedia en buen balance.", "Me gusto porque las protagonistas tienen estilos distintos.", "Los villanos son memorables.", "La volveria a ver por Mojo Jojo."),
		},
		"Tom y Jerry": {
			personajes: []personajeSeedChen{
				{"Tom", "Gato que siempre intenta atrapar a Jerry.", "Protagonista", "Insistente"},
				{"Jerry", "Raton astuto que se escapa con ingenio.", "Protagonista", "Listo"},
				{"Spike", "Perro fuerte que protege su tranquilidad.", "Secundario", "Gruñon"},
				{"Tyke", "Cachorro pequeno que acompana a Spike.", "Secundario", "Inocente"},
			},
			episodios: episodiosChen("tom-jerry", 1940, []string{
				"La persecucion", "Jerry en la cocina", "Tom pianista", "Spike se enoja", "El queso perdido", "Noche en la casa",
			}),
			ratings: ratingsChen("Comedia visual que no envejece.", "Me gusto porque casi no necesita dialogos.", "Tiene capitulos simples pero muy entretenidos.", "La volveria a ver por las persecuciones."),
		},
		"Masha y el Oso": {
			personajes: []personajeSeedChen{
				{"Masha", "Nina inquieta que siempre visita al Oso.", "Protagonista", "Traviesa"},
				{"Oso", "Oso tranquilo que cuida su casa.", "Protagonista", "Paciente"},
				{"Osa", "Vecina que llama la atencion del Oso.", "Secundaria", "Elegante"},
				{"Liebre", "Amigo del bosque que suele quedar atrapado en juegos.", "Amigo", "Nervioso"},
			},
			episodios: episodiosChen("masha", 2009, []string{
				"Visita inesperada", "La receta de Masha", "Oso quiere dormir", "Juego en el bosque", "Fiesta de invierno", "El cuadro arruinado",
			}),
			ratings: ratingsChen("Muy graciosa para ratos cortos.", "Me gusto porque Masha siempre causa situaciones nuevas.", "El Oso hace que la serie sea tierna.", "La volveria a ver por sus capitulos ligeros."),
		},
		"Craig del Arroyo": {
			personajes: []personajeSeedChen{
				{"Craig", "Nino explorador que dibuja mapas del arroyo.", "Protagonista", "Curioso"},
				{"Kelsey", "Amiga aventurera que imagina misiones epicas.", "Amiga", "Valiente"},
				{"JP", "Amigo noble que acompana cualquier plan.", "Amigo", "Relajado"},
				{"Jessica", "Hermana menor de Craig.", "Familia", "Inteligente"},
			},
			episodios: episodiosChen("craig", 2018, []string{
				"El mapa del arroyo", "La aventura secreta", "El fuerte perdido", "Carrera entre amigos", "Misterio en el bosque", "La reunion del puente",
			}),
			ratings: ratingsChen("Me gusto porque convierte un arroyo en un mundo enorme.", "Los personajes se sienten como amigos reales.", "Tiene capitulos simples pero entretenidos.", "La volveria a ver por su estilo de aventura."),
		},
		"Gravity Falls": {
			personajes: []personajeSeedChen{
				{"Dipper Pines", "Nino curioso que investiga misterios.", "Protagonista", "Analitico"},
				{"Mabel Pines", "Hermana alegre con sueteres inolvidables.", "Protagonista", "Optimista"},
				{"Stan Pines", "Tio de los gemelos y dueno de la cabana.", "Familia", "Astuto"},
				{"Soos", "Empleado amable que ayuda en la cabana.", "Amigo", "Leal"},
				{"Wendy", "Amiga tranquila que trabaja con Stan.", "Amiga", "Relajada"},
			},
			episodios: episodiosChen("gravity-falls", 2012, []string{
				"Trampa turistica", "El diario numero tres", "Mabel gana un cerdo", "Misterio en el lago", "La tienda encantada", "El secreto de la cabana",
			}),
			ratings: ratingsChen("Misterio y comedia bien mezclados.", "Me gusto porque cada detalle parece importante.", "Los personajes son memorables y graciosos.", "La volveria a ver por sus secretos."),
		},
		"Steven Universe": {
			personajes: []personajeSeedChen{
				{"Steven", "Nino con poderes de gema y mucho corazon.", "Protagonista", "Empatico"},
				{"Garnet", "Gema fuerte y tranquila que guia al grupo.", "Mentora", "Serena"},
				{"Amatista", "Gema libre que disfruta la diversion.", "Aliada", "Espontanea"},
				{"Perla", "Gema elegante y cuidadosa.", "Aliada", "Ordenada"},
			},
			episodios: episodiosChen("steven", 2013, []string{
				"El brillo de Steven", "Garnet decide", "La espada de Perla", "Amatista se transforma", "Cancion en la playa", "La burbuja magica",
			}),
			ratings: ratingsChen("Tiene canciones y momentos emotivos.", "Me gusto porque habla mucho de amistad.", "Los personajes evolucionan bastante.", "La volveria a ver por sus colores y musica."),
		},
		"Oggy y las Cucarachas": {
			personajes: []personajeSeedChen{
				{"Oggy", "Gato azul que intenta vivir tranquilo.", "Protagonista", "Paciente"},
				{"Joey", "Cucaracha pequena que lidera los planes.", "Antagonista", "Mandona"},
				{"Dee Dee", "Cucaracha glotona que causa desorden.", "Antagonista", "Hambrienta"},
				{"Marky", "Cucaracha relajada que sigue el caos.", "Antagonista", "Despreocupada"},
			},
			episodios: episodiosChen("oggy", 1998, []string{
				"El refrigerador abierto", "Cucarachas al ataque", "Oggy limpia la casa", "La pizza perdida", "Noche de television", "El jardin invadido",
			}),
			ratings: ratingsChen("Es caotica pero muy facil de ver.", "Me gusto porque usa mucho humor fisico.", "Tiene capitulos simples pero entretenidos.", "La volveria a ver por las ocurrencias de Oggy."),
		},
		"Los Rugrats": {
			personajes: []personajeSeedChen{
				{"Tommy", "Bebe lider que imagina grandes aventuras.", "Protagonista", "Valiente"},
				{"Chuckie", "Amigo miedoso pero muy leal.", "Amigo", "Nervioso"},
				{"Angelica", "Prima mayor que suele mandar al grupo.", "Prima", "Mandona"},
				{"Phil", "Bebe que acompana travesuras con su hermana.", "Amigo", "Jugueton"},
				{"Lil", "Hermana de Phil y companera de aventuras.", "Amiga", "Curiosa"},
			},
			episodios: episodiosChen("rugrats", 1991, []string{
				"Aventura en la sala", "El juguete perdido", "Angelica manda", "Chuckie se anima", "Phil y Lil exploran", "Dia en el parque",
			}),
			ratings: ratingsChen("Tiene mucha imaginacion desde la mirada de bebes.", "Me gusto porque los personajes son muy distintos.", "Los capitulos tienen nostalgia y humor.", "La volveria a ver por Tommy y Chuckie."),
		},
		"Hey Arnold!": {
			personajes: []personajeSeedChen{
				{"Arnold", "Nino amable que vive historias de barrio.", "Protagonista", "Solidario"},
				{"Gerald", "Mejor amigo de Arnold y gran narrador.", "Amigo", "Seguro"},
				{"Helga", "Companera intensa que oculta sus sentimientos.", "Companera", "Fuerte"},
				{"Abuelo Phil", "Abuelo de Arnold que cuenta historias raras.", "Familia", "Divertido"},
			},
			episodios: episodiosChen("hey-arnold", 1996, []string{
				"El barrio despierta", "Gerald cuenta una historia", "El secreto de Helga", "Partido en la calle", "La pension de los abuelos", "Arnold ayuda a un amigo",
			}),
			ratings: ratingsChen("Tiene historias urbanas con mucho corazon.", "Me gusto porque mezcla humor con temas cotidianos.", "Los personajes secundarios son memorables.", "La volveria a ver por Arnold y Helga."),
		},
		"El Increible Mundo de Gumball": {
			personajes: []personajeSeedChen{
				{"Gumball", "Gato azul que vive problemas absurdos.", "Protagonista", "Impulsivo"},
				{"Darwin", "Hermano y mejor amigo de Gumball.", "Protagonista", "Noble"},
				{"Anais", "Hermana menor muy inteligente.", "Familia", "Lista"},
				{"Nicole", "Mama de la familia Watterson.", "Madre", "Fuerte"},
				{"Richard", "Papa torpe pero carinoso.", "Padre", "Despistado"},
			},
			episodios: episodiosChen("gumball", 2011, []string{
				"El DVD", "La deuda", "El tercero", "La pintura", "El club", "La consola perdida",
			}),
			ratings: ratingsChen("El humor absurdo funciona muy bien.", "Me gusto porque mezcla estilos visuales diferentes.", "Tiene capitulos rapidos y muy creativos.", "La volveria a ver por Gumball y Darwin."),
		},
	}

	return datos
}

func episodiosChen(slug string, anio int, titulos []string) []episodioSeedChen {
	episodios := make([]episodioSeedChen, 0, len(titulos))
	for i, titulo := range titulos {
		episodios = append(episodios, episodioSeedChen{
			titulo:      titulo,
			temporada:   1 + i/4,
			numero:      i + 1,
			descripcion: descripcionEpisodioToon(titulo),
			duracion:    duracionToonChen(slug),
			fecha:       fechaToonChen(anio, i),
		})
	}
	return episodios
}

func ratingsChen(a, b, c, d string) []ratingSeedChen {
	return []ratingSeedChen{
		{5, a, "2026-01-05 10:00:00"},
		{4, b, "2026-01-06 11:00:00"},
		{5, c, "2026-01-07 12:00:00"},
		{4, d, "2026-01-08 13:00:00"},
	}
}

func descripcionEpisodioToon(titulo string) string {
	return titulo + " presenta una aventura corta con momentos divertidos y faciles de seguir."
}

func duracionToonChen(slug string) int {
	switch slug {
	case "pocoyo", "bluey", "peppa":
		return 7
	case "snoopy":
		return 25
	case "gravity-falls", "steven":
		return 22
	default:
		return 11
	}
}

func fechaToonChen(anio int, indice int) string {
	mes := 1 + indice
	if mes > 12 {
		mes = 12
	}
	return fmt.Sprintf("%04d-%02d-%02d", anio, mes, 5+indice)
}

package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Registro struct {
	Author        string
	Handle        string
	Facultad      string
	Carrera       string
	Fecha         string
	Bibliotecario string
}

var db *sql.DB

// Inicializa la base de datos y crea la tabla si no existe
func InitDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// Crear la tabla si no existe
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS registro (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			author TEXT NOT NULL, 
			handle TEXT NOT NULL, 
			facultad TEXT NOT NULL,
			carrera TEXT NOT NULL,
			fecha TEXT NOT NULL,
			bibliotecario TEXT NOT NULL
		)`,
	)
	if err != nil {
		return fmt.Errorf("error al crear la tabla: %w", err)
	}

	return nil
}

func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}

// Verifica si la base de datos está inicializada
func checkDatabase() error {
	if db == nil {
		return fmt.Errorf("la base de datos no está inicializada")
	}
	return nil
}

// Agrega un registro a la base de datos
func AddRegistro(a *Registro) (int64, error) {
	if err := checkDatabase(); err != nil {
		return 0, err
	}

	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO registro (author, handle, facultad, carrera, fecha, bibliotecario) VALUES (?,?,?,?,?,?);`,
		a.Author, a.Handle, a.Facultad, a.Carrera, a.Fecha, a.Bibliotecario,
	)
	if err != nil {
		return 0, fmt.Errorf("error al insertar el registro: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el ID del último registro insertado: %w", err)
	}

	return id, nil
}

// Realiza una consulta en la base de datos basada en una columna y un valor
func FetchbyQuery(q string, column string) ([]Registro, error) {
	if err := checkDatabase(); err != nil {
		return nil, err
	}

	// Validar el nombre de la columna para evitar inyecciones SQL
	validColumns := map[string]bool{
		"author":        true,
		"handle":        true,
		"facultad":      true,
		"carrera":       true,
		"fecha":         true,
		"bibliotecario": true,
	}
	if !validColumns[column] {
		return nil, fmt.Errorf("columna no válida: %s", column)
	}

	// Agregar los comodines al parámetro
	query := "%" + q + "%"

	// Ejecutar la consulta
	rows, err := db.QueryContext(
		context.Background(),
		fmt.Sprintf(`SELECT author, handle, facultad, carrera, fecha, bibliotecario FROM registro WHERE %s LIKE ?`, column),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %w", err)
	}
	defer rows.Close()

	// Crear un slice para almacenar los registros
	var registros []Registro

	// Iterar sobre los resultados
	for rows.Next() {
		var registro Registro
		err := rows.Scan(
			&registro.Author,
			&registro.Handle,
			&registro.Facultad,
			&registro.Carrera,
			&registro.Fecha,
			&registro.Bibliotecario,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear un registro: %w", err)
		}
		registros = append(registros, registro)
	}

	// Verificar si hubo errores al iterar
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados: %w", err)
	}

	return registros, nil
}

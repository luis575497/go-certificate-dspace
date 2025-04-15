package main

import (
	"encoding/csv"
	"fmt"
	"generador-certificados/database"
	"generador-certificados/scraping"
	"generador-certificados/word"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	VariantDark = theme.VariantDark
)

var bibliotecarios = []string{
	"PAOLA DEL ROCIO AMAYA ARCE",
	"DIANA ALEXANDRA LEON BRAVO",
	"FRANCISCO TEODORO ASTUDILLO SAQUINAULA",
	"DIANA MARLENE FAJARDO PASÁN",
	"WILMAN GONZALO TANDAZO GUEVARA",
	"LOURDES GABRIELA ORELLANA GUERRA",
	"JESSICA ELIZABETH BERMEO SOTAMBA",
	"ERIKA SOFIA PEÑAFIEL VAZQUEZ",
	"NUBE DEL ROCIO SALTO MORQUECHO",
	"HECTOR BLADIMIR CABRERA RODRIGUEZ",
	"JUAN PABLO CRIOLLO SAQUICARAY",
	"VANESSA ALEXANDRA MORALES MARIÑO",
	"DANIEL RAMIRO CARRIÓN ROMÁN",
	"DORIS PATRICIA TENESACA CARDENAS",
}

var facultadComplexivo = map[string]map[string]string{
	"Facultad de Ciencias Agropecuarias": {
		"Carrera 1": "Medicina Veterinaria",
		"Carrera 2": "Agronomía",
	},
	"Facultad de Arquitectura": {
		"Carrera 1": "Arquitectura",
	},
	"Facultad de Artes": {
		"Carrera 1": "Artes Visuales",
		"Carrera 2": "Artes Escénicas",
		"Carrera 3": "Artes Musicales",
		"Carrera 4": "Diseño de Interirores",
		"Carrera 5": "Diseño Gráfico",
	},
	"Facultad de Ciencias Ecnómicas y Administrativas": {
		"Carrera 1": "Administración de Empresas",
		"Carrera 2": "Contabilida y Auditoría",
		"Carrera 3": "Economía",
		"Carrera 4": "Mercadotecnia",
		"Carrera 5": "Sociología",
		"Carrera 6": "Emprendimiento e Innovación",
	},
	"Facultad de Ciencias Médicas": {
		"Carrera 1": "Medicina",
		"Carrera 2": "Enfermería",
		"Carrera 3": "Fonoaudiología",
		"Carrera 4": "Fisioterapia",
		"Carrera 5": "Laboratorio Clínico",
		"Carrera 6": "Nutrición y Dietética",
		"Carrera 7": "Imagenología y Radiología",
		"Carrera 8": "Estimulación Temprana en Salud",
	},
	"Facultad de Psicología": {
		"Carrera 1": "Psicología",
		"Carrera 2": "Psicología Educativa",
	},
	"Facultad de Ciencias Químicas": {
		"Carrera 1": "Ingeniería Ambiental",
		"Carrera 2": "Bioquímica y Farmacia",
		"Carrera 3": "Ingeniería Química",
		"Carrera 4": "Ingeniería Industrial",
	},
	"Facultad de Odontología": {
		"Carrera 1": "Odontología",
	},
	"Facultad de Jurisprudencia y Ciencias Políticas y Sociales": {
		"Carrera 1": "Derecho",
		"Carrera 2": "Trabajo Social",
		"Carrera 3": "Género y Desarrollo",
		"Carrera 4": "Orientación Familiar",
	},
	"Facultad de Ingeniería": {
		"Carrera 1": "Computación",
		"Carrera 2": "Telecomunicaciones",
		"Carrera 3": "Electricidad",
		"Carrera 4": "Ingeniería Civil",
	},
	"Facultad de Ciencias de la Hospitalidad": {
		"Carrera 1": "Gastronomía",
		"Carrera 2": "Turismo",
		"Carrera 3": "Hospitalidad y Hotelería",
	},
	"Facultad de Filosofía, Letras y Ciencias de la Educación": {
		"Carrera 1":  "Cine",
		"Carrera 2":  "Educación Básica",
		"Carrera 3":  "Pedagogía de la Lengua y la Literatura",
		"Carrera 4":  "Pedagogía de la Historia y las Ciencias Sociales",
		"Carrera 5":  "Pedagogía de las Artes y Humanidades",
		"Carrera 6":  "Pedagogía de las Ciencias Experimentales: Matemática y Física",
		"Carrera 7":  "Pedagogía de lso Idiomas Nacionales y Extranjeros",
		"Carrera 8":  "Educación Inicial",
		"Carrera 9":  "Pedagogía de la Actividad Física y el Deporte",
		"Carrera 10": "Comunicación",
		"Carrera 11": "Periodismo",
		"Carrera 12": "Pedagogía de las Ciencias Experimentales: Química y Biología",
	},
}

func main() {
	err := database.InitDatabase("registro.db")
	if err != nil {
		fmt.Println("Error al inicializar la base de datos:", err)
		return
	}
	defer database.CloseDatabase()

	myApp := app.New()
	myWindow := myApp.NewWindow("Generador de Certificado")
	myApp.Settings().SetTheme(theme.DarkTheme())

	// Crear las pestañas principales
	tabs := container.NewAppTabs(
		container.NewTabItem("Tesis", createTesisTab(myApp, myWindow)),
		container.NewTabItem("Complexivos", createComplexivosTab(myApp, myWindow)),
		container.NewTabItem("Masivos", createMasivosTab(myApp, myWindow)),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	myWindow.Resize(fyne.NewSize(300, 600))
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

// Crear la pestaña para Tesis
func createTesisTab(myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	// Card de información
	card := createInfoCard()

	// Entrada de texto
	uuid := widget.NewEntry()
	uuid.SetPlaceHolder("Escribe la URL del trabajo de titulación aquí...")

	// Selectores
	tituloSelect := widget.NewLabel("Selecciona el tipo de estudio:")
	selectStudy := widget.NewSelect([]string{"Pregrado", "Posgrado"}, nil)

	entryPosgrado := widget.NewEntry()
	entryPosgrado.SetPlaceHolder("Escribe el nombre del posgrado aquí...")
	entryPosgrado.Hide()

	selectStudy.OnChanged = func(selected string) {
		if selected == "Posgrado" {
			entryPosgrado.Show()
		} else {
			entryPosgrado.Hide()
		}
	}

	selectStudy.SetSelected("Pregrado")

	tituloPersona := widget.NewLabel("Referencista:")
	selectorPersona := widget.NewSelect(bibliotecarios, nil)

	selectorPersona.SetSelected("PAOLA DEL ROCIO AMAYA ARCE")

	// Botones
	progressbar := widget.NewProgressBarInfinite()
	progressbar.Hide()

	botonCrearCertificado := createCertificadoButton(myApp, uuid, selectStudy, entryPosgrado, selectorPersona, progressbar)
	cerrarButton := widget.NewButton("Cerrar", func() {
		myWindow.Close()
	})

	botonBuscar := widget.NewButton("Buscar", func() {
		searchCertificado(myApp)
	})
	botonExportar := widget.NewButton("Exportar", func() {
		exportarPorFecha(myApp)
	})

	// Layout principal
	botonera := container.NewPadded(container.NewGridWithColumns(2,
		container.NewVBox(botonCrearCertificado),
		container.NewVBox(cerrarButton),
	))

	return container.NewPadded(container.NewVBox(
		card,
		layout.NewSpacer(),
		uuid,
		progressbar,
		layout.NewSpacer(),
		container.NewGridWithColumns(2,
			container.NewVBox(tituloSelect, selectStudy, entryPosgrado, tituloPersona, selectorPersona),
			container.NewVBox(widget.NewLabel(""), botonBuscar, botonExportar),
		),
		layout.NewSpacer(),
		botonera,
	))
}

// Crear la pestaña para Complexivos
func createComplexivosTab(myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	// Card de información
	card := createInfoCard()

	// Entrada de texto
	name := widget.NewEntry()
	name.SetPlaceHolder("Escribe el nombre del estudiante aquí...")

	// Selectores
	tituloSelect := widget.NewLabel("Selecciona el tipo de estudio:")
	selectStudy := widget.NewSelect([]string{"Pregrado", "Posgrado"}, nil)
	selectStudy.SetSelected("Pregrado")

	entryPosgrado := widget.NewEntry()
	entryPosgrado.SetPlaceHolder("Escribe el nombre del posgrado aquí...")
	entryPosgrado.Hide()

	selectStudy.OnChanged = func(selected string) {
		if selected == "Posgrado" {
			entryPosgrado.Show()
		} else {
			entryPosgrado.Hide()
		}
	}

	tituloPersona := widget.NewLabel("Referencista:")
	selectorPersona := widget.NewSelect(bibliotecarios, nil)

	// Botones
	progressbar := widget.NewProgressBarInfinite()
	progressbar.Hide()

	selectorPersona.SetSelected("PAOLA DEL ROCIO AMAYA ARCE")

	titleFacultad := widget.NewLabel("Selecciona la facultad:")
	tituloCarrera := widget.NewLabel("Selecciona la carrera:")
	carreraSelect := widget.NewSelect([]string{}, nil)

	// Obtener todas las facultades dinámicamente desde el mapa facultadComplexivo
	facultades := []string{}
	for facultad := range facultadComplexivo {
		facultades = append(facultades, facultad)
	}

	// Crear el selector de facultades
	facultadSelect := widget.NewSelect(facultades, func(facultadSeleccionada string) {
		if carreras, ok := facultadComplexivo[facultadSeleccionada]; ok {
			carrerasOptions := []string{}
			for _, carrera := range carreras {
				carrerasOptions = append(carrerasOptions, carrera)
			}
			carreraSelect.SetOptions(carrerasOptions)
		}
	})
	facultadSelect.SetSelected("Facultad de Ciencias Médicas")

	botonCrearCertificado := createCertificadoComplexivoButton(myApp, name, selectStudy, entryPosgrado, selectorPersona, progressbar, facultadSelect, carreraSelect)
	cerrarButton := widget.NewButton("Cerrar", func() {
		myWindow.Close()
	})

	botonera := container.NewPadded(container.NewGridWithColumns(2,
		container.NewVBox(botonCrearCertificado),
		container.NewVBox(cerrarButton),
	))

	return container.NewPadded(container.NewVBox(
		card,
		layout.NewSpacer(),
		name,
		progressbar,
		layout.NewSpacer(),
		container.NewGridWithColumns(2,
			container.NewVBox(tituloSelect, selectStudy, entryPosgrado, tituloPersona, selectorPersona),
			container.NewVBox(titleFacultad, facultadSelect, tituloCarrera, carreraSelect),
		),
		layout.NewSpacer(),
		botonera,
	))
}

func createCertificadoComplexivoButton(myApp fyne.App, name *widget.Entry, selectStudy *widget.Select, entryPosgrado *widget.Entry, selectorPersona *widget.Select, progressbar *widget.ProgressBarInfinite, facultad *widget.Select, carrera *widget.Select) *widget.Button {
	return widget.NewButton("Crear Certificado", func() {
		// Obtener los datos ingresados
		nameEstudiante := name.Text
		estudio := selectStudy.Selected
		referencista := selectorPersona.Selected
		facultadSeleccionada := facultad.Selected
		carreraSeleccionada := fmt.Sprintf("de la carrera de %s", carrera.Selected)

		// Validar campos obligatorios
		if nameEstudiante == "" {
			showErrorWindow(myApp, fmt.Errorf("El nombre del estudiante no puede estar vacío"))
			progressbar.Hide()
			return
		}

		if facultadSeleccionada == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, selecciona una facultad"))
			progressbar.Hide()
			return
		}

		if carreraSeleccionada == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, selecciona una carrera"))
			progressbar.Hide()
			return
		}

		if selectStudy.Selected == "Posgrado" && entryPosgrado.Text == "" {
			showErrorWindow(myApp, fmt.Errorf("El campo de posgrado no puede estar vacío"))
			progressbar.Hide()
			return
		}

		// Si es posgrado, usar el texto ingresado en lugar del selector
		if selectStudy.Selected == "Posgrado" {
			estudio = entryPosgrado.Text
		}

		// Mostrar la barra de progreso
		progressbar.Show()

		// Crear el objeto del estudiante
		estudiante := scraping.Person{
			Author:   nameEstudiante,
			URI:      "", // No se utiliza en este caso
			Facultad: facultadSeleccionada,
			Carrera:  carreraSeleccionada,
		}

		// Canal para manejar errores
		errChan := make(chan error, 1)

		// Generar el certificado en un goroutine
		go func() {
			defer progressbar.Hide()

			// Crear el documento Word
			err := word.CreateWordDocument([]scraping.Person{estudiante}, estudio, referencista)
			if err != nil {
				errChan <- fmt.Errorf("Error al crear el documento: %w", err)
				return
			}

			// Agregar el registro a la base de datos
			_, err = database.AddRegistro(&database.Registro{
				Author:        estudiante.Author,
				Handle:        estudiante.URI,
				Facultad:      estudiante.Facultad,
				Carrera:       estudiante.Carrera,
				Fecha:         time.Now().Format("2006-01-02"),
				Bibliotecario: referencista,
			})
			if err != nil {
				errChan <- fmt.Errorf("Error al guardar el registro en la base de datos: %w", err)
				return
			}

			// Si todo es exitoso, enviar nil al canal de errores
			errChan <- nil
		}()

		// Manejar errores desde el canal
		go func() {
			if err := <-errChan; err != nil {
				showErrorWindow(myApp, err)
			} else {
				// Mostrar mensaje de éxito
				successWindow := myApp.NewWindow("Éxito")
				successLabel := widget.NewLabel("El certificado se generó correctamente y se guardó en la base de datos.")
				successWindow.SetContent(container.NewVBox(
					successLabel,
					widget.NewButton("Cerrar", func() {
						successWindow.Close()
					}),
				))
				successWindow.Resize(fyne.NewSize(300, 200))
				successWindow.Show()
			}
		}()
	})
}

// Crear la tarjeta de información
func createInfoCard() *fyne.Container {
	icono := widget.NewFileIcon(nil)
	contenidoCard := widget.NewLabel("Este programa se utiliza para generar certificados de no adeudar material bibliográfico a la biblioteca de la Universidad de Cuenca. Para generar un certificado solo debes poner la URL trabajo de titulación.")
	contenidoCard.Wrapping = fyne.TextWrapWord

	return container.NewGridWithColumns(2,
		icono,
		widget.NewCard("Generador de Certificados", "Versión: 2.0", contenidoCard),
	)
}

// Crear el botón para generar certificados
func createCertificadoButton(myApp fyne.App, uuid *widget.Entry, selectStudy *widget.Select, entryPosgrado *widget.Entry, selectorPersona *widget.Select, progressbar *widget.ProgressBarInfinite) *widget.Button {
	return widget.NewButton("Crear Certificado", func() {
		hdl := uuid.Text
		estudio := selectStudy.Selected
		referencista := selectorPersona.Selected

		if selectStudy.Selected == "Posgrado" && entryPosgrado.Text == "" {
			showErrorWindow(myApp, fmt.Errorf("el campo de posgrado no puede estar vacío"))
			progressbar.Hide()
			return
		}

		if selectStudy.Selected == "Posgrado" {
			estudio = entryPosgrado.Text
		}

		progressbar.Show()

		// Canal para manejar errores
		errChan := make(chan error, 1)

		go func() {
			defer progressbar.Hide()
			resultado, err := scraping.Scrapper(hdl, estudio)
			if err != nil {
				errChan <- err
				return
			}

			err = word.CreateWordDocument(resultado, estudio, referencista)
			if err != nil {
				errChan <- err
				return
			}

			for _, person := range resultado {
				_, err := database.AddRegistro(&database.Registro{
					Author:        person.Author,
					Handle:        person.URI,
					Facultad:      person.Facultad,
					Carrera:       person.Carrera,
					Fecha:         time.Now().Format("2006-01-02"),
					Bibliotecario: referencista,
				})
				if err != nil {
					errChan <- err
					return
				}
			}

			errChan <- nil
		}()

		// Manejar errores desde el canal
		go func() {
			if err := <-errChan; err != nil {
				showErrorWindow(myApp, err)
			}
		}()

		uuid.SetText("") // Limpia el campo de entrada
	})
}

// Mostrar ventana de error
func showErrorWindow(myApp fyne.App, err error) {
	errorWindow := myApp.NewWindow("Error")
	errorLabel := widget.NewLabel(fmt.Sprintf("Ocurrió un error:\n\n%s", err.Error()))
	errorLabel.Wrapping = fyne.TextWrapWord

	cerrarBoton := widget.NewButton("Cerrar", func() {
		errorWindow.Close()
	})

	errorContent := container.NewVBox(
		errorLabel,
		layout.NewSpacer(),
		cerrarBoton,
	)

	errorWindow.SetContent(errorContent)
	errorWindow.Resize(fyne.NewSize(300, 200))
	errorWindow.Show()
}

// Mostrar ventana de búsqueda de certificado
func searchCertificado(myApp fyne.App) {
	searchWindow := myApp.NewWindow("Buscar certificado")
	searchLabel := widget.NewLabel("Busca el certificado en la base de datos")
	entrySearch := widget.NewEntry()
	entrySearch.SetPlaceHolder("Introduce el handle, ejemplo: 15489")

	// Etiqueta de espacio
	espacio := widget.NewLabel("")

	// Crear una tabla vacía inicialmente
	var data [][]string
	data = append(data, []string{"Autor", "Handle", "Facultad", "Carrera", "Fecha", "Bibliotecario"}) // Encabezados

	// Variable para almacenar la fila seleccionada
	var filaSeleccionada []string

	// Crear la tabla
	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextTruncate // Truncar el texto si es demasiado largo
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		},
	)
	list.SetColumnWidth(0, 150) // Ajustar el ancho de las columnas
	list.SetColumnWidth(1, 150)
	list.SetColumnWidth(2, 200)
	list.SetColumnWidth(3, 200)
	list.SetColumnWidth(4, 150)
	list.SetColumnWidth(5, 200)

	// Manejar la selección de filas
	list.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 { // Ignorar la fila de encabezados
			filaSeleccionada = data[id.Row]
			fmt.Println("Fila seleccionada:", filaSeleccionada)
		}
	}

	scrollableList := container.NewScroll(list)
	scrollableList.SetMinSize(fyne.NewSize(700, 400))

	// Botón para buscar y actualizar la tabla
	botonBuscar := widget.NewButton("Buscar", func() {
		handle := entrySearch.Text // Obtener el texto ingresado
		if handle == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, introduce un handle para buscar"))
			return
		}

		// Consultar los datos en la base de datos
		registros, err := database.FetchbyQuery(handle, "handle")
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al buscar en la base de datos: %w", err))
			return
		}

		// Actualizar los datos de la tabla
		data = [][]string{{"Autor", "Handle", "Facultad", "Carrera", "Fecha", "Bibliotecario"}} // Reiniciar encabezados
		for _, registro := range registros {
			data = append(data, []string{
				registro.Author,
				registro.Handle,
				registro.Facultad,
				registro.Carrera,
				registro.Fecha,
				registro.Bibliotecario,
			})
		}

		// Refrescar la tabla
		list.Refresh()
	})

	// Botón para cerrar la ventana
	cerrarBoton := widget.NewButton("Cerrar", func() {
		searchWindow.Close()
	})

	// Botón para crear el certificado
	crearcertificadoButton := widget.NewButton("Crear Certificado", func() {
		// Validar que haya una fila seleccionada
		if filaSeleccionada == nil || len(filaSeleccionada) == 0 {
			showErrorWindow(myApp, fmt.Errorf("Por favor, selecciona una fila para crear el certificado"))
			return
		}

		// Crear el struct con los datos de la fila seleccionada
		registro := database.Registro{
			Author:        filaSeleccionada[0],
			Handle:        filaSeleccionada[1],
			Facultad:      filaSeleccionada[2],
			Carrera:       filaSeleccionada[3],
			Fecha:         filaSeleccionada[4],
			Bibliotecario: filaSeleccionada[5],
		}

		// Generar el documento de Word
		person := scraping.Person{
			Author:   registro.Author,
			URI:      registro.Handle,
			Facultad: registro.Facultad,
			Carrera:  registro.Carrera,
		}

		err := word.CreateWordDocument([]scraping.Person{person}, registro.Facultad, registro.Bibliotecario)
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al crear el certificado: %w", err))
			return
		}

		// Agregar el registro a la base de datos
		_, err = database.AddRegistro(&registro)
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al agregar el registro a la base de datos: %w", err))
			return
		}

		// Mostrar mensaje de éxito
		successWindow := myApp.NewWindow("Éxito")
		successLabel := widget.NewLabel("El certificado se generó correctamente y se agregó a la base de datos.")
		successWindow.SetContent(container.NewVBox(
			successLabel,
			widget.NewButton("Cerrar", func() {
				successWindow.Close()
			}),
		))
		successWindow.Resize(fyne.NewSize(300, 200))
		successWindow.Show()
	})

	// Contenedor para los botones con espaciado adecuado
	botonera := container.NewGridWithColumns(3,
		container.NewPadded(botonBuscar),
		container.NewPadded(cerrarBoton),
		container.NewPadded(crearcertificadoButton),
	)

	// Contenido de la ventana
	searchContent := container.NewBorder(
		container.NewVBox(searchLabel, entrySearch, espacio), // Parte superior
		botonera,       // Parte inferior
		nil,            // Izquierda
		nil,            // Derecha
		scrollableList, // Centro
	)

	searchWindow.SetContent(searchContent)
	searchWindow.Resize(fyne.NewSize(900, 600))
	searchWindow.Show()
}

// Exportar por fecha
func exportarPorFecha(myApp fyne.App) {
	exportWindow := myApp.NewWindow("Exportar por Fecha")
	exportLabel := widget.NewLabel("Introduce la fecha para exportar (Formato: YYYY, YYYY-MM o YYYY-MM-DD):")
	entryFecha := widget.NewEntry()
	entryFecha.SetPlaceHolder("Ejemplo: 2025, 2025-05, 2025-05-15")

	// Validar el formato de la fecha
	validarFecha := func(fecha string) bool {
		if len(fecha) == 4 { // Año
			return true
		} else if len(fecha) == 7 && fecha[4] == '-' { // Año-Mes
			return true
		} else if len(fecha) == 10 && fecha[4] == '-' && fecha[7] == '-' { // Año-Mes-Día
			return true
		}
		return false
	}
	// Botón para exportar los datos
	botonExportar := widget.NewButton("Exportar", func() {
		fecha := entryFecha.Text // Obtener el texto ingresado
		if fecha == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, introduce una fecha para exportar"))
			return
		}

		// Validar el formato de la fecha
		if !validarFecha(fecha) {
			showErrorWindow(myApp, fmt.Errorf("Formato de fecha no válido. Usa YYYY, YYYY-MM o YYYY-MM-DD"))
			return
		}

		// Consultar los datos en la base de datos
		registros, err := database.FetchbyQuery(fecha, "fecha")
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al buscar en la base de datos: %w", err))
			return
		}

		// Crear el archivo CSV
		fileName := fmt.Sprintf("export_%s.csv", fecha)
		file, err := os.Create(fileName)
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al crear el archivo CSV: %w", err))
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Escribir encabezados
		writer.Write([]string{"Autor", "Handle", "Facultad", "Carrera", "Fecha", "Bibliotecario"})

		// Escribir los registros
		for _, registro := range registros {
			writer.Write([]string{
				registro.Author,
				registro.Handle,
				registro.Facultad,
				registro.Carrera,
				registro.Fecha,
				registro.Bibliotecario,
			})
		}

		// Mostrar mensaje de éxito
		successWindow := myApp.NewWindow("Éxito")
		successLabel := widget.NewLabel(fmt.Sprintf("Datos exportados correctamente a %s", fileName))
		successWindow.SetContent(container.NewVBox(
			successLabel,
			widget.NewButton("Cerrar", func() {
				successWindow.Close()
			}),
		))
		successWindow.Resize(fyne.NewSize(300, 200))
		successWindow.Show()

		// Limpiar la entrada
		entryFecha.SetText("")
	})

	// Botón para cerrar la ventana
	botonCerrar := widget.NewButton("Cerrar", func() {
		exportWindow.Close()
	})

	// Contenedor para los botones
	botonera := container.NewGridWithColumns(2,
		container.NewPadded(botonExportar),
		container.NewPadded(botonCerrar),
	)

	// Contenido de la ventana
	exportContent := container.NewVBox(
		exportLabel,
		entryFecha,
		layout.NewSpacer(),
		botonera,
	)

	exportWindow.SetContent(exportContent)
	exportWindow.Resize(fyne.NewSize(400, 100))
	exportWindow.Show()
}

func createMasivosTab(myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	card := createInfoCard()

	fileLabel := widget.NewLabel("")
	fileLabel.Hide()
	fileButton := widget.NewButtonWithIcon("Seleccionar Archivo", theme.FileIcon(), func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorWindow(myApp, fmt.Errorf("error al abrir el archivo: %w", err))
				return
			}
			if reader == nil {
				return // El usuario canceló la selección
			}
			defer reader.Close()

			fileLabel.SetText(reader.URI().Path())
			fileLabel.Show()
		}, myWindow)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		// Set the title by showing a custom dialog with a label
		dialog.ShowCustom("Seleccionar archivo CSV", "Cerrar", container.NewVBox(), myWindow)
		fileDialog.Show()
	})

	fileSelector := container.NewGridWithColumns(2,
		container.NewPadded(fileButton),
		container.NewPadded(fileLabel),
	)

	// Selectores
	tituloSelect := widget.NewLabel("Selecciona el tipo de estudio:")
	selectStudy := widget.NewSelect([]string{"Pregrado", "Posgrado"}, nil)
	selectStudy.SetSelected("Pregrado")

	entryPosgrado := widget.NewEntry()
	entryPosgrado.SetPlaceHolder("Escribe el nombre del posgrado aquí...")
	entryPosgrado.Hide()

	selectStudy.OnChanged = func(selected string) {
		if selected == "Posgrado" {
			entryPosgrado.Show()
		} else {
			entryPosgrado.Hide()
		}
	}

	tituloPersona := widget.NewLabel("Referencista:")
	selectorPersona := widget.NewSelect(bibliotecarios, nil)

	// Botones
	progressbar := widget.NewProgressBarInfinite()
	progressbar.Hide()

	selectorPersona.SetSelected("PAOLA DEL ROCIO AMAYA ARCE")

	titleFacultad := widget.NewLabel("Selecciona la facultad:")
	tituloCarrera := widget.NewLabel("Selecciona la carrera:")
	carreraSelect := widget.NewSelect([]string{}, nil)

	// Obtener todas las facultades dinámicamente desde el mapa facultadComplexivo
	facultades := []string{}
	for facultad := range facultadComplexivo {
		facultades = append(facultades, facultad)
	}

	// Crear el selector de facultades
	facultadSelect := widget.NewSelect(facultades, func(facultadSeleccionada string) {
		if carreras, ok := facultadComplexivo[facultadSeleccionada]; ok {
			carrerasOptions := []string{}
			for _, carrera := range carreras {
				carrerasOptions = append(carrerasOptions, carrera)
			}
			carreraSelect.SetOptions(carrerasOptions)
		}
	})
	facultadSelect.SetSelected("Facultad de Ciencias Médicas")

	buttonCertificado := certificadosMasivosButton(myApp, entryPosgrado, selectStudy, entryPosgrado, selectorPersona, progressbar, facultadSelect, carreraSelect, fileLabel)
	botonera := container.NewPadded(container.NewGridWithColumns(2,
		container.NewVBox(buttonCertificado),
		container.NewVBox(widget.NewButton("Cerrar", func() {
			myWindow.Close()
		})),
	))

	return container.NewVBox(
		card,
		layout.NewSpacer(),
		fileSelector,
		layout.NewSpacer(),
		container.NewGridWithColumns(2,
			container.NewVBox(tituloSelect, selectStudy, entryPosgrado, tituloPersona, selectorPersona),
			container.NewVBox(titleFacultad, facultadSelect, tituloCarrera, carreraSelect),
		),
		layout.NewSpacer(),
		botonera,
	)
}

func certificadosMasivosButton(myApp fyne.App, name *widget.Entry, selectStudy *widget.Select, entryPosgrado *widget.Entry, selectorPersona *widget.Select, progressbar *widget.ProgressBarInfinite, facultad *widget.Select, carrera *widget.Select, fileURI *widget.Label) *widget.Button {
	return widget.NewButton("Crear Certificados Masivos", func() {
		// Obtener los datos ingresados
		file := fileURI.Text
		estudio := selectStudy.Selected
		referencista := selectorPersona.Selected
		facultadSeleccionada := facultad.Selected
		carreraSeleccionada := fmt.Sprintf("de la carrera de %s", carrera.Selected)

		// Validar campos obligatorios
		if file == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, selecciona un archivo CSV"))
			progressbar.Hide()
			return
		}

		if facultadSeleccionada == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, selecciona una facultad"))
			progressbar.Hide()
			return
		}

		if carreraSeleccionada == "" {
			showErrorWindow(myApp, fmt.Errorf("Por favor, selecciona una carrera"))
			progressbar.Hide()
			return
		}

		if selectStudy.Selected == "Posgrado" && entryPosgrado.Text == "" {
			showErrorWindow(myApp, fmt.Errorf("El campo de posgrado no puede estar vacío"))
			progressbar.Hide()
			return
		}

		// Si es posgrado, usar el texto ingresado en lugar del selector
		if selectStudy.Selected == "Posgrado" {
			estudio = entryPosgrado.Text
		}

		// Mostrar barra de progreso
		progressbar.Show()

		// Leer el archivo CSV
		fileHandle, err := os.Open(file)
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al abrir el archivo: %w", err))
			progressbar.Hide()
			return
		}
		defer fileHandle.Close()

		reader := csv.NewReader(fileHandle)

		// Leer todas las filas del archivo
		rows, err := reader.ReadAll()
		if err != nil {
			showErrorWindow(myApp, fmt.Errorf("Error al leer el archivo CSV: %w", err))
			progressbar.Hide()
			return
		}

		// Crear el slice de personas
		personas := []scraping.Person{}

		// Procesar las filas
		for i, row := range rows {
			// Saltar la primera fila si es un encabezado
			if i == 0 {
				continue
			}

			// Validar que la fila tenga al menos un nombre
			if len(row) < 1 || row[0] == "" {
				showErrorWindow(myApp, fmt.Errorf("La fila %d no tiene un nombre válido", i+1))
				progressbar.Hide()
				return
			}

			// Crear un objeto scraping.Person
			personas = append(personas, scraping.Person{
				Author:   row[0],               // Nombre del estudiante
				Facultad: facultadSeleccionada, // Facultad seleccionada
				Carrera:  carreraSeleccionada,  // Carrera seleccionada
				URI:      "",                   // No se utiliza en este caso
			})
		}

		// Canal para manejar errores
		errChan := make(chan error, 1)

		// Generar certificados en un goroutine
		go func() {
			defer progressbar.Hide()

			// Crear los certificados
			err := word.CreateWordDocument(personas, estudio, referencista)
			if err != nil {
				errChan <- fmt.Errorf("Error al crear los certificados: %w", err)
				return
			}

			// Agregar los registros a la base de datos
			for _, persona := range personas {
				_, err := database.AddRegistro(&database.Registro{
					Author:        persona.Author,
					Handle:        persona.URI,
					Facultad:      persona.Facultad,
					Carrera:       persona.Carrera,
					Fecha:         time.Now().Format("2006-01-02"),
					Bibliotecario: referencista,
				})
				if err != nil {
					errChan <- fmt.Errorf("Error al guardar el registro en la base de datos: %w", err)
					return
				}
			}

			// Si todo es exitoso, enviar nil al canal de errores
			errChan <- nil
		}()

		// Manejar errores desde el canal
		go func() {
			if err := <-errChan; err != nil {
				showErrorWindow(myApp, err)
			} else {
				// Mostrar mensaje de éxito
				successWindow := myApp.NewWindow("Éxito")
				successLabel := widget.NewLabel("Los certificados se generaron correctamente y se guardaron en la base de datos.")
				successWindow.SetContent(container.NewVBox(
					successLabel,
					widget.NewButton("Cerrar", func() {
						successWindow.Close()
					}),
				))
				successWindow.Resize(fyne.NewSize(300, 200))
				successWindow.Show()
			}
		}()
	})
}

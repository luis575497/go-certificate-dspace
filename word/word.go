package word

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"generador-certificados/scraping" // Paquete scraping

	"github.com/fumiama/go-docx" // Paquete para trabajar con documentos .docx
)

func CreateWordDocument(persons []scraping.Person, estudio string, referencista string) error {
	// Crear el directorio de salida si no existe
	err := os.MkdirAll("./certificados", os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear el directorio de certificados: %w", err)
	}

	// Obtener la fecha actual y formatearla en español
	currentDate := time.Now()
	dateStr := currentDate.Format("02 de January de 2006")
	monthMap := map[string]string{
		"January":   "enero",
		"February":  "febrero",
		"March":     "marzo",
		"April":     "abril",
		"May":       "mayo",
		"June":      "junio",
		"July":      "julio",
		"August":    "agosto",
		"September": "septiembre",
		"October":   "octubre",
		"November":  "noviembre",
		"December":  "diciembre",
	}
	for eng, esp := range monthMap {
		dateStr = strings.ReplaceAll(dateStr, eng, esp)
	}
	fullDate := "Cuenca, " + dateStr

	// Para cada estudiante, se crea un documento independiente
	for _, person := range persons {
		// Crear un nuevo documento para cada estudiante
		w := docx.New().WithDefaultTheme()

		// Crear una tabla con 1 fila y 2 columnas
		encabezado := w.AddTable(1, 2, 9000, nil) // 1 fila, 2 columnas, ancho ajustado

		// Eliminar bordes usando métodos de la librería (si están disponibles)
		encabezado.TableProperties.TableBorders = nil // Eliminar bordes de la tabla

		// Celda izquierda: Imagen
		celdaImagen := encabezado.TableRows[0].TableCells[0] // Primera fila, primera columna
		pImage := celdaImagen.AddParagraph()
		r, err := pImage.AddInlineDrawingFrom("logoucuenca.png")
		if err != nil {
			log.Fatalf("Error al añadir la imagen: %v", err)
		}

		// Ajustar el tamaño de la imagen manteniendo la proporción
		drawing, ok := r.Children[0].(*docx.Drawing)
		if !ok {
			log.Fatalf("El primer hijo no es un Drawing")
		}
		if drawing.Inline == nil {
			log.Fatalf("El Drawing no tiene Inline definido")
		}

		// Establecer el tamaño de la imagen (ancho x alto en EMUs)
		ancho := int64(6 * 360000)           // 6 cm
		alto := int64(float64(ancho) / 3.33) // Mantener proporción original (1754x526)
		drawing.Inline.Size(ancho, alto)
		// Celda derecha: Texto
		celdaTexto := encabezado.TableRows[0].TableCells[1] // Primera fila, segunda columna
		p1 := celdaTexto.AddParagraph()
		p1.Justification("end") // Justificar al extremo derecho
		p1.AddText("FORMATO DE NO ADEUDAR MATERIAL BIBLIOGRÁFICO A LA BIBLIOTECA").
			Bold().
			Size("15").
			Font("Arial", "Arial", "Arial", "Arial")
		p2 := celdaTexto.AddParagraph()
		p2.Justification("end") // Justificar al extremo derecho
		p2.AddText("UC-CDRJVB-FOR-020").
			Size("15").
			Font("Arial", "Arial", "Arial", "Arial")
		p3 := celdaTexto.AddParagraph()
		p3.Justification("end") // Justificar al extremo derecho
		p3.AddText("Página 1 de 1").
			Size("15").
			Font("Arial", "Arial", "Arial", "Arial")

		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")

		// 1. Título centrado: "C E R T I F I C A"
		pTitle := w.AddParagraph()
		pTitle.Justification("center")
		pTitle.AddText("CERTIFICADO DE NO ADEUDAR").
			Bold().
			Size("24").
			Font("Arial", "Arial", "Arial", "Arial")
		// Espacio adicional
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")

		// 2. Cuerpo del certificado con datos del estudiante
		pBody := w.AddParagraph()
		pBody.Justification("both")

		// Texto completo en un solo bloque
		pBody.AddText(fmt.Sprintf(
			"El Centro de Documentación Regional \"Juan Bautista Vázquez\" certifica que %s, portador de la cédula de ciudadanía No. XXXXXXXX, estudiante %s, de la %s, no adeuda ningún bien, ni material bibliográfico en esta dependencia.",
			strings.ToUpper(person.Author),
			person.Carrera,
			person.Facultad,
		)).
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")

		// Espaciado posterior (3 líneas vacías)
		w.AddParagraph().AddText("")
		w.AddParagraph().AddText("")
		w.AddParagraph().AddText("")

		// 3. Fecha (se usa la fecha actual generada)
		pDate := w.AddParagraph()
		pDate.Justification("end")
		pDate.AddText(fullDate).
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")

		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")

		// 4. Bloque de firma (referencista)
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial") // Espacio
		pAtt := w.AddParagraph()
		pAtt.Justification("center")
		pAtt.AddText("Atentamente,").
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")

		pLine := w.AddParagraph()
		pLine.Justification("center")
		pLine.AddText("________________________________________").
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		pRef := w.AddParagraph()
		pRef.Justification("center")
		pRef.AddText(strings.ToTitle(referencista)).
			Bold().
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")
		pRole := w.AddParagraph()
		pRole.Justification("center")
		pRole.AddText("Bibliotecario 2").
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")
		pCDR := w.AddParagraph()
		pCDR.Justification("center")
		pCDR.AddText("CDR \"Juan Bautista Vázquez\"").
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")

		// 5. Enlace al repositorio
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		pLink := w.AddParagraph()
		pLink.Justification("start")
		pLink.AddText(fmt.Sprintf("Link: %s", person.URI)).
			Size("22").
			Font("Arial", "Arial", "Arial", "Arial")

		// Control de Versiones
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")
		w.AddParagraph().AddText("").
			Font("Arial", "Arial", "Arial", "Arial")

		pversion := w.AddParagraph()
		pversion.Justification("end")
		pversion.AddText("Version: 2.0").
			Size("15").
			Font("Arial", "Arial", "Arial", "Arial")

		// 6. Definir el nombre del archivo usando el nombre del estudiante (Author)
		fileName := fmt.Sprintf("%s_certificado.docx", strings.ReplaceAll(person.Author, " ", "_"))
		outputFile := fmt.Sprintf("./certificados/%s", fileName)

		// 6.1 Verificar si el archivo ya existe y eliminarlo si es necesario
		if _, err := os.Stat(outputFile); err == nil {
			err = os.Remove(outputFile)
			if err != nil {
				return fmt.Errorf("no se pudo eliminar el archivo existente para %s: %w", person.Author, err)
			}
		}

		// 7. Crear el archivo y escribir el documento
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("error al crear el archivo para %s: %w", person.Author, err)
		}
		_, err = w.WriteTo(f)
		f.Close()
		if err != nil {
			return fmt.Errorf("error al escribir en el archivo para %s: %w", person.Author, err)
		}

		fmt.Println("Documento creado exitosamente para:", person.Author)
	}
	return nil
}

package scraping

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

type Person struct {
	Author   string
	URI      string
	Facultad string
	Carrera  string
}

// Scrapper hace scraping de la URL y devuelve una lista de Person
func Scrapper(url string, estudio string) ([]Person, error) {
	// 1. Crea un contexto base con timeout
	baseCtx, cancelBase := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelBase()

	// 2. Crea un contexto de navegador a partir del anterior
	ctx, cancelBrowser := chromedp.NewContext(baseCtx)
	defer cancelBrowser()

	var tbodyHTML string
	var breadcrumbHTML string

	// 3. Ejecuta las tareas dentro del contexto con timeout
	err := chromedp.Run(ctx,
		chromedp.Navigate(url+"/full"),
		chromedp.Sleep(2*time.Second), // mejor aún: usa WaitVisible si sabes qué esperar
		chromedp.OuterHTML("tbody", &tbodyHTML),
		chromedp.OuterHTML("ol.container.breadcrumb.my-0", &breadcrumbHTML),
	)

	// 4. Maneja los errores
	if err != nil {
		if baseCtx.Err() == context.DeadlineExceeded {
			return nil, errors.New("el scraping tardó demasiado (timeout)")
		}
		return nil, errors.New("error al ejecutar chromedp: " + err.Error())
	}

	if tbodyHTML == "" {
		return nil, errors.New("no se encontró el elemento <tbody> en la página")
	}

	if breadcrumbHTML == "" {
		return nil, errors.New("no se encontró el breadcrumb en la página")
	}
	facultad, carrera := parseBreadcrumb(breadcrumbHTML)

	if estudio != "Pregrado" {
		carrera = fmt.Sprintf("de la Maestría de %s", estudio)
	} else {
		carrera = fmt.Sprintf("de la Carrera de %s", carrera)
	}

	persons := parseTable(tbodyHTML, facultad, carrera)

	if len(persons) == 0 {
		return persons, errors.New("no se encontraron personas en la tabla")
	}

	return persons, nil
}

// Ananliza el Breadcrumb y obten los valores de facultad y carrera
func parseBreadcrumb(breadcrumbHTML string) (string, string) {
	doc, err := html.Parse(strings.NewReader(breadcrumbHTML))
	if err != nil {
		fmt.Println("Error al parsear HTML:", err)
		return "", ""
	}

	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Buscar nodo de texto dentro del <a>
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					text := strings.TrimSpace(c.Data)
					if text != "" {
						links = append(links, text)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	// Mostrar todos los textos encontrados en los <a>
	fmt.Println("Links:", links)

	var facultad, carrera string
	if len(links) >= 3 {
		facultad = links[len(links)-3]
		carrera = links[len(links)-2]
	}

	return facultad, carrera
}

func parseTable(tbodyHTML string, facultad string, carrera string) []Person {
	body := strings.NewReader(tbodyHTML)
	z := html.NewTokenizer(body)

	var uri string
	var authors []string
	var currentRow []string
	var inRow bool

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		switch {
		case tt == html.StartTagToken && t.Data == "tr":
			inRow = true
			currentRow = nil
		case tt == html.EndTagToken && t.Data == "tr":
			inRow = false
			if len(currentRow) >= 2 {
				key := currentRow[0]
				value := currentRow[1]
				if key == "dc.identifier.uri" {
					uri = value
				} else if key == "dc.contributor.author" {
					authors = append(authors, value)
				}
			}
		case inRow && tt == html.StartTagToken && t.Data == "td":
			z.Next()
			text := strings.TrimSpace(string(z.Text()))
			currentRow = append(currentRow, text)
		}
	}

	var persons []Person
	for _, author := range authors {
		persons = append(persons, Person{Author: author, URI: uri, Facultad: facultad, Carrera: carrera})
	}

	return persons
}

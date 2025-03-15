package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"
)

func main() {
	port := 8080
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/calculate", handleCalculation)
	log.Printf("Server started on :%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../../static/templates/form.html"))
	tmpl.Execute(w, nil)
	//123
}

func handleCalculation(w http.ResponseWriter, r *http.Request) {
	// Обработка данных формы
	r.ParseForm()
	birthDate := r.Form.Get("birthdate")
	birthTime := r.Form.Get("birthtime")
	location := r.Form.Get("location")

	// Геокодирование (упрощенно)
	lat, lon := geocode(location)

	// Расчет позиций планет (заглушка)
	positions := calculatePlanetaryPositions(birthDate, birthTime, lat, lon)

	// Генерация бодиграфа
	bodygraph := generateBodygraph(positions)

	// Отображение результатов
	tmpl := template.Must(template.ParseFiles("../../static/templates/result.html"))
	data := struct {
		SVG         template.HTML
		Description string
	}{
		SVG:         template.HTML(generateSVG(bodygraph)),
		Description: getBodygraphDescription(bodygraph),
	}
	tmpl.Execute(w, data)
}

func geocode(location string) (float64, float64) {
	//Реальная реализация будет использовать API геокодирования
	return 52.5200, 13.4050 // Пример: координаты Берлина
}

func calculatePlanetaryPositions(date, timeStr string, lat, lon float64) map[string]float64 {
	// Заглушка: реальная реализация использует Swiss Ephemeris
	return map[string]float64{
		"Sun":  25.4,
		"Moon": 12.7,
	}
}

func generateBodygraph(positions map[string]float64) map[int]bool {
	// Заглушка: преобразование позиций в активированные ворота
	return map[int]bool{
		1:  true,
		5:  true,
		10: true,
	}
}

func generateSVG(bodygraph map[int]bool) string {
	// Упрощенный SVG с примерами элементов
	//return fmt.Sprintf(`
	//<svg width="400" height="400">
	//<circle cx="200" cy="100" r="20" fill="%s"/>
	//<circle cx="200" cy="300" r="20" fill="%s"/>
	//<line x1="200" y1="100" x2="200" y2="300" stroke="black"/>
	//</svg>
	//`, getColor(bodygraph[1]), getColor(bodygraph[5]))

	//Первая версия расширенного
	var sb strings.Builder

	// Стили и фон
	sb.WriteString(`<svg viewBox="0 0 400 680" xmlns="http://www.w3.org/2000/svg">`)
	sb.WriteString(`<rect width="100%" height="100%" fill="#fff"/>`)
	log.Println(sb.String())
	log.Println("'1'")

	// Определение центров и каналов
	centers := map[string]struct {
		X, Y  float64
		Shape string
		Gates []int
	}{
		"Head":   {200, 80, "circle", []int{64, 61, 63}},
		"Ajna":   {200, 150, "circle", []int{47, 24, 4}},
		"Throat": {200, 220, "square", []int{62, 23, 56, 35, 12, 45, 33, 8, 31, 20, 16}},
		"G":      {200, 350, "diamond", []int{1, 2, 7, 10, 13, 15, 25}},
		"Sacral": {200, 480, "circle", []int{5, 14, 29, 9, 3, 42, 27}},
	}

	channels := []struct {
		Path string
	}{
		{Path: "M200,105 L200,145"},
		{Path: "M200,170 L200,210"},
		{Path: "M200,240 L200,340"},
	}

	// Отрисовка каналов
	for _, ch := range channels {
		sb.WriteString(fmt.Sprintf(
			`<path d="%s" stroke="#eee" stroke-width="8"/>`,
			ch.Path,
		))
	}

	// Отрисовка центров
	for name, center := range centers {
		switch center.Shape {
		case "circle":
			sb.WriteString(fmt.Sprintf(
				`<circle cx="%.0f" cy="%.0f" r="35" fill="#fff" stroke="#ddd" stroke-width="2"/>`,
				center.X, center.Y,
			))
		case "square":
			sb.WriteString(fmt.Sprintf(
				`<rect x="%.0f" y="%.0f" width="70" height="70" fill="#fff" stroke="#ddd" stroke-width="2"/>`,
				center.X-35, center.Y-35,
			))
		}
		sb.WriteString(fmt.Sprintf(
			`<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle">%s</text>`,
			center.X, center.Y+5, name,
		))
	}

	// Отрисовка ворот
	for _, center := range centers {
		totalGates := len(center.Gates)
		for i, gate := range center.Gates {
			angle := float64(i) * (2 * math.Pi / float64(totalGates))
			x := center.X + 70*math.Cos(angle)
			y := center.Y + 70*math.Sin(angle)

			color := "#ff4444"
			if !bodygraph[gate] {
				color = "#eee"
			}

			sb.WriteString(fmt.Sprintf(
				`<circle cx="%.1f" cy="%.1f" r="15" fill="%s" stroke="#ddd"/>`,
				x, y, color,
			))
			sb.WriteString(fmt.Sprintf(
				`<text x="%.1f" cy="%.1f" font-size="12" text-anchor="middle" dy=".3em">%d</text>`,
				x, y, gate,
			))
		}
	}

	sb.WriteString("</svg>")
	log.Println(sb.String())
	log.Println("'2'")
	return sb.String()
}

func getColor(active bool) string {
	if active {
		return "red"
	}
	return "gray"
}

func getBodygraphDescription(bodygraph map[int]bool) string {
	// Пример описания
	return "Активированы ворота 1, 5 и 10. Это указывает на..."
}

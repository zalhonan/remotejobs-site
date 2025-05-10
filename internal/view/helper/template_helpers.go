package helper

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// TemplateFuncs возвращает карту функций для использования в шаблонах
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"truncate":     truncate,
		"formatDate":   formatDate,
		"safeHTML":     safeHTML,
		"iterate":      iterate,
		"add":          add,
		"subtract":     subtract,
		"escapeJS":     escapeJS,
		"formatNumber": formatNumber,
		"split":        strings.Split,
	}
}

// truncate обрезает строку до указанной длины и добавляет многоточие, если строка была обрезана
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	// Ищем последний пробел перед maxLen
	lastSpace := strings.LastIndex(s[:maxLen], " ")
	if lastSpace > 0 {
		return s[:lastSpace] + "..."
	}

	return s[:maxLen] + "..."
}

// formatDate форматирует время в читаемую строку
func formatDate(t time.Time) string {
	return t.Format("02.01.2006")
}

// safeHTML помечает строку как безопасный HTML
func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

// iterate создает слайс целых чисел от start до end (включительно)
// Полезно для создания циклов в шаблонах
func iterate(start, end int) []int {
	if start > end {
		return []int{}
	}

	result := make([]int, end-start+1)
	for i := range result {
		result[i] = start + i
	}

	return result
}

// add добавляет два числа
func add(a, b int) int {
	return a + b
}

// subtract вычитает b из a
func subtract(a, b int) int {
	return a - b
}

// escapeJS экранирует строку для безопасного использования в JavaScript
func escapeJS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")

	return s
}

// formatNumber форматирует число с разделителями
func formatNumber(n int) string {
	return fmt.Sprintf("%d", n)
}

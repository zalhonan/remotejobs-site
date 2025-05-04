package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/zalhonan/remotejobs-site/internal/view/helper"
	"go.uber.org/zap"
)

// TemplateRenderer отвечает за рендеринг HTML шаблонов
type TemplateRenderer struct {
	templates    map[string]*template.Template
	templateDir  string
	baseTemplate string
	logger       *zap.Logger
}

// NewTemplateRenderer создает новый рендерер шаблонов
func NewTemplateRenderer(templateDir, baseTemplate string, logger *zap.Logger) (*TemplateRenderer, error) {
	renderer := &TemplateRenderer{
		templates:    make(map[string]*template.Template),
		templateDir:  templateDir,
		baseTemplate: baseTemplate,
		logger:       logger,
	}

	if err := renderer.precompileTemplates(); err != nil {
		return nil, fmt.Errorf("не удалось предварительно скомпилировать шаблоны: %w", err)
	}

	return renderer, nil
}

// precompileTemplates предварительно компилирует шаблоны при инициализации
func (tr *TemplateRenderer) precompileTemplates() error {
	// Шаблоны для страниц
	pageTemplates := []string{
		"pages/home.html",
		"pages/job_details.html",
		"errors/error.html",
	}

	// Общие компоненты
	commonTemplates := []string{
		"layout/components/header.html",
		"layout/components/footer.html",
		"layout/components/pagination.html",
	}

	// Базовый шаблон
	basePath := filepath.Join(tr.templateDir, tr.baseTemplate)

	// Загружаем каждый шаблон страницы
	for _, page := range pageTemplates {
		// Создаем новый шаблон с функциями
		tmpl := template.New(filepath.Base(page)).Funcs(helper.TemplateFuncs())

		// Загружаем базовый шаблон
		tmpl, err := tmpl.ParseFiles(basePath)
		if err != nil {
			return fmt.Errorf("не удалось загрузить базовый шаблон %s: %w", basePath, err)
		}

		// Загружаем общие компоненты
		for _, component := range commonTemplates {
			componentPath := filepath.Join(tr.templateDir, component)
			tmpl, err = tmpl.ParseFiles(componentPath)
			if err != nil {
				return fmt.Errorf("не удалось загрузить компонент %s: %w", componentPath, err)
			}
		}

		// Загружаем шаблон страницы
		pagePath := filepath.Join(tr.templateDir, page)
		tmpl, err = tmpl.ParseFiles(pagePath)
		if err != nil {
			return fmt.Errorf("не удалось загрузить шаблон страницы %s: %w", pagePath, err)
		}

		// Сохраняем шаблон в кэше
		tr.templates[page] = tmpl
	}

	return nil
}

// Render рендерит шаблон с заданными данными
func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := tr.templates[name]
	if !ok {
		return fmt.Errorf("шаблон %s не найден", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return tmpl.ExecuteTemplate(w, "base", data)
}

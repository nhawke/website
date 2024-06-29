package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	recipes "github.com/nhawke/recipe-handler"
)

var (
	recipeFolderPath = flag.String("recipe_path", "", "path to recipe folder")

	//go:embed pages
	pagesFS embed.FS

	baseTemplate *template.Template
)

func init() {
}

func main() {
	flag.Parse()

	http.HandleFunc("/", servePage)

	recipeHandler := &recipes.Handler{
		Path:               *recipeFolderPath,
		RecipePageTemplate: pageTemplate("recipes.html.tmpl"),
		RecipeListTemplate: recipes.DefaultRecipeListTemplate,
	}
	http.Handle("/recipes/", http.StripPrefix("/recipes", recipeHandler))

	fmt.Fprintln(os.Stderr, "Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func pageTemplate(pageName string) *template.Template {
	tmpl, err := template.ParseFS(pagesFS, "pages/base.html.tmpl", "pages/"+pageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find template at path: pages/%v\n", pageName)
		return nil
	}
	return tmpl
}

func servePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stderr, "Serving page: %v\n", req.URL.Path)

	pageName := strings.TrimPrefix(req.URL.Path, "/")
	if pageName == "" {
		pageName = "index"
	}
	pageName += ".html.tmpl"
	fmt.Fprintf(os.Stderr, "Page %v resolved to %v\n", req.URL.Path, pageName)

	t := pageTemplate(pageName)
	if t == nil {
		http.NotFound(w, req)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

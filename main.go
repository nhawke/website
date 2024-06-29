package main

import (
	"bytes"
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
	//go:embed layouts
	templatesFS embed.FS

	baseTemplate *template.Template
)

func init() {
}

func main() {
	flag.Parse()

	http.HandleFunc("/", servePage)

	recipePageTmpl := template.Must(template.ParseFS(templatesFS, "layouts/base.html.tmpl", "layouts/recipe.html.tmpl"))
	recipeListTmpl := template.Must(template.ParseFS(templatesFS, "layouts/base.html.tmpl", "layouts/recipe_list.html.tmpl"))

	recipeHandler := &recipes.Handler{
		Path:               *recipeFolderPath,
		RecipePageTemplate: recipePageTmpl,
		RecipeListTemplate: recipeListTmpl,
	}
	http.Handle("/recipes/", http.StripPrefix("/recipes", recipeHandler))

	fmt.Fprintln(os.Stderr, "Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func servePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stderr, "Serving page: %v\n", req.URL.Path)

	pageName := strings.TrimPrefix(req.URL.Path, "/")
	if pageName == "" {
		pageName = "index"
	}
	pageName += ".html.tmpl"
	fmt.Fprintf(os.Stderr, "Page %v resolved to %v\n", req.URL.Path, pageName)

	tmpl, err := template.ParseFS(templatesFS, "layouts/base.html.tmpl", "pages/"+pageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find template at path: pages/%v\n", pageName)
	}
	if tmpl == nil {
		http.NotFound(w, req)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template for page %v: %v\n", pageName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(buf.Bytes())
}

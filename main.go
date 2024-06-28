package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	recipes "github.com/nhawke/recipe-handler"
)

var (
	recipeFolderPath = flag.String("recipe_path", "", "path to recipe folder")

	//go:embed pages
	pages embed.FS
)

func main() {
	flag.Parse()

	recipeHandler := recipes.NewHandler(recipes.Config{
		Path: *recipeFolderPath,
	})
	http.Handle("/recipes/", http.StripPrefix("/recipes", recipeHandler))

	root, err := fs.Sub(pages, "pages")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	http.Handle("/", http.FileServerFS(root))

	fmt.Fprintln(os.Stderr, "Listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

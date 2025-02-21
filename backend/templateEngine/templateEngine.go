package templateengine

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/robrotheram/gogallery/backend/embeds"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

const HomeTemplate = "index"
const AlbumTemplate = "albums"
const CollectionTemplate = "collections"
const PhotoTemplate = "photo"
const PaginationTemplate = "pagination"

func (te *TemplateEngine) LoadFromEmbed() error {

	te.Cache = newTeamplateCache()

	base, err := template.ParseFS(embeds.ThemeFS, "eastnor/default.hbs")
	if err != nil {
		return err
	}
	base.Funcs(template.FuncMap{"ImgSizes": func() map[string]int { return ImageSizes }})
	base, _ = base.ParseFS(embeds.ThemeFS, "eastnor/partials/*.hbs")

	items, err := embeds.ThemeFS.ReadDir("eastnor/pages")
	if err != nil {
		return err
	}
	for _, item := range items {
		name := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
		pageTemplate := template.Must(base.Clone())
		pageTemplate = template.Must(pageTemplate.ParseFS(embeds.ThemeFS, "eastnor/pages/"+item.Name()))
		te.Cache.Add(name, pageTemplate)
	}
	return nil
}

func (te *TemplateEngine) Load(basePath string) error {
	if basePath == "default" {
		return te.LoadFromEmbed()
	}
	return te.LoadFromPath(basePath)
}

func (te *TemplateEngine) LoadFromPath(basePath string) error {
	te.Cache = newTeamplateCache()
	pagePath := "pages"
	base, err := template.ParseFiles(filepath.Join(basePath, "default.hbs"))
	if err != nil {
		return err
	}
	base.Funcs(template.FuncMap{"ImgSizes": func() map[string]int { return ImageSizes }})
	base, err = base.ParseGlob(filepath.Join(basePath, "partials/*.hbs"))
	if err != nil {
		return err
	}
	items, err := ioutil.ReadDir(filepath.Join(basePath, pagePath))
	if err != nil {
		return err
	}
	for _, item := range items {
		name := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
		pageTemplate := template.Must(base.Clone())
		pageTemplate = template.Must(pageTemplate.ParseGlob(filepath.Join(basePath, pagePath, item.Name())))
		te.Cache.Add(name, pageTemplate)
	}
	return nil
}

type TemplateEngine struct {
	Cache *TemplateCache
	m     *minify.M
}

var Templates = NewTemplateEngine()

func NewTemplateEngine() *TemplateEngine {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	return &TemplateEngine{
		m: m,
	}
}

func (te *TemplateEngine) RenderPage(w io.Writer, pageName string, data Page) {
	var tpl bytes.Buffer
	err := te.Cache.Get(pageName).Execute(&tpl, data)
	if err != nil {
		fmt.Println(err)
	}
	b, err := te.m.Bytes("text/html", tpl.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	w.Write(b)
}

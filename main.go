package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"os"
	"strings"
	"regexp"
	"path/filepath"
)

type FuncMap map[string]interface{}

var AlertManagerFuncs = FuncMap{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title":   strings.Title,
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"match": regexp.MatchString,
	"safeHtml": func(text string) template.HTML {
		return template.HTML(text)
	},
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
}

func main() {
	var (
		flTemplate       = flag.String("template", "", "path to template")
		flJSON           = flag.String("json", "", "json input to use in template")
		flTemplateString = flag.Bool("string", false, "template is a string, not a file")
		flEnableAlertmanager = flag.Bool("alertmanager", false, "enable alertmanager template functions")
	)
	flag.Parse()

	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(*flJSON), &dict); err != nil {
		log.Fatal(err)
	}

	var (
		err  error
		tmpl *template.Template
		name string
	)

	name = "template"
	tmpl = template.New(name)

	if *flEnableAlertmanager {
		tmpl = tmpl.Funcs(template.FuncMap(AlertManagerFuncs))
	}

	if *flTemplateString {
		tmpl, err = tmpl.Parse(*flTemplate)
	} else {
		name = filepath.Base(*flTemplate)
		tmpl, err = tmpl.ParseFiles(*flTemplate)
	}


	if err != nil {
		log.Fatal(err)
	}

	if err := tmpl.ExecuteTemplate(os.Stdout, name, dict); err != nil {
		log.Fatal(err)
	}
}

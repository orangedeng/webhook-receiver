package tmpl

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"
)

func ExecuteTextString(data interface{}, notificationTmpl string) (string, error) {
	if notificationTmpl == "" {
		return "", nil
	}

	tmpl := template.New("").Option("missingkey=zero")
	tmpl.Funcs(template.FuncMap(DefaultFuncs))
	tpl, err := tmpl.Parse(notificationTmpl)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	if err := tpl.ExecuteTemplate(&buf, "title.text.list", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type FuncMap map[string]interface{}

var DefaultFuncs = FuncMap{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title":   strings.Title,
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"match": regexp.MatchString,
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
	"stringSlice": func(s ...string) []string {
		return s
	},
}

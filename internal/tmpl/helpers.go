package tmpl

import (
	"bytes"
	"html/template"
	"net/http"

	i18nctx "tournament_manager/internal/i18n"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ctxKey string

const localizerKey ctxKey = "localizer"

func GetLocalizer(r *http.Request) *i18n.Localizer {
	return i18nctx.GetLocalizer(r.Context())
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"T": func(msgID string, data ...map[string]interface{}) string {
			return msgID // dummy, will be replaced per request
		},
		"dict": func(values ...string) map[string]string {
			m := map[string]string{}
			for i := 0; i < len(values); i += 2 {
				m[values[i]] = values[i+1]
			}
			return m
		},
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data interface{}) error {
	loc := GetLocalizer(r)
	t, err := tmpl.Clone()
	if err != nil {
		return err
	}

	t = t.Funcs(template.FuncMap{
		"T": func(msgID string, data ...map[string]interface{}) string {
			if loc == nil {
				return msgID
			}
			str, _ := loc.Localize(&i18n.LocalizeConfig{
				MessageID:    msgID,
				TemplateData: firstMap(data),
			})
			return str
		},
	})
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}

func firstMap(data []map[string]interface{}) map[string]interface{} {
	if len(data) > 0 && data[0] != nil {
		return data[0]
	}
	return nil
}

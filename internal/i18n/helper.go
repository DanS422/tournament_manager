package i18nhelper

import (
	"context"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ctxKey string

const localizerKey ctxKey = "localizer"

var bundle *i18n.Bundle

// LocalizerMiddleware injects a Localizer into request context
func I18nMiddleware(bundle *i18n.Bundle, defaultLang string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("lang")
			if lang == "" {
				lang = r.Header.Get("Accept-Language")
			}
			if lang == "" {
				lang = defaultLang
			}

			loc := i18n.NewLocalizer(bundle, lang)
			ctx := setLocalizer(r.Context(), loc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func setLocalizer(ctx context.Context, l *i18n.Localizer) context.Context {
	return context.WithValue(ctx, localizerKey, l)
}

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	if loc, ok := ctx.Value(localizerKey).(*i18n.Localizer); ok {
		return loc
	}

	return nil
}

package i18n

import "strings"

var defaultLang string

type Messages struct {
	Unauthorized string
	InvalidID    string
	InvalidName  string
}

var translations = map[string]Messages{
	"en": {
		Unauthorized: "unauthorized",
		InvalidID:    "invalid id",
		InvalidName:  "invalid name",
	},
	"ru": {
		Unauthorized: "неавторизован",
		InvalidID:    "неверный id",
		InvalidName:  "неверное имя",
	},
	"kk": {
		Unauthorized: "авторизацияланбаған",
		InvalidID:    "жарамсыз id",
		InvalidName:  "жарамсыз аты",
	},
}

func Init(lang string) {
	lang = strings.ToLower(lang)
	if _, ok := translations[lang]; !ok {
		lang = "en"
	}
	defaultLang = lang
}

func Get(lang string) Messages {
	lang = strings.ToLower(lang)
	if lang != "" {
		if msg, ok := translations[lang]; ok {
			return msg
		}
	}
	return translations[defaultLang]
}

func ParseAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return ""
	}
	parts := strings.Split(acceptLang, ",")
	if len(parts) > 0 {
		lang := strings.TrimSpace(parts[0])
		if idx := strings.Index(lang, ";"); idx != -1 {
			lang = lang[:idx]
		}
		if idx := strings.Index(lang, "-"); idx != -1 {
			lang = lang[:idx]
		}
		return strings.ToLower(lang)
	}
	return ""
}

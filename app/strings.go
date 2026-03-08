package app

import (
	"net/http"
	"strings"
)

// STRINGS holds all localized text in the application.
var STRINGS = map[string]Translation{
	"Title": {
		English: "Vaudreuil-Soulanges News",
		French:  "Nouvelles de Vaudreuil-Soulanges",
	},
	"Subtitle": {
		English: "Local news shouldn't be hard to find.",
		French:  "Les nouvelles locales ne devraient pas être difficiles à trouver.",
	},
	"ThemeLight": {
		English: "Light",
		French:  "Clair",
	},
	"ThemeDark": {
		English: "Dark",
		French:  "Sombre",
	},
	"ThemeAuto": {
		English: "Auto",
		French:  "Auto",
	},
	"LanguageEnglish": {
		English: "English",
		French:  "English",
	},
	"LanguageFrench": {
		English: "Français",
		French:  "Français",
	},
	"LanguageBilingual": {
		English: "Les Deux",
		French:  "Les Deux",
	},
	"ChooseSources": {
		English: "Choose your sources",
		French:  "Choisissez vos sources",
	},
	"ChoicesSaved": {
		English: "Your choices will be saved for your next visit.",
		French:  "Vos choix seront enregistrés pour votre prochaine visite.",
	},
}

// Translation contains equivalent text in English and French.
type Translation struct {
	English string
	French  string
}

// translations is a helper that gets the translations for the given language
func translations(r *http.Request) map[string]string {
	inputLang := "fr"
	if c, err := r.Cookie("lang"); err == nil {
		if c.Value == "en" {
			inputLang = c.Value
		}
	} else {
		accept := r.Header.Get("Accept-Language")
		if strings.Contains(accept, "en") {
			inputLang = "en"
		}
	}

	translatedStrings := map[string]string{}

	for name, translation := range STRINGS {
		if inputLang == "en" {
			translatedStrings[name] = translation.English
		} else {
			translatedStrings[name] = translation.French
		}
	}

	return translatedStrings
}

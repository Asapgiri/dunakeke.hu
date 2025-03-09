package apis

import (
	"dunakeke/dictionary"
	"net/http"
)

func SetLanguage(w http.ResponseWriter, r *http.Request) {
    lang := r.PathValue("lang")
    dictionary.SetLanguage(w, lang)
    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

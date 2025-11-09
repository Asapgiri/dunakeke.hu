package pages

import (
	"dunakeke/config"
	"dunakeke/logic"
	"net/http"

	"github.com/asapgiri/golib/renderer"
)

func User(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    userName := r.PathValue("name")
    user := logic.User{}
    user.FindByUsername(userName)

    session.UpdateTitle(config.Config.Site, user.Username)
    fil, _ := renderer.ReadArtifact("user/page.html", w.Header())
    renderer.Render(session, w, fil, user)
}

func UserEdit(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    userId := r.PathValue("id")
    user := logic.User{}
    user.Find(userId)

    if (!session.Auth.IsMod && session.Auth.Username != user.Username) {
        NotFound(w, r)
        return
    }

    session.UpdateTitle(config.Config.Site, user.Username)
    fil, _ := renderer.ReadArtifact("user/Edit.html", w.Header())
    renderer.Render(session, w, fil, user)
}

func UserEditSave(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    userId := r.PathValue("id")
    user := logic.User{}
    user.Find(userId)

    if (!session.Auth.IsMod && session.Auth.Username != user.Username) {
        NotFound(w, r)
        return
    }

    user.Name           = r.FormValue("form[name]")
    user.Email          = r.FormValue("form[email]")
    user.Phone          = r.FormValue("form[phone]")
    user.EmailVisible   = "on" == r.FormValue("form[email_visible]")
    user.PhoneVisible   = "on" == r.FormValue("form[phone_visible]")

    user.Update()

    http.Redirect(w, r, "/user/"+user.Username, http.StatusSeeOther)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
    NotFound(w, r)
}

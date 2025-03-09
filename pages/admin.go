package pages

import (
	"dunakeke/logic"
	"dunakeke/session"
	"net/http"
	"slices"
)

func checkAdminPageAccess(session session.Sessioner) bool {
    return slices.Contains(session.Auth.Roles, logic.ROLES.ADMIN) ||
           slices.Contains(session.Auth.Roles, logic.ROLES.EDITOR) ||
           slices.Contains(session.Auth.Roles, logic.ROLES.MODERATOR)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    fil, _ := read_artifact("admin/index.html", w.Header())
    Render(session, w, fil, nil)
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    user := logic.User{}
    users := user.List()

    fil, _ := read_artifact("admin/users.html", w.Header())
    Render(session, w, fil, users)
}

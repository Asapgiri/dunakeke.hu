package pages

import (
	"dunakeke/logic"
	"dunakeke/session"
	"net/http"
	"slices"
)

func PostShow(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    //post := logic.Post{}
    //post.Select()...

    fil, _ := read_artifact("post/show.html", w.Header())
    Render(session, w, fil, nil)
}

func checkEditorAccess(session session.Sessioner) bool {
    return slices.Contains(session.Auth.Roles, logic.ROLES.ADMIN) ||
           slices.Contains(session.Auth.Roles, logic.ROLES.EDITOR)
}

func PostEdit(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    fil, _ := read_artifact("post/edit.html", w.Header())
    Render(session, w, fil, nil)
}

func PostDelete(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    fil, _ := read_artifact("post/delete.html", w.Header())
    Render(session, w, fil, nil)
}

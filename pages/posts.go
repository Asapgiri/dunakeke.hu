package pages

import (
	"dunakeke/session"
	"net/http"
)

func PostShow(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    //post := logic.Post{}
    //post.Select()...

    fil, _ := read_artifact("post/show.html", w.Header())
    Render(session, w, fil, nil)
}

func PostEdit(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("post/edit.html", w.Header())
    Render(session, w, fil, nil)
}

func PostDelete(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("post/delete.html", w.Header())
    Render(session, w, fil, nil)
}

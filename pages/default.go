package pages

import (
	"dunakeke/logger"
	"dunakeke/logic"
	"dunakeke/session"
	"io"
	"net/http"
)

var log = logger.Logger {
    Color: logger.Colors.Red,
    Pretext: "pages",
}

func Unexpected(session session.Sessioner, w http.ResponseWriter, r *http.Request) {
    alternative := logic.LinkFind(r.URL.Path)

    if "" !=  alternative {
        http.Redirect(w, r, alternative, http.StatusSeeOther)
    }

    fil, typ := read_artifact(r.URL.Path, w.Header())
    if "" == fil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    if "text" == typ {
        log.Println(r.URL.Path)
        Render(session, w, fil, nil)
    } else {
        // TODO: Check if file type/path needs auth..
        // If it is in artifacts tho is shouldn't..
        io.WriteString(w, fil)
    }
}

func Root(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if "/" == r.URL.Path {
        dto := DtoRoot{
            Main: DtoMain{},
            Posts: logic.PostList(),
        }

        fil, _ := read_artifact("index.html", w.Header())
        Render(session, w, fil, dto)
    } else {
        Unexpected(session, w, r)
    }
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("not-found.html", w.Header())
    Render(session, w, fil, nil)
}

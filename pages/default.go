package pages

import (
	"asapgiri/golib/logger"
	"asapgiri/golib/renderer"
	"asapgiri/golib/session"
	"dunakeke/logic"
	"io"
	"net/http"
	"sort"
)

var log = logger.Logger {
    Color: logger.Colors.Red,
    Pretext: "pages",
}

func Unexpected(session session.Sessioner, w http.ResponseWriter, r *http.Request) {
    link := logic.Link{}
    link.SelectByAlternative(r.URL.Path)

    if "" !=  link.Original {
        r.URL.Path = link.Original
        http.DefaultServeMux.ServeHTTP(w, r)
        return
    }

    fil, typ := renderer.ReadArtifact(r.URL.Path, w.Header())
    if "" == fil {
        // FIXME: Redirect due to request type...
        //http.Error(w, "File not found", http.StatusNotFound)

        NotFound(w, r)
        return
    }

    if "text" == typ {
        log.Println(r.URL.Path)
        renderer.Render(session, w, fil, nil)
    } else {
        // TODO: Check if file type/path needs auth..
        // If it is in artifacts tho is shouldn't..
        io.WriteString(w, fil)
    }
}

func Root(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    post := logic.Post{}
    plist := post.List(checkEditorAccess(session))

    // FIXME: Check if post is public or not..
    sort.Slice(plist, func(i, j int) bool { return plist[i].EditDate.After(plist[j].EditDate) })

    if "/" == r.URL.Path {
        dto := DtoRoot{
            Main: DtoMain{},
            Posts: plist,
        }

        fil, _ := renderer.ReadArtifact("index.html", w.Header())
        renderer.Render(session, w, fil, dto)
    } else {
        Unexpected(session, w, r)
    }
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    fil, _ := renderer.ReadArtifact("not-found.html", w.Header())
    renderer.Render(session, w, fil, nil)
}

func AccessViolation(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    fil, _ := renderer.ReadArtifact("auth/access-violation.html", w.Header())
    renderer.Render(session, w, fil, nil)
}

func renderPageWithAccessViolation(w http.ResponseWriter, r *http.Request) {
    AccessViolation(w, r)
}

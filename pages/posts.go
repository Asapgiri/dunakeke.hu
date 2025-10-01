package pages

import (
	"dunakeke/config"
	"dunakeke/logic"
	"dunakeke/session"
	"encoding/json"
	"io"
	"net/http"
	"slices"
)

func checkEditorAccess(session session.Sessioner) bool {
    return slices.Contains(session.Auth.Roles, logic.ROLES.ADMIN) ||
           slices.Contains(session.Auth.Roles, logic.ROLES.EDITOR)
}

func PostShow(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    post := logic.Post{}
    err := post.Select(r.PathValue("id"))
    if nil != err || (!post.Public && !checkEditorAccess(session)) {
        NotFound(w, r)
        return
    }

    fil, _ := read_artifact("post/show.html", w.Header())
    Render(session, w, fil, post)
}

func PostNew(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    user := logic.User{}
    user.FindByUsername(session.Auth.Username)

    id := logic.PostNew(session.Dictionary, user)

    http.Redirect(w, r, "/post/edit/" + id, http.StatusSeeOther)
}

func PostEdit(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    post := logic.Post{}
    err := post.Select(r.PathValue("id"))
    if nil != err {
        NotFound(w, r)
        return
    }

    fil, _ := read_artifact("post/edit.html", w.Header())
    Render(session, w, fil, post)
}

func PostDelete(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    post := logic.Post{}
    err := post.Select(r.PathValue("id"))
    if nil != err {
        NotFound(w, r)
        return
    }

    err = post.Delete()
    if nil != err {
        // FIXME: better..
        NotFound(w, r)
        return
    }

    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func PostSave(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    ps := logic.PostSave{}
    de := json.NewDecoder(r.Body)
    de.DisallowUnknownFields()

    err := de.Decode(&ps)
    if nil != err {
        io.WriteString(w, "NOK")
        return
    }

    err = logic.PostUpdate(ps)
    if nil != err {
        io.WriteString(w, "NOK")
        return
    }

    io.WriteString(w, "OK")
}

func PostSaveImage(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    psir := PostSaveImageResponse{Success: 0}
    err_ret, _ := json.Marshal(psir)

    err := r.ParseMultipartForm(config.Config.Site.MaxImgUploadMB << 20)
    if nil != err {
        io.WriteString(w, string(err_ret))
        return
    }
    
    file, header, err := r.FormFile("editormd-image-file")
    if nil != err {
        io.WriteString(w, string(err_ret))
        return
    }
    defer file.Close()

    // FIXME: Save which photos are being used...
    save_name := "/photos/" + header.Filename
    err = save_artifact(save_name, file)
    if nil != err {
        io.WriteString(w, string(err_ret))
        return
    }

    psir.Url = save_name
    psir.Success = 1
    success_ret, _ := json.Marshal(psir)
    io.WriteString(w, string(success_ret))
}

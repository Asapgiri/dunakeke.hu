package pages

import (
	"asapgiri/golib/renderer"
	"asapgiri/golib/session"
	"dunakeke/config"
	"dunakeke/dictionary"
	"dunakeke/logic"
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
    session := GetCurrentSession(r)

    post := logic.Post{}
    err := post.Select(r.PathValue("id"))
    if nil != err || (!post.Public && !checkEditorAccess(session)) {
        NotFound(w, r)
        return
    }

    fil, _ := renderer.ReadArtifact("post/show.html", w.Header())
    renderer.Render(session, w, fil, post)
}

func PostNew(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    user := logic.User{}
    user.FindByUsername(session.Auth.Username)

    id := logic.PostNew(session.Dictionary.(dictionary.Dictionary), user)

    http.Redirect(w, r, "/post/edit/" + id, http.StatusSeeOther)
}

func PostEdit(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

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

    fil, _ := renderer.ReadArtifact("post/edit.html", w.Header())
    renderer.Render(session, w, fil, post)
}

func PostDelete(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

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
    session := GetCurrentSession(r)

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

func PostPublish(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

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

    post.Public = "public" == r.PathValue("val")
    err = post.Update()

    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func PostSaveImage(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

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
    err = renderer.SaveArtifact(save_name, file)
    if nil != err {
        io.WriteString(w, string(err_ret))
        return
    }

    psir.Url = save_name
    psir.Success = 1
    success_ret, _ := json.Marshal(psir)
    io.WriteString(w, string(success_ret))
}

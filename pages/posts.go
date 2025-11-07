package pages

import (
	"dunakeke/config"
	"dunakeke/dictionary"
	"dunakeke/logic"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/asapgiri/golib/renderer"
	"github.com/asapgiri/golib/session"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

    session.UpdateTitle(config.Config.Site, post.Title)
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

func postEdit(session session.Sessioner, w http.ResponseWriter, r *http.Request) {
    dp := DtoEditor{}
    err := dp.Post.Select(r.PathValue("id"))
    if nil != err {
        NotFound(w, r)
        return
    }

    tag := logic.Tag{}
    tags, err := tag.List()

    dp.Tags = make([]DtoTag, len(tags))
    for i, t := range(tags) {
        dp.Tags[i].Tag = t
        dp.Tags[i].Selected = slices.ContainsFunc(dp.Post.Tags, func(tf logic.Tag) bool {
            return t.Id == tf.Id
        })
    }

    session.UpdateTitle(config.Config.Site, dp.Post.Title)
    fil, _ := renderer.ReadArtifact("post/edit.html", w.Header())
    renderer.Render(session, w, fil, dp)
}

func PostEdit(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    postEdit(session, w, r)
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
        io.WriteString(w, "NOK - Decode")
        return
    }

    user := logic.User{}
    user.FindByUsername(session.Auth.Username)

    err = logic.PostUpdate(ps, user)
    if nil != err {
        io.WriteString(w, "NOK - Update")
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

func saveFileFromForm(dict dictionary.Dictionary, fileTag string, r *http.Request) (logic.File, error) {
    err := r.ParseMultipartForm(config.Config.Site.MaxImgUploadMB << 20)
    if nil != err {
        log.Println(err)
        return logic.File{}, errors.New(dict.Editor.ErrorFileIsLargerThan+strconv.FormatInt(config.Config.Site.MaxImgUploadMB, 10)+"MB")
    }

    file, header, err := r.FormFile(fileTag)
    if nil != err {
        log.Println(err)
        return logic.File{}, errors.New(dict.Editor.ErrorFromFile)
    }
    defer file.Close()

    new_id := primitive.NewObjectID().Hex()
    sfile := logic.File{
        Id:         new_id,
        Name:       header.Filename,
        SaveName:   "/files/" + new_id + "-" + header.Filename,
    }
    parts := strings.Split(sfile.Name, ".")
    sfile.Extension = strings.ToLower(parts[len(parts)-1])

    err = renderer.SaveArtifact(sfile.SaveName, file)
    if nil != err {
        log.Println(err)
        return logic.File{}, errors.New(dict.Editor.ErrorFromFile)
    }

    sfile.Add()

    return sfile, nil
}

func PostSaveImage(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    psir := PostSaveImageResponse{Success: 0}
    err_ret, _ := json.Marshal(psir)

    file, err := saveFileFromForm(session.Dictionary.(dictionary.Dictionary), "editormd-image-file", r)
    if nil != err {
        io.WriteString(w, string(err_ret))
        return
    }

    psir.Url = file.SaveName
    psir.Alt = file.Name
    psir.Success = 1
    success_ret, _ := json.Marshal(psir)
    io.WriteString(w, string(success_ret))
}

func PostEditPhotoSave(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    file, err := saveFileFromForm(session.Dictionary.(dictionary.Dictionary), "image-input", r)
    if nil != err {
        session.Error = err.Error()
        postEdit(session, w, r)
        return
    }

    post := logic.Post{}
    err = post.Select(r.PathValue("id"))
    if nil != err {
        session.Error = err.Error()
        postEdit(session, w, r)
        return
    }
    post.Image = file.SaveName
    post.Update()

    postEdit(session, w, r)
}

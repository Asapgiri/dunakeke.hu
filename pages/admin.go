package pages

import (
	"asapgiri/golib/renderer"
	"asapgiri/golib/session"
	"dunakeke/logic"
	"encoding/json"
	"io"
	"net/http"
	"slices"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func checkAdminPageAccess(session session.Sessioner) bool {
    return slices.Contains(session.Auth.Roles, logic.ROLES.ADMIN)  ||
           slices.Contains(session.Auth.Roles, logic.ROLES.EDITOR) ||
           slices.Contains(session.Auth.Roles, logic.ROLES.MODERATOR)
}

func adminRender(session session.Sessioner, w http.ResponseWriter, temp string, dto any) {
    renderer.RenderMultiTemplate(session, w, []string{temp, "admin/base.html"}, dto)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    adminRender(session, w, "admin/index.html", nil)
}

func adminRenderUsers(session session.Sessioner, w http.ResponseWriter) {
    user := logic.User{}

    dto := DtoAdminUsers{
        Users: user.List(),
        Roles: logic.RolePerms,
    }

    adminRender(session, w, "admin/users.html", dto)
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    adminRenderUsers(session, w)
}

func AdminPosts(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    post := logic.Post{}
    posts := post.List(true, nil, true)

    adminRender(session, w, "admin/posts.html", posts)
}

func AdminTags(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    tag := logic.Tag{}
    tags, _ := tag.List()

    adminRender(session, w, "admin/tags.html", tags)
}

func AdminDonations(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    user := logic.User{}
    don := logic.Donation{}
    ad := DtoAdminDonations{
        Donations: don.List(),
        Sum: 0.0,
    }
    ad.Username = make([]string, len(ad.Donations))
    for i, d := range(ad.Donations) {
        ad.Sum += d.Amount * float64(len(d.Occurences))

        // check invalid "000000000000000000000000" id
        _, err := primitive.ObjectIDFromHex(d.UserId)
        if  nil == err {
            user.Find(d.UserId)
            ad.Username[i] = user.Username
        }
    }

    adminRender(session, w, "admin/donations.html", ad)
}

func AdminUserSetRole(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    id := r.PathValue("id")
    trole := r.PathValue("role")

    eperm := logic.FindPermsFor(trole)

    log.Println(id)
    log.Println(trole)
    log.Println(eperm)

    if !renderer.Subset(session.Auth.Roles, eperm.EditPerm) {
        session.Error = "You cannot edit user!"
        // FIXME: should be a better way to redirect and show errors

        log.Println(session.Error)
        adminRenderUsers(session, w)
        return
    }

    user := logic.User{}
    user.Find(id)
    if "" == user.Username {
        session.Error = "User not found!"
        log.Println(session.Error)
        adminRenderUsers(session, w)
        return
    }

    log.Println(user)



    if slices.Contains(user.Roles, trole) {
        i := slices.Index(user.Roles, trole)
        user.Roles = append(user.Roles[:i], user.Roles[i+1:]...)
    } else {
        user.Roles = append(user.Roles, trole)
    }

    log.Println(user)
    user.Update()

    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func AdminLinks(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    link := logic.Link{}
    links := link.List()

    adminRender(session, w, "admin/links.html", links)
}

func AdminLinksUpdate(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    lu := LinkUpdate{}
    de := json.NewDecoder(r.Body)
    de.DisallowUnknownFields()
    err := de.Decode(&lu)
    if nil != err {
        return
    }

    user := logic.User{}
    user.FindByUsername(session.Auth.Username)

    logic.AlternativeUpdate(lu.Original, lu.Alternative, user)

    io.WriteString(w, "OK")
}

func AdminLinksDelete(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    id := r.PathValue("id")
    link := logic.Link{}
    link.Select(id)
    link.Delete()

    http.Redirect(w, r, "/admin/links", http.StatusSeeOther)
}

func AdminTagsUpdate(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    lu := TagUpdate{}
    de := json.NewDecoder(r.Body)
    de.DisallowUnknownFields()
    err := de.Decode(&lu)
    if nil != err {
        return
    }

    log.Println(lu)

    tag := logic.Tag{}
    err = tag.Select(lu.Id)

    tag.Name = lu.Name
    tag.Color = lu.Color

    if nil != err {
        tag.Listable = true
        tag.Add()
    } else {
        tag.Update()
    }

    log.Println(tag)

    io.WriteString(w, "OK")
}

func AdminTagsToggleListable(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    tagId := r.PathValue("id")
    tag := logic.Tag{}
    tag.Select(tagId)

    tag.Listable = !tag.Listable
    tag.Update()

    http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func AdminTagsDelete(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    id := r.PathValue("id")
    tag := logic.Tag{}
    tag.Select(id)
    tag.Delete()

    http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

package pages

import (
    "asapgiri/golib/renderer"
    "asapgiri/golib/session"
	"dunakeke/logic"
	"net/http"
	"slices"
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

func AdminUsers(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    user := logic.User{}
    users := user.List()

    adminRender(session, w, "admin/users.html", users)
}

func AdminPosts(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    post := logic.Post{}
    posts := post.List(true)

    adminRender(session, w, "admin/posts.html", posts)
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
        if "000000000000000000000000" != d.UserId {
            user.Find(d.UserId)
            ad.Username[i] = user.Username
        }
    }

    adminRender(session, w, "admin/donations.html", ad)
}

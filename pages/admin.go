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

func AdminPage(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    fil, _ := renderer.ReadArtifact("admin/index.html", w.Header())
    renderer.Render(session, w, fil, nil)
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    user := logic.User{}
    users := user.List()

    fil, _ := renderer.ReadArtifact("admin/users.html", w.Header())
    renderer.Render(session, w, fil, users)
}

func AdminPosts(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkAdminPageAccess(session) {
        NotFound(w, r)
        return
    }

    post := logic.Post{}
    posts := post.List(true)

    fil, _ := renderer.ReadArtifact("admin/posts.html", w.Header())
    renderer.Render(session, w, fil, posts)
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

    fil, _ := renderer.ReadArtifact("admin/donations.html", w.Header())
    renderer.Render(session, w, fil, ad)
}

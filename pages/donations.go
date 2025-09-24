package pages

import (
	"dunakeke/session"
	"dunakeke/logic"
	"net/http"
)

func DonationRoot(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    do := logic.DonationOption{}
    dos := do.List()

    fil, _ := read_artifact("donate/root.html", w.Header())
    Render(session, w, fil, dos)
}

func DonationInProgress(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("donate/in-progress.html", w.Header())
    Render(session, w, fil, nil)
}

func DonationSuccess(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)


    // if "" != session.Auth.Username {
    //     ddon.User, _ = primitive.ObjectIDFromHex(session.Auth.Id)
    // }


    fil, _ := read_artifact("donate/success.html", w.Header())
    Render(session, w, fil, nil)
}

func DonationError(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("donate/error.html", w.Header())
    Render(session, w, fil, nil)
}

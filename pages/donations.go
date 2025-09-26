package pages

import (
	"dunakeke/logic"
	"dunakeke/session"
	"net/http"
	"strconv"
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

    // FIXME: Sanitize further::
    f_amount := r.FormValue("form[amount]")
    // FIXME: Use form values..
    f_name := r.FormValue("form[name]")
    f_email := r.FormValue("form[email]")
    f_subscribeToNewsletter := r.FormValue("form[subscribeToNewsletter]")
    f_gdprAgreed := r.FormValue("form[gdprAgreed]")
    f_csrf := r.FormValue("form[csrf]")

    log.Println(f_amount)
    log.Println(f_name)
    log.Println(f_email)
    log.Println(f_subscribeToNewsletter)
    log.Println(f_gdprAgreed)
    log.Println(f_csrf)
    amount, err := strconv.ParseFloat(f_amount, 64)

    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := read_artifact("donate/error.html", w.Header())
        Render(session, w, fil, err)
        return
    }

    otp_ret, err := logic.RedirectToOtpApi(session.Dictionary, amount, f_email)

    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := read_artifact("donate/error.html", w.Header())
        Render(session, w, fil, err)
    } else {
        log.Printf("Redirect URL: %s\n", otp_ret.PaymentUrl)
        http.Redirect(w, r, otp_ret.PaymentUrl, http.StatusSeeOther)
    }
}

func DonationReturn(w http.ResponseWriter, r *http.Request) {
    if logic.ProgressOtpReply(r.URL.Query().Get("r"), r.URL.Query().Get("s")) {
        http.Redirect(w, r, "/donate/success", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/donate/failure", http.StatusSeeOther)
    }
}

func DonationSuccess(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("donate/success.html", w.Header())
    Render(session, w, fil, nil)
}

func DonationFailure(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("donate/fail.html", w.Header())
    Render(session, w, fil, nil)
}

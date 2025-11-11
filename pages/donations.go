package pages

import (
	"dunakeke/config"
	"dunakeke/dictionary"
	"dunakeke/logic"
	"net/http"
	"strconv"
	"time"

	"github.com/asapgiri/golib/renderer"
	"github.com/asapgiri/golib/session"
)

func DonationRoot(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    do := logic.DonationOption{}
    dos := do.List()

    session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.Header)
    fil, _ := renderer.ReadArtifact("donate/root.html", w.Header())
    renderer.Render(session, w, fil, dos)
}

func DonationInProgress(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if "1" != r.FormValue("form[gdprAgreed]") {
        log.Printf("GDPR not accepted!\n")
        fil, _ := renderer.ReadArtifact("donate/error.html", w.Header())
        renderer.Render(session, w, fil, nil)
        return
    }

    amount, err := strconv.ParseFloat(r.FormValue("form[amount]"), 64)
    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := renderer.ReadArtifact("donate/error.html", w.Header())
        renderer.Render(session, w, fil, err)
        return
    }

    // FIXME: Sanitize further::
    // Also check for mandatory fields
    donation := logic.Donation{
        UserId: session.Auth.Id,
        Date: time.Now(),
        Name: r.FormValue("form[name]"),
        Message: r.FormValue("form[message]"),
        Email: r.FormValue("form[email]"),
        Amount: amount,
        Newsletter: "1" == r.FormValue("form[subscribeToNewsletter]"),
        GDPR: "1" == r.FormValue("form[gdprAgreed]"),
        Recurring: "1" == r.FormValue("form[recurring]"),
    }

    log.Println(donation)

    // FIXME: undo after testing..
    otp_ret, err := logic.RedirectToOtpApi(session.Dictionary.(dictionary.Dictionary), donation)
    //otp_ret := logic.OtpJsonResponse{PaymentUrl: "/donate"}

    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := renderer.ReadArtifact("donate/error.html", w.Header())
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.Header)
        renderer.Render(session, w, fil, err)
    } else {
        log.Printf("Redirect URL: %s\n", otp_ret.PaymentUrl)
        http.Redirect(w, r, otp_ret.PaymentUrl, http.StatusSeeOther)

        // FIXME: Implementation should occasionally check if the request finished,
        //        if the user closed the SimplePay site..
    }
}

func DonationReturn(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    donation, err := logic.ProgressOtpReply(r.URL.Query().Get("r"), r.URL.Query().Get("s"))
    if donation.Successful && nil == err {
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionSuccess)
        http.Redirect(w, r, "/donate/" + donation.Id, http.StatusSeeOther)
    } else {
        // TODO: Handle errors, and their passing ...
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionFailed)
        http.Redirect(w, r, "/donate/" + donation.Id, http.StatusSeeOther)
    }

    donationEmail(session, donation)
}

func DonationShowStatus(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    donation := logic.Donation{Id: r.PathValue("id")}
    donation.Select()

    fil, _ := renderer.ReadArtifact("donate/success.html", w.Header())
    if donation.Successful {
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionSuccess)
    } else {
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionFailed)
    }
    renderer.Render(session, w, fil, donation)
}

func donationEmail(session session.Sessioner, donation logic.Donation) {
    fil, _ := renderer.ReadArtifact("donate/email.html", nil)
    session.MainDto = config.Config
    session.Dto = donation
    message := renderer.PreRender(fil, session)

    logic.SendEmail(session.Config.Title, message, logic.Messagee{Name: donation.Name, Email: donation.Email})
}

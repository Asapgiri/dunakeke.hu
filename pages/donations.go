package pages

import (
	"dunakeke/logic"
	"dunakeke/session"
	"net/http"
	"strconv"
	"time"
)

func DonationRoot(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    do := logic.DonationOption{}
    dos := do.List()

    fil, _ := read_artifact("donate/root.html", w.Header())
    Render(session, w, fil, dos)
}

func checkCSFR(csfr string) bool {
    return true
}

func DonationInProgress(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    if !checkCSFR(r.FormValue("form[csrf]")) {
        log.Printf("CFSR ERR!\n")
        fil, _ := read_artifact("donate/error.html", w.Header())
        Render(session, w, fil, nil)
        return
    }

    if "1" != r.FormValue("form[gdprAgreed]") {
        log.Printf("GDPR not accepted!\n")
        fil, _ := read_artifact("donate/error.html", w.Header())
        Render(session, w, fil, nil)
        return
    }

    amount, err := strconv.ParseFloat(r.FormValue("form[amount]"), 64)
    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := read_artifact("donate/error.html", w.Header())
        Render(session, w, fil, err)
        return
    }

    // FIXME: Sanitize further::
    // Also check for mandatory fields
    donation := logic.Donation{
        UserId: session.Auth.Id,
        Date: time.Now(),
        Name: r.FormValue("form[name]"),
        Email: r.FormValue("form[email]"),
        Amount: amount,
        Newsletter: "1" == r.FormValue("form[subscribeToNewsletter]"),
        GDPR: "1" == r.FormValue("form[gdprAgreed]"),
        InvoiceNeeded: "1" == r.FormValue("form[invoiceneeded]"),
        Recurring: "1" == r.FormValue("form[recurring]"),
    }

    if donation.InvoiceNeeded {
        donation.Invoice = logic.Invoice{
            Name:       donation.Name,
            Company:    r.FormValue("form[invoicecompany]"),
            Country:    r.FormValue("form[invoicecountry]"),
            State:      r.FormValue("form[invoicestate]"),
            City:       r.FormValue("form[invoicecity]"),
            Zip:        r.FormValue("form[invoicezip]"),
            Address:    r.FormValue("form[invoiceaddress]"),
            Address2:   r.FormValue("form[invoiceaddress2]"),
            Phone:      r.FormValue("form[invoicephone]"),
            TaxNumber:  r.FormValue("form[invoictaxnumber]"),
        }
    }

    log.Println(donation)

    // FIXME: undo after testing..
    otp_ret, err := logic.RedirectToOtpApi(session.Dictionary, donation)
    //otp_ret := logic.OtpJsonResponse{PaymentUrl: "/donate"}

    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := read_artifact("donate/error.html", w.Header())
        session.UpdateTitle(session.Dictionary.Donate.Header)
        Render(session, w, fil, err)
    } else {
        log.Printf("Redirect URL: %s\n", otp_ret.PaymentUrl)
        http.Redirect(w, r, otp_ret.PaymentUrl, http.StatusSeeOther)
    }
}

func DonationReturn(w http.ResponseWriter, r *http.Request) {
    id, success, err := logic.ProgressOtpReply(r.URL.Query().Get("r"), r.URL.Query().Get("s"))
    if success && nil == err {
        http.Redirect(w, r, "/donate/" + id, http.StatusSeeOther)
    } else {
        // TODO: Handle errors, and their passing ...
        http.Redirect(w, r, "/donate/" + id, http.StatusSeeOther)
    }
}

func DonationShowStatus(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    donation := logic.Donation{Id: r.PathValue("id")}
    donation.Select()

    fil, _ := read_artifact("donate/success.html", w.Header())
    if donation.Successful {
        session.UpdateTitle(session.Dictionary.Donate.TransactionSuccess)
    } else {
        session.UpdateTitle(session.Dictionary.Donate.TransactionFailed)
    }
    Render(session, w, fil, donation)
}

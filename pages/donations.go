package pages

import (
	"dunakeke/config"
	"dunakeke/dictionary"
	"dunakeke/logic"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asapgiri/golib/renderer"
	"github.com/asapgiri/golib/session"
)

var trustedNets []*net.IPNet

func DonationInit() {
    for _, cidr := range config.Config.Donation.SimplePayTrustedIPs {
        _, network, err := net.ParseCIDR(cidr)
        if err != nil {
            log.Printf("Invalid CIDR %s: %v\n", cidr, err)
        }
        trustedNets = append(trustedNets, network)
    }
    log.Println(trustedNets)
}

func DonationRoot(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(w, r)

    do := logic.DonationOption{}
    dos := do.List()

    session.Meta["title"] = session.Dictionary.(dictionary.Dictionary).Donate.Header
    session.Meta["description"] = session.Dictionary.(dictionary.Dictionary).Donate.Description
    session.Meta["twitter:card"] = "summary_large_image"
    session.Meta["theme-color"] = "#0d6efd"

    session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.Header)
    fil, _ := renderer.ReadArtifact("donate/root.html", w.Header())
    renderer.Render(session, w, fil, dos)
}

func DonationInProgress(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(w, r)

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
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.Header)
        renderer.Render(session, w, fil, err)
        return
    }

    if amount < config.Config.Donation.MinAmount {
        session.Notice.Set(NOTICE.DANGER, fmt.Sprint("Amount ", amount, " is smaller than minimum amount ", config.Config.Donation.MinAmount))
        http.Redirect(w, r, "/donate", http.StatusSeeOther)
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
    otp_ret, err := logic.RedirectToOtpApi(session.Dictionary.(dictionary.Dictionary), &donation)
    //otp_ret := logic.OtpJsonResponse{PaymentUrl: "/donate"}

    if nil != err {
        log.Printf("Redirect ERR: %s\n", err)
        fil, _ := renderer.ReadArtifact("donate/root.html", w.Header())
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.Header)
        renderer.Render(session, w, fil, err)
    } else {
        log.Printf("Redirect URL: %s\n", otp_ret.PaymentUrl)
        http.Redirect(w, r, otp_ret.PaymentUrl, http.StatusSeeOther)

        // go logic.CheckTransactionProgress(donation, func(d logic.Donation) {
        //     donationEmail(session, d)
        // })
    }
}

func DonationReturn(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(w, r)

    donation, err := logic.ProgressOtpReply(r.URL.Query().Get("r"), r.URL.Query().Get("s"))
    if donation.Successful && nil == err {
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionSuccess)
    } else {
        // TODO: Handle errors, and their passing ...
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionFailed)
    }
    http.Redirect(w, r, "/donate/" + donation.Id, http.StatusSeeOther)
}

func checkSimpepayIPRanges(r *http.Request) bool {
    ip := r.RemoteAddr
    xf := r.Header.Get("X-Forwarded-For")
    if "" != xf {
        parts := strings.Split(xf, ",")
        ip = strings.TrimSpace(parts[0])
    } else {
        ip, _, _ = net.SplitHostPort(ip)
    }

    parsedip := net.ParseIP(ip)
    if nil == parsedip {
        return false
    }

    for _, network := range trustedNets {
        if network.Contains(parsedip) {
            return true
        }
    }

    return false
}

func DonationIpn(w http.ResponseWriter, r *http.Request) {
    if !checkSimpepayIPRanges(r) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    resp, err := io.ReadAll(r.Body)
    if nil != err {
        http.Error(w, "Cannot read message body", http.StatusBadRequest)
        return
    }

    response, err := logic.DonationIpn(w, r.Header.Get("Signature"), resp, func(d logic.Donation) {
        donationEmail(GetCurrentSession(w, r), d)
    })
    if nil != err {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    io.WriteString(w, string(response))
}

func DonationShowStatus(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(w, r)

    donation := logic.Donation{Id: r.PathValue("id")}
    donation.Select()

    fil, _ := renderer.ReadArtifact("donate/success.html", w.Header())
    if donation.Successful {
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionSuccess)
    } else {
        session.UpdateTitle(config.Config.Site, session.Dictionary.(dictionary.Dictionary).Donate.TransactionFailed)
    }

    public := PublicDonation{
        Status: logic.DonationGetPublicStatus(donation, session.Dictionary.(dictionary.Dictionary)),
        Donation: donation,
    }

    renderer.Render(session, w, fil, public)
}

func donationEmail(session session.Sessioner, donation logic.Donation) {
    fil, _ := renderer.ReadArtifact("donate/email.html", nil)

    public := PublicDonation{
        Status: logic.DonationGetPublicStatus(donation, session.Dictionary.(dictionary.Dictionary)),
        Config: config.Config,
        Donation: donation,
    }

    session.Dto = public
    message := renderer.PreRender(fil, session)

    logic.SendEmail(session.Config.Title, message, logic.Messagee{Name: donation.Name, Email: donation.Email})
}

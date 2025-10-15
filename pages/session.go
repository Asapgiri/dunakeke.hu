package pages

import (
	"asapgiri/golib/session"
	"dunakeke/config"
	"dunakeke/dictionary"
	"dunakeke/logic"
	"net/http"
)

func GetCurrentSession(r *http.Request) session.Sessioner {
    session := session.Sessioner{}
    session.Authenticate(r)
    logic.Authenticate(&session.Auth)

    session.Dictionary = dictionary.GetLanguage(r)
    session.Config = config.Config.Site

    // FIXME: Put this to somewhere
    logic.SaveStatistics(r, session.Auth.Id)

    return session
}

package pages

import (
	"github.com/asapgiri/golib/session"
	"dunakeke/config"
	"dunakeke/dictionary"
	"dunakeke/logic"
	"net/http"
)

func GetCurrentSession(w http.ResponseWriter, r *http.Request) session.Sessioner {
    sess := session.Sessioner{}
    sess.Authenticate(w, r)
    logic.Authenticate(&sess.Auth)

    sess.Dictionary = dictionary.GetLanguage(r)
    sess.Config = config.Config.Site
    sess.Path = r.URL.String()

    // FIXME: Put this to somewhere
    logic.SaveStatistics(r, sess.Auth.Id)

    // FIXME: Should be somewhere else...
    supp := logic.Supporter{}
    supporters, _ := supp.List()
    sess.MainDto = Footer{
        Supporters: supporters,
        Sections: []Section{
            Section{
                Title: "Connections",
                Sites: []SiteA{
                    SiteA{
                        Title: "Facebook",
                        Icon: "https://images.seeklogo.com/logo-png/29/2/facebook-icon-logo-png_seeklogo-290338.png",
                        Url: "https://www.facebook.com/DUNAKEKE2021/",
                        Blank: true,
                    },
                    SiteA{Title: "Telefon: +36/30/179-58-86"},
                },
            },
            Section{
                Title: "About",
                Sites: []SiteA{
                    SiteA{Title: "Adatvedelem", Icon: "", Url: "/adatkezelesi-tajekoztato"},
                    SiteA{Title: "Impresszum",  Icon: "", Url: "/impressum"},
                    SiteA{Title: "Beszamolok",  Icon: "", Url: "/beszamolok"},
                },
            },
        },
    }
    sess.Meta = session.MetaData{}

    return sess
}

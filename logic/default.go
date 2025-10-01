package logic

import (
	"dunakeke/dbase"
	"dunakeke/logger"
)

var log = logger.Logger {
    Color: logger.Colors.Cyan,
    Pretext: "logic",
}

func LinkFind(link string) string {
    dblink := dbase.Link{}
    err := dblink.FindByLink(link)

    if nil != err {
        return ""
    }

    return dblink.Alternative
}

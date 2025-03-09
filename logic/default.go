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

func PostList() []Post {
    post := dbase.Post{}
    dbposts, err := post.List()

    if nil != err {
        return []Post{}
    }

    lposts := make([]Post, len(dbposts))

    for i, p := range(dbposts) {
        lposts[i].Map(p)
    }

    return lposts
}

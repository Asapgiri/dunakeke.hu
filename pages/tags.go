package pages

import (
	"asapgiri/golib/renderer"
	"dunakeke/config"
	"dunakeke/logic"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type TagSave struct {
    Name    string `json:"name"`
}

type TagRespons struct {
    Id      string `json:"id"`
    Name    string `json:"name"`
    Color   string `json:"color"`
}

func TagList(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    tagName := r.PathValue("tagname")
    page, err := strconv.ParseInt(r.PathValue("page"), 10, 32)
    if nil != err {
        page = 0
    }
    post_per_page, err := strconv.ParseInt(r.PathValue("ppp"), 10, 32)
    if nil != err {
        post_per_page = 25
    }

    tag := logic.Tag{}
    tag.SelectByName(tagName)

    post := logic.Post{}
    plist, pages := post.List(checkEditorAccess(session), &tag.Id, int(page), int(post_per_page), false)

    dto := DtoRoot{
        Main: DtoMain{Title: "Tag: "+tag.Name},
        Posts: plist,
        Page: Pages{
            Current: int(page),
            Count: pages,
            Ppp: int(post_per_page),
        },
    }

    session.UpdateTitle(config.Config.Site, dto.Main.Title)
    fil, _ := renderer.ReadArtifact("index.html", w.Header())
    renderer.Render(session, w, fil, dto)
}

func TagAdd(w http.ResponseWriter, r *http.Request) {
    session := GetCurrentSession(r)

    if !checkEditorAccess(session) {
       renderPageWithAccessViolation(w, r)
       return
    }

    ps := TagSave{}
    de := json.NewDecoder(r.Body)
    de.DisallowUnknownFields()

    err := de.Decode(&ps)
    if nil != err {
        io.WriteString(w, "NOK - Decode")
        return
    }

    tag := logic.Tag{}
    err = tag.SelectByName(ps.Name)
    if nil == err {
        io.WriteString(w, "NOK - Tagname exists")
        return
    }

    tag.Listable = true
    tag.Name = ps.Name
    err = tag.Add()
    if nil != err {
        io.WriteString(w, "NOK - Failed to add tag")
        return
    }

    resp, err := json.Marshal(TagRespons{
        Id: tag.Id,
        Name: tag.Name,
        Color: tag.Color,
    })

    io.WriteString(w, string(resp))
}

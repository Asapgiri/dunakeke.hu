package pages

import (
	"asapgiri/golib/renderer"
	"dunakeke/logic"
	"encoding/json"
	"io"
	"net/http"
	"sort"
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

    tag := logic.Tag{}
    tag.SelectByName(tagName)

    post := logic.Post{}
    plist := post.List(checkEditorAccess(session), &tag.Id, false)

    // FIXME: Check if post is public or not..
    sort.Slice(plist, func(i, j int) bool { return plist[i].EditDate.After(plist[j].EditDate) })

    dto := DtoRoot{
        Main: DtoMain{Title: "Tag: "+tag.Name},
        Posts: plist,
    }

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

package logic

import (
	"dunakeke/dbase"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AlternativeUpdate(original string, alternative string, user User) error {
    link := Link{}
    link.SelectByOrigin(original)

    link.Alternative = alternative
    link.Date = time.Now()

    // check invalid "000000000000000000000000" id
    _, err := primitive.ObjectIDFromHex(link.Id)
    if nil == err {
        return link.Update()
    } else {
        link.Author = user
        link.Original = original
        return link.Add()
    }
}

func (link *Link) List() []Link {
    dlink := dbase.Link{}
    dlinks, _ := dlink.List()

    links := make([]Link, len(dlinks))
    for i, dl := range(dlinks) {
        links[i].Map(dl)
    }

    return links
}

func (link *Link) Add() error {
    dlink := link.UnMap()
    dlink.Id = primitive.NewObjectID()
    link.Id = dlink.Id.Hex()
    return dlink.Add()
}

func (link *Link) Update() error {
    dlink := link.UnMap()
    return dlink.Update()
}

func (link *Link) Select(id string) error {
    dlink := dbase.Link{}
    oid, _ := primitive.ObjectIDFromHex(id)
    err := dlink.Select(oid)
    if nil != err {
        return err
    }

    link.Map(dlink)
    return nil
}

func (link *Link) SelectByOrigin(alter string) error {
    dlink := dbase.Link{}
    err := dlink.FindByOriginal(alter)
    if nil != err {
        return err
    }
    link.Map(dlink)
    return nil
}

func (link *Link) SelectByAlternative(alter string) error {
    dlink := dbase.Link{}
    err := dlink.FindByAlternative(alter)
    if nil != err {
        return err
    }
    link.Map(dlink)
    return nil
}

func (link *Link) Delete() error {
    dlink := link.UnMap()
    return dlink.Delete()
}

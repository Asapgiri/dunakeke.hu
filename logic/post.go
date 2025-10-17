package logic

import (
	"dunakeke/dbase"
	"dunakeke/dictionary"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostNew(dict dictionary.Dictionary, author User) string {
    post := Post{
        Author: author,
        Date: time.Now(),
        Markdown: dict.Editor.StartMessage,
        Image: "/placeholder.jpg",
    }
    post.EditDate = post.Date

    post.Add()

    return post.Id
}

func PostUpdate(ps PostSave) error {
    log.Println("Post Save Body:")
    log.Println(ps)

    post := Post{}
    err := post.Select(ps.Id)
    if nil != err {
        log.Println(err)
        return err
    }

    post.EditDate   = time.Now()
    post.Title      = ps.Title
    post.Markdown   = ps.Markdown
    post.Html       = ps.Html

    return post.Update()
}


func (post *Post) List(show_private bool) []Post {
    dpost := dbase.Post{}
    var dposts []dbase.Post
    var err error
    if show_private {
        dposts, err = dpost.List()
    } else {
        dposts, err = dpost.ListPublic()
    }
    if nil != err {
        log.Println(err)
        return []Post{}
    }

    posts := make([]Post, len(dposts))
    for i, p := range(dposts) {
        posts[i].Map(p)
    }

    return posts
}

func (post *Post) Add() error {
    dpost := post.UnMap()
    dpost.Id = primitive.NewObjectID()
    post.Id = dpost.Id.Hex()
    return dpost.Add()
}

func (post *Post) Select(id string) error {
    dpost := dbase.Post{}
    oid, _ := primitive.ObjectIDFromHex(id)
    err := dpost.Select(oid)
    if nil != err {
        return err
    }

    post.Map(dpost)
    return nil
}

func (post *Post) Update() error {
    dpost := post.UnMap()
    return dpost.Update()
}

func (post *Post) Delete() error {
    dpost := post.UnMap()
    // FIXME: Delete comments and stuff as well..
    return dpost.Delete()
}

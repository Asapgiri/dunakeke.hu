package logic

import (
	"dunakeke/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (user *User) Map(duser dbase.User) {
    user.Id         = duser.Id.Hex()
    user.RegDate    = duser.RegDate
    user.EditDate   = duser.EditDate
    user.Username   = duser.Username
    user.Name       = duser.Name
    user.Email      = duser.Email
    user.Phone      = duser.Phone
    user.Roles      = duser.Roles
}

func (user *User) UnMap() dbase.User {
    duser := dbase.User{}

    duser.Id, _      = primitive.ObjectIDFromHex(user.Id)
    duser.RegDate    = user.RegDate
    duser.EditDate   = user.EditDate
    duser.Username   = user.Username
    duser.Name       = user.Name
    duser.Email      = user.Email
    duser.Phone      = user.Phone
    duser.Roles      = user.Roles

    return duser
}

func (post *Post) Map(dpost dbase.Post) {
    author := User{}
    author.Find(dpost.Author)

    post.Id         = dpost.Id.Hex()
    post.Author     = author
    post.Date       = dpost.Date
    post.EditDate   = dpost.EditDate
    post.Title      = dpost.Title
    post.Short      = dpost.Short
    post.Image      = dpost.Image
    post.Content    = dpost.Content

    // Comments should be loaded separately
    // post.Comments   = comments
}

package logic

import (
	"dunakeke/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (user *User) Find(id primitive.ObjectID) {
    duser := dbase.User{}
    err := duser.Select(id)

    if nil != err {
        user.Id = ""
        return
    }

    user.Map(duser)
}

func (user *User) FindByUsername(username string) {
    duser := dbase.User{}
    err := duser.FindByUsername(username)

    if nil != err {
        user.Username = ""
        user.Id = ""
        return
    }

    user.Map(duser)
}

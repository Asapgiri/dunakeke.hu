package logic

import (
	"dunakeke/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (tag *Tag) List() ([]Tag, error) {
    dtag := dbase.Tag{}
    dtags, err := dtag.List()
    if nil != err {
        return []Tag{}, err
    }

    tags := make([]Tag, len(dtags))
    for i, dt := range(dtags) {
        tags[i].Map(dt)
    }

    return tags, nil
}

func (tag *Tag) Select(id string) error {
    dtag := dbase.Tag{}
    oid, _ := primitive.ObjectIDFromHex(id)
    err := dtag.Select(oid)
    if nil != err {
        return err
    }

    tag.Map(dtag)
    return nil
}

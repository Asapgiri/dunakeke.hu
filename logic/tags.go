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

func (tag *Tag) SelectByName(name string) error {
    dtag := dbase.Tag{}
    err := dtag.SelectByName(name)
    if nil != err {
        return err
    }

    tag.Map(dtag)
    return nil
}

func (tag *Tag) Add() error {
    dtag := tag.UnMap()
    dtag.Id = primitive.NewObjectID()
    tag.Id = dtag.Id.Hex()
    return dtag.Add()
}

func (tag *Tag) Update() error {
    dtag := tag.UnMap()
    return dtag.Update()
}

func (tag *Tag) Delete() error {
    dtag := tag.UnMap()
    return dtag.Delete()
}

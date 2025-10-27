package logic

import (
	"dunakeke/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (supporter *Supporter) List() ([]Supporter, error) {
    dsupporter := dbase.Supporter{}
    dsupporters, err := dsupporter.List()
    if nil != err {
        return []Supporter{}, err
    }

    supporters := make([]Supporter, len(dsupporters))
    for i, dt := range(dsupporters) {
        supporters[i].Map(dt)
    }

    return supporters, nil
}

func (supporter *Supporter) Select(id string) error {
    dsupporter := dbase.Supporter{}
    oid, _ := primitive.ObjectIDFromHex(id)
    err := dsupporter.Select(oid)
    if nil != err {
        return err
    }

    supporter.Map(dsupporter)
    return nil
}

func (supporter *Supporter) Add() error {
    dsupporter := supporter.UnMap()
    dsupporter.Id = primitive.NewObjectID()
    supporter.Id = dsupporter.Id.Hex()
    return dsupporter.Add()
}

func (supporter *Supporter) Update() error {
    dsupporter := supporter.UnMap()
    return dsupporter.Update()
}

func (supporter *Supporter) Delete() error {
    dsupporter := supporter.UnMap()
    return dsupporter.Delete()
}

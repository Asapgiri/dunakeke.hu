package logic

import (
	"dunakeke/dbase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (donation *Donation) List() []Donation {
    ddon := dbase.Donation{}
    ddons, _ := ddon.List()

    donations := make([]Donation, len(ddons))
    for i, d := range(ddons) {
        donations[i].Map(d)
    }

    return donations
}

func (donation *Donation) Add() error {
    ddon := donation.UnMap()
    ddon.Id = primitive.NewObjectID()

    return ddon.Add()
}

func (donation *Donation) Update() error {
    ddon := donation.UnMap()
    return ddon.Update()
}

func (donation *Donation) Select() error {
    ddon := dbase.Donation{}
    oid, _ := primitive.ObjectIDFromHex(donation.Id)
    err := ddon.Select(oid)
    if nil != err {
        return err
    }

    donation.Map(ddon)
    return nil
}


func (do *DonationOption) List() []DonationOption {
    ddon := dbase.DonationOption{}
    ddons, err := ddon.List()
    if nil != err {
        log.Println(err)
        return []DonationOption{}
    }

    donations := make([]DonationOption, len(ddons))
    for i, d := range(ddons) {
        donations[i].Map(d)
    }

    return donations
}

func (do *DonationOption) Add() error {
    ddon := do.UnMap()
    ddon.Id = primitive.NewObjectID()
    return ddon.Add()
}

func (do *DonationOption) Update() error {
    ddon := do.UnMap()
    return ddon.Update()
}

func (do *DonationOption) Select() error {
    ddon := dbase.DonationOption{}
    oid, _ := primitive.ObjectIDFromHex(do.Id)
    err := ddon.Select(oid)
    if nil != err {
        return err
    }

    do.Map(ddon)
    return nil
}

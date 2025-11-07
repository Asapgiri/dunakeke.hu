package logic

import (
	"dunakeke/dbase"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (file *File) List() ([]File, error) {
    return file.ListByExtension("")
}

func (file *File) ListByExtension(ext string) ([]File, error) {
    dfile := dbase.File{}
    var dfiles []dbase.File
    var err error

    if "" == ext {
        dfiles, err = dfile.List()
    } else {
        dfiles, err = dfile.ListByExtension(ext)
    }

    if nil != err {
        return []File{}, err
    }

    files := make([]File, len(dfiles))
    for i, dt := range(dfiles) {
        files[i].Map(dt)
    }

    return files, nil
}

func (file *File) Select(id string) error {
    dfile := dbase.File{}
    oid, _ := primitive.ObjectIDFromHex(id)
    err := dfile.Select(oid)
    if nil != err {
        return err
    }

    file.Map(dfile)
    return nil
}

func (file *File) Add() error {
    dfile := file.UnMap()
    return dfile.Add()
}

func (file *File) Update() error {
    dfile := file.UnMap()
    return dfile.Update()
}

func (file *File) Delete() error {
    dfile := file.UnMap()
    err := os.Remove("artifacts"+file.SaveName)
    if nil != err {
        return err
    }
    return dfile.Delete()
}

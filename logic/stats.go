package logic

import (
	"dunakeke/dbase"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveStatistics(r *http.Request, userid string) {
    stat := SiteStatistic{
        Date        : time.Now(),
        UserId      : userid,

        Method      : r.Method,
        Url         : r.URL.Path,

        Proto       : r.Proto,
        ProtoMajor  : r.ProtoMajor,
        ProtoMinor  : r.ProtoMinor,
        Header      : r.Header,

        Host        : r.Host,
        RemoteAddr  : r.RemoteAddr,
        RequestURI  : r.RequestURI,

        Referer     : r.Referer(),
        Pattern     : r.Pattern,
    }

    stat.Add()
}

func (s *SiteStatistic) List() []SiteStatistic {
    ds := dbase.SiteStatistic{}

    dsls, err := ds.List()
    if nil != err {
        log.Println(err)
        return []SiteStatistic{}
    }

    sls := make([]SiteStatistic, len(dsls))
    for i, s := range(dsls) {
        sls[i].Map(s)
    }

    return sls
}

func (s *SiteStatistic) Select(id string) error {
    ds := dbase.SiteStatistic{}
    oid, _ := primitive.ObjectIDFromHex(id)
    err := ds.Select(oid)
    if nil != err {
        return err
    }

    s.Map(ds)
    return nil}

func (s *SiteStatistic) Add() error {
    ds := s.UnMap()
    ds.Id = primitive.NewObjectID()
    s.Id = ds.Id.Hex()
    return ds.Add()
}

func (s *SiteStatistic) Update() error {
    ds := s.UnMap()
    return ds.Update()
}

func (s *SiteStatistic) Delete() error {
    ds := s.UnMap()
    return ds.Delete()
}

package logic

import (
	"net/http"
	"time"
    "dunakeke/dbase"
)

type SiteStatistic struct {
    _db         dbase.SiteStatistic

    Id          string
    Date        time.Time
    UserId      string

    Method      string
    Url         string

	Proto       string
	ProtoMajor  int
	ProtoMinor  int
	Header      http.Header

	Host        string
	RemoteAddr  string
	RequestURI  string

    Referer     string
    Pattern     string
}

type User struct {
    _db             dbase.User
    Id              string
    RegDate         time.Time
    EditDate        time.Time
    Username        string
    Name            string
    Email           string
    Phone           string
    Roles           []string
}

type Comment struct {
    _db             dbase.Comment
    Id              string
    Author          User
    Date            time.Time
    EditDate        time.Time
    Name            string
    Content         string
}

type Post struct {
    _db             dbase.Post
    Id              string
    Author          User
    Date            time.Time
    EditDate        time.Time
    Public          bool
    Path            string
    Title           string
    Short           string
    Image           string
    Markdown        string
    Html            string
    Alternative     Link
    Tags            []Tag
}

type Tag struct {
    _db             dbase.Tag
    Id              string
    Name            string
    Listable        bool
    Color           string
}

type Link struct {
    _db             dbase.Link
    Id              string
    Author          User
    Date            time.Time
    Original        string
    Alternative     string
}

type Newsletter struct {
    _db             dbase.Newsletter
    Id              string
    User            User
    RegDate         time.Time
    Email           string
}

type Donation struct {
    _db             dbase.Donation
    Id              string
    UserId          string
    Tokens          []string
    Name            string
    Message         string
    Email           string
    Date            time.Time
    Amount          float64
    Status          string
    Successful      bool
    Recurring       bool
    RecurringActive bool
    Occurences      []time.Time
    Newsletter      bool
    GDPR            bool
    TransactionId   int
}

type DonationOption struct {
    _db             dbase.DonationOption
    Id              string
    Amount          float64
}

type Stat struct {
    Id              string
    User            User
    Date            time.Time
    Ip              string
    Route           string
    Post            bool
}

type PostSave struct {
    Id          string      `json:"id"`
    Title       string      `json:"title"`
    Markdown    string      `json:"markdown"`
    Html        string      `json:"html"`
    Alternative string      `json:"alternative"`
    Tags        []string    `json:"tags"`
}

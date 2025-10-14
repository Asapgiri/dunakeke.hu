package dbase

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SiteStatistic struct {
    Id          primitive.ObjectID  `bson:"_id"`
    UserId      primitive.ObjectID
    Date        time.Time
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
    Id              primitive.ObjectID `bson:"_id"`
    RegDate         time.Time
    EditDate        time.Time
    Username        string             `bson:"username"`
    PasswordHash    string
    Name            string
    Email           string             `bson:"email"`
    Phone           string
    Roles           []string
}

type Post struct {
    Id              primitive.ObjectID `bson:"_id"`
    Author          primitive.ObjectID
    Date            time.Time
    EditDate        time.Time
    Public          bool
    Path            string
    Title           string
    Short           string
    Image           string
    Markdown        string
    Html            string
}

type Photo struct {
    Id              primitive.ObjectID `bson:"_id"`
}

type Comment struct {
    Id              primitive.ObjectID `bson:"_id"`
    Post            primitive.ObjectID
    Author          primitive.ObjectID
    Date            time.Time
    EditDate        time.Time
    Name            string
    Content         string
}

type Link struct {
    Id              primitive.ObjectID `bson:"_id"`
    Author          primitive.ObjectID
    Date            time.Time
    Link            string
    Alternative     string
}

type Newsletter struct {
    Id              primitive.ObjectID `bson:"_id"`
    User            primitive.ObjectID
    RegDate         time.Time
    Email           string
}

type DonationInvoice struct {
    Id          primitive.ObjectID `bson:"_id"`
    Name        string
    Company     string
    Country     string
    State       string
    City        string
    Zip         string
    Address     string
    Address2    string
    Phone       string
    TaxNumber   string
}

type Donation struct {
    Id              primitive.ObjectID `bson:"_id"`
    User            primitive.ObjectID
    Tokens          []string
    Name            string
    Email           string
    Date            time.Time
    Amount          float64
    Status          string
    Successful      bool
    Recurring       bool
    RecurringActive bool
    Occurences      []time.Time
    Newsletter      bool
    Gdpr            bool
    InvoiceNeeded   bool
    Invoice         primitive.ObjectID
    TransactionId   int
}

type DonationOption struct {
    Id              primitive.ObjectID `bson:"_id"`
    Date            time.Time
    Amount          float64
}

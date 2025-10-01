package logic

import (
	"time"
)

type User struct {
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
    Id              string
    Author          User
    Date            time.Time
    EditDate        time.Time
    Name            string
    Content         string
}

type Post struct {
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
}

type Link struct {
    Id              string
    Author          User
    Date            time.Time
    Link            string
    Alternative     string
}

type Newsletter struct {
    Id              string
    User            User
    RegDate         time.Time
    Email           string
}

type Invoice struct {
    Id          string
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
    Id              string
    UserId          string
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
    GDPR            bool
    InvoiceNeeded   bool
    Invoice         Invoice
    TransactionId   int
}

type DonationOption struct {
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
    Id          string  `json:"id"`
    Title       string  `json:"title"`
    Markdown    string  `json:"markdown"`
    Html        string  `json:"html"`
}

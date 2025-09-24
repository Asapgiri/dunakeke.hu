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
    Title           string
    Short           string
    Image           string
    Content         string
    Comments        []Comment
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

type Donation struct {
    Id              string
    User            User
    Name            string
    Email           string
    Date            time.Time
    Amount          float64
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

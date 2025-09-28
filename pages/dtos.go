package pages

import "dunakeke/logic"

type DtoMain struct {
    Session string
    // etc...
}

type DtoRoot struct {
    Main    DtoMain
    Posts   []logic.Post
}

type DtoAdminDonations struct {
    Sum         float64
    Donations   []logic.Donation
    Username    []string
}

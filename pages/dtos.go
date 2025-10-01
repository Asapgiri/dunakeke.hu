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

type DtoEditor struct {
    Id          string
    Markdown    string
}

type DtoPostShow struct {
    Post        logic.Post
    Comments    []logic.Comment
}

type PostSaveImageResponse struct {
    Success int     `json:"success"`
    Url     string  `json:"url"`
}

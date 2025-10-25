package pages

import "dunakeke/logic"

type DtoMain struct {
    Title   string
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

type DtoAdminUsers struct {
    Roles   []logic.RolePerm
    Users   []logic.User
}

type DtoTag struct {
    Tag         logic.Tag
    Selected    bool
}

type DtoEditor struct {
    Post    logic.Post
    Tags    []DtoTag
}

type DtoPostShow struct {
    Post        logic.Post
    Comments    []logic.Comment
}

type PostSaveImageResponse struct {
    Success int     `json:"success"`
    Url     string  `json:"url"`
}

type LinkUpdate struct {
    Original    string  `json:"original"`
    Alternative string  `json:"alternative"`
}

type TagUpdate struct {
    Id      string `json:"id"`
    Name    string `json:"name"`
    Color   string `json:"color"`
}

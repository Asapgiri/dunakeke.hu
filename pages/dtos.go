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

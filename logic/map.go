package logic

import (
	"dunakeke/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *SiteStatistic) Map(dstat dbase.SiteStatistic) {
    s._db           = dstat
    s.Id            = dstat.Id.Hex()
    s.Date          = dstat.Date
    s.UserId        = dstat.UserId.Hex()
    s.Method        = dstat.Method
    s.Url           = dstat.Url
	s.Proto         = dstat.Proto
	s.ProtoMajor    = dstat.ProtoMajor
	s.ProtoMinor    = dstat.ProtoMinor
	s.Header        = dstat.Header
	s.Host          = dstat.Host
	s.RemoteAddr    = dstat.RemoteAddr
	s.RequestURI    = dstat.RequestURI
    s.Referer       = dstat.Referer
    s.Pattern       = dstat.Pattern
}

func (s *SiteStatistic) UnMap() dbase.SiteStatistic {
    dstat := s._db

    dstat.Id, _         = primitive.ObjectIDFromHex(s.Id)
    dstat.Date          = s.Date
    dstat.UserId, _     = primitive.ObjectIDFromHex(s.UserId)
    dstat.Method        = s.Method
    dstat.Url           = s.Url
	dstat.Proto         = s.Proto
	dstat.ProtoMajor    = s.ProtoMajor
	dstat.ProtoMinor    = s.ProtoMinor
	dstat.Header        = s.Header
	dstat.Host          = s.Host
	dstat.RemoteAddr    = s.RemoteAddr
	dstat.RequestURI    = s.RequestURI
    dstat.Referer       = s.Referer
    dstat.Pattern       = s.Pattern

    return dstat
}

func (user *User) Map(duser dbase.User) {
    user._db        = duser
    user.Id         = duser.Id.Hex()
    user.RegDate    = duser.RegDate
    user.EditDate   = duser.EditDate
    user.Username   = duser.Username
    user.Name       = duser.Name
    user.Email      = duser.Email
    user.Phone      = duser.Phone
    user.Roles      = duser.Roles
}

func (user *User) UnMap() dbase.User {
    duser := user._db

    duser.RegDate    = user.RegDate
    duser.EditDate   = user.EditDate
    duser.Username   = user.Username
    duser.Name       = user.Name
    duser.Email      = user.Email
    duser.Phone      = user.Phone
    duser.Roles      = user.Roles

    return duser
}

func (post *Post) Map(dpost dbase.Post) {
    author := User{}
    author.Find(dpost.Author.Hex())
    alternative := Link{}
    alternative.Select(dpost.Alternative.Hex())

    post._db            = dpost
    post.Id             = dpost.Id.Hex()
    post.Author         = author
    post.Date           = dpost.Date
    post.EditDate       = dpost.EditDate
    post.Public         = dpost.Public
    post.Path           = alternative.Alternative
    post.Title          = dpost.Title
    post.Short          = dpost.Short
    post.Image          = dpost.Image
    post.Markdown       = dpost.Markdown
    post.Html           = dpost.Html
    post.Alternative    = alternative
}

func (post *Post) UnMap() dbase.Post {
    dpost := post._db

    dpost.Id, _             = primitive.ObjectIDFromHex(post.Id)
    dpost.Author, _         = primitive.ObjectIDFromHex(post.Author.Id)
    dpost.Date              = post.Date
    dpost.EditDate          = post.EditDate
    dpost.Public            = post.Public
    dpost.Title             = post.Title
    dpost.Short             = post.Short
    dpost.Image             = post.Image
    dpost.Markdown          = post.Markdown
    dpost.Html              = post.Html
    dpost.Alternative, _    = primitive.ObjectIDFromHex(post.Alternative.Id)

    return dpost
}

func (donation *Donation) Map(ddon dbase.Donation) {
    donation._db                = ddon
    donation.Id                 = ddon.Id.Hex()
    donation.UserId             = ddon.User.Hex()
    donation.Tokens             = ddon.Tokens
    donation.Name               = ddon.Name
    donation.Email              = ddon.Email
    donation.Date               = ddon.Date
    donation.Amount             = ddon.Amount
    donation.Status             = ddon.Status
    donation.Successful         = ddon.Successful
    donation.Recurring          = ddon.Recurring
    donation.RecurringActive    = ddon.RecurringActive
    donation.Occurences         = ddon.Occurences
    donation.Newsletter         = ddon.Newsletter
    donation.GDPR               = ddon.Gdpr
    donation.TransactionId      = ddon.TransactionId
}

func (donation *Donation) UnMap() dbase.Donation {
    ddon := donation._db

    ddon.Id, _              = primitive.ObjectIDFromHex(donation.Id)
    ddon.User, _            = primitive.ObjectIDFromHex(donation.UserId)
    ddon.Tokens             = donation.Tokens
    ddon.Name               = donation.Name
    ddon.Email              = donation.Email
    ddon.Date               = donation.Date
    ddon.Amount             = donation.Amount
    ddon.Status             = donation.Status
    ddon.Successful         = donation.Successful
    ddon.Recurring          = donation.Recurring
    ddon.RecurringActive    = donation.RecurringActive
    ddon.Occurences         = donation.Occurences
    ddon.Newsletter         = donation.Newsletter
    ddon.Gdpr               = donation.GDPR
    ddon.TransactionId      = donation.TransactionId

    return ddon
}

func (do *DonationOption) Map(ddon dbase.DonationOption) {
    do._db    = ddon
    do.Id     = ddon.Id.Hex()
    do.Amount = ddon.Amount
}

func (do *DonationOption) UnMap() dbase.DonationOption {
    ddon := do._db

    ddon.Id, _      = primitive.ObjectIDFromHex(do.Id)
    ddon.Amount     = do.Amount

    return ddon
}

func (link *Link)Map(dlink dbase.Link) {
    author := User{}
    author.Find(dlink.Author.Hex())

    link._db            = dlink
    link.Id             = dlink.Id.Hex()
    link.Author         = author
    link.Date           = dlink.Date
    link.Original       = dlink.Original
    link.Alternative    = dlink.Alternative
}

func (link *Link)UnMap() dbase.Link {
    dlink := link._db

    dlink.Id, _         = primitive.ObjectIDFromHex(link.Id)
    dlink.Author, _     = primitive.ObjectIDFromHex(link.Author.Id)
    dlink.Date          = link.Date
    dlink.Original      = link.Original
    dlink.Alternative   = link.Alternative

    return dlink
}

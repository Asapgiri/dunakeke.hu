package logic

import (
	"dunakeke/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (user *User) Map(duser dbase.User) {
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
    duser := dbase.User{}

    duser.Id, _      = primitive.ObjectIDFromHex(user.Id)
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

    post.Id         = dpost.Id.Hex()
    post.Author     = author
    post.Date       = dpost.Date
    post.EditDate   = dpost.EditDate
    post.Public     = dpost.Public
    post.Path       = dpost.Path
    post.Title      = dpost.Title
    post.Short      = dpost.Short
    post.Image      = dpost.Image
    post.Markdown   = dpost.Markdown
    post.Html       = dpost.Html
}

func (post *Post) UnMap() dbase.Post {
    dpost := dbase.Post{}

    dpost.Id, _     = primitive.ObjectIDFromHex(post.Id)
    dpost.Author, _ = primitive.ObjectIDFromHex(post.Author.Id)
    dpost.Date      = post.Date
    dpost.EditDate  = post.EditDate
    dpost.Public    = post.Public
    dpost.Path      = post.Path
    dpost.Title     = post.Title
    dpost.Short     = post.Short
    dpost.Image     = post.Image
    dpost.Markdown  = post.Markdown
    dpost.Html      = post.Html

    return dpost
}

func (donation *Donation) Map(ddon dbase.Donation) {
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
    donation.InvoiceNeeded      = ddon.InvoiceNeeded
    donation.TransactionId      = ddon.TransactionId

    invoice := dbase.DonationInvoice{}
    invoice.Select(ddon.Invoice)
    donation.Invoice.Map(invoice)
}

func (donation *Donation) UnMap() dbase.Donation {
    ddon := dbase.Donation{}

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
    ddon.InvoiceNeeded      = donation.InvoiceNeeded
    ddon.Invoice, _         = primitive.ObjectIDFromHex(donation.Invoice.Id)
    ddon.TransactionId      = donation.TransactionId

    return ddon
}

func (do *DonationOption) Map(ddon dbase.DonationOption) {
    do.Id     = ddon.Id.Hex()
    do.Amount = ddon.Amount
}

func (do *DonationOption) UnMap() dbase.DonationOption {
    ddon := dbase.DonationOption{}

    ddon.Id, _      = primitive.ObjectIDFromHex(do.Id)
    ddon.Amount     = do.Amount

    return ddon
}

func (iv *Invoice) Map(div dbase.DonationInvoice) {
    iv.Id          = div.Id.Hex()
    iv.Name        = div.Name
    iv.Company     = div.Company
    iv.Country     = div.Country
    iv.State       = div.State
    iv.City        = div.City
    iv.Zip         = div.Zip
    iv.Address     = div.Address
    iv.Address2    = div.Address2
    iv.Phone       = div.Phone
    iv.TaxNumber   = div.TaxNumber
}

func (iv *Invoice) UnMap() dbase.DonationInvoice {
    div := dbase.DonationInvoice{}

    div.Id, _       = primitive.ObjectIDFromHex(iv.Id)
    div.Name        = iv.Name
    div.Company     = iv.Company
    div.Country     = iv.Country
    div.State       = iv.State
    div.City        = iv.City
    div.Zip         = iv.Zip
    div.Address     = iv.Address
    div.Address2    = iv.Address2
    div.Phone       = iv.Phone
    div.TaxNumber   = iv.TaxNumber

    return div
}

package dbase

import (
	"context"
	"dunakeke/config"
	"dunakeke/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongo_client *mongo.Client
var db *mongo.Database

var dbUSERS         *mongo.Collection
var dbPOSTS         *mongo.Collection
var dbCOMMENTS      *mongo.Collection
var dbLINKS         *mongo.Collection
var dbNEWSLETTER    *mongo.Collection
var dbDONATIONS     *mongo.Collection
var dbDONATIONOPTS  *mongo.Collection
var dbSTATISTICS    *mongo.Collection

var log = logger.Logger {
    Color: logger.Colors.Purple,
    Pretext: "database",
}

// =====================================================================================================================
// Basic connect and stuff

func Connect() error {
    var err error

    // Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Config.Dbase.Url).SetServerAPIOptions(serverAPI)

    // Create a new client and connect to the server
    mongo_client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
        return err
	}
    db = mongo_client.Database(config.Config.Dbase.Name)

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

    dbUSERS         = db.Collection("users")
    dbPOSTS         = db.Collection("posts")
    dbCOMMENTS      = db.Collection("comments")
    dbLINKS         = db.Collection("links")
    dbNEWSLETTER    = db.Collection("newsletter")
    dbDONATIONS     = db.Collection("donations")
    dbDONATIONOPTS  = db.Collection("donation-options")
    dbSTATISTICS    = db.Collection("statistics")

    return nil
}

// =====================================================================================================================
// Internal User Listing CRUD

func (user *User) List() ([]User, error) {
    var anyime []User

    cursor, err := dbUSERS.Find(context.TODO(), bson.D{{}})
    if nil != err {
        return anyime, err
    }

    err = cursor.All(context.TODO(), &anyime)

    return anyime, err
}

func (user *User) Select(id primitive.ObjectID) error {
    return dbUSERS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(user)
}

func (user *User) FindByUsername(username string) error {
    return dbUSERS.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(user)
}

func (user *User) FindByEmail(email string) error {
    return dbUSERS.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(user)
}

func (user *User) Add() error {
    _, err := dbUSERS.InsertOne(context.TODO(), user)
    return err
}

func (user *User) Update() error {
    _, err := dbUSERS.ReplaceOne(context.TODO(), bson.D{{"_id", user.Id}}, user)
    return err
}

func (user *User) Delete() error {
    filter := bson.D{{"_id", user.Id}}
    _, err := dbUSERS.DeleteOne(context.TODO(), filter)
    return err
}

// =====================================================================================================================
// Internal Post CRUD

func (post *Post) List() ([]Post, error) {
    var posts []Post
    cursor, err := dbPOSTS.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return posts, err
    }
    err = cursor.All(context.TODO(), &posts)
    return posts, err
}

func (post *Post) Select(id primitive.ObjectID) error {
    return dbPOSTS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(post)
}

func (post *Post) Add() error {
    _, err := dbPOSTS.InsertOne(context.TODO(), post)
    return err
}

func (post *Post) Update() error {
    _, err := dbPOSTS.ReplaceOne(context.TODO(), bson.D{{"_id", post.Id}}, post)
    return err
}

func (post *Post) Delete() error {
    _, err := dbPOSTS.DeleteOne(context.TODO(), bson.D{{"_id", post.Id}})
    return err
}

// =====================================================================================================================
// Internal Comment CRUD

func (comment *Comment) List() ([]Comment, error) {
    var comments []Comment
    cursor, err := dbCOMMENTS.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return comments, err
    }
    err = cursor.All(context.TODO(), &comments)
    return comments, err
}

func (comment *Comment) Select(id primitive.ObjectID) error {
    return dbCOMMENTS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(comment)
}

func (comment *Comment) ListByPost(postID primitive.ObjectID) ([]Comment, error) {
    var comments []Comment
    cursor, err := dbCOMMENTS.Find(context.TODO(), bson.D{{"Post", postID}})
    if err != nil {
        return comments, err
    }
    err = cursor.All(context.TODO(), &comments)
    return comments, err
}

func (comment *Comment) Add() error {
    _, err := dbCOMMENTS.InsertOne(context.TODO(), comment)
    return err
}

func (comment *Comment) Update() error {
    _, err := dbCOMMENTS.ReplaceOne(context.TODO(), bson.D{{"_id", comment.Id}}, comment)
    return err
}

func (comment *Comment) Delete() error {
    _, err := dbCOMMENTS.DeleteOne(context.TODO(), bson.D{{"_id", comment.Id}})
    return err
}

// =====================================================================================================================
// Internal Link CRUD

func (link *Link) List() ([]Link, error) {
    var links []Link
    cursor, err := dbLINKS.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return links, err
    }
    err = cursor.All(context.TODO(), &links)
    return links, err
}

func (link *Link) Select(id primitive.ObjectID) error {
    return dbLINKS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(link)
}

func (link *Link) FindByLink(qlink string) error {
    return dbLINKS.FindOne(context.TODO(), bson.D{{"link", qlink}}).Decode(link)
}

func (link *Link) Add() error {
    _, err := dbLINKS.InsertOne(context.TODO(), link)
    return err
}

func (link *Link) Update() error {
    _, err := dbLINKS.ReplaceOne(context.TODO(), bson.D{{"_id", link.Id}}, link)
    return err
}

func (link *Link) Delete() error {
    _, err := dbLINKS.DeleteOne(context.TODO(), bson.D{{"_id", link.Id}})
    return err
}

// =====================================================================================================================
// Internal Newsletter CRUD

func (newsletter *Newsletter) List() ([]Newsletter, error) {
    var newsletters []Newsletter
    cursor, err := dbNEWSLETTER.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return newsletters, err
    }
    err = cursor.All(context.TODO(), &newsletters)
    return newsletters, err
}

func (newsletter *Newsletter) Select(id primitive.ObjectID) error {
    return dbNEWSLETTER.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(newsletter)
}

func (newsletter *Newsletter) Add() error {
    _, err := dbNEWSLETTER.InsertOne(context.TODO(), newsletter)
    return err
}

func (newsletter *Newsletter) Update() error {
    _, err := dbNEWSLETTER.ReplaceOne(context.TODO(), bson.D{{"_id", newsletter.Id}}, newsletter)
    return err
}

func (newsletter *Newsletter) Delete() error {
    _, err := dbNEWSLETTER.DeleteOne(context.TODO(), bson.D{{"_id", newsletter.Id}})
    return err
}

// =====================================================================================================================
// Internal Donation CRUD

func (donation *Donation) List() ([]Donation, error) {
    var donations []Donation
    cursor, err := dbDONATIONS.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return donations, err
    }
    err = cursor.All(context.TODO(), &donations)
    return donations, err
}

func (donation *Donation) Select(id primitive.ObjectID) error {
    return dbDONATIONS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(donation)
}

func (donation *Donation) Add() error {
    _, err := dbDONATIONS.InsertOne(context.TODO(), donation)
    return err
}

func (donation *Donation) Update() error {
    _, err := dbDONATIONS.ReplaceOne(context.TODO(), bson.D{{"_id", donation.Id}}, donation)
    return err
}

func (donation *Donation) Delete() error {
    _, err := dbDONATIONS.DeleteOne(context.TODO(), bson.D{{"_id", donation.Id}})
    return err
}

func (do *DonationOption) List() ([]DonationOption, error) {
    var donations []DonationOption
    cursor, err := dbDONATIONOPTS.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return donations, err
    }
    err = cursor.All(context.TODO(), &donations)
    return donations, err
}

func (do *DonationOption) Select(id primitive.ObjectID) error {
    return dbDONATIONOPTS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(do)
}

func (do *DonationOption) Add() error {
    _, err := dbDONATIONOPTS.InsertOne(context.TODO(), do)
    return err
}

func (do *DonationOption) Update() error {
    _, err := dbDONATIONOPTS.ReplaceOne(context.TODO(), bson.D{{"_id", do.Id}}, do)
    return err
}

func (do *DonationOption) Delete() error {
    _, err := dbDONATIONOPTS.DeleteOne(context.TODO(), bson.D{{"_id", do.Id}})
    return err
}

// =====================================================================================================================
// Internal Stat CRUD

func (stat *Stat) List() ([]Stat, error) {
    var stats []Stat
    cursor, err := dbSTATISTICS.Find(context.TODO(), bson.D{{}})
    if err != nil {
        return stats, err
    }
    err = cursor.All(context.TODO(), &stats)
    return stats, err
}

func (stat *Stat) Select(id primitive.ObjectID) error {
    return dbSTATISTICS.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(stat)
}

func (stat *Stat) Add() error {
    _, err := dbSTATISTICS.InsertOne(context.TODO(), stat)
    return err
}

func (stat *Stat) Update() error {
    _, err := dbSTATISTICS.ReplaceOne(context.TODO(), bson.D{{"_id", stat.Id}}, stat)
    return err
}

func (stat *Stat) Delete() error {
    _, err := dbSTATISTICS.DeleteOne(context.TODO(), bson.D{{"_id", stat.Id}})
    return err
}

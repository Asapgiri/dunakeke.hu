package dbase

import (
	"github.com/asapgiri/golib/logger"
	"context"
	"dunakeke/config"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongo_client *mongo.Client
var db *mongo.Database

var dbUSERS             *mongo.Collection
var dbPOSTS             *mongo.Collection
var dbTAGS              *mongo.Collection
var dbFILES             *mongo.Collection
var dbCOMMENTS          *mongo.Collection
var dbLINKS             *mongo.Collection
var dbNEWSLETTER        *mongo.Collection
var dbDONATIONS         *mongo.Collection
var dbDONATIONOPTS      *mongo.Collection
var dbSTATISTICS        *mongo.Collection
var dbSUPPORTERS        *mongo.Collection

var log = logger.Logger {
    Color: logger.Colors.Purple,
    Pretext: "database",
}

func filter[T any](s []T, keep func(T) bool) []T {
    var result []T
    for _, v := range(s) {
        if keep(v) {
            result = append(result, v)
        }
    }
    return result
}

func count(db *mongo.Collection, pipeline mongo.Pipeline) int {
    pipeline = append(pipeline, bson.D{{Key: "$count", Value: "count"}})
    log.Println(pipeline)

    cursor, err := db.Aggregate(context.Background(), pipeline)
    if nil != err {
        log.Println(err)
        return 0
    }
    defer cursor.Close(context.Background())

    var result []bson.M
    err = cursor.All(context.Background(), &result)
    if nil != err {
        return 0
    }
    log.Println(result)

    var count int
    if len(result) > 0 {
        log.Println(result[0])
        log.Println(result[0]["count"])
        count = int(result[0]["count"].(int32))
    }

    return count
}

// =====================================================================================================================
// Basic connect and stuff

func Connect() error {
    var err error

    // Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Config.Dbase.Url).SetServerAPIOptions(serverAPI)

    // Create a new client and connect to the server
    mongo_client, err = mongo.Connect(context.Background(), opts)
	if err != nil {
        return err
	}
    db = mongo_client.Database(config.Config.Dbase.Name)

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := db.RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

    dbUSERS             = db.Collection("users")
    dbPOSTS             = db.Collection("posts")
    dbTAGS              = db.Collection("tags")
    dbFILES             = db.Collection("files")
    dbCOMMENTS          = db.Collection("comments")
    dbLINKS             = db.Collection("links")
    dbNEWSLETTER        = db.Collection("newsletter")
    dbDONATIONS         = db.Collection("donations")
    dbDONATIONOPTS      = db.Collection("donation-options")
    dbSTATISTICS        = db.Collection("statistics")
    dbSUPPORTERS        = db.Collection("supporters")

    return nil
}

// func (coll *mongo.Collection) List[T interface{Id primitive.ObjectID}](s T) ([]T, error) {
//     var ret []T
//     cursor, err := coll.Find(context.Background(), bson.D{{}})
//     if nil != err {
//         return ret, err
//     }
//     err = cursor.All(context.Background(), &ret)
//     return ret, err
// }

// =====================================================================================================================
// Internal User Listing CRUD

func (user *User) List() ([]User, error) {
    var anyime []User

    cursor, err := dbUSERS.Find(context.Background(), bson.D{{}})
    if nil != err {
        return anyime, err
    }
    defer cursor.Close(context.Background())

    err = cursor.All(context.Background(), &anyime)

    return anyime, err
}

func (user *User) Select(id primitive.ObjectID) error {
    return dbUSERS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(user)
}

func (user *User) FindByUsername(username string) error {
    return dbUSERS.FindOne(context.Background(), bson.D{{"username", username}}).Decode(user)
}

func (user *User) FindByEmail(email string) error {
    return dbUSERS.FindOne(context.Background(), bson.D{{"email", email}}).Decode(user)
}

func (user *User) Add() error {
    _, err := dbUSERS.InsertOne(context.Background(), user)
    return err
}

func (user *User) Update() error {
    _, err := dbUSERS.ReplaceOne(context.Background(), bson.D{{"_id", user.Id}}, user)
    return err
}

func (user *User) Delete() error {
    filter := bson.D{{"_id", user.Id}}
    _, err := dbUSERS.DeleteOne(context.Background(), filter)
    return err
}

// =====================================================================================================================
// Internal Post CRUD

func (post *Post) List(public_only bool, tagId *primitive.ObjectID, page int, limit int, admin bool) ([]Post, int, error) {
    var posts []Post
    var query = bson.D{}
    var pipeline = mongo.Pipeline{}

    if public_only {
        query = append(query, bson.E{Key: "public", Value: true})
    }
    if nil != tagId {
        query = append(query, bson.E{Key: "tags", Value: *tagId})
    }

    if len(query) > 0 {
        pipeline = append(pipeline, bson.D{{Key: "$match", Value: query}})
    }

    if !admin && nil == tagId {
        // Lookup tags
        pipeline = append(pipeline, bson.D{{
            Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "tags"},
                {Key: "localField", Value: "tags"},
                {Key: "foreignField", Value: "_id"},
                {Key: "as", Value: "tag_docs"},
            },
        }})

        // Match posts that have at least one listable tag
        pipeline = append(pipeline, bson.D{{
            Key: "$match", Value: bson.D{
                {Key: "$or", Value: bson.A{
                    bson.D{{Key: "tag_docs.listable", Value: true}},
                    bson.D{{Key: "tag_docs", Value: bson.D{{Key: "$size", Value: 0}}}},
                }},
            },
        }})
    }

    // Count before the filtering..
    full_len := count(dbPOSTS, pipeline)
    page_count := full_len / limit
    if full_len % limit > 0 {
        page_count++
    }

    cursor_start := page * limit
    if full_len < cursor_start {
        return posts, page_count, errors.New("Page does not exists")
    }

    pipeline = append(pipeline, bson.D{{Key: "$sort",   Value: bson.D{{Key: "date", Value: -1}}}},
                                bson.D{{Key: "$skip",   Value: cursor_start}},
                                bson.D{{Key: "$limit",  Value: limit}})

    posts = []Post{}
    cursor, err := dbPOSTS.Aggregate(context.Background(), pipeline)
    if err != nil {
        return posts, 0, err
    }
    defer cursor.Close(context.Background())

    err = cursor.All(context.Background(), &posts)
    return posts, page_count, err
}

func (post *Post) Select(id primitive.ObjectID) error {
    return dbPOSTS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(post)
}

func (post *Post) Add() error {
    _, err := dbPOSTS.InsertOne(context.Background(), post)
    return err
}

func (post *Post) Update() error {
    _, err := dbPOSTS.ReplaceOne(context.Background(), bson.D{{"_id", post.Id}}, post)
    return err
}

func (post *Post) Delete() error {
    _, err := dbPOSTS.DeleteOne(context.Background(), bson.D{{"_id", post.Id}})
    return err
}

// =====================================================================================================================
// Internal Tags CRUD

func (tag *Tag) List() ([]Tag, error) {
    var tags []Tag
    cursor, err := dbTAGS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return tags, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &tags)
    return tags, err
}

func (tag *Tag) Select(id primitive.ObjectID) error {
    return dbTAGS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(tag)
}

func (tag *Tag) SelectByName(name string) error {
    return dbTAGS.FindOne(context.Background(), bson.D{{"name", name}}).Decode(tag)
}

func (tag *Tag) Add() error {
    _, err := dbTAGS.InsertOne(context.Background(), tag)
    return err
}

func (tag *Tag) Update() error {
    _, err := dbTAGS.ReplaceOne(context.Background(), bson.D{{"_id", tag.Id}}, tag)
    return err
}

func (tag *Tag) Delete() error {
    _, err := dbTAGS.DeleteOne(context.Background(), bson.D{{"_id", tag.Id}})
    return err
}

// =====================================================================================================================
// Internal Photos CRUD

func (file *File) List() ([]File, error) {
    return file.ListByExtension("")
}

func (file *File) Select(id primitive.ObjectID) error {
    return dbFILES.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(file)
}

func (file *File) ListByExtension(ext string) ([]File, error) {
    var files []File
    query := bson.D{{}}
    if "" != ext {
        query = bson.D{{"extension", ext}}
    }

    cursor, err := dbFILES.Find(context.Background(), query)
    if err != nil {
        return files, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &files)
    return files, err
}

func (file *File) Add() error {
    _, err := dbFILES.InsertOne(context.Background(), file)
    return err
}

func (file *File) Update() error {
    _, err := dbFILES.ReplaceOne(context.Background(), bson.D{{"_id", file.Id}}, file)
    return err
}

func (file *File) Delete() error {
    _, err := dbFILES.DeleteOne(context.Background(), bson.D{{"_id", file.Id}})
    return err
}

// =====================================================================================================================
// Internal Comment CRUD

func (comment *Comment) List() ([]Comment, error) {
    var comments []Comment
    cursor, err := dbCOMMENTS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return comments, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &comments)
    return comments, err
}

func (comment *Comment) Select(id primitive.ObjectID) error {
    return dbCOMMENTS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(comment)
}

func (comment *Comment) ListByPost(postID primitive.ObjectID) ([]Comment, error) {
    var comments []Comment
    cursor, err := dbCOMMENTS.Find(context.Background(), bson.D{{"Post", postID}})
    if err != nil {
        return comments, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &comments)
    return comments, err
}

func (comment *Comment) Add() error {
    _, err := dbCOMMENTS.InsertOne(context.Background(), comment)
    return err
}

func (comment *Comment) Update() error {
    _, err := dbCOMMENTS.ReplaceOne(context.Background(), bson.D{{"_id", comment.Id}}, comment)
    return err
}

func (comment *Comment) Delete() error {
    _, err := dbCOMMENTS.DeleteOne(context.Background(), bson.D{{"_id", comment.Id}})
    return err
}

// =====================================================================================================================
// Internal Link CRUD

func (link *Link) List() ([]Link, error) {
    var links []Link
    cursor, err := dbLINKS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return links, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &links)
    return links, err
}

func (link *Link) Select(id primitive.ObjectID) error {
    return dbLINKS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(link)
}

func (link *Link) FindByOriginal(qlink string) error {
    return dbLINKS.FindOne(context.Background(), bson.D{{"original", qlink}}).Decode(link)
}

func (link *Link) FindByAlternative(alternative string) error {
    return dbLINKS.FindOne(context.Background(), bson.D{{"alternative", alternative}}).Decode(link)
}

func (link *Link) Add() error {
    _, err := dbLINKS.InsertOne(context.Background(), link)
    return err
}

func (link *Link) Update() error {
    _, err := dbLINKS.ReplaceOne(context.Background(), bson.D{{"_id", link.Id}}, link)
    return err
}

func (link *Link) Delete() error {
    _, err := dbLINKS.DeleteOne(context.Background(), bson.D{{"_id", link.Id}})
    return err
}

// =====================================================================================================================
// Internal Newsletter CRUD

func (newsletter *Newsletter) List() ([]Newsletter, error) {
    var newsletters []Newsletter
    cursor, err := dbNEWSLETTER.Find(context.Background(), bson.D{{}})
    if err != nil {
        return newsletters, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &newsletters)
    return newsletters, err
}

func (newsletter *Newsletter) Select(id primitive.ObjectID) error {
    return dbNEWSLETTER.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(newsletter)
}

func (newsletter *Newsletter) Add() error {
    _, err := dbNEWSLETTER.InsertOne(context.Background(), newsletter)
    return err
}

func (newsletter *Newsletter) Update() error {
    _, err := dbNEWSLETTER.ReplaceOne(context.Background(), bson.D{{"_id", newsletter.Id}}, newsletter)
    return err
}

func (newsletter *Newsletter) Delete() error {
    _, err := dbNEWSLETTER.DeleteOne(context.Background(), bson.D{{"_id", newsletter.Id}})
    return err
}

// =====================================================================================================================
// Internal Donation CRUD

func (donation *Donation) List() ([]Donation, error) {
    var donations []Donation
    cursor, err := dbDONATIONS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return donations, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &donations)
    return donations, err
}

func (donation *Donation) Select(id primitive.ObjectID) error {
    return dbDONATIONS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(donation)
}

func (donation *Donation) Add() error {
    _, err := dbDONATIONS.InsertOne(context.Background(), donation)
    return err
}

func (donation *Donation) Update() error {
    _, err := dbDONATIONS.ReplaceOne(context.Background(), bson.D{{"_id", donation.Id}}, donation)
    return err
}

func (donation *Donation) Delete() error {
    _, err := dbDONATIONS.DeleteOne(context.Background(), bson.D{{"_id", donation.Id}})
    return err
}

func (do *DonationOption) List() ([]DonationOption, error) {
    var donations []DonationOption
    cursor, err := dbDONATIONOPTS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return donations, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &donations)
    return donations, err
}

func (do *DonationOption) Select(id primitive.ObjectID) error {
    return dbDONATIONOPTS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(do)
}

func (do *DonationOption) Add() error {
    _, err := dbDONATIONOPTS.InsertOne(context.Background(), do)
    return err
}

func (do *DonationOption) Update() error {
    _, err := dbDONATIONOPTS.ReplaceOne(context.Background(), bson.D{{"_id", do.Id}}, do)
    return err
}

func (do *DonationOption) Delete() error {
    _, err := dbDONATIONOPTS.DeleteOne(context.Background(), bson.D{{"_id", do.Id}})
    return err
}

// =====================================================================================================================
// Internal Stat CRUD

func (stat *SiteStatistic) List() ([]SiteStatistic, error) {
    var stats []SiteStatistic
    cursor, err := dbSTATISTICS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return stats, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &stats)
    return stats, err
}

func (stat *SiteStatistic) Select(id primitive.ObjectID) error {
    return dbSTATISTICS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(stat)
}

func (stat *SiteStatistic) Add() error {
    _, err := dbSTATISTICS.InsertOne(context.Background(), stat)
    return err
}

func (stat *SiteStatistic) Update() error {
    _, err := dbSTATISTICS.ReplaceOne(context.Background(), bson.D{{"_id", stat.Id}}, stat)
    return err
}

func (stat *SiteStatistic) Delete() error {
    _, err := dbSTATISTICS.DeleteOne(context.Background(), bson.D{{"_id", stat.Id}})
    return err
}

// =====================================================================================================================
// Internal Supporters CRUD

func (supporter *Supporter) List() ([]Supporter, error) {
    var supporters []Supporter
    cursor, err := dbSUPPORTERS.Find(context.Background(), bson.D{{}})
    if err != nil {
        return supporters, err
    }
    defer cursor.Close(context.Background())
    err = cursor.All(context.Background(), &supporters)
    return supporters, err
}

func (supporter *Supporter) Select(id primitive.ObjectID) error {
    return dbSUPPORTERS.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(supporter)
}

func (supporter *Supporter) Add() error {
    _, err := dbSUPPORTERS.InsertOne(context.Background(), supporter)
    return err
}

func (supporter *Supporter) Update() error {
    _, err := dbSUPPORTERS.ReplaceOne(context.Background(), bson.D{{"_id", supporter.Id}}, supporter)
    return err
}

func (supporter *Supporter) Delete() error {
    _, err := dbSUPPORTERS.DeleteOne(context.Background(), bson.D{{"_id", supporter.Id}})
    return err
}

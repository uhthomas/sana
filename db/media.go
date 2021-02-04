package db

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
	uppdb "upper.io/db"
)

var MediaCollection uppdb.Collection

func init() {
	coll, err := Session.Collection("media")
	if err != nil && err != uppdb.ErrCollectionDoesNotExist {
		panic(err)
	}
	MediaCollection = coll
}

type Media struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	OriginalId  string        `json:"originalId" bson:"originalId"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Duration    time.Duration `json:"duration" bson:"duration"`
	Uploaded    time.Time     `json:"uploaded" bson:"uploaded"`
	Provider    string        `json:"provider" bson:"provider"`
	ProviderId  string        `json:"providerId" bson:"providerId"`
	Source      string        `json:"source" bson:"source"`
	Etag        string        `json:"etag" bson:"etag"`
}

func NewMedia() *Media {
	return &Media{
		Id: bson.NewObjectId(),
	}
}

func GetMedia(query interface{}) (m Media, err error) {
	err = MediaCollection.Find(query).One(&m)
	return
}

func GetMultiMedia(max int, query interface{}) (m []Media, err error) {
	q := MediaCollection.Find(query).Sort("uploaded")
	if max < 0 {
		err = q.All(&m)
	} else {
		err = q.Limit(uint(max)).All(&m)
	}
	return
}

func (m Media) Save() (err error) {
	_, err = MediaCollection.Append(m)
	if err != nil {
		return
	}
	_, err = ElasticClient.Index().
		Index("sana").
		Type("media").
		Id(fmt.Sprintf("%x", string(m.Id))).
		BodyJson(m).
		Do()
	if err != nil {
		return
	}
	_, err = ElasticClient.Flush().
		Index("sana").
		Do()
	return
}

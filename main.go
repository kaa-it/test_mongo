package main

import (
	"context"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://astro:Q123456q@light.tvhelp.ru:27017/admin?authSource=astro"
const dbName = "astro"

type City struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string
}

type Segment struct {
	Id     primitive.ObjectID `bson:"_id"`
	CityId primitive.ObjectID `bson:"city_id"`
	Name   string
}

type Controller struct {
	Id     primitive.ObjectID `bson:"_id"`
	Master bool
}

type Lamp struct {
	Id              primitive.ObjectID `bson:"_id"`
	SegmentId       *primitive.ObjectID
	City            *string
	Segment         *string
	Mac             string
	Source          string
	Dir             int
	Level           int
	Nid             *int
	Group           int
	Smac            string
	Rssi            int
	Devt            int
	Devm            int
	Eblk            *int
	Cycles          *string
	Runh            *int
	Nvsc            *int
	Lpwm            *int
	Cpwm            *int
	Mrssi           int
	Rfch            *int
	Rfpwr           *int
	Pwm             *int
	Pwmct           *int
	Pow             *int
	Lux             *int
	Temp            *int
	Energy          *int
	Rng             *int
	Tlevel          int
	Date            *int64
	Lat             *string
	Lng             *string `bson:"lon"`
	Val             *int
	Rise            *string
	Set             *string
	ProfileId       *int `bson: "id"`
	Scdtm           *int
	Rfps            *int
	Twil            *int
	Received        primitive.DateTime
	SimpleProfiles  [3]*SimpleProfile
	ComplexProfiles [3]*ComplexProfile
}

type SimpleProfile struct {
	D1 int
	P1 int
	D2 int
	P2 int
}

type ComplexProfile struct {
	Pwm0  int
	Time1 string
	Pwm1  int
	Time2 string
	Pwm2  int
	Time3 string
	Pwm3  int
	Time4 string
	Pwm4  int
}

func main() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	//cities := client.Database(dbName).Collection("cities")
	//segments := client.Database(dbName).Collection("segments")
	//controllers := client.Database(dbName).Collection("controllers")
	lamps := client.Database(dbName).Collection("lamps")

	/*options := options.Find()
	options.SetLimit(20)
	options.SetSort(bson.M{"received": 1})*/

	cur, err := lamps.Find(context.Background(), bson.D{{}} /*, options*/)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		var result Lamp // bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}

		if result.Lat != nil {
			lat, err := strconv.ParseFloat(*result.Lat, 64)
			if err == nil {
				log.Printf("%f", lat)
			}
		}

		if result.Lng != nil {
			lng, err := strconv.ParseFloat(*result.Lng, 64)
			if err == nil {
				log.Printf("%f", lng)
			}
		}

		log.Printf("%s %s", result.Mac, result.Received.Time().String())

	}

	cur.Close(context.Background())

	client.Disconnect(context.TODO())
}

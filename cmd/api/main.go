package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	db struct {
		uri         string
		mongoClient *mongo.Client
	}
}

type application struct {
	config config
}

func main() {
	r := gin.Default()
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config

	flag.StringVar(&cfg.db.uri, "mongodb-uri", os.Getenv("MONGODB_URI"), "MongoDB URI")

	flag.Parse()

	app := &application{
		config: cfg,
	}

	// Connect to MongoDB
	if err := connect_to_mongodb(&app.config); err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}
	// Router
	Router(r, app)

	r.Run()
}

func connect_to_mongodb(cfg *config) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.db.uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	cfg.db.mongoClient = client
	return nil
}

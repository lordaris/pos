package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogEntry struct {
	Timestamp    time.Time          `bson:"timestamp"`
	Method       string             `bson:"method"`
	Path         string             `bson:"path"`
	Status       int                `bson:"status"`
	Latency      string             `bson:"latency"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Username     string             `bson:"username"`
	RequestBody  string             `bson:"request_body"`
	ResponseBody string             `bson:"response_body"`
	QueryParams  string             `bson:"query_params"`
	ResourceID   string             `bson:"resource_id"`
}

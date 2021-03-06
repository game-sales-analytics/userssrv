package db

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/game-sales-analytics/users-service/internal/db/repository"
)

type DB struct {
	client *mongo.Client
	logger *logrus.Entry
	Repo   repository.Repo
}

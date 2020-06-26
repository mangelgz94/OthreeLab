package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Customer struct {
	Name        string    `json:"name,omitempty" bson:"name,omitempty"`
	LastName    string    `json:"last_name, omitempty" bson:"last_name, omitempty"`
	BirthDate   time.Time `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	Email       string    `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
}

func (customer *Customer) createCustomer(db *mongo.Client) error {
	collection := db.Database("customers_db").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.InsertOne(ctx, customer)
	return err
}

func getCustomers(db *mongo.Client) ([]Customer, error) {
	var customers []Customer
	collection := db.Database("customers_db").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var customer Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

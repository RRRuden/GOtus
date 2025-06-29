package repository

import (
	"context"
	"gotus/internal/model/reservation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoReservationRepository struct {
	collection *mongo.Collection
}

func NewMongoReservationRepository(db *mongo.Database) ReservationRepository {
	return &mongoReservationRepository{collection: db.Collection("reservations")}
}

func (r *mongoReservationRepository) StoreReservation(res *reservation.Reservation) {
	if res.Id == 0 {
		newID, _ := r.getNextID()
		res.SetID(newID)
	}

	r.collection.InsertOne(context.TODO(), res)
}

func (r *mongoReservationRepository) GetReservations() ([]*reservation.Reservation, int) {
	cursor, _ := r.collection.Find(context.TODO(), bson.M{})
	var result []*reservation.Reservation
	for cursor.Next(context.TODO()) {
		var res reservation.Reservation
		cursor.Decode(&res)
		result = append(result, &res)
	}
	return result, len(result)
}

func (r *mongoReservationRepository) UpdateReservationById(id int, updated *reservation.Reservation) bool {
	res, _ := r.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": updated})
	return res.ModifiedCount > 0
}

func (r *mongoReservationRepository) FindReservationById(id int) (*reservation.Reservation, bool) {
	var res reservation.Reservation
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&res)
	return &res, err == nil
}

func (r *mongoReservationRepository) DeleteReservationById(id int) bool {
	res, _ := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	return res.DeletedCount > 0
}

func (r *mongoReservationRepository) HasActiveReservation(bookInstanceID int) bool {
	filter := bson.M{"bookinstanceid": bookInstanceID, "reservationstatusid": bson.M{"$in": []int{1, 2}}} // 1: Booked, 2: Extended
	count, _ := r.collection.CountDocuments(context.TODO(), filter)
	return count > 0
}

func (r *mongoReservationRepository) getNextID() (int, error) {
	var result struct {
		Seq int `bson:"seq"`
	}
	filter := bson.M{"_id": "reservationid"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	err := r.collection.Database().Collection("counters").
		FindOneAndUpdate(context.TODO(), filter, update, opts).
		Decode(&result)

	return result.Seq, err
}

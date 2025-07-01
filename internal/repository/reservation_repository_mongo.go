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

func (r *mongoReservationRepository) StoreReservation(res *reservation.Reservation) (int, error) {
	if res.Id == 0 {
		newID, err := r.getNextID()
		if err != nil {
			return 0, err
		}
		res.SetID(newID)
	}

	_, err := r.collection.InsertOne(context.TODO(), res)
	return res.Id, err
}

func (r *mongoReservationRepository) GetReservations() ([]*reservation.Reservation, int, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var result []*reservation.Reservation
	for cursor.Next(context.TODO()) {
		var res reservation.Reservation
		if err := cursor.Decode(&res); err != nil {
			return nil, 0, err
		}
		result = append(result, &res)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return result, len(result), nil
}

func (r *mongoReservationRepository) UpdateReservationById(id int, updated *reservation.Reservation) (bool, error) {
	res, err := r.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": updated})
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *mongoReservationRepository) FindReservationById(id int) (*reservation.Reservation, bool, error) {
	var res reservation.Reservation
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &res, true, nil
}

func (r *mongoReservationRepository) DeleteReservationById(id int) (bool, error) {
	res, err := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

func (r *mongoReservationRepository) HasActiveReservation(bookInstanceID int) (bool, error) {
	filter := bson.M{"bookinstanceid": bookInstanceID, "reservationstatusid": bson.M{"$in": []int{1, 2}}} // 1: Booked, 2: Extended
	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

package repository

import (
	"context"
	"errors"
	"time"

	domain "taskmanagement/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoNoteRepository struct {
	Collection *mongo.Collection
}

func (r *MongoNoteRepository) GetAllNote(userID primitive.ObjectID) ([]domain.Note, error) {
	filter := bson.M{"user_id": userID}

	cursor, err := r.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notes []domain.Note
	if err := cursor.All(context.Background(), &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *MongoNoteRepository) GetNoteByID(noteID, userID primitive.ObjectID) (domain.Note, error) {
	var note domain.Note
	filter := bson.M{"_id": noteID, "user_id": userID}

	err := r.Collection.FindOne(context.Background(), filter).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Note{}, errors.New("note not found")
		}
		return domain.Note{}, err
	}
	return note, nil
}

func (r *MongoNoteRepository) CreateNote(note domain.Note, userID primitive.ObjectID) (domain.Note, error) {
	note.UserID = userID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(context.Background(), note)
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (r *MongoNoteRepository) UpdateNote(note domain.Note, userID primitive.ObjectID) (domain.Note, error) {
	note.UpdatedAt = time.Now()

	filter := bson.M{"_id": note.ID, "user_id": userID}
	_, err := r.Collection.UpdateOne(context.Background(), filter, bson.M{"$set": note})
	if err != nil {
		return domain.Note{}, err
	}

	return r.GetNoteByID(note.ID, userID)
}

func (r *MongoNoteRepository) DeleteNote(noteID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": noteID, "user_id": userID}
	result, err := r.Collection.DeleteOne(context.Background(), filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New("note not found or unauthorized")
	}
	return nil
}

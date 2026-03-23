package usecase

import (
	domain "taskmanagement/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteUsecase struct {
	Repo domain.NoteRepository
}

func (u *NoteUsecase) GetAllNote(userId primitive.ObjectID) ([]domain.Note, error) {
	return u.Repo.GetAllNote(userId)
}

func (u *NoteUsecase) GetNoteByID(noteId, userId primitive.ObjectID) (domain.Note, error) {
	return u.Repo.GetNoteByID(noteId, userId)
}

func (u *NoteUsecase) CreateNote(note domain.Note, userId primitive.ObjectID) (domain.Note, error) {
	if note.Title == "" {
		note.Title = "Untitled"
	}
	return u.Repo.CreateNote(note, userId)
}

func (u *NoteUsecase) UpdateNote(note domain.Note, userId primitive.ObjectID) (domain.Note, error) {
	return u.Repo.UpdateNote(note, userId)
}

func (u *NoteUsecase) DeleteNote(noteId, userId primitive.ObjectID) error {
	return u.Repo.DeleteNote(noteId, userId)
}

package storage

import (
	"errors"
	"notes-api/models"
	"sync"
)

// necessary variables
var (
	notes = make(map[string]models.Note)
	mu    sync.RWMutex
)

// add a note to the storage
func Create(note models.Note) {
	mu.Lock()
	defer mu.Unlock()
	notes[note.ID] = note
}

// get all notes from storage --> return the JSON formats
func AllNotes() []models.Note {
	mu.RLock()
	defer mu.RUnlock()
	answer := make([]models.Note, 0, len(notes))
	for _, note := range notes {
		answer = append(answer, note)
	}
	return answer
}

func GetNoteByID(id string) (models.Note, error) {
	mu.RLock()
	defer mu.RUnlock()
	note, err := notes[id]
	if !err {
		return models.Note{}, errors.New("not found")
	}
	return note, nil
}

func UpdateNote(note models.Note) error {
	mu.RLock()
	defer mu.RUnlock()
	notes[note.ID] = note
	if _, exists := notes[note.ID]; !exists {
		return errors.New("note not found")
	}
	return nil
}

func DeleteNoteById(id string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := notes[id]; !exists {
		return errors.New("note not found")
	}
	delete(notes, id)
	return nil
}

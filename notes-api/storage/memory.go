package storage

import (
	"errors"
	"notes-api/models"
	"sync"
)

var (
	notes = make(map[string]models.Note)
	mu    sync.RWMutex
)

func Create(note models.Note) {
	mu.Lock()
	defer mu.Unlock()
	notes[note.ID] = note
}

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
	note, exists := notes[id]
	if !exists {
		return models.Note{}, errors.New("not found")
	}
	return note, nil
}

func UpdateNote(note models.Note) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := notes[note.ID]; !exists {
		return errors.New("note not found")
	}
	notes[note.ID] = note
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

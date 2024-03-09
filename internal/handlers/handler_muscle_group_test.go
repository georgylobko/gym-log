package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/mappers"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetMuscleGroups(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	queries := database.New(d)
	apiCfg := ApiConfig{
		DB: queries,
	}

	t.Run("should return muscle groups", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "photo_url"}).
			AddRow(1, "Biceps", "https://placeholder/biceps.png")
		mock.ExpectQuery("SELECT id, name, photo_url FROM muscle_groups").
			WillReturnRows(rows)

		r, _ := http.NewRequest("GET", "/muscle-groups", nil)
		w := httptest.NewRecorder()
		apiCfg.HandlerGetMuscleGroups(w, r, mappers.User{})

		assert.Equal(t, http.StatusOK, w.Code)

		var response []mappers.MuscleGroup
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		expectedMuscleGroups := []mappers.MuscleGroup{
			{ID: 1, Name: "Biceps", PhotoUrl: "https://placeholder/biceps.png"},
		}
		assert.Equal(t, expectedMuscleGroups, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

func TestHandlerCreateMuscleGroup(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	queries := database.New(d)
	apiCfg := ApiConfig{
		DB: queries,
	}

	t.Run("should return 400 error code when user input is invalid", func(t *testing.T) {
		requestBody := `{"name": "Triceps"}`
		r, _ := http.NewRequest("POST", "/create-muscle-group", bytes.NewBufferString(requestBody))
		w := httptest.NewRecorder()

		apiCfg.HandlerCreateMuscleGroup(w, r, mappers.User{})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should create muscle group if user input is valid", func(t *testing.T) {
		requestBody := `{"name": "Triceps", "photo_url": "https://placeholder/triceps.png"}`
		r, _ := http.NewRequest("POST", "/create-muscle-group", bytes.NewBufferString(requestBody))
		w := httptest.NewRecorder()

		mockMuscleGroup := database.MuscleGroup{
			ID:       1,
			Name:     "Triceps",
			PhotoUrl: "https://placeholder/triceps.png",
		}

		mock.ExpectQuery("INSERT INTO muscle_groups").
			WithArgs("Triceps", "https://placeholder/triceps.png").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "photo_url"}).
				AddRow(mockMuscleGroup.ID, mockMuscleGroup.Name, mockMuscleGroup.PhotoUrl))

		apiCfg.HandlerCreateMuscleGroup(w, r, mappers.User{})

		assert.Equal(t, http.StatusOK, w.Code)

		var response mappers.MuscleGroup
		err = json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		expectedMuscleGroup := mappers.DatabaseMuscleGroupToMuscleGroup(mockMuscleGroup)
		assert.Equal(t, expectedMuscleGroup, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

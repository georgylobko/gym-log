package handlers

import (
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

	t.Run("return muscle groups", func(t *testing.T) {
		// Mock database query result
		rows := sqlmock.NewRows([]string{"id", "name", "photo_url"}).
			AddRow(1, "Biceps", "https://placeholder/biceps.png")
		mock.ExpectQuery("SELECT id, name, photo_url FROM muscle_groups").
			WillReturnRows(rows)

		// Create a mock HTTP request
		req, err := http.NewRequest("GET", "/muscle-groups", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP response writer
		w := httptest.NewRecorder()

		// Call the handler function
		apiCfg := ApiConfig{
			DB: queries,
		}
		apiCfg.HandlerGetMuscleGroups(w, req, mappers.User{})

		// Check the HTTP response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse the response body (assuming it's JSON) for further assertions
		var response []mappers.MuscleGroup
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		// Assert the expected values based on the mocked database result
		expectedMuscleGroups := []mappers.MuscleGroup{
			{ID: 1, Name: "Biceps", PhotoUrl: "https://placeholder/biceps.png"},
		}
		assert.Equal(t, expectedMuscleGroups, response)

		// Assert that the expected database query was executed
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

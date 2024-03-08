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
		rows := sqlmock.NewRows([]string{"id", "name", "photo_url"}).
			AddRow(1, "Biceps", "https://placeholder/biceps.png")
		mock.ExpectQuery("SELECT id, name, photo_url FROM muscle_groups").
			WillReturnRows(rows)

		req, err := http.NewRequest("GET", "/muscle-groups", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()

		apiCfg := ApiConfig{
			DB: queries,
		}
		apiCfg.HandlerGetMuscleGroups(w, req, mappers.User{})

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

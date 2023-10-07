package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestTemplateAPI(t *testing.T) {
	t.Parallel()
	t.Run("GetRepos", func(t *testing.T) {
		t.Parallel()
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/templates", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var templates []*models.TemplateRepo
		err = json.NewDecoder(response.Body).Decode(&templates)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, templates) {
			return
		}
		hasLocal := false
		hasCommunity := false
		for _, v := range templates {
			if v.IsLocal {
				hasLocal = true
			} else if v.Name == "community" {
				hasCommunity = true
			}
		}

		assert.True(t, hasLocal, "No local repo")
		assert.True(t, hasCommunity, "No community template repo")
	})

	t.Run("GetCommunityRepo", func(t *testing.T) {
		t.Parallel()
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/templates/1", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var templates []*models.Template
		err = json.NewDecoder(response.Body).Decode(&templates)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, templates) {
			return
		}
	})

	t.Run("GetTemplateFromCommunity", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/templates/1/minecraft", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var template models.Template
		err = json.NewDecoder(response.Body).Decode(&template)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, template) && !assert.NotEmpty(t, template.Name) {
			return
		}
	})

	t.Run("AddTemplateToLocal", func(t *testing.T) {
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}

		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("PUT", "/api/templates/0/minecraft-vanilla", strings.NewReader(TemplateData), session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}

		mo := &models.Template{
			Name: "minecraft-vanilla",
		}
		var count int64
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, int64(1), count) {
			return
		}
	})

	t.Run("DeleteTemplateFromLocal", func(t *testing.T) {
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}

		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		mo := &models.Template{
			Name: "minecraft-vanilla",
		}
		var count int64
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, int64(1), count) {
			return
		}

		response := CallAPI("DELETE", "/api/templates/0/minecraft-vanilla", strings.NewReader(TemplateData), session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}

		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, int64(0), count) {
			return
		}
	})
}
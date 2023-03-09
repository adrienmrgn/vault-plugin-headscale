package headscale

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"

)

func TestNewClient(t *testing.T) {
	c := newClient()
	// Vérifie que les champs ont bien été initialisés avec les valeurs par défaut.
	assert.Equal(t, "", c.ApiURL, "ApiURL doit être vide")
	assert.Equal(t, "", c.ApiKey, "ApiKey doit être vide")
	assert.Equal(t, time.Minute, c.HTTP.Timeout, "Timeout doit être d'une minute")

	// Vérifie que l'objet *Client est bien créé avec le type attendu.
	assert.IsType(t, &Client{}, c, "newClient doit retourner un objet de type *Client")

}
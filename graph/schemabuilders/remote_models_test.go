package schemabuilders

import (
	"testing"
)

func TestCreateRemoteModelSchemaBuilder(t *testing.T) {
	sb := CreateRemoteModelSchemaBuilder("https://tequila-mantle.terra.dev")
}

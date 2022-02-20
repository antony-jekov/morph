package morph

import (
	"github.com/stretchr/testify/require"
	"testing"
)

//region New

func TestServer_New_Nothing(t *testing.T) {
	require.NotNil(t, New())
}

//endregion New

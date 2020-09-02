package handler

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_writeError(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeError(recorder, "some error message", 400)

	result := recorder.Result()
	t.Cleanup(func() { result.Body.Close() })

	require.Equal(t, 400, result.StatusCode)

	body, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err)
	require.Equal(t, `{"error":"some error message"}`, string(body))
}

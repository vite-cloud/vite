package manifest

import (
	"encoding/json"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func TestManifest_MarshalJSON(t *testing.T) {
	data := map[string]any{
		"hello": []string{"world"},
		"foo":   []any{"bar", 4},
	}

	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	got, err := json.Marshal(m)
	assert.NilError(t, err)

	want, _ := json.Marshal(data)

	assert.Equal(t, string(got), string(want))
}

func TestManifest_Add(t *testing.T) {
	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	got, ok := m.Resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "world")

	got, ok = m.Resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "bar")
	assert.Equal(t, got.([]any)[1], 4)
}

func TestManifest_Save(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	err = m.Save()
	assert.NilError(t, err)

	got, err := os.ReadFile(datadir.Dir() + "/" + Store.String() + "/testing.json")
	assert.NilError(t, err)

	want, err := json.Marshal(m)
	assert.NilError(t, err)

	assert.Equal(t, string(got), string(want))
}

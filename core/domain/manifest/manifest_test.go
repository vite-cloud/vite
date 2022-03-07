package manifest

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
)

func TestManifest_MarshalJSON(t *testing.T) {
	data := manifestJSON{
		Version: "testing",
		Resources: map[string]any{
			"hello": []string{"world"},
			"foo":   []any{"bar", 4},
		},
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

	got, ok := m.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "world")

	got, ok = m.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "bar")
	assert.Equal(t, got.([]any)[1], 4)
}

func TestManifest_Save(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	err = m.Save()
	assert.NilError(t, err)

	dir, err := Store.Dir()
	assert.NilError(t, err)

	got, err := os.ReadFile(dir + "/testing.json")
	assert.NilError(t, err)

	want, err := json.Marshal(m)
	assert.NilError(t, err)

	assert.Equal(t, string(got), string(want))
}

func TestManifest_Get(t *testing.T) {
	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	got, err := m.Get("hello")
	assert.NilError(t, err)
	assert.Equal(t, got[0], "world")

	got, err = m.Get("foo")
	assert.NilError(t, err)
	assert.Equal(t, got[0], "bar")
	assert.Equal(t, got[1], 4)
}

func TestList(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	err = m.Save()
	assert.NilError(t, err)

	got, err := List()
	assert.NilError(t, err)

	assert.Equal(t, len(got), 1)
}

func TestManifest_UnmarshalJSON(t *testing.T) {
	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	marshaled, err := json.Marshal(m)
	assert.NilError(t, err)

	var unmarshaled Manifest
	err = json.Unmarshal(marshaled, &unmarshaled)
	assert.NilError(t, err)

	got, ok := unmarshaled.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "world")

	got, ok = unmarshaled.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "bar")
	assert.Equal(t, got.([]any)[1], 4.0)
}

func TestGet(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	err = m.Save()
	assert.NilError(t, err)

	found, err := Get("testing")
	assert.NilError(t, err)

	assert.Equal(t, found.Version, "testing")

	got, ok := found.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "world")

	got, ok = found.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "bar")
	assert.Equal(t, got.([]any)[1], 4.0)
}

func TestDelete(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "world")
	m.Add("foo", "bar")
	m.Add("foo", 4)

	err = m.Save()
	assert.NilError(t, err)

	found, err := Get("testing")
	assert.NilError(t, err)

	assert.Equal(t, found.Version, "testing")

	got, ok := found.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "world")

	got, ok = found.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]any)[0], "bar")
	assert.Equal(t, got.([]any)[1], 4.0)

	err = Delete("testing")
	assert.NilError(t, err)

	_, err = Get("testing")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDelete2(t *testing.T) {
	datadir.SetHomeDir("/nop")

	err := Delete("testing")
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestDelete3(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	err = Delete("does_not_exist")
	assert.NilError(t, err)
}

func TestDelete4(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	err = Delete("\000x")
	assert.ErrorContains(t, err, "invalid argument")
}

func TestList2(t *testing.T) {
	datadir.SetHomeDir("/nop")

	_, err := List()
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestList3(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/whatever", 0644)
	assert.NilError(t, err)

	_, err = List()
	assert.ErrorContains(t, err, "manifest store is corrupted: whatever is a directory")
}

func TestList4(t *testing.T) {
	// ensure that an invalid json file returns an error
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.WriteFile(dir+"/whatever.json", []byte("not JSON"), 0644)
	assert.NilError(t, err)

	_, err = List()
	assert.ErrorContains(t, err, "invalid character")
}

func TestManifest_UnmarshalJSON2(t *testing.T) {
	m := &Manifest{Version: "testing"}

	err := m.UnmarshalJSON([]byte("not JSON"))
	assert.ErrorContains(t, err, "invalid character")
}

func TestManifest_Get2(t *testing.T) {
	m := &Manifest{Version: "testing"}

	_, err := m.Get("does not exist")
	assert.ErrorContains(t, err, "no resources found matching given key")

}

func TestManifest_Save2(t *testing.T) {
	home, err := os.MkdirTemp("", "manifest-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	m := &Manifest{Version: "testing"}

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/testing.json", 0644)
	assert.NilError(t, err)

	err = m.Save()
	assert.ErrorContains(t, err, "is a directory")
}

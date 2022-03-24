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
		Resources: map[string][]LabeledValue{
			"hello": {{"label", "world"}},
			"foo":   {{"label1", "bar"}, {"label2", 4}},
		},
	}

	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	got, err := json.Marshal(m)
	assert.NilError(t, err)

	want, _ := json.Marshal(data)

	assert.Equal(t, string(got), string(want))
}

func TestManifest_Add(t *testing.T) {
	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	got, ok := m.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "world")

	got, ok = m.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label1")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "bar")
	assert.Equal(t, got.([]LabeledValue)[1].Label, "label2")
	assert.Equal(t, got.([]LabeledValue)[1].Value, 4)
}

func TestManifest_Save(t *testing.T) {
	datadir.UseTestHome(t)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	err := m.Save()
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

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	got, err := m.Get("hello")
	assert.NilError(t, err)
	assert.Equal(t, got[0].Label, "label")
	assert.Equal(t, got[0].Value, "world")

	got, err = m.Get("foo")
	assert.NilError(t, err)
	assert.Equal(t, got[0].Label, "label1")
	assert.Equal(t, got[0].Value, "bar")
	assert.Equal(t, got[1].Label, "label2")
	assert.Equal(t, got[1].Value, 4)
}

func TestList(t *testing.T) {
	datadir.UseTestHome(t)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	err := m.Save()
	assert.NilError(t, err)

	got, err := List()
	assert.NilError(t, err)

	assert.Equal(t, len(got), 1)

	key, err := got[0].Get("hello")
	assert.NilError(t, err)

	assert.Equal(t, key[0].Label, "label")
	assert.Equal(t, key[0].Value, "world")

	key, err = got[0].Get("foo")
	assert.NilError(t, err)

	assert.Equal(t, key[0].Label, "label1")
	assert.Equal(t, key[0].Value, "bar")
}

func TestManifest_UnmarshalJSON(t *testing.T) {
	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	marshaled, err := json.Marshal(m)
	assert.NilError(t, err)

	var unmarshaled Manifest
	err = json.Unmarshal(marshaled, &unmarshaled)
	assert.NilError(t, err)

	got, ok := unmarshaled.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "world")

	got, ok = unmarshaled.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label1")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "bar")
	assert.Equal(t, got.([]LabeledValue)[1].Label, "label2")
	assert.Equal(t, got.([]LabeledValue)[1].Value, 4.0)
}

func TestGet(t *testing.T) {
	datadir.UseTestHome(t)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	err := m.Save()
	assert.NilError(t, err)

	found, err := Get("testing")
	assert.NilError(t, err)

	assert.Equal(t, found.Version, "testing")

	got, ok := found.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "world")

	got, ok = found.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label1")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "bar")
	assert.Equal(t, got.([]LabeledValue)[1].Label, "label2")
	assert.Equal(t, got.([]LabeledValue)[1].Value, 4.0)
}

func TestDelete(t *testing.T) {
	datadir.UseTestHome(t)

	m := &Manifest{Version: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	err := m.Save()
	assert.NilError(t, err)

	found, err := Get("testing")
	assert.NilError(t, err)

	assert.Equal(t, found.Version, "testing")

	got, ok := found.resources.Load("hello")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "world")

	got, ok = found.resources.Load("foo")
	assert.Assert(t, ok)
	assert.Equal(t, got.([]LabeledValue)[0].Label, "label1")
	assert.Equal(t, got.([]LabeledValue)[0].Value, "bar")
	assert.Equal(t, got.([]LabeledValue)[1].Label, "label2")
	assert.Equal(t, got.([]LabeledValue)[1].Value, 4.0)

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
	datadir.UseTestHome(t)

	err := Delete("does_not_exist")
	assert.NilError(t, err)
}

func TestDelete4(t *testing.T) {
	datadir.UseTestHome(t)

	err := Delete("\000x")
	assert.ErrorContains(t, err, "invalid argument")
}

func TestList2(t *testing.T) {
	datadir.SetHomeDir("/nop")

	_, err := List()
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestList3(t *testing.T) {
	datadir.UseTestHome(t)

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/whatever", 0644)
	assert.NilError(t, err)

	_, err = List()
	assert.ErrorContains(t, err, "manifest store is corrupted: whatever is a directory")
}

func TestList4(t *testing.T) {
	datadir.UseTestHome(t)

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
	datadir.UseTestHome(t)

	m := &Manifest{Version: "testing"}

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/testing.json", 0644)
	assert.NilError(t, err)

	err = m.Save()
	assert.ErrorContains(t, err, "is a directory")
}

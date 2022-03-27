package deployment

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestManifest_MarshalJSON(t *testing.T) {
	data := deploymentJSON{
		ID: "testing",
		Resources: map[string][]LabeledValue{
			"hello": {{"label", "world"}},
			"foo":   {{"label1", "bar"}, {"label2", 4}},
		},
	}

	m := &Deployment{ID: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	got, err := json.Marshal(m)
	assert.NilError(t, err)

	want, _ := json.Marshal(data)

	assert.Equal(t, string(got), string(want))
}

func TestManifest_Add(t *testing.T) {
	m := &Deployment{ID: "testing"}

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

func TestManifest_Get(t *testing.T) {
	m := &Deployment{ID: "testing"}

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

func TestManifest_UnmarshalJSON(t *testing.T) {
	m := &Deployment{ID: "testing"}

	m.Add("hello", "label", "world")
	m.Add("foo", "label1", "bar")
	m.Add("foo", "label2", 4)

	marshaled, err := json.Marshal(m)
	assert.NilError(t, err)

	var unmarshaled Deployment
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

func TestManifest_UnmarshalJSON2(t *testing.T) {
	m := &Deployment{ID: "testing"}

	err := m.UnmarshalJSON([]byte("not JSON"))
	assert.ErrorContains(t, err, "invalid character")
}

func TestManifest_Get2(t *testing.T) {
	m := &Deployment{ID: "testing"}

	_, err := m.Get("does not exist")
	assert.ErrorContains(t, err, "no resources found matching given key")

}

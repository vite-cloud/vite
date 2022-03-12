package deployment

import (
	"github.com/vite-cloud/vite/core/domain/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestFlatten5(t *testing.T) {
	t.Parallel()

	example := &config.Service{
		Name: "example",
	}
	services := map[string]*config.Service{
		"example": example,
	}

	layers, err := Layered(services)
	assert.NilError(t, err)
	assert.Equal(t, len(layers), 1)
	assert.Equal(t, len(layers[0]), 1)
	assert.Equal(t, layers[0][0], example)
}

func TestFlatten4(t *testing.T) {
	t.Parallel()

	a := &config.Service{
		Name: "a",
	}
	b := &config.Service{
		Name: "b",
	}
	a.Requires = []*config.Service{b}
	b.Requires = []*config.Service{a}

	services := map[string]*config.Service{
		"a": a,
		"b": b,
	}

	_, err := Layered(services)
	assert.ErrorContains(t, err, "circular dependency detected:")
}

func TestFlatten3(t *testing.T) {
	t.Parallel()

	c := service("c")
	b := service("b", c)
	a := service("a", b)
	// circular dependency
	c.Requires = []*config.Service{a}

	services := map[string]*config.Service{
		"a": a,
		"b": b,
		"c": c,
	}

	_, err := Layered(services)
	assert.ErrorContains(t, err, "circular dependency detected:")
}

func TestFlatten2(t *testing.T) {
	t.Parallel()

	fs := service("fs")
	minio := service("minio", fs)
	mysql := service("mysql", fs)
	redis := service("redis")
	elastic := service("elastic", minio)
	laravel := service("laravel", redis, elastic, mysql)

	services := map[string]*config.Service{
		"fs":      fs,
		"minio":   minio,
		"mysql":   mysql,
		"redis":   redis,
		"elastic": elastic,
		"laravel": laravel,
	}

	layers, err := Layered(services)
	assert.NilError(t, err)

	assert.Equal(t, len(layers), 4)

	// layer 0: fs
	assert.Equal(t, len(layers[0]), 1)
	assert.Equal(t, layers[0][0], fs)

	// layer 1: minio
	assert.Equal(t, len(layers[1]), 1)
	assert.Equal(t, layers[1][0], minio)

	// layer 2: elastic, mysql, redis
	assert.Equal(t, len(layers[2]), 3)
	assert.Equal(t, layers[2][0], elastic)
	assert.Equal(t, layers[2][1], mysql)
	assert.Equal(t, layers[2][2], redis)

	// layer 3: laravel
	assert.Equal(t, len(layers[3]), 1)
	assert.Equal(t, layers[3][0], laravel)
}

func service(name string, requires ...*config.Service) *config.Service {
	s := &config.Service{
		Name: name,
	}
	for _, r := range requires {
		s.Requires = append(s.Requires, r)
	}
	return s
}

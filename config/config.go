// Package config based on method loading github.com/gobuffalo/pop of configs.
// But here we can save config in ConnectionDetails store and free to use like we need.
package config

import (
	"bytes"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

// ConfigName is the name of the YAML databases config file
var ConfigName = "database.yml"

// ConnectionDetails all required information for connection
var ConnectionDetails = map[string]*pop.ConnectionDetails{}

// LoadConfigFile loads a POP config file from the configured lookup paths
func LoadConfigFile() error {
	path := "./" + ConfigName
	log.Debug().Msgf("Loading config file from %s", path)

	f, err := os.Open("./" + ConfigName)
	if err != nil {
		return errors.WithStack(err)
	}
	return LoadFrom(f)
}

// LoadFrom reads a configuration from the reader and sets up the connections
func LoadFrom(r io.Reader) error {
	envy.Load()
	tmpl := template.New("test")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(s1, s2 string) string {
			return envy.Get(s1, s2)
		},
		"env": func(s1 string) string {
			return envy.Get(s1, "")
		},
	})
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.WithStack(err)
	}
	t, err := tmpl.Parse(string(b))
	if err != nil {
		return errors.Wrap(err, "couldn't parse config template")
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't execute config template")
	}

	err = yaml.Unmarshal(bb.Bytes(), &ConnectionDetails)
	if err != nil {
		return errors.Wrap(err, "couldn't unmarshal config to yaml")
	}

	return nil
}

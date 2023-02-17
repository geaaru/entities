/*
Copyright © 2022 Funtoo Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package entities

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	UserKind    = "user"
	ShadowKind  = "shadow"
	GroupKind   = "group"
	GShadowKind = "gshadow"
)

type EntitiesParser interface {
	ReadEntity(entity string) (Entity, error)
}

type Signature struct {
	Kind string `yaml:"kind"`
}

type Parser struct{}

func (p Parser) ReadEntityFromBytes(yamlFile []byte) (Entity, error) {

	var signature Signature
	err := yaml.Unmarshal(yamlFile, &signature)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while parsing entity file")
	}

	switch signature.Kind {
	case UserKind:
		var user UserPasswd

		err = yaml.Unmarshal(yamlFile, &user)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return user, nil
	case ShadowKind:
		var shad Shadow

		err = yaml.Unmarshal(yamlFile, &shad)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return shad, nil
	case GroupKind:
		var group Group

		err = yaml.Unmarshal(yamlFile, &group)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return group, nil

	case GShadowKind:
		var group GShadow

		err = yaml.Unmarshal(yamlFile, &group)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return group, nil
	}

	return nil, errors.New("Unsupported format")
}
func (p Parser) ReadEntity(entity string) (Entity, error) {
	yamlFile, err := ioutil.ReadFile(entity)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while reading entity file")
	}
	return p.ReadEntityFromBytes(yamlFile)

}

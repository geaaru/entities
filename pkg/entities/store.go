/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package entities

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type EntitiesStore struct {
	Users    map[string]UserPasswd
	Groups   map[string]Group
	Shadows  map[string]Shadow
	GShadows map[string]GShadow
}

func NewEntitiesStore() *EntitiesStore {
	return &EntitiesStore{
		Users:    make(map[string]UserPasswd, 0),
		Groups:   make(map[string]Group, 0),
		Shadows:  make(map[string]Shadow, 0),
		GShadows: make(map[string]GShadow, 0),
	}
}

func (s *EntitiesStore) Load(dir string) error {
	var regexConfs = regexp.MustCompile(`.yml$|.yaml$`)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	p := &Parser{}

	for _, file := range files {
		if file.IsDir() {
			err = s.Load(filepath.Join(dir, file.Name()))
			// Ignoring errors. Maybe print a warning?
			continue
		}

		if !regexConfs.MatchString(file.Name()) {
			continue
		}

		entity, err := p.ReadEntity(filepath.Join(dir, file.Name()))
		if err == nil {
			s.AddEntity(entity)
		}
	}

	return nil
}

func (s *EntitiesStore) AddEntity(e Entity) error {
	var err error
	switch e.GetKind() {
	case UserKind:
		err = s.AddUser((e.(UserPasswd)))
	case GroupKind:
		err = s.AddGroup((e.(Group)))
	case ShadowKind:
		err = s.AddShadow((e.(Shadow)))
	case GShadowKind:
		err = s.AddGShadow((e.(GShadow)))
	default:
		err = errors.New("Invalid entity")
	}

	return err
}

func (s *EntitiesStore) AddUser(u UserPasswd) error {
	if u.Username == "" {
		return errors.New("Invalid username field")
	}
	if e, ok := s.Users[u.Username]; ok {
		newEntity, err := e.Merge(u)
		if err != nil {
			return err
		}
		s.Users[u.Username] = newEntity.(UserPasswd)
	} else {
		s.Users[u.Username] = u
	}

	return nil
}

func (s *EntitiesStore) AddGroup(g Group) error {
	if g.Name == "" {
		return errors.New("Invalid group name field")
	}

	if e, ok := s.Groups[g.Name]; ok {
		newEntity, err := e.Merge(g)
		if err != nil {
			return err
		}
		s.Groups[g.Name] = newEntity.(Group)

	} else {
		s.Groups[g.Name] = g
	}

	return nil
}

func (s *EntitiesStore) AddShadow(e Shadow) error {
	if e.Username == "" {
		return errors.New("Invalid username field")
	}

	if ne, ok := s.Shadows[e.Username]; ok {
		newEntity, err := ne.Merge(e)
		if err != nil {
			return err
		}
		s.Shadows[e.Username] = newEntity.(Shadow)
	} else {
		s.Shadows[e.Username] = e
	}

	return nil
}

func (s *EntitiesStore) AddGShadow(e GShadow) error {
	if e.Name == "" {
		return errors.New("Invalid name field")
	}

	if ne, ok := s.GShadows[e.Name]; ok {
		newEntity, err := ne.Merge(e)
		if err != nil {
			return err
		}
		s.GShadows[e.Name] = newEntity.(GShadow)
	} else {
		s.GShadows[e.Name] = e
	}

	return nil
}

func (s *EntitiesStore) GetShadow(name string) (Shadow, bool) {
	if e, ok := s.Shadows[name]; ok {
		return e, true
	} else {
		return Shadow{}, false
	}
}

func (s *EntitiesStore) GetGShadow(name string) (GShadow, bool) {
	if e, ok := s.GShadows[name]; ok {
		return e, true
	} else {
		return GShadow{}, false
	}
}

func (s *EntitiesStore) GetUser(name string) (UserPasswd, bool) {
	if e, ok := s.Users[name]; ok {
		return e, true
	} else {
		return UserPasswd{}, false
	}
}

func (s *EntitiesStore) GetGroup(name string) (Group, bool) {
	if e, ok := s.Groups[name]; ok {
		return e, true
	} else {
		return Group{}, false
	}
}

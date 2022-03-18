/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/

package entities

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tredoe/osutil/user/crypt/sha512_crypt"

	permbits "github.com/phayes/permbits"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ParseShadow opens the file and parses it into a map from usernames to Entries
func ParseShadow(path string) (map[string]Shadow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ParseReader(file)
}

// ParseReader consumes the contents of r and parses it into a map from
// usernames to Entries
func ParseReader(r io.Reader) (map[string]Shadow, error) {
	lines := bufio.NewReader(r)
	entries := make(map[string]Shadow)
	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			break
		}
		name, entry, err := parseLine(string(copyBytes(line)))
		if err != nil {
			return nil, err
		}
		entries[name] = entry
	}
	return entries, nil
}

func parseLine(line string) (string, Shadow, error) {
	fs := strings.Split(line, ":")
	if len(fs) != 9 {
		return "", Shadow{}, errors.New("Unexpected number of fields in /etc/shadow: found " + strconv.Itoa(len(fs)))
	}

	return fs[0], Shadow{fs[0], fs[1], fs[2], fs[3], fs[4], fs[5], fs[6], fs[7], fs[8]}, nil
}

func copyBytes(x []byte) []byte {
	y := make([]byte, len(x))
	copy(y, x)
	return y
}

type Shadow struct {
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
	LastChanged    string `yaml:"last_changed" json:"last_changed"`
	MinimumChanged string `yaml:"minimum_changed" json:"minimum_changed"`
	MaximumChanged string `yaml:"maximum_changed" json:"maximum_changed"`
	Warn           string `yaml:"warn" json:"warn"`
	Inactive       string `yaml:"inactive" json:"inactive"`
	Expire         string `yaml:"expire" json:"expire"`
	Reserved       string `yaml:"reserved" json:"reserved"`
}

func (u Shadow) GetKind() string { return ShadowKind }

func (u Shadow) String() string {
	return strings.Join([]string{u.Username,
		u.Password,
		u.LastChanged,
		u.MinimumChanged,
		u.MaximumChanged,
		u.Warn,
		u.Inactive,
		u.Expire,
		u.Reserved,
	}, ":")
}

func ShadowDefault(s string) string {
	if s == "" {
		s = os.Getenv(ENTITY_ENV_DEF_SHADOW)
		if s == "" {
			s = "/etc/shadow"
		}
	}
	return s
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func encryptPassword(userPassword string) (string, error) {
	salt := []byte(fmt.Sprintf("$6$%s", randStringBytes(8)))
	c := sha512_crypt.New()
	hash, err := c.Generate([]byte(userPassword), salt)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (u Shadow) prepare() Shadow {
	if u.LastChanged == "now" {
		// POST: Set in last_changed the current days from 1970
		now := time.Now()
		days := now.Unix() / 24 / 60 / 60
		u.LastChanged = fmt.Sprintf("%d", days)
	}
	/*
	 A password field which starts with an exclamation mark means
	 that the password is locked. The remaining characters on the
	 line represent the password field before the password was
	 locked.

	 Refer to crypt(3) for details on how this string is
	 interpreted.

	 If the password field contains some string that is not a
	 valid result of crypt(3), for instance ! or *, the user will
	 not be able to use a unix password to log in (but the user
	 may log in the system by other means).
	*/
	if !strings.HasPrefix(u.Password, "$") && u.Password != "" &&
		!strings.HasPrefix(u.Password, "!") && u.Password != "*" {
		if pwd, err := encryptPassword(u.Password); err == nil {
			u.Password = pwd
		}
	}
	return u
}

// FIXME: Delete can be shared across all of the supported Entities
func (u Shadow) Delete(s string) error {
	s = ShadowDefault(s)
	input, err := ioutil.ReadFile(s)
	if err != nil {
		return errors.Wrap(err, "Could not read input file")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}
	lines := bytes.Replace(input, []byte(u.String()+"\n"), []byte(""), 1)

	err = ioutil.WriteFile(s, []byte(lines), os.FileMode(permissions))
	if err != nil {
		return errors.Wrap(err, "Could not write")
	}

	return nil
}

// FIXME: Create can be shared across all of the supported Entities
func (u Shadow) Create(s string) error {
	var f *os.File

	s = ShadowDefault(s)

	u = u.prepare()

	_, err := os.Stat(s)
	if err == nil {
		current, err := ParseShadow(s)
		if err != nil {
			return errors.Wrap(err, "Failed parsing passwd")
		}
		if _, ok := current[u.Username]; ok {
			return errors.New("Entity already present")
		}
		permissions, err := permbits.Stat(s)
		if err != nil {
			return errors.Wrap(err, "Failed getting permissions")
		}
		f, err = os.OpenFile(s, os.O_APPEND|os.O_WRONLY, os.FileMode(permissions))
		if err != nil {
			return errors.Wrap(err, "Could not read")
		}
	} else if os.IsNotExist(err) {
		f, err = os.OpenFile(s, os.O_RDWR|os.O_CREATE, 0400)
		if err != nil {
			return errors.Wrap(err, "Could not create the file")
		}
	} else {
		return errors.Wrap(err, "Error on stat file")
	}

	defer f.Close()

	if _, err = f.WriteString(u.String() + "\n"); err != nil {
		return errors.Wrap(err, "Could not write")
	}
	return nil
}

func (u Shadow) Apply(s string, safe bool) error {
	s = ShadowDefault(s)

	u = u.prepare()

	_, err := os.Stat(s)
	if err == nil {
		current, err := ParseShadow(s)
		if err != nil {
			return errors.Wrap(err, "Failed parsing passwd")
		}
		permissions, err := permbits.Stat(s)
		if err != nil {
			return errors.Wrap(err, "Failed getting permissions")
		}

		if _, ok := current[u.Username]; ok {
			input, err := ioutil.ReadFile(s)
			if err != nil {
				return errors.Wrap(err, "Could not read input file")
			}

			lines := strings.Split(string(input), "\n")

			for i, line := range lines {
				if entityIdentifier(line) == u.Username && !safe {
					lines[i] = u.String()
				}
			}
			output := strings.Join(lines, "\n")
			err = ioutil.WriteFile(s, []byte(output), os.FileMode(permissions))
			if err != nil {
				return errors.Wrap(err, "Could not write")
			}

		} else {
			// Add it
			return u.Create(s)
		}
	} else if os.IsNotExist(err) {
		return u.Create(s)
	} else {
		return errors.Wrap(err, "Could not stat file")
	}

	return nil
}

func (s Shadow) Merge(e Entity) (Entity, error) {
	if e.GetKind() != ShadowKind {
		return s, errors.New("merge possible only for entities of the same kind")
	}

	toMerge := e.(Shadow)

	if toMerge.MinimumChanged != "" && toMerge.MinimumChanged != "0" &&
		(s.MinimumChanged == "0" || s.MinimumChanged == "") {
		s.MinimumChanged = toMerge.MinimumChanged
	}

	if toMerge.MaximumChanged != "" && toMerge.MaximumChanged != "0" {
		s.MaximumChanged = toMerge.MaximumChanged
	}

	if toMerge.Warn != "" {
		s.Warn = toMerge.Warn
	}

	if toMerge.Inactive != "" {
		s.Inactive = toMerge.Inactive
	}

	// NOTE: i avoid to change current password.
	return s, nil
}

func (s Shadow) ToMap() map[interface{}]interface{} {
	ans := make(map[interface{}]interface{}, 0)
	d, _ := yaml.Marshal(&s)
	yaml.Unmarshal(d, &ans)
	ans["kind"] = s.GetKind()
	return ans
}

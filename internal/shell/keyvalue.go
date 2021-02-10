package shell

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

type KeyValueStore struct {
	EnvVars map[string]string `json:"env"`
	Aliases map[string]string `json:"alias"`

	path string
	sync.Mutex
}

func newKeyValueStore(path string) *KeyValueStore {
	return &KeyValueStore{
		make(map[string]string),
		make(map[string]string),
		path,
		sync.Mutex{},
	}
}

func LoadKeyValueStore(path string) (*KeyValueStore, error) {
	var err error

	if path == "" {
		path, err = utils.GetKeyValueStorePath()
		if err != nil {
			return nil, err
		}
	}

	if !utils.FileExists(path) {
		tui.Debug("Key value store doesn't exist. Creating in %s", path)
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}

		kv := newKeyValueStore(path)
		err = kv.write()
		if err != nil {
			return nil, err
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}

		return newKeyValueStore(path), nil
	} else {
		tui.Debug("Loading Key value store from %s", path)
		kv := newKeyValueStore(path)
		if err := kv.load(); err != nil {
			return nil, err
		}

		return kv, nil
	}
}

func (kv *KeyValueStore) load() error {
	bytes, err := ioutil.ReadFile(kv.path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, kv); err != nil {
		return err
	}

	return nil
}

func (kv *KeyValueStore) write() error {
	bytes, err := json.MarshalIndent(kv, "", "  ")
	if err != nil {
		return err
	}

	if err = utils.AtomicWrite(bytes, kv.path); err != nil {
		return err
	}

	return nil
}

func (kv *KeyValueStore) addToNamespace(ns, key, value string) error {
	var m map[string]string

	switch ns {
	case "env":
		m = kv.EnvVars
	case "alias":
		m = kv.Aliases
	default:
		panic(fmt.Sprintf("No namespace %s in key value store", ns))
	}

	kv.Lock()
	defer kv.Unlock()

	if err := kv.load(); err != nil {
		return err
	}

	if existingVal, ok := m[key]; ok {
		if value == existingVal {
			return nil
		}
		return fmt.Errorf("Key %s already exists with different value", key)
	}

	m[key] = value
	if err := kv.write(); err != nil {
		return err
	}

	return nil
}

func (kv *KeyValueStore) AddEnv(key, value string) error {
	return kv.addToNamespace("env", key, value)
}

func (kv *KeyValueStore) AddAlias(alias, command string) error {
	return kv.addToNamespace("alias", alias, command)
}

func (kv *KeyValueStore) ListAliases() (map[string]string, error) {
	kv.Lock()
	defer kv.Unlock()

	if err := kv.load(); err != nil {
		return nil, err
	}

	return kv.Aliases, nil
}

func (kv *KeyValueStore) ListEnvVars() (map[string]string, error) {
	kv.Lock()
	defer kv.Unlock()

	if err := kv.load(); err != nil {
		return nil, err
	}

	return kv.EnvVars, nil
}

func (kv *KeyValueStore) Path() string {
	return kv.path
}

type KeyValueExist struct {
	Key string
}

func (k *KeyValueExist) Error() string {
	return fmt.Sprintf("Key '%s' already exist with the requested value", k.Key)
}

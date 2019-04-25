/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


	Copyright 2017 - 2019 OneLedger

		Encapsulate the underlying storage from our app. Currently using:
		Tendermint's memdb (just an in-memory Merkle Tree)
		Tendermint's persistent kvstore (with Merkle Trees & Proofs)
			- Can only be opened by one process...

	Only one connection can occur to LevelDB at a time...

*/

package storage

import (
	"errors"
	"github.com/Oneledger/protocol/data"
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type DatabaseKey = []byte // Database key
type Message = []byte     // TODO: Maybe replaced by something better named?


var k Storage = KeyValue{}


var ErrNilData = errors.New("data is nil")

// Wrap the underlying usage
type KeyValue struct {
	Type StorageType

	Name string
	File string

	memory   *db.MemDB
	tree     *iavl.MutableTree
	database db.DB

	version int64
}

type KeyValueSession struct {
	store *KeyValue
}

// TODO: Should be moved to some common/shared/utils directory
// Test to see if this exists already
func fileExists(name string, dir string) bool {
	dbPath := filepath.Join(dir, name+".db")
	info, err := os.Stat(dbPath)
	if err != nil {
		return false
	}
	_ = info
	return true
}

// Convert Data headed for persistence
func convertData(data interface{}) ([]byte, error) {
	buffer, err := pSzlr.Serialize(data)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// Unconvert Data from persistence
func unconvertData(data []byte) (interface{}, error) {
	if data == nil || string(data) == "" {
		return nil, ErrNilData
	}

	var result interface{}
	err := pSzlr.Deserialize(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NewKeyValue initializes a new application
func NewKeyValue(name, dbDir, configDB string, newType StorageType) *KeyValue {
	switch newType {

	case MEMORY:
		// TODO: No Merkle tree?
		return &KeyValue{
			Type:   newType,
			Name:   name,
			memory: db.NewMemDB(),
		}

	case PERSISTENT:
		storage, err := getDatabase(name, dbDir, configDB)
		if err != nil {
			log.Error("Database create failed", "err", err)
			panic("Can't create a database " + dbDir + "/" + name)
		}

		tree := iavl.NewMutableTree(storage, 100)

		// Note: the tree is empty, until at least one version is loaded
		tree.LoadVersion(0)

		return &KeyValue{
			Type:     newType,
			Name:     name,
			File:     name,
			tree:     tree,
			database: storage,
			version:  tree.Version(),
		}
	default:
		panic("Unknown Type")

	}
	return nil
}

// Begin a new writable session
func (store KeyValue) Begin() StorageSession {
	return NewKeyValueSession(&store)
}

// Dump out debugging information from the KeyValue datastore
func (store KeyValue) Dump() {
	// TODO: Dump out debugging information here
	texts := store.database.Stats()
	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

	iter := store.database.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		hash := iter.Key()
		node := iter.Value()
		log.Debug("Row", "hash", hash, "node", node)
	}
}

// Print out the error details
func (store KeyValue) Errors() string {
	return ""
}

// Close the database
func (store KeyValue) Close() {
	switch store.Type {

	case MEMORY:
		store.memory = nil

	case PERSISTENT:
		store.tree = nil
		store.database.Close()
		store.database = nil

	default:
		panic("Unknown Type")
	}
}

// Close and reopen the datastore
func (store KeyValue) Reopen() {
}

// FindAll of the keys in the database
func (store KeyValue) FindAll() []data.StoreKey {
	return store.list()
}

// Test to see if a key exists
func (store KeyValue) Exists(key data.StoreKey) (bool, error) {
	return store.tree.Has(key), nil

}

// Get a key from the database
func (store KeyValue) Get(key data.StoreKey) ([]byte, error) {

	version :=  store.tree.Version()
	index, value := store.tree.GetVersioned(key, version)
	if index == -1 {
		return nil, ErrNotFound
	}
	return value, nil

}


// List all of the keys
func (store KeyValue) list() (keys []data.StoreKey) {
	switch store.Type {

	case PERSISTENT:
		//store.tree.
		size := store.tree.Size()
		results := make([]data.StoreKey, size, size)
		for i := int64(0); i < store.tree.Size(); i++ {
			key, _ := store.tree.GetByIndex(i)
			results[i] = data.StoreKey(key)
		}
		return results

	default:
		panic("Invalid Op")
	}
}

// Empty out all rows from the database
func (store KeyValue) empty() {
	switch store.Type {
	case MEMORY:
	case PERSISTENT:
	default:
		panic("Unknown Type")
	}
}


/*
		KeyValueSession begins here
 */
// Create a new session
func NewKeyValueSession(store *KeyValue) StorageSession {
	return &KeyValueSession{store: store}
}

// Find all of the keys in the datastore
func (session KeyValueSession) FindAll() []data.StoreKey {
	return session.store.list()
}

// Store inserts or updates a value under a key
func (session KeyValueSession) Set(key data.StoreKey, dat []byte) error {
	ok := session.store.tree.Set(key, dat)
	if !ok {
		return ErrSetFailed
	}
	return nil
}

// Test to see if a key exists
func (session KeyValueSession) Exists(key data.StoreKey) (bool, error) {
	version := session.store.tree.Version()
	index, _ := session.store.tree.GetVersioned(key, version)
	if index == -1 {
		return false, nil
	}
	return true, nil
}

// Load return the stored value
func (session KeyValueSession) Get(key data.StoreKey) ([]byte, error) {
	version := session.store.tree.Version()
	index, value := session.store.tree.GetVersioned(key, version)
	if index == -1 {
		return nil, ErrNotFound
	}
	return value, nil
}

// Delete a key from the datastore
func (session KeyValueSession) Delete(key data.StoreKey) (bool, error) {
	_, deleted := session.store.tree.Remove(key)
	return deleted, nil
}

// List out the errors
func (session KeyValueSession) Errors() string {
	return ""
}

// Commit the changes to persistence
func (session KeyValueSession) Commit() bool {
	_, version, err := session.store.tree.SaveVersion()
	if err != nil {
		log.Fatal("Database Error", "err", err)
	}
	session.store.version = version

	return true
}

// Rollback any changes since the last commit
func (session KeyValueSession) Rollback() bool {
	return false
}

// Dump out the contents of the database
func (session KeyValueSession) Dump() {
	texts := session.store.database.Stats()
	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

	iter := session.store.database.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		hash := iter.Key()
		node := iter.Value()
		log.Debug("Row", "hash", hash, "node", node)
	}
}


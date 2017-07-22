package helpers

import (
	"path"

	"github.com/boltdb/bolt"
)

// Init: create an connexion to the database
var db *bolt.DB

var transports = []byte("Transports")

func init() {
	var err error
	// Open db connection
	db, err = bolt.Open("./Paris.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	// Create Transports bucket if it does not exists
	err = db.Update(func(tx *bolt.Tx) error {
		_, suberr := tx.CreateBucketIfNotExists([]byte("Transports"))
		return suberr
	})
	if err != nil {
		panic(err)
	}
}

// StoreTransports - Store the Transports for a given agency ID
// @param agencyID, the ID of the agency
// @param transports, the Transports to store
// @return an error if one occures
func StoreTransports(agencyID string, buffer []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.Bucket(transports).Put([]byte(agencyID), buffer)
		return nil
	})
}

// GetTransports - get the Transports for a given agency ID
// @param agencyID, the ID of the agency
// @return raw bytes, because it would not be easy to create a custom slice
// for every custom Transports struct
func GetTransports(agencyID string) ([]byte, error) {
	var rawTransports []byte
	err := db.View(func(tx *bolt.Tx) error {
		rawTransports = tx.Bucket(transports).Get([]byte(agencyID))
		return nil
	})
	return rawTransports, err
}

// CheckTransportsExists - check that some Transports are allready stored
// @param agencyID, the ID of the agency
// @return boolean true if the Transports are present
func CheckTransportsExists(agencyID string) (bool, error) {
	var rawTransports []byte
	err := db.View(func(tx *bolt.Tx) error {
		rawTransports = tx.Bucket(transports).Get([]byte(agencyID))
		return nil
	})
	return rawTransports != nil, err
}

// TmpDir - The tmp dir will be cleaned before starting the server
// @param agencyID, the ID of the agency
// @return the tmp folder for the given agency ID
func TmpDir(agencyID string) string {
	return path.Join("tmp", agencyID)
}

// MediaDir - The media folder will be server at run time
// @param agencyID, the ID of the agency
// @return the media folder for the given agency ID
func MediaDir(agencyID string) string {
	return path.Join("medias", agencyID)
}

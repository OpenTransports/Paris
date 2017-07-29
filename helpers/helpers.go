package helpers

import (
	"os"
	"path"
)

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
	return "medias"
}

// ServerURL - URL of the server
var ServerURL = os.Getenv("SERVER_URL")

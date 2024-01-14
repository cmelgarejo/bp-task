package models

type IPFSToken map[string]interface{}

// IPFSMetadata struct, the model for a record of the tokens table
type IPFSMetadata struct {
	Cid   string    `json:"cid" db:"cid"`
	Token IPFSToken `json:"token" db:"token"`
}

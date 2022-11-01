package bw2kp

import (
	"testing"
)

const vaultJson = `
{
"encrypted": false,
"folders": [
	{
	"id": "100",
	"name": "Test Folder"
	}
],
"items": [
	{
	"id": "200",
	"organizationId": null,
	"folderId": "100",
	"type": 1,
	"reprompt": 0,
	"name": "Test Item 1",
	"notes": "This is a note",
	"favorite": false,
	"login": {
		"username": "username1",
		"password": "password1",
		"totp": null
	},
	"collectionIds": null
	},
	{
	"id": "201",
	"organizationId": null,
	"folderId": null,
	"type": 1,
	"reprompt": 0,
	"name": "Test Item 2",
	"notes": null,
	"favorite": false,
	"login": {
		"username": "username2",
		"password": "password2",
		"totp": null
	},
	"collectionIds": null
	}
]
}
`

func getBitwardenVault() (bitwardenVault, error) {
	return marshalBitwardenVault([]byte(vaultJson))
}

func TestCreateKeepassDatabase(t *testing.T) {
	vault, err := getBitwardenVault()

	if err != nil {
		t.Fatalf("Failed to get Bitwarden Vault: %v", err)
	}

	if err = createKeepassDatabase(vault, "../artifacts/db.kdbx", "password"); err != nil {
		t.Fatalf("Failed to create Keepass database: %v", err)
	}
}

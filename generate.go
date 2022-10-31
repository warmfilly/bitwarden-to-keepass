package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

func GenerateKeepassDatabase(opts options) error {
	vault, err := exportBitwardenVault(opts.BitwardenSession)

	if err != nil {
		return err
	}
	createKeepassDatabase(vault, opts.DatabasePath, opts.DatabasePassword)

	return nil
}

func exportBitwardenVault(bitwardenSession string) (BitwardenDatabase, error) {
	out, err := exec.Command("/usr/bin/bw", "export", "--format", "json", "--raw", "--session", bitwardenSession).Output()

	if err != nil || len(out) <= 1 {
		return BitwardenDatabase{}, errors.New("failed to export Bitwarden vault")
	}

	var vault BitwardenDatabase

	if err = json.Unmarshal(out, &vault); err != nil {
		return BitwardenDatabase{}, errors.New("invalid vault data")
	}

	return vault, nil
}

func createKeepassDatabase(vault BitwardenDatabase, path string, password string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root"

	subgroups := createSubgroups(vault)

	for _, item := range vault.Items {
		entry := getEntry(item)

		if item.FolderId == "" {
			rootGroup.Entries = append(rootGroup.Entries, entry)
		} else {
			subgroup := subgroups[item.FolderId]
			subgroup.Entries = append(subgroups[item.FolderId].Entries, entry)

			subgroups[item.FolderId] = subgroup
		}
	}

	for _, subgroup := range subgroups {
		rootGroup.Groups = append(rootGroup.Groups, subgroup)
	}

	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(password),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	db.LockProtectedEntries()

	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(db); err != nil {
		panic(err)
	}

	fmt.Printf("Wrote kdbx file at %v\n", path)
}

func createSubgroups(vault BitwardenDatabase) map[string]gokeepasslib.Group {
	subgroups := make(map[string]gokeepasslib.Group)

	// Add all subgroups
	for _, folder := range vault.Folders {
		subgroup := gokeepasslib.NewGroup()
		subgroup.Name = folder.Name
		subgroups[folder.Id] = subgroup
	}

	return subgroups
}

func getEntry(item BitwardenItem) gokeepasslib.Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, mkValue("Title", item.Name))
	entry.Values = append(entry.Values, mkValue("UserName", item.Login.Username))
	entry.Values = append(entry.Values, mkProtectedValue("Password", item.Login.Password))

	return entry
}

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value},
	}
}

func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value, Protected: w.NewBoolWrapper(true)},
	}
}

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

type Message struct {
	Name string
	Body string
	Time int64
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

func GenerateKeepassDatabase(opts options) {
	vault, exportErr := exportBitwardenVault(opts.BitwardenSession)

	if exportErr != nil {
		panic(exportErr)
	}

	createKeepassDatabase(vault, opts.DatabasePath, opts.DatabasePassword)

}

func exportBitwardenVault(bitwardenSession string) (BitwardenDatabase, error) {
	out, err := exec.Command("/usr/bin/bw", "export", "--format", "json", "--raw", "--session", bitwardenSession).Output()

	if err == nil {
		var db BitwardenDatabase
		err = json.Unmarshal(out, &db)
		return db, err
	}

	return BitwardenDatabase{}, err
}

func createKeepassDatabase(vault BitwardenDatabase, path string, password string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root"

	//var subgroups map[string]gokeepasslib.Group
	subgroups := make(map[string]gokeepasslib.Group)

	// Add all subgroups
	for _, folder := range vault.Folders {
		subgroup := gokeepasslib.NewGroup()
		subgroup.Name = folder.Name
		subgroups[folder.Id] = subgroup
	}

	for _, item := range vault.Items {
		entry := getEntry(item)

		if item.FolderId == "" {
			//add to root
			rootGroup.Entries = append(rootGroup.Entries, entry)
		} else {
			//add to subgroup
			subgroup := subgroups[item.FolderId]
			subgroup.Entries = append(subgroups[item.FolderId].Entries, entry)

			subgroups[item.FolderId] = subgroup
		}

	}

	for _, subgroup := range subgroups {
		rootGroup.Groups = append(rootGroup.Groups, subgroup)
	}

	// now create the database containing the root group
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

	// Lock entries using stream cipher
	db.LockProtectedEntries()

	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(db); err != nil {
		panic(err)
	}

	log.Printf("Wrote kdbx file: %s", path)
}

func getEntry(item BitwardenItem) gokeepasslib.Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, mkValue("Title", item.Name))
	entry.Values = append(entry.Values, mkValue("UserName", item.Login.Username))
	entry.Values = append(entry.Values, mkProtectedValue("Password", item.Login.Password))

	return entry
}

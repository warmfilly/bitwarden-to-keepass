package main

type BitwardenDatabase struct {
	Encrypted bool
	Folders   []BitwardenFolder
	Items     []BitwardenItem
}

type BitwardenFolder struct {
	Id   string
	Name string
}

type BitwardenItem struct {
	Id       string
	Name     string
	FolderId string
	Notes    string
	Login    BitwardenLogin
}

type BitwardenLogin struct {
	Username string
	Password string
}

package bw2kp

type bitwardenDatabase struct {
	Encrypted bool
	Folders   []bitwardenFolder
	Items     []bitwardenItem
}

type bitwardenFolder struct {
	Id   string
	Name string
}

type bitwardenItem struct {
	Id       string
	Name     string
	FolderId string
	Notes    string
	Login    bitwardenLogin
}

type bitwardenLogin struct {
	Username string
	Password string
}

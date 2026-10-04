package enums

type Status uint8

const (
	StatusInActive Status = iota + 1
	StatusActive
	StatusSuspended
	StatusBanned
)

package protocol

type Location = int

const (
	UnavailableStatus Location = 0
	LocalStatus       Location = 1
	RemoteStatus      Location = 2
)

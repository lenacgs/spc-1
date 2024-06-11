package protocol

type Status = int

const (
	UnavailableStatus Status = 0
	LocalStatus       Status = 1
	RemoteStatus      Status = 2
)

package id

//AccountID defines account id object.
type AccountID string

func (a AccountID) String() string {
	return string(a)
}

type TripID string

func (t TripID) String() string {
	return string(t)
}

//IdentityID defines identity id object.
type IdentityID string

func (t IdentityID) String() string {
	return string(t)
}


type CarID string

func (t CarID) String() string {
	return string(t)
}

type BlobID string

func (t BlobID) String() string {
	return string(t)
}

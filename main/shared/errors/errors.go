package errors

type EmailError struct {
	Description string
}

type PhoneNumberError struct {
	Description string
}

type DateError struct {
	Description string
}

type RentError struct {
	Description string
}

type MembershipStatusError struct {
	Description string
}

type WaitingListError struct {
	Description string
}

type NotFoundError struct {
	Description string
}

type RepositoryError struct {
	Description string
}

func (e EmailError) Error() string {
	return e.Description
}

func (p PhoneNumberError) Error() string {
	return p.Description
}

func (d DateError) Error() string {
	return d.Description
}

func (r RentError) Error() string {
	return r.Description
}

func (i MembershipStatusError) Error() string {
	return i.Description
}

func (w WaitingListError) Error() string {
	return w.Description
}

func (n NotFoundError) Error() string {
	return n.Description
}

func (r RepositoryError) Error() string {
	return r.Description
}

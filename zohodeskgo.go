package zohodeskgo

// Version # of zohodeskgo
const Version = "0.1.0"

func NewTickets(
	clientID string,
	clientSecret string,
	refreshToken string,
	pathAccessToken string,
	organizationID uint64,
	deparmentID uint64,
	contactID uint64,
) ZohoTickets {
	t := newTickets(
		clientID,
		clientSecret,
		refreshToken,
		pathAccessToken,
		organizationID,
		deparmentID,
		contactID,
	)

	t.doRefreshToken()

	return t
}

func NewUploads(
	clientID string,
	clientSecret string,
	refreshToken string,
	pathAccessToken string,
	organizationID uint64,
) ZohoUploads {
	u := newUploads(
		clientID,
		clientSecret,
		refreshToken,
		pathAccessToken,
		organizationID,
	)

	u.doRefreshToken()

	return u
}

type ZohoTickets interface {
	Create(payload map[string]interface{}) (ticket *Ticket, err error)

	doRefreshToken() (err error)
}

type ZohoUploads interface {
	Upload(filename string, file []byte) (upload *Upload, err error)

	doRefreshToken() (err error)
}

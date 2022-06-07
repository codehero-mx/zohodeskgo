package zohodeskgo

// Version # of resty
const Version = "0.0.1"

func NewTickets(
	clientID string,
	clientSecret string,
	refreshToken string,
	pathAccessToken string,
	organizationID uint64,
	deparmentID uint64,
	contactID uint64,
) *tickets {
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
) *uploads {
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

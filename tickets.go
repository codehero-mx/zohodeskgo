package zohodeskgo

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type ticket struct {
	ID string `json:"id"`
}

type tickets struct {
	zoho

	deparmentID uint64
	contactID   uint64

	ticket ticket
}

const ticketsUri = "/tickets"

func newTickets(
	clientID string,
	clientSecret string,
	refreshToken string,
	pathAccessToken string,
	organizationID uint64,
	deparmentID uint64,
	contactID uint64,
) *tickets {
	return &tickets{
		zoho: zoho{
			clientID:        clientID,
			clientSecret:    clientSecret,
			refreshToken:    refreshToken,
			pathAccessToken: pathAccessToken,
			organizationID:  organizationID,
		},
		deparmentID: deparmentID,
		contactID:   contactID,
	}
}

func (t *tickets) Create(payload map[string]interface{}) (err error) {
	payload["departmentId"] = t.deparmentID
	payload["contactId"] = t.contactID

	resp, err := t.wrapRequestor(func() (*resty.Response, error) {
		rest, err := t.requestor()
		if err != nil {
			return nil, err
		}

		return rest.
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			Post(ticketsUri)
	})
	if err != nil {
		return
	}

	re := new(respError)

	if resp.StatusCode() != 200 {
		json.Unmarshal(resp.Body(), &re)
	} else {
		json.Unmarshal(resp.Body(), &t.ticket)
	}

	fmt.Printf("ID: %s\n", t.ticket.ID)

	return nil
}

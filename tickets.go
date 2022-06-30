package zohodeskgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Ticket struct {
	ID           string                 `json:"id,omitempty"`
	TicketNumber string                 `json:"ticketNumber,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Subcategory  string                 `json:"subCategory,omitempty"`
	Subject      string                 `json:"subject,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	Channel      string                 `json:"channel,omitempty"`
	DepartmentID uint64                 `json:"departmentId,omitempty"`
	ContactID    uint64                 `json:"contactId,omitempty"`
	Status       string                 `json:"status,omitempty"`
	StatusType   string                 `json:"statusType,omitempty"`
	DueDate      time.Time              `json:"dueDate,omitempty"`
	Language     string                 `json:"language,omitempty"`
	Resolution   string                 `json:"resolution,omitempty"`
	Priority     string                 `json:"priority,omitempty"`
	CF           map[string]interface{} `json:"cf,omitempty"`
	Deleted      string                 `json:"isDeleted,omitempty"`
	Trashed      string                 `json:"isTrashed,omitempty"`
	ClosedTime   time.Time              `json:"closedTime,omitempty"`
	WebURL       string                 `json:"webUrl,omitempty"`
	Uploads      []string               `json:"uploads,omitempty"`
	ModifiedBy   string                 `json:"modifiedBy,omitempty"`
	ModifiedTime time.Time              `json:"modifiedTime,omitempty"`
	CreatedTime  time.Time              `json:"createdTime,omitempty"`
}

type tickets struct {
	zoho

	deparmentID uint64
	contactID   uint64
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
) ZohoTickets {
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

func (t *tickets) Create(payload map[string]interface{}) (ticket *Ticket, err error) {
	payload["departmentId"] = t.deparmentID

	if t.contactID > 0 {
		payload["contactId"] = t.contactID
	}

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
		fmt.Printf("%v\n", err.Error())
		return
	}

	re := new(respError)

	if resp.StatusCode() != 200 {
		fmt.Printf("%v\n", string(resp.Body()))

		json.Unmarshal(resp.Body(), &re)

		err = errors.New("an error occurred while creating a ticket")
	} else {
		json.Unmarshal(resp.Body(), &ticket)
	}

	return
}

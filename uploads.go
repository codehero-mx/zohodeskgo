package zohodeskgo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Upload struct {
	ID          string    `json:"id"`
	CreatorId   string    `json:"creatorId"`
	CreatedTime time.Time `json:"createdTime"`
	Name        string    `json:"name"`
	IsPublic    bool      `json:"isPublic"`
	Size        string    `json:"size"`
	Href        string    `json:"href"`
}

type uploads struct {
	zoho
}

const uploadsUri = "/uploads"

func newUploads(
	clientID string,
	clientSecret string,
	refreshToken string,
	pathAccessToken string,
	organizationID uint64,
) *uploads {
	return &uploads{
		zoho: zoho{
			clientID:        clientID,
			clientSecret:    clientSecret,
			refreshToken:    refreshToken,
			pathAccessToken: pathAccessToken,
			organizationID:  organizationID,
		},
	}
}

func (u *uploads) Upload(file []byte) (upload *Upload, err error) {
	resp, err := u.wrapRequestor(func() (*resty.Response, error) {
		rest, err := u.requestor()
		if err != nil {
			return nil, err
		}

		return rest.
			SetBody(file).
			Post(uploadsUri)
	})
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	re := new(respError)

	if resp.StatusCode() != 200 {
		json.Unmarshal(resp.Body(), &re)
	} else {
		json.Unmarshal(resp.Body(), &upload)
	}

	return
}

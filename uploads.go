package zohodeskgo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Upload struct {
	Size        string    `json:"size"`
	CreatorId   string    `json:"creatorId"`
	Name        string    `json:"name"`
	CreatedTime time.Time `json:"createdTime"`
	IsPublic    bool      `json:"isPublic"`
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
		return
	}

	re := new(respError)

	if resp.StatusCode() != 200 {
		json.Unmarshal(resp.Body(), &re)
	} else {
		json.Unmarshal(resp.Body(), &upload)
	}

	fmt.Printf("%v\n", upload)

	return
}

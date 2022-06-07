package zohodeskgo

import (
	"bytes"
	"encoding/json"
	"errors"
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
) ZohoUploads {
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

func (u *uploads) Upload(filename string, file []byte) (upload *Upload, err error) {
	resp, err := u.wrapRequestor(func() (*resty.Response, error) {
		rest, err := u.requestor()
		if err != nil {
			return nil, err
		}

		return rest.
			SetFileReader("file", filename, bytes.NewReader(file)).
			Post(uploadsUri)
	})
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	re := new(respError)

	if resp.StatusCode() != 200 {
		fmt.Printf("%v\n", string(resp.Body()))

		json.Unmarshal(resp.Body(), &re)

		err = errors.New("an error occurred while uploading the file")
	} else {
		json.Unmarshal(resp.Body(), &upload)
	}

	return
}

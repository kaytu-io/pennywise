package auth0

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/pkg"
	"io"
	"net/http"
)

type ResponseAbout struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email-verified"`
}

func RequestAbout(accessToken string) (ResponseAbout, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/userinfo", pkg.Auth0Hostname), nil)
	if err != nil {
		return ResponseAbout{}, fmt.Errorf("[requestAbout] : %v", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ResponseAbout{}, fmt.Errorf("[requestAbout] : %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ResponseAbout{}, fmt.Errorf("[requestAbout] : %v", err)
	}
	err = res.Body.Close()
	if err != nil {
		return ResponseAbout{}, fmt.Errorf("[requestAbout] : %v", err)
	}
	response := ResponseAbout{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ResponseAbout{}, fmt.Errorf("[requestAbout] : %v", err)
	}
	return response, nil
}

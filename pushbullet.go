package akari

import (
	"bytes"
	"errors"
	"net/http"
)

type PushbulletPush struct {
	pushType    string
	title       string
	body        string
	accessToken string
}

// Push sends a POST request to Pushbullet server in order to make a Pushbullet push notification.
func (p *PushbulletPush) Push() (err error) {
	jsonStr := `{"body":"` + p.body + `","title":"` + p.title + `","type":"` + p.pushType + `"}`
	url := "https://api.pushbullet.com/v2/pushes"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", p.accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return
}

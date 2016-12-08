// Copyright © 2016 nrechn <nrechn@gmail.com>
//
// This file is part of akari.
//
// akari is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// akari is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with akari. If not, see <http://www.gnu.org/licenses/>.
//

package akari

import (
	"bytes"
	"errors"
	"net/http"
)

// Push sends a POST request to Pushbullet server in order to make a Pushbullet push notification.
func (p *PushbulletPush) Push() (err error) {
	jsonStr := `{"body":"` + p.Body + `","title":"` + p.Title + `","type":"` + p.PushType + `"}`
	url := "https://api.pushbullet.com/v2/pushes"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", p.AccessToken)
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

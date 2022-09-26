package providers

import (
	"encoding/json"
	"io"
	"net/http"
)

type Linode struct {
	Provider
}

func (l Linode) GetRegions() ([]Region, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.linode.com/v4/regions", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Dekamik/autovpn")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				panic(err)
			}
		}(res.Body)
	}

	body := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	data := body["data"].([]interface{})
	regions := make([]Region, len(data))

	for i, region := range data {
		regionData := region.(map[string]interface{})
		regions[i] = Region{Id: regionData["id"].(string), Country: regionData["country"].(string)}
	}

	return regions, nil
}

func (l Linode) CreateServer() (Instance, error) {
	// TODO: Implement
	return Instance{}, nil
}

func (l Linode) DestroyServer(instance Instance) error {
	// TODO: Implement
	return nil
}

package providers

import (
	"autovpn/helpers"
	"autovpn/options"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Linode struct {
	Provider
}

type regionRes struct {
	Data []struct {
		Id      string
		Country string
	}
}

type createRes struct {
	Id   float64
	Ipv4 []string
}

func (l Linode) GetRegions(silent bool) ([]Region, error) {
	if !silent {
		fmt.Print("Getting regions... ")
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.linode.com/v4/regions", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Dekamik/autovpn")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("region download caused error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("region download returned %d %s", res.StatusCode, res.Status)
	}

	body := &regionRes{}
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return nil, fmt.Errorf("region download caused error: %w", err)
	}

	regions := make([]Region, len(body.Data))
	for i, region := range body.Data {
		regions[i] = Region{Id: region.Id, Country: region.Id}
	}

	if !silent {
		fmt.Println("OK")
	}
	return regions, nil
}

func (l Linode) CreateServer(arguments options.Arguments, config options.Config) (*Instance, error) {
	fmt.Print("Creating server... ")

	client := http.Client{}
	rootPass, err := helpers.GeneratePassword(64)
	if err != nil {
		return nil, err
	}

	conf := config.Providers[arguments.Provider]
	var jsonData = []byte(
		fmt.Sprintf("{\"image\":\"%s\",\"region\":\"%s\",\"root_pass\":\"%s\",\"type\":\"%s\"}",
			conf.Image, arguments.Region, rootPass, conf.TypeSlug))

	req, err := http.NewRequest(http.MethodPost, "https://api.linode.com/v4/linode/instances", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("server creation caused error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+conf.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Dekamik/autovpn")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server creation returned %d %s", res.StatusCode, res.Status)
	}

	body := &createRes{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("server creation caused error: %w", err)
	}

	instance := &Instance{
		Id:        fmt.Sprintf("%f", body.Id),
		IpAddress: body.Ipv4[0],
	}

	fmt.Println("OK")

	// TODO: Await provisioning and booting

	return instance, nil
}

func (l Linode) DestroyServer(instance Instance, token string) error {
	fmt.Print("Destroying server... ")

	client := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "https://api.linode.com/v4/linode/instances/"+instance.Id, nil)
	if err != nil {
		return fmt.Errorf("server destruction caused error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Dekamik/autovpn")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("server destruction caused error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server creation returned %d %s", res.StatusCode, res.Status)
	}

	fmt.Println("OK")
	return nil
}

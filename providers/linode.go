package providers

import (
	"autovpn/data"
	"autovpn/helpers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Linode struct {
	Client
}

type regionRes struct {
	Data []struct {
		Id      string
		Country string
	}
}

type instanceRes struct {
	Id     float64
	Ipv4   []string
	Status string
	Tags   []string
}

type listRes struct {
	Data []instanceRes
}

func (l Linode) getRegions(args data.ArgsBundle) ([]Region, error) {
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
		regions[i] = Region{Id: region.Id, Country: region.Country}
	}

	return regions, nil
}

func (l Linode) getInstances(args data.ArgsBundle) ([]data.Instance, error) {
	client := http.Client{}
	conf := args.Config.Providers["linode"]

	req, err := http.NewRequest(http.MethodGet, "https://api.linode.com/v4/linode/instances/", nil)
	if err != nil {
		return nil, fmt.Errorf("server check caused error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+conf.Key)
	req.Header.Set("User-Agent", "Dekamik/autovpn")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server list returned %d %s", res.StatusCode, res.Status)
	}

	body := &listRes{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("server list caused error: %w", err)
	}

	instances := make([]data.Instance, len(body.Data))
	for i, instance := range body.Data {
		instances[i] = data.Instance{
			Id:        fmt.Sprintf("%f", instance.Id),
			IpAddress: "",
			User:      "root",
			Pass:      "",
			SshPort:   22,
			Tags:      instance.Tags,
		}
	}
	return instances, nil
}

func (l Linode) createServer(args data.ArgsBundle) (*data.Instance, error) {
	client := http.Client{}
	config := args.Config.Providers[args.Arguments.Provider]

	var rootPass string
	if len(args.Config.Overrides.RootPass) != 0 {
		rootPass = args.Config.Overrides.RootPass
	} else {
		var err error
		rootPass, err = helpers.GeneratePassword(64)
		if err != nil {
			return nil, err
		}
	}

	var jsonData = []byte(
		fmt.Sprintf("{\"image\":\"%s\",\"region\":\"%s\",\"root_pass\":\"%s\",\"tags\":[\"%s\"],\"type\":\"%s\"}",
			config.Image, args.Arguments.Region, rootPass, InstanceTag, config.TypeSlug))

	req, err := http.NewRequest(http.MethodPost, "https://api.linode.com/v4/linode/instances", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("server creation caused error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Key)
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

	body := &instanceRes{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("server creation caused error: %w", err)
	}

	instance := &data.Instance{
		Id:        fmt.Sprintf("%.0f", body.Id),
		IpAddress: body.Ipv4[0],
		User:      "root",
		Pass:      rootPass,
		SshPort:   22,
	}

	return instance, nil
}

func (l Linode) awaitProvisioning(args data.ArgsBundle) error {
	token := args.Config.Providers["linode"].Key
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://api.linode.com/v4/linode/instances/"+args.Instance.Id, nil)
	if err != nil {
		return fmt.Errorf("server check caused error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Dekamik/autovpn")

	for {
		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("server check caused error: %w", err)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("server check returned %d %s", res.StatusCode, res.Status)
		}

		body := &instanceRes{}
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return fmt.Errorf("server check caused error: %w", err)
		}

		if body.Status == "running" {
			return nil
		}
		time.Sleep(time.Second * 5)
	}
}

func (l Linode) destroyServer(args data.ArgsBundle) error {
	token := args.Config.Providers["linode"].Key
	client := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "https://api.linode.com/v4/linode/instances/"+args.Instance.Id, nil)
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
		return fmt.Errorf("server destruction returned %d %s", res.StatusCode, res.Status)
	}

	return nil
}

func (l Linode) connect(_ data.ArgsBundle) error {
	// Nothing needs to be done
	return nil
}

func (l Linode) failSafeSetup(args data.ArgsBundle) ([]string, error) {
	commands := []string{
		fmt.Sprintf(
			"echo \"$(date +%%M) $(($(($(date +%%H) + %d)) %% 24)) * * * /usr/bin/env curl -H 'Authorization: Bearer %s' -X DELETE https://api.linode.com/v4/linode/instances/%s\" > /etc/crontab",
			args.Config.Agent.ServerTtlHours, args.Config.Providers["linode"].Key, args.Instance.Id),
		"crontab /etc/crontab",
	}
	return commands, nil
}

package transcode

import (
	"encoding/json"
	"net/http"

	"k8s.io/helm/pkg/proto/hapi/services"
)

// List proxies a list request
func (p *Proxy) List(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.ListReleasesRequest{}

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	var res *services.ListReleasesResponse
	err = p.do(func(rlc services.ReleaseServiceClient) error {
		s, err := rlc.ListReleases(NewContext(), req)
		if err != nil {
			return err
		}
		res, err = s.Recv()
		return err
	})

	if err != nil {
		return err
	}

	data, err = json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// History lists the history for a named release.
func (p *Proxy) History(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.GetHistoryRequest{}
	var res *services.GetHistoryResponse

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	err = p.do(func(rlc services.ReleaseServiceClient) error {
		var err error
		res, err = rlc.GetHistory(NewContext(), req)
		return err
	})

	if err != nil {
		return err
	}

	data, err = json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

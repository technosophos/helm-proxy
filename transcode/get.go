package transcode

import (
	"encoding/json"
	"net/http"

	//"golang.org/x/net/context"
	"k8s.io/helm/pkg/proto/hapi/services"
)

// Get retrieves a release record.
func (p *Proxy) Get(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.GetReleaseContentRequest{}
	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	var res *services.GetReleaseContentResponse
	err = p.do(func(rlc services.ReleaseServiceClient) error {
		ctx := NewContext()
		var err error
		res, err = rlc.GetReleaseContent(ctx, req)
		if err != nil {
			return err
		}
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

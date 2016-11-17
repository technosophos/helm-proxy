package transcode

import (
	"encoding/json"
	"net/http"

	//"golang.org/x/net/context"
	"k8s.io/helm/pkg/proto/hapi/services"
)

// Install loads a package into Tiller for installation,
func (p *Proxy) Install(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.InstallReleaseRequest{}
	var res *services.InstallReleaseResponse

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	err = p.do(func(rlc services.ReleaseServiceClient) error {
		ctx := NewContext()
		var err error
		res, err = rlc.InstallRelease(ctx, req)
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

// Upgrade loads a new version of a release to Tiller for upgrading and existing release
func (p *Proxy) Upgrade(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.UpdateReleaseRequest{}
	var res *services.UpdateReleaseResponse

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	err = p.do(func(rlc services.ReleaseServiceClient) error {
		ctx := NewContext()
		var err error
		res, err = rlc.UpdateRelease(ctx, req)
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

// Uninstall destroys a release
func (p *Proxy) Uninstall(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.UninstallReleaseRequest{}
	var res *services.UninstallReleaseResponse

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	err = p.do(func(rlc services.ReleaseServiceClient) error {
		ctx := NewContext()
		var err error
		res, err = rlc.UninstallRelease(ctx, req)
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

// Rollback rolls back to a previous release version.
func (p *Proxy) Rollback(w http.ResponseWriter, r *http.Request) error {
	data, err := body(r)
	if err != nil {
		return err
	}

	req := &services.RollbackReleaseRequest{}
	var res *services.RollbackReleaseResponse

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	err = p.do(func(rlc services.ReleaseServiceClient) error {
		ctx := NewContext()
		var err error
		res, err = rlc.RollbackRelease(ctx, req)
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

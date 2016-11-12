/*Package transcode proxies a request between local and a remote gRPC server.

It returns JSON data.
*/
package transcode

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func NewContext() context.Context {
	md := metadata.Pairs("x-helm-api-client", "2.0.0")
	return metadata.NewContext(context.TODO(), md)
}

type Proxy struct {
	// Fields on here should be read-only to avoid races, or else we should add
	// a mutex on here.

	host string
}

func New(host string) *Proxy {
	return &Proxy{
		host: host,
	}
}

func body(r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	return b, err
}

type doFunc func(rlc services.ReleaseServiceClient) error

// do executes a particular function with a gRPC client and context.
func (p *Proxy) do(fn doFunc) error {
	c, err := grpc.Dial(p.host, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()
	rlc := services.NewReleaseServiceClient(c)
	return fn(rlc)
}

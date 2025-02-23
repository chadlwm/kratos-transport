package fasthttp

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/valyala/fasthttp"

	kHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"

	api "github.com/tx7do/kratos-transport/_example/api/protobuf"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		WithAddress(":8800"),
	)

	srv.GET("/login/*param", func(c *fasthttp.RequestCtx) {
		_, _ = c.WriteString("Hello World!")
	})

	srv.GET("/hygrothermograph", func(c *fasthttp.RequestCtx) {
		var out api.Hygrothermograph
		out.Humidity = strconv.FormatInt(int64(rand.Intn(100)), 10)
		out.Temperature = strconv.FormatInt(int64(rand.Intn(100)), 10)
		_ = json.NewEncoder(c.Response.BodyWriter()).Encode(&out)
	})

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("expected nil got %v", err)
		}
	}()
}

func TestClient(t *testing.T) {
	ctx := context.Background()

	cli, err := kHttp.NewClient(ctx,
		kHttp.WithEndpoint("127.0.0.1:8800"),
	)
	assert.Nil(t, err)
	assert.NotNil(t, cli)

	resp, err := GetHygrothermograph(ctx, cli, nil, kHttp.EmptyCallOption{})
	assert.Nil(t, err)
	t.Log(resp)
}

func GetHygrothermograph(ctx context.Context, cli *kHttp.Client, in *api.Hygrothermograph, opts ...kHttp.CallOption) (*api.Hygrothermograph, error) {
	var out api.Hygrothermograph

	pattern := "/hygrothermograph"
	path := binding.EncodeURL(pattern, in, true)

	opts = append(opts, kHttp.Operation("/GetHygrothermograph"))
	opts = append(opts, kHttp.PathTemplate(pattern))

	err := cli.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

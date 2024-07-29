package ngrok

import (
	"context"
	"net/http"

	"golang.ngrok.com/ngrok"
	ngrokcfg "golang.ngrok.com/ngrok/config"

	"github.com/lakrizz/logsync/internal/config"
)

// 2dXpKgxBUVNpB6hVtm1rM6vtPv0_5DAGrYvBLAN3ytQNN6wQH
func Start(ctx context.Context, cfg *config.Config, handler http.HandlerFunc) (string, chan error, error) {
	listener, err := ngrok.Listen(ctx,
		ngrokcfg.HTTPEndpoint(),
		ngrok.WithAuthtoken(cfg.Ngrok.AuthToken),
	)
	if err != nil {
		return "", nil, err
	}

	ch := make(chan error, 1)
	go listen(listener, handler, ch)
	return listener.URL(), ch, nil
}

func listen(listener ngrok.Tunnel, handler http.HandlerFunc, ch chan error) {
	err := http.Serve(listener, http.HandlerFunc(handler))
	if err != nil {
		ch <- err
	}
}

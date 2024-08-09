package ngrok

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
	nlog "golang.ngrok.com/ngrok/log"
)

// Simple logger that forwards to the Go standard logger.
type logger struct {
	lvl nlog.LogLevel
}

func (l *logger) log(ctx context.Context, lvl nlog.LogLevel, msg string, data map[string]any) {
	if lvl > l.lvl {
		return
	}
	lvlName, _ := nlog.StringFromLogLevel(lvl)
	log.Printf("[%s] %s %v", lvlName, msg, data)
}

var l *logger = &logger{
	lvl: nlog.LogLevelDebug,
}

func Listen(ctx context.Context, token string) error {
	u := url.URL{Host: "0.0.0.0:5432"}
	forwarder, err := ngrok.ListenAndForward(ctx,
		&u,
		config.TCPEndpoint(),
		ngrok.WithAuthtoken(token),
	)
	if err != nil {
		return fmt.Errorf("could not create lister: %w", err)
	}

	return run(ctx, forwarder)
}

func run(ctx context.Context, fwd ngrok.Forwarder) error {
	for {
		l.log(ctx, nlog.LogLevelInfo, "ingress established", map[string]any{
			"url": fwd.URL(),
		})

		err := fwd.Wait()
		if err == nil {
			return nil
		}
		l.log(ctx, nlog.LogLevelWarn, "accept error. now setting up a new forwarder.",
			map[string]any{"err": err})
	}
}

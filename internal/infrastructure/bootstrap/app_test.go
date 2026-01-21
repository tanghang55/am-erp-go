package bootstrap

import "testing"

func TestBootstrapBuildsEngine(t *testing.T) {
    app, err := Build()
    if err != nil || app.Engine == nil {
        t.Fatalf("expected engine, got err=%v", err)
    }
}

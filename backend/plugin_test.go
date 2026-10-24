package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"antennae.space/plugin-sdk/backend"
)

// ide plays the IDE: it sends requests to the backend, answers its core calls with a fixed
// project and collects its events.
type ide struct {
	t      *testing.T
	w      io.Writer
	r      *bufio.Scanner
	events []string
}

func start(t *testing.T) *ide {
	in, toB := io.Pipe()
	fromB, out := io.Pipe()
	go func() { _ = backend.ServeIO(context.Background(), plugin(), in, out); out.Close() }()
	t.Cleanup(func() { toB.Close() })
	f := &ide{t: t, w: toB, r: bufio.NewScanner(fromB)}
	f.ask("initialize", map[string]any{"apiVersion": 1, "plugin": "example", "grants": []string{"core:fs.list"}, "project": map[string]any{"id": "p", "local": true}})
	return f
}

type msg struct {
	ID     json.RawMessage `json:"id,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *backend.Error  `json:"error,omitempty"`
}

func (f *ide) send(m any) {
	b, _ := json.Marshal(m)
	_, _ = f.w.Write(append(b, '\n'))
}

// ask sends a request and answers what the backend asks until its answer comes.
func (f *ide) ask(method string, params any) msg {
	f.t.Helper()
	f.send(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	for f.r.Scan() {
		var m msg
		if err := json.Unmarshal(f.r.Bytes(), &m); err != nil {
			f.t.Fatal(err)
		}
		switch m.Method {
		case "core":
			f.send(map[string]any{"jsonrpc": "2.0", "id": m.ID, "result": []map[string]any{
				{"name": "README.md", "path": "README.md", "size": 120},
				{"name": "src", "path": "src", "dir": true},
				{"name": "go.mod", "path": "go.mod", "size": 30},
			}})
		case "emit":
			f.events = append(f.events, string(m.Params))
		case "":
			return m
		}
	}
	f.t.Fatal("the backend ended")
	return msg{}
}

func TestEntries(t *testing.T) {
	f := start(t)
	m := f.ask("call", map[string]any{"method": "example.entries"})
	var list []Entry
	if err := json.Unmarshal(m.Result, &list); err != nil || len(list) != 3 || list[0].Name != "src" {
		t.Fatalf("entries, folders first: %s %v", m.Result, m.Error)
	}
}

func TestSummaryTool(t *testing.T) {
	f := start(t)
	m := f.ask("tool", map[string]any{"name": "example_summary", "args": map[string]any{}})
	var r backend.Result
	_ = json.Unmarshal(m.Result, &r)
	if r.Error || r.Content != "The root of the project holds 1 folder and 2 files (150 bytes)." {
		t.Fatalf("tool: %+v", r)
	}
}

func TestRefreshEmits(t *testing.T) {
	f := start(t)
	f.ask("call", map[string]any{"method": "example.refresh"})
	f.ask("call", map[string]any{"method": "example.summary"})
	if len(f.events) != 1 || !strings.Contains(f.events[0], `"event":"example.changed"`) {
		t.Fatalf("events: %v", f.events)
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"antennae.space/plugin-sdk/backend"
)

// Entry is an entry of the project, as the IDE's fs.list gives it.
type Entry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Dir  bool   `json:"dir"`
	Size int64  `json:"size"`
}

// Summary is what the tool shows and the assistant gets.
type Summary struct {
	Folders int   `json:"folders"`
	Files   int   `json:"files"`
	Bytes   int64 `json:"bytes"`
}

// plugin is what the backend serves: the methods of its views (example.*) and its assistant
// tools, as named in plugin.json.
func plugin() backend.Plugin {
	return backend.Plugin{
		Methods: map[string]backend.Method{
			"example.entries": func(ctx context.Context, _ json.RawMessage) (any, error) { return entries(ctx) },
			"example.summary": func(ctx context.Context, _ json.RawMessage) (any, error) { return summary(ctx) },
			// Tells every window of the project that the entries changed: their views reload.
			"example.refresh": func(ctx context.Context, _ json.RawMessage) (any, error) {
				return nil, backend.Emit(ctx, "example.changed", nil)
			},
		},
		Tools: map[string]backend.Tool{
			"example_summary": func(ctx context.Context, _ json.RawMessage) (backend.Result, error) {
				s, err := summary(ctx)
				if err != nil {
					return backend.Result{}, err
				}
				// The model reads English; the summary is shown to the user in the assistant.
				return backend.Result{
					Content: fmt.Sprintf("The root of the project holds %s and %s (%d bytes).", plural(s.Folders, "folder"), plural(s.Files, "file"), s.Bytes),
					Summary: plural(s.Folders, "folder") + ", " + plural(s.Files, "file"),
				}, nil
			},
		},
	}
}

// entries are those at the root of the project, folders first.
func entries(ctx context.Context) ([]Entry, error) {
	var list []Entry
	if err := backend.Core(ctx, "fs.list", map[string]any{"path": ""}, &list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].Dir && !list[j].Dir })
	return list, nil
}

func summary(ctx context.Context) (Summary, error) {
	list, err := entries(ctx)
	if err != nil {
		return Summary{}, err
	}
	var s Summary
	for _, e := range list {
		if e.Dir {
			s.Folders++
		} else {
			s.Files++
			s.Bytes += e.Size
		}
	}
	return s, nil
}

func plural(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return fmt.Sprintf("%d %ss", n, word)
}

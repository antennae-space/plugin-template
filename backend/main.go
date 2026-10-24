// The backend of the plugin: a program the IDE starts confined (Landlock and seccomp) and
// speaks to in JSON-RPC over stdio, through the SDK (antennae.space/plugin-sdk/backend). It
// reaches its own folder, its data folder (its HOME), the project as granted (files:read), and
// asks the IDE for the rest: here, the entries of the project through fs.list (core:fs.list),
// which also works for a project over SSH.
package main

import "antennae.space/plugin-sdk/backend"

func main() {
	backend.Serve(plugin())
}

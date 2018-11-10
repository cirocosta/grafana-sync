grafana-sync - Keeps your Grafana dashboards in sync

HOW IT WORKS

	At each time that it's run, `grafana-sync` gathers information
	about dashboards from a particular source of truth and then
	updates the state of the other.

	What a source of truth is is configurable:
	- a remote grafana instance, or
	- the local filesystem.

	Once differences are perceived, these are highlighted and
	asked for user's input to proceed with applying necessary mutations.


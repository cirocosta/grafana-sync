package main

type pushCommand struct{
	Datasource string `long:"data-source" short:"d" required:"true" description:"datasource used by the dashboards"`
}

func (p *pushCommand) Execute(args []string) (err error) {
	return
}

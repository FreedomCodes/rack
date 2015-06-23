package instancecommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", util.IDOrNameUsage("instance")),
	Description: "Retrieves an existing server",
	Action:      actionGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysGet = []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor", "Public IPv4", "Public IPv6", "Private IPv4", "KeyName"}

type paramsGet struct {
	server string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).server = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	id, err := command.Context().IDOrName(osServers.IDFromName)
	resource.Params.(*paramsGet).server = id
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	serverID := resource.Params.(*paramsGet).server
	server, err := servers.Get(command.Context().ServiceClient, serverID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = serverSingle(server)
}

func (command *commandGet) StdinField() string {
	return "id"
}

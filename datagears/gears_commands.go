package datagears

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// NOTE: Dump registration command
var CmdDumpRegistrations = "RG.DUMPREGISTRATIONS"

// NOTE: Dump executions command
var CmdDumpExecutions = "RG.DUMPEXECUTIONS"

// NOTE: Execution command
var CmdExecute = "RG.PYEXECUTE"
var ArgUnblockingExecution = "UNBLOCKING"
var ArgRequirements = "REQUIREMENTS"

// NOTE: unregister
var CmdUnregister = "RG.UNREGISTER"

//GearsCommands
type GearsCommands struct {
	gearsInstance *GearsInstance
}

//NewGearsCommands
func NewGearsCommands(gearsInstance *GearsInstance) *GearsCommands {
	return &GearsCommands{gearsInstance: gearsInstance}
}

//Execute
func (cmd *GearsCommands) Execute(gear *DGGear) *redis.Cmd {
	cmdParts := make([]interface{}, 0)
	cmdParts = append(cmdParts, CmdExecute, gear.Script())

	if !gear.IsBlocking() {
		cmdParts = append(cmdParts, ArgUnblockingExecution)
	}

	if len(gear.Requirements) > 0 {
		cmdParts = append(cmdParts, ArgRequirements)
		for _, requirement := range gear.Requirements {
			cmdParts = append(cmdParts, requirement)
		}
	}

	ctx := context.Background()
	return cmd.gearsInstance.client.Do(
		ctx,
		cmdParts...,
	)
}

//DumpRegistrations
func (cmd *GearsCommands) DumpRegistrations() ([]Registration, error) {
	ctx := context.Background()
	redisCmd := cmd.gearsInstance.client.Do(ctx, CmdDumpRegistrations)

	serializer := &DumpRegistrationSerializer{redisCmd}
	return serializer.Parse()
}

func (cmd *GearsCommands) Unregister(dgId string) (*redis.Cmd, error) {
	registrations, err := cmd.DumpRegistrations()
	if err != nil {
		return nil, err
	}

	for _, reg := range registrations {
		if dgId == reg.Desc.ID {
			ctx := context.Background()
			redisCmd := cmd.gearsInstance.client.Do(ctx, []interface{}{CmdUnregister, reg.ID}...)
			return redisCmd, nil
		}
	}

	return nil, ErrNotFound
}

package datagears

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

//DataGears
type DataGears struct {
	manifest *DGManifest

	remote   *GearsInstance
	commands *GearsCommands
}

//NewDataGears constructor.
func NewDataGears(remoteName string, manifestPath string) (*DataGears, error) {
	manifest := NewDGManifest(manifestPath)
	remote, err := manifest.GetRemote(remoteName)
	if err != nil {
		return nil, err
	}

	gearsInstance := NewGearsInstance(fmt.Sprintf("%s:%d", remote.Host, remote.Port)).SetDB(remote.Database)
	gearsInstance.Build()

	return &DataGears{
		manifest: manifest,
		remote:   gearsInstance,
		commands: NewGearsCommands(gearsInstance),
	}, nil
}

//ListRegistrations for all attached trigger functions.
func (dg *DataGears) ListRegistrations() ([]Registration, error) {
	return dg.commands.DumpRegistrations()
}

//RemoveRegistration of an attached trigger function.
func (dg *DataGears) RemoveRegistration(dgId string) (*redis.Cmd, error) {
	return dg.commands.Unregister(dgId)
}

//DeployGear deploy specified gear to the remote.
func (dg *DataGears) DeployGear(gearName string) (*redis.Cmd, error) {
	gear, err := dg.manifest.GetGear(gearName)
	if err != nil {
		return nil, err
	}

	err = gear.Build()
	if err != nil {
		return nil, err
	}

	return dg.commands.Execute(gear), nil
}

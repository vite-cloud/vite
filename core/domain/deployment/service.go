package deployment

import (
	"context"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/runtime"
)

// Deployment holds the information needed to deploy a service.
type Deployment struct {
	ID     string
	Docker *runtime.Client
}

// Deploy deploys a service.
func (d *Deployment) Deploy(ctx context.Context, service *config.Service) error {

	return nil
	//if service.HasDependencies() {
	//	net, err := d.CreateServiceNetwork(service)
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = d.ConnectServiceNetwork(service, net)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//err := d.PullImage(service.Image)
	//if err != nil {
	//	return err
	//}
	//
	//id, err := d.CreateContainer(service)
	//if err != nil {
	//	return err
	//}
	//
	//err = d.RunHooks(id, service.Hooks.Prestart)
	//if err != nil {
	//	return err
	//}
	//
	//err = d.StartContainer(id)
	//if err != nil {
	//	return err
	//}
	//
	//err = d.RunHooks(id, service.Hooks.Poststart)
	//if err != nil {
	//	return err
	//}
	//
	//err = d.EnsureContainerIsRunning(id)
	//if err != nil {
	//	if err2 := d.Docker.ContainerStop(ctx, id); err2 != nil {
	//		return fmt.Errorf("%w (cleanup failed: %s)", err, err2)
	//	}
	//	return err
	//}
	//
	//return nil
}
package factory

import (
	"github.com/appscode/wheel/pkg/extpoints"
	"golang.org/x/net/context"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type KubeconfigConnector struct {
	*TunnelConnector

	Context string
}

var _ extpoints.Connector = &KubeconfigConnector{}

const (
	UIDKubeconfigConnector = "kubeconfig"
)

func (c *KubeconfigConnector) UID() string {
	return UIDKubeconfigConnector
}

func (c *KubeconfigConnector) Connect(ctx context.Context) (rls.ReleaseServiceClient, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	addr, err := c.GetTillerAddr(client, config)
	if err != nil {
		return nil, err
	}
	conn, err := Connect(addr)
	if err != nil {
		return nil, err
	}
	return rls.NewReleaseServiceClient(conn), nil
}

func (c *KubeconfigConnector) getConfig() (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{
		CurrentContext:  c.Context,
		ClusterDefaults: clientcmd.ClusterDefaults,
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
}

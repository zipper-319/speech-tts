package nacos

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"speech-tts/internal/conf"
)

type RegistryDiscovery interface {
	registry.Registrar
	registry.Discovery
}

func NewRegister(config *conf.Data) RegistryDiscovery {
	serviceConfig := []constant.ServerConfig{
		{
			IpAddr:      config.Nacos.Addr,
			Port:        uint64(config.Nacos.Port),
			ContextPath: config.Nacos.ContextPath,
		},
	}
	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(config.Nacos.NamespaceId), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(uint64(config.Nacos.TimeoutMs)),
		constant.WithNotLoadCacheAtStart(config.Nacos.NotLoadCacheAtStart),
		constant.WithLogDir(config.Nacos.LogDir),
		constant.WithCacheDir(config.Nacos.CacheDir),
		constant.WithLogLevel(config.Nacos.LogLevel),
	)

	client, _ := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serviceConfig,
	})

	return nacos.New(client, nacos.WithGroup(config.Nacos.Group))
}

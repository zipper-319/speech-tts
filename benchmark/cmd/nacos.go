package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"speech-tts/internal/conf"
)

var bc conf.Bootstrap

func init() {
	c := config.New(
		config.WithSource(
			file.NewSource("./configs"),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
}

func getNacosService(config *conf.Data) error {
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

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serviceConfig,
	})
	if err != nil {
		return err
	}
	//if _, err := client.RegisterInstance(vo.RegisterInstanceParam{
	//	ServiceName: "speech-tts",
	//	Ip:          util.LocalIP(),
	//	Port:        3012,
	//	Healthy:     true,
	//	Ephemeral:   true,
	//	Enable:      true,
	//	GroupName:   config.Nacos.Group,
	//}); err != nil {
	//	return errors.Wrap(err, "fail to RegisterInstance")
	//}

	//if _, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
	//	ServiceName: "speech-tts",
	//	Ip:          util.LocalIP(),
	//	Port:        3012,
	//	Ephemeral:   true,
	//	GroupName:   config.Nacos.Group,
	//}); err != nil {
	//	return errors.Wrap(err, "fail to RegisterInstance")
	//}

	result, err := client.GetService(vo.GetServiceParam{
		ServiceName: "speech-tts.http",
		GroupName:   config.Nacos.Group,
	})
	for _, host := range result.Hosts {
		fmt.Printf("ip:%s, port:%d, serviceName:%s, ClusterName:%s", host.Ip, host.Port, host.ServiceName, host.ClusterName)
	}

	return err
}

func main() {
	fmt.Println(getNacosService(bc.Data))
}

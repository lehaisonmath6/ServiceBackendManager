package transports

import (
	// "github.com/OpenStars/backendclients/go//gen-go/OpenStars/TrustKeys/ObjectStore" //Todo: Fix this
	"TrustKeys/Decentralization/ServiceBackendManager/ObjectStore" //Todo: Fix this
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/OpenStars/thriftpool"
)

var (
	mTObjectStoreServiceBinaryMapPool = thriftpool.NewMapPool(1000, 3600, 3600,
		thriftpool.GetThriftClientCreatorFunc(func(c thrift.TClient) interface{} { return (ObjectStore.NewTObjectServiceClient(c)) }),
		thriftpool.DefaultClose)

	mTObjectStoreServiceCommpactMapPool = thriftpool.NewMapPool(1000, 3600, 3600,
		thriftpool.GetThriftClientCreatorFuncCompactProtocol(func(c thrift.TClient) interface{} { return (ObjectStore.NewTObjectServiceClient(c)) }),
		thriftpool.DefaultClose)
)

func init() {
	fmt.Println("init thrift TObjectService client ")
}

//GetTObjectStoreServiceBinaryClient client by host:port
func GetTObjectStoreServiceBinaryClient(aHost, aPort string) *thriftpool.ThriftSocketClient {
	client, _ := mTObjectStoreServiceBinaryMapPool.Get(aHost, aPort).Get()
	return client
}

//GetTObjectStoreServiceCompactClient get compact client by host:port
func GetTObjectStoreServiceCompactClient(aHost, aPort string) *thriftpool.ThriftSocketClient {
	client, _ := mTObjectStoreServiceCommpactMapPool.Get(aHost, aPort).Get()
	return client
}

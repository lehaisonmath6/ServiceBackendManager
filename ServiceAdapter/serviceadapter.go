package ServiceAdapter

import (
	"TrustKeys/Decentralization/ServiceBackendManager/ObjectInterface"
	"TrustKeys/Decentralization/ServiceBackendManager/ObjectStore"
	"TrustKeys/Decentralization/ServiceBackendManager/transports"
	"context"
	"errors"
	"fmt"
)

var GlobalServiceManager *ServiceAdapter

type ServiceAdapter struct {
	MapService map[string]*ObjectStore.TObjectServiceClient
}

func InitServiceManager() {
	GlobalServiceManager = &ServiceAdapter{
		MapService: make(map[string]*ObjectStore.TObjectServiceClient),
	}
}

func (o *ServiceAdapter) AddService(nameService, host, port string) bool {
	aClient := transports.GetTObjectStoreServiceCompactClient(host, port)
	aClient.BackToPool()
	remoteClient := aClient.Client.(*ObjectStore.TObjectServiceClient)
	o.MapService[nameService] = remoteClient
	return true
}

func (o *ServiceAdapter) PutObject(nameService string, ob ObjectInterface.Serializableif) error {
	aClient, ok := o.MapService[nameService]
	if !ok {
		return errors.New("NOT_FOUND_SERVICE")
	}

	key, value := ob.Serialization()
	if key == "" && value == "" {
		return errors.New("ERROR_SERIALIZATION")
	}
	eCode, err := aClient.PutObject(context.Background(), &ObjectStore.TKObject{Key: key, Value: value})
	if err == nil && eCode == ObjectStore.TErrorCode_EGood {
		return nil
	}
	if err != nil {
		return err
	}
	switch eCode {
	case ObjectStore.TErrorCode_EDataExisted:
		return errors.New("ERROR_DATA_EXISTED")
	default:
		return errors.New("ERROR_UNKNOWN")
	}
}

func (o *ServiceAdapter) GetObject(nameService string, key string, oif ObjectInterface.Serializableif) error {
	aClient, ok := o.MapService[nameService]
	if !ok {
		return errors.New("NOT_FOUND_SERVICE")
	}
	ob, err := aClient.GetObject(context.Background(), key)

	fmt.Println("object is ", ob)
	if err != nil {
		return err
	}
	if ob.Code == ObjectStore.TErrorCode_ENotFound {
		return errors.New("ERROR_NOT_FOUND")
	}
	if ob.Code == ObjectStore.TErrorCode_EUnknown {
		return errors.New("ERROR_UNKNOWN")
	}
	fmt.Println("key is ", ob.Data.Key, "value is ", ob.Data.Value)
	err = oif.Deserialization(ob.Data.Key, ob.Data.Value)
	if err != nil {
		return err
	}
	return nil
}

func (o *ServiceAdapter) RemoveObject(nameService string, key string) error {
	aClient, ok := o.MapService[nameService]
	if !ok {
		return errors.New("NOT_FOUND_SERVICE")
	}
	ob, err := aClient.RemoveObject(context.Background(), key)
	if err != nil {
		return err
	}
	switch ob {
	case ObjectStore.TErrorCode_ENotFound:
		return errors.New("ERROR_NOT_FOUND")
	case ObjectStore.TErrorCode_EGood:
		return nil
	default:
		return errors.New("ERROR_UNKNOWN")
	}
}

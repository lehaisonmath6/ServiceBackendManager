package ObjectInterface

type Serializableif interface {
	Serialization() (key, value string)
	Deserialization(key, value string) error
}

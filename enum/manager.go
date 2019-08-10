package kernel

const E_WARN int = 2
const E_OK int = 1
const E_ERROR int = 3
const E_NORMAL int = 0

type EnumItem struct {
	Value interface{} `json:"value"`
	Label string      `json:"label"`
	Level int         `json:"level"`
}
type EnumManager struct {
	KeyName string
	m       []EnumItem
}

var EnumContainers []*EnumManager = []*EnumManager{}

func NewEnumManager(key string, m []EnumItem) *EnumManager {

	instance := &EnumManager{
		KeyName: key,
		m:       m,
	}
	EnumContainers = append(EnumContainers, instance)
	return instance
}

func (e *EnumManager) ContainsLabel(value string) bool {
	for _, val := range e.m {
		if value == val.Label {
			return true
		}
	}
	return false
}
func (e *EnumManager) GetData() []EnumItem {
	return e.m
}
func (e *EnumManager) ContainsValue(value interface{}) bool {
	for _, item := range e.m {
		if item.Value == value {
			return true
		}
	}
	return false
}

func (e *EnumManager) GetLabel(key interface{}) string {
	for _, item := range e.m {
		if key == item.Value {
			return item.Label
		}
	}
	return "[ERROR]:LABEL_NOT_FOUND"

}

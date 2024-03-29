package msgqueue

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type DynamicEventMapper struct {
	typeMap map[string]reflect.Type
}

func NewDynamicEventMapper() EventMapper {
	return &DynamicEventMapper{
		typeMap: make(map[string]reflect.Type),
	}
}

func (e *DynamicEventMapper) MapEvent(eventName string, serialized interface{}) (Event, error) {
	typ, ok := e.typeMap[eventName]
	if !ok {
		return nil, fmt.Errorf("no mapping configured for event %s", eventName)
	}
	// reflect
	instance := reflect.New(typ)
	iface := instance.Interface()
	// from type typ get Event interface
	event, ok := iface.(Event)
	if !ok {
		return nil, fmt.Errorf("type %T does not implement the Event interface", iface)
	}
	// serialized interface{} cast to type
	switch s := serialized.(type) {
	case []byte:
		// bytes, event
		err := json.Unmarshal(s, event)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
		}
	default:
		cfg := mapstructure.DecoderConfig{
			Result:  event,
			TagName: "json",
		}
		// page 335
		// get decoder
		dec, err := mapstructure.NewDecoder(&cfg)
		if err != nil {
			return nil, fmt.Errorf("could not initialize decoder for event %s: %s", eventName, err)
		}
		// decode
		err = dec.Decode(s)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
		}
	}

	return event, nil
}

func (e *DynamicEventMapper) RegisterMapping(eventType reflect.Type) error {
	instance := reflect.New(eventType).Interface()
	event, ok := instance.(Event)
	if !ok {
		return fmt.Errorf("type %T does not implement the Event interface", instance)
	}

	e.typeMap[event.EventName()] = eventType
	return nil
}

package matrixcomm

import "errors"

type Filter struct {
	AccountData FilterPart `json:"account_data,omitempty"`
	EventFields []string   `json:"event_fields,omitempty"`
	EventFormat string     `json:"event_format,omitempty"`
	Presence    FilterPart `json:"presence,omitempty"`
	Room        RoomFilter `json:"room,omitempty"`
}

type RoomFilter struct {
	AccountData  FilterPart `json:"account_data,omitempty"`
	Ephemeral    FilterPart `json:"ephemeral,omitempty"`
	IncludeLeave bool       `json:"include_leave,omitempty"`
	NotRooms     []string   `json:"not_rooms,omitempty"`
	Rooms        []string   `json:"rooms,omitempty"`
	State        FilterPart `json:"state,omitempty"`
	Timeline     FilterPart `json:"timeline,omitempty"`
}

type FilterPart struct {
	NotRooms    []string `json:"not_rooms,omitempty"`
	Rooms       []string `json:"rooms,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	NotSenders  []string `json:"not_senders,omitempty"`
	NotTypes    []string `json:"not_types,omitempty"`
	Senders     []string `json:"senders,omitempty"`
	Types       []string `json:"types,omitempty"`
	ContainsURL *bool    `json:"contains_url,omitempty"`
}

func (filter *Filter) Validate() error {
	if filter.EventFormat != "client" && filter.EventFormat != "federation" {
		return errors.New("Bad event_format value. Must be one of [\"client\", \"federation\"]")
	}
	return nil
}

func DefaultFilter() Filter {
	return Filter{
		AccountData: DefaultFilterPart(),
		EventFields: nil,
		EventFormat: "client",
		Presence:    DefaultFilterPart(),
		Room: RoomFilter{
			AccountData:  DefaultFilterPart(),
			Ephemeral:    DefaultFilterPart(),
			IncludeLeave: false,
			NotRooms:     nil,
			Rooms:        nil,
			State:        DefaultFilterPart(),
			Timeline:     DefaultFilterPart(),
		},
	}
}

func DefaultFilterPart() FilterPart {
	return FilterPart{
		NotRooms:   nil,
		Rooms:      nil,
		Limit:      20,
		NotSenders: nil,
		NotTypes:   nil,
		Senders:    nil,
		Types:      nil,
	}
}

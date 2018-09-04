package matrixcomm

type ReqAccountData struct {
	Addresshex string
	Roomid string
}

type ReqPresenceUser struct {
	Presence	string 		`json:"presence"`
	StatusMsg    string 	`json:"status_msg"`
}

type ReqPresenceList struct {
	Drop 	[]string	`json:"drop"`
	Invite 	[]string 	`json:"invite"`
}

type ReqRegister struct {
	Username                 string      `json:"username,omitempty"`
	BindEmail                bool        `json:"bind_email,omitempty"`
	Password                 string      `json:"password,omitempty"`
	DeviceID                 string      `json:"device_id,omitempty"`
	InitialDeviceDisplayName string      `json:"initial_device_display_name"`
	Auth                     AuthDict	 `json:"auth,omitempty"`
	Type                     string      `json:"type,omitempty"`
	Admin                    bool        `json:"admin"`
}

type AuthDict struct {
	Type    string         				`json:"type"`
	Session string                      `json:"session"`
	Mac     []byte 						`json:"mac"`
	Response string 					`json:"response"`
}

type reqPresenceList struct {
	Drop	string `json:"drop"`
	Invite  string `json:"invite"`
}

type ReqLogin struct {
	Type                     string `json:"type"`
	Password                 string `json:"password,omitempty"`
	Medium                   string `json:"medium,omitempty"`
	User                     string `json:"user,omitempty"`
	Address                  string `json:"address,omitempty"`
	Token                    string `json:"token,omitempty"`
	DeviceID                 string `json:"device_id,omitempty"`
	InitialDeviceDisplayName string `json:"initial_device_display_name,omitempty"`
}

type ReqCreateRoom struct {
	Visibility      string                 `json:"visibility,omitempty"`
	RoomAliasName   string                 `json:"room_alias_name,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Topic           string                 `json:"topic,omitempty"`
	Invite          []string               `json:"invite,omitempty"`
	Invite3PID      []ReqInvite3PID        `json:"invite_3pid,omitempty"`
	CreationContent map[string]interface{} `json:"creation_content,omitempty"`
	InitialState    []Event                `json:"initial_state,omitempty"`
	Preset          string                 `json:"preset,omitempty"`
	IsDirect        bool                   `json:"is_direct,omitempty"`
}

type ReqRedact struct {
	Reason string `json:"reason,omitempty"`
}

type ReqInvite3PID struct {
	IDServer string `json:"id_server"`
	Medium   string `json:"medium"`
	Address  string `json:"address"`
}

type ReqInviteUser struct {
	UserID string `json:"user_id"`
}

type ReqKickUser struct {
	Reason string `json:"reason,omitempty"`
	UserID string `json:"user_id"`
}

type ReqBanUser struct {
	Reason string `json:"reason,omitempty"`
	UserID string `json:"user_id"`
}

type ReqUnbanUser struct {
	UserID string `json:"user_id"`
}

type ReqTyping struct {
	Typing  bool  `json:"typing"`
	Timeout int64 `json:"timeout"`
}

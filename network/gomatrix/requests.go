package gomatrix

//ReqUserSearch is the JSON request for UserSearch
type ReqUserSearch struct {
	Limit      int    `json:"limit,omitempty"`
	SearchTerm string `json:"search_term"`
}

//ReqAccountData is the JSON request for AccountData
type ReqAccountData struct {
	//Addresshex string
	//Roomid     []string
	AccountData map[string]interface{} `json:"account_data"`
}

//ReqPresenceUser is the JSON request for PresenceUser
type ReqPresenceUser struct {
	Presence  string `json:"presence"`
	StatusMsg string `json:"status_msg"`
}

//ReqPresenceList is the JSON request for PresenceList
type ReqPresenceList struct {
	Drop   []string `json:"drop"`
	Invite []string `json:"invite"`
}

//ReqRegister is the JSON request for Register
type ReqRegister struct {
	Username                 string   `json:"username,omitempty"`
	BindEmail                bool     `json:"bind_email,omitempty"`
	Password                 string   `json:"password,omitempty"`
	DeviceID                 string   `json:"device_id,omitempty"`
	InitialDeviceDisplayName string   `json:"initial_device_display_name"`
	Auth                     AuthDict `json:"auth,omitempty"`
	Type                     string   `json:"type,omitempty"`
	Admin                    bool     `json:"admin"`
}

//AuthDict is the JSON request for AuthDict
type AuthDict struct {
	Type     string `json:"type"`
	Session  string `json:"session"`
	Mac      []byte `json:"mac"`
	Response string `json:"response"`
}

//ReqLogin is the JSON request for Login
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

//ReqCreateRoom is the JSON request for CreateRoom
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

//ReqRedact is the JSON request for Redact
type ReqRedact struct {
	Reason string `json:"reason,omitempty"`
}

//ReqInvite3PID is the JSON request for Invite3PID
type ReqInvite3PID struct {
	IDServer string `json:"id_server"`
	Medium   string `json:"medium"`
	Address  string `json:"address"`
}

//ReqInviteUser is the JSON request for InviteUser
type ReqInviteUser struct {
	UserID string `json:"user_id"`
}

//ReqKickUser is the JSON request for KickUser
type ReqKickUser struct {
	Reason string `json:"reason,omitempty"`
	UserID string `json:"user_id"`
}

//ReqBanUser is the JSON request for BanUser
type ReqBanUser struct {
	Reason string `json:"reason,omitempty"`
	UserID string `json:"user_id"`
}

//ReqUnbanUser is the JSON request for UnbanUser
type ReqUnbanUser struct {
	UserID string `json:"user_id"`
}

// ReqTyping is the JSON request for Typing
type ReqTyping struct {
	Typing  bool  `json:"typing"`
	Timeout int64 `json:"timeout"`
}

package matrixcomm

// RespAccountData is the JSON for AccountData
type RespAccountData struct {
	Events []Event `json:"events"`
}

// RespUserSearch is the JSON for UserSearch
type RespUserSearch struct {
	Limited bool     `json:"limited"`
	Results []UserInfo `json:"results"`
}

// UserInfo is a polular JSON for UserSearch
type UserInfo struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	UserID      string `json:"user_id"`
}

// RespError is the standard JSON error response from Homeservers. It also implements the Golang "error" interface.
type RespError struct {
	ErrCode string `json:"errcode"`
	Err     string `json:"error"`
}

// Error returns the errcode and error message.
func (e RespError) Error() string {
	return e.ErrCode + ": " + e.Err
}

// RespCreateFilter is the JSON response for CreateFilter
type RespCreateFilter struct {
	FilterID string `json:"filter_id"`
}

// RespVersions is the JSON response for get c-s-api version
type RespVersions struct {
	Versions []string `json:"versions"`
}

// RespJoinRoom is the JSON response for JoinRoom
type RespJoinRoom struct {
	RoomID string `json:"room_id"`
}

// RespLeaveRoom is the JSON response for LeaveRoom
type RespLeaveRoom struct{}

// RespForgetRoom is the JSON response for ForgetRoom
type RespForgetRoom struct{}

// RespInviteUser is the JSON response for InviteUser
type RespInviteUser struct{}

// RespKickUser is the JSON response for KickUser
type RespKickUser struct{}

// RespBanUser is the JSON response for BanUser
type RespBanUser struct{}

// RespUnbanUser is the JSON response for UnbanUser
type RespUnbanUser struct{}

// RespTyping is the JSON response for Typing
type RespTyping struct{}

// RespJoinedRooms is the JSON response for JoinedRooms
type RespJoinedRooms struct {
	JoinedRooms []string `json:"joined_rooms"`
}

// RespJoinedMembers is the JSON response for JoinedMembers
type RespJoinedMembers struct {
	Joined map[string]struct {
		DisplayName *string `json:"display_name"`
		AvatarURL   *string `json:"avatar_url"`
	} `json:"joined"`
}

// RespMessages is the JSON response for Messages
type RespMessages struct {
	Start string  `json:"start"`
	Chunk []Event `json:"chunk"`
	End   string  `json:"end"`
}

// RespSendEvent is the JSON response for SendEvent
type RespSendEvent struct {
	EventID string `json:"event_id"`
}

// RespMediaUpload is the JSON response for MediaUpload
type RespMediaUpload struct {
	ContentURI string `json:"content_uri"`
}

// RespUserInteractive is the JSON response for UserInteractive
type RespUserInteractive struct {
	Flows []struct {
		Stages []string `json:"stages"`
	} `json:"flows"`
	Params    map[string]interface{} `json:"params"`
	Session   string                 `json:"string"`
	Completed []string               `json:"completed"`
	ErrCode   string                 `json:"errcode"`
	Error     string                 `json:"error"`
}

// HasSingleStageFlow returns true if there exists at least 1 Flow with a single stage of stageName.
func (r RespUserInteractive) HasSingleStageFlow(stageName string) bool {
	for _, f := range r.Flows {
		if len(f.Stages) == 1 && f.Stages[0] == stageName {
			return true
		}
	}
	return false
}

//RespUserDisplayName is the JSON response for UserDisplayName
type RespUserDisplayName struct {
	DisplayName string `json:"displayname"`
}

// RespRegister is the JSON response for Register
type RespRegister struct {
	AccessToken  string `json:"access_token"`
	DeviceID     string `json:"device_id"`
	HomeServer   string `json:"home_server"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

// RespLogin is the JSON response for Login
type RespLogin struct {
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
}

// RespPresenceUser is the JSON response for PresenceUser
type RespPresenceUser struct {
	LastActiveAgo   int    `json:"last_active_ago"`
	StatusMsg       string `json:"status_msg"`
	CurrentlyActive bool   `json:"currently_active"`
	UserID          string `json:"user_id"`
	Presence        string `json:"presence"`
}

//RespPresenceList is the JSON response for PresenceList
type RespPresenceList struct {
	Accepted        interface{} `json:"accepted"`
	LastActiveAgo   int         `json:"last_active_ago"`
	StatusMsg       string      `json:"status_msg"`
	CurrentlyActive bool        `json:"currently_active"`
	UserID          string      `json:"user_id"`
	Presence        string      `json:"presence"`
}

// RespLogout is the JSON response for Logout
type RespLogout struct{}

// RespCreateRoom is the JSON response for CreateRoom
type RespCreateRoom struct {
	RoomID string `json:"room_id"`
}

// RespSync is the JSON response of sync
type RespSync struct {
	NextBatch   string `json:"next_batch"`
	AccountData struct {
		Events []Event `json:"events"`
	} `json:"account_data"`
	Presence struct {
		Events []Event `json:"events"`
	} `json:"presence"`
	Rooms struct {
		Leave map[string]struct {
			State struct {
				Events []Event `json:"events"`
			} `json:"state"`
			Timeline struct {
				Events    []Event `json:"events"`
				Limited   bool    `json:"limited"`
				PrevBatch string  `json:"prev_batch"`
			} `json:"timeline"`
		} `json:"leave"`
		Join map[string]struct {
			State struct {
				Events []Event `json:"events"`
			} `json:"state"`
			Timeline struct {
				Events    []Event `json:"events"`
				Limited   bool    `json:"limited"`
				PrevBatch string  `json:"prev_batch"`
			} `json:"timeline"`
		} `json:"join"`
		Invite map[string]struct {
			State struct {
				Events []Event
			} `json:"invite_state"`
		} `json:"invite"`
	} `json:"rooms"`
}

// RespTurnServer is the JSON response of turn server
type RespTurnServer struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	TTL      int      `json:"ttl"`
	URIs     []string `json:"uris"`
}

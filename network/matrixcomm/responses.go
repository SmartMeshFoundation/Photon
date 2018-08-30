package matrixcomm

type RespError struct {
	ErrCode string `json:"errcode"`
	Err     string `json:"error"`
}

func (e RespError) Error() string {
	return e.ErrCode + ": " + e.Err
}



type RespCreateFilter struct {
	FilterID string `json:"filter_id"`
}

type RespVersions struct {
	Versions []string `json:"versions"`
}

type RespJoinRoom struct {
	RoomID string `json:"room_id"`
}

type RespLeaveRoom struct{}

type RespForgetRoom struct{}

type RespInviteUser struct{}

type RespKickUser struct{}

type RespBanUser struct{}

type RespUnbanUser struct{}

type RespTyping struct{}

type RespJoinedRooms struct {
	JoinedRooms []string `json:"joined_rooms"`
}

type RespJoinedMembers struct {
	Joined map[string]struct {
		DisplayName *string `json:"display_name"`
		AvatarURL   *string `json:"avatar_url"`
	} `json:"joined"`
}

type RespMessages struct {
	Start string  `json:"start"`
	Chunk []Event `json:"chunk"`
	End   string  `json:"end"`
}

type RespSendEvent struct {
	EventID string `json:"event_id"`
}

type RespMediaUpload struct {
	ContentURI string `json:"content_uri"`
}

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

func (r RespUserInteractive) HasSingleStageFlow(stageName string) bool {
	for _, f := range r.Flows {
		if len(f.Stages) == 1 && f.Stages[0] == stageName {
			return true
		}
	}
	return false
}

type RespUserDisplayName struct {
	DisplayName string `json:"displayname"`
}

type RespRegister struct {
	AccessToken  string `json:"access_token"`
	DeviceID     string `json:"device_id"`
	HomeServer   string `json:"home_server"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

type RespLogin struct {
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
}


type RespPresenceUser struct {
	LastActiveAgo int 		`json:"last_active_ago"`
	StatusMsg    string 	`json:"status_msg"`
	CurrentlyActive  bool 	`json:"currently_active"`
	UserID      string 		`json:"user_id"`
	Presence	string 		`json:"presence"`
}

type RespPresenceList struct {
	Accepted	interface{} `json:"accepted"`
	LastActiveAgo int 		`json:"last_active_ago"`
	StatusMsg    string 	`json:"status_msg"`
	CurrentlyActive  bool 	`json:"currently_active"`
	UserID      string 		`json:"user_id"`
	Presence	string 		`json:"presence"`

} 

type RespLogout struct{}

type RespCreateRoom struct {
	RoomID string `json:"room_id"`
}

/*type Respwhois struct{
	Devices struct{
		Teapot struct{
			map[string]
		}`json:"teapot"`
	} `json:"next_batch"`
	UserId string `json:"user_id"`
}*/

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

type RespTurnServer struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	TTL      int      `json:"ttl"`
	URIs     []string `json:"uris"`
}

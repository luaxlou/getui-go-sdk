package getui

// AuthDTO 认证请求DTO
type AuthDTO struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	AppKey    string `json:"appkey"`
}

// PushDTO 推送请求DTO
type PushDTO struct {
	RequestID   string       `json:"request_id"`
	TaskName    string       `json:"task_name,omitempty"`
	GroupName   string       `json:"group_name,omitempty"`
	Settings    *Settings    `json:"settings,omitempty"`
	Audience    interface{}  `json:"audience"`
	PushMessage *PushMessage `json:"push_message"`
	PushChannel *PushChannel `json:"push_channel,omitempty"`
}

// PushBatchDTO 批量推送请求DTO
type PushBatchDTO struct {
	RequestID   string       `json:"request_id"`
	TaskName    string       `json:"task_name,omitempty"`
	GroupName   string       `json:"group_name,omitempty"`
	Settings    *Settings    `json:"settings,omitempty"`
	Audience    interface{}  `json:"audience"`
	PushMessage *PushMessage `json:"push_message"`
	PushChannel *PushChannel `json:"push_channel,omitempty"`
}

// AudienceDTO 受众DTO
type AudienceDTO struct {
	RequestID string      `json:"request_id"`
	TaskName  string      `json:"task_name,omitempty"`
	GroupName string      `json:"group_name,omitempty"`
	Settings  *Settings   `json:"settings,omitempty"`
	Audience  interface{} `json:"audience"`
}

// Audience 受众
type Audience struct {
	CIDs   []string `json:"cid,omitempty"`
	Alias  []string `json:"alias,omitempty"`
	Tag    []string `json:"tag,omitempty"`
	All    string   `json:"all,omitempty"`
	FileID string   `json:"file_id,omitempty"`
}

// PushMessage 推送消息
type PushMessage struct {
	NetworkType  int           `json:"network_type,omitempty"`
	Duration     string        `json:"duration,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
	Transmission string        `json:"transmission,omitempty"`
	Revoke       *RevokeBean   `json:"revoke,omitempty"`
}

// Notification 通知消息
type Notification struct {
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	ClickType    string            `json:"click_type"`
	URL          string            `json:"url,omitempty"`
	Intent       string            `json:"intent,omitempty"`
	Payload      string            `json:"payload,omitempty"`
	Badge        int               `json:"badge,omitempty"`
	Ring         int               `json:"ring,omitempty"`
	Buzz         int               `json:"buzz,omitempty"`
	Logo         string            `json:"logo,omitempty"`
	LogoURL      string            `json:"logo_url,omitempty"`
	ChannelID    string            `json:"channel_id,omitempty"`
	ChannelName  string            `json:"channel_name,omitempty"`
	ChannelLevel int               `json:"channel_level,omitempty"`
	MultiPkg     bool              `json:"multi_pkg,omitempty"`
	NotifyID     int               `json:"notify_id,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
}

// RevokeBean 撤回消息
type RevokeBean struct {
	OldTaskID string `json:"old_task_id"`
}

// PushChannel 推送通道
type PushChannel struct {
	IOS     *IOSDTO     `json:"ios,omitempty"`
	Android *AndroidDTO `json:"android,omitempty"`
	Harmony *HarmonyDTO `json:"harmony,omitempty"`
}

// IOSDTO iOS推送参数
type IOSDTO struct {
	Type             string `json:"type,omitempty"`
	APNS             *APNS  `json:"apns,omitempty"`
	APNSCollapseID   string `json:"apns_collapse_id,omitempty"`
	AutoBadge        string `json:"auto_badge,omitempty"`
	MutableContent   int    `json:"mutable_content,omitempty"`
	ContentAvailable int    `json:"content_available,omitempty"`
	Category         string `json:"category,omitempty"`
	Alert            *Alert `json:"alert,omitempty"`
}

// APNS APNS配置
type APNS struct {
	Alert            *Alert            `json:"alert,omitempty"`
	Badge            int               `json:"badge,omitempty"`
	Sound            string            `json:"sound,omitempty"`
	ContentAvailable int               `json:"content_available,omitempty"`
	MutableContent   int               `json:"mutable_content,omitempty"`
	Category         string            `json:"category,omitempty"`
	CustomData       map[string]string `json:"custom_data,omitempty"`
}

// Alert iOS通知内容
type Alert struct {
	Title    string   `json:"title,omitempty"`
	Body     string   `json:"body,omitempty"`
	Subtitle string   `json:"subtitle,omitempty"`
	Action   string   `json:"action,omitempty"`
	Args     []string `json:"args,omitempty"`
}

// AndroidDTO Android推送参数
type AndroidDTO struct {
	UPS *UPS `json:"ups,omitempty"`
}

// UPS 统一推送服务
type UPS struct {
	Notification *ThirdNotification `json:"notification,omitempty"`
	Options      map[string]string  `json:"options,omitempty"`
	Transmission string             `json:"transmission,omitempty"`
}

// ThirdNotification 第三方通知
type ThirdNotification struct {
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	ClickType    string            `json:"click_type"`
	URL          string            `json:"url,omitempty"`
	Intent       string            `json:"intent,omitempty"`
	Payload      string            `json:"payload,omitempty"`
	NotifyID     string            `json:"notify_id,omitempty"`
	ChannelID    string            `json:"channel_id,omitempty"`
	ChannelName  string            `json:"channel_name,omitempty"`
	ChannelLevel int               `json:"channel_level,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
}

// HarmonyDTO 鸿蒙推送参数
type HarmonyDTO struct {
	Notification *HarmonyNotification `json:"notification,omitempty"`
}

// HarmonyNotification 鸿蒙通知
type HarmonyNotification struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Category  string `json:"category,omitempty"`
	ClickType string `json:"click_type,omitempty"`
	Want      string `json:"want,omitempty"`
}

// Settings 推送设置
type Settings struct {
	TTL          int       `json:"ttl,omitempty"`
	Strategy     *Strategy `json:"strategy,omitempty"`
	Speed        int       `json:"speed,omitempty"`
	ScheduleTime string    `json:"schedule_time,omitempty"`
}

// Strategy 推送策略
type Strategy struct {
	Default int `json:"default,omitempty"`
	IOS     int `json:"ios,omitempty"`
	St      int `json:"st,omitempty"`
	Hw      int `json:"hw,omitempty"`
	Xm      int `json:"xm,omitempty"`
	Vv      int `json:"vv,omitempty"`
	Op      int `json:"op,omitempty"`
	Fcm     int `json:"fcm,omitempty"`
}

// TaskIDDTO 任务ID响应
type TaskIDDTO struct {
	TaskID string `json:"task_id"`
}

// ScheduleTaskDTO 定时任务
type ScheduleTaskDTO struct {
	TaskID       string `json:"task_id"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
	ScheduleTime string `json:"schedule_time"`
}

// CidStatusDTO CID状态
type CidStatusDTO struct {
	CID    string `json:"cid"`
	Status string `json:"status"`
}

// StatisticDTO 统计数据
type StatisticDTO struct {
	TaskID       string `json:"task_id"`
	SendCount    int    `json:"send_count"`
	ReceiveCount int    `json:"receive_count"`
	DisplayCount int    `json:"display_count"`
	ClickCount   int    `json:"click_count"`
}

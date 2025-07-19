package getui

import (
	"fmt"
	"time"
)

// StatisticAPI 统计分析API接口
type StatisticAPI struct {
	client *Client
}

// QueryPushResultByTaskIDs 根据任务ID查询推送结果
func (api *StatisticAPI) QueryPushResultByTaskIDs(taskIDs []string) (*ApiResult, error) {
	if len(taskIDs) == 0 {
		return nil, fmt.Errorf("task_ids cannot be empty")
	}

	requestBody := map[string]interface{}{
		"task_id_list": taskIDs,
	}

	return api.client.DoRequest("POST", "/report/push/result", requestBody)
}

// QueryPushResultByDate 根据日期查询推送结果
func (api *StatisticAPI) QueryPushResultByDate(date string) (*ApiResult, error) {
	if date == "" {
		// 默认查询今天的日期
		date = time.Now().Format("2006-01-02")
	}

	url := fmt.Sprintf("/report/push/date/%s", date)
	return api.client.DoRequest("GET", url, nil)
}

// QueryPushResultByTaskID 根据单个任务ID查询推送结果
func (api *StatisticAPI) QueryPushResultByTaskID(taskID string) (*ApiResult, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id cannot be empty")
	}

	url := fmt.Sprintf("/report/push/task/%s", taskID)
	return api.client.DoRequest("GET", url, nil)
}

// QueryUserData 查询用户数据
func (api *StatisticAPI) QueryUserData(date string) (*ApiResult, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	url := fmt.Sprintf("/report/user/date/%s", date)
	return api.client.DoRequest("GET", url, nil)
}

// QueryPerformanceData 查询性能数据
func (api *StatisticAPI) QueryPerformanceData(date string) (*ApiResult, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	url := fmt.Sprintf("/report/performance/date/%s", date)
	return api.client.DoRequest("GET", url, nil)
}

// QueryOnlineUserCount 查询在线用户数
func (api *StatisticAPI) QueryOnlineUserCount() (*ApiResult, error) {
	return api.client.DoRequest("GET", "/report/online_user", nil)
}

// QueryAppData 查询应用数据
func (api *StatisticAPI) QueryAppData(date string) (*ApiResult, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	url := fmt.Sprintf("/report/app/date/%s", date)
	return api.client.DoRequest("GET", url, nil)
}

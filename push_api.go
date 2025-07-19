package getui

import (
	"fmt"
)

// PushAPI 推送API接口
type PushAPI struct {
	client *Client
}

// PushToSingleByCID 根据CID单推
func (api *PushAPI) PushToSingleByCID(pushDTO *PushDTO) (*ApiResult, error) {
	if err := api.validatePushDTO(pushDTO); err != nil {
		return nil, err
	}

	// 确保RequestID不为空
	if pushDTO.RequestID == "" {
		pushDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/single/cid", pushDTO)
}

// PushToSingleByAlias 根据别名单推
func (api *PushAPI) PushToSingleByAlias(pushDTO *PushDTO) (*ApiResult, error) {
	if err := api.validatePushDTO(pushDTO); err != nil {
		return nil, err
	}

	if pushDTO.RequestID == "" {
		pushDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/single/alias", pushDTO)
}

// PushBatchByCID 根据CID批量推送
func (api *PushAPI) PushBatchByCID(batchDTO *PushBatchDTO) (*ApiResult, error) {
	if err := api.validatePushBatchDTO(batchDTO); err != nil {
		return nil, err
	}

	if batchDTO.RequestID == "" {
		batchDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/single/batch/cid", batchDTO)
}

// PushBatchByAlias 根据别名批量推送
func (api *PushAPI) PushBatchByAlias(batchDTO *PushBatchDTO) (*ApiResult, error) {
	if err := api.validatePushBatchDTO(batchDTO); err != nil {
		return nil, err
	}

	if batchDTO.RequestID == "" {
		batchDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/single/batch/alias", batchDTO)
}

// PushAll 群推
func (api *PushAPI) PushAll(pushDTO *PushDTO) (*ApiResult, error) {
	if err := api.validatePushDTO(pushDTO); err != nil {
		return nil, err
	}

	if pushDTO.RequestID == "" {
		pushDTO.RequestID = api.client.GenerateRequestID()
	}

	// 群推时Audience设置为"all"
	pushDTO.Audience = "all"

	return api.client.DoRequest("POST", "/push/all", pushDTO)
}

// PushByTag 根据标签推送
func (api *PushAPI) PushByTag(pushDTO *PushDTO) (*ApiResult, error) {
	if err := api.validatePushDTO(pushDTO); err != nil {
		return nil, err
	}

	if pushDTO.RequestID == "" {
		pushDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/tag", pushDTO)
}

// PushByFastCustomTag 使用标签快速推送
func (api *PushAPI) PushByFastCustomTag(pushDTO *PushDTO) (*ApiResult, error) {
	if err := api.validatePushDTO(pushDTO); err != nil {
		return nil, err
	}

	if pushDTO.RequestID == "" {
		pushDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/fast_custom_tag", pushDTO)
}

// CreateMsg 创建消息体
func (api *PushAPI) CreateMsg(pushDTO *PushDTO) (*ApiResult, error) {
	if err := api.validatePushDTO(pushDTO); err != nil {
		return nil, err
	}

	if pushDTO.RequestID == "" {
		pushDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/list/message", pushDTO)
}

// PushListByCID 根据CID列表推送
func (api *PushAPI) PushListByCID(audienceDTO *AudienceDTO) (*ApiResult, error) {
	if err := api.validateAudienceDTO(audienceDTO); err != nil {
		return nil, err
	}

	if audienceDTO.RequestID == "" {
		audienceDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/list/cid", audienceDTO)
}

// PushListByAlias 根据别名列表推送
func (api *PushAPI) PushListByAlias(audienceDTO *AudienceDTO) (*ApiResult, error) {
	if err := api.validateAudienceDTO(audienceDTO); err != nil {
		return nil, err
	}

	if audienceDTO.RequestID == "" {
		audienceDTO.RequestID = api.client.GenerateRequestID()
	}

	return api.client.DoRequest("POST", "/push/list/alias", audienceDTO)
}

// StopPush 停止推送任务
func (api *PushAPI) StopPush(taskID string) (*ApiResult, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id cannot be empty")
	}

	return api.client.DoRequest("DELETE", fmt.Sprintf("/task/%s", taskID), nil)
}

// QueryScheduleTask 查询定时任务
func (api *PushAPI) QueryScheduleTask(taskID string) (*ApiResult, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id cannot be empty")
	}

	return api.client.DoRequest("GET", fmt.Sprintf("/task/schedule/%s", taskID), nil)
}

// DeleteScheduleTask 删除定时任务
func (api *PushAPI) DeleteScheduleTask(taskID string) (*ApiResult, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id cannot be empty")
	}

	return api.client.DoRequest("DELETE", fmt.Sprintf("/task/schedule/%s", taskID), nil)
}

// validatePushDTO 验证推送DTO
func (api *PushAPI) validatePushDTO(pushDTO *PushDTO) error {
	if pushDTO == nil {
		return fmt.Errorf("push_dto cannot be nil")
	}

	if pushDTO.RequestID != "" && (len(pushDTO.RequestID) < 10 || len(pushDTO.RequestID) > 32) {
		return ErrInvalidRequestID
	}

	if pushDTO.Audience == nil {
		return ErrEmptyAudience
	}

	if pushDTO.PushMessage == nil {
		return ErrEmptyPushMessage
	}

	return nil
}

// validatePushBatchDTO 验证批量推送DTO
func (api *PushAPI) validatePushBatchDTO(batchDTO *PushBatchDTO) error {
	if batchDTO == nil {
		return fmt.Errorf("batch_dto cannot be nil")
	}

	if batchDTO.RequestID != "" && (len(batchDTO.RequestID) < 10 || len(batchDTO.RequestID) > 32) {
		return ErrInvalidRequestID
	}

	if batchDTO.Audience == nil {
		return ErrEmptyAudience
	}

	if batchDTO.PushMessage == nil {
		return ErrEmptyPushMessage
	}

	return nil
}

// validateAudienceDTO 验证受众DTO
func (api *PushAPI) validateAudienceDTO(audienceDTO *AudienceDTO) error {
	if audienceDTO == nil {
		return fmt.Errorf("audience_dto cannot be nil")
	}

	if audienceDTO.RequestID != "" && (len(audienceDTO.RequestID) < 10 || len(audienceDTO.RequestID) > 32) {
		return ErrInvalidRequestID
	}

	if audienceDTO.Audience == nil {
		return ErrEmptyAudience
	}

	return nil
}

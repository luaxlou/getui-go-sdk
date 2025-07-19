package getui

import (
	"fmt"
)

// UserAPI 用户管理API接口
type UserAPI struct {
	client *Client
}

// QueryUserStatus 查询用户状态
func (api *UserAPI) QueryUserStatus(cids []string) (*ApiResult, error) {
	if len(cids) == 0 {
		return nil, fmt.Errorf("cids cannot be empty")
	}

	requestBody := map[string]interface{}{
		"cid": cids,
	}

	return api.client.DoRequest("POST", "/user/status", requestBody)
}

// QueryAliasByCID 根据CID查询别名
func (api *UserAPI) QueryAliasByCID(cid string) (*ApiResult, error) {
	if cid == "" {
		return nil, ErrInvalidCID
	}

	return api.client.DoRequest("GET", fmt.Sprintf("/user/alias/%s", cid), nil)
}

// QueryCIDByAlias 根据别名查询CID
func (api *UserAPI) QueryCIDByAlias(alias string) (*ApiResult, error) {
	if alias == "" {
		return nil, ErrInvalidAlias
	}

	return api.client.DoRequest("GET", fmt.Sprintf("/user/cid/%s", alias), nil)
}

// BindAlias 绑定别名
func (api *UserAPI) BindAlias(alias string, cid string) (*ApiResult, error) {
	if alias == "" {
		return nil, ErrInvalidAlias
	}
	if cid == "" {
		return nil, ErrInvalidCID
	}

	requestBody := map[string]string{
		"alias": alias,
		"cid":   cid,
	}

	return api.client.DoRequest("POST", "/user/alias", requestBody)
}

// UnbindAlias 解绑别名
func (api *UserAPI) UnbindAlias(alias string, cid string) (*ApiResult, error) {
	if alias == "" {
		return nil, ErrInvalidAlias
	}
	if cid == "" {
		return nil, ErrInvalidCID
	}

	requestBody := map[string]string{
		"alias": alias,
		"cid":   cid,
	}

	return api.client.DoRequest("DELETE", "/user/alias", requestBody)
}

// BindAliasBatch 批量绑定别名
func (api *UserAPI) BindAliasBatch(aliasCidList []map[string]string) (*ApiResult, error) {
	if len(aliasCidList) == 0 {
		return nil, fmt.Errorf("alias_cid_list cannot be empty")
	}

	requestBody := map[string]interface{}{
		"data_list": aliasCidList,
	}

	return api.client.DoRequest("POST", "/user/alias/batch", requestBody)
}

// UnbindAliasBatch 批量解绑别名
func (api *UserAPI) UnbindAliasBatch(aliasCidList []map[string]string) (*ApiResult, error) {
	if len(aliasCidList) == 0 {
		return nil, fmt.Errorf("alias_cid_list cannot be empty")
	}

	requestBody := map[string]interface{}{
		"data_list": aliasCidList,
	}

	return api.client.DoRequest("DELETE", "/user/alias/batch", requestBody)
}

// QueryUserDetail 查询用户详情
func (api *UserAPI) QueryUserDetail(cid string) (*ApiResult, error) {
	if cid == "" {
		return nil, ErrInvalidCID
	}

	return api.client.DoRequest("GET", fmt.Sprintf("/user/detail/%s", cid), nil)
}

// SetUserTag 设置用户标签
func (api *UserAPI) SetUserTag(cid string, tags []string) (*ApiResult, error) {
	if cid == "" {
		return nil, ErrInvalidCID
	}

	requestBody := map[string]interface{}{
		"cid":  cid,
		"tags": tags,
	}

	return api.client.DoRequest("POST", "/user/tag", requestBody)
}

// GetUserTag 获取用户标签
func (api *UserAPI) GetUserTag(cid string) (*ApiResult, error) {
	if cid == "" {
		return nil, ErrInvalidCID
	}

	return api.client.DoRequest("GET", fmt.Sprintf("/user/tag/%s", cid), nil)
}

// DeleteUserTag 删除用户标签
func (api *UserAPI) DeleteUserTag(cid string, tags []string) (*ApiResult, error) {
	if cid == "" {
		return nil, ErrInvalidCID
	}

	requestBody := map[string]interface{}{
		"cid":  cid,
		"tags": tags,
	}

	return api.client.DoRequest("DELETE", "/user/tag", requestBody)
}

// GetUserCount 获取用户数量
func (api *UserAPI) GetUserCount() (*ApiResult, error) {
	return api.client.DoRequest("GET", "/user/count", nil)
}

// GetUserList 获取用户列表
func (api *UserAPI) GetUserList(page int, size int) (*ApiResult, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 1000 {
		size = 100
	}

	url := fmt.Sprintf("/user/list?page=%d&size=%d", page, size)
	return api.client.DoRequest("GET", url, nil)
}

package getui

import (
	"encoding/json"
)

// ApiResult API响应结果
type ApiResult struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// IsSuccess 判断请求是否成功
func (r *ApiResult) IsSuccess() bool {
	return r.Code == 0
}

// GetData 获取响应数据
func (r *ApiResult) GetData() json.RawMessage {
	return r.Data
}

// UnmarshalData 解析响应数据到指定结构
func (r *ApiResult) UnmarshalData(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

// String 返回字符串表示
func (r *ApiResult) String() string {
	return r.Msg
}

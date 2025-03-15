package response

import "net/http"

type SuccessResponse struct {
	Code       string      `json:"code"`
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Count      *int        `json:"count,omitempty"`
	Page       *int        `json:"page,omitempty"`
	PageSize   *int        `json:"page_size,omitempty"`
	TotalPage  *int        `json:"total_page,omitempty"`
}

func Success(code string) *SuccessResponse {
	msg := Code(code)

	res := &SuccessResponse{
		Code:       code,
		StatusCode: http.StatusOK,
		Message:    msg.Name(),
	}

	return res
}

func (s *SuccessResponse) WithData(data interface{}) *SuccessResponse {
	s.Data = data
	return s
}

func (s *SuccessResponse) WithCount(count int) *SuccessResponse {
	s.Count = &count
	return s
}

func (s *SuccessResponse) WithPage(page int) *SuccessResponse {
	s.Page = &page
	return s
}

func (s *SuccessResponse) WithPageSize(pageSize int) *SuccessResponse {
	s.PageSize = &pageSize
	return s
}

func (s *SuccessResponse) WithTotalPage(totalPage int) *SuccessResponse {
	s.TotalPage = &totalPage
	return s
}

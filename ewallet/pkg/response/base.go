package response

var codeNames = map[Code]string{
	// Error Code 22101 - 22150
	"22101": "internal server error",
	"22102": "field %s is %s",
	"22103": "username is already exist",
	"22104": "your account exceeds the price limit",
	"22149": "not found",
	"22150": "unauthorized",
	// Success Code 22151 - 22200
	"22151": "data has been created",
	"22152": "data has been loaded",
	"22153": "data has been deleted",
	"22154": "data has been updated",
	"22155": "customer registered",
	"22156": "login success",
	"22157": "successful transfer",
	"22200": "",
}

type BaseResponse struct {
	StatusCode     int             `json:"status_code"`
	Message        string          `json:"message"`
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty"`
}

type AdditionalInfo struct {
	Usecase string `json:"usecase"`
	Info    string `json:"info"`
}

type Message string

type Code string

func (c *Code) Name() string {
	return codeNames[*c]
}

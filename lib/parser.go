package lib

type Pas struct {
	Result
	Request
}

type Result struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    string `json:"data"`
	Requuid string `json:"uuid"`
}

type Request struct {
	Action     string `json:"action"`
	Rspuuid    string `json:"uuid"`
	Desc       string `json:"desc"`
	CameraCode string `json:"cameraCode"`
	ImosCode   string `json:"imosCode"`
	MediaCode  string `json:"mediaCode"`
	RtspUrl    string `json:"rtspUrl"`
	Resolution string `json:"resolution"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	FlowName   string `json:"flowName"`
	SeekTime   string `json:"seekTime"`
	SendTime   int64  `json:"sendTime"`
}

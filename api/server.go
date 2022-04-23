package api

import (
	"encoding/json"
	"net/http"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"go.uber.org/zap"
)

type Server struct {
	listenAddress string
	stsClient     *sts.Client
	stsOption     *sts.CredentialOptions
}

func NewServer(listenAddress string, client *sts.Client, option *sts.CredentialOptions) *Server {
	return &Server{
		listenAddress: listenAddress,
		stsClient:     client,
		stsOption:     option,
	}
}

func (s Server) Run() {
	http.HandleFunc("/sts", s.getSts)

	err := http.ListenAndServe(s.listenAddress, nil)
	if err != nil {
		zap.L().Fatal("http server stopped", zap.Error(err))
	}
}

func (s Server) getSts(w http.ResponseWriter, r *http.Request) {
	zap.L().Debug("start get session token")

	result, err := s.stsClient.GetCredential(s.stsOption)
	if err != nil {
		zap.L().Error("get credential failed", zap.Error(err))
	}

	response := &StsResponse{
		TempSecretID:  result.Credentials.TmpSecretID,
		TempSecretKey: result.Credentials.TmpSecretKey,
		SessionToken:  result.Credentials.SessionToken,
		StartTime:     result.StartTime,
		ExpiredTime:   result.ExpiredTime,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		zap.L().Error("json encode failed", zap.Error(err), zap.Any("response", response))
	}

	zap.L().Info("get session token successfully", zap.Any("response", response))
}

type StsResponse struct {
	TempSecretID  string `json:"temp_secret_id"`
	TempSecretKey string `json:"temp_secret_key"`
	SessionToken  string `json:"session_token"`
	StartTime     int    `json:"start_time"`
	ExpiredTime   int    `json:"expired_time"`
}

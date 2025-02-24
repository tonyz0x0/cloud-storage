package config

const (
	// UploadServiceHost : 上传服务监听的地址
	UploadServiceHost = "0.0.0.0:8080"
	// UploadLBHost: 上传服务LB地址
	//UploadLBHost = "http://upload.fileserver.com"
	UploadLBHost = "http://localhost:28080"
	// DownloadLBHost: 下载服务LB地址
	//DownloadLBHost = "http://download.fileserver.com"
	DownloadLBHost = "http://localhost:38080"
	// TracerAgentHost: tracing agent地址
	TracerAgentHost = "127.0.0.1:6831"
)

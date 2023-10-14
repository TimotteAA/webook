package oauth2

import "webook/internal/service"

func InitWeChatService() service.WeChatService {
	return service.NewWeChatService("", "")
}

package cache

import "fmt"

const (
	UserSessionsKey = "user_sessions:%s"
	SessionUserKey  = "session_user:%s"
	UserInfoKey     = "user_info:%d"
)

// GetUserSessionKey 获取用户会话key, sessionId:session映射[hash]
func GetUserSessionKey(sessionId string) string {
	return fmt.Sprintf(UserSessionsKey, sessionId)
}

// GetSessionUserKey 获取会话用户key, sessionId:userId映射[string]
func GetSessionUserKey(sessionId string) string {
	return fmt.Sprintf(SessionUserKey, sessionId)
}

// GetUserInfoKey 获取用户信息key
func GetUserInfoKey(userId int) string {
	return fmt.Sprintf(UserInfoKey, userId)
}

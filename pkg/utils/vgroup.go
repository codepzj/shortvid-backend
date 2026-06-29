package utils

import "strings"

// 根据vgroup获取pool路径
func GetVidPoolPathFromVgroup(vgroup string) string {
	split := strings.Split(vgroup, "_")

	md5 := strings.ToLower(split[0])
	size := split[1]
	cnt := len(md5)

	return "pool/pub/" + md5[cnt-6:cnt-4] + "/" + md5[cnt-4:cnt-2] + "/" + md5[cnt-2:] + "/" + md5 + "/" + size
}

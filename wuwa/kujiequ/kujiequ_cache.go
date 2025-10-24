package kujiequ

import (
	"iter"
	"maps"
)

// 库街区全局缓存

var (

	// 缓存map中的最大数量
	MaxCache = 10

	// 库街区角色信息缓存
	globalRoles = map[string][]RoleInfo{}
)

func AddGlobalRoles(token string, roles []RoleInfo) {
	if len(globalRoles) >= MaxCache {
		deleteKey := ""
		maps.Keys(globalRoles)(func(v string) bool {
			deleteKey = v
			return false
		})
		delete(globalRoles, deleteKey)
	}
	globalRoles[token] = roles
}

func GetGlobalRoles(token string) (roles []RoleInfo, ok bool) {
	roles, ok = globalRoles[token]
	if len(roles) <= 0 {
		ok = false
	}
	return
}
func KeysGlobalRoles() iter.Seq[string] {
	return maps.Keys(globalRoles)
}

func DeleteGlobalRoles(token string) {
	delete(globalRoles, token)
}

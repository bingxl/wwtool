package kujiequ_test

import (
	"log/slog"
	"os"
	"strconv"
	"testing"
	"wuwa/kujiequ"
)

// 需要设置环境变量 k_token k_devcode, k_roleID, k_userID 分别是库街区的 token 与devcode 角色的roleID与userID
// powershell 中 $env:k_token="xxxx";$env:k_devcode="xxxx"; ....
func getKujiequ(t *testing.T) (k kujiequ.KujieQu) {
	envToken := os.Getenv("k_token")
	envDevcode := os.Getenv("k_devcode")
	if envToken == "" || envDevcode == "" {
		t.Error("请设置环境变量 k_token 与 k_devcode")
		return
	}

	token := kujiequ.Token{
		Token:   envToken,
		Devcode: envDevcode,
	}
	headers := map[string]string{
		"ip":          "192.168.100.133",
		"version":     "2.2.5",
		"versioncode": "2250",
		"distinct_id": "da0b5c51-4627-4ea6-8c54-7ddce5cb6c31",
		"model":       "Redmi K30",
		"user-agent":  "okhttp/3.11.0",
	}
	k = kujiequ.NewKujieQu(token, headers)
	return k
}
func getRole() kujiequ.RoleInfo {
	gameIdEnv := os.Getenv("k_gameID")
	id, err := strconv.Atoi(gameIdEnv)
	if err != nil {
		slog.Error("error to get gameID from env " + err.Error())
		os.Exit(1)
	}

	return kujiequ.RoleInfo{
		GameId: id, RoleId: os.Getenv("K_roleID"),
		RoleName: "bing", UserId: os.Getenv("k_userID"),
		ServerId:   os.Getenv("k_serverID"),
		ServerName: "鸣潮",
	}
}

// 显示log信息 go test -v  -run ^TestKujieQu_GetWidget$ wuwa/kujiequ
func TestKujieQu_GetWidget(t *testing.T) {
	// k := getKujiequ()
	// roles := k.FilterRoles(3)
	// if len(roles) <= 0 {
	// 	t.Error("kujiequ FindAllRoles error")
	// 	return
	// }
	role := getRole()
	t.Log("getRole return ", role)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		role    kujiequ.RoleInfo
		want    kujiequ.WidgetResponseData
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: role.RoleName, role: role},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			k := getKujiequ(t)
			got, gotErr := k.GetWidget(tt.role)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetWuwaCard() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetWuwaCard() succeeded unexpectedly")
			}
			if got.UserId == 0 {
				t.Error("got userId is 0")
			}
			slog.Info("GetWuwaCard Return value", "value", got)
			// TODO: update the condition below to compare got with tt.want.
			// if true {
			// 	t.Errorf("GetWuwaCard() = %v, want %v", got, tt.want)
			// }
		})
	}
}

// 显示log信息 go test -v  -run ^TestKujieQu_GetAllWidgets$ wuwa/kujiequ
func TestKujieQu_GetAllWidgets(t *testing.T) {
	k := getKujiequ(t)

	tests := []struct {
		name    string // description of this test case
		want    []kujiequ.WidgetResponseData
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "获取所有角色小组件测试", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := k.GetAllWidgets()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetAllWidgets() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetAllWidgets() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			t.Log("get widgets", got)
			// if true {
			// 	t.Errorf("GetAllWidgets() = %v, want %v", got, tt.want)
			// }
		})
	}
}

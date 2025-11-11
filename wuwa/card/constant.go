package card

var (
	// 国服抽卡请求链接
	RECORD_REQUEST_URL_CN = "https://gmserver-api.aki-game2.com/gacha/record/query"
	// 国际服抽卡请求链接
	RECORD_REQUEST_URL_OVERSEA = "https://gmserver-api.aki-game2.net/gacha/record/"

	// 抽卡类型 map[key:请求参数值]value:卡池类型名称
	CARD_TYPE = map[int]string{
		1: T("角色活动唤取"), // key 为传递到请求的参数, value	为卡池类型放置到存储文件中
		2: T("武器活动唤取"),
		3: T("角色常驻唤取"),
		4: T("武器常驻唤取"),
		5: T("新手唤取"),
		6: T("新手自选唤取"),
		7: T("新手自选唤取（感恩定向唤取）"),
		8: T("角色新旅唤取"),
		9: T("武器新旅唤取"),
	}
)

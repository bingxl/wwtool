package card

import (
	"bufio"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// startCardRecodeRequest 开始抽卡记录请求
//
// 参数:
//
//	gamePath string 游戏本体目录, 以Wuthering Waves Game 结尾
func StartCardRecodeRequest(gamePath ...string) {
	// 鸣潮log文件， 位于游戏本体目录下
	gameFold := "C://program/games/Wuthering Waves/Wuthering Waves Game"
	wuwaLogFile := "Client/Saved/Logs/Client.log"
	if len(gamePath) != 0 {
		gameFold = gamePath[0]
	}

	gachaUrl, err := GetLinkFromLog(filepath.Join(gameFold, wuwaLogFile))
	slog.Info("抽卡链接：" + gachaUrl)

	if err == nil && gachaUrl != "" {
		// 从抽卡链接中提取参数
		u, err := url.ParseQuery(strings.Split(gachaUrl, "?")[1])
		if err != nil {
			slog.Info(err.Error())
			return
		}
		slog.Info("抽卡链接参数", "params", u)
		// u.Get("") record_id  resources_id  svr_id  player_id
		// lang  gacha_id  gacha_type  svr_area
		cardPool := NewCardPool(u)
		slog.Info(cardPool.Country_URL())
	}
}

// 判断文件是否存在
func IsExists(filenameOrDir string) bool {
	_, err := os.Stat(filenameOrDir)
	return err == nil
}

// 从log文件中获取抽卡链接
func GetLinkFromLog(file string) (string, error) {
	pattern := regexp.MustCompile(`https://.*aki/gacha/index\.html#\/record[?=&\w\-]+`)
	wuwaLogFile := "Client/Saved/Logs/Client.log"
	if !strings.Contains(file, wuwaLogFile) {
		file = filepath.Join(file, wuwaLogFile)
	}
	f, err := os.Open(file)
	if err != nil {
		slog.Info("打开文件失败", "err", err)
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var link string
	// 寻找日志文件中最新的链接
	for scanner.Scan() {
		line := scanner.Text()
		txt := pattern.FindString(line)

		if pattern.MatchString(line) {
			// slog.Info(pattern.FindString(line))
			link = txt

		}
	}

	if err := scanner.Err(); err != nil {
		slog.Info("读取文件失败：%v", "err", err)
	}
	// slog.Info("all links", "links", link)
	return link, err
}

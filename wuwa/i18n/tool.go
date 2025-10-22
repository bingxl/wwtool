package i18n

//go:generate gotext -srclang=zh-Hans -dir=locales update -out=catalog.go -lang=en,zh-Hans wuwa/card

import (
	"log/slog"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// zh-Hans default lang
var lang = language.SimplifiedChinese

// 支持的language
var support = []language.Tag{
	lang,
	language.English,
}
var match = language.NewMatcher(support)

var prints = map[string]*message.Printer{
	lang.String(): message.NewPrinter(lang),
}

// 全局默认语言
var Printer = prints[lang.String()]

// 获取指定语言的 *message.Printer
func GetPrinter(langStr string) *message.Printer {
	return getLangPrinter(langStr)
}

// 翻译字符串，参数与 fmt.Sprintf 一致，
func Sprintf(msgId string, options ...any) string {
	return Printer.Sprintf(msgId, options...)
}

// 翻译字符串，参数与 fmt.Sprintf 一致，
var T = Sprintf

func getLangPrinter(str string) *message.Printer {
	l, _ := language.MatchStrings(match, str)
	v, ok := prints[str]
	if !ok {
		v = message.NewPrinter(l)
		prints[str] = v
	}
	return v
}

// 设置全局默认语言
func SetLang(str string) {
	// 只改变printer指向的值， Printer本身不变
	*Printer = *getLangPrinter(str)
	slog.Info("change global lang to" + str)
}

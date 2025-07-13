package i18n

//go:generate gotext -srclang=zh-Hans -dir=locales update -out=catalog.go -lang=en,zh-Hans wwtool/

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// zh-Hans default lang
var lang = language.SimplifiedChinese

// 支持的language
var support = []language.Tag{
	language.SimplifiedChinese,
	language.English,
}

// 全局默认语言
var printer = message.NewPrinter(lang)

// 翻译字符串，参数与 fmt.Sprintf 一致，
func T(msgId string, options ...any) string {
	return printer.Sprintf(msgId, options...)
}

// 设置全局默认语言
func SetLang(tag language.Tag) {
	if tag == lang {
		return
	} else {
		lang = tag
		printer = message.NewPrinter(tag)
	}

}

// 获取支持的翻译名
func GetSupportLang() []language.Tag {
	return support
}

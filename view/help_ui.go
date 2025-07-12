package view

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type helpDocT struct {
	title   string
	content []string
}

func helpUI() fyne.CanvasObject {

	docs := getHelpDocs()
	scrollsUI := make([]*widget.AccordionItem, len(docs))

	for i := 0; i < len(docs); i++ {
		doc := docs[i]
		content := widget.NewLabel(strings.Join(doc.content, "\n"))
		content.Wrapping = fyne.TextWrapWord
		accItem := widget.NewAccordionItem(doc.title, content)
		scrollsUI[i] = accItem
	}
	render := widget.NewAccordion(scrollsUI...)

	return render
}

func getHelpDocs() []helpDocT {
	docs := []helpDocT{
		{title: "启动游戏",
			content: []string{
				"如果游戏不是安装在默认路径上`C:/Program Files/Wuthering Waves/Wuthering Waves Game`",
				"则点击`添加游戏路径`后选择游戏路径添加（注: 游戏本体的路径，不是启动器的路径）",
				"下拉框中选择正确的游戏路径后点击启动",
			},
		},
		{title: "B服官服切换",
			content: []string{
				"首次运行时点击切换到官服，运行官方启动器，修复游戏；修复完成后点击切换到bilibili服，运行b服启动器，修复游戏",
				"之后的切换只需要点击对应的服就可以",
			},
		},
		{title: "抽卡链接获取",
			content: []string{
				"启动游戏，游戏里打开抽卡记录后点击`获取抽卡链接`， 获取成功后会自动复制到剪切板里",
			},
		},
	}

	return docs
}

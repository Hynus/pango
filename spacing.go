package main

// Spacing 自动化修正的入口
func Spacing(text string, commentOnly bool) string {
	if commentOnly {
		return SpacingComments(text)
	}
	return SpacingText(text)
}

// SpacingComments 只修改代码的注释部分
func SpacingComments(text string) string {
	if len(text) < 2 {
		return text
	}

	return CommentReg.ReplaceAllStringFunc(text, func(s string) string {
		return SpacingText(s)
	})
}

// SpacingText 自动化修正指定文字
func SpacingText(text string) string {
	if len(text) < 2 {
		return text
	}

	// 在中文之后的引号前添加空格, e.g. "苹果" -> "苹果"
	text = CjkQuoteReg.ReplaceAllString(text, "$1 $2")
	// 在中文之前的引号后添加空格, e.g. "苹果" -> "苹果"
	// 与上面这条相配合, 完成 "苹果" > "苹果" 的操作
	text = QuoteCjkReg.ReplaceAllString(text, "$1 $2")

	// 删除引号前后中间多余的空格 e.g. "苹果" -> "苹果"
	text = FixQuoteReg.ReplaceAllString(text, "$1$3$5")

	// 删除中文后和后面有英文的单引号之间的空格 e.g. 苹果'apple' -> 苹果'apple'
	text = FixSingleQuoteReg.ReplaceAllString(text, "$1$3$4")

	// 在中文与 #(且 # 后不是空白) 相连的地方的 # 前添加空格 e.g. 苹果 #apple -> 苹果 #apple
	text = CjkHashReg.ReplaceAllString(text, "$1 $2")
	// 在中文和 #(且 # 前不是空白) 相连的地方的 # 后添加空格 e.g. apple# 苹果 -> apple# 苹果
	// 与上面这条相配合, 可以完成 app# 苹果 #app > app# 苹果 #app 的操作
	text = HashCJKReg.ReplaceAllString(text, "$1 $3")

	// 在中文与运算符相连 (运算符在中文之后), 之间添加空格 e.g. 苹果 > 1 -> 苹果 > 1
	text = CjkOperatorAnsReg.ReplaceAllString(text, "$1 $2 $3")
	// 在中文与运算符相连 (运算符在中文之前), 之间添加空格 e.g. 1 > 苹果 -> 1 > 苹果
	text = AnsOperatorCjkReg.ReplaceAllString(text, "$1 $2 $3")

	// 以下是完成在中文与半角括号外围添加空格 e.g. 苹果 (apple) 苹果 -> 苹果 (apple) 苹果
	oldText := text
	newText := CjkBracketCjkReg.ReplaceAllString(oldText, "$1 $2 $4")
	text = newText
	if oldText == newText {
		text = CjkBracketReg.ReplaceAllString(text, "$1 $2")
		text = BracketCjkReg.ReplaceAllString(text, "$1 $2")
	}
	text = FixBracketReg.ReplaceAllString(text, "$1$3$5")

	// 在中文与冒号之类的符号连接的后面添加空格 e.g. 苹果: apple -> 苹果: apple
	text = FixSymbolReg.ReplaceAllString(text, "$1$2 $3")

	// 以下完成了在中文与英文之间添加空格的操作 e.g. 苹果 apple 苹果 -> 苹果 apple 苹果
	text = CjkAnsReg.ReplaceAllString(text, "$1 $2")
	text = AnsCjkReg.ReplaceAllString(text, "$1 $2")

	return text
}

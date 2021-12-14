package main

import (
	"bytes"
	"regexp"
	"text/template"
)

// 在当前的语义中，中文泛指中日韩文

// CJK is short for Chinese, Japanese and Korean.
//
// The constant cjk contains following Unicode blocks:
// 	\u2e80-\u2eff CJK Radicals Supplement
// 	\u2f00-\u2fdf Kangxi Radicals
// 	\u3040-\u309f Hiragana
// 	\u30a0-\u30ff Katakana
// 	\u3100-\u312f Bopomofo
// 	\u3200-\u32ff Enclosed CJK Letters and Months
// 	\u3400-\u4dbf CJK Unified Ideographs Extension A
// 	\u4e00-\u9fff CJK Unified Ideographs
// 	\uf900-\ufaff CJK Compatibility Ideographs
//
// For more information about Unicode blocks, see
// 	http://unicode-table.com/en/
const cjk = "" +
	"\u2e80-\u2eff" +
	"\u2f00-\u2fdf" +
	"\u3040-\u309f" +
	"\u30a0-\u30ff" +
	"\u3100-\u312f" +
	"\u3200-\u32ff" +
	"\u3400-\u4dbf" +
	"\u4e00-\u9fff" +
	"\uf900-\ufaff"

// ANS is short for Alphabets, Numbers
// and Symbols (`~!@#$%^&*()-_=+[]{}\|;:'",<.>/?).
//
// The constant ans doesn't contain all symbols above.
const ans = "A-Za-z0-9`\\$%\\^&\\*\\-=\\+\\\\|/\u00a1-\u00ff\u2022\u2027\u2150-\u218f"

// 会用到的正则表达式, 提前加载在内存中, 具体使用要结合操作查看, 在spacing.go中调用的地方有注释
var (
	CommentReg = regexp.MustCompile(`(/{2}).*`)

	CjkQuoteReg = regexp.MustCompile(re("([{{ .CJK}}])" + "([\"'])"))
	QuoteCjkReg = regexp.MustCompile(re("([\"'])" + "([{{.CJK}}])"))

	FixQuoteReg       = regexp.MustCompile(re("([\"'\\(\\[\\{<\u201c])" + "(\\s*)" + "(.+?)" + "(\\s*)" + "([\"'\\)\\]\\}>\u201d])"))
	FixSingleQuoteReg = regexp.MustCompile(re("([{{ .CJK}}])" + "()" + "(')" + "([A-Za-z])"))

	CjkHashReg = regexp.MustCompile(re("([{{ .CJK}}])" + "(#(\\S+))"))
	HashCJKReg = regexp.MustCompile(re("((\\S+)#)" + "([{{.CJK}}])"))

	CjkOperatorAnsReg = regexp.MustCompile(re("([{{ .CJK}}])" + "([\\+\\-\\*/=&\\|<>])" + "([A-Za-z0-9])"))
	AnsOperatorCjkReg = regexp.MustCompile(re("([A-Za-z0-9])" + "([\\+\\-\\*/=&\\|<>])" + "([{{.CJK}}])"))

	CjkBracketCjkReg = regexp.MustCompile(re("([{{ .CJK}}])" + "([\\(\\[\\{<\u201c]+(.*?)[\\)\\]\\}>\u201d]+)" + "([{{.CJK}}])"))
	CjkBracketReg    = regexp.MustCompile(re("([{{ .CJK}}])" + "([\\(\\[\\{<\u201c>])"))
	BracketCjkReg    = regexp.MustCompile(re("([\\)\\]\\}>\u201d<])" + "([{{.CJK}}])"))

	FixBracketReg = regexp.MustCompile(re("([\\(\\[\\{<\u201c]+)" + "(\\s*)" + "(.+?)" + "(\\s*)" + "([\\)\\]\\}>\u201d]+)"))

	FixSymbolReg = regexp.MustCompile(re("([{{ .CJK}}])" + "([~!;:,\\.\\?\u2026])" + "([A-Za-z0-9])"))

	CjkAnsReg = regexp.MustCompile(re("([{{ .CJK}}])([{{.ANS}}@])"))
	AnsCjkReg = regexp.MustCompile(re("([{{ .ANS}}~!;:,\\.\\?\u2026])([{{.CJK}}])"))

	context = map[string]string{
		"CJK": cjk,
		"ANS": ans,
	}
)

// 将预定义的内容嵌入模板, 形成正则表达式
func re(exp string) string {
	var (
		buf  bytes.Buffer
		tmpl = template.New("pango")
	)

	tmpl, _ = tmpl.Parse(exp)
	_ = tmpl.Execute(&buf, context)

	return buf.String()
}

# pango
####The project is modified from "https://github.com/vinta/pangu", so it's called pango
主要用于自动化给代码中中英文之间添加空格

目录下的二进制文件就是编译得到的,可以直接运行

修正终端输入的字符串<br>
`./pango t [input your test]`<br>
or<br>
`./pango text [input your test]`


修正文件或目录内所有文件 (注: 只会处理 .go 文件) 里的所有内容<br>
-w 代表覆盖原来的文件；-c 代表只修正文件代码中单行注释之后的内容<br>
`./pango f -w -c [filename|dir]`<br>
or <br>
`./pango file -w -c [filename|dir]`


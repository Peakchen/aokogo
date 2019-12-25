### 修改exe 详细信息步骤：
0. 链接:https://xuri.me/2016/06/15/embedded-icon-in-go-windows-application.html

1. 所需文件 xxx.manifest, xxx.rc, xxx.ico;
2. win 下 cmd 执行 windres -o xxx-res.syso xxx.rc  生成xxx-res.syso;
3. 编辑xxx.rc 文件时注意将文件编码方式改为ANSI 编码，这样就不会出现中文乱码;
4. 将xxx-res.syso  放置在项目下（一般在有main.go 文件同级目录），编译所得exe就能显示出详细信息.

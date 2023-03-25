# go-projects
Some projects developed with Go

### 项目1：parsemd
这个项目是为了把我之前用Markdown格式写的博客全都转成Json格式，并且生成博客条目内容（main.json文件中的内容）。

**转换成Json**的目的是因为我开发的另一个基于**Vue3**的博客系统读取json格式的文件更加方便。

具体转换过程，是把post目录里面的所有.md格式的文件全都转成.json格式文件，并放在json目录中。同时生成一个main.json文件，记录博客的条目信息。两个目录中都了有示例文件。

### 项目2：fileserv
这个项目只是简单的提供一个静态文件服务，其目的是把**parsemd**项目转成的json格式的博客文件，以文件服务的方式提供给我的基于**Vue3**的博客系统访问。
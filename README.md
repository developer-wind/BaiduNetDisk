# BaiduNetDisk

百度网盘API

# 操作说明
- 1.登录https://pan.baidu.com
- 2.抓取request头部的Cookie存入本地文件
- 3.调用BaiduNetDisk.ImportCookie("/tmp/baiduNetDisk/cookie")装载cookie
  - cookie导入目录可随意指定，请确保有可读权限
- 4.后续就可以正常的增删改查了

# 第一版实现接口如下：
- 1.获取目录
- 2.转存（支持提取码）
- 3.删除
- 4.创建目录

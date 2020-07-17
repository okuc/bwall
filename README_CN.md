# 壁纸自动设置器

# 软件说明

每日自动下载bing当日高清壁纸，设置成个人的桌面，然后从文本中读取个人设置好的座右铭、名人名言、精彩词句等，根据设置的时间轮流添加到壁纸上。

软件英文说明：https://github.com/okuc/bwall/blob/master/README.md

软件下载地址：https://gitee.com/okuc/bwall/raw/master/release/bwall.7z

# 配置介绍

```
{
  "CurrentText": "宠辱莫惊　闲看庭前花开花落$去留无意　漫随天外云卷云舒",
  "MottoFileName": "mymotto.txt",
  "Interval": 1
}
```

- `CurrentText：当前座右铭。忽略即可，系统会自动配置。`
- `MottoFileName`：座右铭文本文件。格式为一行或多行。不同的座右铭，用空行隔开即可。
- `Interval`：座右铭切换时间，单位为分钟，默认为1分钟。

# 感谢
本项目基于以下项目演进而来，特别表示感谢。
- https://github.com/zhimiaoli/bwall.git
- https://github.com/chkhetiani/bwall.git
---
name: bilibili-transcript
description: 给定一个 B站视频链接，通过 WebFetch 抓取字幕并整理成连贯的讲稿文本（.txt）。如果字幕是非中文语言，同时输出原文版和中文翻译版。当用户提供 B站链接并要求获取讲稿、字幕、文字稿时触发。
---

# Bilibili 视频讲稿提取

用户提供了一个 B站视频链接，按以下步骤提取字幕并整理成讲稿。

## 执行步骤

### Step 1：从 URL 提取 BV 号

从用户提供的链接中解析出 `bvid`，格式为 `BV` 开头的字符串，例如 `BV1xxxxxxx`。

### Step 2：获取视频基本信息（CID）

使用 WebFetch 请求：
```
https://api.bilibili.com/x/web-interface/view?bvid=<bvid>
```

从响应中提取：
- `data.cid` — 第一分P的 Content ID
- `data.title` — 视频标题
- `data.desc` — 视频简介
- `data.pages` — 如果是多P视频，列出所有分P

### Step 3：获取字幕列表

使用 WebFetch 请求（带 Referer header）：
```
https://api.bilibili.com/x/player/v2?bvid=<bvid>&cid=<cid>
```

从响应中提取 `data.subtitle.subtitles` 数组，获取每条字幕的：
- `lan` — 语言代码（如 `zh-CN`、`en`、`ai-zh`）
- `lan_doc` — 语言名称
- `subtitle_url` — 字幕文件路径（补全为 `https:` 开头）

如果 `subtitles` 为空数组，说明该视频没有字幕，告知用户并停止。

### Step 4：下载字幕 JSON

对每个字幕文件，使用 WebFetch 请求完整 URL：
```
https:{subtitle_url}
```

提取 `body` 数组中每条的 `content` 字段，按顺序拼接成原始字幕文本。

### Step 5：整理成连贯讲稿

将原始字幕文本整理为连贯的文章形式：
- 去除重复片段
- 合并断句��补全标点
- 按语义分段，每段之间空一行
- 不保留时间戳

### Step 6：处理多语言

- 如果字幕语言是 `zh-CN` 或 `ai-zh`（AI自动中文）：只输出一份中文讲稿
- 如果字幕语言是其他语言（英文、日文等）：
  1. 输出原文讲稿
  2. 将原文翻译成中文，输出中文讲稿
- 优先使用人工字幕（非 `ai-` 前缀），其次使用 AI 字幕

### Step 7：输出文件

将整理好的讲稿以代码块形式展示给用户，并说明文件名建议：
- 单语言：`<视频标题>.txt`
- 双语言：`<视频标题>_原文.txt` 和 `<视频标题>_中文.txt`

询问用户是否需要保存到本地文件（如需要，使用 Write 工具写入当前工作目录）。

## 错误处理

- **API 返回 -352 或 412**：B站反爬触发，告知用户需要登录 Cookie，请求用户提供 `SESSDATA` 值
- **subtitles 为空**：该视频无字幕，告知用户
- **subtitle_url 为空字符串**：字幕文件暂不可用，告知用户
- **多P视频**：默认处理第一P，询问用户是否需要处理其他分P

## 注意事项

- WebFetch 请求 B站 API 时，提取 JSON 中的关键字段即可，不需要解析整个响应
- 字幕整理时保持原意，不要过度改写
- 翻译时保持专业术语准确

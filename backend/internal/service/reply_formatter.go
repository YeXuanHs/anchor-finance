package service

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// 允许的 HTML 标签
var allowedTags = map[string]bool{
	"p": true, "br": true, "strong": true, "b": true, "em": true, "i": true, "u": true,
	"ul": true, "ol": true, "li": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"blockquote": true, "pre": true, "code": true,
	"a": true, "span": true, "div": true, "hr": true,
	"table": true, "thead": true, "tbody": true, "tr": true, "th": true, "td": true,
}

// MarkdownToHTML 将 AI 的 Markdown 回复转为工单可渲染的 HTML
func MarkdownToHTML(markdown string) string {
	text := strings.TrimSpace(markdown)
	if text == "" {
		return ""
	}
	if looksLikeHTML(text) {
		return SanitizeHTML(text)
	}
	return SanitizeHTML(markdownToHTML(text))
}

// SanitizeHTML 清理 HTML，只保留安全标签
func SanitizeHTML(htmlStr string) string {
	// 解码 HTML 实体
	htmlStr = html.UnescapeString(htmlStr)

	// 移除不在白名单中的标签
	tagPattern := regexp.MustCompile(`</?(\w+)(?:\s[^>]*)?>`)
	htmlStr = tagPattern.ReplaceAllStringFunc(htmlStr, func(match string) string {
		// 提取标签名
		parts := tagPattern.FindStringSubmatch(match)
		if len(parts) < 2 {
			return ""
		}
		tagName := strings.ToLower(parts[1])
		if !allowedTags[tagName] {
			return ""
		}
		return match
	})

	// 移除 on* 事件属性
	onAttrPattern := regexp.MustCompile(`\s+on\w+\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]*)`)
	htmlStr = onAttrPattern.ReplaceAllString(htmlStr, "")

	// 移除 style 和 class 属性
	stylePattern := regexp.MustCompile(`\s+(?:style|class)\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]*)`)
	htmlStr = stylePattern.ReplaceAllString(htmlStr, "")

	// 清理 a 标签的 href
	linkPattern := regexp.MustCompile(`<a\s+([^>]*?)href\s*=\s*(?:"([^"]*?)"|'([^']*?)')([^>]*)>`)
	htmlStr = linkPattern.ReplaceAllStringFunc(htmlStr, func(match string) string {
		parts := linkPattern.FindStringSubmatch(match)
		if len(parts) < 4 {
			return "<span>"
		}
		href := parts[2]
		if href == "" {
			href = parts[3]
		}
		// 只允许 http/https/mailto/tel/相对路径
		if !isAllowedHref(href) {
			return "<span>"
		}
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">`, html.EscapeString(href))
	})

	return strings.TrimSpace(htmlStr)
}

// AppendAIDisclaimer 追加 AI 免责声明
func AppendAIDisclaimer(htmlStr string) string {
	htmlStr = strings.TrimSpace(htmlStr)
	if htmlStr == "" {
		return htmlStr
	}
	return htmlStr + `<hr /><p><em>—— 本条回复由 AI 自动生成</em></p>`
}

// ─── 内部实现 ───

var htmlTagPattern = regexp.MustCompile(`<(?:p|br|ul|ol|li|strong|em|h[1-6]|div|table|blockquote|pre)\b`)

func looksLikeHTML(text string) bool {
	return htmlTagPattern.MatchString(text)
}

func isAllowedHref(href string) bool {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return true
	}
	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") {
		return true
	}
	if strings.HasPrefix(href, "/") {
		return true
	}
	return false
}

// markdownToHTML 将 Markdown 转为 HTML（不处理嵌套，简单实现）
func markdownToHTML(text string) string {
	lines := strings.Split(text, "\n")
	var out []string
	i := 0
	count := len(lines)

	for i < count {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// 空行跳过
		if trimmed == "" {
			i++
			continue
		}

		// 代码块 ```lang ... ```
		if strings.HasPrefix(trimmed, "```") {
			i++
			var codeLines []string
			for i < count && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
				codeLines = append(codeLines, lines[i])
				i++
			}
			if i < count {
				i++
			}
			out = append(out, "<pre><code>"+escapeHTML(strings.Join(codeLines, "\n"))+"</code></pre>")
			continue
		}

		// 分隔线 --- or *** or ___
		if isHorizontalRule(trimmed) {
			out = append(out, "<hr />")
			i++
			continue
		}

		// 标题 # ~ ######
		if headingMatch := matchHeading(trimmed); headingMatch != nil {
			level := headingMatch[0]
			content := headingMatch[1]
			out = append(out, fmt.Sprintf("<h%d>%s</h%d>", level, processInline(content), level))
			i++
			continue
		}

		// 表格 |...|...|
		if isTableRow(trimmed) {
			var tableLines []string
			for i < count && isTableRow(strings.TrimSpace(lines[i])) {
				tableLines = append(tableLines, lines[i])
				i++
			}
			out = append(out, renderTable(tableLines))
			continue
		}

		// 引用块 >
		if strings.HasPrefix(trimmed, ">") {
			var quoteLines []string
			for i < count && strings.HasPrefix(strings.TrimSpace(lines[i]), ">") {
				qLine := strings.TrimSpace(lines[i])
				qLine = strings.TrimPrefix(qLine, ">")
				qLine = strings.TrimSpace(qLine)
				quoteLines = append(quoteLines, qLine)
				i++
			}
			out = append(out, "<blockquote><p>"+processInline(strings.Join(quoteLines, "\n"))+"</p></blockquote>")
			continue
		}

		// 无序列表 - / * / +
		if matchUnorderedList(trimmed) {
			var items []string
			for i < count && matchUnorderedList(strings.TrimSpace(lines[i])) {
				item := regexp.MustCompile(`^\s*[-*+]\s+`).ReplaceAllString(lines[i], "")
				items = append(items, item)
				i++
			}
			out = append(out, renderList("ul", items))
			continue
		}

		// 有序列表 1. 2. ...
		if matchOrderedList(trimmed) {
			var items []string
			for i < count && matchOrderedList(strings.TrimSpace(lines[i])) {
				item := regexp.MustCompile(`^\s*\d+\.\s+`).ReplaceAllString(lines[i], "")
				items = append(items, item)
				i++
			}
			out = append(out, renderList("ol", items))
			continue
		}

		// 段落（收集连续非空行直到遇到特殊语法）
		var paraLines []string
		for i < count {
			cur := strings.TrimSpace(lines[i])
			if cur == "" ||
				strings.HasPrefix(cur, "```") ||
				matchHeading(cur) != nil ||
				isTableRow(cur) ||
				strings.HasPrefix(cur, ">") ||
				matchUnorderedList(cur) ||
				matchOrderedList(cur) ||
				isHorizontalRule(cur) {
				break
			}
			paraLines = append(paraLines, lines[i])
			i++
		}
		if len(paraLines) > 0 {
			inner := processInline(strings.Join(paraLines, "\n"))
			inner = strings.ReplaceAll(inner, "\n", "<br />")
			out = append(out, "<p>"+inner+"</p>")
		}
	}

	return strings.Join(out, "")
}

// ─── 内联格式处理 ───

func processInline(text string) string {
	s := escapeHTML(text)

	// 加粗 ***text*** 或 **text**
	s = regexp.MustCompile(`\*\*\*([^*\n]+?)\*\*\*`).ReplaceAllString(s, "<strong>$1</strong>")
	s = regexp.MustCompile(`\*\*([^*\n]+?)\*\*`).ReplaceAllString(s, "<strong>$1</strong>")

	// 斜体 *text*
	s = regexp.MustCompile(`(?<!\*)\*([^*\n]+?)\*(?!\*)`).ReplaceAllString(s, "<em>$1</em>")

	// 加粗 __text__
	s = regexp.MustCompile(`__([^_\n]+?)__`).ReplaceAllString(s, "<strong>$1</strong>")

	// 斜体 _text_
	s = regexp.MustCompile(`(?<!\w)_([^_\n]+?)_(?!\w)`).ReplaceAllString(s, "<em>$1</em>")

	// 删除线 ~~text~~
	s = regexp.MustCompile(`~~([^~\n]+?)~~`).ReplaceAllString(s, "<del>$1</del>")

	// 行内代码 `text`
	s = regexp.MustCompile("`([^`\\n]+?)`").ReplaceAllString(s, "<code>$1</code>")

	// 链接 [text](url)
	linkRe := regexp.MustCompile(`\[([^\]\n]+?)\]\(([^)\s]+)\)`)
	s = linkRe.ReplaceAllStringFunc(s, func(match string) string {
		parts := linkRe.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		text := parts[1]
		url := parts[2]
		if !isAllowedHref(url) {
			return match
		}
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, html.EscapeString(url), text)
	})

	// 自动链接 https://...
	autoLinkRe := regexp.MustCompile(`https?://[^\s<>&]+`)
	s = autoLinkRe.ReplaceAllStringFunc(s, func(url string) string {
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, html.EscapeString(url), url)
	})

	return s
}

// ─── 辅助函数 ───

func escapeHTML(text string) string {
	return html.EscapeString(text)
}

func isHorizontalRule(line string) bool {
	re := regexp.MustCompile(`^\s*[-*_]{3,}\s*$`)
	return re.MatchString(line)
}

func matchHeading(line string) []string {
	re := regexp.MustCompile(`^\s*(#{1,6})\s+(.+?)\s*$`)
	parts := re.FindStringSubmatch(line)
	if len(parts) < 3 {
		return nil
	}
	level := len(parts[1])
	if level > 6 {
		level = 6
	}
	return []string{fmt.Sprintf("%d", level), parts[2]}
}

func isTableRow(line string) bool {
	return strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|")
}

func matchUnorderedList(line string) bool {
	re := regexp.MustCompile(`^\s*[-*+]\s+`)
	return re.MatchString(line)
}

func matchOrderedList(line string) bool {
	re := regexp.MustCompile(`^\s*\d+\.\s+`)
	return re.MatchString(line)
}

func renderList(tag string, items []string) string {
	var lis []string
	for _, item := range items {
		lis = append(lis, "<li>"+processInline(item)+"</li>")
	}
	return fmt.Sprintf("<%s>%s</%s>", tag, strings.Join(lis, ""), tag)
}

func renderTable(lines []string) string {
	var rows [][]string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 跳过分隔行 |---|---|
		if isSeparatorRow(trimmed) {
			continue
		}
		// 移除首尾 |
		trimmed = strings.TrimPrefix(trimmed, "|")
		trimmed = strings.TrimSuffix(trimmed, "|")
		cells := strings.Split(trimmed, "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		rows = append(rows, cells)
	}

	if len(rows) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("<table><thead><tr>")
	for _, cell := range rows[0] {
		b.WriteString("<th>" + processInline(cell) + "</th>")
	}
	b.WriteString("</tr></thead><tbody>")
	for _, row := range rows[1:] {
		b.WriteString("<tr>")
		for _, cell := range row {
			b.WriteString("<td>" + processInline(cell) + "</td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table>")
	return b.String()
}

func isSeparatorRow(line string) bool {
	re := regexp.MustCompile(`^\|[\s\-:|]+\|$`)
	return re.MatchString(line)
}

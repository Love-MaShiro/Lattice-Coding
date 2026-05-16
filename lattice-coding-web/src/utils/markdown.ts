const blockTokenPrefix = '\u0000MD_BLOCK_'

function escapeHtml(value: string) {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function escapeAttribute(value: string) {
  return escapeHtml(value).replace(/`/g, '&#96;')
}

function isSafeUrl(value: string) {
  return /^(https?:\/\/|mailto:)/i.test(value)
}

function renderInline(value: string) {
  return value
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/__([^_]+)__/g, '<strong>$1</strong>')
    .replace(/\*([^*]+)\*/g, '<em>$1</em>')
    .replace(/_([^_]+)_/g, '<em>$1</em>')
    .replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_match, label: string, url: string) => {
      if (!isSafeUrl(url)) return label
      return `<a href="${escapeAttribute(url)}" target="_blank" rel="noopener noreferrer">${label}</a>`
    })
}

function renderList(lines: string[], ordered: boolean) {
  const tag = ordered ? 'ol' : 'ul'
  const items = lines
    .map((line) => {
      const content = ordered ? line.replace(/^\d+\.\s+/, '') : line.replace(/^[-*]\s+/, '')
      return `<li>${renderInline(content)}</li>`
    })
    .join('')

  return `<${tag}>${items}</${tag}>`
}

function renderParagraph(lines: string[]) {
  return `<p>${renderInline(lines.join('<br>'))}</p>`
}

function renderBlock(block: string) {
  const lines = block.split('\n')
  const firstLine = lines[0] || ''

  if (/^#{1,6}\s+/.test(firstLine) && lines.length === 1) {
    const level = firstLine.match(/^#{1,6}/)?.[0].length || 1
    return `<h${level}>${renderInline(firstLine.replace(/^#{1,6}\s+/, ''))}</h${level}>`
  }

  if (lines.every((line) => /^[-*]\s+/.test(line))) {
    return renderList(lines, false)
  }

  if (lines.every((line) => /^\d+\.\s+/.test(line))) {
    return renderList(lines, true)
  }

  if (lines.every((line) => /^>\s?/.test(line))) {
    const quote = lines.map((line) => line.replace(/^>\s?/, '')).join('<br>')
    return `<blockquote>${renderInline(quote)}</blockquote>`
  }

  return renderParagraph(lines)
}

export function renderMarkdown(markdown: string) {
  if (!markdown) return ''

  const codeBlocks: string[] = []
  const escaped = escapeHtml(markdown).replace(/```([a-zA-Z0-9_-]+)?\n?([\s\S]*?)```/g, (_match, language = '', code = '') => {
    const index = codeBlocks.push(
      `<pre><code${language ? ` class="language-${escapeAttribute(language)}"` : ''}>${code.trim()}</code></pre>`
    ) - 1
    return `${blockTokenPrefix}${index}\u0000`
  })

  const rendered = escaped
    .split(/\n{2,}/)
    .map((block) => {
      const trimmed = block.trim()
      if (!trimmed) return ''

      const tokenMatch = trimmed.match(new RegExp(`^${blockTokenPrefix}(\\d+)\\u0000$`))
      if (tokenMatch) return codeBlocks[Number(tokenMatch[1])] || ''

      return renderBlock(trimmed)
    })
    .filter(Boolean)
    .join('')

  return rendered.replace(new RegExp(`${blockTokenPrefix}(\\d+)\\u0000`, 'g'), (_match, index: string) => {
    return codeBlocks[Number(index)] || ''
  })
}

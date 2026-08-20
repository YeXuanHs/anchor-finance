/**
 * 导出数据为 CSV 文件
 * @param data 要导出的数据数组
 * @param columns 列定义 [{ key: '字段名', title: '列标题' }]
 * @param filename 文件名（不含扩展名）
 */
export function exportToCSV(data: any[], columns: { key: string; title: string }[], filename: string) {
  if (!data || data.length === 0) {
    return
  }
  const header = columns.map(c => `"${c.title}"`).join(',')
  const rows = data.map(row =>
    columns.map(c => {
      const val = row[c.key]
      if (val === null || val === undefined) return '""'
      return `"${String(val).replace(/"/g, '""')}"`
    }).join(',')
  )
  const csv = '\uFEFF' + [header, ...rows].join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `${filename}_${new Date().toISOString().slice(0, 10)}.csv`
  link.click()
  URL.revokeObjectURL(link.href)
}
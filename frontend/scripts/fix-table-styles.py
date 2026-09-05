import re, glob

files = glob.glob('src/features/*/components/*Table.tsx')
changed = []

# Common old patterns to replace with unified style
old_wrapper_patterns = [
    # Pattern 1: rounded-md border (old generated)
    (r'<div className="rounded-md border">\n      <table className="w-full text-sm">\n      <thead className="border-b bg-gray-50">',
     '<div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">\n      <table className="w-full text-sm text-left">\n      <thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">'),
    # Pattern 2: rounded-xl border border-black/5 (FloorTable, RoomTable style)
    (r'<div className="rounded-xl border border-black/5 overflow-hidden bg-white">\n      <table className="w-full text-sm">\n      <thead className="border-b border-black/5 bg-black/\[0\.02\]">',
     '<div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">\n      <table className="w-full text-sm text-left">\n      <thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">'),
]

# Common th class replacements
th_replacements = [
    (r'<th key=\{header\.id\} className="px-4 py-3 text-left font-medium">',
     '<th key={header.id} className="px-5 py-3.5">'),
    (r'<th key=\{header\.id\} className="px-4 py-3 text-left font-semibold text-black/60 text-xs uppercase tracking-wider">',
     '<th key={header.id} className="px-5 py-3.5">'),
]

# tbody and tr styles
tbody_tr_patterns = [
    (r'<tbody>\n          \{table\.getRowModel\(\)\.rows\.map\(\(row\) => \(\n            <tr\n              key=\{row\.id\}\n              className="border-b border-black/5 hover:bg-black/\[0\.02\] cursor-pointer transition-colors"',
     '<tbody className="divide-y divide-slate-100">\n          {table.getRowModel().rows.map((row) => (\n            <tr\n              key={row.id}\n              className="hover:bg-orange/5 cursor-pointer transition-colors"'),
    (r'<tbody>\n          \{table\.getRowModel\(\)\.rows\.map\(\(row\) => \(\n            <tr\n              key=\{row\.id\}\n              className="border-b hover:bg-gray-50 cursor-pointer"',
     '<tbody className="divide-y divide-slate-100">\n          {table.getRowModel().rows.map((row) => (\n            <tr\n              key={row.id}\n              className="hover:bg-orange/5 cursor-pointer transition-colors"'),
]

# td class
td_replacement = (
    r'<td key=\{cell\.id\} className="px-4 py-3">',
    '<td key={cell.id} className="px-5 py-4">',
)

# Empty state
empty_patterns = [
    (r'<td colSpan=\{columns\.length\} className="px-4 py-8 text-center text-gray-500">\n                No data found\.\n              </td>',
     '<td colSpan={columns.length} className="px-5 py-12 text-center text-slate-400">\n                No data found.\n              </td>'),
    (r'<td colSpan=\{columns\.length\} className="px-4 py-12 text-center text-black/30">\n                <.*?className="mx-auto mb-2 opacity-50" />\n                No \w+ found\. Create one to get started\.\n              </td>',
     '<td colSpan={columns.length} className="px-5 py-12 text-center text-slate-400">\n                No data found.\n              </td>'),
]

for f in files:
    src = open(f).read()
    orig = src

    for old, new in old_wrapper_patterns:
        src = re.sub(old, new, src)

    for old, new in th_replacements:
        src = re.sub(old, new, src)

    for old, new in tbody_tr_patterns:
        src = re.sub(old, new, src)

    src = re.sub(td_replacement[0], td_replacement[1], src)

    for old, new in empty_patterns:
        src = re.sub(old, new, src, flags=re.DOTALL)

    # Also fix data: data -> data: data || [] in useReactTable
    src = re.sub(r'const table = useReactTable\(\{\n    data,\n    columns,', 'const table = useReactTable({\n    data: data || [],\n    columns,', src)

    if src != orig:
        open(f, 'w').write(src)
        changed.append(f)

print(f"Changed {len(changed)} files")
for f in changed:
    print(f)

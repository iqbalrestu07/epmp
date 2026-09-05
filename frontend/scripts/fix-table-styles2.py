import glob, re

changed = []

table_files = glob.glob('src/features/*/components/*Table.tsx')

for f in table_files:
    src = open(f).read()
    orig = src

    # Fix table wrapper
    src = src.replace(
        '<div className="rounded-md border">',
        '<div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">'
    )
    src = src.replace(
        '<div className="rounded-xl border border-black/5 overflow-hidden bg-white">',
        '<div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">'
    )
    src = src.replace(
        '<div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">',
        '<div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">'
    )

    # Fix table element
    src = src.replace(
        '<table className="w-full text-sm">',
        '<table className="w-full text-sm text-left">'
    )

    # Fix thead
    src = src.replace(
        '<thead className="border-b bg-slate-50">',
        '<thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">'
    )
    src = src.replace(
        '<thead className="border-b border-black/5 bg-black/[0.02]">',
        '<thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">'
    )

    # Fix tbody if it doesn't have divide-y
    src = src.replace(
        '<tbody>',
        '<tbody className="divide-y divide-slate-100">'
    )
    # But don't double-add if already has divide-y
    src = src.replace(
        '<tbody className="divide-y divide-slate-100" className="divide-y divide-slate-100">',
        '<tbody className="divide-y divide-slate-100">'
    )

    if src != orig:
        open(f, 'w').write(src)
        changed.append(f)

print(f"Changed {len(changed)} files")
for f in sorted(changed):
    print(f)

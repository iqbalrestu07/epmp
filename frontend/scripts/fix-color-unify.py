import glob, re

files = glob.glob('src/**/*.tsx', recursive=True) + glob.glob('src/**/*.ts', recursive=True)
changed = []

# Color unification mappings (order matters - more specific first)
replacements = [
    # Black opacity classes -> slate equivalents
    ('bg-black/[0.02]', 'bg-slate-50/75'),
    ('bg-black/5', 'bg-slate-100'),
    ('bg-black/10', 'bg-slate-200'),
    ('border-black/5', 'border-slate-200'),
    ('border-black/10', 'border-slate-300'),
    ('divide-black/5', 'divide-slate-100'),
    ('text-black/30', 'text-slate-400'),
    ('text-black/40', 'text-slate-500'),
    ('text-black/50', 'text-slate-500'),
    ('text-black/60', 'text-slate-600'),
    ('text-black/70', 'text-slate-600'),
    ('text-black/80', 'text-slate-800'),
    ('text-black/90', 'text-slate-900'),
    ('hover:text-black', 'hover:text-slate-900'),
    ('hover:bg-black/5', 'hover:bg-slate-100'),
    ('hover:bg-black/[0.02]', 'hover:bg-slate-50'),
    ('text-black', 'text-slate-900'),
    ('bg-black', 'bg-slate-900'),

    # Gray classes -> slate equivalents
    ('bg-gray-50', 'bg-slate-50'),
    ('bg-gray-100', 'bg-slate-100'),
    ('bg-gray-200', 'bg-slate-200'),
    ('hover:bg-gray-50', 'hover:bg-slate-50'),
    ('hover:bg-gray-100', 'hover:bg-slate-100'),
    ('text-gray-400', 'text-slate-400'),
    ('text-gray-500', 'text-slate-500'),
    ('text-gray-600', 'text-slate-600'),
    ('text-gray-700', 'text-slate-700'),
    ('text-gray-800', 'text-slate-800'),
    ('text-gray-900', 'text-slate-900'),
    ('border-gray-200', 'border-slate-200'),
    ('border-gray-300', 'border-slate-300'),
    ('border-gray-400', 'border-slate-400'),
    ('divide-gray-100', 'divide-slate-100'),
    ('divide-gray-200', 'divide-slate-200'),
    ('ring-gray-300', 'ring-slate-300'),
    ('ring-gray-400', 'ring-slate-400'),

    # Green/Red/Yellow/Blue 100 -> 50 for badge backgrounds (more subtle)
    ('bg-green-100 text-green-700', 'bg-green-50 text-green-700 border border-green-200'),
    ('bg-red-100 text-red-700', 'bg-red-50 text-red-700 border border-red-200'),
    ('bg-blue-100 text-blue-700', 'bg-blue-50 text-blue-700 border border-blue-200'),
    ('bg-yellow-100 text-yellow-700', 'bg-yellow-50 text-yellow-700 border border-yellow-200'),
]

for f in files:
    src = open(f).read()
    orig = src
    for old, new in replacements:
        src = src.replace(old, new)
    if src != orig:
        open(f, 'w').write(src)
        changed.append(f)

print(f"Changed {len(changed)} files")
for f in sorted(changed):
    print(f)

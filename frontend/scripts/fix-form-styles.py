import glob, re

changed = []

# ─── 1. Unify form styles across all form components ───
form_files = glob.glob('src/features/*/components/*Form.tsx')
form_files += glob.glob('src/features/*/components/Interactive*.tsx')

# Unified select className pattern
OLD_SELECT_PATTERNS = [
    # shadcn default style
    r'className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"',
    # basic style
    r'className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm"',
    # variant without bg-white
    r'className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"',
    # variant with py-2.5
    r'className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2\.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"',
]
NEW_SELECT_CLASS = 'className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"'

# Unified error text
OLD_ERROR = 'className="text-sm text-red-500"'
NEW_ERROR = 'className="text-xs text-red-600 mt-1"'

# Unified form spacing
OLD_FORM_SPACING = 'className="space-y-4"'
NEW_FORM_SPACING = 'className="space-y-5"'

# Remove max-w-md from forms (the create page wrapper handles width)
OLD_MAXWMD = 'className="space-y-4 max-w-md"'
NEW_MAXWMD = 'className="space-y-5"'

# Unified submit button
OLD_SUBMIT_PATTERNS = [
    r'<Button type="submit" disabled=\{isSubmitting\}>\n        \{isSubmitting \? "Saving\.\.\." : "Save"\}\n      </Button>',
    r'<Button type="submit" disabled=\{isSubmitting\}>\n        \{isSubmitting \? "Saving\.\.\." : "Save Contract"\}\n      </Button>',
]
NEW_SUBMIT = '<Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">\n        {isSubmitting ? "Saving..." : "Save"}\n      </Button>'

for f in form_files:
    src = open(f).read()
    orig = src

    # Fix select classNames
    for pat in OLD_SELECT_PATTERNS:
        src = re.sub(pat, NEW_SELECT_CLASS, src)

    # Fix error text
    src = src.replace(OLD_ERROR, NEW_ERROR)

    # Fix form spacing
    src = src.replace(OLD_MAXWMD, NEW_MAXWMD)
    src = src.replace(OLD_FORM_SPACING, NEW_FORM_SPACING)

    # Fix submit buttons - match any "Save X" text
    src = re.sub(
        r'<Button type="submit" disabled=\{isSubmitting\}>\n        \{isSubmitting \? "Saving\.\.\." : "Save[^"]*"\}\n      </Button>',
        '<Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">\n        {isSubmitting ? "Saving..." : "Save"}\n      </Button>',
        src
    )

    if src != orig:
        open(f, 'w').write(src)
        changed.append(f)

# ─── 2. Fix create page layouts (remove h-full flex flex-col that causes overflow) ───
create_files = glob.glob('src/features/*/pages/*CreatePage.tsx')
create_files += glob.glob('src/features/*/pages/*EditPage.tsx')

for f in create_files:
    src = open(f).read()
    orig = src

    # Replace h-full flex flex-col with normal flow (for pages with basic form mode)
    src = src.replace('className="space-y-6 h-full flex flex-col"', 'className="space-y-6"')

    # Replace flex-1 min-h-[600px] with normal div (for basic form content)
    src = src.replace('className="flex-1 min-h-[600px]"', 'className="min-h-[400px]"')

    # For pages without interactive mode, also fix the form container
    # Add proper card wrapper styling if missing
    src = src.replace(
        'className="max-w-2xl bg-white p-6 rounded-xl border border-slate-200 shadow-sm"',
        'className="max-w-2xl bg-white p-8 rounded-xl border border-slate-200 shadow-sm"'
    )

    if src != orig:
        open(f, 'w').write(src)
        if f not in changed:
            changed.append(f)

# ─── 3. Fix detail pages - better loading and not found states ───
detail_files = glob.glob('src/features/*/pages/*DetailPage.tsx')

for f in detail_files:
    src = open(f).read()
    orig = src

    # Better loading state
    src = src.replace(
        'if (isLoading) return <div>Loading...</div>;',
        'if (isLoading) return <div className="text-center py-12 text-slate-400">Loading...</div>;'
    )
    # Better not found state
    src = src.replace(
        'if (!data) return <div>Not found</div>;',
        'if (!data) return <div className="text-center py-12 text-slate-400">Not found</div>;'
    )

    # Better detail card wrapper - wrap the dl in a card
    src = src.replace(
        '<dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">',
        '<dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 bg-white rounded-xl border border-slate-200 shadow-sm p-6">'
    )

    # Better dt styling
    src = src.replace(
        'className="text-sm font-medium text-slate-500"',
        'className="text-xs font-semibold text-slate-500 uppercase tracking-wider"'
    )

    # Better dd styling
    src = src.replace(
        '<dd className="text-sm">',
        '<dd className="text-sm text-slate-800 mt-1">'
    )

    if src != orig:
        open(f, 'w').write(src)
        if f not in changed:
            changed.append(f)

# ─── 4. Fix list page loading states ───
list_files = glob.glob('src/features/*/pages/*ListPage.tsx')

for f in list_files:
    src = open(f).read()
    orig = src

    # Better loading state
    src = src.replace(
        '<div className="text-center py-8">Loading...</div>',
        '<div className="text-center py-12 text-slate-400">Loading...</div>'
    )
    src = src.replace(
        '<div className="text-center py-8 text-slate-400">Loading...</div>',
        '<div className="text-center py-12 text-slate-400">Loading...</div>'
    )

    if src != orig:
        open(f, 'w').write(src)
        if f not in changed:
            changed.append(f)

print(f"Changed {len(changed)} files")
for f in sorted(changed):
    print(f)

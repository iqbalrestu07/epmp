import glob, re

changed = []

# Target all generated list pages (skip Property, Building, Floor which are hand-crafted)
list_files = glob.glob('src/features/*/pages/*ListPage.tsx')

# Skip the reference pages we don't want to change
SKIP = [
    'src/features/property/pages/PropertyListPage.tsx',
    'src/features/building/pages/BuildingListPage.tsx',
    'src/features/floor/pages/FloorListPage.tsx',
    'src/features/asset/pages/AssetListPage.tsx',  # has custom style already
    'src/features/organization/pages/OrganizationListPage.tsx',  # special page
]

for f in list_files:
    if f in SKIP:
        continue

    src = open(f).read()
    orig = src

    # Extract module name from file path for dynamic breadcrumb/icon
    # e.g. src/features/tenant/pages/TenantListPage.tsx -> Tenant
    parts = f.split('/')
    feature_name = parts[2]  # 'tenant', 'contract', etc.
    # Capitalize first letter
    display_name = feature_name[0].upper() + feature_name[1:]
    # For compound names like 'tenantidentity' -> 'Tenant Identity'
    display_name = re.sub(r'([a-z])([A-Z])', r'\1 \2', display_name)
    # For 'assetassignment' etc
    display_name = re.sub(r'([A-Z])', r' \1', display_name).strip()
    display_name = ' '.join(w[0].upper() + w[1:] for w in display_name.split())

    # Pluralize
    plural_name = display_name + 's'

    # Dashboard path
    dashboard_path = '/dashboard/' + feature_name.replace('_', '-')

    # ─── Fix 1: Add breadcrumb + icon to header ───
    # Replace bare h1 header with breadcrumb + icon header
    # Pattern: <div className="flex items-center justify-between">\n        <h1 className="text-2xl font-bold">X</h1>\n        <Button onClick={() => navigate("...")}>\n          New X\n        </Button>\n      </div>

    # Add Plus icon to New button
    src = re.sub(
        r'<Button onClick=\{\(\) => navigate\("(/dashboard/[^"]+/new)"\)}>\n          New (\w+)',
        r'<Button onClick={() => navigate("\1")}>\n          <Plus className="w-4 h-4 mr-1" />\n          New \2',
        src
    )

    # Replace bare header with breadcrumb + icon
    # Match: <h1 className="text-2xl font-bold">X</h1>
    # Replace with breadcrumb pattern
    old_header = f'<h1 className="text-2xl font-bold">{plural_name}</h1>'
    new_header = f'''<div>
          <div className="flex items-center gap-2 text-sm text-slate-500 mb-1">
            <Link to="/dashboard" className="hover:text-slate-900">Dashboard</Link>
            <span>/</span>
            <span className="text-slate-900 font-medium">{plural_name}</span>
          </div>
          <h1 className="text-2xl font-bold text-slate-900">{plural_name}</h1>
        </div>'''
    src = src.replace(old_header, new_header)

    # Also handle singular names that might not have been pluralized
    old_header2 = f'<h1 className="text-2xl font-bold">{display_name}</h1>'
    new_header2 = f'''<div>
          <div className="flex items-center gap-2 text-sm text-slate-500 mb-1">
            <Link to="/dashboard" className="hover:text-slate-900">Dashboard</Link>
            <span>/</span>
            <span className="text-slate-900 font-medium">{display_name}</span>
          </div>
          <h1 className="text-2xl font-bold text-slate-900">{display_name}</h1>
        </div>'''
    src = src.replace(old_header2, new_header2)

    # ─── Fix 2: Add Link import if not present ───
    if 'Link' not in src and 'from "react-router-dom"' in src:
        src = src.replace(
            'import { useNavigate } from "react-router-dom";',
            'import { useNavigate, Link } from "react-router-dom";'
        )

    # ─── Fix 3: Add Plus import if not present ───
    if 'Plus' not in src and 'lucide-react' not in src:
        # Add lucide-react import
        src = src.replace(
            'import type {',
            'import { Plus } from "lucide-react";\nimport type {'
        )
    elif 'Plus' not in src and 'lucide-react' in src:
        # Add Plus to existing lucide-react import
        src = re.sub(
            r'import \{([^}]+)\} from "lucide-react";',
            lambda m: 'import {' + m.group(1) + ', Plus } from "lucide-react";' if 'Plus' not in m.group(1) else m.group(0),
            src
        )

    # ─── Fix 4: Unify container spacing ───
    src = src.replace('className="space-y-4"', 'className="space-y-6"')

    # ─── Fix 5: Better loading text ───
    src = src.replace(
        '<div className="text-center py-12 text-slate-400">Loading...</div>',
        f'<div className="text-center py-12 text-slate-400">Loading {display_name.lower()}s...</div>'
    )

    # ─── Fix 6: Unify search input placeholder ───
    src = src.replace(
        'placeholder="Search..."',
        f'placeholder="Search {display_name.lower()}s..."'
    )

    # ─── Fix 7: Add text-slate-900 to h1 ───
    src = src.replace('className="text-2xl font-bold"', 'className="text-2xl font-bold text-slate-900"')

    if src != orig:
        open(f, 'w').write(src)
        changed.append(f)

print(f"Changed {len(changed)} files")
for f in sorted(changed):
    print(f)

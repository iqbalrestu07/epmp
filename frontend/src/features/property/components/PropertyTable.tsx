import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import type { Property } from "../types";
import { Building2 } from "lucide-react";

const columnHelper = createColumnHelper<Property>();

const columns = [
  columnHelper.accessor("name", {
    header: "Name",
    cell: (info) => (
      <div className="flex items-center gap-2">
        <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center text-orange shrink-0">
          <Building2 size={16} />
        </div>
        <span className="font-semibold text-slate-800">{info.getValue() || "Unnamed Property"}</span>
      </div>
    ),
  }),
  columnHelper.accessor("property_type", {
    header: "Type",
    cell: (info) => {
      const val = info.getValue();
      const formatted = val ? String(val).replace(/_/g, " ") : "Unknown";
      return (
        <span className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-slate-100 text-slate-700 capitalize border border-slate-200">
          {formatted}
        </span>
      );
    },
  }),
  columnHelper.accessor("address", {
    header: "Address",
    cell: (info) => (
      <span className="text-slate-600 truncate max-w-xs block">
        {info.getValue() || <span className="text-slate-400">—</span>}
      </span>
    ),
  }),
  columnHelper.accessor("is_active", {
    header: "Status",
    cell: (info) => {
      const active = info.getValue() !== false;
      return (
        <span
          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold border ${
            active
              ? "bg-green-50 text-green-700 border-green-200"
              : "bg-red-50 text-red-700 border-red-200"
          }`}
        >
          {active ? "Active" : "Inactive"}
        </span>
      );
    },
  }),
];

interface PropertyTableProps {
  data: Property[];
  onRowClick?: (row: Property) => void;
}

export function PropertyTable({ data, onRowClick }: PropertyTableProps) {
  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">
      <table className="w-full text-sm text-left">
        <thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="px-5 py-3.5">
                  {flexRender(
                    header.column.columnDef.header,
                    header.getContext()
                  )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-slate-100">
          {table.getRowModel().rows.map((row) => (
            <tr
              key={row.id}
              className="hover:bg-orange/5 cursor-pointer transition-colors"
              onClick={() => onRowClick?.(row.original)}
            >
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className="px-5 py-4">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
          {table.getRowModel().rows.length === 0 && (
            <tr>
              <td colSpan={columns.length} className="px-5 py-12 text-center text-slate-400">
                <Building2 size={36} className="mx-auto mb-2 opacity-30" />
                <p className="font-medium text-slate-600">No properties found.</p>
                <p className="text-xs text-slate-400 mt-1">Get started by creating your first property.</p>
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

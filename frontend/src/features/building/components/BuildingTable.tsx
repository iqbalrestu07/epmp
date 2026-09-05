import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import type { Building } from "../types";
import type { Property } from "../../property/types";
import { Building2 } from "lucide-react";

const columnHelper = createColumnHelper<Building>();

interface BuildingTableProps {
  data: Building[];
  properties?: Property[];
  onRowClick?: (row: Building) => void;
}

export function BuildingTable({ data, properties = [], onRowClick }: BuildingTableProps) {
  const propertyNameById = new Map(properties.map((p) => [p.id, p.name]));

  const columns = [
    columnHelper.accessor("name", {
      header: "Name",
      cell: (info) => (
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center text-orange shrink-0">
            <Building2 size={16} />
          </div>
          <span className="font-semibold text-slate-800">{info.getValue() || "Unnamed Building"}</span>
        </div>
      ),
    }),
    columnHelper.accessor("property_id", {
      header: "Property",
      cell: (info) => (
        <span className="text-slate-600">
          {propertyNameById.get(info.getValue()) ?? "—"}
        </span>
      ),
    }),
    columnHelper.accessor("total_floors", {
      header: "Floors",
      cell: (info) => (
        <span className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700 border border-blue-200">
          {info.getValue()} floors
        </span>
      ),
    }),
  ];

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
                <p className="font-medium text-slate-600">No buildings found.</p>
                <p className="text-xs text-slate-400 mt-1">Create one to get started.</p>
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

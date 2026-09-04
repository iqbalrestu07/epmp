import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import { Layers, Eye, Pencil, Trash2 } from "lucide-react";
import type { Floor } from "../types";

const columnHelper = createColumnHelper<Floor>();

interface FloorTableProps {
  data: Floor[];
  onRowClick?: (row: Floor) => void;
  onEdit?: (row: Floor) => void;
  onDelete?: (row: Floor) => void;
}

export function FloorTable({ data, onRowClick, onEdit, onDelete }: FloorTableProps) {
  const columns = [
    columnHelper.accessor("name", {
      header: "Name",
      cell: (info) => (
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center flex-shrink-0">
            <Layers size={16} className="text-orange" />
          </div>
          <span className="font-medium">{info.getValue()}</span>
        </div>
      ),
    }),
    columnHelper.accessor("floor_number", {
      header: "Floor #",
      cell: (info) => (
        <span className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-black/5 text-sm font-semibold">
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("building_id", {
      header: "Building",
      cell: (info) => (
        <span className="text-sm text-black/50 font-mono truncate max-w-32 inline-block">
          {info.getValue() ? info.getValue().slice(0, 12) + "…" : "—"}
        </span>
      ),
    }),
    columnHelper.accessor("is_active", {
      header: "Status",
      cell: (info) => (
        <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${
          info.getValue()
            ? "bg-green-100 text-green-700"
            : "bg-red-100 text-red-700"
        }`}>
          {info.getValue() ? "Active" : "Inactive"}
        </span>
      ),
    }),
    columnHelper.display({
      id: "actions",
      header: "Actions",
      cell: (info) => (
        <div className="flex items-center gap-1">
          <button
            onClick={(e) => { e.stopPropagation(); onRowClick?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-black/5 text-black/50 hover:text-black transition-colors"
            title="View"
          >
            <Eye size={16} />
          </button>
          <button
            onClick={(e) => { e.stopPropagation(); onEdit?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-black/5 text-black/50 hover:text-orange transition-colors"
            title="Edit"
          >
            <Pencil size={16} />
          </button>
          <button
            onClick={(e) => { e.stopPropagation(); onDelete?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-black/5 text-black/50 hover:text-red-500 transition-colors"
            title="Delete"
          >
            <Trash2 size={16} />
          </button>
        </div>
      ),
    }),
  ];

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-xl border border-black/5 overflow-hidden bg-white">
      <table className="w-full text-sm">
        <thead className="border-b border-black/5 bg-black/[0.02]">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="px-4 py-3 text-left font-semibold text-black/60 text-xs uppercase tracking-wider">
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr
              key={row.id}
              className="border-b border-black/5 hover:bg-black/[0.02] cursor-pointer transition-colors"
              onClick={() => onRowClick?.(row.original)}
            >
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className="px-4 py-3">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
          {table.getRowModel().rows.length === 0 && (
            <tr>
              <td colSpan={columns.length} className="px-4 py-12 text-center text-black/30">
                <Layers size={32} className="mx-auto mb-2 opacity-50" />
                No floors found. Create one to get started.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

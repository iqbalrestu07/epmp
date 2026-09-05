import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import { DoorOpen, Eye, Pencil, Trash2, Users, DollarSign } from "lucide-react";
import type { Room } from "../types";
import type { Floor } from "../../floor/types";

const columnHelper = createColumnHelper<Room>();

interface RoomTableProps {
  data: Room[];
  floors?: Floor[];
  onRowClick?: (row: Room) => void;
  onEdit?: (row: Room) => void;
  onDelete?: (row: Room) => void;
}

export function RoomTable({ data, floors = [], onRowClick, onEdit, onDelete }: RoomTableProps) {
  const floorNameById = new Map(floors.map((f) => [f.id, f.name]));
  const columns = [
    columnHelper.accessor("name", {
      header: "Name",
      cell: (info) => (
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center flex-shrink-0">
            <DoorOpen size={16} className="text-orange" />
          </div>
          <span className="font-medium">{info.getValue()}</span>
        </div>
      ),
    }),
    columnHelper.accessor("capacity", {
      header: "Capacity",
      cell: (info) => (
        <span className="inline-flex items-center gap-1 text-sm text-slate-600">
          <Users size={14} className="text-slate-400" />
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("price", {
      header: "Price",
      cell: (info) => (
        <span className="inline-flex items-center gap-1 text-sm font-medium">
          <DollarSign size={14} className="text-slate-400" />
          {info.getValue().toLocaleString()}
        </span>
      ),
    }),
    columnHelper.accessor("floor_id", {
      header: "Floor",
      cell: (info) => (
        <span className="text-sm text-slate-600">
          {floorNameById.get(info.getValue()) ?? "—"}
        </span>
      ),
    }),
    columnHelper.accessor("is_available", {
      header: "Status",
      cell: (info) => (
        <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${
          info.getValue()
            ? "bg-green-50 text-green-700 border border-green-200"
            : "bg-red-50 text-red-700 border border-red-200"
        }`}>
          {info.getValue() ? "Available" : "Occupied"}
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
            className="p-1.5 rounded-lg hover:bg-slate-100 text-slate-500 hover:text-slate-900 transition-colors"
            title="View"
          >
            <Eye size={16} />
          </button>
          <button
            onClick={(e) => { e.stopPropagation(); onEdit?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-slate-100 text-slate-500 hover:text-orange transition-colors"
            title="Edit"
          >
            <Pencil size={16} />
          </button>
          <button
            onClick={(e) => { e.stopPropagation(); onDelete?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-slate-100 text-slate-500 hover:text-red-500 transition-colors"
            title="Delete"
          >
            <Trash2 size={16} />
          </button>
        </div>
      ),
    }),
  ];

  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-xl border border-slate-200 overflow-hidden bg-white">
      <table className="w-full text-sm text-left">
        <thead className="border-b border-slate-200 bg-slate-50/75">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="px-5 py-3.5">
                  {flexRender(header.column.columnDef.header, header.getContext())}
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
                No data found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

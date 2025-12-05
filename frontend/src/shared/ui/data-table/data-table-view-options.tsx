"use client"

import { Table } from "@tanstack/react-table"
import { Settings2 } from "lucide-react"

import { Button } from "@/shared/shadcn/ui/button"

interface DataTableViewOptionsProps<TData> {
  table: Table<TData>
}

export function DataTableViewOptions<TData>({
  table,
}: DataTableViewOptionsProps<TData>) {
  return (
    <div className="relative">
      <Button
        variant="outline"
        size="sm"
        className="ml-auto hidden h-8 lg:flex"
      >
        <Settings2 />
        Колонки
      </Button>
    </div>
  )
}

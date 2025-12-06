"use client";

import { useState } from "react";
import { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, Pencil, Trash } from "lucide-react";
import { DataTableColumnHeader } from "@/shared/ui/data-table/data-table-column-header";
import { Button } from "@/shared/shadcn/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/shadcn/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/shadcn/ui/alert-dialog";
import { useRouter } from "next/navigation";
import { useDeleteTransaction } from "@/features/transaction/hooks";
import type { Transaction } from "@/entities/transaction/model/transaction.schema";

const ActionsCell = ({ transaction }: { transaction: Transaction }) => {
  const router = useRouter();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const deleteMutation = useDeleteTransaction();

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync({ id: transaction.id });
      setShowDeleteDialog(false);
    } catch (error) {
      console.error("Failed to delete transaction:", error);
    }
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="h-8 w-8 p-0">
            <span className="sr-only">Открыть меню</span>
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>Действия</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => {
              router.push(`/transactions/edit/${transaction.id}`);
            }}
          >
            <Pencil className="mr-2 h-4 w-4" />
            Редактировать
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => setShowDeleteDialog(true)}
            className="text-destructive focus:text-destructive"
          >
            <Trash className="mr-2 h-4 w-4" />
            Удалить
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить транзакцию?</AlertDialogTitle>
            <AlertDialogDescription>
              Это действие нельзя отменить. Транзакция будет удалена навсегда.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="space-y-2 text-sm my-4">
            <p>
              <strong>Тип:</strong>{" "}
              {transaction.type === "Пополнение/Доход" ? "Доход" : "Расход"}
            </p>
            <p>
              <strong>Категория:</strong> {transaction.kategoria}
            </p>
            <p>
              <strong>Сумма:</strong>{" "}
              {new Intl.NumberFormat("ru-RU", {
                style: "currency",
                currency: "RUB",
              }).format(transaction.amount)}
            </p>
            <p>
              <strong>Описание:</strong> {transaction.description || "—"}
            </p>
          </div>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteMutation.isPending ? "Удаление..." : "Удалить"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
};

export const columns: ColumnDef<Transaction>[] = [
  {
    accessorKey: "date",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Дата" />
    ),
    cell: ({ row }) => {
      const timestamp = row.getValue("date") as number;
      return new Date(timestamp).toLocaleDateString("ru-RU", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
      });
    },
  },
  {
    accessorKey: "description",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Описание" />
    ),
  },
  {
    accessorKey: "kategoria",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Категория" />
    ),
  },
  {
    accessorKey: "amount",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Сумма" />
    ),
    cell: ({ row }) => {
      const amount = row.getValue("amount") as number;
      const type = row.getValue("type") as string;
      const formatted = new Intl.NumberFormat("ru-RU", {
        style: "currency",
        currency: "RUB",
      }).format(Math.abs(amount));
      const isIncome = type === "Пополнение/Доход";
      return (
        <span className={isIncome ? "text-green-600" : "text-red-600"}>
          {isIncome ? "+" : "-"}
          {formatted}
        </span>
      );
    },
  },
  {
    accessorKey: "type",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Тип" />
    ),
    cell: ({ row }) => {
      const type = row.getValue("type") as string;
      return type === "Пополнение/Доход" ? "Доход" : "Расход";
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const transaction = row.original;
      return <ActionsCell transaction={transaction} />;
    },
  },
];

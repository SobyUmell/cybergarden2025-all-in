"use client"
import { DataTable } from "@/shared/ui/data-table"
import { columns } from "../model/category.columns"

type Transaction = {
  id: string;
  date: string;
  description: string;
  category: string;
  amount: number;
  type: "income" | "expense";
};

const data: Transaction[] = [
  {
    id: "1",
    date: "2025-12-01",
    description: "Зарплата",
    category: "Доход",
    amount: 150000,
    type: "income",
  },
  {
    id: "2",
    date: "2025-12-02",
    description: "Продукты",
    category: "Питание",
    amount: -5200,
    type: "expense",
  },
  {
    id: "3",
    date: "2025-12-03",
    description: "Транспорт",
    category: "Транспорт",
    amount: -1500,
    type: "expense",
  },
  {
    id: "4",
    date: "2025-12-04",
    description: "Кафе",
    category: "Развлечения",
    amount: -2800,
    type: "expense",
  },
  {
    id: "5",
    date: "2025-12-05",
    description: "Фриланс",
    category: "Доход",
    amount: 35000,
    type: "income",
  },
  {
    id: "6",
    date: "2025-12-01",
    description: "Зарплата",
    category: "Доход",
    amount: 150000,
    type: "income",
  },
  {
    id: "7",
    date: "2025-12-02",
    description: "Продукты",
    category: "Питание",
    amount: -5200,
    type: "expense",
  },
  {
    id: "8",
    date: "2025-12-03",
    description: "Транспорт",
    category: "Транспорт",
    amount: -1500,
    type: "expense",
  },
  {
    id: "9",
    date: "2025-12-04",
    description: "Кафе",
    category: "Развлечения",
    amount: -2800,
    type: "expense",
  },
  {
    id: "10",
    date: "2025-12-05",
    description: "Фриланс",
    category: "Доход",
    amount: 35000,
    type: "income",
  },
  {
    id: "11",
    date: "2025-12-01",
    description: "Зарплата",
    category: "Доход",
    amount: 150000,
    type: "income",
  },
  {
    id: "12",
    date: "2025-12-02",
    description: "Продукты",
    category: "Питание",
    amount: -5200,
    type: "expense",
  },
  {
    id: "13",
    date: "2025-12-03",
    description: "Транспорт",
    category: "Транспорт",
    amount: -1500,
    type: "expense",
  },
  {
    id: "14",
    date: "2025-12-04",
    description: "Кафе",
    category: "Развлечения",
    amount: -2800,
    type: "expense",
  },
  {
    id: "15",
    date: "2025-12-05",
    description: "Фриланс",
    category: "Доход",
    amount: 35000,
    type: "income",
  },
];

export const CategoryTable = () => {
  return (
    <DataTable
      columns={columns}
      data={data}
      searchKey="description"
      searchPlaceholder="Поиск по описанию..."
    />
  )
}

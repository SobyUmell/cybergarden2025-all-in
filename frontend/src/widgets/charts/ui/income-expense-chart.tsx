"use client"

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"

const data = [
  { month: "Янв", income: 150000, expense: 98000 },
  { month: "Фев", income: 155000, expense: 102000 },
  { month: "Мар", income: 160000, expense: 95000 },
  { month: "Апр", income: 150000, expense: 110000 },
  { month: "Май", income: 170000, expense: 105000 },
  { month: "Июн", income: 165000, expense: 98000 },
]

export const IncomeExpenseChart = () => {
  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Доходы и расходы</h3>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey="income" fill="#82ca9d" name="Доходы" />
          <Bar dataKey="expense" fill="#ff8042" name="Расходы" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

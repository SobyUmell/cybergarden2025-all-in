"use client"

import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Legend, Line } from "recharts"

const data = [
  { month: "Янв", amount: 98000 },
  { month: "Фев", amount: 102000 },
  { month: "Мар", amount: 95000 },
  { month: "Апр", amount: 110000 },
  { month: "Май", amount: 105000 },
  { month: "Июн", amount: 98000 },
]

export const MonthlySpendingChart = () => {
  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Ежемесячные расходы</h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="amount" stroke="#8884d8" strokeWidth={2} name="Расходы" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}

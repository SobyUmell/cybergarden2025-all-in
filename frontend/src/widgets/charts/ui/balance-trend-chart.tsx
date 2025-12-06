"use client"

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"

const data = [
  { month: "Янв", balance: 50000 },
  { month: "Фев", balance: 53000 },
  { month: "Мар", balance: 68000 },
  { month: "Апр", balance: 58000 },
  { month: "Май", balance: 73000 },
  { month: "Июн", balance: 90000 },
]

export const BalanceTrendChart = () => {
  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Динамика баланса</h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="balance" stroke="#8884d8" strokeWidth={2} name="Баланс" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}

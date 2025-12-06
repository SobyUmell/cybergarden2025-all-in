"use client"

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"

const data = [
  { month: "Янв", rate: 34.7 },
  { month: "Фев", rate: 34.2 },
  { month: "Мар", rate: 40.6 },
  { month: "Апр", rate: 26.7 },
  { month: "Май", rate: 38.2 },
  { month: "Июн", rate: 40.6 },
]

export const SavingsRateChart = () => {
  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Процент сбережений</h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="rate" stroke="#00C49F" strokeWidth={2} name="Процент (%)" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}

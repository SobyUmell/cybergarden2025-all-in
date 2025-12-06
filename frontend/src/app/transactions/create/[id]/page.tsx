import { CreateTransactionForm } from "@/features/transaction/create-transaction/ui/create-transaction-form"

export default async function CreateTransactionPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params

  return (
    <div>
      <CreateTransactionForm id={id} />
    </div>
  )
}

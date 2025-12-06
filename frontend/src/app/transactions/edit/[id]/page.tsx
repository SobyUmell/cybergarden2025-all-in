import { EditTransactionClient } from "@/features/transaction/edit-transaction/ui/edit-transaction-client";

export default async function EditTransactionPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return (
    <div className="space-y-4">
      <h1 className="text-2xl md:text-3xl font-bold">
        Редактировать транзакцию
      </h1>
      <EditTransactionClient id={Number(id)} />
    </div>
  );
}

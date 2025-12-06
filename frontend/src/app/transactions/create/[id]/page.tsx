import { CreateTransactionForm } from "@/features/transaction/create-transaction/ui/create-transaction-form";

export default function CreateTransactionPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl md:text-3xl font-bold">Создать транзакцию</h1>
      <CreateTransactionForm />
    </div>
  );
}

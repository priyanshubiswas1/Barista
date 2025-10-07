import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { IoClose } from "react-icons/io5";
import { useOverlay } from "../../store/overlayStore";
import Input from "../Input";
import Button from "../Button";
import { addItem, fetchTablesByShopId } from "../../api/api";

const AddNewOrder = ({ header }) => {
  const { setOverlay } = useOverlay();

  const [formData, setFormData] = useState({
    name: "",
    categoryId: "",
    price: 0.0,
  });

  const [err, setErr] = useState({
    name: "",
    categoryId: "",
    price: 0.0,
  });

  const {
    data,
    isLoading,
    error: tableError,
  } = useQuery({
    queryKey: ["tables"],
    queryFn: () => fetchTablesByShopId(localStorage.getItem("shopId")),
  });

  const queryClient = useQueryClient();

  if (isLoading) return <p>Loading item...</p>;
  if (tableError) return <p>Error: {tableError.message}</p>;

  console.log(data);
  const tables = data.tables;

  const {
    mutate,
    isPending,
    isError,
    isSuccess,
    error: mutationError,
  } = useMutation({
    mutationFn: addItem,
    onSuccess: () => {
      alert("Item added✅");
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
    },
    onError: (error) => {
      console.log("Add Item failed: ", error);
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const newError = {
      name: formData.name ? "" : "Name is required",
      categoryId: formData.categoryId ? "" : "Category is required",
      price: formData.price ? "" : "Price is required",
    };

    setErr(newError);

    if (newError.name || newError.categoryId || newError.price) return;

    const payload = {
      ...formData,
      price: parseFloat(formData.price),
      shopId: localStorage.getItem("shopId"),
    };

    console.log(payload);
    mutate(payload);
  };

  return (
    <div className="flex flex-col w-full h-full">
      <div className="px-6 pt-6 pb-4 flex justify-between items-start">
        {header && <h2 className="font-semibold text-2xl">{header}</h2>}
        <button
          onClick={() => setOverlay("")}
          className="text-gray-600 hover:text-gray-900"
          aria-label="Close overlay"
        >
          <IoClose size={24} />
        </button>
      </div>

      <div className="w-full h-px bg-gray-200" />

      <form className="flex flex-col gap-2 px-6 pt-4">
        <div className="flex flex-col">
          <label className="text-xs font-medium text-gray-800 mb-1">
            Table
          </label>
          <select
            className="w-full rounded-md border px-3 py-2 bg-gray-50 border-gray-200 outline-none text-sm text-gray-800"
            name="categoryId"
            value={formData.categoryId}
            onChange={handleChange}
          >
            <option value="" disabled className="text-gray-400">
              Select table
            </option>
            {tables.map((table) => (
              <option key={table.tableId} value={table.tableId}>
                {table.tableNumber}
              </option>
            ))}
          </select>
          {err.categoryId && (
            <p className="mt-1 text-xs text-red-700">{err.categoryId}</p>
          )}
        </div>
      </form>

      <div className="flex justify-between px-6 mt-4">
        <Button
          onClick={() =>
            setFormData({
              name: "",
              categoryId: "",
              price: 0.0,
            })
          }
          variant="outline"
        >
          Cancel
        </Button>
        <Button onClick={handleSubmit} type="submit" variant="primary">
          {isPending ? "Saving..." : "Save"}
        </Button>
      </div>
    </div>
  );
};

export default AddNewOrder;

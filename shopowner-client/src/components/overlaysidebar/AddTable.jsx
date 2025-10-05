import React, { useState } from "react";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { IoClose } from "react-icons/io5";
import Input from "../Input";
import Button from "../Button";
import { createTableByShopId } from "../../api/api";

const AddTable = ({ header }) => {
  const queryClient = useQueryClient();
  const { setOverlay } = useOverlay();
  const [tableNumber, setTableNumber] = useState("");

  const createMutation = useMutation({
    mutationFn: createTableByShopId,
    onSuccess: () => {
      alert("Table added successfully");
      setTableNumber("");
      queryClient.invalidateQueries({ queryKey: ["tables"] });
      setOverlay("");
    },
    onError: (error) => {
      alert("Failed to add table: " + error.message);
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!tableNumber.trim()) {
      alert("Please enter a table number");
      return;
    }

    createMutation.mutate({
      shopId: localStorage.getItem("shopId"),
      tableNumber: tableNumber,
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-lg w-full max-w-md p-6 relative">
      <div className="flex justify-between items-center mb-6">
        <h2 className="font-semibold text-2xl">{header}</h2>
        <button
          onClick={() => setOverlay("")}
          className="text-gray-600 hover:text-gray-900"
          aria-label="Close popup"
        >
          <IoClose size={24} />
        </button>
      </div>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <Input
          label="Table Number"
          type="text"
          name="tableNumber"
          placeholder="Enter table number"
          value={tableNumber}
          onChange={(e) => setTableNumber(e.target.value)}
          required
        />

        <Button
          type="submit"
          variant="primary"
          disabled={createMutation.isPending}
        >
          {createMutation.isPending ? "Adding..." : "Add"}
        </Button>
      </form>
    </div>
  );
};

export default AddTable;

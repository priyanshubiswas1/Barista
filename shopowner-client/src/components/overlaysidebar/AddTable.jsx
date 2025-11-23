import React, { useState } from "react";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { IoClose } from "react-icons/io5";
import Input from "../Input";
import Button from "../Button";
import { createTableByShopId } from "../../api/api";
import { toast } from "react-toastify";

const AddTable = ({ header }) => {
  const queryClient = useQueryClient();
  const { setOverlay } = useOverlay();
  const [tableNumber, setTableNumber] = useState("");

  const createMutation = useMutation({
    mutationFn: createTableByShopId,
    onSuccess: () => {
      toast("Table added successfully");
      setTableNumber("");
      queryClient.invalidateQueries({ queryKey: ["tables"] });
      setOverlay("");
    },
    onError: (error) => {
      toast.error("Failed to add table: " + error.message);
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!tableNumber.trim()) {
      toast.error("Please enter a table number");
      return;
    }

    createMutation.mutate({
      shopId: localStorage.getItem("shopId"),
      tableNumber: tableNumber,
    });
  };

  return (
    <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg p-8 relative">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="font-bold text-2xl text-gray-900">{header}</h2>
          <p className="text-sm text-gray-500 mt-1">Enter table details</p>
        </div>
        <button
          onClick={() => setOverlay("")}
          className="text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg p-2 transition-all"
          aria-label="Close popup"
        >
          <IoClose size={24} />
        </button>
      </div>

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <Input
          label="Table Number"
          type="text"
          name="tableNumber"
          placeholder="Enter table number"
          value={tableNumber}
          onChange={(e) => setTableNumber(e.target.value)}
          required
        />

        <div className="flex gap-3 pt-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => setOverlay("")}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="primary"
            disabled={createMutation.isPending}
          >
            {createMutation.isPending ? "Adding..." : "Add Table"}
          </Button>
        </div>
      </form>
    </div>
  );
};

export default AddTable;

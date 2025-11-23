import React, { useEffect, useState } from "react";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { IoClose } from "react-icons/io5";
import Input from "../Input";
import Button from "../Button";
import { updateTableByShopId } from "../../api/api";
import { toast } from "react-toastify";

const UpdateTable = ({ header, table }) => {
  const queryClient = useQueryClient();
  const { setOverlay } = useOverlay();
  const [tableNumber, setTableNumber] = useState(table?.tableNumber || "");
  const [tableId, setTableId] = useState(table?.tableId || "");

  useEffect(() => {
    setTableNumber(table?.tableNumber || "");
    setTableId(table?.tableId || "");
  }, [table]);

  const createMutation = useMutation({
    mutationFn: updateTableByShopId,
    onSuccess: () => {
      toast("Table updated successfully");
      setTableNumber("");
      queryClient.invalidateQueries({ queryKey: ["tables"] });
      setOverlay("");
    },
    onError: (error) => {
      toast.error("Failed to update table: " + error.message);
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!tableId) {
      toast.error("Missing table information");
      return;
    }

    if (!tableNumber.trim()) {
      toast.error("Please enter a table number");
      return;
    }

    createMutation.mutate({
      tableId: tableId,
      shopId: localStorage.getItem("shopId"),
      tableNumber: tableNumber,
    });
  };

  return (
    <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden relative">
      {/* Header with gradient */}
      <div className="bg-gradient-to-br from-emerald-600 to-emerald-700 text-white p-8 relative overflow-hidden">
        <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16" />
        <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full -ml-12 -mb-12" />
        
        <div className="relative z-10 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="bg-white rounded-full h-12 w-12 flex items-center justify-center shadow-lg">
              <span className="text-2xl">☕</span>
            </div>
            <div>
              <h2 className="font-bold text-2xl">{header}</h2>
              <p className="text-emerald-100 text-sm mt-1">Modify table details</p>
            </div>
          </div>
          <button
            onClick={() => setOverlay("")}
            className="text-white hover:bg-white/20 rounded-lg p-2 transition-all"
            aria-label="Close popup"
          >
            <IoClose size={24} />
          </button>
        </div>
      </div>

      {/* Form Content */}
      <div className="p-8 bg-gradient-to-b from-gray-50 to-white">
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
              {createMutation.isPending ? "Updating..." : "Update Table"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default UpdateTable;

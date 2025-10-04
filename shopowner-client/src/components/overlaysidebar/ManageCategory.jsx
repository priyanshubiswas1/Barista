import { IoClose } from "react-icons/io5";
import { LuCirclePlus } from "react-icons/lu";
import { AiOutlineDelete } from "react-icons/ai";
import { PiPencilSimpleLight } from "react-icons/pi";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createCategory,
  deleteCategoryByCategoryId,
  fetchCategoryItemsByShopId,
} from "../../api/api";
import Input from "../Input";
import { useState } from "react";
import Button from "../Button";

const CategoryItem = ({ name, id, count, onDelete }) => {
  return (
    <div className="w-full pt-3">
      <div className="flex justify-between items-center px-6 pb-3">
        <div className="flex text-xl text-gray-900 gap-3">
          <div>{name}</div>
          <div className="bg-gray-100 border border-gray-200 w-7 h-7 rounded flex items-center justify-center">
            <h2>{count}</h2>
          </div>
        </div>
        <div className="flex text-gray-500 text-xl gap-3">
          <PiPencilSimpleLight />
          <AiOutlineDelete onClick={() => onDelete(id)} />
        </div>
      </div>
      <div className="w-full h-px bg-gray-200" />
    </div>
  );
};

const ManageCategory = ({ header }) => {
  const { setOverlay } = useOverlay();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["categories"],
    queryFn: () => fetchCategoryItemsByShopId(localStorage.getItem("shopId")),
  });

  const [addCategory, setAddCategory] = useState(false);
  const [addItemName, setAddItemName] = useState("");

  const createMutation = useMutation({
    mutationFn: createCategory,
    onSuccess: () => {
      alert("Category added✅");
      setAddCategory(false);
      queryClient.invalidateQueries({ queryKey: ["categories"] });
    },
    onError: () => {},
  });

  const deleteMutation = useMutation({
    mutationFn: deleteCategoryByCategoryId,
    onSuccess: () => {
      alert("Category deleted");
      queryClient.invalidateQueries({ queryKey: ["categories"] });
    },
    onError: () => {
      alert("Deletion failed");
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    createMutation.mutate({
      name: addItemName,
      shopId: localStorage.getItem("shopId"),
    });
  };

  const handleDeleteCategory = (id) => {
    deleteMutation.mutate({
      categoryId: id,
      shopId: localStorage.getItem("shopId"),
    });
  };

  if (isLoading) return <p>Loading items...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const categories = data.categories;

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
      <div className="flex flex-col w-full h-full">
        <div className="pt-5">
          {categories.map((category, index) => (
            <CategoryItem
              key={index}
              name={category.name}
              id={category.categoryId}
              count={category.items.length}
              onDelete={handleDeleteCategory}
            />
          ))}
        </div>
        <div className="pt-3">
          {addCategory ? (
            <form
              onSubmit={handleSubmit}
              className="flex justify-between items-center gap-2 px-6 pt-4"
            >
              <Input
                label="New Category Name"
                type="text"
                name="add category"
                placeholder="Enter category name"
                value={addItemName}
                onChange={(e) => setAddItemName(e.target.value)}
              />
              <Button type="submit" variant="primary">
                {createMutation.isPending ? "Adding..." : "Add"}
              </Button>
            </form>
          ) : (
            <button
              onClick={() => setAddCategory(true)}
              className="flex items-center gap-3 px-6"
            >
              <LuCirclePlus />
              <h2 className="text-lg font-medium">Add New Category</h2>
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default ManageCategory;

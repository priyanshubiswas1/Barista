import { X, PencilSimple, TrashSimple, PlusCircle } from "phosphor-react";
import { LuCirclePlus } from "react-icons/lu";
import { useOverlay } from "../../store/overlayStore";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createCategory,
  deleteCategoryByCategoryId,
  fetchCategoryItemsByShopId,
  updateCategory,
} from "../../api/api";
import Input from "../Input";
import { useState } from "react";
import Button from "../Button";
import { toast } from "react-toastify";

const CategoryItem = ({ name, id, count, onDelete, onEdit }) => {
  return (
    <div className="w-full">
      <div className="flex justify-between items-center px-8 py-3 border-b-1 border-gray-200">
        <div className="flex  gap-3">
          <div className="text-sm text-gray-600 font-medium">{name}</div>
          <div className="bg-gray-100 border border-gray-200 w-5 h-5 rounded flex items-center justify-center">
            <div className="text-xs text-gray-700 font-medium">{count}</div>
          </div>
        </div>
        <div className="flex text-gray-400 gap-3">
          <PencilSimple
            size={16}
            className="hover:text-gray-800 cursor-pointer"
            onClick={() => onEdit(id, name)}
          />
          <TrashSimple
            size={16}
            className="hover:text-red-800 cursor-pointer"
            onClick={() => onDelete(id)}
          />
        </div>
      </div>
    </div>
  );
};

const ManageCategory = ({ header }) => {
  const { setOverlay } = useOverlay();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["categoryItems"],
    queryFn: () => fetchCategoryItemsByShopId(localStorage.getItem("shopId")),
  });

  const [addCategory, setAddCategory] = useState(false);
  const [addItemName, setAddItemName] = useState("");
  const [editingCategoryId, setEditingCategoryId] = useState(null);
  const [editItemName, setEditItemName] = useState("");

  const createMutation = useMutation({
    mutationFn: createCategory,
    onSuccess: () => {
      toast("Category added✅");
      setAddCategory(false);
      queryClient.invalidateQueries({ queryKey: ["categoryItems"] });
    },
    onError: () => {},
  });

  const deleteMutation = useMutation({
    mutationFn: deleteCategoryByCategoryId,
    onSuccess: () => {
      toast("Category deleted");
      queryClient.invalidateQueries({ queryKey: ["categoryItems"] });
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
    },
    onError: () => {
      toast.error("Deletion failed");
    },
  });

  const updateMutation = useMutation({
    mutationFn: updateCategory,
    onSuccess: () => {
      toast("Category updated✅");
      setEditingCategoryId(null);
      setEditItemName("");
      queryClient.invalidateQueries({ queryKey: ["categoryItems"] });
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
    },
    onError: () => {
      toast.error("Update failed");
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

  const handleEditCategory = (id, name) => {
    setEditingCategoryId(id);
    setEditItemName(name);
  };

  const handleUpdateSubmit = (e) => {
    e.preventDefault();
    updateMutation.mutate({
      name: editItemName,
      shopId: localStorage.getItem("shopId"),
      categoryId: editingCategoryId,
    });
  };

  const handleCancelEdit = () => {
    setEditingCategoryId(null);
    setEditItemName("");
  };

  if (isLoading) return <p>Loading items...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const categories = data.categories;

  return (
    <div className="flex flex-col w-full h-full">
      {/* Header with gradient */}
      <div className="bg-gradient-to-br from-emerald-600 to-emerald-700 text-white px-8 py-6 relative overflow-hidden shadow-lg">
        <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16" />
        <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full -ml-12 -mb-12" />
        
        <div className="relative z-10 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="bg-white rounded-full h-12 w-12 flex items-center justify-center shadow-lg">
              <span className="text-2xl">☕</span>
            </div>
            <div>
              {header && <h2 className="font-bold text-2xl">{header}</h2>}
              <p className="text-emerald-100 text-sm mt-1">Organize your menu</p>
            </div>
          </div>
          <button
            onClick={() => setOverlay("")}
            className="text-white hover:bg-white/20 rounded-lg p-2 transition-all"
            aria-label="Close overlay"
          >
            <X size={24} weight="bold" />
          </button>
        </div>
      </div>
      <div className="flex flex-col w-full h-full">
        <div className="pt-6">
          {categories.map((category, index) => (
            <div key={index}>
              {editingCategoryId === category.categoryId ? (
                <form
                  onSubmit={handleUpdateSubmit}
                  className="flex justify-between items-center gap-3 px-8 py-2"
                >
                  <Input
                    label="New Category Name"
                    type="text"
                    name="edit category"
                    placeholder="Enter category name"
                    value={editItemName}
                    onChange={(e) => setEditItemName(e.target.value)}
                  />
                  <div className="flex gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleCancelEdit}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" variant="primary">
                      {updateMutation.isPending ? "Saving..." : "Save"}
                    </Button>
                  </div>
                </form>
              ) : (
                <CategoryItem
                  name={category.name}
                  id={category.categoryId}
                  count={category.items.length}
                  onDelete={handleDeleteCategory}
                  onEdit={handleEditCategory}
                />
              )}
            </div>
          ))}
        </div>
        <div className="pt-3">
          {addCategory ? (
            <form
              onSubmit={handleSubmit}
              className="flex justify-between items-center gap-3 px-8 py-2"
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
              className="flex items-center gap-2 px-8 cursor-pointer"
            >
              <PlusCircle />
              <p className="text-sm font-semibold">Add New Category</p>
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default ManageCategory;

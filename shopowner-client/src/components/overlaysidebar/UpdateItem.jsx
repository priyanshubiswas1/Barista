import { useState, useEffect, useRef } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { X } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import Input from "../Input";
import Button from "../Button";
import {
  addItem,
  updateItem,
  fetchCategoriesByShopId,
  createCategory,
} from "../../api/api";
import Radio from "../Radio";
import ImageUploader from "../ImageUploader";
import SearchableSelect from "../SearchableSelect";
import { toast } from "react-toastify";
import CircleLoading from "../CircleLoading";

const UpdateItem = ({ header, item }) => {
  const { setOverlay } = useOverlay();
  const isEditing = !!item;

  const [formData, setFormData] = useState({
    name: "",
    categoryId: "",
    price: "0.0",
  });

  const [selected, setSelected] = useState("veg");
  const [selectedCategoryName, setSelectedCategoryName] = useState("");

  const [err, setErr] = useState({
    name: "",
    categoryId: "",
    price: "0.0",
  });

  const {
    data,
    isLoading,
    error: categoryError,
  } = useQuery({
    queryKey: ["categories"],
    queryFn: () => fetchCategoriesByShopId(localStorage.getItem("shopId")),
  });

  const queryClient = useQueryClient();
  const initializedItemIdRef = useRef(null);

  // Populate form when item is provided, or reset when item is null
  // Only initialize when item changes, not when data refetches
  useEffect(() => {
    if (item && data?.categories) {
      const currentItemId = item.itemId;
      // Only initialize if this is a different item or hasn't been initialized yet
      if (initializedItemIdRef.current !== currentItemId) {
        const categories = data.categories;
        const nameToId = Object.fromEntries(
          categories.map((c) => [c.name, c.categoryId])
        );

        setFormData({
          name: item.itemName || "",
          categoryId: nameToId[item.category] || "",
          price: item.price || "0.0",
        });
        setSelected(item.veg ? "veg" : "nonveg");
        setSelectedCategoryName(item.category || "");
        initializedItemIdRef.current = currentItemId;
      }
    } else if (!item) {
      // Reset form when no item (add mode)
      setFormData({
        name: "",
        categoryId: "",
        price: "0.0",
      });
      setSelected("veg");
      setSelectedCategoryName("");
      initializedItemIdRef.current = null;
    }
  }, [item, data]);

  const {
    mutate,
    isPending,
    isError,
    isSuccess,
    error: mutationError,
  } = useMutation({
    mutationFn: isEditing ? updateItem : addItem,
    onSuccess: () => {
      toast(isEditing ? "Item has been updated" : "Item has been added");
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
      setOverlay("");
    },
    onError: (error) => {
      console.log(`${isEditing ? "Update" : "Add"} Item failed: `, error);
      toast.error(
        `Failed to ${isEditing ? "update" : "add"} item. Please try again.`
      );
    },
  });

  const createCategoryMutation = useMutation({
    mutationFn: createCategory,
    onSuccess: (res) => {
      const newCat = res?.category || res;
      toast("Category added✅");
      queryClient.invalidateQueries({ queryKey: ["categories"] });
      queryClient.invalidateQueries({ queryKey: ["categoryItems"] });
      setSelectedCategoryName(newCat.name);
      setFormData((prev) => ({ ...prev, categoryId: newCat.categoryId }));
    },
    onError: () => {
      toast("Failed to add category");
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

    const payload = isEditing
      ? {
          itemId: item.itemId,
          shopId: localStorage.getItem("shopId"),
          itemName: formData.name,
          categoryId: formData.categoryId,
          price: parseFloat(formData.price),
          veg: selected === "veg",
        }
      : {
          ...formData,
          name: formData.name,
          price: parseFloat(formData.price),
          shopId: localStorage.getItem("shopId"),
          veg: selected === "veg",
        };

    console.log(payload);
    mutate(payload);
  };

  if (isLoading) return <CircleLoading text="Loading item..." />;
  if (categoryError) return <p>Error: {categoryError.message}</p>;

  const categories = data.categories;

  const nameToId = Object.fromEntries(
    categories.map((c) => [c.name, c.categoryId])
  );
  const categoryOptions = categories.map((c) => c.name);

  return (
    <div className="flex flex-col w-full h-full bg-gradient-to-b from-gray-50 to-white">
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
              {header && <h2 className="text-2xl font-bold">{header}</h2>}
              <p className="text-emerald-100 text-sm mt-1">
                Update item details
              </p>
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

      {/* Form Content - Scrollable */}
      <div className="flex-1 overflow-y-auto">
        <form className="flex flex-col gap-4 px-8 py-8">
          {/* Image Upload Section */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-200">
            <h3 className="text-sm font-semibold text-gray-700 mb-4 uppercase tracking-wide">
              Item Image
            </h3>
            <ImageUploader label="Image" />
          </div>

          {/* Basic Details Section */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-200">
            <h3 className="text-sm font-semibold text-gray-700 mb-5 uppercase tracking-wide">
              Basic Details
            </h3>
            <div className="flex flex-col">
              <Input
                label="Name"
                type="text"
                name="name"
                placeholder="Enter item name"
                helperText={err.name}
                error={!!err.name}
                value={formData.name}
                onChange={handleChange}
              />
              <SearchableSelect
                label="Category"
                placeholder="Search or select category"
                initialOptions={categoryOptions}
                defaultValue={selectedCategoryName}
                onSelect={(name) => {
                  setSelectedCategoryName(name);
                  setFormData((prev) => ({
                    ...prev,
                    categoryId: nameToId[name] || "",
                  }));
                }}
                onChange={(value) => {
                  // If the input value doesn't match the selected category name,
                  // clear the categoryId
                  if (value !== selectedCategoryName) {
                    setSelectedCategoryName("");
                    setFormData((prev) => ({
                      ...prev,
                      categoryId: "",
                    }));
                  }
                }}
                onAdd={(name) => {
                  createCategoryMutation.mutate({
                    name,
                    shopId: localStorage.getItem("shopId"),
                  });
                }}
                helperText={err.categoryId}
                error={!!err.categoryId}
              />
            </div>
          </div>

          {/* Pricing & Type Section */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-200">
            <h3 className="text-sm font-semibold text-gray-700 mb-5 uppercase tracking-wide">
              Pricing & Type
            </h3>
            <div className="flex flex-col">
              <Input
                label="Price"
                type="number"
                name="price"
                placeholder="0.00"
                prefix="₹"
                helperText={err.price}
                error={!!err.price}
                value={formData.price}
                onChange={handleChange}
              />
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">
                  Food Type
                </label>
                <div className="flex gap-8">
                  <Radio
                    name="food"
                    value="veg"
                    label="Veg"
                    checked={selected === "veg"}
                    onChange={() => setSelected("veg")}
                  />
                  <Radio
                    name="food"
                    value="nonveg"
                    label="Non-Veg"
                    checked={selected === "nonveg"}
                    onChange={() => setSelected("nonveg")}
                  />
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>

      {/* Footer Actions - Fixed */}
      <div className="border-t-2 border-gray-200 bg-gradient-to-r from-gray-50 to-white px-8 py-6 flex justify-end gap-3 shadow-lg">
        <Button
          onClick={() => {
            setOverlay("");
          }}
          variant="outline"
        >
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          type="submit"
          variant="primary"
          disabled={isPending}
        >
          {isPending ? "Updating..." : "Update Item"}
        </Button>
      </div>
    </div>
  );
};

export default UpdateItem;

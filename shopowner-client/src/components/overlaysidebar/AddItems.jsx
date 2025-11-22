import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { X } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import Input from "../Input";
import Button from "../Button";
import {
  addItem,
  fetchCategoriesByShopId,
  createCategory,
} from "../../api/api";
import Radio from "../Radio";
import ImageUploader from "../ImageUploader";
import SearchableSelect from "../SearchableSelect";
import { toast } from "react-toastify";

const AddItems = ({ header }) => {
  const { setOverlay } = useOverlay();

  const [formData, setFormData] = useState({
    name: "",
    categoryId: "",
    price: 0.0,
  });

  const [selected, setSelected] = useState("veg");
  const [selectedCategoryName, setSelectedCategoryName] = useState("");
  const [imageFile, setImageFile] = useState(null);

  const [err, setErr] = useState({
    name: "",
    categoryId: "",
    price: 0.0,
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

  const {
    mutate,
    isPending,
    isError,
    isSuccess,
    error: mutationError,
  } = useMutation({
    mutationFn: addItem,
    onSuccess: () => {
      toast("Item has been added");
      queryClient.invalidateQueries({ queryKey: ["menuItems"] });
      setFormData({
        name: "",
        categoryId: "",
        price: 0.0,
      });
      setSelectedCategoryName("");
      setImageFile(null);
    },
    onError: (error) => {
      console.error("Add Item failed: ", error);
      toast.error(
        error?.response?.data?.error || "Failed to add item. Please try again."
      );
    },
  });

  const createCategoryMutation = useMutation({
    mutationFn: createCategory,
    onSuccess: (res) => {
      const newCat = res?.category || res;
      toast("Category added✅");
      setSelectedCategoryName(newCat.name);
      setFormData((prev) => ({ ...prev, categoryId: newCat.categoryId }));
      queryClient.invalidateQueries({ queryKey: ["categories"] });
      queryClient.invalidateQueries({ queryKey: ["categoryItems"] });
    },
    onError: () => {
      toast("Failed to add category");
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleImageChange = (file) => {
    setImageFile(file);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const newError = {
      name: formData.name ? "" : "Name is required",
      categoryId: formData.categoryId ? "" : "Category is required",
      price: formData.price ? "" : "Price is required",
    };

    setErr(newError);

    if (newError.name || newError.categoryId || newError.price) return;

    const payload = {
      name: formData.name,
      shopId: localStorage.getItem("shopId"),
      categoryId: formData.categoryId,
      price: parseFloat(formData.price),
      veg: selected === "veg",
      imageFile: imageFile,
    };

    console.log(payload);
    mutate(payload);
  };

  if (isLoading) return <p>Loading item...</p>;
  if (categoryError) return <p>Error: {categoryError.message}</p>;

  const categories = data.categories;

  const nameToId = Object.fromEntries(
    categories.map((c) => [c.name, c.categoryId])
  );
  const categoryOptions = categories.map((c) => c.name);

  return (
    <div className="flex flex-col w-full h-full">
      <div className="px-8 py-4 flex justify-between items-center">
        {header && <h2>{header}</h2>}
        <button
          onClick={() => setOverlay("")}
          className="text-gray-600 hover:text-gray-800 cursor-pointer"
          aria-label="Close overlay"
        >
          <X size={16} weight="bold" />
        </button>
      </div>

      <div className="w-full h-px bg-gray-200" />

      <form className="flex flex-col gap-6 px-8 py-6">
        <ImageUploader
          label="Image"
          details="Supports JPEGs & PNGS under 5mb."
          onChange={handleImageChange}
          value={imageFile}
        />
        <div className="flex flex-col gap-2">
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
          <div className="flex gap-20">
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
      </form>

      <div className="flex justify-between px-8">
        <Button
          onClick={() => {
            setFormData({
              name: "",
              categoryId: "",
              price: 0.0,
            });
            setSelectedCategoryName("");
            setImageFile(null);
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
          {isPending ? "Saving..." : "Save"}
        </Button>
      </div>
    </div>
  );
};

export default AddItems;

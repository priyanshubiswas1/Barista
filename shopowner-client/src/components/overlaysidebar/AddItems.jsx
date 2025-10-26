import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { X } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import Input from "../Input";
import Button from "../Button";
import { addItem, fetchCategoriesByShopId } from "../../api/api";
import Radio from "../Radio";
import ImageUploader from "../ImageUploader";
import SearchableSelect from "../SearchableSelect";

const AddItems = ({ header }) => {
  const { setOverlay } = useOverlay();

  const [formData, setFormData] = useState({
    name: "",
    categoryId: "",
    price: 0.0,
  });

  const [selected, setSelected] = useState("veg");

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

  if (isLoading) return <p>Loading item...</p>;
  if (categoryError) return <p>Error: {categoryError.message}</p>;

  const categories = data.categories;

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
          ></SearchableSelect>

          <div className="flex flex-col">
            <label className="text-xs font-medium text-gray-800 mb-1">
              Category
            </label>
            <select
              className="w-full rounded-md border px-3 py-2 bg-gray-50 border-gray-200 outline-none text-sm text-gray-800"
              name="categoryId"
              value={formData.categoryId}
              onChange={handleChange}
            >
              <option value="" disabled className="text-gray-400">
                Select category
              </option>
              {categories.map((cat) => (
                <option key={cat.categoryId} value={cat.categoryId}>
                  {cat.name}
                </option>
              ))}
            </select>
            {err.categoryId && (
              <p className="mt-1 text-xs text-red-700">{err.categoryId}</p>
            )}
          </div>

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

export default AddItems;

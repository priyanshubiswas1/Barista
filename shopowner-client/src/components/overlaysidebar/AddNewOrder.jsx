import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { IoAdd, IoRemove, IoClose } from "react-icons/io5";
import { IoMdArrowDropright, IoMdArrowDropdown } from "react-icons/io";
import { IoEllipseSharp, IoTriangleSharp } from "react-icons/io5";
import { Plus } from "phosphor-react";
import { useOverlay } from "../../store/overlayStore";
import { PiNotepadBold, PiQrCodeBold } from "react-icons/pi";
import Button from "../Button";
import {
  createOrder,
  fetchItemsByShopId,
  fetchTablesByShopId,
} from "../../api/api";
import Searchbar from "../Searchbar";
import { toast } from "react-toastify";

const AddNewOrder = ({ header }) => {
  const { setOverlay } = useOverlay();
  const [foodType, setFoodType] = useState("");
  const [tableId, setTableId] = useState("");
  const [formData, setFormData] = useState({
    tableId: "",
    orderItems: [],
  });
  const [err, setErr] = useState({
    tableId: "",
    orderItems: "",
  });

  const [viewOrder, setViewOrder] = useState(false);

  const [orderItems, setOrderItems] = useState([]);
  const [expandedCategories, setExpandedCategories] = useState({});

  const handleViewOrderClick = () => {
    if (!tableId) {
      setErr((prev) => ({ ...prev, tableId: "Table ID is required" }));
      toast.error("Please select a table.");
      return;
    }
    setViewOrder(true);
  };

  const queryClient = useQueryClient();

  const {
    data,
    isLoading,
    error: tableError,
  } = useQuery({
    queryKey: ["tables"],
    queryFn: () => fetchTablesByShopId(localStorage.getItem("shopId")),
  });

  const mutation = useMutation({
    mutationFn: createOrder,
    onSuccess: (data, variables) => {
      toast("Order added✅");
      queryClient.invalidateQueries({ queryKey: ["sessions"] });
      queryClient.invalidateQueries({
        queryKey: ["order-items", variables.tableId],
      });
      setFormData({ tableId: "", orderItems: [] });
      setOrderItems([]);
      setTableId("");
      setOverlay("");
    },
    onError: (error) => {
      console.log("Add Order failed: ", error);
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name === "tableId") {
      setTableId(value);
    }
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const newError = {
      tableId: tableId ? "" : "Table ID is required",
      orderItems:
        formData.orderItems.length > 0 ? "" : "OrderItems are required",
    };

    setErr(newError);

    if (newError.tableId || newError.orderItems) return;

    const payload = {
      ...formData,
      shopId: localStorage.getItem("shopId"),
      tableId: tableId,
      specialNotes: "",
    };

    mutation.mutate(payload);
  };

  if (isLoading) return <p>Loading tables...</p>;
  if (tableError) return <p>Error: {tableError.message}</p>;

  const tables = data.tables;

  const TableSelect = () => (
    <div>
      <label className="text-xs font-medium text-gray-800 mb-1">Table</label>
      <div className="flex items-start gap-2">
        <select
          className="flex-1 rounded-md border px-3 py-2 bg-gray-50 border-gray-200 outline-none text-sm text-gray-800"
          name="tableId"
          value={tableId}
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
        <button
          className="aspect-square h-[39px] border border-gray-200 rounded-md flex justify-center items-center text-3xl"
          type="button"
        >
          <PiQrCodeBold />
        </button>
      </div>
      {err.tableId && (
        <p className="text-xs text-red-700 mt-1">{err.tableId}</p>
      )}
    </div>
  );

  const VegNonVegSelect = ({ value, onChange }) => {
    const toggle = (type) => {
      if (value === type) onChange("");
      else onChange(type);
    };

    return (
      <div className="flex gap-4">
        <button
          type="button"
          onClick={() => toggle("veg")}
          className={`text-md font-medium rounded-4xl border flex items-center justify-around px-3 py-1 gap-2 ${
            value === "veg"
              ? "bg-green-100 border-green-600 text-green-900"
              : "border-gray-200 text-gray-800"
          }`}
        >
          <div className="text-green-800 border border-green-800 rounded-sm p-1">
            <IoEllipseSharp size={8} />
          </div>
          Veg
        </button>

        <button
          type="button"
          onClick={() => toggle("non-veg")}
          className={`text-md font-medium rounded-4xl border flex items-center justify-around px-3 py-1 gap-2 ${
            value === "non-veg"
              ? "bg-red-100 border-red-600 text-red-900"
              : "border-gray-200 text-gray-800"
          }`}
        >
          <div className="text-red-800 border border-red-800 rounded-sm p-1">
            <IoTriangleSharp size={8} />
          </div>
          Non Veg
        </button>
      </div>
    );
  };

  const Menu = ({
    orderItems,
    setOrderItems,
    expandedCategories,
    setExpandedCategories,
    foodType,
  }) => {
    const { data, isLoading, error } = useQuery({
      queryKey: ["menuItems"],
      queryFn: () => fetchItemsByShopId(localStorage.getItem("shopId")),
    });

    const toggleCategory = (category) => {
      setExpandedCategories((prev) => ({
        ...prev,
        [category]: !prev[category],
      }));
    };

    if (isLoading) return <p>Loading menu...</p>;
    if (error) return <p>Error: {error.message}</p>;

    let items = data.items;

    if (foodType === "veg") items = items.filter((item) => item.veg === true);
    else if (foodType === "non-veg")
      items = items.filter((item) => item.veg === false);

    const groupedItems = items.reduce((acc, item) => {
      if (!acc[item.category]) acc[item.category] = [];
      acc[item.category].push(item);
      return acc;
    }, {});

    const Item = ({ item }) => {
      const currentItem = orderItems.find((i) => i.itemId === item.itemId);
      const count = currentItem ? currentItem.quantity : 0;

      const increment = () => {
        setOrderItems((prev) => {
          const exists = prev.find((i) => i.itemId === item.itemId);
          const newOrderItems = exists
            ? prev.map((i) =>
                i.itemId === item.itemId
                  ? {
                      ...i,
                      itemName: item.itemName,
                      quantity: i.quantity + 1,
                      priceAtOrderTime: item.price,
                    }
                  : i
              )
            : [
                ...prev,
                {
                  itemId: item.itemId,
                  itemName: item.itemName,
                  quantity: 1,
                  priceAtOrderTime: item.price,
                },
              ];

          setFormData((prevForm) => ({
            ...prevForm,
            orderItems: newOrderItems,
          }));

          return newOrderItems;
        });
      };

      const decrement = () => {
        setOrderItems((prev) => {
          const exists = prev.find((i) => i.itemId === item.itemId);
          if (!exists) return prev;

          const newOrderItems =
            exists.quantity === 1
              ? prev.filter((i) => i.itemId !== item.itemId)
              : prev.map((i) =>
                  i.itemId === item.itemId
                    ? { ...i, quantity: i.quantity - 1 }
                    : i
                );

          setFormData((prevForm) => ({
            ...prevForm,
            orderItems: newOrderItems,
          }));

          return newOrderItems;
        });
      };

      return (
        <div className="flex justify-between border border-gray-300 rounded-xl p-3">
          <div className="flex flex-col justify-between">
            <div className="flex items-center gap-2">
              <div
                className={`${
                  item.veg
                    ? "text-green-800 border-green-800"
                    : "text-red-800 border-red-800"
                } border rounded-sm w-4 h-4 flex items-center justify-center`}
              >
                {item.veg ? (
                  <IoEllipseSharp size={8} />
                ) : (
                  <IoTriangleSharp size={8} />
                )}
              </div>
              <h2 className="font-medium text-sm text-gray-800">
                {item.itemName}
              </h2>
            </div>
            <h2 className="text-gray-600 text-sm">₹ {item.price}</h2>
          </div>

          <div
            className="relative w-20 h-20 rounded-md flex items-center justify-center"
            style={{
              backgroundImage: `url(${item.image})`,
              backgroundSize: "cover",
              backgroundPosition: "center",
            }}
          >
            {count === 0 ? (
              <button
                type="button"
                onClick={increment}
                className="absolute bottom-[-8px] bg-rose-700 text-white rounded-md w-8 h-8 flex items-center justify-center shadow-md"
              >
                <IoAdd size={20} />
              </button>
            ) : (
              <div className="absolute bottom-[-8px] bg-rose-100 border border-rose-300 rounded-md flex justify-between items-center px-1 py-1 gap-2 shadow-sm w-full">
                <button
                  type="button"
                  onClick={decrement}
                  className="text-rose-700 text-xl leading-none"
                >
                  <IoRemove />
                </button>
                <span className="text-rose-700 text-sm font-medium">
                  {count}
                </span>
                <button
                  type="button"
                  onClick={increment}
                  className="text-rose-700 text-xl leading-none"
                >
                  <IoAdd />
                </button>
              </div>
            )}
          </div>
        </div>
      );
    };

    const Category = ({ category, items }) => {
      const isExpanded = expandedCategories[category] ?? false;

      return (
        <div className="flex flex-col gap-2">
          <button
            type="button"
            onClick={() => toggleCategory(category)}
            className="bg-gray-200 border border-gray-300 rounded-md py-1 px-3 text-black font-medium flex justify-between items-center w-full"
          >
            <div className="flex gap-4 items-center">
              {category}
              <div className="bg-white w-5 h-5 rounded-md flex items-center justify-center text-sm">
                {items.length}
              </div>
            </div>
            {isExpanded ? (
              <IoMdArrowDropdown size={25} className="text-gray-800" />
            ) : (
              <IoMdArrowDropright size={25} className="text-gray-800" />
            )}
          </button>

          {isExpanded && (
            <div className="flex flex-col gap-2">
              {items.map((item) => (
                <Item key={item.itemId} item={item} />
              ))}
            </div>
          )}
        </div>
      );
    };

    return (
      <div className="pt-5 flex-grow flex flex-col gap-3 overflow-auto max-h-[400px]">
        {Object.entries(groupedItems).map(([category, categoryItems], index) =>
          categoryItems.length > 0 ? (
            <Category key={index} category={category} items={categoryItems} />
          ) : null
        )}
      </div>
    );
  };

  const AddOrder = () => {
    const navigate = useNavigate();
    return (
      <form className="flex flex-col gap-2 px-6 pt-4">
        <div className="flex flex-col gap-5">
          <TableSelect />
          <div className="-mx-6 h-px bg-gray-200" />
          <div className="flex gap-3">
            <Searchbar />
            <Button
              variant="primary"
              icon={PiNotepadBold}
              type="button"
              onClick={() => navigate("/menu")}
            >
              Menu
            </Button>
          </div>
          <VegNonVegSelect
            value={foodType}
            onChange={(val) => setFoodType(val)}
          />
          <Menu
            orderItems={orderItems}
            setOrderItems={setOrderItems}
            expandedCategories={expandedCategories}
            setExpandedCategories={setExpandedCategories}
            foodType={foodType}
          />
          {err.orderItems && (
            <p className="mt-1 text-xs text-red-700">{err.orderItems}</p>
          )}
        </div>
      </form>
    );
  };

  const ViewOrder = () => {
    const increment = (item) => {
      setOrderItems((prev) =>
        prev.map((i) =>
          i.itemId === item.itemId ? { ...i, quantity: i.quantity + 1 } : i
        )
      );
    };

    const decrement = (item) => {
      setOrderItems((prev) =>
        prev
          .map((i) =>
            i.itemId === item.itemId ? { ...i, quantity: i.quantity - 1 } : i
          )
          .filter((i) => i.quantity > 0)
      );
    };

    const BasicDetails = () => {
      return (
        <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-2">
          <h2 className="font-semibold">Basic Details</h2>
          <div className="flex justify-between">
            <p>Order Status</p>
            <div className="bg-blue-400 text-white border rounded-md px-4">
              new
            </div>
          </div>
          <div className="flex justify-between">
            <p>Table No.</p>
            <p>
              {tables.find((t) => t.tableId === tableId)?.tableNumber || "-"}
            </p>
          </div>
          <div className="flex justify-between">
            <p>Date</p>
            <p>{new Date().toLocaleDateString()}</p>
          </div>
        </div>
      );
    };

    const OrderDetails = () => {
      return (
        <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-3">
          <h2 className="font-semibold">Order Details</h2>

          {orderItems.map((item) => (
            <div
              key={item.itemId}
              className="flex justify-between items-center"
            >
              <p>{item.itemName}</p>

              <div className="flex items-center gap-3">
                <div className="relative w-22 flex items-center justify-center">
                  <div className="bg-rose-100 border border-rose-300 rounded-md flex justify-between items-center px-2 py-1 gap-2 shadow-sm w-full">
                    <button
                      type="button"
                      onClick={() => decrement(item)}
                      className="text-rose-700 text-xl leading-none"
                    >
                      <IoRemove />
                    </button>
                    <span className="text-rose-700 text-sm font-medium">
                      {item.quantity}
                    </span>
                    <button
                      type="button"
                      onClick={() => increment(item)}
                      className="text-rose-700 text-xl leading-none"
                    >
                      <IoAdd />
                    </button>
                  </div>
                </div>

                <p>₹ {item.priceAtOrderTime * item.quantity}</p>
              </div>
            </div>
          ))}

          <div className="mt-4">
            <Button
              variant="outline"
              icon={Plus}
              type="button"
              onClick={() => setViewOrder(false)}
            >
              Add More Items
            </Button>
          </div>
        </div>
      );
    };

    const BillingDetails = () => {
      const total = orderItems.reduce(
        (acc, item) => acc + item.priceAtOrderTime * item.quantity,
        0
      );
      return (
        <div className="flex flex-col justify-between border border-gray-200 rounded-xl p-3 gap-2">
          <h2 className="font-semibold">Billing Details</h2>
          <div className="flex justify-between">
            <p>Subtotal</p>
            <p>₹ {total}</p>
          </div>
        </div>
      );
    };
    return (
      <div className="w-full h-full flex flex-col p-4 gap-3">
        <BasicDetails />
        <OrderDetails />
        <BillingDetails />
      </div>
    );
  };

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex flex-col w-full h-full overflow-auto">
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
                <p className="text-emerald-100 text-sm mt-1">
                  Create a new order
                </p>
              </div>
            </div>
            <button
              onClick={() => setOverlay("")}
              className="text-white hover:bg-white/20 rounded-lg p-2 transition-all"
              aria-label="Close overlay"
            >
              <IoClose size={24} />
            </button>
          </div>
        </div>

        {viewOrder ? (
          <ViewOrder
            orderItems={orderItems}
            setOrderItems={setOrderItems}
            setViewOrder={setViewOrder}
          />
        ) : (
          <AddOrder />
        )}
      </div>

      <div className="bg-emerald-700 text-white flex justify-around items-center py-8">
        <h2>
          Items added | ₹{" "}
          {orderItems.reduce(
            (acc, item) => acc + item.priceAtOrderTime * item.quantity,
            0
          )}
        </h2>
        {!viewOrder ? (
          <button onClick={handleViewOrderClick}>View Order</button>
        ) : (
          <button onClick={handleSubmit} disabled={mutation.isLoading}>
            {mutation.isLoading ? "Submitting..." : "Place Order"}
          </button>
        )}
      </div>
    </div>
  );
};

export default AddNewOrder;

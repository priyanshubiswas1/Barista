import { useState, useRef, useEffect } from "react";
import { CaretDown, PlusCircle } from "phosphor-react";

const SearchableSelect = ({
  initialOptions = [],
  onAdd,
  label,
  placeholder = "Select or type...",
  helperText = "",
  error = false,
  disabled = false,
}) => {
  const [options, setOptions] = useState(initialOptions);
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const wrapperRef = useRef(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const filtered = options.filter((opt) =>
    opt.toLowerCase().includes(query.toLowerCase())
  );

  const handleSelect = (value) => {
    setQuery(value);
    setIsOpen(false);
  };

  const handleAdd = () => {
    if (query.trim() && !options.includes(query)) {
      setOptions([...options, query]);
      onAdd?.(query);
      handleSelect(query);
    }
  };

  return (
    <div className="w-full relative" ref={wrapperRef}>
      {/* Label */}
      {label && (
        <label
          className={`block text-xs font-medium mb-1 ${
            disabled ? "text-gray-400" : "text-gray-800"
          }`}
        >
          {label}
        </label>
      )}

      {/* Input container */}
      <div
        className={`relative flex items-center w-full rounded-md px-3 h-8
          ${
            disabled
              ? "bg-gray-100 border border-gray-200 cursor-not-allowed"
              : error
              ? "border border-red-700 focus-within:ring-2 focus-within:ring-red-700"
              : "border border-gray-200 focus-within:ring-2 focus-within:ring-gray-800"
          }
          bg-gray-50`}
        onClick={() => {
          if (!disabled) setIsOpen(true);
        }}
      >
        <input
          type="text"
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setIsOpen(true);
          }}
          onFocus={() => setIsOpen(true)}
          placeholder={placeholder}
          disabled={disabled}
          className="flex-1 bg-transparent outline-none text-sm text-gray-800 placeholder:text-gray-400 disabled:placeholder:text-gray-300 h-full"
        />

        {/* CaretDown icon */}
        <CaretDown
          size={14}
          weight="bold"
          className={`absolute right-3 top-1/2 -translate-y-1/2 ${
            isOpen ? "text-gray-800" : "text-gray-400"
          }`}
        />
      </div>

      {/* Dropdown */}
      {isOpen && !disabled && (
        <div
          className={`absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-md shadow-sm max-h-48 overflow-auto p-1`}
        >
          {filtered.length > 0 ? (
            filtered.map((opt) => (
              <div
                key={opt}
                onClick={() => handleSelect(opt)}
                className="px-3 py-2 text-sm text-gray-600 rounded hover:bg-gray-100 hover:text-gray-800 cursor-pointer"
              >
                {opt}
              </div>
            ))
          ) : (
            <div
              onClick={handleAdd}
              className="flex items-center px-3 py-2 text-sm text-gray-600 rounded hover:bg-gray-100 hover:text-gray-800 cursor-pointer"
            >
              <PlusCircle
                size={16}
                weight="regular"
                className="mr-2 text-gray-800"
              />
              Add "{query}"
            </div>
          )}
        </div>
      )}

      {/* Helper / error text */}
      <p
        className={`mt-1 text-xs h-4 ${
          error ? "text-red-700" : disabled ? "text-gray-300" : "text-gray-500"
        }`}
      >
        {helperText || " "}
      </p>
    </div>
  );
};

export default SearchableSelect;

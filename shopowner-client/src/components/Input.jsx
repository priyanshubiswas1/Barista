export default function Input({
  label,
  placeholder,
  value,
  onChange,
  helperText = "",
  error = false,
  disabled = false,
  prefix,
  type = "text",
  autocomplete = "",
  ...props
}) {
  return (
    <div className="w-full">
      {label && (
        <label
          className={`block text-xs font-medium ${
            disabled ? "text-gray-400" : "text-gray-800"
          } mb-1`}
        >
          {label}
        </label>
      )}

      <div
        className={`flex items-center w-full h-8 rounded-md border px-3 py-2 bg-gray-50
          ${
            disabled
              ? "bg-gray-100 border-gray-200 cursor-not-allowed"
              : error
              ? "border-red-700"
              : "border-gray-200"
          }`}
      >
        {prefix && <span className="mr-2 text-gray-800 text-sm">{prefix}</span>}

        <input
          type={type}
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          disabled={disabled}
          autoComplete={autocomplete}
          className={`flex-1 w-full bg-transparent outline-none text-sm text-gray-800
            placeholder:text-gray-400
            disabled:placeholder:text-gray-300
            focus-visible:outline-none`}
          {...props}
        />
      </div>

      {/* Always reserve space for helper text */}
      <p
        className={`mt-1 text-xs h-4 ${
          error ? "text-red-700" : disabled ? "text-gray-300" : "text-gray-500"
        }`}
      >
        {helperText || " "}
      </p>
    </div>
  );
}

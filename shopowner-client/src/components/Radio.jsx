const Radio = ({ name, value, label, checked, onChange }) => (
  <label className="flex items-center space-x-2 cursor-pointer select-none">
    <input
      type="radio"
      name={name}
      value={value}
      checked={checked}
      onChange={onChange}
      className="hidden"
    />
    <div
      className={`w-4 h-4 rounded-full border-2 flex items-center justify-center ${
        checked ? "border-gray-800" : "border-gray-300"
      }`}
    >
      {checked && <div className="w-2 h-2 rounded-full bg-gray-800"></div>}
    </div>
    {label && <span className="text-sm text-gray-800">{label}</span>}
  </label>
);

export default Radio;

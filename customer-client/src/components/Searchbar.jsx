import { useState } from "react";
import { MagnifyingGlass } from "phosphor-react";

const Searchbar = () => {
  const [query, setQuery] = useState("");

  const handleChange = (e) => setQuery(e.target.value);

  return (
    <div className="group flex items-center gap-2 h-10 bg-gray-200 text-gray-800 border border-gray-300 rounded-md shadow-xs px-3 focus-within:ring-2 focus-within:ring-gray-800">
      <MagnifyingGlass
        size={16}
        className="text-gray-400 group-focus-within:text-gray-800"
      />
      <input
        type="text"
        value={query}
        onChange={handleChange}
        placeholder="Search Item"
        className="bg-transparent outline-none border-none text-gray-800 placeholder-gray-400 w-full text-sm h-full"
      />
    </div>
  );
};

export default Searchbar;

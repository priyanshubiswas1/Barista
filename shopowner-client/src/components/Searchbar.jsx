import { useState } from "react";
// import { IoIosSearch } from "react-icons/io";
import { MagnifyingGlass } from "phosphor-react";

const Searchbar = () => {
  const [query, setQuery] = useState("");

  const handleChange = (e) => setQuery(e.target.value);

  return (
    <div className="flex items-center gap-2 bg-white text-gray-800 border border-gray-200 hover:border-gray-300 rounded-md shadow-xs disabled:shadow-none px-3 py-1.5 focus-within:ring-2 focus-within:ring-gray-800">
      <MagnifyingGlass size={16} color="currentColor" />
      <input
        type="text"
        value={query}
        onChange={handleChange}
        placeholder="Search"
        className="bg-transparent outline-none border-none text-gray-800 placeholder-gray-400 w-full"
      />
    </div>
  );
};

export default Searchbar;

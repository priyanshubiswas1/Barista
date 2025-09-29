import { useState } from "react";
import { IoIosSearch } from "react-icons/io";

const Searchbar = () => {
  const [query, setQuery] = useState("");

  const handleChange = (e) => setQuery(e.target.value);

  return (
    <div className="flex items-center gap-3 bg-white border border-gray-200 rounded-lg px-2 py-2">
      <IoIosSearch className="text-gray-400 text-2xl flex-shrink-0" />
      <input
        type="text"
        value={query}
        onChange={handleChange}
        placeholder="Search"
        className="bg-transparent outline-none border-none text-gray-600 placeholder-gray-400"
      />
    </div>
  );
};

export default Searchbar;

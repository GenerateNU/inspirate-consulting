import { Search } from "lucide-react";

export function SearchBar({
  value,
  onChangeText,
  placeholder,
  className = "",
}: {
  value: string;
  onChangeText: (text: string) => void;
  placeholder?: string;
  className?: string;
}) {
  return (
    <div className={`relative w-full ${className}`}>
      <Search
        size={16}
        aria-hidden="true"
        className="pointer-events-none absolute top-1/2 left-[13px] -translate-y-1/2 text-gray-400"
      />
      <input
        type="search"
        aria-label={placeholder}
        value={value}
        onChange={(e) => onChangeText(e.target.value)}
        placeholder={placeholder}
        className="h-11 w-full rounded-none border border-[#76777d] bg-white py-0 pr-[13px] pl-[37px] font-inter text-sm text-text placeholder:text-gray-400"
      />
    </div>
  );
}

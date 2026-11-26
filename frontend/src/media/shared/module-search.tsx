import { Search } from "lucide-react";

import { Input } from "@/components/ui/input";

interface ModuleSearchProps {
  value: string;
  onChange: (value: string) => void;
}

// the pill-shaped "Search specific modules" box at the top of a Modules list
export default function ModuleSearch({ value, onChange }: ModuleSearchProps) {
  return (
    <div className="relative w-[264px] max-md:w-full">
      <Search
        aria-hidden="true"
        className="pointer-events-none absolute top-1/2 left-3 size-3.5 -translate-y-1/2 text-warm-grey"
        strokeWidth={2}
      />
      <Input
        type="search"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Search specific modules"
        aria-label="Search modules"
        className="h-8 rounded-full border-cream bg-white pr-3 pl-8 font-host text-xs text-charcoal shadow-none placeholder:text-warm-grey md:text-xs"
      />
    </div>
  );
}

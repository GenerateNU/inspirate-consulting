import { useListMedia } from "../api/endpoints";

export type ConnectToType = "" | "module" | "essay" | "extracurricular";

export type ConnectToValue = {
  type: ConnectToType;
  itemIds: string[];
};

type ConnectToSelectProps = {
  value: ConnectToValue;
  onChange: (value: ConnectToValue) => void;
};

// PLACEHOLDER: essays are per-student, and it's undecided how this works when assigning to many students
const PLACEHOLDER_ESSAYS = [
  { id: "essay-1", label: "Georgetown University: Supplemental" },
  { id: "essay-2", label: "Amherst College: Supplemental" },
  { id: "essay-3", label: "Northeastern University: Supplemental" },
];

// PLACEHOLDER: replace once there's an extracurriculars list endpoint/hook
const PLACEHOLDER_EXTRACURRICULARS = [
  { id: "ec-1", label: "Acapella" },
  { id: "ec-2", label: "Soccer" },
  { id: "ec-3", label: "Mock Trial" },
  { id: "ec-4", label: "Baseball" },
];

function toggle(list: string[], item: string) {
  return list.includes(item) ? list.filter((x) => x !== item) : [...list, item];
}

function OptionList({
  options,
  selected,
  onToggle,
  note,
}: {
  options: { id: string; label: string }[];
  selected: string[];
  onToggle: (id: string) => void;
  note?: string;
}) {
  return (
    <div className="flex flex-col gap-1 border p-2">
      {note && <p className="text-xs text-gray-500">{note}</p>}
      {options.map((option) => (
        <label
          key={option.id}
          className="flex items-center justify-between text-sm"
        >
          {option.label}
          <input
            type="checkbox"
            checked={selected.includes(option.id)}
            onChange={() => onToggle(option.id)}
          />
        </label>
      ))}
    </div>
  );
}

export default function ConnectToSelect({
  value,
  onChange,
}: ConnectToSelectProps) {
  // Real data: videos/modules
  const mediaQuery = useListMedia(
    { limit: 100, offset: 0 },
    { query: { enabled: value.type === "module" } },
  );
  const mediaOptions = Array.isArray(mediaQuery.data)
    ? mediaQuery.data.map((media) => ({ id: media.id, label: media.title }))
    : [];

  const toggleItem = (id: string) =>
    onChange({ ...value, itemIds: toggle(value.itemIds, id) });

  return (
    <div className="flex flex-col gap-2">
      <label className="text-sm">
        Connect to
        <select
          value={value.type}
          // Reset selected items when switching type
          onChange={(e) =>
            onChange({ type: e.target.value as ConnectToType, itemIds: [] })
          }
          className="mt-1 block w-full border px-2 py-1"
        >
          <option value="">Select an option</option>
          <option value="module">Module</option>
          <option value="essay">Essay</option>
          <option value="extracurricular">Extracurricular</option>
        </select>
      </label>

      {value.type === "module" &&
        (mediaQuery.isLoading ? (
          <p className="text-sm">Loading videos...</p>
        ) : (
          <OptionList
            options={mediaOptions}
            selected={value.itemIds}
            onToggle={toggleItem}
          />
        ))}

      {value.type === "essay" && (
        <OptionList
          options={PLACEHOLDER_ESSAYS}
          selected={value.itemIds}
          onToggle={toggleItem}
          note="Placeholder essays"
        />
      )}

      {value.type === "extracurricular" && (
        <OptionList
          options={PLACEHOLDER_EXTRACURRICULARS}
          selected={value.itemIds}
          onToggle={toggleItem}
          note="Placeholder extracurriculars"
        />
      )}
    </div>
  );
}

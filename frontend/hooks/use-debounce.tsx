import { useEffect, useState } from "react";
import { type FetchStudentsParams } from "../src/api/models";

interface DebounceProps {
  Input: string;
  setInput: React.Dispatch<React.SetStateAction<string>>;
}

export function useDebounce({ Input, setInput }: DebounceProps) {
  type Timer = ReturnType<typeof setTimeout>;
  const [timer, setTimer] = useState<Timer>();
  const delay: number = 8000;
  useEffect(() => {
    const newTimer = setTimeout(() => {
      setInput(Input);
    }, delay);

    clearTimeout(timer);
    setTimer(newTimer);
  }, [Input]);

  const params: FetchStudentsParams = { name: Input };
  return { Input: Input, params: params };
}

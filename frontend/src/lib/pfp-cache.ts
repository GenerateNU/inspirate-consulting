import { useEffect, useState } from "react";

const CACHE_NAME = "pfps";
const KEY_HEADER = "X-Pfp-Key";

const cacheRequest = (userId: string) => `/__pfp-cache/${userId}`;

const openCache = () =>
  "caches" in window ? caches.open(CACHE_NAME) : Promise.resolve(null);

interface Entry {
  key: string;
  src: Promise<string | null>;
  resolved?: string | null;
}

function makeEntry(key: string, src: Promise<string | null>): Entry {
  const entry: Entry = { key, src };
  void src.then((s) => (entry.resolved = s));
  return entry;
}

const entries = new Map<string, Entry>();

function replaceEntry(userId: string, entry: Entry | undefined) {
  const previous = entries.get(userId);
  if (previous) {
    void previous.src.then((src) => {
      if (src?.startsWith("blob:")) URL.revokeObjectURL(src);
    });
  }
  if (entry) entries.set(userId, entry);
  else entries.delete(userId);
}

async function load(userId: string, key: string, url: string | null) {
  try {
    const cache = await openCache();
    const cached = await cache?.match(cacheRequest(userId));
    if (cached?.headers.get(KEY_HEADER) === key) {
      return URL.createObjectURL(await cached.blob());
    }

    if (!url) return null;
    const response = await fetch(url);
    if (!response.ok) throw new Error(`pfp fetch failed: ${response.status}`);
    const blob = await response.blob();
    await cache?.put(
      cacheRequest(userId),
      new Response(blob, {
        headers: { "Content-Type": blob.type, [KEY_HEADER]: key },
      }),
    );
    return URL.createObjectURL(blob);
  } catch {
    return url;
  }
}

export async function seedPfpCache(userId: string, key: string, blob: Blob) {
  const src = URL.createObjectURL(blob);
  replaceEntry(userId, { key, src: Promise.resolve(src), resolved: src });
  try {
    const cache = await openCache();
    await cache?.put(
      cacheRequest(userId),
      new Response(blob, {
        headers: { "Content-Type": blob.type, [KEY_HEADER]: key },
      }),
    );
  } catch {
    // Best effort
  }
}

// Returns an image src for the user's profile picture, downloading it only
// when nothing is cached for them or their pfp_key has changed. `url` is the
// presigned URL from the API, used only on a cache miss.
export function usePfpSrc(
  userId: string | undefined,
  key: string | null | undefined,
  url: string | null | undefined,
): string | null {
  const current = userId ? entries.get(userId) : undefined;
  const [src, setSrc] = useState<{ key: string; src: string | null } | null>(
    null,
  );

  useEffect(() => {
    if (!userId) return;
    if (!key) {
      // Picture removed
      replaceEntry(userId, undefined);
      void openCache().then((cache) => cache?.delete(cacheRequest(userId)));
      return;
    }
    let entry = entries.get(userId);
    if (entry?.key !== key) {
      entry = makeEntry(key, load(userId, key, url ?? null));
      replaceEntry(userId, entry);
    }
    let cancelled = false;
    void entry.src.then((s) => {
      if (!cancelled) setSrc({ key, src: s });
    });
    return () => {
      cancelled = true;
    };
  }, [userId, key, url]);

  if (!key) return null;
  if (src?.key === key) return src.src;
  return current?.key === key ? (current.resolved ?? null) : null;
}

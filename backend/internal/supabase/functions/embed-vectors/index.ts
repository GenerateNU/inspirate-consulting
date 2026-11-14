import "jsr:@supabase/functions-js/edge-runtime.d.ts";
import { createClient } from "jsr:@supabase/supabase-js@2";

// Called by the notify_embed_vectors() trigger (pg_net) when a chat message or todo item is
// inserted, or its text changes. Embeds the text with Cloudflare

const CLOUDFLARE_ACCOUNT_ID = Deno.env.get("CLOUDFLARE_ACCOUNT_ID")!;
const CLOUDFLARE_API_TOKEN = Deno.env.get("CLOUDFLARE_API_TOKEN")!;
const SUPABASE_URL = Deno.env.get("SUPABASE_URL")!;
const SUPABASE_SERVICE_ROLE_KEY = Deno.env.get("SUPABASE_SERVICE_ROLE_KEY")!;

const MODEL = "@cf/qwen/qwen3-embedding-0.6b";
const DIMENSIONS = 1024;

const supabase = createClient(SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY);

// Which column holds the text to embed, and which vector_relationship column links back to the row.
const TABLES: Record<string, { textColumn: string; link: "message" | "todo" }> = {
  chat_messages: { textColumn: "message", link: "message" },
  todo_items: { textColumn: "todo_description", link: "todo" },
};

async function embed(text: string): Promise<number[]> {
  const res = await fetch(
    `https://api.cloudflare.com/client/v4/accounts/${CLOUDFLARE_ACCOUNT_ID}/ai/run/${MODEL}`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${CLOUDFLARE_API_TOKEN}`,
      },
      body: JSON.stringify({ text }),
    },
  );
  if (!res.ok) {
    throw new Error(`Cloudflare embedding failed (${res.status}): ${await res.text()}`);
  }

  // Response: { success, errors, result: { shape: [1, 1024], data: [[...]] } }
  const body = await res.json();
  const embedding = body?.result?.data?.[0];
  if (!Array.isArray(embedding) || embedding.length !== DIMENSIONS) {
    throw new Error(`Unexpected embedding shape: ${JSON.stringify(body?.result?.shape)}`);
  }
  return embedding;
}

async function upsertVector(embedding: number[], link: "message" | "todo", id: string) {
  const { error } = await supabase.rpc("upsert_vector", {
    // pgvector text format "[0.1,0.2,...]", cast to vector inside the SQL function
    p_embedding: JSON.stringify(embedding),
    p_message_id: link === "message" ? id : null,
    p_todo_id: link === "todo" ? id : null,
  });
  if (error) {
    // 23503 = foreign key violation: the message/todo was deleted before its embedding finished.
    // Nothing to store, and the rollback left no orphaned vector, so this isn't a failure.
    if (error.code === "23503") {
      console.log(`skipping ${link} ${id}: deleted before embedding finished`);
      return;
    }
    throw new Error(`upsert_vector failed: ${error.message}`);
  }
}

async function checkConnection(): Promise<boolean> {
  const res = await fetch(
    `https://api.cloudflare.com/client/v4/accounts/${CLOUDFLARE_ACCOUNT_ID}/ai/models/search?search=qwen3-embedding`,
    { headers: { "Authorization": `Bearer ${CLOUDFLARE_API_TOKEN}` } },
  );
  return res.ok;
}

Deno.serve(async (req) => {
  // Health/connection check
  if (req.method === "GET") {
    try {
      const healthy = await checkConnection();
      return new Response(
        JSON.stringify({ connected: healthy, model: MODEL }),
        { status: healthy ? 200 : 503, headers: { "Content-Type": "application/json" } },
      );
    } catch (err) {
      return new Response(JSON.stringify({ connected: false, error: String(err) }), {
        status: 503,
        headers: { "Content-Type": "application/json" },
      });
    }
  }

  try {
    const payload = await req.json();
    const { type, table, record, old_record } = payload;

    const config = TABLES[table];
    if (!config) {
      return new Response(`Unknown table: ${table}`, { status: 400 });
    }

    switch (type) {
      case "INSERT":
      case "UPDATE": {
        const text = record[config.textColumn];
        
        if (type === "UPDATE" && old_record?.[config.textColumn] === text) {
          return new Response("unchanged", { status: 200 });
        }
        if (typeof text !== "string" || text.trim() === "") {
          return new Response("no text to embed", { status: 200 });
        }

        const embedding = await embed(text);
        await upsertVector(embedding, config.link, record.id);
        break;
      }

      default:
        return new Response(`Unknown event type: ${type}`, { status: 400 });
    }

    return new Response("ok", { status: 200 });
  } catch (err) {
    console.error(err);
    return new Response(String(err), { status: 500 });
  }
});

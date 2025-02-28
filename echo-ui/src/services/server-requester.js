"use server";

import { cookies } from "next/headers";

// TODO: move to an environment variable
const BASE_URL = "http://localhost:8000";

export async function getPresignedUrl({ filename, eventType, contentType }) {
  const c = await cookies();
  const authToken = c.get("accessToken");
  console.log("authToken", authToken);

  const response = await fetch(`${BASE_URL}/uploads/presigned-url`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: authToken?.value,
    },
    body: JSON.stringify({ filename, eventType, contentType }),
    credentials: "include",
  });

  return await response.json();
}

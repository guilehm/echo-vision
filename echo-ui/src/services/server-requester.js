"use server";

import { cookies } from "next/headers";
import { API_BASE_URL } from "@/settings";

export async function serverRequester(url, options = {}, next = {}) {
  const { method = "GET", cache = "no-store", params, body } = options;

  if (params) {
    const urlParams = new URLSearchParams(params);
    url += `?${urlParams.toString()}`;
  }
  const requestURL = API_BASE_URL + "/" + url;

  const c = await cookies();

  const h = new Headers();
  h.append("Content-Type", "application/json");
  h.append("Cookie", c.toString());
  h.append("Accept", "application/json");

  const accessToken = c.get("accessToken")?.value;
  if (accessToken) {
    h.append("Authorization", accessToken);
  }

  const config = {
    method,
    headers: h,
    credentials: "include",
    ...(body && { body: JSON.stringify(body) }),
    ...(next ? { next } : { cache }),
  };
  return fetch(requestURL, config);
}

export async function getPresignedUrl({ filename, eventType, contentType }) {
  const c = await cookies();
  const authToken = c.get("accessToken");

  const response = await fetch(`${API_BASE_URL}/uploads/presigned-url`, {
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

export async function createEvent({
  eventType,
  subType,
  filename,
  filepath,
  contentType,
  filesize,
}) {
  const c = await cookies();
  const authToken = c.get("accessToken");

  const response = await fetch(`${API_BASE_URL}/events`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: authToken?.value,
    },
    body: JSON.stringify({
      eventType,
      subType,
      filename,
      filepath,
      contentType,
      filesize,
    }),
    credentials: "include",
  });
  return await response.json();
}

export async function getOwnEvents() {
  const c = await cookies();
  const authToken = c.get("accessToken");

  const response = await fetch(`${API_BASE_URL}/events`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: authToken?.value,
    },
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch own events");
  }

  return await response.json();
}

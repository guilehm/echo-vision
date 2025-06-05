"use client";

const BASE_URL = "http://localhost:8000";

// TODO: move to an environment variable

export async function signIn({ email, password }) {
  const response = await fetch(`${BASE_URL}/users/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
    credentials: "include",
  });
  return await response.json();
}

export async function signUp({ firstName, lastName, email, password }) {
  const response = await fetch(`${BASE_URL}/users/`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ firstName, lastName, email, password }),
    credentials: "include",
  });
  return await response.json();
}

export async function uploadS3File({ file, presignedURL }) {
  const response = await fetch(presignedURL, {
    method: "PUT",
    headers: { "Content-Type": file.type },
    body: file,
  });
  if (!response.ok) {
    throw new Error("failed to upload file");
  }
  return response;
}

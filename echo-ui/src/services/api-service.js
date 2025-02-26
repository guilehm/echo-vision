const BASE_URL = "http://localhost:8000";

export async function signIn({ email, password }) {
  const response = await fetch(`${BASE_URL}/users/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
    credentials: "include",
  });
  return await response.json();
}

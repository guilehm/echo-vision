export async function clientRequester(url, options = {}, next = {}) {
  const { method = "GET", cache = "no-store", params, body } = options;

  if (params) {
    const urlParams = new URLSearchParams(params);
    url += `?${urlParams.toString()}`;
  }

  const requestURL = "/api" + url;

  const headers = new Headers();
  headers.append("Content-Type", "application/json");
  headers.append("Accept", "application/json");

  const accessToken = document.cookie
    .split("; ")
    .find((row) => row.startsWith("accessToken="))
    ?.substring("accessToken=".length);

  if (accessToken) {
    headers.append("Authorization", `${accessToken}`);
  }

  const config = {
    method,
    headers,
    credentials: "include",
    ...(body && { body: JSON.stringify(body) }),
    ...(next ? { next } : { cache }),
  };
  return await fetch(requestURL, config);
}

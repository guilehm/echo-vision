export async function makeRequest(requester, url, options, next) {
  const response = await requester(url, options, next);
  if (response.ok) {
    if (response.status === 204) return {};
    return await response.json();
  }
  throw new Error(
    `Request failed with status ${response.status}: ${response.statusText}`,
  );
}

export async function getOwnEvents(requester, limit, cursor) {
  const params = {};
  if (limit) params.limit = limit;
  if (cursor) params.cursor = cursor;

  return await makeRequest(requester, "/events", {
    method: "GET",
    params,
  });
}

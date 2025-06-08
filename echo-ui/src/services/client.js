export async function makeRequest(requester, url, options, next) {
  const response = await requester(url, options, next);
  if (response.ok) {
    if (response.status === 204) return {};
    return await response.json();
  }
  const errorData = {
    status: response.status,
    errorMessage: response.statusText || "An error occurred",
  };
  try {
    const jsonResponse = await response.json();
    errorData.errorMessage = jsonResponse.error;
  } catch { }

  return errorData;
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

export async function signUp(requester, data = {}) {
  return await makeRequest(requester, "/users", {
    method: "POST",
    body: data,
  });
}

export async function signIn(requester, { email, password }) {
  return await makeRequest(requester, "/users/login", {
    method: "POST",
    body: { email, password },
  });
}

export async function uploadS3File({ file, presignedURL }) {
  const response = await fetch(presignedURL, {
    method: "PUT",
    headers: { "Content-Type": file.type },
    body: file,
  });
  if (!response.ok) {
    throw new Error("Failed to upload file");
  }
  return response;
}

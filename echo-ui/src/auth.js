"use server";

import { jwtDecode } from "jwt-decode";
import { cookies } from "next/headers";

export async function getSession() {
  const cookiesList = await cookies();
  const accessTokenCookie = cookiesList.get("accessToken");
  const accessToken = accessTokenCookie?.value;
  if (!accessToken) {
    console.log("no access token found in cookies");
    return null;
  }

  try {
    const decoded = jwtDecode(accessToken);
    return decoded;
  } catch (error) {
    console.error("could not decode session", error);
    return null;
  }
}

"use server";

import { jwtDecode } from "jwt-decode";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

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

export async function withAuth() {
  const session = await getSession();
  if (!session) {
    redirect("/sign-in");
  }
  return session;
}

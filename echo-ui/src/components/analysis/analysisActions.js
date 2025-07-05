"use server";

import { serverRequester } from "@/services/server-requester";
import { getOwnEvents } from "@/services/client";

export async function fetchMoreData({ limit, cursor }) {
  const response = await getOwnEvents(serverRequester, limit, cursor);
  if (response.status !== 200) {
    const { errorMessage } = response || {};
    throw new Error(errorMessage);
  }
  const newAnalyses = response.data?.results || [];
  const nextCursor = response.data?.nextCursor || null;
  return { newAnalyses, nextCursor };
}

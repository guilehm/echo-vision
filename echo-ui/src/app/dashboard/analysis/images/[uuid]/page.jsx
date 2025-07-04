import { notFound } from "next/navigation";
import { getEventById } from "@/services/client";
import { serverRequester } from "@/services/server-requester";

export default async function ImageDetailsPage({ params }) {
  const { uuid } = params;
  const response = await getEventById(serverRequester, uuid);
  if (!response || response.status !== 200) {
    notFound();
  }

  const data = response.data;
  if (!data) {
    notFound();
  }

  return (
    <div>
      <h1>Image Details for UUID: {uuid}</h1>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}

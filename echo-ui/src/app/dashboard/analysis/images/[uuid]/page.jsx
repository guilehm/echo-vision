import DashboardHeader from "@/components/headers/dashboard-header";
import ImageAnalysisDetail from "@/components/analysis/imageDetails";
import { getEventById } from "@/services/client";
import { serverRequester } from "@/services/server-requester";
import { notFound } from "next/navigation";
import { redirect } from "next/navigation";

const headerProps = {
  breadcrumbData: [
    { label: "Dashboard", href: "/dashboard/" },
    { label: "Analysis", href: null },
    { label: "Images", href: "/dashboard/analysis/images" },
    { label: "Details", href: null },
  ],
};

export default async function ImageAnalysisDetailPage({ params }) {
  const { uuid } = await params;
  const response = await getEventById(serverRequester, uuid);

  if (!response || response.status === 404) {
    notFound();
  }

  const data = response.data;
  if (!data) {
    notFound();
  }

  if (response.status === 403 || response.status === 401) {
    // TODO: handle unauthorized access
    redirect("/sign-in");
  }

  return (
    <>
      <DashboardHeader {...headerProps} />
      <ImageAnalysisDetail event={data} />
    </>
  );
}

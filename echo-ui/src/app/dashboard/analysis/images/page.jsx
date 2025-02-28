import DashboardHeader from "@/components/headers/dashboard-header";
import ImageAnalysis from "@/components/analysis/images";

const headerProps = {
  breadcrumbData: [
    { label: "Dashboard", href: "/dashboard/" },
    { label: "Analysis", href: "" },
    { label: "Images", href: "" },
  ],
};

export default function ImageAnalysisPage() {
  return (
    <>
      <DashboardHeader {...headerProps} />
      <ImageAnalysis />
    </>
  );
}

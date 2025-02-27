import DashboardHeader from "@/components/headers/dashboard-header";
import { Button } from "@/components/ui/button";

const breadcrumbData = [
  { label: "Dashboard", href: "/dashboard/" },
  { label: "Analysis", href: "" },
  { label: "Images", href: "" },
];

export default function AnalysisPage() {
  return (
    <>
      <DashboardHeader breadcrumbData={breadcrumbData} />
      <div className="flex flex-col gap-4 p-4 pt-0">
        <div className="flex-1 space-y-4 pt-6">
          <div className="flex items-center justify-between space-y-2">
            <h2 className="text-3xl font-bold tracking-tight">
              Image Analysis
            </h2>
            {/* <div className="flex items-center space-x-2"> */}
            {/*   <Button>New</Button> */}
            {/* </div> */}
          </div>
        </div>

        <div className="grid auto-rows-min gap-4 md:grid-cols-3">
          <div className="bg-muted/50 aspect-video rounded-xl" />
          <div className="bg-muted/50 aspect-video rounded-xl" />
          <div className="bg-muted/50 aspect-video rounded-xl" />
        </div>
      </div>
    </>
  );
}

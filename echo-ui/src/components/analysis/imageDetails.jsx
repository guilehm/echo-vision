import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ImageUpload from "@/components/uploads/image-upload";
import { AnalysisListTable } from "./analysisTable";

export default async function ImageAnalysisDetail({ event }) {
  return (
    <div className="flex flex-col gap-4 p-4 pt-0">
      <div className="flex-1 space-y-4 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">
            Image Analysis Detail
          </h2>
        </div>
      </div>

      <div className="h-full flex-col md:flex">
        <div className="container flex flex-col items-start justify-between space-y-2 py-4 sm:flex-row sm:items-center sm:space-y-0 md:h-16">
          <h2 className="text-lg font-semibold">Some Info here</h2>
        </div>
        <Separator />
      </div>
    </div>
  );
}

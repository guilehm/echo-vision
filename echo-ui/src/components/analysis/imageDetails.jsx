import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ImageUpload from "@/components/uploads/image-upload";
import { format } from "date-fns";
import {
  CalendarDays,
  Clock,
  Hash,
  Layers,
  RefreshCw,
  ShieldAlert,
  ShieldCheck,
  Type,
} from "lucide-react";
import { AnalysisListTable } from "./analysisTable";

export default async function ImageAnalysisDetail({ event }) {
  // Helper function to parse and display the result data if it exists
  const renderResultData = () => {
    if (!event.result || event.result.length === 0) return null;

    return (
      <div className="mt-6">
        <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
          <Layers className="h-5 w-5" />
          Analysis Results
        </h3>
        <div className="bg-muted rounded-lg p-4 text-sm overflow-x-auto">
          {JSON.stringify(event.result, null, 2)}
        </div>
      </div>
    );
  };

  // Status badge with appropriate colors
  const getStatusBadge = () => {
    const baseClasses =
      "inline-flex items-center rounded-full px-3 py-1 text-xs font-medium";

    switch (event.status.toLowerCase()) {
      case "completed":
        return (
          <span className={`${baseClasses} bg-green-100 text-green-800`}>
            <ShieldCheck className="h-3 w-3 mr-1" />
            Completed
          </span>
        );
      case "failed":
        return (
          <span className={`${baseClasses} bg-red-100 text-red-800`}>
            <ShieldAlert className="h-3 w-3 mr-1" />
            Failed
          </span>
        );
      case "processing":
        return (
          <span className={`${baseClasses} bg-blue-100 text-blue-800`}>
            <RefreshCw className="h-3 w-3 mr-1 animate-spin" />
            Processing
          </span>
        );
      default:
        return (
          <span className={`${baseClasses} bg-gray-100 text-gray-800`}>
            {event.Status}
          </span>
        );
    }
  };

  return (
    <div className="flex flex-col gap-4 p-4 pt-0">
      <div className="flex-1 space-y-4 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">
            Image Analysis Detail
          </h2>
          {getStatusBadge()}
        </div>
      </div>

      <div className="h-full flex-col md:flex">
        <div className="container flex flex-col items-start justify-between space-y-2 py-4 sm:flex-row sm:items-center sm:space-y-0 md:h-16">
          <h2 className="text-lg font-semibold">Event Information</h2>
        </div>
        <Separator />

        <div className="grid gap-6 py-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Event ID */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex flex-col space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Hash className="h-4 w-4" />
                  Event ID
                </div>
                <p className="font-medium overflow-hidden text-ellipsis">
                  {event.id}
                </p>
              </div>
            </div>

            {/* User ID */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex flex-col space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Hash className="h-4 w-4" />
                  User ID
                </div>
                <p className="font-medium overflow-hidden text-ellipsis">
                  {event.userID}
                </p>
              </div>
            </div>

            {/* Event Type */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex flex-col space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Type className="h-4 w-4" />
                  Event Type
                </div>
                <p className="font-medium">{event.eventType}</p>
              </div>
            </div>

            {/* Sub Type */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex flex-col space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Type className="h-4 w-4" />
                  Sub Type
                </div>
                <p className="font-medium">{event.subType}</p>
              </div>
            </div>

            {/* Created At */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex flex-col space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <CalendarDays className="h-4 w-4" />
                  Created At
                </div>
                <p className="font-medium">
                  {format(new Date(event.createdAt), "PPP")}
                  <span className="flex items-center gap-1 text-sm text-muted-foreground mt-1">
                    <Clock className="h-3 w-3" />
                    {format(new Date(event.createdAt), "pp")}
                  </span>
                </p>
              </div>
            </div>

            {/* Updated At */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex flex-col space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <RefreshCw className="h-4 w-4" />
                  Updated At
                </div>
                <p className="font-medium">
                  {format(new Date(event.updatedAt), "PPP")}
                  <span className="flex items-center gap-1 text-sm text-muted-foreground mt-1">
                    <Clock className="h-3 w-3" />
                    {format(new Date(event.updatedAt), "pp")}
                  </span>
                </p>
              </div>
            </div>
          </div>

          {/* Result Data */}
          {renderResultData()}
        </div>
      </div>
    </div>
  );
}

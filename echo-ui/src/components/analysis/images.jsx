import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ImageUpload from "@/components/uploads/image-upload";
import { getOwnEvents } from "@/services/server-requester";
import { formatDate, statusStyles } from "@/utils";

async function fetchUserAnalyses() {
  const response = await getOwnEvents();
  return response.data?.results;
}

export default async function ImageAnalysis() {
  const results = await fetchUserAnalyses();
  return (
    <div className="flex flex-col gap-4 p-4 pt-0">
      <div className="flex-1 space-y-4 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">Image Analysis</h2>
        </div>
      </div>

      <div className="h-full flex-col md:flex">
        <div className="container flex flex-col items-start justify-between space-y-2 py-4 sm:flex-row sm:items-center sm:space-y-0 md:h-16">
          <h2 className="text-lg font-semibold">New Analysis</h2>
        </div>
        <Separator />

        <Tabs defaultValue="edit" className="flex-1">
          <div className="container h-full py-6">
            <div className="grid h-full items-stretch gap-6 md:grid-cols-[1fr_50px]">
              <div className="hidden flex-col space-y-4 sm:flex md:order-2">
                <div className="grid gap-2">
                  <TabsList className="grid grid-cols-1">
                    <TabsTrigger value="edit">
                      <span className="sr-only">Edit</span>
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 20 20"
                        fill="none"
                        className="h-5 w-5"
                      >
                        <rect
                          x="4"
                          y="3"
                          width="12"
                          height="2"
                          rx="1"
                          fill="currentColor"
                        />
                        <rect
                          x="4"
                          y="7"
                          width="12"
                          height="2"
                          rx="1"
                          fill="currentColor"
                        />
                        <rect
                          x="4"
                          y="11"
                          width="3"
                          height="2"
                          rx="1"
                          fill="currentColor"
                        />
                        <rect
                          x="4"
                          y="15"
                          width="4"
                          height="2"
                          rx="1"
                          fill="currentColor"
                        />
                        <rect
                          x="8.5"
                          y="11"
                          width="3"
                          height="2"
                          rx="1"
                          fill="currentColor"
                        />
                        <path
                          d="M17.154 11.346a1.182 1.182 0 0 0-1.671 0L11 15.829V17.5h1.671l4.483-4.483a1.182 1.182 0 0 0 0-1.671Z"
                          fill="currentColor"
                        />
                      </svg>
                    </TabsTrigger>
                  </TabsList>
                </div>
              </div>

              <div className="md:order-1 space-y-6">
                <TabsContent value="edit" className="mt-0 border-0 p-0">
                  <ImageUpload eventSubType={"image_analysis"} />
                </TabsContent>
                <AnalysisListTable analyses={results} />
              </div>
            </div>
          </div>
        </Tabs>
      </div>
    </div>
  );
}

export function AnalysisListTable({ analyses }) {
  if (!analyses || analyses.length === 0) {
    return;
  }
  return (
    <div className="mt-8">
      <Card>
        <CardHeader>
          <CardTitle>Previous Analyses</CardTitle>
        </CardHeader>
        <CardContent>
          {analyses.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-8">
              <p className="text-muted-foreground">No analyses found</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Analysis ID</TableHead>
                  {/* <TableHead>Type</TableHead> */}
                  <TableHead>Type</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {analyses.map((analysis) => (
                  <TableRow key={analysis.id}>
                    <TableCell className="font-medium">
                      {analysis.id.substring(0, 8)}...
                    </TableCell>
                    {/* <TableCell className="capitalize"> */}
                    {/*   {analysis.eventType.replace("_", " ")} */}
                    {/* </TableCell> */}
                    <TableCell className="capitalize">
                      {analysis.subType.replace("_", " ")}
                    </TableCell>
                    <TableCell>
                      <Badge
                        style={{
                          backgroundColor:
                            statusStyles[analysis.status]?.bg ||
                            statusStyles.default.bg,
                          color:
                            statusStyles[analysis.status]?.text ||
                            statusStyles.default.text,
                          borderColor:
                            statusStyles[analysis.status]?.border ||
                            statusStyles.default.border,
                        }}
                        className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border"
                      >
                        {analysis.status}
                      </Badge>
                    </TableCell>
                    <TableCell>{formatDate(analysis.createdAt)}</TableCell>
                    <TableCell>
                      <Button variant="outline" size="sm">
                        View
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

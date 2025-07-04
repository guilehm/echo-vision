"use client";
import * as badge from "@/components/ui/badge";
import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { clientRequester } from "@/services/client-requester";
import { getOwnEvents } from "@/services/client";
import { formatDate, statusStyles } from "@/utils";
import { useState } from "react";
import InfiniteScroll from "react-infinite-scroll-component";
import { toast } from "sonner";

export function AnalysisListTable({ initialData, initialCursor }) {
  const [analyses, setAnalyses] = useState(initialData);
  const [cursor, setCursor] = useState(initialCursor);
  const [hasMore, setHasMore] = useState(!!initialCursor);
  const limit = 10;

  if (!initialData || initialData.length === 0) {
    return;
  }

  async function fetchMoreData() {
    const response = await getOwnEvents(clientRequester, limit, cursor);
    if (response.status !== 200) {
      toast.error("Failed to load more analyses");
      return;
    }
    const newAnalyses = response.data?.results || [];
    setAnalyses((prev) => [...prev, ...newAnalyses]);
    setCursor(response.data.nextCursor || null);
    setHasMore(!!response.data.nextCursor);
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
            <InfiniteScroll
              dataLength={analyses.length}
              next={fetchMoreData}
              hasMore={hasMore}
              loader={
                <p className="mt-2 text-center text-muted-foreground">
                  Loading more analyses...
                </p>
              }
              endMessage={
                initialData.length !== analyses.length && (
                  <p className="mt-2 text-center text-muted-foreground">
                    {"Yay! You've seen it all!"}
                  </p>
                )
              }
            >
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Analysis ID</TableHead>
                    {/* <TableHead>Type</TableHead> */}
                    <TableHead>Type</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Created At</TableHead>
                    {/* <TableHead>Actions</TableHead> */}
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {analyses.map((analysis) => (
                    <TableRow key={analysis.id}>
                      <TableCell className="font-medium">
                        <Link href={`images/${analysis.id}`}>
                          {analysis.id.substring(0, 8)}...
                        </Link>
                      </TableCell>
                      {/* <TableCell className="capitalize"> */}
                      {/*   {analysis.eventType.replace("_", " ")} */}
                      {/* </TableCell> */}
                      <TableCell className="capitalize">
                        {analysis.subType.replace("_", " ")}
                      </TableCell>
                      <TableCell>
                        <badge.Badge
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
                        </badge.Badge>
                      </TableCell>
                      <TableCell>{formatDate(analysis.createdAt)}</TableCell>
                      {/* <TableCell> */}
                      {/*   <Button variant="outline" size="sm"> */}
                      {/*     View */}
                      {/*   </Button> */}
                      {/* </TableCell> */}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </InfiniteScroll>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

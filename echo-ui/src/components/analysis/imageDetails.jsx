"use client";

import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import {
  Copy,
  Download,
  Eye,
  Calendar,
  Clock,
  User,
  Hash,
  Tag,
  Activity,
  ImageIcon,
} from "lucide-react";
import { useState } from "react";
import { imagePlaceholder } from "@/utils";

export default function ImageAnalysisDetail({ event }) {
  const [copiedField, setCopiedField] = useState(null);

  const isImageFile = (contentType) => {
    return contentType.startsWith("image/");
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (
      Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
    );
  };

  const copyToClipboard = async (text, field) => {
    await navigator.clipboard.writeText(text);
    setCopiedField(field);
    setTimeout(() => setCopiedField(null), 2000);
  };

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case "completed":
      case "success":
        return "bg-green-100 text-green-800 border-green-200";
      case "pending":
      case "processing":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "failed":
      case "error":
        return "bg-red-100 text-red-800 border-red-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  };

  const formatResult = (result) => {
    if (!result) return null;
    return JSON.stringify(result, null, 2);
  };

  return (
    <div className="flex flex-col gap-4 p-4 pt-0">
      <div className="flex-1 space-y-4 pt-6">
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between space-y-2 sm:space-y-0">
          <h2 className="text-3xl font-bold tracking-tight">
            Image Analysis Detail
          </h2>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" size="sm">
              <Download className="w-4 h-4 mr-2" />
              Export
            </Button>
            <Button variant="outline" size="sm">
              <Eye className="w-4 h-4 mr-2" />
              View Raw
            </Button>
          </div>
        </div>
      </div>

      <div className="h-full flex-col md:flex">
        <div className="container flex flex-col items-start justify-between space-y-2 py-4 sm:flex-row sm:items-center sm:space-y-0 md:h-16">
          <div className="flex items-center gap-2">
            <h2 className="text-lg font-semibold">Event Information</h2>
            <Badge className={getStatusColor(event.status)}>
              {event.status}
            </Badge>
          </div>
        </div>
        <Separator />

        {/* Status Summary */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="w-4 h-4" />
              Processing Summary
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3">
              <div className="text-center p-4 bg-muted/50 rounded-lg">
                <div className="text-2xl font-bold">{event.eventType}</div>
                <div className="text-sm text-muted-foreground">Event Type</div>
              </div>
              <div className="text-center p-4 bg-muted/50 rounded-lg">
                <div
                  className={`text-2xl font-bold ${
                    event.status.toLowerCase() === "completed"
                      ? "text-green-600"
                      : event.status.toLowerCase() === "pending"
                        ? "text-yellow-600"
                        : event.status.toLowerCase() === "failed"
                          ? "text-red-600"
                          : "text-gray-600"
                  }`}
                >
                  {event.status}
                </div>
                <div className="text-sm text-muted-foreground">
                  Current Status
                </div>
              </div>
              <div className="text-center p-4 bg-muted/50 rounded-lg">
                <div className="text-2xl font-bold">
                  {Math.round(
                    (new Date(event.updatedAt).getTime() -
                      new Date(event.createdAt).getTime()) /
                      1000,
                  )}
                  s
                </div>
                <div className="text-sm text-muted-foreground">
                  Processing Time
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 py-6">
          {/* Event Metadata */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base">
                <Hash className="w-4 h-4" />
                Event Metadata
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-muted-foreground">
                    Event ID
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => copyToClipboard(event.id, "id")}
                    className="h-6 px-2"
                  >
                    <Copy className="w-3 h-3" />
                  </Button>
                </div>
                <p className="font-mono bg-muted p-2 rounded text-xs break-all">
                  {event.id}
                </p>
                {copiedField === "id" && (
                  <p className="text-xs text-green-600">Copied!</p>
                )}
              </div>

              <div className="space-y-2">
                <span className="text-sm font-medium text-muted-foreground">
                  Event Type
                </span>
                <div className="flex gap-2">
                  <Badge variant="secondary">{event.eventType}</Badge>
                  {event.subType && (
                    <Badge variant="outline">{event.subType}</Badge>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* User Information */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base">
                <User className="w-4 h-4" />
                User Information
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-muted-foreground">
                    User ID
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => copyToClipboard(event.userID, "userID")}
                    className="h-6 px-2"
                  >
                    <Copy className="w-3 h-3" />
                  </Button>
                </div>
                <p className="text-sm font-mono bg-muted p-2 rounded text-xs break-all">
                  {event.userID}
                </p>
                {copiedField === "userID" && (
                  <p className="text-xs text-green-600">Copied!</p>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Timestamps */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base">
                <Clock className="w-4 h-4" />
                Timeline
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Calendar className="w-3 h-3 text-muted-foreground" />
                  <span className="text-sm font-medium text-muted-foreground">
                    Created
                  </span>
                </div>
                <p className="text-sm">{formatDate(event.createdAt)}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Activity className="w-3 h-3 text-muted-foreground" />
                  <span className="text-sm font-medium text-muted-foreground">
                    Updated
                  </span>
                </div>
                <p className="text-sm">{formatDate(event.updatedAt)}</p>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="py-6 space-y-6">
          {/* Image Display Section */}
          {event.file && isImageFile(event.file.contentType) && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <ImageIcon className="w-4 h-4" />
                  Image Preview
                </CardTitle>
                <CardDescription>
                  Original image that was analyzed
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col lg:flex-row gap-6">
                  <div className="flex-1">
                    <div className="relative bg-muted rounded-lg overflow-hidden">
                      <img
                        src={event.file.url || imagePlaceholder}
                        alt={event.file.filename}
                        className="w-full h-auto max-h-96 object-contain"
                        onError={(e) => {
                          e.currentTarget.src =
                            "/placeholder.svg?height=400&width=600";
                        }}
                      />
                    </div>
                  </div>
                  <div className="lg:w-80 space-y-4">
                    <div className="space-y-2">
                      <span className="text-sm font-medium text-muted-foreground">
                        Filename
                      </span>
                      <p className="text-sm font-mono bg-muted p-2 rounded">
                        {event.file.filename}
                      </p>
                    </div>
                    <div className="space-y-2">
                      <span className="text-sm font-medium text-muted-foreground">
                        Content Type
                      </span>
                      <Badge variant="secondary">
                        {event.file.contentType}
                      </Badge>
                    </div>
                    <div className="space-y-2">
                      <span className="text-sm font-medium text-muted-foreground">
                        File Size
                      </span>
                      <p className="text-sm">
                        {formatFileSize(event.file.filesize)}
                      </p>
                    </div>
                    <Button className="w-full bg-transparent" variant="outline">
                      <Download className="w-4 h-4 mr-2" />
                      Download Original
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* File Info for non-images */}
          {event.file && !isImageFile(event.file.contentType) && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="w-4 h-4" />
                  File Information
                </CardTitle>
                <CardDescription>
                  Details about the processed file
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <span className="text-sm font-medium text-muted-foreground">
                      Filename
                    </span>
                    <p className="text-sm font-mono bg-muted p-2 rounded">
                      {event.file.filename}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <span className="text-sm font-medium text-muted-foreground">
                      Content Type
                    </span>
                    <Badge variant="secondary">{event.file.contentType}</Badge>
                  </div>
                  <div className="space-y-2">
                    <span className="text-sm font-medium text-muted-foreground">
                      File Size
                    </span>
                    <p className="text-sm">
                      {formatFileSize(event.file.filesize)}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <span className="text-sm font-medium text-muted-foreground">
                      File Path
                    </span>
                    <p className="font-mono bg-muted p-2 rounded text-xs break-all">
                      {event.file.filepath}
                    </p>
                  </div>
                </div>
                <Button className="mt-4 bg-transparent" variant="outline">
                  <Download className="w-4 h-4 mr-2" />
                  Download File
                </Button>
              </CardContent>
            </Card>
          )}

          {/* Analysis Results */}
          {event.result && (
            <Card className="mt-6">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Tag className="w-4 h-4" />
                  Analysis Results
                </CardTitle>
                <CardDescription>
                  Detailed results from the image analysis process
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="relative">
                  <Button
                    variant="outline"
                    size="sm"
                    className="absolute top-2 right-2 z-10 bg-transparent"
                    onClick={() =>
                      copyToClipboard(
                        formatResult(event.result) || "",
                        "result",
                      )
                    }
                  >
                    <Copy className="w-3 h-3 mr-1" />
                    Copy
                  </Button>
                  <pre className="bg-muted p-4 rounded-lg text-xs overflow-auto max-h-96 font-mono">
                    {formatResult(event.result)}
                  </pre>
                  {copiedField === "result" && (
                    <p className="text-xs text-green-600 mt-2">
                      Results copied to clipboard!
                    </p>
                  )}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}

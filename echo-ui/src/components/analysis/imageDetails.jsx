"use client";

import { Badge } from "@/components/ui/badge";

import { formatDate } from "@/utils";
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
  FileText,
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

// TODO: use types
const isDetectFaces = (event) => event.subType === "detect_faces";

const getEmotionColor = (emotionType) => {
  const colors = {
    HAPPY: "text-yellow-600 bg-yellow-50 border-yellow-200",
    SAD: "text-blue-600 bg-blue-50 border-blue-200",
    ANGRY: "text-red-600 bg-red-50 border-red-200",
    FEAR: "text-purple-600 bg-purple-50 border-purple-200",
    SURPRISED: "text-orange-600 bg-orange-50 border-orange-200",
    DISGUSTED: "text-green-600 bg-green-50 border-green-200",
    CALM: "text-teal-600 bg-teal-50 border-teal-200",
    CONFUSED: "text-gray-600 bg-gray-50 border-gray-200",
  };
  return colors[emotionType] || "text-gray-600 bg-gray-50 border-gray-200";
};

const parseEmotionResults = (result) => {
  if (!result || !Array.isArray(result)) return [];
  return result.map((detection, index) => ({
    id: index,
    emotions: detection.emotions || [],
    confidence: detection.confidence || 0,
    boundingBox: detection.boundingBox || null,
  }));
};

const getTopEmotion = (emotions) => {
  if (!emotions || emotions.length === 0) return null;
  return emotions.reduce((prev, current) =>
    prev.Confidence > current.Confidence ? prev : current,
  );
};

export default function ImageAnalysisDetail({ event }) {
  const [copiedField, setCopiedField] = useState(null);
  const [selectedDetection, setSelectedDetection] = useState(null);

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
                  className={`text-2xl font-bold ${event.status.toLowerCase() === "completed"
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
                <p className="font-mono bg-muted p-2 rounded text-xs break-all">
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
          {/* Enhanced Image Display Section with Bounding Boxes */}
          {event.file && isImageFile(event.file.contentType) && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <ImageIcon className="w-4 h-4" />
                  Image Analysis Preview
                </CardTitle>
                <CardDescription>
                  Original image with detected info
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col xl:flex-row gap-6">
                  <div className="flex-1">
                    <div className="relative bg-muted rounded-lg overflow-hidden">
                      <div className="relative inline-block">
                        <img
                          id="analysis-image"
                          src={event.file.url || imagePlaceholder}
                          alt={event.file.filename}
                          className="w-full h-auto max-h-96 object-contain"
                          onError={(e) => {
                            e.currentTarget.src = imagePlaceholder;
                          }}
                          onLoad={(e) => {
                            const img = e.currentTarget;
                            const container = img.parentElement;
                            if (!container) return;

                            const existingOverlays = container.querySelectorAll(
                              ".bounding-box-overlay",
                            );
                            existingOverlays.forEach((overlay) =>
                              overlay.remove(),
                            );

                            const emotionResults = parseEmotionResults(
                              event.result,
                            );
                            emotionResults.forEach((detection, index) => {
                              if (detection.boundingBox) {
                                const overlay = document.createElement("div");
                                overlay.className =
                                  "bounding-box-overlay absolute border-2 border-red-500 bg-red-500/10 pointer-events-auto cursor-pointer transition-all";
                                overlay.style.left = `${detection.boundingBox.left * 100}%`;
                                overlay.style.top = `${detection.boundingBox.top * 100}%`;
                                overlay.style.width = `${detection.boundingBox.width * 100}%`;
                                overlay.style.height = `${detection.boundingBox.height * 100}%`;

                                overlay.onclick = () => {
                                  setSelectedDetection(index);
                                };

                                const topEmotion = getTopEmotion(
                                  detection.emotions,
                                );
                                if (topEmotion) {
                                  const label = document.createElement("div");
                                  label.className =
                                    "absolute -top-6 left-0 bg-red-500 text-white text-xs px-2 py-1 rounded whitespace-nowrap pointer-events-none";
                                  label.textContent = `${topEmotion.Type} (${topEmotion.Confidence.toFixed(1)}%)`;
                                  overlay.appendChild(label);
                                }

                                container.appendChild(overlay);
                              }
                            });
                          }}
                        />
                      </div>
                    </div>
                    {selectedDetection !== null && (
                      <div className="mt-4">
                        {(() => {
                          const detection = parseEmotionResults(event.result)[
                            selectedDetection
                          ];
                          const topEmotion = getTopEmotion(detection.emotions);
                          const topEmotions = detection.emotions
                            .sort((a, b) => b.Confidence - a.Confidence)
                            .slice(0, 3);

                          return (
                            <Card className="w-full max-w-sm">
                              <CardHeader className="pb-2">
                                <div className="flex items-center justify-between">
                                  <CardTitle className="text-sm">
                                    Face #{selectedDetection + 1}
                                  </CardTitle>
                                  <Badge
                                    variant="outline"
                                    className="text-xs bg-green-50 text-green-700 border-green-200"
                                  >
                                    {detection.confidence.toFixed(1)}%
                                  </Badge>
                                </div>
                              </CardHeader>
                              <CardContent className="pt-0 space-y-3">
                                {topEmotion && (
                                  <div className="text-center p-2 rounded-lg bg-muted/50">
                                    <Badge
                                      className={`${getEmotionColor(topEmotion.Type)} mb-1`}
                                    >
                                      {topEmotion.Type}
                                    </Badge>
                                    <div className="text-lg font-bold">
                                      {topEmotion.Confidence.toFixed(1)}%
                                    </div>
                                  </div>
                                )}

                                <div className="space-y-1">
                                  {topEmotions.map((emotion, emotionIndex) => (
                                    <div
                                      key={emotionIndex}
                                      className="flex items-center justify-between text-xs"
                                    >
                                      <span
                                        className={`px-2 py-1 rounded ${getEmotionColor(emotion.Type)}`}
                                      >
                                        {emotion.Type}
                                      </span>
                                      <div className="flex items-center gap-1 flex-1 ml-2">
                                        <div className="flex-1 bg-muted rounded-full h-1.5">
                                          <div
                                            className="bg-primary h-1.5 rounded-full"
                                            style={{
                                              width: `${Math.min(emotion.Confidence, 100)}%`,
                                            }}
                                          />
                                        </div>
                                        <span className="text-xs text-muted-foreground min-w-[2.5rem] text-right">
                                          {emotion.Confidence.toFixed(1)}%
                                        </span>
                                      </div>
                                    </div>
                                  ))}
                                </div>

                                {detection.boundingBox && (
                                  <div className="text-xs text-muted-foreground bg-muted/30 p-2 rounded">
                                    <div className="grid grid-cols-2 gap-1">
                                      <div>
                                        X:{" "}
                                        {(
                                          detection.boundingBox.left * 100
                                        ).toFixed(0)}
                                        %
                                      </div>
                                      <div>
                                        Y:{" "}
                                        {(
                                          detection.boundingBox.top * 100
                                        ).toFixed(0)}
                                        %
                                      </div>
                                      <div>
                                        W:{" "}
                                        {(
                                          detection.boundingBox.width * 100
                                        ).toFixed(0)}
                                        %
                                      </div>
                                      <div>
                                        H:{" "}
                                        {(
                                          detection.boundingBox.height * 100
                                        ).toFixed(0)}
                                        %
                                      </div>
                                    </div>
                                  </div>
                                )}
                              </CardContent>
                            </Card>
                          );
                        })()}
                      </div>
                    )}
                  </div>
                  <div className="xl:w-80 space-y-4">
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

          {/* Emotion Analysis Results */}
          {isDetectFaces(event) &&
            event.result &&
            parseEmotionResults(event.result).length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Tag className="w-4 h-4" />
                    Emotion Analysis Results {isDetectFaces}
                </CardTitle>
                <CardDescription>
                  {parseEmotionResults(event.result).length} face
                  {parseEmotionResults(event.result).length > 1
                    ? "s"
                    : ""}{" "}
                    detected
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {parseEmotionResults(event.result).map(
                    (detection, index) => {
                      const topEmotion = getTopEmotion(detection.emotions);
                      const topEmotions = detection.emotions
                        .sort((a, b) => b.Confidence - a.Confidence)
                        .slice(0, 3); // Show only top 3 emotions

                      return (
                        <Card key={index} className="relative">
                          <CardHeader className="pb-2">
                            <div className="flex items-center justify-between">
                              <CardTitle className="text-sm">
                                  Face #{index + 1}
                              </CardTitle>
                              <Badge
                                variant="outline"
                                className="text-xs bg-green-50 text-green-700 border-green-200"
                              >
                                {detection.confidence.toFixed(1)}%
                              </Badge>
                            </div>
                          </CardHeader>
                          <CardContent className="pt-0 space-y-3">
                            {/* Primary Emotion */}
                            {topEmotion && (
                              <div className="text-center p-2 rounded-lg bg-muted/50">
                                <Badge
                                  className={`${getEmotionColor(topEmotion.Type)} mb-1`}
                                >
                                  {topEmotion.Type}
                                </Badge>
                                <div className="text-lg font-bold">
                                  {topEmotion.Confidence.toFixed(1)}%
                                </div>
                              </div>
                            )}

                            {/* Top 3 Emotions */}
                            <div className="space-y-1">
                              {topEmotions.map((emotion, emotionIndex) => (
                                <div
                                  key={emotionIndex}
                                  className="flex items-center justify-between text-xs"
                                >
                                  <span
                                    className={`px-2 py-1 rounded text-xs ${getEmotionColor(emotion.Type)}`}
                                  >
                                    {emotion.Type}
                                  </span>
                                  <div className="flex items-center gap-1 flex-1 ml-2">
                                    <div className="flex-1 bg-muted rounded-full h-1.5">
                                      <div
                                        className="bg-primary h-1.5 rounded-full transition-all"
                                        style={{
                                          width: `${Math.min(emotion.Confidence, 100)}%`,
                                        }}
                                      />
                                    </div>
                                    <span className="text-xs text-muted-foreground min-w-[2.5rem] text-right">
                                      {emotion.Confidence.toFixed(1)}%
                                    </span>
                                  </div>
                                </div>
                              ))}
                            </div>

                            {/* Bounding Box Info */}
                            {detection.boundingBox && (
                              <div className="text-xs text-muted-foreground bg-muted/30 p-2 rounded">
                                <div className="grid grid-cols-2 gap-1">
                                  <div>
                                      X:{" "}
                                    {(
                                      detection.boundingBox.left * 100
                                    ).toFixed(0)}
                                      %
                                  </div>
                                  <div>
                                      Y:{" "}
                                    {(
                                      detection.boundingBox.top * 100
                                    ).toFixed(0)}
                                      %
                                  </div>
                                  <div>
                                      W:{" "}
                                    {(
                                      detection.boundingBox.width * 100
                                    ).toFixed(0)}
                                      %
                                  </div>
                                  <div>
                                      H:{" "}
                                    {(
                                      detection.boundingBox.height * 100
                                    ).toFixed(0)}
                                      %
                                  </div>
                                </div>
                              </div>
                            )}

                            {/* Show all emotions button */}
                            {/* {detection.emotions.length > 3 && ( */}
                            {/*   <Button */}
                            {/*     variant="ghost" */}
                            {/*     size="sm" */}
                            {/*     className="w-full h-6 text-xs" */}
                            {/*     onClick={() => { */}
                            {/*       // Toggle detailed view - you can implement this functionality */}
                            {/*       console.log( */}
                            {/*         "Show all emotions for detection", */}
                            {/*         index, */}
                            {/*       ); */}
                            {/*     }} */}
                            {/*   > */}
                            {/*     +{detection.emotions.length - 3} more emotions */}
                            {/*   </Button> */}
                            {/* )} */}
                          </CardContent>
                        </Card>
                      );
                    },
                  )}
                </div>

                {/* Summary Stats */}
                {parseEmotionResults(event.result).length > 1 && (
                  <div className="mt-4 p-3 bg-muted/30 rounded-lg">
                    <div className="text-sm font-medium mb-2">
                        Overall Summary
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-xs">
                      <div className="text-center">
                        <div className="font-semibold">
                          {parseEmotionResults(event.result).length}
                        </div>
                        <div className="text-muted-foreground">Faces</div>
                      </div>
                      <div className="text-center">
                        <div className="font-semibold">
                          {(
                            parseEmotionResults(event.result).reduce(
                              (sum, d) => sum + d.confidence,
                              0,
                            ) / parseEmotionResults(event.result).length
                          ).toFixed(1)}
                            %
                        </div>
                        <div className="text-muted-foreground">
                            Avg Confidence
                        </div>
                      </div>
                      <div className="text-center">
                        <div className="font-semibold">
                          {
                            parseEmotionResults(event.result)
                              .map((d) => getTopEmotion(d.emotions)?.Type)
                              .filter(
                                (emotion, index, arr) =>
                                  arr.indexOf(emotion) === index,
                              ).length
                          }
                        </div>
                        <div className="text-muted-foreground">
                            Unique Emotions
                        </div>
                      </div>
                      <div className="text-center">
                        <div className="font-semibold">
                          {
                            parseEmotionResults(event.result)
                              .map((d) => getTopEmotion(d.emotions)?.Type)
                              .reduce((acc, emotion) => {
                                acc[emotion] = (acc[emotion] || 0) + 1;
                                return acc;
                              }, {}).mostCommonEmotion
                          }
                        </div>
                        <div className="text-muted-foreground">
                            Most Common
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      {selectedDetection !== null && (
        <div
          onClick={() => setSelectedDetection(null)}
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-sm"
        >
          <div className="bg-white rounded-lg shadow-lg max-w-sm w-full p-4 relative">
            {/* <button */}
            {/*   className="absolute top-2 right-2 text-muted-foreground hover:text-black" */}
            {/*   onClick={() => setSelectedDetection(null)} */}
            {/* > */}
            {/*   ✕ */}
            {/* </button> */}

            {(() => {
              const detection = parseEmotionResults(event.result)[
                selectedDetection
              ];
              const topEmotion = getTopEmotion(detection.emotions);
              const topEmotions = detection.emotions
                .sort((a, b) => b.Confidence - a.Confidence)
                .slice(0, 3);

              return (
                <>
                  <h3 className="text-sm font-semibold mb-2">
                    Face #{selectedDetection + 1}
                  </h3>

                  <div className="space-y-3">
                    {topEmotion && (
                      <div className="text-center p-2 rounded-lg bg-muted/50">
                        <Badge
                          className={`${getEmotionColor(topEmotion.Type)} mb-1`}
                        >
                          {topEmotion.Type}
                        </Badge>
                        <div className="text-lg font-bold">
                          {topEmotion.Confidence.toFixed(1)}%
                        </div>
                      </div>
                    )}

                    <div className="space-y-1">
                      {topEmotions.map((emotion, idx) => (
                        <div
                          key={idx}
                          className="flex items-center justify-between text-xs"
                        >
                          <span
                            className={`px-2 py-1 rounded ${getEmotionColor(emotion.Type)}`}
                          >
                            {emotion.Type}
                          </span>
                          <div className="flex items-center gap-1 flex-1 ml-2">
                            <div className="flex-1 bg-muted rounded-full h-1.5">
                              <div
                                className="bg-primary h-1.5 rounded-full"
                                style={{
                                  width: `${Math.min(emotion.Confidence, 100)}%`,
                                }}
                              />
                            </div>
                            <span className="text-xs text-muted-foreground min-w-[2.5rem] text-right">
                              {emotion.Confidence.toFixed(1)}%
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>

                    {detection.boundingBox && (
                      <div className="text-xs text-muted-foreground bg-muted/30 p-2 rounded">
                        <div className="grid grid-cols-2 gap-1">
                          <div>
                            X: {(detection.boundingBox.left * 100).toFixed(0)}%
                          </div>
                          <div>
                            Y: {(detection.boundingBox.top * 100).toFixed(0)}%
                          </div>
                          <div>
                            W: {(detection.boundingBox.width * 100).toFixed(0)}%
                          </div>
                          <div>
                            H: {(detection.boundingBox.height * 100).toFixed(0)}
                            %
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                </>
              );
            })()}
          </div>
        </div>
      )}
    </div>
  );
}

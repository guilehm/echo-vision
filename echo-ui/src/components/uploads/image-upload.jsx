"use client";

import { useState } from "react";
import { useDropzone } from "react-dropzone";
import Image from "next/image";
import { toast } from "sonner";
import { createEvent, getPresignedUrl } from "@/services/server-requester";
import { Button } from "@/components/ui/button";
import { uploadS3File } from "@/services/client-requester";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const handleErrors = (error) => {
  console.log(error);
  toast.error("An error occurred. Please try again later.");
};

export default function ImageUpload() {
  const [file, setFile] = useState(null);
  const [eventSubType, setEventSubType] = useState("detect_labels");

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: { "image/*": [] },
    multiple: false,
    onDrop: (acceptedFiles) => {
      setFile(acceptedFiles[0]);
    },
  });

  const handleTypeChange = (value) => {
    setEventSubType(value);
  };

  const handleSubmit = () => {
    if (!file) {
      toast.warning("Please upload an image");
      return;
    }

    const data = {
      filename: file.name,
      // TODO: do not hardcode the eventType
      eventType: "image_analysis",
      contentType: file.type,
    };

    const handlePresignedSuccess = (presignedResponse) => {
      if (presignedResponse.status !== 200) {
        console.log(presignedResponse.data);
        toast.error("An error occurred. Please try again later.");
        return;
      }
      const presignedURL = presignedResponse.data?.url;
      if (!presignedURL || !presignedURL.length) {
        toast.error("An error occurred. Please try again later.");
        return;
      }

      uploadS3File({ file, presignedURL })
        .then((res) => handleUploadSuccess(res, presignedResponse))
        .catch(handleErrors);
    };

    const handleUploadSuccess = (response, presignedResponse) => {
      if (response.status !== 200) {
        console.log(response.data);
        toast.error("Could not upload the file. Please try again later.");
        return;
      }

      // TODO: do not hardcode values
      const eventData = {
        filename: presignedResponse.data?.filename,
        filepath: presignedResponse.data?.filepath,
        eventType: "image_analysis",
        subType: eventSubType,
        contentType: file.type,
        filesize: file.size,
      };

      createEvent(eventData).then(handleCreateEventSuccess).catch(handleErrors);
    };

    const handleCreateEventSuccess = (response) => {
      if (response.status !== 200) {
        console.log(response.data);
        toast.error("Could not create the event. Please try again later.");
        return;
      }
      toast.success("Event successfully created");
    };

    getPresignedUrl(data).then(handlePresignedSuccess).catch(handleErrors);

    // Implement API upload logic here
  };

  return (
    <>
      <div {...getRootProps()}>
        <input {...getInputProps()} />
        {file ? (
          <Image
            src={URL.createObjectURL(file)}
            alt="Uploaded Image"
            width={400}
            height={400}
            className="rounded-md"
          />
        ) : (
          <p className="flex items-center justify-center h-full min-h-[200px] rounded-md border bg-muted cursor-pointer hover:bg-gray-200 transition">
            {isDragActive
              ? "Drop the image here..."
              : "Drag & drop the image or click to upload"}
          </p>
        )}
      </div>
      <div className="mt-4">
        <AnalysisTypeSelect value={eventSubType} onChange={handleTypeChange} />
      </div>
      <Button className="my-4" onClick={handleSubmit}>
        Submit
      </Button>
    </>
  );
}

export function AnalysisTypeSelect({ value, onChange }) {
  const analysisSubTypes = [
    // TODO: do not hardcode the values
    { value: "detect_labels", label: "Label Detection" },
    { value: "detect_faces", label: "Face Detection" },
  ];

  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <Label htmlFor="analysis-type">Analysis Type</Label>
      <Select value={value} onValueChange={onChange}>
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Select analysis type" />
        </SelectTrigger>
        <SelectContent>
          {analysisSubTypes.map((type) => (
            <SelectItem key={type.value} value={type.value}>
              {type.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}

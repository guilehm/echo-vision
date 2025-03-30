"use client";

import { useState } from "react";
import { useDropzone } from "react-dropzone";
import Image from "next/image";
import { toast } from "sonner";
import { createEvent, getPresignedUrl } from "@/services/server-requester";
import { Button } from "@/components/ui/button";
import { uploadS3File } from "@/services/client-requester";

const handleErrors = (error) => {
  console.log(error);
  toast.error("An error occurred. Please try again later.");
};

export default function ImageUpload() {
  const [file, setFile] = useState(null);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: { "image/*": [] },
    multiple: false,
    onDrop: (acceptedFiles) => {
      setFile(acceptedFiles[0]);
    },
  });

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
    getPresignedUrl(data)
      .then((presignedResponse) => {
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
          .then((response) => {
            console.log("uploadResponse", response);
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
              subType: "detect_labels",
              contentType: file.type,
              filesize: file.size,
            };

            console.log("eventData", eventData);

            createEvent(eventData)
              .then((response) => {
                console.log("eventResponse", response);
                if (response.status !== 200) {
                  console.log(response.data);
                  toast.error(
                    "Could not create the event. Please try again later.",
                  );
                  return;
                }
                toast.success("Event successfully created");
              })
              .catch(handleErrors);
          })
          .catch(handleErrors);
      })
      .catch(handleErrors);

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
      <Button className="my-4" onClick={handleSubmit}>
        Submit
      </Button>
    </>
  );
}

export const formatDate = (dateString) => {
  const options = {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  };
  return new Date(dateString).toLocaleDateString("en-US", options);
};

export const statusStyles = {
  pending: {
    bg: "#fef9c3",
    text: "#92400e",
    border: "#fde68a",
  },
  processing: {
    bg: "#dbeafe",
    text: "#1d4ed8",
    border: "#93c5fd",
  },
  completed: {
    bg: "#dcfce7",
    text: "#166534",
    border: "#86efac",
  },
  failed: {
    bg: "#fee2e2",
    text: "#991b1b",
    border: "#fca5a5",
  },
  default: {
    bg: "#f4f4f5",
    text: "#3f3f46",
    border: "#d4d4d8",
  },
};

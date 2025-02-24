import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function FeatureCard({ icon, title, description }) {
  return (
    <Card className="text-center p-6">
      <CardHeader className="flex flex-col items-center space-y-2">
        {icon}
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-gray-600">{description}</p>
      </CardContent>
    </Card>
  );
}

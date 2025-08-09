import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function DashboardPage() {
    return (
        <div>
            <Card>
                <CardHeader>
                    <CardTitle>داشبورد</CardTitle>
                </CardHeader>
                <CardContent>
                    <p>به پنل مدیریت خوش آمدید!</p>
                </CardContent>
            </Card>
        </div>
    );
}
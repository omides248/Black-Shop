// file: frontend/app/admin/layout.tsx
import Link from "next/link";

export default function AdminLayout({
                                        children,
                                    }: {
    children: React.ReactNode;
}) {
    return (
        <div className="flex min-h-screen">
            {/* Sidebar */}
            <aside className="w-64 bg-gray-800 text-white p-6 space-y-6 flex-shrink-0">
                <h2 className="text-2xl font-bold mb-6 text-indigo-400">پنل مدیریت</h2>
                <nav>
                    <ul className="space-y-3">
                        <li>
                            <Link href="/admin" className="block p-3 rounded-md hover:bg-gray-700 transition-colors">
                                داشبورد
                            </Link>
                        </li>
                        <li>
                            <Link href="/admin/categories" className="block p-3 rounded-md hover:bg-gray-700 transition-colors">
                                مدیریت دسته‌بندی‌ها
                            </Link>
                        </li>
                        <li>
                            <Link href="/admin/products" className="block p-3 rounded-md hover:bg-gray-700 transition-colors">
                                مدیریت محصولات
                            </Link>
                        </li>
                        {/* لینک‌های دیگر */}
                    </ul>
                </nav>
            </aside>

            {/* Main Content */}
            <div className="flex-grow p-8 bg-gray-100 text-gray-800">
                {children}
            </div>
        </div>
    );
}
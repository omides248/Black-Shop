"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Package, Home, ShoppingCart, Users, Folder } from "lucide-react";

const NavLink = ({ href, children }: { href: string; children: React.ReactNode }) => {
    const pathname = usePathname();
    const isActive = pathname === href;

    return (
        <Link
            href={href}
            className={`flex items-center gap-3 rounded-lg px-3 py-2 text-gray-500 transition-all hover:text-gray-900 ${
                isActive ? "bg-gray-200 text-gray-900" : ""
            }`}
        >
            {children}
        </Link>
    );
};

export function Sidebar() {
    return (
        <div className="hidden border-l bg-gray-100/40 lg:block">
            <div className="flex h-full max-h-screen flex-col gap-2">
                <div className="flex h-[60px] items-center border-b px-6">
                    <Link href="/admin/dashboard" className="flex items-center gap-2 font-semibold">
                        <Package className="h-6 w-6" />
                        <span>پنل مدیریت</span>
                    </Link>
                </div>
                <div className="flex-1 overflow-auto py-2">
                    <nav className="grid items-start px-4 text-sm font-medium">
                        <NavLink href="/admin/dashboard">
                            <Home className="h-4 w-4" />
                            داشبورد
                        </NavLink>
                        <NavLink href="/admin/orders">
                            <ShoppingCart className="h-4 w-4" />
                            سفارشات
                        </NavLink>
                        <NavLink href="/admin/products">
                            <Package className="h-4 w-4" />
                            محصولات
                        </NavLink>
                        <NavLink href="/admin/categories">
                            <Folder className="h-4 w-4" />
                            دسته‌بندی‌ها
                        </NavLink>
                        <NavLink href="/admin/users">
                            <Users className="h-4 w-4" />
                            کاربران
                        </NavLink>
                    </nav>
                </div>
            </div>
        </div>
    );
}
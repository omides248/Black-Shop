"use client";

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";

export function Header() {
    return (
        <header className="flex h-14 items-center gap-4 border-b bg-gray-100/40 px-6 lg:h-[60px]">
            <Button
                variant="outline"
                size="icon"
                className="shrink-0 lg:hidden"
            >
                <Menu className="h-5 w-5" />
                <span className="sr-only">Toggle navigation menu</span>
            </Button>
            <div className="w-full flex-1">
                {/* می‌توانید جستجوی کلی را اینجا اضافه کنید */}
            </div>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="rounded-full border w-8 h-8">
                        <img
                            src="https://placehold.co/32x32/E2E8F0/4A5568?text=A"
                            className="rounded-full"
                            alt="Avatar"
                        />
                        <span className="sr-only">Toggle user menu</span>
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>حساب کاربری من</DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>تنظیمات</DropdownMenuItem>
                    <DropdownMenuItem>پشتیبانی</DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>خروج</DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </header>
    );
}
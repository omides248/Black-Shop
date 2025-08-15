// مسیر: lib/actions/category-actions.ts
"use server";

import { revalidatePath } from "next/cache";

// آدرس API بک‌اند شما
const API_URL = "http://192.168.8.140:8080/v1/categories";

/**
 * یک دسته‌بندی جدید در سرور ایجاد می‌کند.
 * @param formData داده‌های فرم شامل name و parentId.
 * @returns آبجکتی شامل وضعیت موفقیت و یا پیام خطا.
 */
export async function createCategory(formData: FormData) {
    const name = formData.get("name") as string;
    const parentId = formData.get("parentId") as string; // مقدار value از Select

    // ساختار بدنه درخواست دقیقاً مطابق با نیاز شما
    const body = {
        name: name,
        parent_id: parentId || null, // اگر parentId رشته خالی بود، null ارسال می‌شود
    };

    console.log(body);

    try {
        const response = await fetch(API_URL, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body),
        });

        if (!response.ok) {
            // خواندن پیام خطا از API بک‌اند
            const errorData = await response.json();
            throw new Error(errorData.error || "خطا در ایجاد دسته‌بندی");
        }

        // ✅ مهم: این خط باعث می‌شود لیست دسته‌بندی‌ها در UI به‌روز شود
        revalidatePath("/admin/categories");

        return { success: true };
    } catch (error: any) {
        console.error("Error creating category:", error);
        return { success: false, error: error.message };
    }
}

// توابع دیگر مانند getCategories, updateCategory, deleteCategory اینجا قرار می‌گیرند...
export interface Category {
    id: string;
    name: string;
    parentId?: string | null;
    imageUrl?: string | null;
    depth: number;
    createdAt: string;
    subcategory?: Category[];
}

export async function getCategories(): Promise<Category[]> {
    try {
        const res = await fetch(API_URL, { cache: "no-store" });
        if (!res.ok) throw new Error("Failed to fetch categories");
        const data = await res.json();
        return data.categories || [];
    } catch (error) {
        console.error("Error fetching categories:", error);
        return [];
    }
}
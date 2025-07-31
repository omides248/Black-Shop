// file: frontend/app/admin/categories/actions.tsx
"use server";

import { catalogAPI } from "@/lib/api/client";

interface Category {
    id: string;
    name: string;
    imageUrl?: string | null;
    parentId?: string | null;
    depth: number;
}

interface ListCategoriesResult {
    categories: Category[];
    error: string | null;
}

// <<-- تغییر در اینجا: تعریف رابط برای خطای ساختاریافته
interface APIError {
    message: string;
    code: number; // کد gRPC (مثلاً 9 برای FailedPrecondition, 6 برای AlreadyExists)
    details: any[];
    status: number; // کد وضعیت HTTP
}

interface CreateCategoryResult {
    category: Category | null;
    error: APIError | null; // خطا حالا از نوع APIError است
}

/**
 * Server Action for listing categories.
 * Fetches categories from the backend API.
 * @returns {ListCategoriesResult} An object containing categories array or an error message.
 */
export async function listCategoriesServerAction(): Promise<ListCategoriesResult> {
    try {
        const data = await catalogAPI.listCategories();
        return { categories: data.categories || [], error: null };
    } catch (error: any) {
        console.error("Server Action: Failed to fetch categories:", error);
        // اگر خطای ساختاریافته باشد، پیام آن را برمی‌گردانیم
        return { categories: [], error: error.message || "خطا در بارگذاری دسته‌بندی‌ها از سرور." };
    }
}

/**
 * Server Action for creating a new category.
 * Creates a category via the backend API.
 * @param {object} data - The category data (name, imageUrl, parentId).
 * @returns {CreateCategoryResult} An object containing the created category or an error message.
 */
export async function createCategoryServerAction(data: { name: string; imageUrl?: string | null; parentId?: string | null }): Promise<CreateCategoryResult> {
    try {
        const response = await catalogAPI.createCategory(data);
        return { category: response.category, error: null };
    } catch (error: any) {
        console.error("Server Action: Failed to create category:", error);
        // <<-- تغییر در اینجا: آبجکت خطای کامل را برمی‌گردانیم
        return { category: null, error: error as APIError }; // type assertion
    }
}
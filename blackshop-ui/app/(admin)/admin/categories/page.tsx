// مسیر: app/(admin)/admin/categories/page.tsx
import React from 'react';
import CategoryManager from '@/components/admin/categories/category-manager';
import { getCategories, Category } from '@/lib/actions/category-actions';

export default async function CategoryPage() {
    const categories: Category[] = await getCategories();

    return <CategoryManager initialCategories={categories} />;
}
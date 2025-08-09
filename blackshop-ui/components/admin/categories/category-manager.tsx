// مسیر: components/admin/categories/category-manager.tsx
"use client";

import React, { useState } from 'react';
import CategoryForm from './category-form';
import CategoryTree from './category-tree';
import { Category } from '@/lib/actions/category-actions';

interface CategoryManagerProps {
    initialCategories: Category[];
}

const flattenCategories = (categories: Category[]): Category[] => {
    const flatList: Category[] = [];

    function recurse(categoryList: Category[], depth: number) {
        for (const category of categoryList) {
            flatList.push({ ...category, name: `${'—'.repeat(depth)} ${category.name}` });
            if (category.subcategory) {
                recurse(category.subcategory, depth + 1);
            }
        }
    }

    recurse(categories, 0);
    return flatList;
};


export default function CategoryManager({ initialCategories }: CategoryManagerProps) {
    const [editingCategory, setEditingCategory] = useState<Category | null>(null);

    const flatCategoryList = flattenCategories(initialCategories);

    const handleEdit = (category: Category) => {
        setEditingCategory(category);
    };

    const handleDelete = (category: Category) => {
        // TODO: Implement delete confirmation
        console.log("Delete:", category);
    };

    const handleFinish = () => {
        setEditingCategory(null);
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1">
                <CategoryForm
                    categories={flatCategoryList}
                    editingCategory={editingCategory}
                    onFinish={handleFinish}
                />
            </div>
            <div className="lg:col-span-2">
                <CategoryTree
                    categories={initialCategories}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                />
            </div>
        </div>
    );
}
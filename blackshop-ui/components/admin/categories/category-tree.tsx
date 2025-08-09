"use client";

import React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Edit, Trash2, ChevronDown, ChevronRight } from "lucide-react";
import { Category } from "@/lib/actions/category-actions";

interface CategoryTreeProps {
    categories: Category[];
    onEdit: (cat: Category) => void;
    onDelete: (cat: Category) => void;
}

const CategoryTreeItem = ({
                              category,
                              onEdit,
                              onDelete,
                              level = 0,
                          }: {
    category: Category;
    onEdit: (cat: Category) => void;
    onDelete: (cat: Category) => void;
    level: number;
}) => {
    const [open, setOpen] = React.useState(true);

    const hasChildren = category.subcategory && category.subcategory.length > 0;

    return (
        <div className="relative">
            <div
                className="flex items-center justify-between rounded-md hover:bg-gray-100 relative"
                style={{ paddingRight: `${level * 20 + 8}px` }}
            >
                {/* خطوط راهنما */}
                {level > 0 && (
                    <div
                        className="absolute top-0 bottom-0 border-r border-gray-300"
                        style={{ right: `${level * 20 - 10}px` }}
                    />
                )}

                <div className="flex items-center gap-1 cursor-pointer"
                     onClick={() => hasChildren && setOpen(!open)}>
                    {hasChildren ? (
                        <span className="p-1 rounded hover:bg-gray-200 transition">
            {open ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
        </span>
                    ) : (
                        <span className="w-[18px]" /> // فاصله برای هم‌ترازی
                    )}

                    <span
                        className={`${
                            level === 0
                                ? "font-bold text-gray-900"
                                : level === 1
                                    ? "text-gray-700"
                                    : "text-gray-500"
                        }`}
                    >
        {category.name}
    </span>
                </div>


                <div className="flex items-center gap-2 pr-2">
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onEdit(category)}
                    >
                        <Edit className="h-4 w-4 text-blue-600" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onDelete(category)}
                    >
                        <Trash2 className="h-4 w-4 text-red-600" />
                    </Button>
                </div>
            </div>

            {/* نمایش بازگشتی زیرشاخه‌ها */}
            {hasChildren && open && (
                <div className="transition-all duration-200 ease-in-out">
                    {category.subcategory!.map((subCat) => (
                        <CategoryTreeItem
                            key={subCat.id}
                            category={subCat}
                            onEdit={onEdit}
                            onDelete={onDelete}
                            level={level + 1}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

const CategoryTree = ({ categories, onEdit, onDelete }: CategoryTreeProps) => {
    return (
        <Card className="shadow-lg">
            <CardHeader>
                <CardTitle className="text-xl font-bold text-slate-800">
                    لیست دسته‌بندی‌ها
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-1">
                    {categories.map((category) => (
                        <CategoryTreeItem
                            key={category.id}
                            category={category}
                            onEdit={onEdit}
                            onDelete={onDelete}
                            level={0}
                        />
                    ))}
                </div>
            </CardContent>
        </Card>
    );
};

export default CategoryTree;

// مسیر: components/admin/categories/category-form.tsx
"use client";

import React, { useRef, useState } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { createCategory, Category } from '@/lib/actions/category-actions';

interface CategoryFormProps {
    categories: Category[]; // این لیست باید صاف (flat) باشد
    editingCategory?: Category | null;
    onFinish: () => void;
}

export default function CategoryForm({ categories, editingCategory, onFinish }: CategoryFormProps) {
    const formRef = useRef<HTMLFormElement>(null);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [parentId, setParentId] = useState(editingCategory?.parentId || "");


    // برای ریست کردن فرم هنگام تغییر وضعیت ویرایش
    const formKey = editingCategory ? editingCategory.id : 'new-category';

    const handleFormAction = async (formData: FormData) => {
        // TODO: Handle edit logic here
        console.log("formData", formData);
        const result = await createCategory(formData);

        if (result.success) {
            setMessage({ type: 'success', text: 'دسته‌بندی با موفقیت ایجاد شد.' });
            formRef.current?.reset();
            onFinish();
        } else {
            setMessage({ type: 'error', text: `خطا: ${result.error}` });
        }
    };

    return (
        <Card className="shadow-lg">
            <CardHeader>
                <CardTitle className="text-xl font-bold text-slate-800">
                    {editingCategory ? 'ویرایش دسته‌بندی' : 'افزودن دسته‌بندی جدید'}
                </CardTitle>
            </CardHeader>
            <CardContent>
                <form
                    key={formKey}
                    ref={formRef}
                    action={handleFormAction}
                    className="space-y-6"
                >
                    <div>
                        <Label htmlFor="name" className="mb-2 block text-gray-500">نام دسته‌بندی</Label>
                        <Input
                            id="name"
                            name="name"
                            defaultValue={editingCategory?.name || ''}
                            required
                        />
                    </div>
                    <div>
                        <Label htmlFor="parentId" className="mb-2 block text-gray-500">دسته‌بندی والد</Label>
                        <Select
                            value={parentId}
                            onValueChange={(val) => setParentId(val)}
                        >
                            <SelectTrigger className="flex-row-reverse justify-between text-right">
                                <SelectValue placeholder="-- بدون والد (اصلی) --" />
                            </SelectTrigger>
                            <SelectContent className="select-rtl bg-white dark:bg-slate-900 shadow-lg border border-gray-200">
                                {categories.map(cat => (
                                    <SelectItem key={cat.id} value={cat.id} disabled={editingCategory?.id === cat.id}>
                                        {cat.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>

                        <input type="hidden" name="parentId" value={parentId} />
                    </div>

                    <div className="flex gap-2 pt-2">
                        <Button type="submit" className="w-full bg-slate-900 hover:bg-slate-800 text-white">
                            {editingCategory ? 'ذخیره تغییرات' : 'افزودن دسته‌بندی'}
                        </Button>
                        {editingCategory && (
                            <Button type="button" variant="outline" onClick={onFinish}>انصراف</Button>
                        )}
                    </div>

                    {message && (
                        <p className={`text-sm mt-2 ${message.type === 'error' ? 'text-red-600' : 'text-green-600'}`}>
                            {message.text}
                        </p>
                    )}
                </form>
            </CardContent>
        </Card>
    );
}
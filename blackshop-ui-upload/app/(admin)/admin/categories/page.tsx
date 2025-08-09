/*
================================================================================
|                                                                              |
|   فایل شماره ۱: صفحه اصلی دسته‌بندی‌ها (نسخه نهایی)                           |
|   مسیر: app/(admin)/admin/categories/page.tsx                                 |
|                                                                              |
================================================================================
*/
"use client";

import React, { useState, useEffect, useId } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { UploadCloud, X, Edit, Trash2 } from "lucide-react";

// داده‌های نمونه
const initialCategories = [
    { id: '1', name: 'کالای دیجیتال', parentId: null, depth: 0, image: null as File | null },
    { id: '2', name: 'موبایل', parentId: '1', depth: 1, image: null as File | null },
    { id: '3', name: 'گوشی هوشمند', parentId: '2', depth: 2, image: null as File | null },
    { id: '4', name: 'لپتاپ', parentId: '1', depth: 1, image: null as File | null },
    { id: '5', name: 'پوشاک', parentId: null, depth: 0, image: null as File | null },
    { id: '6', name: 'مردانه', parentId: '5', depth: 1, image: null as File | null },
];

// کامپوننت آپلود عکس
const ImageUploader = ({ onFileSelect, existingFile, onRemove }: { onFileSelect: (file: File | null) => void, existingFile: File | null, onRemove: () => void }) => {
    const uniqueId = useId();
    const [isDragging, setIsDragging] = useState(false);

    const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        setIsDragging(true);
    };

    const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        setIsDragging(false);
    };

    const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        setIsDragging(false);
        if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
            onFileSelect(e.dataTransfer.files[0]);
            e.dataTransfer.clearData();
        }
    };

    return (
        <div>
            <div
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                className={`mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-dashed rounded-lg transition-colors ${isDragging ? 'border-indigo-600 bg-indigo-50' : ''}`}
            >
                <div className="space-y-1 text-center">
                    <UploadCloud className="mx-auto h-12 w-12 text-gray-400" />
                    <div className="flex text-sm text-gray-600">
                        <label htmlFor={uniqueId} className="relative cursor-pointer bg-white rounded-md font-medium text-indigo-600 hover:text-indigo-500">
                            <span>آپلود فایل</span>
                            <input id={uniqueId} type="file" className="sr-only" onChange={(e) => e.target.files && onFileSelect(e.target.files[0])} />
                        </label>
                    </div>
                </div>
            </div>
            {existingFile && (
                <div className="mt-4 relative w-24 h-24">
                    <img src={URL.createObjectURL(existingFile)} alt="Preview" className="w-full h-full object-cover rounded-md" />
                    <button type="button" onClick={onRemove} className="absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 leading-none"><X className="w-3 h-3" /></button>
                </div>
            )}
        </div>
    );
};

// کامپوننت برای نمایش یک ردیف دسته‌بندی در درخت
const CategoryTreeItem = ({ category, onEdit, onDelete, level = 0 }: { category: any, onEdit: (cat: any) => void, onDelete: (cat: any) => void, level: number }) => (
    <div className="flex items-center justify-between p-2 rounded-md hover:bg-gray-100">
        <p style={{ marginRight: `${level * 20}px` }}>{category.name}</p>
        <div className="flex items-center gap-2">
            <Button variant="ghost" size="icon" onClick={() => onEdit(category)}><Edit className="h-4 w-4 text-blue-600" /></Button>
            <Button variant="ghost" size="icon" onClick={() => onDelete(category)}><Trash2 className="h-4 w-4 text-red-600" /></Button>
        </div>
    </div>
);

// کامپوننت برای ساخت و نمایش درخت دسته‌بندی‌ها
const CategoryTree = ({ categories, onEdit, onDelete }: { categories: any[], onEdit: (cat: any) => void, onDelete: (cat: any) => void }) => {
    const buildTree = (parentId: string | null = null) => {
        return categories
            .filter(cat => cat.parentId === parentId)
            .map(cat => (
                <div key={cat.id}>
                    <CategoryTreeItem category={cat} onEdit={onEdit} onDelete={onDelete} level={cat.depth} />
                    {buildTree(cat.id)}
                </div>
            ));
    };
    return <div className="space-y-1">{buildTree()}</div>;
};

export default function CategoryPage() {
    const [categories, setCategories] = useState(initialCategories);
    const [editingCategory, setEditingCategory] = useState<any>(null);
    const [formData, setFormData] = useState({ name: '', parentId: '', image: null as File | null });

    useEffect(() => {
        if (editingCategory) {
            setFormData({
                name: editingCategory.name,
                parentId: editingCategory.parentId || '',
                image: editingCategory.image
            });
        } else {
            setFormData({ name: '', parentId: '', image: null });
        }
    }, [editingCategory]);

    const handleFormSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (editingCategory) {
            // منطق ویرایش
            setCategories(categories.map(c => c.id === editingCategory.id ? { ...c, name: formData.name, parentId: formData.parentId || null, image: formData.image } : c));
        } else {
            // منطق افزودن
            const newCategory = {
                id: Date.now().toString(),
                name: formData.name,
                parentId: formData.parentId || null,
                depth: formData.parentId ? categories.find(c => c.id === formData.parentId)!.depth + 1 : 0,
                image: formData.image
            };
            setCategories([...categories, newCategory]);
        }
        setEditingCategory(null);
    };

    const handleCancelEdit = () => {
        setEditingCategory(null);
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* ستون فرم افزودن/ویرایش */}
            <div className="lg:col-span-1">
                <Card className="shadow-lg">
                    <CardHeader>
                        <CardTitle className="text-xl font-bold text-slate-800">{editingCategory ? 'ویرایش دسته‌بندی' : 'افزودن دسته‌بندی جدید'}</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleFormSubmit} className="space-y-6">
                            <div>
                                <Label htmlFor="name" className="mb-2 block text-gray-500">نام دسته‌بندی</Label>
                                <Input id="name" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} required />
                            </div>
                            <div>
                                <Label htmlFor="parentId" className="mb-2 block text-gray-500">دسته‌بندی والد</Label>
                                <Select value={formData.parentId} onValueChange={value => setFormData({...formData, parentId: value})}>
                                    <SelectTrigger>
                                        <SelectValue placeholder="-- بدون والد (اصلی) --" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {/* [FIX]: حذف آیتم با مقدار خالی برای جلوگیری از خطا */}
                                        {categories.map(cat => (
                                            <SelectItem key={cat.id} value={cat.id}>{cat.name}</SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div>
                                <Label className="mb-2 block text-gray-500">تصویر دسته‌بندی</Label>
                                <ImageUploader
                                    onFileSelect={file => setFormData({...formData, image: file})}
                                    existingFile={formData.image}
                                    onRemove={() => setFormData({...formData, image: null})}
                                />
                            </div>
                            <div className="flex gap-2 pt-2">
                                <Button type="submit" className="w-full bg-slate-900 hover:bg-slate-800 text-white">{editingCategory ? 'ذخیره تغییرات' : 'افزودن دسته‌بندی'}</Button>
                                {editingCategory && (
                                    <Button type="button" variant="outline" onClick={handleCancelEdit}>انصراف</Button>
                                )}
                            </div>
                        </form>
                    </CardContent>
                </Card>
            </div>

            {/* ستون لیست دسته‌بندی‌ها */}
            <div className="lg:col-span-2">
                <Card className="shadow-lg">
                    <CardHeader>
                        <CardTitle className="text-xl font-bold text-slate-800">لیست دسته‌بندی‌ها</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <CategoryTree
                            categories={categories}
                            onEdit={setEditingCategory}
                            onDelete={(cat) => console.log("Delete:", cat)}
                        />
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}

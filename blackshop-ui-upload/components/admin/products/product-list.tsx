"use client";

import React, { useState, useEffect, useRef } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog";
import { PlusCircle, MoreHorizontal } from "lucide-react";

// داده‌های نمونه
const allProducts = [
    { id: 1, name: 'گوشی هوشمند مدل Pro X', sku: 'GP-X1-PRO', category: 'کالای دیجیتال', price: '32,500,000', stock: 120, status: 'منتشر شده' },
    { id: 2, name: 'پیراهن مردانه نخی', sku: 'SH-M-COT-L', category: 'پوشاک', price: '890,000', stock: 45, status: 'پیش‌نویس' },
    { id: 3, name: 'تلویزیون هوشمند 55 اینچ', sku: 'TV-S55-4K', category: 'لوازم خانگی', price: '25,000,000', stock: 80, status: 'منتشر شده' },
    { id: 4, name: 'لپتاپ گیمینگ Legion', sku: 'LP-LEG-5', category: 'کالای دیجیتال', price: '78,000,000', stock: 0, status: 'منتشر شده' },
];

// کامپوننت منوی عملیات
const ActionsMenu = ({ product, onDeleteClick }: { product: any, onDeleteClick: (product: any) => void }) => {
    const [isOpen, setIsOpen] = useState(false);
    const menuRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (menuRef.current && !menuRef.current.contains(event.target as Node)) setIsOpen(false);
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);
    return (
        <div className="relative" ref={menuRef}>
            <Button variant="ghost" size="icon" onClick={() => setIsOpen(!isOpen)}><MoreHorizontal className="h-4 w-4" /></Button>
            {isOpen && (
                <div className="absolute left-0 mt-2 w-48 bg-white rounded-md shadow-lg border z-10">
                    <a href="#" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">ویرایش</a>
                    <button onClick={() => { onDeleteClick(product); setIsOpen(false); }} className="block w-full text-right px-4 py-2 text-sm text-red-600 hover:bg-gray-100">حذف</button>
                </div>
            )}
        </div>
    );
};

export function ProductListPage({ setView }: { setView: (view: 'list' | 'add') => void }) {
    const [isDeleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<any>(null);
    const [filters, setFilters] = useState({ search: '', category: 'all', status: 'all' });

    const handleDeleteClick = (product: any) => {
        setSelectedProduct(product);
        setDeleteDialogOpen(true);
    };
    const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFilters(prev => ({ ...prev, [name]: value }));
    };

    const handleSelectFilterChange = (name: 'category' | 'status') => (value: string) => {
        setFilters(prev => ({ ...prev, [name]: value }));
    };

    const filteredProducts = allProducts.filter(product => {
        const searchMatch = product.name.toLowerCase().includes(filters.search.toLowerCase());
        const categoryMatch = filters.category === 'all' || product.category === filters.category;
        const statusMatch = filters.status === 'all' || product.status === filters.status;
        return searchMatch && categoryMatch && statusMatch;
    });

    return (
        <>
            <Card className="w-full">
                <CardHeader className="flex flex-col sm:flex-row items-center justify-between gap-4">
                    <CardTitle>مدیریت محصولات</CardTitle>
                    <Button onClick={() => setView('add')}><PlusCircle className="ml-2 h-4 w-4" /> افزودن محصول</Button>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                        <Input name="search" placeholder="جستجو در محصولات..." value={filters.search} onChange={handleFilterChange} />
                        <Select value={filters.category} onValueChange={handleSelectFilterChange('category')}>
                            <SelectTrigger>
                                <SelectValue placeholder="دسته‌بندی" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">همه دسته‌بندی‌ها</SelectItem>
                                <SelectItem value="کالای دیجیتال">کالای دیجیتال</SelectItem>
                                <SelectItem value="پوشاک">پوشاک</SelectItem>
                                <SelectItem value="لوازم خانگی">لوازم خانگی</SelectItem>
                            </SelectContent>
                        </Select>
                        <Select value={filters.status} onValueChange={handleSelectFilterChange('status')}>
                            <SelectTrigger>
                                <SelectValue placeholder="وضعیت" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">همه وضعیت‌ها</SelectItem>
                                <SelectItem value="منتشر شده">منتشر شده</SelectItem>
                                <SelectItem value="پیش‌نویس">پیش‌نویس</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead className="hidden sm:table-cell">تصویر</TableHead>
                                <TableHead>نام</TableHead>
                                <TableHead>وضعیت</TableHead>
                                <TableHead className="hidden md:table-cell">قیمت</TableHead>
                                <TableHead className="hidden md:table-cell">موجودی</TableHead>
                                <TableHead><span className="sr-only">عملیات</span></TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {filteredProducts.map((product) => (
                                <TableRow key={product.id}>
                                    <TableCell className="hidden sm:table-cell"><img alt="Product" className="aspect-square rounded-md object-cover" height="64" src={`https://placehold.co/64x64/E2E8F0/4A5568?text=P${product.id}`} width="64" /></TableCell>
                                    <TableCell className="font-medium">{product.name}</TableCell>
                                    <TableCell><Badge variant={product.status === 'منتشر شده' ? 'default' : 'secondary'}>{product.status}</Badge></TableCell>
                                    <TableCell className="hidden md:table-cell">{product.price} تومان</TableCell>
                                    <TableCell className="hidden md:table-cell">{product.stock > 0 ? `${product.stock} عدد` : <span className="text-red-500">ناموجود</span>}</TableCell>
                                    <TableCell><ActionsMenu product={product} onDeleteClick={handleDeleteClick} /></TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
            <AlertDialog open={isDeleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>آیا کاملا مطمئن هستید؟</AlertDialogTitle>
                        <AlertDialogDescription>این عملیات قابل بازگشت نیست. محصول "{selectedProduct?.name}" برای همیشه حذف خواهد شد.</AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel onClick={() => setDeleteDialogOpen(false)}>انصراف</AlertDialogCancel>
                        <AlertDialogAction onClick={() => { console.log(`Deleting product: ${selectedProduct?.name}`); setDeleteDialogOpen(false); }}>بله، حذف کن</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    );
};

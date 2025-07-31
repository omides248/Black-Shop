// file: frontend/app/admin/categories/page.tsx
"use client";

import { useEffect, useState, useActionState } from "react";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { catalogAPI } from "@/lib/api/client";
import { listCategoriesServerAction, createCategoryServerAction } from './actions';

interface Category {
    id: string;
    name: string;
    imageUrl?: string | null;
    parentId?: string | null;
    depth: number;
}

interface FormState {
    success: boolean;
    message: string | null;
    errors?: { name?: string; parentId?: string; general?: string };
}

const initialFormState: FormState = {
    success: false,
    message: null,
};

async function createCategoryAction(previousState: FormState, formData: FormData): Promise<FormState> {
    return { success: false, message: "این Server Action فقط یک Placeholder است." };
}

const sortCategoriesForDisplay = (categories: Category[]): Category[] => {
    const categoryMap = new Map<string | null, Category[]>();
    categories.forEach(cat => {
        const parentKey = cat.parentId || null;
        if (!categoryMap.has(parentKey)) {
            categoryMap.set(parentKey, []);
        }
        categoryMap.get(parentKey)?.push(cat);
    });

    categoryMap.forEach(list => {
        list.sort((a, b) => a.name.localeCompare(b.name));
    });

    const sortedList: Category[] = [];

    const addChildren = (parentId: string | null) => {
        const children = categoryMap.get(parentId);
        if (children) {
            for (const child of children) {
                sortedList.push(child);
                addChildren(child.id);
            }
        }
    };

    addChildren(null);

    return sortedList;
};


export default function AdminCategoriesPage() {
    const [categories, setCategories] = useState<Category[]>([]);
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [newCategoryName, setNewCategoryName] = useState("");
    const [newCategoryImageUrl, setNewCategoryImageUrl] = useState<string | null>(null);
    const [newCategoryParentId, setNewCategoryParentId] = useState<string | null>(null);
    const [formState, formAction] = useActionState(createCategoryAction, initialFormState);

    const [localError, setLocalError] = useState<string | null>(null);

    const parentCategoriesForSelection = categories.filter(cat => cat.depth === 0 || cat.depth === 1);


    const fetchCategories = async () => {
        try {
            const result = await listCategoriesServerAction();
            if (result.error) {
                setLocalError(result.error);
                setCategories([]);
            } else {
                const sortedCategories = sortCategoriesForDisplay(result.categories);
                setCategories(sortedCategories);
                setLocalError(null);
            }
        } catch (error: any) {
            console.error("خطا در فراخوانی Server Action برای دریافت دسته‌بندی‌ها:", error);
            setLocalError(error.message || "خطا در بارگذاری دسته‌بندی‌ها از سرور.");
        }
    };

    useEffect(() => {
        fetchCategories();
    }, []);


    const handleCreateCategory = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLocalError(null);

        if (!newCategoryName.trim()) {
            setLocalError("نام دسته‌بندی نمی‌تواند خالی باشد.");
            return;
        }

        try {
            const result = await createCategoryServerAction({
                name: newCategoryName,
                imageUrl: newCategoryImageUrl,
                parentId: newCategoryParentId,
            });

            if (result.error) {
                // <<-- تغییر در اینجا: هندل کردن خطا بر اساس کد gRPC
                // کدهای gRPC:
                // codes.AlreadyExists (6)
                // codes.NotFound (5)
                // codes.FailedPrecondition (9)
                // codes.Internal (13)

                switch (result.error.code) {
                    case 6: // AlreadyExists
                        setLocalError("دسته‌بندی با این نام در این سطح تکراری است.");
                        break;
                    case 5: // NotFound (برای والد)
                        setLocalError("دسته‌بندی والد یافت نشد.");
                        break;
                    case 9: // FailedPrecondition (برای محدودیت عمق یا وجود محصول)
                        // در اینجا باید پیام دقیق‌تر را از details یا message خود خطا بگیریم
                        if (result.error.message.includes("depth limit exceeded")) {
                            setLocalError("عمق دسته‌بندی بیش از حد مجاز است (حداکثر ۲ سطح).");
                        } else if (result.error.message.includes("contains products")) {
                            setLocalError("نمی‌توانید به دسته‌بندی‌ای که محصول دارد، زیردسته اضافه کنید.");
                        } else {
                            setLocalError(result.error.message || "شرط اولیه برآورده نشد.");
                        }
                        break;
                    case 13: // Internal (خطای داخلی سرور)
                        setLocalError("خطای داخلی سرور در ایجاد دسته‌بندی.");
                        break;
                    default:
                        setLocalError(result.error.message || "خطا در ایجاد دسته‌بندی.");
                }
            } else if (result.category) {
                fetchCategories();
                setIsDialogOpen(false);
                setNewCategoryName("");
                setNewCategoryImageUrl(null);
                setNewCategoryParentId(null);
                setLocalError(null);
            } else {
                setLocalError("پاسخ نامعتبر از سرور.");
            }
        } catch (error: any) {
            console.error("خطا در فراخوانی Server Action برای ایجاد دسته‌بندی:", error);
            setLocalError(error.message || "خطا در ایجاد دسته‌بندی.");
        }
    };

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">مدیریت دسته‌بندی‌ها</h1>
                <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                    <DialogTrigger asChild>
                        <Button>ایجاد دسته‌بندی جدید</Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-[425px]">
                        <DialogHeader>
                            <DialogTitle>ایجاد دسته‌بندی جدید</DialogTitle>
                            <DialogDescription>
                                نام و جزئیات دسته‌بندی جدید را وارد کنید.
                            </DialogDescription>
                        </DialogHeader>
                        <form onSubmit={handleCreateCategory} className="grid gap-4 py-4">
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="name" className="text-right">
                                    نام
                                </Label>
                                <Input
                                    id="name"
                                    value={newCategoryName}
                                    onChange={(e) => setNewCategoryName(e.target.value)}
                                    className="col-span-3"
                                    required
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="imageUrl" className="text-right">
                                    تصویر URL (اختیاری)
                                </Label>
                                <Input
                                    id="imageUrl"
                                    value={newCategoryImageUrl || ""}
                                    onChange={(e) => setNewCategoryImageUrl(e.target.value || null)}
                                    className="col-span-3"
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="parentId" className="text-right">
                                    دسته‌بندی والد (اختیاری)
                                </Label>
                                <Select
                                    onValueChange={(value) => setNewCategoryParentId(value === "none-parent" ? null : value)}
                                    value={newCategoryParentId || "none-parent"}
                                >
                                    <SelectTrigger className="col-span-3">
                                        <SelectValue placeholder="انتخاب والد" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="none-parent">بدون والد (ریشه)</SelectItem>
                                        {parentCategoriesForSelection.map(cat => (
                                            <SelectItem key={cat.id} value={cat.id}>
                                                {cat.name} (عمق: {cat.depth})
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            {localError && (
                                <p className="text-red-500 text-sm col-span-full text-center">
                                    {localError}
                                </p>
                            )}

                            <DialogFooter>
                                <Button type="submit">ایجاد</Button>
                            </DialogFooter>
                        </form>
                    </DialogContent>
                </Dialog>
            </div>

            {categories.length === 0 ? (
                <p className="text-gray-600">دسته‌بندی وجود ندارد.</p>
            ) : (
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="text-right">نام</TableHead>
                            <TableHead className="text-right">تصویر</TableHead>
                            <TableHead className="text-right">عملیات</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {categories.map((category) => (
                            <TableRow key={category.id}>
                                <TableCell className="text-right">
                                    <div dir="rtl" style={{ paddingRight: `${category.depth * 40}px` }} className="flex items-center justify-start">
                                        {category.depth > 0 && (
                                            <span className="ml-2 text-gray-400">←</span>
                                        )}
                                        {category.name}
                                    </div>
                                </TableCell>
                                <TableCell className="text-right">
                                    <div className="flex justify-start">
                                        {category.imageUrl ? (
                                            <img src={category.imageUrl} alt={category.name} className="h-8 w-8 object-cover rounded-full" />
                                        ) : "—"}
                                    </div>
                                </TableCell>
                                <TableCell className="text-right">
                                    <Button variant="outline" size="sm">ویرایش</Button>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            )}
        </div>
    );
}
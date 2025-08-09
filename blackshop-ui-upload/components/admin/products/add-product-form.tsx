"use client";

import React, {useId, useState} from 'react';
import {Button} from "@/components/ui/button";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {Input} from "@/components/ui/input";
import {Label} from "@/components/ui/label";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {Separator} from "@/components/ui/separator";
import {Textarea} from "@/components/ui/textarea";
import {ArrowRight, PlusCircle, Trash2, UploadCloud, X} from "lucide-react";

// =================================================================================
// کامپوننت‌های کمکی (بدون تغییر)
// =================================================================================

const ProgressBar = ({currentStep}: { currentStep: number }) => {
    const steps = ["اطلاعات اصلی", "انواع محصول", "بازبینی"];
    return (
        <div className="flex items-center justify-center mb-10">
            {steps.map((step, index) => {
                const stepNumber = index + 1;
                let statusClass = 'bg-gray-200 text-gray-600';
                let textClass = 'text-gray-500';
                if (stepNumber < currentStep) {
                    statusClass = 'bg-indigo-600 text-white';
                    textClass = 'text-indigo-700';
                } else if (stepNumber === currentStep) {
                    statusClass = 'bg-indigo-600 text-white';
                    textClass = 'text-indigo-700 font-bold';
                }
                return (
                    <React.Fragment key={step}>
                        <div className="flex items-center">
                            <div
                                className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-lg transition-colors duration-300 ${statusClass}`}>
                                {stepNumber}
                            </div>
                            <p className={`mr-3 font-semibold transition-colors duration-300 ${textClass}`}>{step}</p>
                        </div>
                        {index < steps.length - 1 && <div className="flex-auto border-t-2 border-gray-300 mx-4"></div>}
                    </React.Fragment>
                );
            })}
        </div>
    );
};

const ImageUploader = ({onFilesSelect, existingFiles = [], onRemove}: {
    onFilesSelect: (files: File[]) => void,
    existingFiles: File[],
    onRemove: (index: number) => void
}) => {
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
            onFilesSelect(Array.from(e.dataTransfer.files));
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
                    <UploadCloud className="mx-auto h-12 w-12 text-gray-400"/>
                    <div className="flex text-sm text-gray-600">
                        <label htmlFor={uniqueId}
                               className="relative cursor-pointer bg-white rounded-md font-medium text-indigo-600 hover:text-indigo-500">
                            <span>آپلود فایل</span>
                            <input id={uniqueId} type="file" className="sr-only" multiple
                                   onChange={(e) => e.target.files && onFilesSelect(Array.from(e.target.files))}/>
                        </label>
                        <p className="pr-1">یا بکشید و رها کنید</p>
                    </div>
                </div>
            </div>
            {existingFiles.length > 0 && (<div
                className="mt-4 grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-4">{existingFiles.map((file, index) => (
                <div key={index} className="relative"><img src={URL.createObjectURL(file)} alt={`Preview ${index}`}
                                                           className="w-full h-24 object-cover rounded-md"/>
                    <button type="button" onClick={() => onRemove(index)}
                            className="absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 leading-none"><X
                        className="w-3 h-3"/></button>
                </div>))}</div>)}
        </div>
    );
};


// =================================================================================
// کامپوننت‌های مراحل (منتقل شده به بیرون)
// =================================================================================

// Props for Step1
type Step1Props = {
    formData: { name: string; description: string; category: string; brand: string; primaryImage: File[] };
    handleChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
    handleSelectChange: (name: 'category' | 'brand') => (value: string) => void;
    handleFileChange: (files: File[]) => void;
    removeImage: (index: number) => void;
};

const Step1 = ({formData, handleChange, handleSelectChange, handleFileChange, removeImage}: Step1Props) => (
    <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
                <Label htmlFor="name" className="mb-2 block text-gray-700">نام محصول</Label>
                <Input id="name" name="name" value={formData.name} onChange={handleChange}/>
            </div>
            <div>
                <Label className="mb-2 block text-gray-700">دسته‌بندی</Label>
                <Select value={formData.category} onValueChange={handleSelectChange('category')}>
                    <SelectTrigger className="w-full flex-row-reverse justify-between text-right">
                        <SelectValue placeholder="انتخاب دسته‌بندی"/>
                    </SelectTrigger>
                    <SelectContent className="select-rtl bg-white dark:bg-slate-900 shadow-lg border border-gray-200">
                        <SelectItem value="digital">کالای دیجیتال</SelectItem>
                        <SelectItem value="home">خانه و آشپزخانه</SelectItem>
                    </SelectContent>
                </Select>
            </div>
            <div>
                <Label className="mb-2 block text-gray-700">برند</Label>
                <Select value={formData.brand} onValueChange={handleSelectChange('brand')}>
                    <SelectTrigger className="w-full flex-row-reverse justify-between text-right">
                        <SelectValue placeholder="انتخاب برند"/>
                    </SelectTrigger>
                    <SelectContent className="select-rtl bg-white dark:bg-slate-900 shadow-lg border border-gray-200">
                        <SelectItem value="apple">اپل</SelectItem>
                        <SelectItem value="samsung">سامسونگ</SelectItem>
                        <SelectItem value="xiaomi">شیائومی</SelectItem>
                    </SelectContent>
                </Select>
            </div>
        </div>
        <div>
            <Label htmlFor="description" className="mb-2 block text-gray-700">توضیحات</Label>
            <Textarea id="description" name="description" value={formData.description} onChange={handleChange}/>
        </div>
        <div>
            <Label className="mb-2 block text-gray-700">تصویر اصلی محصول</Label>
            <ImageUploader
                onFilesSelect={handleFileChange}
                existingFiles={formData.primaryImage}
                onRemove={removeImage}
            />
        </div>
    </div>
);

// Props for Step2
type Step2Props = {
    variants: {
        sku: string;
        price: string;
        stock: string;
        attributes: { name: string; value: string }[];
        images: File[]
    }[];
    handleChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>, index: number, subIndex?: number) => void;
    handleFileChange: (files: File[], variantIndex: number) => void;
    removeImage: (index: number, variantIndex: number) => void;
    addVariant: () => void;
    removeVariant: (index: number) => void;
    addAttribute: (variantIndex: number) => void;
    removeAttribute: (variantIndex: number, attrIndex: number) => void;
};

const Step2 = ({
                   variants,
                   handleChange,
                   handleFileChange,
                   removeImage,
                   addVariant,
                   removeVariant,
                   addAttribute,
                   removeAttribute
               }: Step2Props) => (
    <div className="space-y-6">
        {variants.map((variant, index) => (
            <div key={index} className="border rounded-xl p-4 bg-gray-50/80">
                <div className="flex flex-row items-center justify-between mb-4">
                    <h4 className="text-lg font-semibold text-gray-700">نوع محصول #{index + 1}</h4>
                    <Button type="button" variant="destructive" size="icon" onClick={() => removeVariant(index)}>
                        <Trash2 className="h-4 w-4"/>
                    </Button>
                </div>
                <div className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div>
                            <Label htmlFor={`sku-${index}`} className="text-gray-500">SKU</Label>
                            <Input id={`sku-${index}`} name="sku" value={variant.sku}
                                   onChange={(e) => handleChange(e, index)}/>
                        </div>
                        <div>
                            <Label htmlFor={`price-${index}`} className="text-gray-500">قیمت</Label>
                            <Input id={`price-${index}`} name="price" type="number" value={variant.price}
                                   onChange={(e) => handleChange(e, index)}/>
                        </div>
                        <div>
                            <Label htmlFor={`stock-${index}`} className="text-gray-500">موجودی</Label>
                            <Input id={`stock-${index}`} name="stock" type="number" value={variant.stock}
                                   onChange={(e) => handleChange(e, index)}/>
                        </div>
                    </div>
                    <Separator/>
                    <div>
                        <Label className="mb-2 block text-gray-500">ویژگی‌ها</Label>
                        {variant.attributes.map((attr, attrIndex) => (
                            <div key={attrIndex} className="flex items-center gap-2 mt-2">
                                <Input name="name" value={attr.name}
                                       onChange={(e) => handleChange(e, index, attrIndex)}/>
                                <Input name="value" value={attr.value}
                                       onChange={(e) => handleChange(e, index, attrIndex)}/>
                                <Button type="button" variant="outline" size="icon"
                                        onClick={() => removeAttribute(index, attrIndex)}>
                                    <Trash2 className="h-4 w-4 text-red-500"/>
                                </Button>
                            </div>
                        ))}
                        <Button type="button" variant="outline" size="sm" className="mt-2"
                                onClick={() => addAttribute(index)}>
                            افزودن ویژگی
                        </Button>
                    </div>
                    <Separator/>
                    <div>
                        <Label className="mb-1 block text-gray-500">تصاویر این نوع</Label>
                        <ImageUploader
                            onFilesSelect={(files) => handleFileChange(files, index)}
                            existingFiles={variant.images}
                            onRemove={(imgIndex) => removeImage(imgIndex, index)}
                        />
                    </div>
                </div>
            </div>
        ))}
        <Button type="button" className="w-full border-dashed" variant="outline" onClick={addVariant}>
            <PlusCircle className="ml-2 h-4 w-4"/> افزودن نوع جدید
        </Button>
    </div>
);

// Props for Step3
type Step3Props = {
    formData: {
        name: string;
        primaryImage: File[];
        variants: { sku: string; images: File[] }[];
    };
};

const Step3 = ({formData}: Step3Props) => (
    <div className="space-y-6 bg-gray-50 p-6 rounded-lg">
        <div>
            <h4 className="font-semibold text-lg text-gray-800 border-b pb-2 mb-2">اطلاعات اصلی</h4>
            <p><strong>نام:</strong> {formData.name || '-'}</p>
            <div className="mt-2">
                <strong>تصویر اصلی:</strong>
                {formData.primaryImage.length > 0 ? (
                    <img src={URL.createObjectURL(formData.primaryImage[0])}
                         className="h-20 w-20 rounded-md object-cover mt-1" alt="preview"/>
                ) : (
                    <p className="text-sm text-gray-500">انتخاب نشده</p>
                )}
            </div>
        </div>
        <div>
            <h4 className="font-semibold text-lg text-gray-800 border-b pb-2 mb-2">انواع محصول</h4>
            {formData.variants.map((variant, index) => (
                <div key={index} className="mt-2 p-2 border-t">
                    <p><strong>نوع #{index + 1}:</strong></p>
                    <p className="text-sm">SKU: {variant.sku}</p>
                    <div className="mt-2">
                        <strong>تصاویر نوع:</strong>
                        <div className="flex gap-2 mt-1">
                            {variant.images.length > 0 ? (
                                variant.images.map((img, i) => (
                                    <img key={i} src={URL.createObjectURL(img)}
                                         className="h-16 w-16 rounded-md object-cover" alt="preview"/>
                                ))
                            ) : (
                                <p className="text-sm text-gray-500">انتخاب نشده</p>
                            )}
                        </div>
                    </div>
                </div>
            ))}
        </div>
    </div>
);


// =================================================================================
// کامپوننت اصلی
// =================================================================================

export function AddProductPage({setView}: { setView: (view: 'list' | 'add') => void }) {
    const [currentStep, setCurrentStep] = useState(1);
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        category: '',
        brand: '',
        primaryImage: [] as File[],
        variants: [{sku: '', price: '', stock: '', attributes: [{name: 'رنگ', value: ''}], images: [] as File[]}]
    });

    // --- Handlers ---
    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>, index?: number, subIndex?: number) => {
        const {name, value} = e.target;
        if (index !== undefined) {
            const newVariants = [...formData.variants];
            if (subIndex !== undefined) {
                newVariants[index].attributes[subIndex][name as 'name' | 'value'] = value;
            } else {
                (newVariants[index] as any)[name] = value;
            }
            setFormData({...formData, variants: newVariants});
        } else {
            setFormData({...formData, [name]: value});
        }
    };

    const handleSelectChange = (name: 'category' | 'brand') => (value: string) => {
        setFormData(prev => ({...prev, [name]: value}));
    };

    const handleFileChange = (files: File[], variantIndex?: number) => {
        if (variantIndex !== undefined) {
            const newVariants = [...formData.variants];
            newVariants[variantIndex].images = [...newVariants[variantIndex].images, ...files];
            setFormData({...formData, variants: newVariants});
        } else {
            setFormData({...formData, primaryImage: [...formData.primaryImage, ...files]});
        }
    };

    const removeImage = (index: number, variantIndex?: number) => {
        if (variantIndex !== undefined) {
            const newVariants = [...formData.variants];
            newVariants[variantIndex].images = newVariants[variantIndex].images.filter((_, i) => i !== index);
            setFormData({...formData, variants: newVariants});
        } else {
            setFormData({...formData, primaryImage: formData.primaryImage.filter((_, i) => i !== index)});
        }
    };

    const addVariant = () => setFormData({
        ...formData,
        variants: [...formData.variants, {
            sku: '',
            price: '',
            stock: '',
            attributes: [{name: 'رنگ', value: ''}],
            images: []
        }]
    });

    const removeVariant = (index: number) => setFormData({
        ...formData,
        variants: formData.variants.filter((_, i) => i !== index)
    });

    const addAttribute = (variantIndex: number) => {
        const newVariants = [...formData.variants];
        newVariants[variantIndex].attributes.push({name: '', value: ''});
        setFormData({...formData, variants: newVariants});
    };

    const removeAttribute = (variantIndex: number, attrIndex: number) => {
        const newVariants = [...formData.variants];
        newVariants[variantIndex].attributes = newVariants[variantIndex].attributes.filter((_, i) => i !== attrIndex);
        setFormData({...formData, variants: newVariants});
    };

    const nextStep = () => setCurrentStep(prev => prev < 3 ? prev + 1 : prev);
    const prevStep = () => setCurrentStep(prev => prev > 1 ? prev - 1 : prev);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        console.log("Final Data:", formData);
        alert("محصول ثبت شد!");
        setView('list');
    };

    return (
        <Card className="w-full max-w-4xl mx-auto shadow-lg">
            <CardHeader>
                <div className="flex justify-between items-center">
                    <div>
                        <CardTitle className="text-2xl font-bold text-slate-800">افزودن محصول جدید</CardTitle>
                        <CardDescription className="mt-2 text-gray-500">اطلاعات محصول را در سه مرحله ساده وارد
                            کنید.</CardDescription>
                    </div>
                    <Button variant="outline" onClick={() => setView('list')}>
                        <ArrowRight className="ml-2 h-4 w-4"/>
                        بازگشت به لیست
                    </Button>
                </div>
            </CardHeader>
            <CardContent>
                <ProgressBar currentStep={currentStep}/>
                <form onSubmit={handleSubmit}>
                    <div className="mt-4">
                        {currentStep === 1 && <Step1
                            formData={formData}
                            handleChange={handleChange}
                            handleSelectChange={handleSelectChange}
                            handleFileChange={(files) => handleFileChange(files)}
                            removeImage={(index) => removeImage(index)}
                        />}
                        {currentStep === 2 && <Step2
                            variants={formData.variants}
                            handleChange={handleChange}
                            handleFileChange={handleFileChange}
                            removeImage={removeImage}
                            addVariant={addVariant}
                            removeVariant={removeVariant}
                            addAttribute={addAttribute}
                            removeAttribute={removeAttribute}
                        />}
                        {currentStep === 3 && <Step3 formData={formData}/>}
                    </div>
                    <div className={`mt-10 flex ${currentStep === 1 ? 'justify-end' : 'justify-between'}`}>
                        {currentStep > 1 && (
                            <Button type="button" variant="outline" onClick={prevStep}>&rarr; مرحله قبل</Button>)}
                        {currentStep < 3 && (
                            <Button type="button" className="bg-slate-900 hover:bg-slate-800 text-white"
                                    onClick={nextStep}>مرحله بعد &larr;</Button>)}
                        {currentStep === 3 && (
                            <Button type="submit" className="bg-green-600 hover:bg-green-700">ثبت نهایی محصول</Button>)}
                    </div>
                </form>
            </CardContent>
        </Card>
    );
}

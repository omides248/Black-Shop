// مسیر: components/admin/categories/image-uploader.tsx
"use client";

import React, { useState, useId } from 'react';
import { UploadCloud, X } from "lucide-react";

// (کد این کامپوننت را از فایل page.tsx اصلی کپی کنید)
// ...
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

export default ImageUploader;
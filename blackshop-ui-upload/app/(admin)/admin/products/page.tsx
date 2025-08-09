"use client";

import {useState} from "react";
import {ProductListPage} from "@/components/admin/products/product-list";
import {AddProductPage} from "@/components/admin/products/add-product-form";

export default function ProductsPage() {
    const [view, setView] = useState<'list' | 'add'>('list');

    if (view === 'add') {
        return <AddProductPage setView={setView}/>;
    }

    return <ProductListPage setView={setView}/>;
}
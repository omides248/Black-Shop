const API_PROTOCOL = process.env.NEXT_PUBLIC_API_PROTOCOL || 'http';

function getBaseUrl(service: 'catalog' | 'identity'): string {
    if (service === 'catalog') {
        const baseUrl = process.env.NEXT_PUBLIC_CATALOG_API_BASE_URL;
        if (!baseUrl) throw new Error("Catalog API URL is not configured.");
        return `${API_PROTOCOL}://${baseUrl}/v1`;
    }
    // service === 'identity'
    const baseUrl = process.env.NEXT_PUBLIC_IDENTITY_API_BASE_URL;
    if (!baseUrl) throw new Error("Identity API URL is not configured.");
    return `${API_PROTOCOL}://${baseUrl}/v1`;
}


async function apiFetch(url: string, options: RequestInit = {}) {
    try {
        const res = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            cache: 'no-store',
        });

        if (!res.ok) {
            const errorData = await res.json();
            throw new Error(errorData.message || 'An API error occurred');
        }

        const text = await res.text();
        return text ? JSON.parse(text) : null;

    } catch (error) {
        console.error(`API fetch error for ${url}:`, error);
        throw error;
    }
}

export const catalogAPI = {
    getProduct: (id: string) => apiFetch(`${getBaseUrl('catalog')}/products/${id}`),
    getProducts: () => apiFetch(`${getBaseUrl('catalog')}/products`),
};

export const identityAPI = {
    register: (data: any) => apiFetch(`${getBaseUrl('identity')}/auth/register`, {
        method: 'POST',
        body: JSON.stringify(data),
    }),
    login: (data: any) => apiFetch(`${getBaseUrl('identity')}/auth/login`, {
        method: 'POST',
        body: JSON.stringify(data),
    }),
    getProfile: (token: string) => apiFetch(`${getBaseUrl('identity')}/users/me`, {
        headers: {'Authorization': `Bearer ${token}`},
    }),
};
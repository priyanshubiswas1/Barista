import axios from "axios"

const api = axios.create({
    baseURL: "http://localhost:8000/api",
    headers: {
        'Content-Type': 'application/json',
    },
    credentials: 'include'
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

export const registerUser = async (credentials) => {
    const response = await api.post(`/auth/onboard`, credentials);
    return response.status === 200 ? response.data : null
}

export const loginUser = async (credentials) => {
    const response = await api.post(`/auth/login`, credentials);
    return response.status === 200 ? response.data : null
}

export const verifyUser = async () => {
    const response = await api.get(`/auth/verify`);
    return response.status === 200 ? true : false
}

export const fetchTablesByShopId = async (shopId) => {
    const response = await api.get(`/tables/${shopId}`);
    return response.status === 200 ? response.data : null
};

export const createTableByShopId = async ({ shopId, tableNumber }) => {
    const response = await api.post(`/table`, {
        shopId: shopId,
        tableNumber: tableNumber
    });
    return response.status === 200 ? response.data : null
};

export const fetchCategoryItemsByShopId = async (id) => {
    const response = await api.get(`/category-items/${id}`);
    return response.status === 200 ? response.data : null
};

export const fetchCategoriesByShopId = async (id) => {
    const response = await api.get(`/category/${id}`);
    return response.status === 200 ? response.data : null
};

export const deleteCategoryByCategoryId = async ({ categoryId, shopId }) => {
    const response = await api.delete(`/category`, {
        data: { categoryId, shopId }
    });
    return response.status === 200 ? response.data : null
};

export const createCategory = async ({ name, shopId }) => {
    const response = await api.post(`/category`, { name: name, shopId: shopId });
    return response.status === 200 ? response.data : null
}

export const fetchItemsByShopId = async (id) => {
    const response = await api.get(`/items/${id}`);
    return response.status === 200 ? response.data : null
};

export const addItem = async ({ name, shopId, categoryId, price }) => {
    const response = await api.post(`/item`, { name: name, shopId: shopId, description: "test", categoryId: categoryId, price: price, image: "url", available: true });
    return response.status === 200 ? response.data : null
}

export const fetchActiveSessions = async (shopId) => {
    const response = await api.get(`/active-sessions?shopid=${shopId}`);
    return response.status === 200 ? response.data : null
};

export const fetchOrdersBySessionId = async (sessionId) => {
    const response = await api.get(`/orders/${sessionId}`);
    return response.status === 200 ? response.data : null
};

export const createOrder = async (payload) => {
    const response = await api.post(`/order-admin`, payload);
    return response.status === 200 ? response.data : null
};

export const getSessionSummary = async (payload) => {
    const response = await api.get(`/history?shopid=${payload.shopId}&offset=${payload.offset}&limit=${payload.limit}`);
    return response.status === 200 ? response.data : null
};
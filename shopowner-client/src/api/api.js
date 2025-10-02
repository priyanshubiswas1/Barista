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

export const fetchItemsByShopId = async (id) => {
    const response = await api.get(`/items/${id}`);
    return response.status === 200 ? response.data : null
};

export const fetchCategoriesByShopId = async (id) => {
    const response = await api.get(`/category-items/${id}`);
    return response.status === 200 ? response.data : null
};

export const deleteCategoryByCategoryId = async ({ categoryId, shopId }) => {
    const response = await api.delete(`/category`, {
        data: { categoryId, shopId }
    });
    return response.status === 200 ? response.data : null
};
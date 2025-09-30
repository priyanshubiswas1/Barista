import axios from "axios"

const api = axios.create({
    baseURL: "http://localhost:8000/api",
    headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkzNDcxNTUsImlhdCI6MTc1OTI2MDc1NSwianRpIjoiYWEzMzYxMzEtNWYxYi00NTFiLThiODItYzlmZWFiOTA0NDA2In0.WTvuBXuuDxWE9kzTK-C6ueJKcIgQKUOVzl6Y363cIBs',
    },
    credentials: 'include'
});

export const fetchItemsByShopId = async (id) => {
    const response = await api.get(`/items/${id}`);
    return response.status === 200 ? response.data : null
};
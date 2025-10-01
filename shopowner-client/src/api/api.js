import axios from "axios"

const api = axios.create({
    baseURL: "http://localhost:8000/api",
    headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkzOTA4NjEsImlhdCI6MTc1OTMwNDQ2MSwianRpIjoiMTlhYmFhOGYtOGQ0ZC00ZjEyLTg5N2QtMjQ1Y2IwYTA4YTZhIn0.q2h7jXd5jRMkl8kcQGgjus3ryI8VT-MD6igRtqIJ2ek',
    },
    credentials: 'include'
});

export const fetchItemsByShopId = async (id) => {
    const response = await api.get(`/items/${id}`);
    return response.status === 200 ? response.data : null
};
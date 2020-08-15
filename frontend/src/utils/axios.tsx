import axios from 'axios';
import Cookies from 'universal-cookie';


const apiURL = "http://localhost:8080";
const cookies = new Cookies();


export const axiosInstance = axios.create({
    baseURL: apiURL,
    timeout: 3000,
})

export const axiosAuth = axios.create({
    baseURL: apiURL,
    headers: {
        Authorization: `Bearer ${cookies.get("token")}`
    }
})
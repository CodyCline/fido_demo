import React, { useState } from 'react';
import axios from 'axios';

const apiUrl = 'http://localhost:8080';

export const axiosInstance = axios.create({
    baseURL: apiUrl,
    timeout: 3000,
})
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 2, // Количество виртуальных пользователей (подберите под вашу инфраструктуру)
    duration: '3600s', // Длительность теста
    rps: 20, // Ограничение по RPS
};

const BASE_URL = 'http://localhost:1337'; // Замените на реальный адрес exchange_service

export default function () {
    // Случайно выбираем эндпоинт
    const endpoints = [
        '/exchange/rate',    // замените на ваш реальный эндпоинт
    ];
    const endpoint = endpoints[0];

    const payload = JSON.stringify({
        base_currency: 'USD',
        target_currency: 'RUB',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRzdEBnZ21haWwuY29tIiwiZXhwIjoxNzU1OTM2MzE1LCJpZCI6M30.CNAA5M4DzsTfy6UtsXvcV5CMBkvWcm3gACbvW9uns7U'
        },
    };

    let res = http.post(`${BASE_URL}${endpoint}`, payload, params);
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

}
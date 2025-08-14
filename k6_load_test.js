import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 100, // Количество виртуальных пользователей (подберите под вашу инфраструктуру)
    duration: '3600s', // Длительность теста
    rps: 200, // Ограничение по RPS
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
        },
    };

    let res = http.post(`${BASE_URL}${endpoint}`, payload, params);
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

}
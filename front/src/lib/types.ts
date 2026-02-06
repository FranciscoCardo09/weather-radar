export interface Cities {
    id: string;
    name: string;
    latitude: number;
    longitude: number;
}

export interface WeatherData {
    city_id: string;
    city_name: string;
    temperature: number;
    humidity: number;
    wind_speed: number;
    condition: string;
}